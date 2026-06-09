package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	ai "github.com/HiIamJeff67/shift-hero-backend/app/ai"
	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	inputs "github.com/HiIamJeff67/shift-hero-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	validation "github.com/HiIamJeff67/shift-hero-backend/app/validation"
)

const (
	maxScheduleInsightRange = 31 * 24 * time.Hour
	scheduleInsightTimeout  = 90 * time.Second
)

type SchedulingInsightServiceInterface interface {
	GenerateScheduleInsights(
		ctx context.Context,
		reqDto *dtos.GenerateScheduleInsightsReqDto,
		stream func(context.Context, []byte) error,
	) (*dtos.ScheduleInsightResDto, *exceptions.Exception)
}

type SchedulingInsightService struct {
	db                    *gorm.DB
	userAccountRepository repositories.UserAccountRepositoryInterface
	generator             ai.ScheduleInsightGeneratorInterface
	generatorInit         error
}

func NewSchedulingInsightService(
	db *gorm.DB,
	userAccountRepository repositories.UserAccountRepositoryInterface,
	generator ai.ScheduleInsightGeneratorInterface,
	generatorInit error,
) SchedulingInsightServiceInterface {
	if db == nil {
		db = models.DB
	}
	if userAccountRepository == nil {
		userAccountRepository = repositories.NewUserAccountRepository()
	}
	return &SchedulingInsightService{
		db:                    db,
		userAccountRepository: userAccountRepository,
		generator:             generator,
		generatorInit:         generatorInit,
	}
}

func (s *SchedulingInsightService) GenerateScheduleInsights(
	ctx context.Context,
	reqDto *dtos.GenerateScheduleInsightsReqDto,
	stream func(context.Context, []byte) error,
) (*dtos.ScheduleInsightResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	startAt := reqDto.Body.StartAt.UTC()
	endAt := reqDto.Body.EndAt.UTC()
	if exception := validateTimeRange(startAt, endAt); exception != nil {
		return nil, exception
	}
	if endAt.Sub(startAt) > maxScheduleInsightRange {
		return nil, exceptions.Scheduling.InsightRangeTooLarge()
	}
	if s.generator == nil {
		exceptions.Scheduling.AIUnavailable().WithOrigin(s.generatorInit).Log()
		return nil, exceptions.Scheduling.AIUnavailable()
	}

	analysisCtx, cancel := context.WithTimeout(ctx, scheduleInsightTimeout)
	defer cancel()

	db := s.db.WithContext(analysisCtx)
	if _, exception := requireCompanyManager(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	snapshot, metrics, companyName, timezone, exception := s.buildScheduleInsightSnapshot(
		db,
		reqDto.Param.CompanyId,
		startAt,
		endAt,
	)
	if exception != nil {
		return nil, exception
	}

	locale := reqDto.Body.Locale
	if locale == "" {
		locale = "zh-TW"
	}
	focus := reqDto.Body.Focus
	if focus == "" {
		focus = "Provide a balanced operational overview and prioritize the most actionable risks."
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		exceptions.Scheduling.AIGenerationFailed().WithOrigin(err).Log()
		return nil, exceptions.Scheduling.AIGenerationFailed()
	}

	aiUsage, exception := s.reserveAIUsage(analysisCtx, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	summary, err := s.generator.Generate(analysisCtx, ai.ScheduleInsightPromptInput{
		Locale:       locale,
		Focus:        focus,
		SnapshotJSON: string(snapshotJSON),
	}, stream)
	if err != nil {
		s.releaseAIUsageReservation(reqDto.ContextFields.UserId, aiUsage)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			exceptions.Scheduling.AIUnavailable().WithOrigin(err).Log()
			return nil, exceptions.Scheduling.AIUnavailable()
		}
		exceptions.Scheduling.AIGenerationFailed().WithOrigin(err).Log()
		return nil, exceptions.Scheduling.AIGenerationFailed()
	}

	return &dtos.ScheduleInsightResDto{
		CompanyId:   reqDto.Param.CompanyId,
		CompanyName: companyName,
		StartAt:     startAt,
		EndAt:       endAt,
		Timezone:    timezone,
		Locale:      locale,
		Model:       s.generator.Model(),
		Workflow: []string{
			"deterministic_schedule_analyzer",
			"workforce_auditor_agent",
			"executive_narrator_agent",
		},
		Metrics: metrics,
		AIUsage: dtos.ScheduleInsightAIUsageResDto{
			Used:      aiUsage.Used,
			Limit:     aiUsage.Limit,
			Remaining: max(aiUsage.Limit-aiUsage.Used, 0),
			ResetAt:   aiUsage.ResetAt,
		},
		Summary:     summary,
		GeneratedAt: time.Now().UTC(),
	}, nil
}

type aiUsageReservation struct {
	Used        int32
	Limit       int32
	PeriodStart time.Time
	ResetAt     time.Time
}

type MonthlyAIUsageDecision struct {
	Allowed     bool
	Used        int32
	Limit       int32
	PeriodStart time.Time
	ResetAt     time.Time
}

func EvaluateMonthlyAIUsage(
	used int32,
	limit int32,
	storedPeriodStart time.Time,
	now time.Time,
) MonthlyAIUsageDecision {
	periodStart := startOfUTCMonth(now)
	if !startOfUTCMonth(storedPeriodStart).Equal(periodStart) {
		used = 0
	}
	decision := MonthlyAIUsageDecision{
		Allowed:     used < limit,
		Used:        used,
		Limit:       limit,
		PeriodStart: periodStart,
		ResetAt:     periodStart.AddDate(0, 1, 0),
	}
	if decision.Allowed {
		decision.Used++
	}
	return decision
}

func (s *SchedulingInsightService) reserveAIUsage(
	ctx context.Context,
	userId uuid.UUID,
) (*aiUsageReservation, *exceptions.Exception) {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(tx.Error)
	}
	transactionClosed := false
	defer func() {
		if !transactionClosed {
			tx.Rollback()
		}
	}()

	quota, exception := s.userAccountRepository.GetAIUsageQuotaByUserIdForUpdate(
		userId,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		return nil, exception
	}

	decision := EvaluateMonthlyAIUsage(
		quota.MonthlyUsageCount,
		quota.MonthlyLimit,
		quota.PeriodStart,
		time.Now().UTC(),
	)
	if !decision.Allowed {
		return nil, exceptions.Scheduling.AIUsageLimitExceeded(
			decision.Used,
			decision.Limit,
			decision.ResetAt,
		)
	}

	if exception := s.userAccountRepository.UpdateAIUsageByUserId(
		userId,
		inputs.UpdateUserAIUsageInput{
			AIMonthlyUsageCount: decision.Used,
			AIUsagePeriodStart:  decision.PeriodStart,
		},
		options.WithTransactionDB(tx),
	); exception != nil {
		return nil, exception
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(err)
	}
	transactionClosed = true

	return &aiUsageReservation{
		Used:        decision.Used,
		Limit:       decision.Limit,
		PeriodStart: decision.PeriodStart,
		ResetAt:     decision.ResetAt,
	}, nil
}

func (s *SchedulingInsightService) releaseAIUsageReservation(
	userId uuid.UUID,
	reservation *aiUsageReservation,
) {
	if reservation == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if exception := s.userAccountRepository.ReleaseAIUsageReservationByUserId(
		userId,
		reservation.PeriodStart,
		options.WithDB(s.db.WithContext(ctx)),
	); exception != nil {
		exception.Log()
	}
}

func startOfUTCMonth(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}

type scheduleInsightSnapshot struct {
	CompanyName string                       `json:"companyName"`
	StartAt     time.Time                    `json:"startAt"`
	EndAt       time.Time                    `json:"endAt"`
	Timezone    string                       `json:"timezone"`
	Policy      scheduleInsightPolicy        `json:"policy"`
	Metrics     scheduleInsightPromptMetrics `json:"metrics"`
}

type scheduleInsightPolicy struct {
	MaxWeeklyHours int32 `json:"maxWeeklyHours"`
	MinRestHours   int32 `json:"minRestHours"`
}

type scheduleInsightPromptMetrics struct {
	RequiredHeadcount     int64                           `json:"requiredHeadcount"`
	AssignedHeadcount     int64                           `json:"assignedHeadcount"`
	UnfilledHeadcount     int64                           `json:"unfilledHeadcount"`
	CoverageRate          float64                         `json:"coverageRate"`
	OpenSwapRequestCount  int64                           `json:"openSwapRequestCount"`
	AverageHours          float64                         `json:"averageHours"`
	WorkloadSpreadHours   float64                         `json:"workloadSpreadHours"`
	EmployeesAtRisk       int                             `json:"employeesAtRisk"`
	AvailabilityConflicts int                             `json:"availabilityConflicts"`
	Employees             []scheduleInsightPromptEmployee `json:"employees"`
}

type scheduleInsightPromptEmployee struct {
	DisplayName            string             `json:"displayName"`
	EmployeeRole           enums.EmployeeRole `json:"employeeRole"`
	ShiftCount             int                `json:"shiftCount"`
	TotalHours             float64            `json:"totalHours"`
	LongestShiftHours      float64            `json:"longestShiftHours"`
	NightShiftCount        int                `json:"nightShiftCount"`
	WeekendShiftCount      int                `json:"weekendShiftCount"`
	ShortRestCount         int                `json:"shortRestCount"`
	AvailabilityConflicts  int                `json:"availabilityConflicts"`
	OpenSwapRequestCount   int                `json:"openSwapRequestCount"`
	MaxConsecutiveWorkDays int                `json:"maxConsecutiveWorkDays"`
	OvertimeWeekCount      int                `json:"overtimeWeekCount"`
	RiskScore              int                `json:"riskScore"`
	RiskLevel              string             `json:"riskLevel"`
}

type ScheduleInsightMember struct {
	UserId       uuid.UUID          `gorm:"column:user_id"`
	DisplayName  string             `gorm:"column:display_name"`
	EmployeeRole enums.EmployeeRole `gorm:"column:employee_role"`
}

func (s *SchedulingInsightService) buildScheduleInsightSnapshot(
	db *gorm.DB,
	companyId uuid.UUID,
	startAt time.Time,
	endAt time.Time,
) (*scheduleInsightSnapshot, dtos.ScheduleInsightMetricsResDto, string, string, *exceptions.Exception) {
	company := schemas.Company{}
	if err := db.Model(&schemas.Company{}).Where("id = ?", companyId).First(&company).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	settings := schemas.CompanySettings{
		CompanyId:      companyId,
		MaxWeeklyHours: 40,
		MinRestHours:   8,
		Timezone:       "Asia/Taipei",
	}
	if err := db.Model(&schemas.CompanySettings{}).Where("company_id = ?", companyId).First(&settings).Error; err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	location, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		location = time.UTC
		settings.Timezone = "UTC"
	}

	members := []ScheduleInsightMember{}
	memberTable := schemas.UsersToCompanies{}.TableName()
	userTable := schemas.User{}.TableName()
	if err := db.Model(&schemas.UsersToCompanies{}).
		Select(fmt.Sprintf(
			"%q.user_id, users.display_name, %q.employee_role",
			memberTable,
			memberTable,
		)).
		Joins(fmt.Sprintf("JOIN %q AS users ON users.id = %q.user_id", userTable, memberTable)).
		Where(fmt.Sprintf("%q.company_id = ?", memberTable), companyId).
		Order("users.display_name ASC").
		Scan(&members).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	requirements := []schemas.ShiftRequirement{}
	if err := db.Model(&schemas.ShiftRequirement{}).
		Where("company_id = ? AND start_at < ? AND end_at > ?", companyId, endAt, startAt).
		Order("start_at ASC").
		Find(&requirements).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	assignments := []schemas.ShiftAssignment{}
	if err := db.Model(&schemas.ShiftAssignment{}).
		Where("company_id = ? AND start_at < ? AND end_at > ?", companyId, endAt, startAt).
		Order("user_id ASC, start_at ASC").
		Find(&assignments).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	availabilitySlots := []schemas.AvailabilitySlot{}
	if err := db.Model(&schemas.AvailabilitySlot{}).
		Where("company_id = ? AND start_at < ? AND end_at > ?", companyId, endAt, startAt).
		Order("user_id ASC, start_at ASC").
		Find(&availabilitySlots).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	activeSwaps := []schemas.SwapRequest{}
	swapTable := schemas.SwapRequest{}.TableName()
	assignmentTable := schemas.ShiftAssignment{}.TableName()
	if err := db.Model(&schemas.SwapRequest{}).
		Select(fmt.Sprintf("%q.*", swapTable)).
		Joins(fmt.Sprintf(
			"JOIN %q AS assignment ON assignment.id = %q.shift_assignment_id",
			assignmentTable,
			swapTable,
		)).
		Where(fmt.Sprintf("%q.company_id = ? AND %q.status IN ?", swapTable, swapTable), companyId, []enums.SwapRequestStatus{
			enums.SwapRequestStatus_Open,
			enums.SwapRequestStatus_Claimed,
		}).
		Where("assignment.start_at < ? AND assignment.end_at > ?", endAt, startAt).
		Find(&activeSwaps).Error; err != nil {
		return nil, dtos.ScheduleInsightMetricsResDto{}, "", "", exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	metrics := CalculateScheduleInsightMetrics(
		members,
		requirements,
		assignments,
		availabilitySlots,
		activeSwaps,
		settings,
		location,
		startAt,
		endAt,
	)
	snapshot := &scheduleInsightSnapshot{
		CompanyName: company.Name,
		StartAt:     startAt,
		EndAt:       endAt,
		Timezone:    settings.Timezone,
		Policy: scheduleInsightPolicy{
			MaxWeeklyHours: settings.MaxWeeklyHours,
			MinRestHours:   settings.MinRestHours,
		},
		Metrics: scheduleMetricsForPrompt(metrics),
	}
	return snapshot, metrics, company.Name, settings.Timezone, nil
}

func CalculateScheduleInsightMetrics(
	members []ScheduleInsightMember,
	requirements []schemas.ShiftRequirement,
	assignments []schemas.ShiftAssignment,
	availabilitySlots []schemas.AvailabilitySlot,
	activeSwaps []schemas.SwapRequest,
	settings schemas.CompanySettings,
	location *time.Location,
	startAt time.Time,
	endAt time.Time,
) dtos.ScheduleInsightMetricsResDto {
	assignmentsByUser := make(map[uuid.UUID][]schemas.ShiftAssignment)
	assignedByRequirement := make(map[uuid.UUID]int64)
	for _, assignment := range assignments {
		assignmentsByUser[assignment.UserId] = append(assignmentsByUser[assignment.UserId], assignment)
		assignedByRequirement[assignment.ShiftRequirementId]++
	}

	availabilityByUser := make(map[uuid.UUID][]schemas.AvailabilitySlot)
	for _, slot := range availabilitySlots {
		if slot.IsAvailable {
			availabilityByUser[slot.UserId] = append(availabilityByUser[slot.UserId], slot)
		}
	}
	openSwapsByUser := make(map[uuid.UUID]int)
	for _, swap := range activeSwaps {
		openSwapsByUser[swap.RequesterUserId]++
	}

	var requiredHeadcount int64
	var assignedHeadcount int64
	for _, requirement := range requirements {
		required := int64(requirement.RequiredCount)
		assigned := assignedByRequirement[requirement.Id]
		requiredHeadcount += required
		assignedHeadcount += min(assigned, required)
	}

	employeeMetrics := make([]dtos.ScheduleInsightEmployeeResDto, 0, len(members))
	totalHours := 0.0
	minHours := math.MaxFloat64
	maxHours := 0.0
	totalAvailabilityConflicts := 0
	employeesAtRisk := 0
	for _, member := range members {
		metric := calculateEmployeeScheduleInsight(
			member,
			assignmentsByUser[member.UserId],
			availabilityByUser[member.UserId],
			openSwapsByUser[member.UserId],
			settings,
			location,
			startAt,
			endAt,
		)
		employeeMetrics = append(employeeMetrics, metric)
		totalHours += metric.TotalHours
		minHours = math.Min(minHours, metric.TotalHours)
		maxHours = math.Max(maxHours, metric.TotalHours)
		totalAvailabilityConflicts += metric.AvailabilityConflicts
		if metric.RiskLevel == "high" || metric.RiskLevel == "critical" {
			employeesAtRisk++
		}
	}
	if len(employeeMetrics) == 0 {
		minHours = 0
	}

	coverageRate := 1.0
	if requiredHeadcount > 0 {
		coverageRate = float64(assignedHeadcount) / float64(requiredHeadcount)
	}
	averageHours := 0.0
	if len(employeeMetrics) > 0 {
		averageHours = totalHours / float64(len(employeeMetrics))
	}

	sort.SliceStable(employeeMetrics, func(i int, j int) bool {
		if employeeMetrics[i].RiskScore == employeeMetrics[j].RiskScore {
			return employeeMetrics[i].TotalHours > employeeMetrics[j].TotalHours
		}
		return employeeMetrics[i].RiskScore > employeeMetrics[j].RiskScore
	})

	return dtos.ScheduleInsightMetricsResDto{
		RequiredHeadcount:     requiredHeadcount,
		AssignedHeadcount:     assignedHeadcount,
		UnfilledHeadcount:     max(requiredHeadcount-assignedHeadcount, 0),
		CoverageRate:          round(coverageRate, 3),
		OpenSwapRequestCount:  int64(len(activeSwaps)),
		AverageHours:          round(averageHours, 1),
		WorkloadSpreadHours:   round(maxHours-minHours, 1),
		EmployeesAtRisk:       employeesAtRisk,
		AvailabilityConflicts: totalAvailabilityConflicts,
		Employees:             employeeMetrics,
	}
}

func calculateEmployeeScheduleInsight(
	member ScheduleInsightMember,
	assignments []schemas.ShiftAssignment,
	availability []schemas.AvailabilitySlot,
	openSwapCount int,
	settings schemas.CompanySettings,
	location *time.Location,
	rangeStart time.Time,
	rangeEnd time.Time,
) dtos.ScheduleInsightEmployeeResDto {
	sort.Slice(assignments, func(i int, j int) bool {
		return assignments[i].StartAt.Before(assignments[j].StartAt)
	})

	totalHours := 0.0
	longestShift := 0.0
	nightShifts := 0
	weekendShifts := 0
	shortRests := 0
	availabilityConflicts := 0
	workDates := make(map[string]time.Time)
	hoursByWeek := make(map[string]float64)

	for index, assignment := range assignments {
		startAt := latestTime(assignment.StartAt, rangeStart)
		endAt := earliestTime(assignment.EndAt, rangeEnd)
		hours := math.Max(endAt.Sub(startAt).Hours(), 0)
		totalHours += hours
		longestShift = math.Max(longestShift, hours)

		localStart := assignment.StartAt.In(location)
		localEnd := assignment.EndAt.In(location)
		if localStart.Hour() < 6 || localStart.Hour() >= 22 || localEnd.Hour() >= 22 || localEnd.Hour() < 6 {
			nightShifts++
		}
		if localStart.Weekday() == time.Saturday || localStart.Weekday() == time.Sunday {
			weekendShifts++
		}
		dateKey := localStart.Format("2006-01-02")
		workDates[dateKey] = time.Date(localStart.Year(), localStart.Month(), localStart.Day(), 0, 0, 0, 0, location)
		year, week := localStart.ISOWeek()
		hoursByWeek[fmt.Sprintf("%04d-W%02d", year, week)] += hours

		if index > 0 {
			restHours := assignment.StartAt.Sub(assignments[index-1].EndAt).Hours()
			if restHours < float64(settings.MinRestHours) {
				shortRests++
			}
		}
		if !assignmentCoveredByAvailability(assignment, availability) {
			availabilityConflicts++
		}
	}

	overtimeWeeks := 0
	for _, weeklyHours := range hoursByWeek {
		if weeklyHours > float64(settings.MaxWeeklyHours) {
			overtimeWeeks++
		}
	}
	maxConsecutiveDays := calculateMaxConsecutiveDays(workDates)

	riskScore := overtimeWeeks*3 +
		shortRests*2 +
		availabilityConflicts*2 +
		openSwapCount*2 +
		nightShifts/2
	if maxConsecutiveDays >= 6 {
		riskScore += 3
	} else if maxConsecutiveDays >= 5 {
		riskScore++
	}

	return dtos.ScheduleInsightEmployeeResDto{
		UserId:                 member.UserId,
		DisplayName:            member.DisplayName,
		EmployeeRole:           member.EmployeeRole,
		ShiftCount:             len(assignments),
		TotalHours:             round(totalHours, 1),
		LongestShiftHours:      round(longestShift, 1),
		NightShiftCount:        nightShifts,
		WeekendShiftCount:      weekendShifts,
		ShortRestCount:         shortRests,
		AvailabilityConflicts:  availabilityConflicts,
		OpenSwapRequestCount:   openSwapCount,
		MaxConsecutiveWorkDays: maxConsecutiveDays,
		OvertimeWeekCount:      overtimeWeeks,
		RiskScore:              riskScore,
		RiskLevel:              scheduleRiskLevel(riskScore),
	}
}

func assignmentCoveredByAvailability(
	assignment schemas.ShiftAssignment,
	availability []schemas.AvailabilitySlot,
) bool {
	for _, slot := range availability {
		if !slot.StartAt.After(assignment.StartAt) && !slot.EndAt.Before(assignment.EndAt) {
			return true
		}
	}
	return false
}

func scheduleMetricsForPrompt(metrics dtos.ScheduleInsightMetricsResDto) scheduleInsightPromptMetrics {
	employees := make([]scheduleInsightPromptEmployee, len(metrics.Employees))
	for i, employee := range metrics.Employees {
		employees[i] = scheduleInsightPromptEmployee{
			DisplayName:            employee.DisplayName,
			EmployeeRole:           employee.EmployeeRole,
			ShiftCount:             employee.ShiftCount,
			TotalHours:             employee.TotalHours,
			LongestShiftHours:      employee.LongestShiftHours,
			NightShiftCount:        employee.NightShiftCount,
			WeekendShiftCount:      employee.WeekendShiftCount,
			ShortRestCount:         employee.ShortRestCount,
			AvailabilityConflicts:  employee.AvailabilityConflicts,
			OpenSwapRequestCount:   employee.OpenSwapRequestCount,
			MaxConsecutiveWorkDays: employee.MaxConsecutiveWorkDays,
			OvertimeWeekCount:      employee.OvertimeWeekCount,
			RiskScore:              employee.RiskScore,
			RiskLevel:              employee.RiskLevel,
		}
	}
	return scheduleInsightPromptMetrics{
		RequiredHeadcount:     metrics.RequiredHeadcount,
		AssignedHeadcount:     metrics.AssignedHeadcount,
		UnfilledHeadcount:     metrics.UnfilledHeadcount,
		CoverageRate:          metrics.CoverageRate,
		OpenSwapRequestCount:  metrics.OpenSwapRequestCount,
		AverageHours:          metrics.AverageHours,
		WorkloadSpreadHours:   metrics.WorkloadSpreadHours,
		EmployeesAtRisk:       metrics.EmployeesAtRisk,
		AvailabilityConflicts: metrics.AvailabilityConflicts,
		Employees:             employees,
	}
}

func calculateMaxConsecutiveDays(workDates map[string]time.Time) int {
	dates := make([]time.Time, 0, len(workDates))
	for _, date := range workDates {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i int, j int) bool {
		return dates[i].Before(dates[j])
	})

	maxDays := 0
	currentDays := 0
	var previous time.Time
	for _, date := range dates {
		if currentDays == 0 || date.Equal(previous.AddDate(0, 0, 1)) {
			currentDays++
		} else {
			currentDays = 1
		}
		maxDays = max(maxDays, currentDays)
		previous = date
	}
	return maxDays
}

func scheduleRiskLevel(score int) string {
	switch {
	case score >= 10:
		return "critical"
	case score >= 6:
		return "high"
	case score >= 3:
		return "medium"
	default:
		return "low"
	}
}

func latestTime(first time.Time, second time.Time) time.Time {
	if first.After(second) {
		return first
	}
	return second
}

func earliestTime(first time.Time, second time.Time) time.Time {
	if first.Before(second) {
		return first
	}
	return second
}

func round(value float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Round(value*factor) / factor
}

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	emails "github.com/HiIamJeff67/shift-hero-backend/app/emails"
	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	validation "github.com/HiIamJeff67/shift-hero-backend/app/validation"
)

type SchedulingServiceInterface interface {
	CreateShiftRequirement(ctx context.Context, reqDto *dtos.CreateShiftRequirementReqDto) (*dtos.ShiftRequirementResDto, *exceptions.Exception)
	GetShiftRequirements(ctx context.Context, reqDto *dtos.GetShiftRequirementsReqDto) ([]dtos.ShiftRequirementResDto, *exceptions.Exception)
	UpdateShiftRequirement(ctx context.Context, reqDto *dtos.UpdateShiftRequirementReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	DeleteShiftRequirement(ctx context.Context, reqDto *dtos.DeleteShiftRequirementReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	UpsertAvailabilitySlots(ctx context.Context, reqDto *dtos.UpsertAvailabilitySlotsReqDto) ([]dtos.AvailabilitySlotResDto, *exceptions.Exception)
	GetAvailabilitySlots(ctx context.Context, reqDto *dtos.GetAvailabilitySlotsReqDto) ([]dtos.AvailabilitySlotResDto, *exceptions.Exception)
	DeleteAvailabilitySlot(ctx context.Context, reqDto *dtos.DeleteAvailabilitySlotReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	GenerateAssignments(ctx context.Context, reqDto *dtos.GenerateAssignmentsReqDto) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception)
	ReplaceAssignments(ctx context.Context, reqDto *dtos.ReplaceAssignmentsReqDto) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception)
	GetAssignments(ctx context.Context, reqDto *dtos.GetAssignmentsReqDto) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception)
	CreateSwapRequest(ctx context.Context, reqDto *dtos.CreateSwapRequestReqDto) (*dtos.SwapRequestResDto, *exceptions.Exception)
	GetSwapRequests(ctx context.Context, reqDto *dtos.GetSwapRequestsReqDto) ([]dtos.SwapRequestResDto, *exceptions.Exception)
	ClaimSwapRequest(ctx context.Context, reqDto *dtos.ClaimSwapRequestReqDto) (*dtos.SwapRequestResDto, *exceptions.Exception)
	ApproveSwapRequest(ctx context.Context, reqDto *dtos.ApproveSwapRequestReqDto) (*dtos.SwapRequestResDto, *exceptions.Exception)
	CancelSwapRequest(ctx context.Context, reqDto *dtos.CancelSwapRequestReqDto) (*dtos.SwapRequestResDto, *exceptions.Exception)
	GetCompanySettings(ctx context.Context, reqDto *dtos.GetCompanySettingsReqDto) (*dtos.CompanySettingsResDto, *exceptions.Exception)
	UpdateCompanySettings(ctx context.Context, reqDto *dtos.UpdateCompanySettingsReqDto) (*dtos.CompanySettingsResDto, *exceptions.Exception)
}

type SchedulingService struct {
	db *gorm.DB
}

func NewSchedulingService(db *gorm.DB) SchedulingServiceInterface {
	if db == nil {
		db = models.DB
	}
	return &SchedulingService{db: db}
}

func (s *SchedulingService) CreateShiftRequirement(
	ctx context.Context,
	reqDto *dtos.CreateShiftRequirementReqDto,
) (*dtos.ShiftRequirementResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	startAt := truncateToMinute(reqDto.Body.StartAt.UTC())
	endAt := truncateToMinute(reqDto.Body.EndAt.UTC())
	if exception := validateTimeRange(startAt, endAt); exception != nil {
		return nil, exception
	}

	entity := schemas.ShiftRequirement{
		CompanyId:     reqDto.Body.CompanyId,
		EmployeeRole:  reqDto.Body.EmployeeRole,
		StartAt:       startAt,
		EndAt:         endAt,
		RequiredCount: reqDto.Body.RequiredCount,
		Note:          reqDto.Body.Note,
	}
	if err := db.Model(&schemas.ShiftRequirement{}).Create(&entity).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(err)
	}

	res := dtos.ShiftRequirementResDto{
		Id:            entity.Id,
		CompanyId:     entity.CompanyId,
		EmployeeRole:  entity.EmployeeRole,
		StartAt:       entity.StartAt,
		EndAt:         entity.EndAt,
		RequiredCount: entity.RequiredCount,
		Note:          entity.Note,
		UpdatedAt:     entity.UpdatedAt,
		CreatedAt:     entity.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) GetShiftRequirements(
	ctx context.Context,
	reqDto *dtos.GetShiftRequirementsReqDto,
) ([]dtos.ShiftRequirementResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	entities := []schemas.ShiftRequirement{}
	if err := db.Model(&schemas.ShiftRequirement{}).
		Where("company_id = ?", reqDto.Param.CompanyId).
		Order("start_at ASC").
		Find(&entities).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := make([]dtos.ShiftRequirementResDto, len(entities))
	for i, entity := range entities {
		res[i] = dtos.ShiftRequirementResDto{
			Id:            entity.Id,
			CompanyId:     entity.CompanyId,
			EmployeeRole:  entity.EmployeeRole,
			StartAt:       entity.StartAt,
			EndAt:         entity.EndAt,
			RequiredCount: entity.RequiredCount,
			Note:          entity.Note,
			UpdatedAt:     entity.UpdatedAt,
			CreatedAt:     entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) UpdateShiftRequirement(
	ctx context.Context,
	reqDto *dtos.UpdateShiftRequirementReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	existing := schemas.ShiftRequirement{}
	if err := db.Model(&schemas.ShiftRequirement{}).
		Where("id = ? AND company_id = ?", reqDto.Body.ShiftRequirementId, reqDto.Body.CompanyId).
		First(&existing).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	updates := map[string]any{}
	startAt := existing.StartAt
	endAt := existing.EndAt
	if reqDto.Body.Values.EmployeeRole != nil {
		updates["employee_role"] = *reqDto.Body.Values.EmployeeRole
	}
	if reqDto.Body.Values.StartAt != nil {
		startAt = truncateToMinute(reqDto.Body.Values.StartAt.UTC())
		updates["start_at"] = startAt
	}
	if reqDto.Body.Values.EndAt != nil {
		endAt = truncateToMinute(reqDto.Body.Values.EndAt.UTC())
		updates["end_at"] = endAt
	}
	if reqDto.Body.Values.RequiredCount != nil {
		updates["required_count"] = *reqDto.Body.Values.RequiredCount
	}
	if reqDto.Body.Values.Note != nil {
		updates["note"] = *reqDto.Body.Values.Note
	}

	if exception := validateTimeRange(startAt, endAt); exception != nil {
		return nil, exception
	}
	if len(updates) == 0 {
		return nil, exceptions.Scheduling.NoChanges()
	}

	result := db.Model(&schemas.ShiftRequirement{}).
		Where("id = ? AND company_id = ?", reqDto.Body.ShiftRequirementId, reqDto.Body.CompanyId).
		Updates(updates)
	if result.Error != nil {
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Scheduling.NoChanges()
	}

	updated := schemas.ShiftRequirement{}
	if err := db.Model(&schemas.ShiftRequirement{}).Where("id = ?", reqDto.Body.ShiftRequirementId).First(&updated).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	return &dtos.MutationUpdatedAtResDto{UpdatedAt: updated.UpdatedAt}, nil
}

func (s *SchedulingService) DeleteShiftRequirement(
	ctx context.Context,
	reqDto *dtos.DeleteShiftRequirementReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	result := db.Model(&schemas.ShiftRequirement{}).
		Where("id = ? AND company_id = ?", reqDto.Body.ShiftRequirementId, reqDto.Body.CompanyId).
		Delete(&schemas.ShiftRequirement{})
	if result.Error != nil {
		return nil, exceptions.Scheduling.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Scheduling.NotFound()
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: truncateToMinute(time.Now().UTC())}, nil
}

func (s *SchedulingService) UpsertAvailabilitySlots(
	ctx context.Context,
	reqDto *dtos.UpsertAvailabilitySlotsReqDto,
) ([]dtos.AvailabilitySlotResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(tx.Error)
	}

	created := []schemas.AvailabilitySlot{}
	for _, slot := range reqDto.Body.Slots {
		startAt := truncateToMinute(slot.StartAt.UTC())
		endAt := truncateToMinute(slot.EndAt.UTC())
		if exception := validateTimeRange(startAt, endAt); exception != nil {
			tx.Rollback()
			return nil, exception
		}

		entity := schemas.AvailabilitySlot{}
		err := tx.Table(schemas.AvailabilitySlot{}.TableName()).
			Where("company_id = ? AND user_id = ? AND start_at = ? AND end_at = ?", reqDto.Body.CompanyId, reqDto.ContextFields.UserId, startAt, endAt).
			First(&entity).Error

		if err == nil {
			if updateErr := tx.Table(schemas.AvailabilitySlot{}.TableName()).
				Where("id = ?", entity.Id).
				Updates(map[string]any{"is_available": slot.IsAvailable}).Error; updateErr != nil {
				tx.Rollback()
				return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(updateErr)
			}
			if reloadErr := tx.Table(schemas.AvailabilitySlot{}.TableName()).Where("id = ?", entity.Id).First(&entity).Error; reloadErr != nil {
				tx.Rollback()
				return nil, exceptions.Scheduling.NotFound().WithOrigin(reloadErr)
			}
		} else {
			entity = schemas.AvailabilitySlot{
				CompanyId:   reqDto.Body.CompanyId,
				UserId:      reqDto.ContextFields.UserId,
				StartAt:     startAt,
				EndAt:       endAt,
				IsAvailable: slot.IsAvailable,
			}
			if createErr := tx.Table(schemas.AvailabilitySlot{}.TableName()).Create(&entity).Error; createErr != nil {
				tx.Rollback()
				return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(createErr)
			}
		}
		created = append(created, entity)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(err)
	}

	res := make([]dtos.AvailabilitySlotResDto, len(created))
	for i, entity := range created {
		res[i] = dtos.AvailabilitySlotResDto{
			Id:          entity.Id,
			CompanyId:   entity.CompanyId,
			UserId:      entity.UserId,
			StartAt:     entity.StartAt,
			EndAt:       entity.EndAt,
			IsAvailable: entity.IsAvailable,
			UpdatedAt:   entity.UpdatedAt,
			CreatedAt:   entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) GetAvailabilitySlots(
	ctx context.Context,
	reqDto *dtos.GetAvailabilitySlotsReqDto,
) ([]dtos.AvailabilitySlotResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	query := db.Model(&schemas.AvailabilitySlot{}).
		Where("company_id = ?", reqDto.Param.CompanyId)
	if reqDto.Body.UserId != nil {
		query = query.Where("user_id = ?", *reqDto.Body.UserId)
	}
	if reqDto.Body.StartAt != nil {
		query = query.Where("end_at >= ?", truncateToMinute(reqDto.Body.StartAt.UTC()))
	}
	if reqDto.Body.EndAt != nil {
		query = query.Where("start_at <= ?", truncateToMinute(reqDto.Body.EndAt.UTC()))
	}

	entities := []schemas.AvailabilitySlot{}
	if err := query.Order("start_at ASC").Find(&entities).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := make([]dtos.AvailabilitySlotResDto, len(entities))
	for i, entity := range entities {
		res[i] = dtos.AvailabilitySlotResDto{
			Id:          entity.Id,
			CompanyId:   entity.CompanyId,
			UserId:      entity.UserId,
			StartAt:     entity.StartAt,
			EndAt:       entity.EndAt,
			IsAvailable: entity.IsAvailable,
			UpdatedAt:   entity.UpdatedAt,
			CreatedAt:   entity.CreatedAt,
		}
	}

	return res, nil
}

func (s *SchedulingService) DeleteAvailabilitySlot(
	ctx context.Context,
	reqDto *dtos.DeleteAvailabilitySlotReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	membership, exception := requireCompanyMember(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	existing := schemas.AvailabilitySlot{}
	if err := db.Model(&schemas.AvailabilitySlot{}).
		Where("id = ? AND company_id = ?", reqDto.Body.AvailabilitySlotId, reqDto.Body.CompanyId).
		First(&existing).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	if membership.EmployeeRole != enums.EmployeeRole_Manager && existing.UserId != reqDto.ContextFields.UserId {
		return nil, exceptions.Scheduling.Forbidden("You can only delete your own availability slot")
	}

	if err := db.Model(&schemas.AvailabilitySlot{}).
		Where("id = ?", reqDto.Body.AvailabilitySlotId).
		Delete(&schemas.AvailabilitySlot{}).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToDelete().WithOrigin(err)
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: truncateToMinute(time.Now().UTC())}, nil
}

func (s *SchedulingService) GenerateAssignments(
	ctx context.Context,
	reqDto *dtos.GenerateAssignmentsReqDto,
) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	requirementsQuery := db.Model(&schemas.ShiftRequirement{}).
		Where("company_id = ?", reqDto.Body.CompanyId)
	if reqDto.Body.StartAt != nil {
		requirementsQuery = requirementsQuery.Where("start_at >= ?", truncateToMinute(reqDto.Body.StartAt.UTC()))
	}
	if reqDto.Body.EndAt != nil {
		requirementsQuery = requirementsQuery.Where("end_at <= ?", truncateToMinute(reqDto.Body.EndAt.UTC()))
	}

	requirements := []schemas.ShiftRequirement{}
	if err := requirementsQuery.Order("start_at ASC").Find(&requirements).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	assignments := []schemas.ShiftAssignment{}
	for _, requirement := range requirements {
		type candidateRow struct {
			UserId uuid.UUID `gorm:"column:user_id"`
		}
		candidates := []candidateRow{}
		err := db.Model(&schemas.AvailabilitySlot{}).
			Select("\"AvailabilitySlotsTable\".user_id").
			Joins("JOIN \"UsersToCompaniesTable\" utc ON utc.user_id = \"AvailabilitySlotsTable\".user_id AND utc.company_id = \"AvailabilitySlotsTable\".company_id").
			Where("\"AvailabilitySlotsTable\".company_id = ?", reqDto.Body.CompanyId).
			Where("\"AvailabilitySlotsTable\".is_available = true").
			Where("\"AvailabilitySlotsTable\".start_at <= ? AND \"AvailabilitySlotsTable\".end_at >= ?", requirement.StartAt, requirement.EndAt).
			Where("utc.employee_role = ?", requirement.EmployeeRole).
			Group("\"AvailabilitySlotsTable\".user_id").
			Order("\"AvailabilitySlotsTable\".user_id ASC").
			Limit(int(requirement.RequiredCount)).
			Find(&candidates).Error
		if err != nil {
			return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
		}

		for _, candidate := range candidates {
			assignment := schemas.ShiftAssignment{
				CompanyId:          reqDto.Body.CompanyId,
				ShiftRequirementId: requirement.Id,
				UserId:             candidate.UserId,
				StartAt:            requirement.StartAt,
				EndAt:              requirement.EndAt,
			}
			if err := db.Model(&schemas.ShiftAssignment{}).Create(&assignment).Error; err != nil {
				return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(err)
			}
			assignments = append(assignments, assignment)
		}
	}

	res := make([]dtos.ShiftAssignmentResDto, len(assignments))
	for i, entity := range assignments {
		res[i] = dtos.ShiftAssignmentResDto{
			Id:                 entity.Id,
			CompanyId:          entity.CompanyId,
			ShiftRequirementId: entity.ShiftRequirementId,
			UserId:             entity.UserId,
			StartAt:            entity.StartAt,
			EndAt:              entity.EndAt,
			UpdatedAt:          entity.UpdatedAt,
			CreatedAt:          entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) ReplaceAssignments(
	ctx context.Context,
	reqDto *dtos.ReplaceAssignmentsReqDto,
) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(tx.Error)
	}

	if err := tx.Table(schemas.ShiftAssignment{}.TableName()).
		Where("company_id = ?", reqDto.Body.CompanyId).
		Delete(&schemas.ShiftAssignment{}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Scheduling.FailedToDelete().WithOrigin(err)
	}

	created := make([]schemas.ShiftAssignment, 0, len(reqDto.Body.Assignments))
	for _, item := range reqDto.Body.Assignments {
		startAt := truncateToMinute(item.StartAt.UTC())
		endAt := truncateToMinute(item.EndAt.UTC())
		if exception := validateTimeRange(startAt, endAt); exception != nil {
			tx.Rollback()
			return nil, exception
		}

		entity := schemas.ShiftAssignment{
			CompanyId:          reqDto.Body.CompanyId,
			ShiftRequirementId: item.ShiftRequirementId,
			UserId:             item.UserId,
			StartAt:            startAt,
			EndAt:              endAt,
		}
		if err := tx.Table(schemas.ShiftAssignment{}.TableName()).Create(&entity).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(err)
		}
		created = append(created, entity)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(err)
	}

	res := make([]dtos.ShiftAssignmentResDto, len(created))
	for i, entity := range created {
		res[i] = dtos.ShiftAssignmentResDto{
			Id:                 entity.Id,
			CompanyId:          entity.CompanyId,
			ShiftRequirementId: entity.ShiftRequirementId,
			UserId:             entity.UserId,
			StartAt:            entity.StartAt,
			EndAt:              entity.EndAt,
			UpdatedAt:          entity.UpdatedAt,
			CreatedAt:          entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) GetAssignments(
	ctx context.Context,
	reqDto *dtos.GetAssignmentsReqDto,
) ([]dtos.ShiftAssignmentResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	query := db.Model(&schemas.ShiftAssignment{}).
		Where("company_id = ?", reqDto.Param.CompanyId)
	if reqDto.Body.UserId != nil {
		query = query.Where("user_id = ?", *reqDto.Body.UserId)
	}
	if reqDto.Body.StartAt != nil {
		query = query.Where("end_at >= ?", truncateToMinute(reqDto.Body.StartAt.UTC()))
	}
	if reqDto.Body.EndAt != nil {
		query = query.Where("start_at <= ?", truncateToMinute(reqDto.Body.EndAt.UTC()))
	}

	entities := []schemas.ShiftAssignment{}
	if err := query.Order("start_at ASC").Find(&entities).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := make([]dtos.ShiftAssignmentResDto, len(entities))
	for i, entity := range entities {
		res[i] = dtos.ShiftAssignmentResDto{
			Id:                 entity.Id,
			CompanyId:          entity.CompanyId,
			ShiftRequirementId: entity.ShiftRequirementId,
			UserId:             entity.UserId,
			StartAt:            entity.StartAt,
			EndAt:              entity.EndAt,
			UpdatedAt:          entity.UpdatedAt,
			CreatedAt:          entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) CreateSwapRequest(
	ctx context.Context,
	reqDto *dtos.CreateSwapRequestReqDto,
) (*dtos.SwapRequestResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	assignment := schemas.ShiftAssignment{}
	if err := db.Model(&schemas.ShiftAssignment{}).
		Where("id = ? AND company_id = ?", reqDto.Body.ShiftAssignmentId, reqDto.Body.CompanyId).
		First(&assignment).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	if assignment.UserId != reqDto.ContextFields.UserId {
		return nil, exceptions.Scheduling.Forbidden("You can only request swap for your own assignment")
	}

	entity := schemas.SwapRequest{
		CompanyId:         reqDto.Body.CompanyId,
		ShiftAssignmentId: reqDto.Body.ShiftAssignmentId,
		RequesterUserId:   reqDto.ContextFields.UserId,
		Status:            enums.SwapRequestStatus_Open,
		Reason:            reqDto.Body.Reason,
	}
	if err := db.Model(&schemas.SwapRequest{}).Create(&entity).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(err)
	}

	res := dtos.SwapRequestResDto{
		Id:                entity.Id,
		CompanyId:         entity.CompanyId,
		ShiftAssignmentId: entity.ShiftAssignmentId,
		RequesterUserId:   entity.RequesterUserId,
		ClaimedByUserId:   entity.ClaimedByUserId,
		Status:            entity.Status,
		Reason:            entity.Reason,
		UpdatedAt:         entity.UpdatedAt,
		CreatedAt:         entity.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) GetSwapRequests(
	ctx context.Context,
	reqDto *dtos.GetSwapRequestsReqDto,
) ([]dtos.SwapRequestResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	query := db.Model(&schemas.SwapRequest{}).
		Where("company_id = ?", reqDto.Param.CompanyId)
	if reqDto.Body.Status != nil {
		query = query.Where("status = ?", *reqDto.Body.Status)
	}

	entities := []schemas.SwapRequest{}
	if err := query.Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := make([]dtos.SwapRequestResDto, len(entities))
	for i, entity := range entities {
		res[i] = dtos.SwapRequestResDto{
			Id:                entity.Id,
			CompanyId:         entity.CompanyId,
			ShiftAssignmentId: entity.ShiftAssignmentId,
			RequesterUserId:   entity.RequesterUserId,
			ClaimedByUserId:   entity.ClaimedByUserId,
			Status:            entity.Status,
			Reason:            entity.Reason,
			UpdatedAt:         entity.UpdatedAt,
			CreatedAt:         entity.CreatedAt,
		}
	}
	return res, nil
}

func (s *SchedulingService) ClaimSwapRequest(
	ctx context.Context,
	reqDto *dtos.ClaimSwapRequestReqDto,
) (*dtos.SwapRequestResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	swap := schemas.SwapRequest{}
	if err := db.Model(&schemas.SwapRequest{}).
		Where("id = ? AND company_id = ?", reqDto.Body.SwapRequestId, reqDto.Body.CompanyId).
		First(&swap).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	if swap.Status != enums.SwapRequestStatus_Open {
		return nil, exceptions.Scheduling.InvalidSwapState("Only open swap requests can be claimed")
	}
	if swap.RequesterUserId == reqDto.ContextFields.UserId {
		return nil, exceptions.Scheduling.Forbidden("Requester cannot claim own swap")
	}

	swap.Status = enums.SwapRequestStatus_Claimed
	swap.ClaimedByUserId = &reqDto.ContextFields.UserId
	if err := db.Model(&schemas.SwapRequest{}).
		Where("id = ?", swap.Id).
		Updates(map[string]any{"status": swap.Status, "claimed_by_user_id": swap.ClaimedByUserId}).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(err)
	}
	if err := db.Model(&schemas.SwapRequest{}).Where("id = ?", swap.Id).First(&swap).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	s.sendSwapClaimedEmail(db, swap)

	res := dtos.SwapRequestResDto{
		Id:                swap.Id,
		CompanyId:         swap.CompanyId,
		ShiftAssignmentId: swap.ShiftAssignmentId,
		RequesterUserId:   swap.RequesterUserId,
		ClaimedByUserId:   swap.ClaimedByUserId,
		Status:            swap.Status,
		Reason:            swap.Reason,
		UpdatedAt:         swap.UpdatedAt,
		CreatedAt:         swap.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) ApproveSwapRequest(
	ctx context.Context,
	reqDto *dtos.ApproveSwapRequestReqDto,
) (*dtos.SwapRequestResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(tx.Error)
	}

	swap := schemas.SwapRequest{}
	if err := tx.Table(schemas.SwapRequest{}.TableName()).
		Where("id = ? AND company_id = ?", reqDto.Body.SwapRequestId, reqDto.Body.CompanyId).
		First(&swap).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	if swap.Status != enums.SwapRequestStatus_Claimed {
		tx.Rollback()
		return nil, exceptions.Scheduling.InvalidSwapState("Only claimed swap requests can be approved")
	}
	if swap.ClaimedByUserId == nil {
		tx.Rollback()
		return nil, exceptions.Scheduling.InvalidSwapState("Cannot approve an unclaimed swap request")
	}

	if err := tx.Table(schemas.ShiftAssignment{}.TableName()).
		Where("id = ? AND company_id = ?", swap.ShiftAssignmentId, reqDto.Body.CompanyId).
		Updates(map[string]any{"user_id": *swap.ClaimedByUserId}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Table(schemas.SwapRequest{}.TableName()).
		Where("id = ?", swap.Id).
		Updates(map[string]any{"status": enums.SwapRequestStatus_Approved}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Scheduling.FailedToCommitTransaction().WithOrigin(err)
	}

	if err := db.Model(&schemas.SwapRequest{}).Where("id = ?", swap.Id).First(&swap).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	s.sendSwapApprovedEmail(db, swap)

	res := dtos.SwapRequestResDto{
		Id:                swap.Id,
		CompanyId:         swap.CompanyId,
		ShiftAssignmentId: swap.ShiftAssignmentId,
		RequesterUserId:   swap.RequesterUserId,
		ClaimedByUserId:   swap.ClaimedByUserId,
		Status:            swap.Status,
		Reason:            swap.Reason,
		UpdatedAt:         swap.UpdatedAt,
		CreatedAt:         swap.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) CancelSwapRequest(
	ctx context.Context,
	reqDto *dtos.CancelSwapRequestReqDto,
) (*dtos.SwapRequestResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	membership, exception := requireCompanyMember(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId)
	if exception != nil {
		return nil, exception
	}

	swap := schemas.SwapRequest{}
	if err := db.Model(&schemas.SwapRequest{}).
		Where("id = ? AND company_id = ?", reqDto.Body.SwapRequestId, reqDto.Body.CompanyId).
		First(&swap).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}
	if swap.Status == enums.SwapRequestStatus_Approved || swap.Status == enums.SwapRequestStatus_Cancelled {
		return nil, exceptions.Scheduling.InvalidSwapState("Approved or cancelled swap requests cannot be cancelled")
	}
	if membership.EmployeeRole != enums.EmployeeRole_Manager && swap.RequesterUserId != reqDto.ContextFields.UserId {
		return nil, exceptions.Scheduling.Forbidden("Only manager or requester can cancel swap request")
	}

	if err := db.Model(&schemas.SwapRequest{}).
		Where("id = ?", swap.Id).
		Updates(map[string]any{"status": enums.SwapRequestStatus_Cancelled}).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(err)
	}
	if err := db.Model(&schemas.SwapRequest{}).Where("id = ?", swap.Id).First(&swap).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := dtos.SwapRequestResDto{
		Id:                swap.Id,
		CompanyId:         swap.CompanyId,
		ShiftAssignmentId: swap.ShiftAssignmentId,
		RequesterUserId:   swap.RequesterUserId,
		ClaimedByUserId:   swap.ClaimedByUserId,
		Status:            swap.Status,
		Reason:            swap.Reason,
		UpdatedAt:         swap.UpdatedAt,
		CreatedAt:         swap.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) GetCompanySettings(
	ctx context.Context,
	reqDto *dtos.GetCompanySettingsReqDto,
) (*dtos.CompanySettingsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMember(db, reqDto.Param.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	settings := schemas.CompanySettings{}
	err := db.Model(&schemas.CompanySettings{}).
		Where("company_id = ?", reqDto.Param.CompanyId).
		First(&settings).Error
	if err != nil {
		settings = schemas.CompanySettings{CompanyId: reqDto.Param.CompanyId}
		if createErr := db.Model(&schemas.CompanySettings{}).Create(&settings).Error; createErr != nil {
			return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(createErr)
		}
	}

	res := dtos.CompanySettingsResDto{
		CompanyId:        settings.CompanyId,
		AutoApproveSwaps: settings.AutoApproveSwaps,
		MaxWeeklyHours:   settings.MaxWeeklyHours,
		MinRestHours:     settings.MinRestHours,
		UpdatedAt:        settings.UpdatedAt,
		CreatedAt:        settings.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) UpdateCompanySettings(
	ctx context.Context,
	reqDto *dtos.UpdateCompanySettingsReqDto,
) (*dtos.CompanySettingsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Scheduling.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManager(db, reqDto.Body.CompanyId, reqDto.ContextFields.UserId); exception != nil {
		return nil, exception
	}

	updates := map[string]any{}
	if reqDto.Body.Values.AutoApproveSwaps != nil {
		updates["auto_approve_swaps"] = *reqDto.Body.Values.AutoApproveSwaps
	}
	if reqDto.Body.Values.MaxWeeklyHours != nil {
		updates["max_weekly_hours"] = *reqDto.Body.Values.MaxWeeklyHours
	}
	if reqDto.Body.Values.MinRestHours != nil {
		updates["min_rest_hours"] = *reqDto.Body.Values.MinRestHours
	}
	if len(updates) == 0 {
		return nil, exceptions.Scheduling.NoChanges()
	}

	settings := schemas.CompanySettings{}
	if err := db.Model(&schemas.CompanySettings{}).
		Where("company_id = ?", reqDto.Body.CompanyId).
		First(&settings).Error; err != nil {
		settings = schemas.CompanySettings{CompanyId: reqDto.Body.CompanyId}
		if createErr := db.Model(&schemas.CompanySettings{}).Create(&settings).Error; createErr != nil {
			return nil, exceptions.Scheduling.FailedToCreate().WithOrigin(createErr)
		}
	}

	if err := db.Model(&schemas.CompanySettings{}).
		Where("company_id = ?", reqDto.Body.CompanyId).
		Updates(updates).Error; err != nil {
		return nil, exceptions.Scheduling.FailedToUpdate().WithOrigin(err)
	}

	if err := db.Model(&schemas.CompanySettings{}).Where("company_id = ?", reqDto.Body.CompanyId).First(&settings).Error; err != nil {
		return nil, exceptions.Scheduling.NotFound().WithOrigin(err)
	}

	res := dtos.CompanySettingsResDto{
		CompanyId:        settings.CompanyId,
		AutoApproveSwaps: settings.AutoApproveSwaps,
		MaxWeeklyHours:   settings.MaxWeeklyHours,
		MinRestHours:     settings.MinRestHours,
		UpdatedAt:        settings.UpdatedAt,
		CreatedAt:        settings.CreatedAt,
	}
	return &res, nil
}

func (s *SchedulingService) sendSwapClaimedEmail(db *gorm.DB, swap schemas.SwapRequest) {
	requester := schemas.User{}
	if err := db.Model(&schemas.User{}).Where("id = ?", swap.RequesterUserId).First(&requester).Error; err != nil {
		return
	}
	company := schemas.Company{}
	if err := db.Model(&schemas.Company{}).Where("id = ?", swap.CompanyId).First(&company).Error; err != nil {
		return
	}
	assignmentSummary := s.getAssignmentSummary(db, swap.ShiftAssignmentId)
	if exception := emails.AsyncSendSwapClaimedEmail(requester.Email, company.Name, assignmentSummary); exception != nil {
		exception.Log()
	}
}

func (s *SchedulingService) sendSwapApprovedEmail(db *gorm.DB, swap schemas.SwapRequest) {
	requester := schemas.User{}
	if err := db.Model(&schemas.User{}).Where("id = ?", swap.RequesterUserId).First(&requester).Error; err != nil {
		return
	}
	company := schemas.Company{}
	if err := db.Model(&schemas.Company{}).Where("id = ?", swap.CompanyId).First(&company).Error; err != nil {
		return
	}
	assignmentSummary := s.getAssignmentSummary(db, swap.ShiftAssignmentId)
	if exception := emails.AsyncSendSwapApprovedEmail(requester.Email, company.Name, assignmentSummary); exception != nil {
		exception.Log()
	}
}

func (s *SchedulingService) getAssignmentSummary(db *gorm.DB, assignmentId uuid.UUID) string {
	assignment := schemas.ShiftAssignment{}
	if err := db.Model(&schemas.ShiftAssignment{}).Where("id = ?", assignmentId).First(&assignment).Error; err != nil {
		return assignmentId.String()
	}
	return fmt.Sprintf("%s ~ %s", assignment.StartAt.Format(time.RFC3339), assignment.EndAt.Format(time.RFC3339))
}

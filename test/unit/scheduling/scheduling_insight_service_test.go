package schedulingunittest

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

func TestCalculateScheduleInsightMetrics(t *testing.T) {
	t.Parallel()

	aliceId := uuid.New()
	bobId := uuid.New()
	requirementId := uuid.New()
	startAt := time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)
	endAt := startAt.Add(7 * 24 * time.Hour)

	members := []services.ScheduleInsightMember{
		{UserId: aliceId, DisplayName: "Alice", EmployeeRole: enums.EmployeeRole_Staff},
		{UserId: bobId, DisplayName: "Bob", EmployeeRole: enums.EmployeeRole_Staff},
	}
	requirements := []schemas.ShiftRequirement{
		{
			Id:            requirementId,
			RequiredCount: 2,
			StartAt:       startAt.Add(20 * time.Hour),
			EndAt:         startAt.Add(23 * time.Hour),
		},
	}
	assignments := []schemas.ShiftAssignment{
		{
			ShiftRequirementId: requirementId,
			UserId:             aliceId,
			StartAt:            startAt.Add(20 * time.Hour),
			EndAt:              startAt.Add(23 * time.Hour),
		},
		{
			ShiftRequirementId: uuid.New(),
			UserId:             aliceId,
			StartAt:            startAt.Add(29 * time.Hour),
			EndAt:              startAt.Add(37 * time.Hour),
		},
	}
	availability := []schemas.AvailabilitySlot{
		{
			UserId:      aliceId,
			StartAt:     startAt.Add(19 * time.Hour),
			EndAt:       startAt.Add(24 * time.Hour),
			IsAvailable: true,
		},
	}
	activeSwaps := []schemas.SwapRequest{
		{
			RequesterUserId: aliceId,
			Status:          enums.SwapRequestStatus_Open,
		},
	}
	settings := schemas.CompanySettings{
		MaxWeeklyHours: 40,
		MinRestHours:   8,
	}

	metrics := services.CalculateScheduleInsightMetrics(
		members,
		requirements,
		assignments,
		availability,
		activeSwaps,
		settings,
		time.UTC,
		startAt,
		endAt,
	)

	require.Equal(t, int64(2), metrics.RequiredHeadcount)
	require.Equal(t, int64(1), metrics.AssignedHeadcount)
	require.Equal(t, int64(1), metrics.UnfilledHeadcount)
	require.Equal(t, 0.5, metrics.CoverageRate)
	require.Equal(t, 5.5, metrics.AverageHours)
	require.Equal(t, 11.0, metrics.WorkloadSpreadHours)
	require.Equal(t, 1, metrics.EmployeesAtRisk)
	require.Equal(t, 1, metrics.AvailabilityConflicts)
	require.Len(t, metrics.Employees, 2)

	alice := metrics.Employees[0]
	require.Equal(t, "Alice", alice.DisplayName)
	require.Equal(t, 2, alice.ShiftCount)
	require.Equal(t, 11.0, alice.TotalHours)
	require.Equal(t, 1, alice.ShortRestCount)
	require.Equal(t, 1, alice.AvailabilityConflicts)
	require.Equal(t, 1, alice.OpenSwapRequestCount)
	require.Equal(t, 2, alice.MaxConsecutiveWorkDays)
	require.Equal(t, "high", alice.RiskLevel)
}

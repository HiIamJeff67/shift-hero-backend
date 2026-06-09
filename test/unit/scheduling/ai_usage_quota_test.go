package schedulingunittest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

func TestEvaluateMonthlyAIUsageReservesOneGeneration(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.June, 10, 12, 0, 0, 0, time.UTC)
	decision := services.EvaluateMonthlyAIUsage(
		4,
		5,
		time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		now,
	)

	require.True(t, decision.Allowed)
	require.Equal(t, int32(5), decision.Used)
	require.Equal(t, time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC), decision.ResetAt)
}

func TestEvaluateMonthlyAIUsageRejectsExhaustedQuota(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.June, 10, 12, 0, 0, 0, time.UTC)
	decision := services.EvaluateMonthlyAIUsage(
		5,
		5,
		time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		now,
	)

	require.False(t, decision.Allowed)
	require.Equal(t, int32(5), decision.Used)
}

func TestEvaluateMonthlyAIUsageResetsForNewMonth(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 2, 12, 0, 0, 0, time.UTC)
	decision := services.EvaluateMonthlyAIUsage(
		5,
		5,
		time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		now,
	)

	require.True(t, decision.Allowed)
	require.Equal(t, int32(1), decision.Used)
	require.Equal(t, time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC), decision.PeriodStart)
}

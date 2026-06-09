package aiintegrationtest

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ai "github.com/HiIamJeff67/shift-hero-backend/app/ai"
)

func TestOpenRouterScheduleInsightGeneratorIntegration(t *testing.T) {
	if os.Getenv("RUN_OPEN_ROUTER_INTEGRATION") != "1" {
		t.Skip("set RUN_OPEN_ROUTER_INTEGRATION=1 to call OpenRouter")
	}

	generator, err := ai.NewOpenRouterScheduleInsightGenerator(
		os.Getenv("OPEN_ROUTER_API_KEY"),
		os.Getenv("OPEN_ROUTER_MODEL"),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	var streamed strings.Builder
	result, err := generator.Generate(
		ctx,
		ai.ScheduleInsightPromptInput{
			Locale: "en",
			Focus:  "coverage and fatigue",
			SnapshotJSON: `{
				"companyName": "Integration Test",
				"metrics": {
					"coverageRate": 0.75,
					"unfilledHeadcount": 2,
					"employeesAtRisk": 1,
					"employees": []
				}
			}`,
		},
		func(_ context.Context, chunk []byte) error {
			streamed.Write(chunk)
			return nil
		},
	)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, result, streamed.String())
}

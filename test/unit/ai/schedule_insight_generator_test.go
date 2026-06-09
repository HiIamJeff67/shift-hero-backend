package aiunittest

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"

	ai "github.com/HiIamJeff67/shift-hero-backend/app/ai"
)

type scheduleInsightFailingLLM struct{}

func (scheduleInsightFailingLLM) Call(
	ctx context.Context,
	prompt string,
	options ...llms.CallOption,
) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, scheduleInsightFailingLLM{}, prompt, options...)
}

func (scheduleInsightFailingLLM) GenerateContent(
	context.Context,
	[]llms.MessageContent,
	...llms.CallOption,
) (*llms.ContentResponse, error) {
	return nil, errors.New("provider unavailable")
}

type scheduleInsightFakeLLM struct {
	mu      sync.Mutex
	calls   int
	prompts []string
}

func (f *scheduleInsightFakeLLM) Call(
	ctx context.Context,
	prompt string,
	options ...llms.CallOption,
) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, f, prompt, options...)
}

func (f *scheduleInsightFakeLLM) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	prompt := messages[0].Parts[0].(llms.TextContent).Text
	f.mu.Lock()
	call := f.calls
	f.calls++
	f.prompts = append(f.prompts, prompt)
	f.mu.Unlock()

	content := "Thought: coverage is the main issue\nFinal Answer: Coverage is 50% and Alice needs attention."
	if call > 0 {
		content = "# AI pulse\nCoverage needs immediate action."
		callOptions := llms.CallOptions{}
		for _, option := range options {
			option(&callOptions)
		}
		if callOptions.StreamingFunc != nil {
			if err := callOptions.StreamingFunc(ctx, []byte("# AI pulse\n")); err != nil {
				return nil, err
			}
			if err := callOptions.StreamingFunc(ctx, []byte("Coverage needs immediate action.")); err != nil {
				return nil, err
			}
		}
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{{Content: content}},
	}, nil
}

func TestScheduleInsightGeneratorRunsAgentThenStreamsOnlyFinalEditor(t *testing.T) {
	t.Parallel()

	fakeLLM := &scheduleInsightFakeLLM{}
	generator := ai.NewLangChainScheduleInsightGenerator(fakeLLM, "test-model")
	var streamed strings.Builder

	result, err := generator.Generate(
		context.Background(),
		ai.ScheduleInsightPromptInput{
			Locale:       "en",
			Focus:        "coverage",
			SnapshotJSON: `{"coverageRate":0.5}`,
		},
		func(_ context.Context, chunk []byte) error {
			streamed.Write(chunk)
			return nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, "# AI pulse\nCoverage needs immediate action.", result)
	require.Equal(t, result, streamed.String())
	require.Equal(t, 2, fakeLLM.calls)
	require.Len(t, fakeLLM.prompts, 2)
	require.Contains(t, fakeLLM.prompts[0], "Workforce Auditor Agent")
	require.Contains(t, fakeLLM.prompts[1], "Coverage is 50% and Alice needs attention.")
	require.NotContains(t, streamed.String(), "Thought:")
}

func TestScheduleInsightGeneratorFallsBackBeforeStreamingStarts(t *testing.T) {
	t.Parallel()

	fallbackLLM := &scheduleInsightFakeLLM{}
	generator := ai.NewLangChainScheduleInsightGeneratorWithFallback(
		scheduleInsightFailingLLM{},
		fallbackLLM,
		"primary",
	)
	var streamed strings.Builder

	result, err := generator.Generate(
		context.Background(),
		ai.ScheduleInsightPromptInput{
			Locale:       "en",
			Focus:        "coverage",
			SnapshotJSON: `{"coverageRate":0.5}`,
		},
		func(_ context.Context, chunk []byte) error {
			streamed.Write(chunk)
			return nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, result, streamed.String())
	require.Equal(t, 2, fallbackLLM.calls)
}

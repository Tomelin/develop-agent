package agentsdk

import (
	"context"
	"time"

	"github.com/develop-agent/backend/pkg/observability"
)

// MetricsProvider decorates an LLM provider and emits Prometheus metrics.
type MetricsProvider struct {
	Base Provider
}

func (m *MetricsProvider) Initialize(ctx context.Context, cfg Config) error {
	return m.Base.Initialize(ctx, cfg)
}

func (m *MetricsProvider) Name() string      { return m.Base.Name() }
func (m *MetricsProvider) ModelName() string { return m.Base.ModelName() }
func (m *MetricsProvider) Close() error      { return m.Base.Close() }

func (m *MetricsProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	started := time.Now()
	resp, err := m.Base.Complete(ctx, req)
	promptTokens := 0
	completionTokens := 0
	if err == nil {
		promptTokens = resp.Usage.InputTokens
		completionTokens = resp.Usage.OutputTokens
	}
	observability.ObserveLLMCall(m.Base.Name(), m.Base.ModelName(), time.Since(started), promptTokens, completionTokens, err)
	if err != nil {
		return CompletionResponse{}, err
	}
	return resp, nil
}

func (m *MetricsProvider) Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error) {
	return m.Base.Stream(ctx, req)
}

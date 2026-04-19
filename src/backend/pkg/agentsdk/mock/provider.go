package mock

import (
	"context"
	"sync"

	"github.com/develop-agent/backend/pkg/agentsdk"
)

// Provider is a deterministic provider for tests.
type Provider struct {
	mu        sync.Mutex
	cfg       agentsdk.Config
	responses []agentsdk.CompletionResponse
	idx       int
}

func New(responses ...agentsdk.CompletionResponse) *Provider {
	return &Provider{responses: responses}
}

func (p *Provider) Initialize(_ context.Context, cfg agentsdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cfg = cfg
	p.idx = 0
	return nil
}

func (p *Provider) Complete(_ context.Context, req agentsdk.CompletionRequest) (agentsdk.CompletionResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.idx < len(p.responses) {
		resp := p.responses[p.idx]
		p.idx++
		return resp, nil
	}

	content := ""
	if len(req.Messages) > 0 {
		content = req.Messages[len(req.Messages)-1].Content
	}
	usage := agentsdk.Usage{InputTokens: len(req.Messages) * 10, OutputTokens: 10}
	usage.TotalTokens = usage.InputTokens + usage.OutputTokens

	return agentsdk.CompletionResponse{
		Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: "mock: " + content},
		Usage:   usage,
	}, nil
}

func (p *Provider) Stream(ctx context.Context, req agentsdk.CompletionRequest) (<-chan agentsdk.StreamEvent, error) {
	ch := make(chan agentsdk.StreamEvent, 2)
	go func() {
		defer close(ch)
		resp, err := p.Complete(ctx, req)
		if err != nil {
			ch <- agentsdk.StreamEvent{Err: err, Done: true}
			return
		}
		ch <- agentsdk.StreamEvent{Delta: resp.Message.Content}
		ch <- agentsdk.StreamEvent{Done: true, StopReason: "end_turn"}
	}()
	return ch, nil
}

func (p *Provider) Name() string {
	return "mock"
}

func (p *Provider) ModelName() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cfg.Model
}

func (p *Provider) Close() error { return nil }

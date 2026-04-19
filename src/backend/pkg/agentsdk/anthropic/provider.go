package anthropic

import (
	"context"
	"fmt"
	"sync"

	"github.com/develop-agent/backend/pkg/agentsdk"
)

type Provider struct {
	mu    sync.RWMutex
	cfg   agentsdk.Config
	name  string
	model string
}

func New() *Provider {
	return &Provider{name: "anthropic"}
}

func (p *Provider) Initialize(_ context.Context, cfg agentsdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cfg = cfg
	p.model = cfg.Model
	return nil
}

func (p *Provider) Complete(_ context.Context, req agentsdk.CompletionRequest) (agentsdk.CompletionResponse, error) {
	p.mu.RLock()
	model := p.model
	name := p.name
	p.mu.RUnlock()

	if len(req.Messages) == 0 {
		return agentsdk.CompletionResponse{}, fmt.Errorf("%s: at least one message is required", name)
	}
	last := req.Messages[len(req.Messages)-1].Content
	usage := agentsdk.Usage{InputTokens: len(req.Messages) * 12, OutputTokens: 20}
	usage.TotalTokens = usage.InputTokens + usage.OutputTokens

	return agentsdk.CompletionResponse{
		Message: agentsdk.Message{Role: agentsdk.RoleAssistant, Content: fmt.Sprintf("[%s/%s] %s", name, model, last)},
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
	return p.name
}

func (p *Provider) ModelName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.model
}

func (p *Provider) Close() error { return nil }

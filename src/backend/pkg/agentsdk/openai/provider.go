package openai

import (
	"context"
	"errors"
	"io"
	"sync"

	openai "github.com/sashabaranov/go-openai"

	"github.com/develop-agent/backend/pkg/agentsdk"
)

type Provider struct {
	mu     sync.RWMutex
	cfg    agentsdk.Config
	name   string
	model  string
	client *openai.Client
}

func New() *Provider {
	return &Provider{name: "openai"}
}

func (p *Provider) Initialize(_ context.Context, cfg agentsdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cfg = cfg
	p.model = cfg.Model
	p.client = openai.NewClient(cfg.Token)
	return nil
}

func (p *Provider) Complete(ctx context.Context, req agentsdk.CompletionRequest) (agentsdk.CompletionResponse, error) {
	p.mu.RLock()
	model := p.model
	client := p.client
	p.mu.RUnlock()

	if len(req.Messages) == 0 {
		return agentsdk.CompletionResponse{}, errors.New("openai: at least one message is required")
	}

	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: mapMessages(req.Messages),
	}

	resp, err := client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return agentsdk.CompletionResponse{}, err
	}

	if len(resp.Choices) == 0 {
		return agentsdk.CompletionResponse{}, errors.New("openai: empty choices returned")
	}

	return agentsdk.CompletionResponse{
		Message: agentsdk.Message{
			Role:    agentsdk.RoleAssistant,
			Content: resp.Choices[0].Message.Content,
		},
		Usage: agentsdk.Usage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}, nil
}

func (p *Provider) Stream(ctx context.Context, req agentsdk.CompletionRequest) (<-chan agentsdk.StreamEvent, error) {
	p.mu.RLock()
	model := p.model
	client := p.client
	p.mu.RUnlock()

	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: mapMessages(req.Messages),
		Stream:   true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	ch := make(chan agentsdk.StreamEvent)
	go func() {
		defer close(ch)
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					ch <- agentsdk.StreamEvent{Done: true, StopReason: "stop"}
				} else {
					ch <- agentsdk.StreamEvent{Err: err, Done: true}
				}
				return
			}
			if len(resp.Choices) > 0 {
				delta := resp.Choices[0].Delta.Content
				if delta != "" {
					ch <- agentsdk.StreamEvent{Delta: delta}
				}
			}
		}
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

func mapMessages(msgs []agentsdk.Message) []openai.ChatCompletionMessage {
	var result []openai.ChatCompletionMessage
	for _, m := range msgs {
		role := openai.ChatMessageRoleUser
		switch m.Role {
		case agentsdk.RoleSystem:
			role = openai.ChatMessageRoleSystem
		case agentsdk.RoleAssistant:
			role = openai.ChatMessageRoleAssistant
		}
		result = append(result, openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		})
	}
	return result
}

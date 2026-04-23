package gemini

import (
	"context"
	"errors"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/develop-agent/backend/pkg/agentsdk"
)

type Provider struct {
	mu     sync.RWMutex
	cfg    agentsdk.Config
	name   string
	model  string
	client *genai.Client
}

func New() *Provider {
	return &Provider{name: "gemini"}
}

func (p *Provider) Initialize(ctx context.Context, cfg agentsdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cfg = cfg
	p.model = cfg.Model

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.Token))
	if err != nil {
		return err
	}
	p.client = client
	return nil
}

func (p *Provider) Complete(ctx context.Context, req agentsdk.CompletionRequest) (agentsdk.CompletionResponse, error) {
	p.mu.RLock()
	modelName := p.model
	client := p.client
	p.mu.RUnlock()

	if len(req.Messages) == 0 {
		return agentsdk.CompletionResponse{}, errors.New("gemini: at least one message is required")
	}

	model := client.GenerativeModel(modelName)
	applySystemPrompt(model, req.Messages)

	history, lastMsg := splitMessages(req.Messages)

	session := model.StartChat()
	session.History = history

	resp, err := session.SendMessage(ctx, genai.Text(lastMsg))
	if err != nil {
		return agentsdk.CompletionResponse{}, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return agentsdk.CompletionResponse{}, errors.New("gemini: empty response returned")
	}

	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content += string(text)
		}
	}

	usage := agentsdk.Usage{}
	if resp.UsageMetadata != nil {
		usage.InputTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.OutputTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return agentsdk.CompletionResponse{
		Message: agentsdk.Message{
			Role:    agentsdk.RoleAssistant,
			Content: content,
		},
		Usage: usage,
	}, nil
}

func (p *Provider) Stream(ctx context.Context, req agentsdk.CompletionRequest) (<-chan agentsdk.StreamEvent, error) {
	p.mu.RLock()
	modelName := p.model
	client := p.client
	p.mu.RUnlock()

	if len(req.Messages) == 0 {
		return nil, errors.New("gemini: at least one message is required")
	}

	model := client.GenerativeModel(modelName)
	applySystemPrompt(model, req.Messages)

	history, lastMsg := splitMessages(req.Messages)

	session := model.StartChat()
	session.History = history

	stream := session.SendMessageStream(ctx, genai.Text(lastMsg))

	ch := make(chan agentsdk.StreamEvent)
	go func() {
		defer close(ch)

		for {
			resp, err := stream.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					ch <- agentsdk.StreamEvent{Done: true, StopReason: "stop"}
				} else {
					ch <- agentsdk.StreamEvent{Err: err, Done: true}
				}
				return
			}

			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				for _, part := range resp.Candidates[0].Content.Parts {
					if text, ok := part.(genai.Text); ok {
						ch <- agentsdk.StreamEvent{Delta: string(text)}
					}
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

func (p *Provider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil {
		p.client.Close()
	}
	return nil
}

func applySystemPrompt(model *genai.GenerativeModel, msgs []agentsdk.Message) {
	var sysContent string
	for _, m := range msgs {
		if m.Role == agentsdk.RoleSystem {
			sysContent += m.Content + "\n"
		}
	}
	if sysContent != "" {
		model.SystemInstruction = genai.NewUserContent(genai.Text(sysContent))
	}
}

func splitMessages(msgs []agentsdk.Message) ([]*genai.Content, string) {
	var history []*genai.Content
	var lastMsg string

	var userMsgs []agentsdk.Message
	for _, m := range msgs {
		if m.Role != agentsdk.RoleSystem {
			userMsgs = append(userMsgs, m)
		}
	}

	if len(userMsgs) == 0 {
		return history, lastMsg
	}

	lastMsg = userMsgs[len(userMsgs)-1].Content

	for i := 0; i < len(userMsgs)-1; i++ {
		m := userMsgs[i]
		role := "user"
		if m.Role == agentsdk.RoleAssistant {
			role = "model"
		}
		history = append(history, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(m.Content)},
		})
	}

	return history, lastMsg
}

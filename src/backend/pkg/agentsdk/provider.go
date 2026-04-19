package agentsdk

import "context"

// Config holds provider connection settings.
type Config struct {
	Token   string
	Model   string
	BaseURL string
}

// Role identifies message authorship.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a single turn within a conversation.
type Message struct {
	Role       Role
	Content    string
	ToolCallID string
	ToolCalls  []ToolCall
}

// ToolCall represents a provider-requested tool invocation.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

// Tool defines a callable function for the model.
type Tool struct {
	Name        string
	Description string
	Parameters  map[string]any
}

// CompletionRequest contains input for a completion.
type CompletionRequest struct {
	Messages    []Message
	MaxTokens   int
	Temperature float64
	Tools       []Tool
}

// CompletionResponse contains the result for a blocking completion.
type CompletionResponse struct {
	Message    Message
	StopReason string
	Usage      Usage
}

// StreamEvent is emitted during provider streaming.
type StreamEvent struct {
	Delta      string
	ToolCall   *ToolCall
	StopReason string
	Done       bool
	Err        error
}

// Usage records token usage for observability and billing.
type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// Provider is the contract implemented by all LLM providers.
type Provider interface {
	Initialize(ctx context.Context, cfg Config) error
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
	Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error)
	Name() string
	ModelName() string
	Close() error
}

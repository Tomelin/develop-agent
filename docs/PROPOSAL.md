me ajude a aprimorar esse projeto. me de dicas de melhorias que podemos ter nesse projeto.  A proposta é ter uma  Agência de Desenvolvimento de Software com Agentes de Inteligência Artificial.
 Uma esteira completa de produção de software onde cada profissional é um agente especializado, com ciclos de revisão obrigatórios e rastreabilidade total da ideia à entrega. 
 
 1. Para desenvolver um software teremos o seguinte fluxo:
    1. Criação do projeto > 2. Engenharia de software > 3. artquitetura do software > 4. criar as fases e as tarefas dentro de cada fase >5.desenvolvimento do software > 6. testes do sofware > 7. avaliação de segurança > 8. documentação

    O primeiro agente **Criação do projeto** o usuário poderá enviar até 10 feedbacks para melhorias.
    para as demais etapas o usuário poderá enviar até 5 feedbacks para melhorias.
    
    Esse será o fluxo para o desenvolvimento do software, sendo que a partir da etapa de engenharia, será separado por frontend e backend

    Nesse fluxo serão 7 grupo de agente, uma para cada etapa.

    e dentro de cada etapa, teremos mais 3, que serão:
    1. o agente que irá desenvolver a etapa inicial
    2. o agente que irá revisar e dar feedbacks
    3. o agente ajustar conforme o feedback e desenvolver e entregar essa etapa

2. Teremos a parte para desenvolvimento de uma landing page.  o usuário poderá vincular o projeto da etapa 1, para ser como base do prompt para o desenvolvimento da landing page.

    Teremos o mesmo padrão de 3 agentes, produtor, revisor e o agente que irá ajustar conforme o feedback e desenvolver e entregar essa etapa.
    

3. A terceira parte de desenvolvimento, é para desenvolver estrategia de marketing, que irá criar campanhas para o linkedin, instagram, google ads, etc.  o usuário poderá vincular o projeto da etapa 1, para ser como base do prompt para o desenvolvimento da estrategia de marketing.

    Teremos o mesmo padrão de 3 agentes, produtor, revisor e o agente que irá ajustar conforme o feedback e desenvolver e entregar essa etapa.
    

    O usuário terá o perfil dele.  Vamos iniciar, apenas um usuário padrão, chamado admin.  e no profile desse usuário, pode ter multiplos prompts para cada grupo de agent.
    A ideia desses multiplos promtps, é serem agregados para o desenvolvimento de cada etapa.  Pense que o usuário pode ter um padrão de cor de uma plataforma, ou landing page.  pode definir linguagem de programação, ....  A ideia é que esses multiplos promtps por fase, sejam agregados para o desenvolvimento de cada etapa.
    

     Entrevistador - **Objetivo:** Transformar uma ideia bruta em visão clara de produto. O agente opera de forma conversacional: faz perguntas, escuta, reformula o que entendeu e pede confirmação. Não avança até ter um entendimento sólido
     
 Teremos uma lista de agentes, e cada agente terá um **nome**, **descrição**, **modelo de AI**, **lista de prompts** e **lista de tarefas**.    e terá uma opção de enabled, se enabled estará habilitado para uso.  O usuário pode esolher irá usar em cada etapa, ou pode escolher dinamico, que irá usar um modelo diferente a cada fase. quando usa moddelo dinamico, o sistema irá escolher aleatoriamente um modelo habilitado para executar a fase.
     
quando usuário tiver a lista de todas as fases com as atividades, o usuário pode escolher, se quer mandar desenvolver toda a fase de uma única vez, ou se será desenvolvido tarefa por tarefa.  isso o usuário irá fazer via portal ou request da API.  Se for tarefa por tarefa, o usuário poderá escolher qual agente irá desenvolver a tarefa, e poderá escolher se irá usar um modelo diferente para cada tarefa.  Se for toda a fase de uma única vez, o sistema irá escolher um modelo para cada etapa, e irá executar a fase completa.


atualize onde for necessário, s eno PROJECT.md, PLAYBOOK.md, em phases e nas tasks que foram necessárias.  
comunicação entre agentes.  A comunicação entre os agentes, deve ser num formato JSON, onde tem o agente de origem, agente de destino, e a mensagem e o status.
Essa comunicação deve ser através de channels.  Por padrão o channel terá um tamanho de 10, mas isso é configuravel, através das configurações.
Os channels serão criados automaticamente, quando um agente for criado.  E serão destruidos automaticamente, quando o agente for destruido.
O dashboard, deve ver o status de cada agente, se está ativo, se está em execução, se está em pausa, se tem mensagem na fila para ser executados, se está em erro, se está em concluído




Para o core business da API que irá conectar com os modelos de AI, quero que seja via interface. 

Pode consultar a estrutura **/home/carlos.tomelin/codes/genai/agencia-desenvolvimento/src/backend/pkg/agentsdk**, para entender um modelo de interface interessante.  E no desenvolvimento da phase desse tópico, especifique que será um modelo de interface, colocando o exemplo de como deve ser.

Esse é o modelo de uma interface
```go
// Package agentsdk defines the common contract and shared types for all
// LLM provider integrations used by the agency's agent layer.
//
// The central abstraction is the [Provider] interface. Every backend
// (OpenAI, Anthropic, Gemini, Ollama, …) must implement this interface so
// that the core layer can drive any provider without coupling to its SDK.
//
// Typical usage:
//
//	var p agentsdk.Provider = openai.New()
//	if err := p.Initialize(ctx, cfg); err != nil { ... }
//	defer p.Close()
//
//	resp, err := p.Complete(ctx, agentsdk.CompletionRequest{
//	    Messages: []agentsdk.Message{
//	        {Role: agentsdk.RoleUser, Content: "Hello!"},
//	    },
//	})
package agentsdk

import "context"

// ─── Configuration ──────────────────────────────────────────────────────────

// Config holds the connection parameters required to initialise an LLM provider.
// It is intentionally flat so it can be populated from any configuration source
// (YAML, environment, Vault, etc.) without coupling this package to the app config.
type Config struct {
	// Token is the API key or bearer token used to authenticate with the provider.
	Token string

	// Model is the identifier of the LLM model to use (e.g. "gpt-4o", "claude-3-5-sonnet").
	Model string

	// BaseURL allows overriding the provider's default API endpoint.
	// It is required for self-hosted providers such as Ollama.
	BaseURL string
}

// ─── Conversation primitives ─────────────────────────────────────────────────

// Role identifies the participant that authored a conversation [Message].
type Role string

const (
	// RoleSystem carries the system prompt that sets the model's behaviour and persona.
	RoleSystem Role = "system"

	// RoleUser represents input coming from the human (or orchestrating agent).
	RoleUser Role = "user"

	// RoleAssistant represents a response previously generated by the model.
	RoleAssistant Role = "assistant"

	// RoleTool carries the result of a tool execution back to the model.
	// When this role is used, Message.ToolCallID must reference the
	// originating ToolCall.ID so the model can correlate the result.
	RoleTool Role = "tool"
)

// Message represents a single turn in a conversation thread.
//
// For tool-result messages (Role == RoleTool), ToolCallID must be set.
// For assistant messages that request tool executions, ToolCalls must be set
// and Content may be empty.
type Message struct {
	// Role identifies who authored this message.
	Role Role

	// Content is the text body of the message.
	Content string

	// ToolCallID links this message to a prior ToolCall.ID.
	// Only populated when Role == RoleTool.
	ToolCallID string

	// ToolCalls contains the tool invocations requested by the model.
	// Only populated when Role == RoleAssistant and the model chose to call tools.
	ToolCalls []ToolCall
}

// ─── Tool-calling ────────────────────────────────────────────────────────────

// ToolCall represents a single tool invocation requested by the model within
// an assistant message. The caller must execute the tool and return the result
// as a follow-up Message with Role == RoleTool and the matching ToolCallID.
type ToolCall struct {
	// ID is a unique, provider-assigned identifier for this specific invocation.
	// It must be echoed back in the corresponding tool-result Message.ToolCallID.
	ID string

	// Name is the name of the tool to invoke, as declared in Tool.Name.
	Name string

	// Arguments is a JSON-encoded object whose fields match the tool's parameter schema.
	Arguments string
}

// Tool defines a callable capability that can be offered to the model.
// Parameters must be a valid JSON Schema object (type "object" with "properties").
// Providers translate this into their native function-calling format.
type Tool struct {
	// Name is the identifier the model uses to refer to this tool.
	// Must be unique within a single CompletionRequest.
	Name string

	// Description explains the tool's purpose to the model.
	// Clear, concise descriptions improve the model's tool-selection accuracy.
	Description string

	// Parameters is a JSON Schema object that describes the tool's input shape.
	// Example: map[string]any{"type": "object", "properties": { ... }}
	Parameters map[string]any
}

// ─── Request / Response ───────────────────────────────────────────────────────

// CompletionRequest holds all inputs for a single completion round-trip.
//
// Messages must contain the full ordered conversation history, ending with the
// turn that should be responded to. An optional system prompt should be placed
// as the first message with Role == RoleSystem.
type CompletionRequest struct {
	// Messages is the ordered conversation history for this request.
	Messages []Message

	// MaxTokens caps the number of tokens the model may generate.
	// Zero means the provider's own default is applied.
	MaxTokens int

	// Temperature controls the randomness of the output (0.0 = deterministic,
	// 1.0 = creative). Zero means the provider's own default is applied.
	Temperature float64

	// Tools declares the callable capabilities available to the model for this
	// request. When non-empty, the model may respond with tool-call requests
	// instead of (or in addition to) plain text.
	Tools []Tool
}

// CompletionResponse holds the result of a completed, non-streaming call.
type CompletionResponse struct {
	// Message is the assistant's reply.
	// When the model requested tool executions, Message.ToolCalls will be
	// populated and Message.Content may be empty.
	Message Message

	// StopReason is the provider-specific string that explains why generation
	// stopped (e.g. "end_turn", "tool_use", "max_tokens", "stop").
	StopReason string

	// Usage records the token consumption for this round-trip.
	Usage Usage
}

// ─── Streaming ───────────────────────────────────────────────────────────────

// StreamEvent is a single incremental event emitted during a streaming call.
//
// Text tokens arrive as non-empty Delta values. Tool-call fragments arrive
// with a non-nil ToolCall pointer. The stream ends when Done is true; if an
// error occurred, Err will also be non-nil. The channel is closed after this
// final event.
type StreamEvent struct {
	// Delta holds the incremental text content for this event.
	// Empty when the event carries a ToolCall or is the final Done event.
	Delta string

	// ToolCall is non-nil when the model is streaming a tool invocation.
	// Multiple events may be needed to stream a single complete ToolCall;
	// the provider implementation must reassemble them before sending.
	ToolCall *ToolCall

	// StopReason is populated only in the final event (Done == true).
	StopReason string

	// Done is true in the last event of the stream.
	// The caller should stop reading from the channel after receiving it.
	Done bool

	// Err is non-nil if a streaming error occurred.
	// When Err is set, Done is also true and the channel will be closed.
	Err error
}

// ─── Observability ───────────────────────────────────────────────────────────

// Usage records token consumption for a completion round-trip.
// It exists to support governance, cost tracking, and observability across providers.
type Usage struct {
	// InputTokens is the number of tokens in the request (prompt + history).
	InputTokens int

	// OutputTokens is the number of tokens generated in the response.
	OutputTokens int

	// TotalTokens is InputTokens + OutputTokens.
	TotalTokens int
}

// ─── Provider interface ───────────────────────────────────────────────────────

// Provider is the contract that every agent SDK backend must satisfy.
// It abstracts over concrete LLM APIs (OpenAI, Anthropic, Gemini, Ollama) so
// the core layer can drive any provider transparently.
//
// Lifecycle: Initialize → Complete / Stream (any number of calls) → Close.
//
// Implementations must be safe for concurrent use after Initialize returns.
type Provider interface {
	// Initialize configures the provider using the supplied Config.
	// It must be called once before any other method.
	// Calling Initialize again on an already-initialised provider must be
	// safe and should refresh the configuration in-place (e.g. token rotation).
	Initialize(ctx context.Context, cfg Config) error

	// Complete performs a blocking, non-streaming completion round-trip.
	// It sends all messages in req to the provider and waits for the full
	// response before returning. Use Stream for incremental delivery.
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

	// Stream begins a streaming completion and returns a read-only event channel.
	// The caller must consume the channel until a StreamEvent where Done == true
	// (or Err != nil) is received; the channel is closed after that event.
	// Abandoning the channel without draining it may leak goroutines inside the
	// provider implementation — always drain or cancel the context.
	Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error)

	// Name returns the canonical lower-case identifier of this provider
	// (e.g. "openai", "anthropic", "gemini", "ollama").
	// It is used by the infra layer for routing, logging, and metrics.
	Name() string

	// ModelName returns the identifier of the model currently configured
	// for this provider instance (e.g. "gpt-4o", "claude-3-5-sonnet-20241022").
	ModelName() string

	// Close releases all resources held by the provider (HTTP connections,
	// background goroutines, etc.). It must be idempotent — calling it more
	// than once must not return an error or panic.
	Close() error
}
```
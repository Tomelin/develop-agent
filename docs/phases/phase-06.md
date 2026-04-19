# Phase 06 — Pipeline de Execução da Tríade de Agentes

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-06 |
| **Título** | Pipeline de Execução da Tríade de Agentes |
| **Tipo** | Backend |
| **Prioridade** | Crítica |
| **Pré-requisitos** | PHASE-01, PHASE-02, PHASE-03, PHASE-04, PHASE-05 concluídas |

---

## Descrição Detalhada

Esta é a fase técnica mais crítica de toda a plataforma. Aqui é implementado o motor de orquestração da **Tríade de Agentes** e o **contrato de interface `agentsdk.Provider`** — a abstração central que permite que qualquer modelo de IA (OpenAI, Anthropic, Gemini, Ollama) seja usado de forma intercambiável pelo core da plataforma.

O **design arquitetural obrigatório** para integração com LLMs é baseado em interface (`agentsdk.Provider`). Nenhum código de domínio ou orquestração pode depender diretamente de um SDK específico. A interface define o contrato de ciclo de vida (`Initialize → Complete/Stream → Close`) e cada provider implementa esse contrato de forma isolada.

O pipeline de execução é assíncrono: quando uma fase é iniciada, o sistema publica uma mensagem no RabbitMQ e um worker consome essa mensagem, executando a Tríade completa em background. O usuário pode acompanhar o progresso em tempo real via SSE (Server-Sent Events). Cada etapa da Tríade gera eventos que são emitidos para o frontend atualizar o status em tempo real.

O motor de execução também gerencia o fluxo de feedback: após a Tríade concluir, o sistema aguarda o feedback do usuário (dentro do limite configurado por fase) e permite que o Refinador aplique melhorias adicionais com base nas sugestões.

Esta fase é o núcleo que, uma vez implementado corretamente, é reutilizado por TODOS os fluxos e fases da plataforma.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Pacote `agentsdk` com a interface `Provider` e todos os tipos (contrato único para todos os LLMs)
- ✅ Implementações concretas da interface para OpenAI, Anthropic, Gemini e Ollama
- ✅ `MockProvider` determinístico para testes unitários e de integração sem custo de API
- ✅ `AgentMessage` — envelope JSON padronizado para toda comunicação entre agentes
- ✅ Sistema de channels com criacao/destruição automática e buffer configurável (padrão: 10)
- ✅ Status de agente (`IDLE | RUNNING | PAUSED | QUEUED | ERROR | COMPLETED`) exibido no dashboard
- ✅ Worker assíncrono que executa a Tríade completa (Produtor → Revisor → Refinador)
- ✅ Producer e Consumer de mensagens RabbitMQ para o pipeline
- ✅ Sistema de eventos em tempo real (SSE) para progresso da Tríade
- ✅ Fluxo de feedback do usuário integrado ao ciclo da Tríade
- ✅ Rastreamento de tokens e custo por execução
- ✅ Sistema de retry para falhas transitórias de LLM

---

## Funcionalidades Entregues

- **Interface `agentsdk.Provider`:** Contrato único que abstrai todos os providers de LLM
- **Providers Concretos:** OpenAI, Anthropic, Gemini, Ollama — todos impl. via interface
- **`AgentMessage`:** Envelope JSON padronizado para toda comunicação entre agentes
- **Channels com Ciclo de Vida Automático:** Criados/destruídos junto com o agente; buffer configurável
- **Dashboard de Status de Agentes:** Visão em tempo real de cada agente com seus 6 estados possíveis
- **Execução Assíncrona:** Pipeline não bloqueia a API; executa em workers background
- **Eventos em Tempo Real:** Frontend acompanha cada step da Tríade ao vivo
- **Ciclo de Feedback:** Usuário pode enviar feedback e iniciar novo ciclo da Tríade
- **Billing por Execução:** Cada chamada de LLM registrada com tokens e custo
- **Retry Inteligente:** Falhas transitórias tratadas sem intervenção manual

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Todas as tasks são executadas sequencialmente pela Tríade de Agentes (Produtor → Revisor → Refinador), sem interrupções entre elas. O sistema avança automaticamente de uma task para a próxima até concluir a phase inteira e aguarda uma única aprovação ao final.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a phase |
| **Velocidade** | Mais rápido — execução contínua sem pausas |
| **Feedback** | Aplicado à phase como um todo |
| **Ideal para** | Phases bem compreendidas onde o usuário confia na execução automática |

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. A Tríade desenvolve apenas a(s) task(s) escolhida(s) e aguarda aprovação explícita antes de prosseguir para a próxima.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — o usuário controla o ritmo |
| **Velocidade** | Mais controlado — requer interação entre tasks |
| **Feedback** | Granular e específico para cada task |
| **Ideal para** | Tasks críticas ou complexas que exigem revisão antes de avançar |

### 🔀 Modo Híbrido

É possível **combinar os dois modos**: inicie a phase automaticamente e pause manualmente em qualquer task que exija atenção especial. Após aprovar aquela task individualmente, a execução automática retoma a partir da próxima task.

> 💡 **Dica:** Para esta phase, o modo recomendado é **task a task** — o TriadOrchestrator é o motor central da plataforma. Cada componente (workers, SSE, billing) deve ser revisado individualmente antes de prosseguir, pois erros aqui afetam toda a esteira.

---

## Tasks

### TASK-06-000 — AgentSDK: Definição da Interface `Provider` (**Fundação do Core**)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o pacote `agentsdk` com a interface `Provider` como contrato único e obrigatório para toda integração com modelos de IA |

**Descrição:**
Criar o pacote `src/backend/pkg/agentsdk/provider.go` com todos os tipos e a interface `Provider`. **Este arquivo é a fundação de toda a integração com LLMs** — nenhuma outra camada (`domain`, `usecase`, `worker`) pode importar SDKs concretos diretamente. O pacote define:

- **`Config`** — parâmetros de conexão (Token, Model, BaseURL)
- **`Role`** + constantes (`RoleSystem`, `RoleUser`, `RoleAssistant`, `RoleTool`)
- **`Message`** — turno de conversa com suporte a tool calling
- **`ToolCall`** e **`Tool`** — suporte a function calling
- **`CompletionRequest`** e **`CompletionResponse`**
- **`StreamEvent`** — evento incremental de streaming
- **`Usage`** — rastreamento de tokens (input, output, total)
- **`Provider`** — a interface central

**Código de Referência Obrigatório:**
```go
// Package agentsdk defines the common contract and shared types for all
// LLM provider integrations used by the agency's agent layer.
//
// The central abstraction is the [Provider] interface. Every backend
// (OpenAI, Anthropic, Gemini, Ollama, ...) must implement this interface so
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

// Config holds the connection parameters required to initialise an LLM provider.
// Intentionally flat so it can be populated from any config source (YAML, env, Vault).
type Config struct {
	Token   string // API key or bearer token
	Model   string // Model identifier (e.g. "gpt-4o", "claude-3-5-sonnet")
	BaseURL string // Override for self-hosted providers like Ollama
}

// Role identifies the participant that authored a conversation Message.
type Role string

const (
	RoleSystem    Role = "system"    // System prompt — sets model behaviour and persona
	RoleUser      Role = "user"      // Input from the human or orchestrating agent
	RoleAssistant Role = "assistant" // Response previously generated by the model
	RoleTool      Role = "tool"      // Result of a tool execution fed back to the model
)

// Message represents a single turn in a conversation thread.
type Message struct {
	Role       Role
	Content    string
	ToolCallID string     // Set when Role == RoleTool
	ToolCalls  []ToolCall // Set when Role == RoleAssistant and model called tools
}

// ToolCall represents a single tool invocation requested by the model.
type ToolCall struct {
	ID        string // Provider-assigned unique ID for this invocation
	Name      string // Tool name as declared in Tool.Name
	Arguments string // JSON-encoded input arguments
}

// Tool defines a callable capability offered to the model (function calling).
type Tool struct {
	Name        string         // Unique identifier used by the model
	Description string         // Purpose description — improves model's tool-selection accuracy
	Parameters  map[string]any // JSON Schema object describing the input shape
}

// CompletionRequest holds all inputs for a single completion round-trip.
type CompletionRequest struct {
	Messages    []Message // Full ordered conversation history
	MaxTokens   int       // Cap on generated tokens; 0 = provider default
	Temperature float64   // Randomness (0.0 deterministic, 1.0 creative); 0 = provider default
	Tools       []Tool    // Callable capabilities available for this request
}

// CompletionResponse holds the result of a completed, non-streaming call.
type CompletionResponse struct {
	Message    Message // Assistant's reply (may have ToolCalls instead of Content)
	StopReason string  // Provider-specific stop reason (e.g. "end_turn", "max_tokens")
	Usage      Usage   // Token consumption for billing
}

// StreamEvent is a single incremental event emitted during a streaming call.
type StreamEvent struct {
	Delta      string    // Incremental text token; empty for tool/done events
	ToolCall   *ToolCall // Non-nil when model is streaming a tool invocation
	StopReason string    // Populated only in the final Done event
	Done       bool      // True in the last event; channel closed after this
	Err        error     // Non-nil if streaming error occurred
}

// Usage records token consumption for billing and observability.
type Usage struct {
	InputTokens  int // Tokens in the request (prompt + history)
	OutputTokens int // Tokens generated in the response
	TotalTokens  int // InputTokens + OutputTokens
}

// Provider is the contract that every agent SDK backend must satisfy.
// It abstracts over OpenAI, Anthropic, Gemini, Ollama so the core layer
// can drive any provider transparently.
//
// Lifecycle: Initialize → Complete / Stream (any number of calls) → Close.
// Implementations MUST be safe for concurrent use after Initialize returns.
type Provider interface {
	// Initialize configures the provider. Must be called once before any
	// other method. Safe to call again for config refresh (e.g. token rotation).
	Initialize(ctx context.Context, cfg Config) error

	// Complete performs a blocking, non-streaming completion round-trip.
	// Use Stream for incremental delivery.
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

	// Stream begins a streaming completion and returns a read-only event channel.
	// The caller MUST consume the channel until Done == true or Err != nil.
	// Abandoning without draining may leak goroutines — cancel the context instead.
	Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error)

	// Name returns the canonical lower-case provider identifier
	// (e.g. "openai", "anthropic", "gemini", "ollama").
	// Used for routing, logging, and metrics.
	Name() string

	// ModelName returns the model currently configured for this instance
	// (e.g. "gpt-4o", "claude-3-5-sonnet-20241022").
	ModelName() string

	// Close releases all resources. Must be idempotent — safe to call multiple times.
	Close() error
}
```

**Estrutura de pacotes criada por esta task:**

```
src/backend/pkg/agentsdk/
├── provider.go          ← Este arquivo (interface + tipos)
├── openai/
│   └── provider.go      ← Implementação OpenAI (TASK-06-000a)
├── anthropic/
│   └── provider.go      ← Implementação Anthropic (TASK-06-000b)
├── gemini/
│   └── provider.go      ← Implementação Gemini (TASK-06-000c)
└── ollama/
    └── provider.go      ← Implementação Ollama (TASK-06-000d)
```

**Regra arquitetural — O que é permitido e proibido:**

```go
// ✅ CORRETO — domínio depende da abstração
var producer agentsdk.Provider = anthropic.New()
producer.Initialize(ctx, agentsdk.Config{Token: os.Getenv("ANTHROPIC_KEY"), Model: "claude-3-5-sonnet"})
defer producer.Close()

resp, err := producer.Complete(ctx, agentsdk.CompletionRequest{
    Messages: []agentsdk.Message{
        {Role: agentsdk.RoleSystem, Content: systemPrompt},
        {Role: agentsdk.RoleUser, Content: userInput},
    },
    Temperature: 0.7,
})

// ❌ PROIBIDO — nunca importe SDKs concretos fora de pkg/agentsdk/*
import sdk "github.com/anthropics/anthropic-sdk-go"
client := sdk.NewClient(...) // viola a separação de camadas
```

**Critério de aceite:** Interface `Provider` definida em `pkg/agentsdk/provider.go`; 4 providers concretos implementados e testados; `MockProvider` determinístico para testes; zero import de SDKs concretos fora de `pkg/agentsdk/*`; documentação godoc completa.

---

### TASK-06-00A — Protocolo de Comunicação entre Agentes (`AgentMessage` + Channels)

| Campo | Valor |
|-------| ------|
| **Camada** | Backend |
| **Objetivo** | Implementar o envelope JSON padronizado e a infraestrutura de Go channels para toda a comunicação entre agentes |

**Descrição:**
Criar o pacote `src/backend/domain/agent/channel.go` com os tipos e a lógica de ciclo de vida dos channels. **A comunicação direta entre agentes é proibida** — todo dado transferido deve usar o `AgentMessage` via channel.

**Tipos a implementar:**

```go
// AgentStatus representa o estado operacional atual de um agente.
// Exibido no dashboard em tempo real.
type AgentStatus string

const (
    AgentStatusIdle      AgentStatus = "IDLE"      // Ativo, aguardando mensagens
    AgentStatusRunning   AgentStatus = "RUNNING"   // Processando uma mensagem agora
    AgentStatusPaused    AgentStatus = "PAUSED"    // Suspenso pelo orquestrador
    AgentStatusQueued    AgentStatus = "QUEUED"    // Há mensagens no buffer pendentes
    AgentStatusError     AgentStatus = "ERROR"     // Falha — requer atenção
    AgentStatusCompleted AgentStatus = "COMPLETED" // Tarefa concluída
)

// AgentMessage é o envelope padrão para toda comunicação entre agentes.
// Transportado exclusivamente via channels — nunca via chamada direta.
type AgentMessage struct {
    ID        string         `json:"id"`             // UUID único por mensagem
    From      string         `json:"from"`           // ID do agente de origem (e.g. "producer-abc123")
    To        string         `json:"to"`             // ID do agente de destino (e.g. "reviewer-xyz456")
    Message   string         `json:"message"`        // Payload (artefato, crítica, refinamento)
    Status    string         `json:"status"`         // "pending" | "processing" | "done" | "error"
    Timestamp time.Time      `json:"timestamp"`      // Momento de criação
    Meta      map[string]any `json:"meta,omitempty"` // Metadados opcionais (phase, token_usage, etc.)
}

// AgentChannel encapsula os canais de entrada e saída de um agente.
// O buffer é configurado via agent.channel_buffer_size (padrão: 10).
type AgentChannel struct {
    In  <-chan AgentMessage // Canal de entrada — mensagens recebidas pelo agente
    Out chan<- AgentMessage // Canal de saída — mensagens enviadas pelo agente
    buf int                // Tamanho do buffer (configurável, não exportado)
}

// NewAgentChannel cria um novo AgentChannel com o buffer especificado.
// Chamado automaticamente pelo orquestrador ao instanciar um agente.
func NewAgentChannel(bufferSize int) (in chan AgentMessage, out chan AgentMessage) {
    ch := make(chan AgentMessage, bufferSize) // canal bidirecional interno
    return ch, ch // in e out apontam para o mesmo channel — separação de responsabilidade via tipos
}
```

**Regras de ciclo de vida:**

| Evento | Comportamento |
|--------|---------------|
| Agente instanciado | `NewAgentChannel(cfg.Agent.ChannelBufferSize)` chamado automaticamente |
| Agente encerrado | `close(ch)` chamado automaticamente pelo orquestrador |
| Buffer cheio | Remetente bloqueia até haver espaço (comportamento natural do channel Go) |
| Agente em `ERROR` | Channel preservado — mensagens pendentes aguardam recuperação |

**Configuração necessária em `config.yaml`:**

```yaml
agent:
  channel_buffer_size: 10   # Tamanho do buffer padrão (configurável)
  channel_drain_timeout: 30  # Segundos para drenar mensagens antes de destruir
```

**Critério de aceite:** `AgentStatus` enum com 6 estados; `AgentMessage` e `AgentChannel` definidos e documentados com godoc; `NewAgentChannel` cria channel com buffer configurado; channel fechado ao encerrar o agente; configuração mapeada no `Config`.

---

### TASK-06-00B — Monitor de Status de Agentes no Dashboard

| Campo | Valor |
|-------| ------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Expor e exibir em tempo real o status operacional de cada agente (`IDLE | RUNNING | PAUSED | QUEUED | ERROR | COMPLETED`) no dashboard |

**Descrição:**

**Backend:** Implementar endpoint `GET /api/v1/agents/status` que retorna o status operacional atual de todos os agentes ativos. Implementar endpoint SSE `GET /api/v1/agents/status/stream` que emite eventos ao vivo sempre que o status de qualquer agente mudar:

```json
// Evento SSE emitido ao mudar status
{
  "agent_id": "producer-abc123",
  "agent_name": "Arquiteto Senior (Anthropic)",
  "status": "RUNNING",
  "queue_size": 0,
  "current_task": "TASK-06-002",
  "timestamp": "2026-04-19T14:30:00Z"
}
```

**Regras de atualização de status:**

| Status | Quando ocorre |
|--------|---------------|
| `IDLE` | Agente criado ou finalizou última mensagem sem novas no buffer |
| `RUNNING` | Iniciou o processamento de uma `AgentMessage` |
| `PAUSED` | Orquestrador emitiu sinal de pausa |
| `QUEUED` | `len(channel) > 0` e agente ainda não iniciou a próxima mensagem |
| `ERROR` | Processamento falhou (após retries esgotados) |
| `COMPLETED` | Tarefa da Tríade concluída — aguarda encerramento |

**Frontend:** Implementar o painel `AgentStatusPanel` no dashboard (`/dashboard`) com:
- Lista de todos os agentes com seus status em tempo real (conectado ao SSE)
- Indicador visual por status: `IDLE` ● cinza, `RUNNING` ● verde pulsante, `PAUSED` ● amarelo, `QUEUED` ● azul, `ERROR` ● vermelho, `COMPLETED` ● verde sólido
- Contador de mensagens na fila para agentes `QUEUED`
- Nome da task atual para agentes `RUNNING`
- Fallback para polling de 5 segundos se SSE não estiver disponível

**Critério de aceite:** Endpoint `/api/v1/agents/status` retorna status correto; SSE emite eventos ao mudar de status; dashboard exibe status em tempo real com indicadores visuais corretos por estado.

---

### TASK-06-001 — Modelo de Execução da Tríade (TriadExecution)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir as structs de controle de execução da Tríade com rastreabilidade completa |

**Descrição:**
Definir as structs de domínio: `TriadExecution` com `ID`, `ProjectID`, `PhaseNumber`, `PhaseTrack` (FULL/FRONTEND/BACKEND), `Status` (PENDING/PRODUCER_RUNNING/REVIEWER_RUNNING/REFINER_RUNNING/AWAITING_FEEDBACK/COMPLETED/FAILED), `Producer` (ref ao agente), `Reviewer` (ref ao agente), `Refiner` (ref ao agente), `ProducerOutput` (artefato gerado pelo Produtor), `ReviewerOutput` (lista de críticas do Revisor), `RefinerOutput` (artefato final do Refinador), `FeedbackHistory` (slice de `UserFeedback{Content, Timestamp}`), `FeedbackCount`, `FeedbackLimit`, `TokenUsage` (por step e total), `StartedAt`, `CompletedAt`. Struct separada `AgentStep` para cada execução individual de agente com `AgentID`, `Model`, `PromptTokens`, `CompletionTokens`, `DurationMs`, `Output`, `Timestamp`.

**Critério de aceite:** Structs com rastreabilidade completa; TokenUsage granular por step; histórico de feedbacks.

---

### TASK-06-002 — Orquestrador da Tríade (TriadOrchestrator)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o motor central que coordena a execução sequencial dos três agentes |

**Descrição:**
Implementar `TriadOrchestrator` em `src/backend/domain/triad/orchestrator.go` com o método `Execute(ctx, TriadExecution) error`. O orquestrador deve: **(1) Executa o Produtor:** compõe o prompt (via PromptAggregator), chama o agente Produtor via AgentSDK, armazena o output, emite evento `PRODUCER_COMPLETED`. **(2) Executa o Revisor:** injeta o output do Produtor + as regras de revisão crítica no prompt do Revisor, chama o agente Revisor, armazena as críticas estruturadas, emite evento `REVIEWER_COMPLETED`. **(3) Executa o Refinador:** injeta o output do Produtor + as críticas do Revisor no prompt do Refinador, chama o agente Refinador, armazena o artefato final, emite evento `REFINER_COMPLETED`, muda status para `AWAITING_FEEDBACK`. Cada step com timeout configurável (default 5 minutos).

**Critério de aceite:** Sequência correta Produtor→Revisor→Refinador; eventos emitidos em cada step; timeout respeitado.

---

### TASK-06-003 — Worker Assíncrono de Execução da Tríade

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o worker que consome mensagens da fila e executa a Tríade em background |

**Descrição:**
Implementar o `TriadWorker` em `src/backend/worker/triad_worker.go` que: registra-se como consumer da fila RabbitMQ `phase.execution`, ao receber uma mensagem, deserializa o `TriadExecution`, executa o `TriadOrchestrator`, em caso de sucesso, dá `ack` na mensagem e persiste o resultado, em caso de falha transitória (timeout, rate limit do LLM), dá `nack` para reprocessamento (com máximo de 3 tentativas), em caso de falha permanente, move para a Dead Letter Queue e emite evento de falha para o usuário. Implementar múltiplas instâncias do worker (configurável via `WORKER_CONCURRENCY`) para processamento paralelo de projetos diferentes.

**Critério de aceite:** Worker consome e processa mensagens; retry em falhas transitórias; DLQ para falhas permanentes; paralelismo configurável.

---

### TASK-06-004 — Handler de Inicialização de Fase

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o endpoint que inicia a execução de uma fase do pipeline |

**Descrição:**
Implementar `POST /api/v1/projects/:id/phases/:phaseNumber/start` que: valida que o projeto existe e pertence ao usuário, valida que a fase pode ser iniciada (fase anterior concluída, nenhuma fase em execução no momento), seleciona a Tríade de agentes (via AgentSelectorService, respeitando Modo Dinâmico), cria o registro de `TriadExecution` no banco, publica mensagem na fila RabbitMQ `phase.execution`, retorna imediatamente com `202 Accepted` e o ID da execução para polling/SSE. Não bloqueia esperando a conclusão da Tríade.

**Critério de aceite:** Retorno imediato `202 Accepted`; mensagem publicada na fila; Tríade selecionada e registrada.

---

### TASK-06-005 — Handler de Submissão de Feedback do Usuário

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o endpoint de feedback que reinicia o ciclo da Tríade com as sugestões do usuário |

**Descrição:**
Implementar `POST /api/v1/projects/:id/phases/:phaseNumber/feedback` que: valida que a fase está no status `AWAITING_FEEDBACK`, valida que o limite de feedbacks não foi atingido (`FeedbackCount < FeedbackLimit`), registra o feedback no histórico (`FeedbackHistory`), incrementa `FeedbackCount`, cria nova `TriadExecution` com o histórico de feedback injetado no contexto, publica nova mensagem na fila para reexecução, retorna `202 Accepted`. Se o limite de feedbacks foi atingido, retorna `400 Bad Request` com mensagem informativa.

**Critério de aceite:** Feedback salvo no histórico; limite de feedbacks validado; reexecução da Tríade iniciada; erro claro quando limite atingido.

---

### TASK-06-006 — Handler de Aprovação de Fase

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o endpoint de aprovação que avança o projeto para a próxima fase |

**Descrição:**
Implementar `POST /api/v1/projects/:id/phases/:phaseNumber/approve` que: valida que a fase está em `AWAITING_FEEDBACK` ou que toda a Tríade foi concluída, muda o status da fase para `COMPLETED`, gera o `SPEC.md` da fase (via um LLM leve — gera um resumo do artefato do Refinador), armazena o SPEC.md no projeto (como contexto RAG para a próxima fase), avança `CurrentPhaseNumber` do projeto, emite evento `PHASE_COMPLETED` via SSE, se for a última fase, muda o status do projeto para `COMPLETED`. Retorna o estado atualizado do projeto.

**Critério de aceite:** Status da fase atualizado; SPEC.md gerado e armazenado; projeto avança para próxima fase; evento emitido.

---

### TASK-06-007 — Sistema de Eventos em Tempo Real (SSE para Progresso da Tríade)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar SSE para que o frontend acompanhe o progresso da Tríade em tempo real |

**Descrição:**
Implementar `GET /api/v1/projects/:id/stream` como endpoint SSE. O backend mantém um canal em memória por projeto (usando Redis Pub/Sub para suportar múltiplas instâncias do servidor). Os eventos emitidos incluem: `PHASE_STARTED {phaseNumber, agentTriad}`, `PRODUCER_STARTED {agentName, model}`, `PRODUCER_COMPLETED {tokenUsage, preview}`, `REVIEWER_STARTED {agentName}`, `REVIEWER_COMPLETED {issuesCount, preview}`, `REFINER_STARTED {agentName}`, `REFINER_COMPLETED {artifactPreview}`, `AWAITING_FEEDBACK {feedbackCount, feedbackLimit}`, `PHASE_COMPLETED`, `PHASE_FAILED {reason}`. Cada evento inclui timestamp e o payload relevante.

**Critério de aceite:** SSE emite todos os eventos da Tríade; Redis Pub/Sub garante funcionamento multi-instância; frontend recebe eventos corretos.

---

### TASK-06-008 — Gerador de SPEC.md (Contexto RAG entre Fases)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o gerador de manifesto de contexto que comprime o output de cada fase para alimentar as seguintes |

**Descrição:**
Implementar o `SpecGenerator` em `src/backend/domain/spec/generator.go`. Ao final de cada fase, o gerador recebe o artefato do Refinador e usa um modelo LLM leve (configurável, preferencialmente Gemini Flash ou GPT-4o-mini por custo/velocidade) para gerar um `SPEC.md` comprimido contendo: resumo executivo dos artefatos gerados, decisões técnicas mais importantes, tecnologias e padrões definidos, pendências ou considerações para fases futuras. O SPEC.md não deve exceder 2000 tokens. Este manifesto é armazenado no projeto e injetado no contexto de cada fase subsequente pelo PromptAggregator.

**Critério de aceite:** SPEC.md gerado com no máximo 2000 tokens; conteúdo relevante e útil para próximas fases; modelo leve utilizado para reduzir custo.

---

### TASK-06-009 — Rastreamento de Tokens e Billing por Execução

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Registrar o consumo de tokens e custo estimado de cada chamada ao LLM para auditoria e billing |

**Descrição:**
Implementar o `BillingTracker` que intercepta todas as chamadas ao AgentSDK e registra: `ProjectID`, `PhaseNumber`, `AgentID`, `Model`, `Provider`, `PromptTokens`, `CompletionTokens`, `TotalTokens`, `EstimatedCostUSD` (cálculo baseado na tabela de preços por provider/modelo), `Timestamp`. Armazenar em coleção separada `billing_records` no MongoDB para facilitar queries de custo por projeto/fase/período. Implementar `GET /api/v1/projects/:id/billing` que retorna o custo acumulado por fase e total do projeto.

**Critério de aceite:** Cada chamada ao LLM registrada; custo estimado calculado por provider/modelo; endpoint de billing funcional.

---

### TASK-06-010 — Sistema de Retry com Backoff Exponencial para Falhas de LLM

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir resiliência do pipeline frente a falhas transitórias dos providers de LLM |

**Descrição:**
Implementar o `RetryableAgentSDK` que envolve o AgentSDK com lógica de retry. Categorias de erro: **Transitório** (rate limit 429, timeout, 503) — retry com backoff exponencial (1s, 2s, 4s, 8s, máximo 5 tentativas), **Permanente** (autenticação 401, modelo não encontrado 404, prompt inválido 400) — falha imediata sem retry, **Desconhecido** — 1 tentativa adicional, depois falha. O estado de retry deve ser visível via evento SSE (`RETRYING {attempt, reason, waitMs}`). Implementar circuit breaker por provider (após 10 falhas consecutivas, pausa novas chamadas por 60 segundos).

**Critério de aceite:** Retry em erros transitórios; falha imediata em erros permanentes; circuit breaker funcional; evento SSE para retry.

---

### TASK-06-011 — Testes de Integração do Pipeline da Tríade

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Validar o fluxo completo de execução da Tríade com mocks dos providers de LLM |

**Descrição:**
Implementar testes de integração usando mocks do AgentSDK para simular respostas dos LLMs sem custo real. Cobrir: execução completa da Tríade do início ao fim (Produtor → Revisor → Refinador), ciclo de feedback (Tríade → aguarda feedback → novo ciclo), aprovação de fase (avanço para próxima fase, geração de SPEC.md), gatilho de rejeição automática (falha catastrófica em Fase 6 retorna para Fase 5), retry em falha transitória (mock retorna 429 nas primeiras 2 tentativas e sucesso na 3ª).

**Critério de aceite:** Todos os fluxos testados com mocks; gatilho de rejeição funcional; retry validado.

---

### TASK-06-012 — Monitoramento e Alertas do Pipeline

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar métricas de observabilidade do pipeline para detectar gargalos e falhas |

**Descrição:**
Instrumentar o pipeline com métricas: taxa de sucesso/falha por fase e por provider, tempo médio de execução da Tríade por fase, fila de mensagens pendentes no RabbitMQ (para detectar backlog), taxa de retry por provider (indica instabilidade), custo médio por execução por fase. Expor métricas no formato Prometheus em `/metrics`. Criar alertas configuráveis: falha de execução consecutiva (3+ falhas), fila acima de 50 mensagens pendentes, custo de um projeto acima do threshold configurável.

**Critério de aceite:** Métricas expostas em `/metrics`; alertas configuráveis; dashboards de observabilidade preparados.

---

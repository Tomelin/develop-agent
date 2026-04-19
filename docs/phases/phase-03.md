# Phase 03 — Catálogo de Agentes de IA

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-03 |
| **Título** | Catálogo de Agentes de IA |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-01, PHASE-02 concluídas |

---

## Descrição Detalhada

O Catálogo de Agentes é o coração da flexibilidade da plataforma. Conforme especificado no PROJECT.md, cada agente é uma entidade configurável com nome, descrição, modelo de IA, system prompts, skills (habilidades), e um toggle de habilitação. O sistema permite que o administrador da plataforma cadastre, configure e gerencie uma biblioteca de agentes especializados que serão utilizados nas Tríades de cada fase do pipeline.

Esta fase implementa o CRUD completo de agentes tanto no backend (API REST) quanto no frontend (interface de gestão), além da lógica de seleção aleatória de agentes para o Modo Dinâmico (multi-modelo). O sistema deve suportar múltiplos providers de LLM (OpenAI, Anthropic, Google Gemini, Ollama) e garantir que somente agentes com `enabled: true` participem das execuções.

A interface de gestão do catálogo deve ser rica e intuitiva, permitindo que o admin visualize todos os agentes cadastrados, filtre por skills e modelo, e configure novos agentes com uma experiência de formulário guiado.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ API REST completa para CRUD de agentes
- ✅ Validação de conectividade com o modelo de IA configurado no agente
- ✅ Lógica de seleção aleatória de agentes por skill (Modo Dinâmico)
- ✅ Interface de listagem de agentes com filtros
- ✅ Interface de criação e edição de agentes com formulário rico
- ✅ Seed com agentes padrão cobrindo todas as fases do pipeline

---

## Funcionalidades Entregues

- **Gestão de Agentes:** CRUD completo com validação de unicidade de nome
- **Configuração de LLM:** Suporte a OpenAI, Anthropic, Google, Ollama com validação de API key
- **Sistema de Skills:** Tags de habilidades que determinam em quais fases o agente pode participar
- **Modo Dinâmico:** Seleção aleatória garantindo diversidade de modelos na Tríade
- **Seed Padrão:** 12 agentes pré-configurados cobrindo todas as especialidades

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

> 💡 **Dica:** Para esta phase, o modo recomendado é **phase completa**, pois o catálogo de agentes e seus prompts são configurações coesas que fazem mais sentido revisadas em bloco.

---

## Tasks

### TASK-03-001 — Modelagem da Entidade Agent no Domínio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir a entidade Agent com todos os atributos necessários para identificação, configuração e uso no pipeline |

**Descrição:**
Definir a struct `Agent` em `src/backend/domain/agent/` com os campos conforme especificado no PROJECT.md: `ID` (ObjectID), `Name` (único), `Description`, `Provider` (enum: `OPENAI`, `ANTHROPIC`, `GOOGLE`, `OLLAMA`), `Model` (string — modelo específico do provider, ex: `claude-3-opus`, `gpt-4o`), `SystemPrompts` (slice de strings — lista de prompts que formam a persona do agente), `Skills` (slice de enums — lista de habilidades como `PROJECT_CREATION`, `ENGINEERING`, `ARCHITECTURE`, `PLANNING`, `DEVELOPMENT_FRONTEND`, `DEVELOPMENT_BACKEND`, `TESTING`, `SECURITY`, `DOCUMENTATION`, `DEVOPS`, `LANDING_PAGE`, `MARKETING`), `Enabled` (boolean), `ApiKeyRef` (referência à chave de API — nunca armazenar a key diretamente), `Status` (enum: `IDLE | RUNNING | PAUSED | QUEUED | ERROR | COMPLETED` — estado operacional do agente, visível no dashboard), `CreatedAt`, `UpdatedAt`. Criar interface `AgentRepository`.

Além da struct, definir o `AgentChannel` em `src/backend/domain/agent/channel.go`:

```go
// AgentMessage é o envelope JSON padrão para toda comunicação entre agentes.
// A comunicação direta entre agentes é proibida — todo dado deve passar por este tipo.
type AgentMessage struct {
    ID        string         `json:"id"`             // UUID único por mensagem
    From      string         `json:"from"`           // ID do agente de origem
    To        string         `json:"to"`             // ID do agente de destino
    Message   string         `json:"message"`        // Payload (artefato, crítica, refinamento)
    Status    string         `json:"status"`         // "pending" | "processing" | "done" | "error"
    Timestamp time.Time      `json:"timestamp"`
    Meta      map[string]any `json:"meta,omitempty"` // Metadados opcionais (phase, token_usage, etc.)
}

// AgentChannel encapsula os canais de entrada e saída de um agente.
// Criado automaticamente quando o agente é instanciado;
// fechado e removido automaticamente quando o agente é encerrado.
type AgentChannel struct {
    In  <-chan AgentMessage // Canal de entrada — mensagens recebidas pelo agente
    Out chan<- AgentMessage // Canal de saída — mensagens enviadas pelo agente
}
```

O tamanho do buffer do channel é configurado via `agent.channel_buffer_size` (padrão: `10`).

**Critério de aceite:** Struct com todos os campos; `AgentStatus` enum com 6 estados; `AgentMessage` e `AgentChannel` definidos; enum de Skills cobrindo todas as fases; interface de repositório definida.

---

### TASK-03-002 — Repositório de Agentes com MongoDB

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar persistência de agentes com suporte a filtros por skill e status |

**Descrição:**
Implementar `AgentMongoRepository` com os métodos: `Create`, `FindByID`, `FindByName`, `Update`, `Delete`, `List` (com filtros: `enabled`, `skill`, `provider`), `FindBySkill` (retorna todos os agentes habilitados para uma skill específica — usado pelo Modo Dinâmico). Criar índices: único em `name`, composto em `{skill, enabled}` para queries de seleção de agentes. Implementar soft delete para preservar histórico.

**Critério de aceite:** Filtros funcionais; índice de skill+enabled criado; `FindBySkill` retorna apenas agentes habilitados.

---

### TASK-03-003 — Serviço de Seleção Aleatória de Agentes (Modo Dinâmico)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a lógica de sorteio de agentes para o Modo Dinâmico garantindo diversidade de modelos na Tríade |

**Descrição:**
Implementar o `AgentSelectorService` em `src/backend/domain/agent/selector.go` com o método `SelectTriad(skill Skill) (Triad, error)`. O método deve: buscar todos os agentes habilitados para a skill solicitada, se o Modo Dinâmico estiver ativado, selecionar 3 agentes diferentes garantindo que nenhum provider se repita (quando possível), se o Modo Dinâmico estiver desativado, usar os agentes configurados fixos para a fase. Registrar o sorteio resultante para auditoria (qual agente foi sorteado para qual papel: Produtor, Revisor, Refinador). Struct de retorno `Triad{Producer, Reviewer, Refiner Agent}`.

**Critério de aceite:** Sorteio garante diversidade de providers; resultado auditável; fallback quando poucos agentes disponíveis.

---

### TASK-03-004 — Integração com Providers de LLM (SDK)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o SDK de integração com múltiplos providers de LLM com interface unificada |

**Descrição:**
Implementar o pacote `src/backend/agentsdk/` com uma interface `Provider` que abstrai a comunicação com qualquer LLM: `Generate(ctx, prompt, systemPrompts []string) (string, TokenUsage, error)`. Implementar os providers concretos: `OpenAIProvider` (usando `github.com/sashabaranov/go-openai`), `AnthropicProvider` (usando o SDK Anthropic Go), `GoogleProvider` (usando Google ADK com Gemini), `OllamaProvider` (HTTP direto para a API local Ollama). Cada provider deve retornar o `TokenUsage` (prompt tokens, completion tokens, total) para o sistema de billing. Implementar retry com backoff exponencial para erros transitórios (rate limit, timeout).

**Critério de aceite:** Interface unificada; todos os providers funcionais; TokenUsage retornado; retry em erros transitórios.

---

### TASK-03-005 — Handler de CRUD de Agentes na API

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor endpoints REST para gestão completa do catálogo de agentes |

**Descrição:**
Implementar handlers em `src/backend/api/handler/agent_handler.go`: `GET /api/v1/agents` (listagem com paginação e filtros: `?enabled=true&skill=ARCHITECTURE&provider=ANTHROPIC`), `GET /api/v1/agents/:id` (detalhe de um agente), `POST /api/v1/agents` (criação — apenas ADMIN), `PUT /api/v1/agents/:id` (atualização — apenas ADMIN), `DELETE /api/v1/agents/:id` (soft delete — apenas ADMIN), `POST /api/v1/agents/:id/test` (testa conectividade do agente com o modelo configurado — envia um prompt simples e valida a resposta). Todas as rotas de escrita requerem role ADMIN.

**Critério de aceite:** CRUD completo funcional; filtros de listagem; endpoint de teste de conectividade; autorização correta.

---

### TASK-03-006 — Gerenciamento Seguro de API Keys dos Providers

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar armazenamento seguro das API keys dos providers de LLM sem expô-las diretamente no banco |

**Descrição:**
Implementar o sistema de referências a segredos. As API keys não devem ser armazenadas diretamente no documento do agente no MongoDB. Em vez disso, armazená-las em Redis com chave `secret:provider:<provider_id>` e TTL longo (ou sem TTL). O campo `ApiKeyRef` no agente é apenas o identificador da referência. Implementar endpoints separados protegidos: `PUT /api/v1/providers/:provider/key` (configura a API key de um provider — ADMIN only), `DELETE /api/v1/providers/:provider/key` (remove a key). A key nunca retorna em nenhuma resposta de API (write-only).

**Critério de aceite:** API key nunca retornada em responses; armazenamento no Redis; write-only por design.

---

### TASK-03-007 — Seed de Agentes Padrão

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Popular o catálogo com agentes pré-configurados cobrindo todas as especialidades do pipeline |

**Descrição:**
Criar seed de 12 agentes padrão no `src/backend/infra/seed/agents_seed.go`, um para cada especialidade: Entrevistador (Fase 1), Engenheiro de Requisitos (Fase 2), Arquiteto de Software (Fase 3), Planejador de Roadmap (Fase 4), Dev Frontend (Fase 5 Front), Dev Backend Golang (Fase 5 Back), QA Engineer (Fase 6), Security Engineer (Fase 7), Tech Writer (Fase 8), DevOps Engineer (Fase 9), Landing Page Designer (Fluxo B), Marketing Strategist (Fluxo C). Cada agente com system prompts bem elaborados que definem sua persona especializada. O seed é idempotente.

**Critério de aceite:** 12 agentes criados com prompts de qualidade; seed idempotente; todos com skills corretas mapeadas.

---

### TASK-03-008 — Interface de Listagem do Catálogo de Agentes

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de listagem do catálogo de agentes com filtros, busca e indicadores visuais de status |

**Descrição:**
Implementar a página `/agents` no frontend com: grid de cards de agentes (nome, provider, skills como tags coloridas, status enabled/disabled com toggle interativo), barra de busca por nome/descrição, filtros dropdown por provider e skill, indicador visual do modelo de IA (badge com cor por provider: verde para Google, azul para OpenAI, laranja para Anthropic, roxo para Ollama), **badge de status operacional do agente** (indicador de `IDLE | RUNNING | PAUSED | QUEUED | ERROR | COMPLETED` com cor e ícone — atualizado via SSE polling), botão "Testar Conexão" por agente com feedback em tempo real, botão "Novo Agente" (apenas para ADMIN). Paginação infinita ao rolar. Layout em grid responsivo 3 colunas no desktop, 1 no mobile.

**Critério de aceite:** Listagem funcional com filtros; toggle de enable/disable funcional; badge de status operacional em tempo real; teste de conexão com feedback visual.

---

### TASK-03-009 — Interface de Criação e Edição de Agente

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar formulário rico para configuração completa de um novo agente ou edição de existente |

**Descrição:**
Implementar o modal/drawer de criação e edição de agente com: campo Nome (com validação de unicidade assíncrona ao digitar), campo Descrição (textarea rico), seletor de Provider (dropdown com ícones), seletor de Modelo (lista dinâmica baseada no provider selecionado), editor de System Prompts (adicionar/remover/reordenar prompts com drag-and-drop), seletor multi-select de Skills (chips clicáveis), toggle Enabled/Disabled. Lógica de preview: ao clicar em "Testar Configuração", enviar um prompt de teste e exibir a resposta do modelo para validação da configuração antes de salvar.

**Critério de aceite:** Formulário com todos os campos; validação assíncrona do nome; preview da resposta do agente funcional.

---

### TASK-03-010 — Página de Detalhe do Agente

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de detalhe do agente exibindo todas as configurações e histórico de uso |

**Descrição:**
Implementar a página `/agents/:id` com: seção de informações gerais (nome, descrição, provider, modelo, status), lista de system prompts com formatação markdown, lista de skills com badges, gráfico de uso (quantas vezes o agente foi utilizado por fase — placeholder para dados reais de fases futuras), histórico das últimas execuções (projeto, fase, papel na Tríade — placeholder), botões de edição e exclusão protegidos por role ADMIN.

**Critério de aceite:** Detalhe completo exibido; seções de histórico com placeholders; botões de ação com controle de role.

---

### TASK-03-011 — Testes Unitários do Seletor de Agentes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que a lógica de seleção aleatória de agentes funciona corretamente em todos os cenários |

**Descrição:**
Implementar testes para o `AgentSelectorService` cobrindo: seleção com providers diversos (verifica que Produtor, Revisor e Refinador têm providers diferentes), seleção com providers limitados (apenas 2 providers disponíveis — deve selecionar 2 diferentes + repetir um), seleção com apenas 1 agente disponível (deve usar o mesmo 3 vezes com log de aviso), seleção com nenhum agente habilitado (deve retornar erro descritivo), modo dinâmico desativado usa configuração fixa.

**Critério de aceite:** Todos os cenários cobertos; diversidade de providers garantida quando possível; erros descritivos.

---

### TASK-03-012 — Documentação Swagger dos Endpoints de Agentes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Documentar todos os endpoints de agentes com Swagger/OpenAPI para facilitar integração |

**Descrição:**
Adicionar anotações Swagger (`swaggo/swag`) em todos os handlers de agentes. Documentar: schemas de request e response, códigos de erro com exemplos, autenticação necessária (Bearer JWT), parâmetros de filtro e paginação. Executar `swag init` para gerar os arquivos de documentação. Disponibilizar a Swagger UI em `/api/v1/docs` (apenas em ambiente de desenvolvimento).

**Critério de aceite:** Swagger UI acessível; todos os endpoints documentados com exemplos; schemas corretos.

---

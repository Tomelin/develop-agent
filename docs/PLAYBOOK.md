# 🏭 PLAYBOOK — Agência de Desenvolvimento de Software com Inteligência Artificial

> **Versão:** 1.0.0  
> **Última atualização:** Abril 2026  
> **Status:** Documento Vivo — sujeito a revisões contínuas

---

## 📋 Índice

1. [Visão Geral e Propósito](#1-visão-geral-e-propósito)
2. [Princípios Fundamentais](#2-princípios-fundamentais)
3. [Arquitetura: A Tríade de Agentes](#3-arquitetura-a-tríade-de-agentes)
4. [Os Três Fluxos de Serviço](#4-os-três-fluxos-de-serviço)
5. [Fluxo A — Desenvolvimento de Software (End-to-End)](#5-fluxo-a--desenvolvimento-de-software-end-to-end)
6. [Fluxo B — Criação de Landing Page](#6-fluxo-b--criação-de-landing-page)
7. [Fluxo C — Estratégia de Marketing](#7-fluxo-c--estratégia-de-marketing)
8. [Sistema de Agentes e Catálogo](#8-sistema-de-agentes-e-catálogo)
9. [Gestão de Perfil e Prompts do Usuário](#9-gestão-de-perfil-e-prompts-do-usuário)
10. [Regras de Interação e Limites de Feedback](#10-regras-de-interação-e-limites-de-feedback)
11. [**Granularidade de Execução: Phase Completa ou Task Individual**](#11-granularidade-de-execução-phase-completa-ou-task-individual)
12. [Otimizações e Recursos Avançados](#12-otimizações-e-recursos-avançados)
13. [Estrutura de Dados e Outputs](#13-estrutura-de-dados-e-outputs)
14. [Stack Tecnológica](#14-stack-tecnológica)
15. [Estrutura de Diretórios do Projeto](#15-estrutura-de-diretórios-do-projeto)
16. [Padrões de Qualidade e Revisão](#16-padrões-de-qualidade-e-revisão)
17. [Glossário](#17-glossário)

---

## 1. Visão Geral e Propósito

A **Agência de IA** é uma plataforma de **esteira de produção de software inteligente e automatizada**, onde cada "profissional" é um agente especialista de Inteligência Artificial. O sistema garante rastreabilidade total desde a ideia inicial até a entrega final, aplicando ciclos de revisão obrigatórios para assegurar consistência, segurança e qualidade de alto nível no ciclo de vida de desenvolvimento de software (SDLC).

### Objetivo Central

Transformar uma ideia bruta em software funcional, landing pages e estratégias de marketing por meio de agentes especializados autônomos, com controle humano nos pontos de decisão mais críticos.

### Por que este sistema existe

- Eliminar gargalos humanos repetitivos no SDLC sem perder qualidade
- Garantir rastreabilidade e auditoria de cada decisão técnica
- Permitir escala sem perda de padronização
- Reduzir viés de modelos únicos com seleção dinâmica multi-modelo

---

## 2. Princípios Fundamentais

| # | Princípio | Descrição |
|---|-----------|-----------|
| 1 | **Revisão Obrigatória** | Nenhum artefato é entregue sem passar pela Tríade completa (Produtor → Revisor → Refinador) |
| 2 | **Rastreabilidade Total** | Cada fase, decisão e interação é registrada e auditável |
| 3 | **Feedback Humano Controlado** | Limites explícitos de feedback por fase para manter eficiência operacional |
| 4 | **Contexto Acumulativo** | O conhecimento gerado em cada fase alimenta as fases seguintes via RAG |
| 5 | **Pluralidade Cognitiva** | Modo dinâmico permite que múltiplos modelos de IA trabalhem em conjunto, eliminando vieses |
| 6 | **Custo Transparente** | Cada token consumido é contabilizado e atribuído ao projeto correspondente |
| 7 | **Separação de Responsabilidades** | Cada agente tem escopo estrito — o Revisor nunca altera, apenas critica |
| 8 | **Controle Granular de Execução** | O usuário decide se executa uma phase inteira ou apenas tasks específicas, combinando velocidade com precisão |
| 9 | **Interface-Driven LLM Integration** | O core de negócio se conecta a qualquer provider de LLM exclusivamente via a interface `agentsdk.Provider` — nunca via SDKs concretos |
| 10 | **Channel-Based Messaging** | Toda comunicação entre agentes é feita via Go channels com envelope JSON padronizado — comunicação direta entre agentes é proibida |

---

## 3. Arquitetura: A Tríade de Agentes e o AgentSDK

O conceito central da plataforma é a **Tríade de Agentes**, aplicada em todas as fases e fluxos de desenvolvimento. Nenhuma entrega sai sem percorrer esse ciclo completo.


### 3.1 — Ciclo da Tríade

```
┌─────────────────────────────────────────────────────────────────┐
│                        TRÍADE DE AGENTES                        │
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐  │
│  │   PRODUTOR   │───▶│   REVISOR    │───▶│    REFINADOR     │  │
│  │              │    │              │    │   (Entregador)   │  │
│  │ Executa a    │    │ Analisa e    │    │ Aplica feedbacks  │  │
│  │ tarefa       │    │ critica —    │    │ e produz a       │  │
│  │ inicial e    │    │ NUNCA altera │    │ versão final     │  │
│  │ gera o       │    │ o artefato   │    │ impecável        │  │
│  │ artefato raiz│    │              │    │                  │  │
│  └──────────────┘    └──────────────┘    └──────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```


### Papel de Cada Agente

#### 🏗️ Agente Produtor
- **Responsabilidade:** Executar a tarefa inicial e gerar o artefato raiz (código, texto, arquitetura)
- **Comportamento:** Foco na execução. Produz o primeiro rascunho completo da entrega
- **Restrição:** Não recebe feedback diretamente do Revisor — o Refinador é o mediador

#### 🔍 Agente Revisor
- **Responsabilidade:** Análise crítica do artefato do Produtor contra regras de negócio, boas práticas e requisitos
- **Comportamento:** Assume **persona de Engenheiro de Qualidade Sênior extremamente cético e paranoico.** Não aprova nada que não esteja perfeito
- **Restrição crítica:** ⚠️ **NUNCA altera o artefato diretamente.** Apenas emite críticas, feedbacks estruturados e correções documentadas
- **Output:** Lista estruturada de problemas encontrados e recomendações de correção

#### ✨ Agente Refinador (Entregador)
- **Responsabilidade:** Receber o artefato do Produtor + as críticas do Revisor, e produzir a versão final
- **Comportamento:** Interage ativamente com cada crítica do Revisor, aplicando os ajustes necessários com raciocínio explícito
- **Output:** Versão final aprovada e auditável do artefato

---

### 3.2 — AgentSDK: Interface `Provider` (Contrato Obrigatório)

Toda a integração com modelos de linguagem é feita exclusivamente através da interface `agentsdk.Provider`. **Nenhum código de domínio ou caso de uso pode depender de um SDK concreto** (OpenAI, Anthropic, Gemini, Ollama). A dependência é sempre em cima da abstração.

#### Por que Interface?

- **Desacoplamento:** Trocar um provider não altera nada no core da plataforma
- **Testabilidade:** Testes usam `MockProvider` determinístico sem custo de API real
- **Modo Dinâmico:** O sorteio da Tríade funciona porque todos satisfazem o mesmo contrato
- **Decoradores:** `BillingTracker` e `RetryWrapper` decoram a interface transparentemente

#### Estrutura de Pacotes Obrigatória

```
src/backend/pkg/agentsdk/
├── provider.go          # Contrato — interface Provider e todos os tipos
├── openai/
│   └── provider.go      # Implementação OpenAI
├── anthropic/
│   └── provider.go      # Implementação Anthropic
├── gemini/
│   └── provider.go      # Implementação Google Gemini
└── ollama/
    └── provider.go      # Implementação Ollama (modelos locais)
```

#### Interface `agentsdk.Provider`

```go
// Provider is the contract that every agent SDK backend must satisfy.
// It abstracts over concrete LLM APIs (OpenAI, Anthropic, Gemini, Ollama) so
// the core layer can drive any provider transparently.
//
// Lifecycle: Initialize → Complete / Stream (any number of calls) → Close.
// Implementations must be safe for concurrent use after Initialize returns.
type Provider interface {
    // Initialize configures the provider using the supplied Config.
    // Must be called once before any other method.
    Initialize(ctx context.Context, cfg Config) error

    // Complete performs a blocking, non-streaming completion round-trip.
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

    // Stream begins a streaming completion and returns a read-only event channel.
    // The caller must drain the channel until Done == true.
    Stream(ctx context.Context, req CompletionRequest) (<-chan StreamEvent, error)

    // Name returns the canonical lower-case identifier of this provider
    // (e.g. "openai", "anthropic", "gemini", "ollama").
    Name() string

    // ModelName returns the model currently configured for this instance.
    ModelName() string

    // Close releases all resources. Must be idempotent.
    Close() error
}
```

#### Regra de Uso — O que é Permitido e Proibido

```go
// ✅ CORRETO — domínio e usecase dependem da abstração
var p agentsdk.Provider = anthropic.New()
p.Initialize(ctx, agentsdk.Config{Token: "sk-ant-...", Model: "claude-3-5-sonnet"})
resp, _ := p.Complete(ctx, agentsdk.CompletionRequest{
    Messages: []agentsdk.Message{
        {Role: agentsdk.RoleSystem, Content: "Você é um arquiteto sênior..."},
        {Role: agentsdk.RoleUser, Content: prompt},
    },
})

// ❌ PROIBIDO — nunca importe SDKs concretos no domínio
import sdk "github.com/anthropics/anthropic-sdk-go"
client := sdk.NewClient(...) // viola a separação de camadas
```

> ⚠️ **Regra de Ouro:** As implementações concretas (`openai`, `anthropic`, `gemini`, `ollama`) só são instanciadas na camada de injeção de dependências (`src/backend/infra/` ou `main.go`). Todo o restante do código depende **exclusivamente** de `agentsdk.Provider`.

---

### 3.3 — Protocolo de Comunicação entre Agentes (Channel-Based Messaging)

Toda troca de informação entre os agentes da Tríade é feita via **Go channels tipados** transportando um **envelope JSON padronizado**. Comunicação direta entre agentes é proibida.

#### Envelope AgentMessage

```go
// AgentMessage é o envelope padrão para toda comunicação entre agentes.
// Todo dado transferido entre agentes DEVE usar este tipo.
type AgentMessage struct {
    ID        string         `json:"id"`                  // UUID único por mensagem
    From      string         `json:"from"`                // ID do agente de origem
    To        string         `json:"to"`                  // ID do agente de destino
    Message   string         `json:"message"`             // Payload (artefato, crítica, refinamento)
    Status    string         `json:"status"`              // "pending" | "processing" | "done" | "error"
    Timestamp time.Time      `json:"timestamp"`           // Momento de criação
    Meta      map[string]any `json:"meta,omitempty"`      // Metadados opcionais
}
```

#### Ciclo de Vida dos Channels

- **Criação:** channel criado automaticamente quando o agente é instanciado (buffer padrão: `10`, configurável via `agent.channel_buffer_size`)
- **Destruição:** channel fechado e removido automaticamente quando o agente é encerrado

#### Status dos Agentes (visível no Dashboard)

| Status | Descrição |
|--------|-----------|
| `IDLE` | Ativo, aguardando mensagens |
| `RUNNING` | Processando uma mensagem |
| `PAUSED` | Suspenso pelo orquestrador |
| `QUEUED` | Há mensagens no buffer aguardando |
| `ERROR` | Falha — requer atenção |
| `COMPLETED` | Tarefa concluída |

> ⚠️ **Regra de Ouro:** A comunicação direta entre agentes é **proibida**. Todo dado transferido entre agentes DEVE passar pelo `AgentMessage` via channel. Isso garante rastreabilidade total e observabilidade da fila de cada agente.

---

## 4. Os Três Fluxos de Serviço

```
┌─────────────────────────────────────────────────────────────────┐
│                    AGÊNCIA DE IA                                │
│                                                                 │
│  ┌─────────────────┐  ┌──────────────────┐  ┌───────────────┐  │
│  │   FLUXO A       │  │    FLUXO B       │  │   FLUXO C     │  │
│  │                 │  │                  │  │               │  │
│  │ Desenvolvimento │  │  Landing Page    │  │  Marketing    │  │
│  │ de Software     │  │                  │  │  Strategy     │  │
│  │ (End-to-End)    │  │  Herda Fluxo A   │  │  Herda Fluxo A│  │
│  │ 9 Fases SDLC    │  │  (opcional)      │  │  (opcional)   │  │
│  └─────────────────┘  └──────────────────┘  └───────────────┘  │
│           │                    │                    │           │
│           └────────────────────┴────────────────────┘           │
│                           TRÍADE DE AGENTES                     │
│                     (aplicada em todos os fluxos)               │
└─────────────────────────────────────────────────────────────────┘
```

---

## 5. Fluxo A — Desenvolvimento de Software (End-to-End)

O Fluxo A é o coração da plataforma. Divide-se em **9 fases estruturais** com execução paralela de Frontend e Backend a partir da Fase 2.

### Diagrama do Fluxo A

```
USUÁRIO
   │
   ▼
┌──────────────────────────────────────────┐
│  FASE 1: Criação do Projeto              │
│  Agente: Entrevistador                   │
│  Feedback: até 10 iterações              │
│  Output: Visão de Produto consolidada    │
└──────────────────────┬───────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────┐
│  FASE 2: Engenharia de Software          │
│  Tríade: Produtor + Revisor + Refinador  │◀── RAG do SPEC.md da Fase 1
│  Feedback: até 5 iterações por fase      │
│  Split: Frontend / Backend               │
└──────────────────────┬───────────────────┘
                       │
          ┌────────────┴───────────┐
          ▼                       ▼
     [FRONTEND]              [BACKEND]
          │                       │
          ▼                       ▼
┌──────────────────┐   ┌──────────────────────┐
│  FASE 3:         │   │  FASE 3:             │
│  Arquitetura     │   │  Arquitetura         │
│  (Front)         │   │  (Back)              │
└────────┬─────────┘   └──────────┬───────────┘
         │                        │
         ▼                        ▼
┌──────────────────────────────────────────┐
│  FASE 4: Planejamento (Roadmap)          │
│  Output: JSON com Tasks + KanBan         │
└──────────────────────┬───────────────────┘
                       │
          ┌────────────┴───────────┐
          ▼                       ▼
     [FRONTEND]              [BACKEND]
┌──────────────────┐   ┌──────────────────────┐
│  FASE 5: Dev     │   │  FASE 5: Dev         │
│  (Front)         │   │  (Back)              │
└────────┬─────────┘   └──────────┬───────────┘
         │                        │
         └────────────┬───────────┘
                      │
                      ▼
┌──────────────────────────────────────────┐
│  FASE 6: Testes de Software              │
│  TDD, Unitários, CI                      │
└──────────────────────┬───────────────────┘
                       │
            ┌──────────┴──────────┐
            │ Falha Catastrófica? │
            │  (Gatilho Auto)     │
            └──────────┬──────────┘
               SIM ────┘ NÃO ──────▶
               │                   │
    Volta Fase 5              FASE 7: Segurança
    (sem descontar            Auditoria OWASP
     feedback manual)
                               │
                               ▼
                   ┌───────────────────────┐
                   │  FASE 8: Documentação │
                   │  README, Manuais, API │
                   └───────────┬───────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │  FASE 9: DevOps/Deploy│
                   │  Dockerfile, K8s,     │
                   │  GitHub Actions       │
                   └───────────────────────┘
```

### Detalhamento das Fases

---

#### 🎯 Fase 1 — Criação do Projeto (Agente Entrevistador)

| Campo | Valor |
|-------|-------|
| **Responsável** | Agente Entrevistador (único agente, sem Tríade) |
| **Objetivo** | Transformar ideia bruta em visão clara de produto |
| **Modo de operação** | Conversacional — faz perguntas, escuta, reformula |
| **Critério de avanço** | Consolidação sólida + aval explícito do usuário |
| **Máximo de feedbacks** | **10 iterações** |
| **Output** | Documento de Visão do Produto validado |

**Comportamento esperado:**
- O agente nunca avança para a próxima fase sem confirmação do usuário
- Reformula o entendimento e pede validação a cada iteração
- Registra todas as decisões de produto como histórico auditável

---

#### ⚙️ Fase 2 — Engenharia de Software

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Definir regras de negócio, requisitos funcionais e não-funcionais |
| **Split** | Inicia a separação Frontend / Backend |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Documento de Engenharia + `SPEC.md` para RAG |

---

#### 🏛️ Fase 3 — Arquitetura de Software

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Modelagem de dados, stack tecnológica, design patterns |
| **Split** | Executada em paralelo: Frontend e Backend |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Diagrama de arquitetura, definição de stack, design patterns adotados |

**Artefatos gerados:**
- Modelagem de dados (entidades, relacionamentos)
- Definição de linguagens e frameworks
- Definição de infraestrutura
- Design patterns adotados (Clean Architecture, DDD, etc.)

---

#### 🗺️ Fase 4 — Planejamento (Roadmap)

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Dividir arquitetura em fases, épicos e tarefas |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | JSON determinístico + KanBan materializado |

**Output estruturado obrigatório (JSON):**

```json
{
  "project_id": "uuid",
  "phases": [
    {
      "id": "uuid",
      "name": "Phase Name",
      "epics": [
        {
          "id": "uuid",
          "title": "Epic Title",
          "tasks": [
            {
              "id": "uuid",
              "title": "Task Title",
              "description": "...",
              "complexity": "LOW|MEDIUM|HIGH|CRITICAL",
              "estimated_hours": 8,
              "type": "FRONTEND|BACKEND|INFRA|TEST|DOC"
            }
          ]
        }
      ]
    }
  ]
}
```

> ⚠️ **Output determinístico:** O backend Golang ingere este JSON e materializa os cards no KanBan visual. Qualquer variação no schema quebrará a integração.

---

#### 💻 Fase 5 — Desenvolvimento

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes (execução ramificada) |
| **Objetivo** | Produção de código funcional |
| **Split** | Frontend e Backend em paralelo |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Código fonte completo por módulo/serviço |

**Gatilho de rejeição automática:** Se a Fase 6 ou 7 identificar falhas catastróficas, o sistema retorna automaticamente para esta fase sem consumir feedbacks manuais do usuário.

---

#### 🧪 Fase 6 — Testes de Software

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Garantia de qualidade e cobertura de testes |
| **Abordagens** | TDD, testes unitários, integração, CI |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Suite de testes + relatório de cobertura |

**Comportamentos especiais:**
- Detecta falhas catastróficas e aciona o **Gatilho de Rejeição Automática**
- Retorna código para Fase 5 sem consumir feedback manual do usuário

---

#### 🔐 Fase 7 — Segurança

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Auditoria de vulnerabilidades |
| **Framework** | OWASP Top 10 + boas práticas de segurança |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Relatório de auditoria de segurança + código corrigido |

**Checklist de segurança (OWASP aplicado):**
- [ ] Injeção (SQL, NoSQL, LDAP)
- [ ] Autenticação e gerenciamento de sessão
- [ ] Exposição de dados sensíveis
- [ ] XML External Entities (XXE)
- [ ] Controle de acesso quebrado
- [ ] Configuração de segurança incorreta
- [ ] Cross-Site Scripting (XSS)
- [ ] Deserialização insegura
- [ ] Uso de componentes com vulnerabilidades conhecidas
- [ ] Log e monitoramento insuficientes

---

#### 📚 Fase 8 — Documentação

| Campo | Valor |
|-------|-------|
| **Responsável** | Tríade de Agentes |
| **Objetivo** | Geração de documentação completa e atualizada |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | README, manuais operacionais, referência de API, tutoriais |

**Artefatos de documentação:**
- `README.md` — Apresentação e guia de início rápido
- `ARCHITECTURE.md` — Visão de arquitetura para desenvolvedores
- `API_REFERENCE.md` — Documentação completa da API REST/gRPC
- `OPERATIONS.md` — Manual operacional para SREs
- Tutoriais de utilização para usuários finais

---

#### 🚀 Fase 9 — DevOps e Deploy (Infra como Código) *(Bônus Evolutivo)*

| Campo | Valor |
|-------|-------|
| **Responsável** | Agente DevOps especializado |
| **Objetivo** | Geração automática de infraestrutura como código |
| **Máximo de feedbacks** | **5 iterações** |
| **Output** | Dockerfiles, docker-compose, GitHub Actions, manifests Kubernetes |

**Artefatos gerados:**
- `Dockerfile` (Frontend e Backend)
- `docker-compose.yml` para desenvolvimento local
- `.github/workflows/` — pipelines de CI/CD
- `k8s/` — manifests YAML para Kubernetes
- `infra/` — configurações de infraestrutura

> 💡 **Como funciona:** O agente DevOps varre a Fase 3 (Arquitetura) e a Fase 5 (código gerado) para inferir automaticamente a configuração de infraestrutura necessária.

---

## 6. Fluxo B — Criação de Landing Page

### Visão Geral

Serviço autônomo focado na geração rápida de landing pages de alta conversão.

```
[Opcional] Herança do Fluxo A
         │
         ▼
┌─────────────────────────────────────┐
│         FLUXO B: LANDING PAGE       │
│                                     │
│  Contexto: Branding + Proposta de   │
│  Valor herdados do Projeto Base     │
│                                     │
│  Tríade de Agentes                  │
│  ┌──────────┐ ┌────────┐ ┌────────┐ │
│  │Produtor  │▶│Revisor │▶│Refina- │ │
│  │(Layout + │ │(Conv.+ │ │dor     │ │
│  │ Copy)    │ │ SEO)   │ │(Final) │ │
│  └──────────┘ └────────┘ └────────┘ │
│                                     │
│  Output: Landing Page completa      │
│  (HTML/CSS/JS ou framework)         │
└─────────────────────────────────────┘
```

### Regras do Fluxo B

- **Herança opcional:** Pode vincular um projeto do Fluxo A para usar como base de branding e proposta de valor
- **Aceleração:** A herança elimina etapas de pesquisa de branding, acelerando a entrega
- **Tríade obrigatória:** Segue o mesmo padrão de Produtor → Revisor → Refinador
- **Máximo de feedbacks:** Seguindo o padrão geral de **5 iterações** por rodada

### Foco de Avaliação do Revisor (Fluxo B)

- Taxa de conversão esperada
- Clareza da proposta de valor (headline, CTA)
- SEO on-page
- Performance de carregamento
- Responsividade (mobile-first)
- Acessibilidade (WCAG 2.1)

---

## 7. Fluxo C — Estratégia de Marketing

### Visão Geral

Serviço para criação de campanhas inteligentes para múltiplos canais de mídia paga e orgânica.

```
[Opcional] Herança do Fluxo A
         │
         ▼
┌─────────────────────────────────────────────┐
│           FLUXO C: MARKETING                │
│                                             │
│  Canais suportados:                         │
│  • LinkedIn Ads + Organic                   │
│  • Instagram Ads + Organic                  │
│  • Google Ads (Search, Display, YouTube)    │
│  • E outros canais de tráfego               │
│                                             │
│  Tríade de Agentes                          │
│  ┌──────────┐ ┌──────────┐ ┌─────────────┐ │
│  │Produtor  │▶│Revisor   │▶│Refinador    │ │
│  │(Strategy │ │(Engagem. │ │(Versão      │ │
│  │+ Copy)   │ │+ ROI)    │ │Final)       │ │
│  └──────────┘ └──────────┘ └─────────────┘ │
│                                             │
│  Output: Campanhas aprovadas por canal      │
└─────────────────────────────────────────────┘
```

### Regras do Fluxo C

- **Herança opcional:** Vincula escopo e arquitetura do projeto raiz do Fluxo A
- **Foco do Revisor:** Engajamento, ROI esperado, compliance de plataforma
- **Tríade obrigatória:** Produtor → Revisor → Refinador com foco em engajamento
- **Máximo de feedbacks:** **5 iterações** por rodada

### Artefatos Gerados

- Estratégia geral da campanha (objetivos, público-alvo, orçamento sugerido)
- Copies por canal e formato (posts, headlines, descrições)
- Calendário editorial sugerido
- KPIs e métricas de sucesso por canal

---

## 8. Sistema de Agentes e Catálogo

### Estrutura de um Agente

Cada agente na plataforma é configurável e possui as seguintes propriedades obrigatórias:

```yaml
agent:
  name: "Bob, Arquiteto Sênior e Especialista Cloud"
  description: "Especialista em design de sistemas distribuídos e cloud-native..."
  model: "claude-3-opus"          # Motor de linguagem que impulsiona o agente
  system_prompts:                 # Identidade base da persona
    - "Você é um arquiteto sênior com 15 anos de experiência..."
    - "Sempre aplique princípios SOLID e Clean Architecture..."
  skills:                         # Etapas do workflow suportadas
    - "software_architecture"
    - "data_modeling"
    - "cloud_design"
  enabled: true                   # Status operacional (Boolean)
```

### Catálogo Padrão de Agentes

| Nome | Especialidade | Fases Aplicáveis |
|------|---------------|-----------------|
| Entrevistador | Product Discovery, UX Research | Fase 1 |
| Engenheiro de Requisitos | Business Rules, Requirements | Fase 2 |
| Arquiteto | System Design, Data Modeling | Fase 3 |
| Planejador | Roadmap, Estimativas | Fase 4 |
| Dev Frontend | UI/UX, React, CSS | Fase 5 (Front) |
| Dev Backend | APIs, Databases, Golang | Fase 5 (Back) |
| QA Engineer | Testing, TDD, CI | Fase 6 |
| Security Engineer | OWASP, Pentest, Auditoria | Fase 7 |
| Tech Writer | Documentação, APIs, README | Fase 8 |
| DevOps Engineer | Docker, K8s, GitHub Actions | Fase 9 |
| Marketing Strategist | Copywriting, Campaigns | Fluxo C |
| Landing Page Designer | Conversion, SEO, HTML/CSS | Fluxo B |

### Modo Dinâmico (Multi-Modelo)

```
┌─────────────────────────────────────────────────────────┐
│                  MODO DINÂMICO ATIVADO                  │
│                                                         │
│  Backend sorteia agentes habilitados aleatoriamente     │
│  para cada posição da Tríade por fase                   │
│                                                         │
│  Exemplo de sorteio:                                    │
│  ┌───────────────┐                                      │
│  │ Produtor      │──▶ GPT-4o (sorteado)                │
│  │ Revisor       │──▶ Claude-3-Opus (sorteado)          │
│  │ Refinador     │──▶ Gemini-1.5-Pro (sorteado)         │
│  └───────────────┘                                      │
│                                                         │
│  Benefício: Visões cognitivas plurais elimina vieses    │
│  de modelos únicos e aumenta robustez do output         │
└─────────────────────────────────────────────────────────┘
```

**Regras do Modo Dinâmico:**
- Somente agentes com `enabled: true` participam do sorteio
- O sistema garante que a Tríade nunca usa o mesmo modelo duas vezes na mesma fase (quando há três ou mais modelos habilitados)
- O sorteio é registrado para auditoria e reprodutibilidade

---

## 9. Gestão de Perfil e Prompts do Usuário

### Usuário Padrão

A plataforma inicia com um usuário padrão chamado **`admin`**.

### Banco de Prompts por Fase

O usuário pode cadastrar múltiplos prompts contextuais vinculados a grupos de agentes (fases). Esses prompts são **agregados** às instruções base de cada fase, funcionando como **regras de domínio persistentes**.

```
PROMPT GLOBAL DO USUÁRIO
        │
        ▼
┌───────────────────────────────────────────────────────────┐
│  + Desenvolvimento utilizando Golang no Backend           │
│  + Framework Next.js no Frontend                          │
│  + Dark Mode como tema padrão                             │
│  + Paleta de cores: Tons escuros com acento Amarelo       │
│  + Testes unitários obrigatórios em todas as funções      │
└───────────────────────────────────────────────────────────┘
        │
        ▼ Injetado automaticamente em cada fase
┌───────────────────────────────────────────────────────────┐
│          BASE DE DIRETRIZES ENVIADAS AO MODELO            │
│  [Prompt do Sistema do Agente] + [Prompts do Usuário]    │
│  + [Contexto RAG da Fase Anterior] + [Input da Fase]     │
└───────────────────────────────────────────────────────────┘
```

### Casos de Uso de Prompts do Usuário

- **Estilização corporativa:** `"Utilize a paleta de cores #1A1A2E, #16213E, #0F3460"`
- **Stack tecnológica:** `"Backend em Golang com GIN, Frontend em React com TailwindCSS"`
- **Padrões de código:** `"Sempre aplique Clean Architecture e SOLID"`
- **Infra:** `"Deploy no Kubernetes com Helm charts"`
- **Banco de dados:** `"Utilize PostgreSQL para dados relacionais e Redis para cache"`
- **Naming conventions:** `"Variáveis em camelCase, funções em PascalCase"`

---

## 10. Regras de Interação e Limites de Feedback

### Tabela de Limites por Fase

| Fase | Serviço | Limite de Feedbacks | Observações |
|------|---------|--------------------:|-------------|
| Fase 1 | Fluxo A | **10** | Fase de descoberta — mais iterações permitidas |
| Fase 2 | Fluxo A | **5** | Padrão para todas as fases técnicas |
| Fase 3 | Fluxo A | **5** | Por trilha (Front e Back separados) |
| Fase 4 | Fluxo A | **5** | |
| Fase 5 | Fluxo A | **5** | + Rejeições automáticas não contam |
| Fase 6 | Fluxo A | **5** | + Rejeições automáticas não contam |
| Fase 7 | Fluxo A | **5** | + Rejeições automáticas não contam |
| Fase 8 | Fluxo A | **5** | |
| Fase 9 | Fluxo A | **5** | |
| — | Fluxo B | **5** | Por sessão de landing page |
| — | Fluxo C | **5** | Por campanha/canal |

> ⚠️ **Importante:** Rejeições automáticas acionadas pelo Gatilho de Rejeição (Fases 6 e 7) **não consomem** o limite de feedbacks manuais do usuário.

---

## 11. Granularidade de Execução: Phase Completa ou Task Individual

O usuário tem **controle total sobre o nível de granularidade** com que deseja acionar a Tríade de Agentes. Ao visualizar o painel de qualquer phase, ele pode escolher entre três modos de execução antes de iniciar.

```
┌─────────────────────────────────────────────────────────────────┐
│            CONTROLE DE GRANULARIDADE DE EXECUÇÃO                │
│                                                                 │
│  [FASE N] — Lista de Tasks                                      │
│                                                                 │
│  ☑ Task 1: Implementar entidade User                            │
│  ☑ Task 2: Criar repositório MongoDB                            │
│  ☑ Task 3: Expor endpoint de criação                            │
│  ☑ Task 4: Testes unitários                                     │
│                                                                 │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────────┐   │
│  │ 🚀 Executar   │  │ 🎯 Executar   │  │ 🔀 Modo          │   │
│  │   Phase       │  │   Task(s)     │  │   Híbrido        │   │
│  │   Completa    │  │   Específica  │  │                  │   │
│  └───────────────┘  └───────────────┘  └───────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 🚀 Modo 1 — Executar a Phase Completa

Todas as tasks da phase são executadas sequencialmente pela Tríade de Agentes (Produtor → Revisor → Refinador), sem interrupção entre elas. O sistema avança automaticamente de uma task para a próxima até concluir a phase inteira.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Uma única aprovação ao final de toda a phase |
| **Velocidade** | Mais rápido — execução contínua sem pausas |
| **Feedback** | Aplicado à phase como um todo (dentro do limite de feedbacks da fase) |
| **Ideal para** | Phases bem compreendidas onde o usuário confia na execução automática |

### 🎯 Modo 2 — Executar Task(s) Específica(s)

O usuário seleciona **uma ou mais tasks individualmente** da lista da phase. A Tríade desenvolve apenas a(s) task(s) selecionada(s) e aguarda aprovação explícita antes de prosseguir para a próxima.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — o usuário controla o ritmo do pipeline |
| **Velocidade** | Mais controlado — requer interação entre tasks |
| **Feedback** | Granular e específico para cada task selecionada |
| **Ideal para** | Tasks críticas, complexas ou que exigem revisão cuidadosa antes de avançar |

> 🔑 **Importante:** O feedback enviado para uma task individual **não é contabilizado** no limite de feedbacks da phase inteira. Cada task tem seu próprio ciclo de aprovação independente.

### 🔀 Modo 3 — Híbrido (Phase Automática + Pausa em Task)

Os dois modos podem ser **combinados livremente** dentro da mesma phase. É possível iniciar a fase em execução automática completa e pausar manualmente em uma task específica que exiga maior atenção. Após aprovar aquela task individualmente, o modo automático retoma a partir da próxima task.

```
Execução Automática ──────▶ [Task 1 ✅] ──▶ [Task 2 ✅] ──▶ PAUSA
                                                                 │
                                    Usuário revisa Task 3 ◀────┘
                                    e envia feedback específico
                                                 │
                              [Task 3 ✅ Aprovada individualmente]
                                                 │
Retoma Automático ────────────────────────────── ▼
                                          [Task 4 ✅] ──▶ FASE CONCLUÍDA
```

### Tabela de Recomendações por Phase

| Phase | Modo Recomendado | Motivo |
|-------|-----------------|--------|
| **Fase 1 — Entrevista** | Individual (conversacional) | Fase interativa por natureza — sempre task a task |
| **Fase 2 — Engenharia** | Phase completa | Output de texto bem estruturado, baixo risco |
| **Fase 3 — Arquitetura** | Híbrido | Revisar arquitetura antes de avançar é recomendado |
| **Fase 4 — Planejamento** | Phase completa | JSON gerado de forma autônoma e validado automaticamente |
| **Fase 5 — Desenvolvimento** | Task a task | Cada task gera código que afeta as tasks seguintes |
| **Fase 6 — Testes** | Phase completa | Testes gerados em bloco são mais consistentes |
| **Fase 7 — Segurança** | Task a task | Issues de segurança são críticos e exigem revisão individual |
| **Fase 8 — Documentação** | Phase completa | Documentação gerada em bloco tem mais coerência |
| **Fase 9 — DevOps** | Phase completa | Infraestrutura como código é gerada de forma holística |
| **Fluxo B — Landing Page** | Híbrido | Revisar copy/design em seções antes de finalizar |
| **Fluxo C — Marketing** | Task a task (por canal) | Cada canal tem particularidades — revisar individualmente |

---

## 12. Otimizações e Recursos Avançados

### ⚡ Handoff e Gestor de Memória (RAG) entre Fases

**Problema resolvido:** Passar o escopo completo da Fase 1 para todas as fases subsequentes causaria sobrecarga de tokens e custos elevados de API.

**Solução:**
```
FASE N finalizada
       │
       ▼
┌──────────────────────────────────┐
│  Modelo Secundário (Leve e Ágil) │
│  Avalia a fase terminada e       │
│  comprime toda a experiência     │
│  em um manifesto resumido        │
│  (SPEC.md)                       │
└──────────────┬───────────────────┘
               │
               ▼
         SPEC.md gerado
               │
               ▼
┌──────────────────────────────────┐
│  FASE N+1                        │
│  Assimila o SPEC.md via RAG      │
│  Mantém o histórico essencial    │
│  Economiza dezenas de milhares   │
│  de Prompt Tokens                │
└──────────────────────────────────┘
```

**Benefícios:**
- Redução significativa de custo por projeto
- Manutenção de contexto semântico preciso entre fases
- Rastreabilidade histórica preservada em formato legível

---

### 🔂 Gatilho de Rejeição Automática Anti-Despesas

**Contexto:** As Fases 6 (Testes) e 7 (Segurança) podem, ao detectar falhas catastróficas, acionar um retorno automático ao código-fonte.

```
FASE 6 ou FASE 7
       │
       ▼
┌─────────────────────────────┐
│  Falha Catastrófica         │
│  detectada?                 │
└─────────────────────────────┘
       │ SIM
       ▼
┌─────────────────────────────────────────────┐
│  Gatilho de Rejeição Automática             │
│                                             │
│  1. Sistema desarma o agente atual          │
│  2. Empacota o relatório de falhas          │
│  3. Retorna para Fase 5 (Dev) em background │
│  4. Injeta correções vitais como instrução  │
│  5. ⚠️ NÃO desconta feedback manual         │
│     do limite de 5 do usuário               │
└─────────────────────────────────────────────┘
       │
       ▼
Fase 5 reexecuta com as correções
```

**Critérios de falha catastrófica:**
- Cobertura de testes abaixo do mínimo configurado
- Vulnerabilidades críticas (CVSS ≥ 9.0)
- Falhas de autenticação ou autorização
- Exposição de dados sensíveis (PII, credenciais)

---

### 📊 Painel de Auditoria de Billing de LLMs

**Objetivo:** Transparência total dos custos de operação da esteira dinâmica.

**Dados coletados por requisição:**

```json
{
  "request_id": "uuid",
  "project_id": "uuid",
  "phase": "ARCHITECTURE",
  "agent_id": "uuid",
  "model": "claude-3-opus",
  "prompt_tokens": 12430,
  "completion_tokens": 3210,
  "total_tokens": 15640,
  "cost_usd": 0.4692,
  "timestamp": "2026-04-19T13:00:00Z"
}
```

**Recursos do painel:**
- Custo total por projeto
- Custo por fase e por agente
- Histórico detalhado de todas as requisições
- Relatório exportável para repasse financeiro aos usuários
- Alertas de orçamento por projeto

---

## 13. Estrutura de Dados e Outputs

### Entidades Principais

#### Projeto
```json
{
  "id": "uuid",
  "name": "Meu SaaS",
  "description": "...",
  "status": "IN_PROGRESS|COMPLETED|PAUSED",
  "current_phase": 3,
  "created_at": "timestamp",
  "user_id": "uuid",
  "spec_md": "...",
  "flows": ["SOFTWARE", "LANDING_PAGE", "MARKETING"]
}
```

#### Fase
```json
{
  "id": "uuid",
  "project_id": "uuid",
  "phase_number": 3,
  "phase_name": "ARCHITECTURE",
  "status": "PENDING|IN_PROGRESS|REVIEW|COMPLETED|REJECTED",
  "track": "FULL|FRONTEND|BACKEND",
  "feedback_count": 2,
  "feedback_limit": 5,
  "artifacts": ["...urls ou conteúdos..."],
  "spec_md": "...",
  "started_at": "timestamp",
  "completed_at": "timestamp"
}
```

#### Agente
```json
{
  "id": "uuid",
  "name": "Bob, Arquiteto Sênior",
  "description": "...",
  "model": "claude-3-opus",
  "system_prompts": ["..."],
  "skills": ["software_architecture"],
  "enabled": true,
  "created_at": "timestamp"
}
```

---

## 14. Stack Tecnológica

### Backend (Core da Esteira)

| Tecnologia | Uso |
|-----------|-----|
| **Golang** | Core da API e orquestração dos agentes |
| **GIN** | Framework HTTP |
| **MongoDB** | Persistência de projetos, fases, agentes e histórico |
| **Redis** | Cache de sessões, filas e estado de orquestração |
| **RabbitMQ** | Mensageria assíncrona entre fases (fase → fila → próxima fase) |

### Frontend (Dashboard)

| Tecnologia | Uso |
|-----------|-----|
| **React** | SPA do painel de controle |
| **TailwindCSS** | Estilização responsiva |
| **Next.js** | Renderização e roteamento |

### Integrações de LLM

| Provider | Modelos disponíveis |
|----------|-------------------|
| **OpenAI** | GPT-4o, GPT-4-Turbo |
| **Anthropic** | Claude-3-Opus, Claude-3-Sonnet, Claude-3-Haiku |
| **Google** | Gemini-1.5-Pro, Gemini-1.5-Flash |
| **Ollama** | Modelos locais (Llama 3, Mistral, etc.) |

### Estrutura de Diretórios

```
develop-agent/
├── src/
│   ├── backend/                  # Core Golang
│   │   ├── config/               # Configuração com Viper
│   │   ├── api/                  # Handlers GIN (REST)
│   │   ├── domain/               # Entidades e regras de negócio
│   │   ├── usecase/              # Casos de uso da aplicação
│   │   ├── infra/                # Adapters (MongoDB, Redis, RabbitMQ)
│   │   ├── agentsdk/             # SDK de integração com LLMs
│   │   └── worker/               # Workers assíncronos por fase
│   └── frontend/                 # React + Next.js + TailwindCSS
│       ├── app/                  # Páginas e rotas
│       ├── components/           # Componentes reutilizáveis
│       └── services/             # Clientes de API
├── _agent/
│   └── skills/                   # Habilidades dos agentes
│       ├── golang-architecture/  # Prompts de arquitetura Golang
│       ├── golang-core/          # Prompts de desenvolvimento Golang
│       └── golang-testing/       # Prompts de testes Golang
├── docs/
│   ├── PROPOSAL.md               # Proposta original do projeto
│   ├── PROJECT.md                # Estruturação oficial do projeto
│   └── PLAYBOOK.md               # Este documento
└── docker-compose.yml            # Ambiente de desenvolvimento
```

---

## 15. Estrutura de Diretórios do Projeto

Para projetos gerados pela plataforma, a estrutura padrão é:

```
<nome-do-projeto>/
├── src/
│   ├── backend/                  # Obrigatório: Todo código Go fica aqui
│   │   ├── cmd/                  # Entrypoints da aplicação
│   │   ├── internal/             # Código privado do serviço
│   │   │   ├── domain/           # Entidades e interfaces
│   │   │   ├── usecase/          # Regras de negócio
│   │   │   └── infra/            # Implementações de infra
│   │   └── pkg/                  # Pacotes exportáveis
│   └── frontend/                 # Obrigatório: Todo código Front fica aqui
│       ├── src/
│       │   ├── app/              # Páginas e rotas
│       │   ├── components/       # Componentes reutilizáveis
│       │   └── hooks/            # Custom hooks
│       └── public/               # Assets estáticos
├── infra/                        # Infraestrutura como código
│   ├── k8s/                      # Manifests Kubernetes
│   └── terraform/                # Terraform (se aplicável)
├── .github/
│   └── workflows/                # GitHub Actions CI/CD
├── docs/                         # Documentação gerada (Fase 8)
├── Dockerfile.backend
├── Dockerfile.frontend
├── docker-compose.yml
└── README.md
```

> ⚠️ **Regra de Ouro:** Todo código backend **obrigatoriamente** reside em `src/backend/` e todo código frontend em `src/frontend/`. Esta convenção é verificada automaticamente pela esteira.

---

## 16. Padrões de Qualidade e Revisão

### Critérios de Aceite por Fase

#### Código (Fase 5)
- [ ] Segue a arquitetura definida na Fase 3
- [ ] Nomenclatura consistente com as convenções de código do usuário
- [ ] Sem código morto ou comentários de debug
- [ ] Tratamento explícito de erros
- [ ] Sem dependências circulares

#### Testes (Fase 6)
- [ ] Cobertura mínima de 80% (configurável por projeto)
- [ ] Testes unitários para toda função de negócio
- [ ] Testes de integração para endpoints críticos
- [ ] Pipeline de CI configurado e funcional

#### Segurança (Fase 7)
- [ ] Zero vulnerabilidades críticas (CVSS ≥ 9.0)
- [ ] Verificação OWASP Top 10 completa
- [ ] Secrets não expostos no código
- [ ] Autenticação e autorização implementadas corretamente

#### Documentação (Fase 8)
- [ ] README com quickstart em menos de 5 minutos
- [ ] API reference completa (todos os endpoints documentados)
- [ ] Diagrama de arquitetura atualizado
- [ ] Guia de contribuição para projetos open-source

### Fluxo de Aprovação por Tríade

```
1. Produtor entrega o artefato
           │
           ▼
2. Revisor analisa e emite feedbacks estruturados
           │
           ├── Nenhum problema? ────▶ Aprova ────▶ Refinador entrega versão final
           │
           └── Problemas encontrados? ──▶ Refinador aplica correções
                       │
                       ▼
           Refinador entrega versão corrigida
                       │
                       ▼
           Revisor re-avalia (se necessário)
```

---

## 17. Glossário

| Termo | Definição |
|-------|-----------|
| **Tríade de Agentes** | Conjunto de 3 agentes (Produtor, Revisor, Refinador) que operam em conjunto para garantir qualidade |
| **Fluxo A/B/C** | Os três serviços disponíveis: Desenvolvimento de Software, Landing Page e Marketing |
| **RAG** | Retrieval-Augmented Generation — técnica de recuperação de contexto para reduzir tokens e custo |
| **SPEC.md** | Manifesto resumido gerado ao final de cada fase para alimentar o contexto da fase seguinte |
| **Modo Dinâmico** | Seleção aleatória de agentes/modelos para cada fase, eliminando vieses de modelo único |
| **Gatilho de Rejeição** | Mecanismo automático que retorna código para a fase anterior sem consumir feedback manual |
| **Feedback Manual** | Interação humana com o sistema para solicitar melhorias — limitada por fase |
| **KanBan** | Visualização das tasks geradas na Fase 4 em formato de cards organizados por status |
| **Billing Panel** | Painel de auditoria de custos de tokens de LLMs por projeto, fase e agente |
| **SDLC** | Software Development Life Cycle — ciclo de vida de desenvolvimento de software |
| **OWASP** | Open Web Application Security Project — framework de segurança para aplicações web |
| **Prompt Tokens** | Tokens enviados ao modelo (input) — cobrados à parte em alguns providers |
| **Completion Tokens** | Tokens gerados pelo modelo (output) — geralmente mais caros |
| **Artefato** | Qualquer entrega concreta de uma fase: código, documento, diagrama, JSON estruturado |
| **Agente Habilitado** | Agente com `enabled: true`, elegível para participar de fases e do Modo Dinâmico |

---

## 📌 Notas Finais e Próximos Passos

### Roadmap de Evolução

- [ ] **v1.0** — Fluxo A completo (Fases 1-8) com Tríade
- [ ] **v1.1** — Integração do Modo Dinâmico multi-modelo
- [ ] **v1.2** — Fluxo B (Landing Page) e Fluxo C (Marketing)
- [ ] **v1.3** — RAG entre fases com SPEC.md
- [ ] **v1.4** — Gatilho de Rejeição Automática
- [ ] **v1.5** — Painel de Billing e Auditoria
- [ ] **v2.0** — Fase 9 (DevOps) e geração de infra como código

### Princípios de Manutenção deste Playbook

1. **Atualize a cada mudança de comportamento** do sistema — este documento é a fonte da verdade operacional
2. **Versione alterações significativas** — use semver para versões do Playbook
3. **Valide com o PROJECT.md** — qualquer divergência entre este Playbook e o PROJECT.md deve ser resolvida imediatamente
4. **Revise limites de feedback** se os usuários reportarem friction excessivo ou insuficiente

---

*Este Playbook é um documento vivo. Versão 1.0.0 — Abril 2026.*

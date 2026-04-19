# Phase 04 — Gestão de Projetos e Dashboard Principal

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-04 |
| **Título** | Gestão de Projetos e Dashboard Principal |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-01, PHASE-02, PHASE-03 concluídas |

---

## Descrição Detalhada

Esta fase implementa o núcleo operacional da plataforma: o gerenciamento de projetos e o dashboard principal do usuário. Um projeto é a entidade central que agrupa todas as fases do pipeline de desenvolvimento, os artefatos gerados, os agentes utilizados e o histórico de interações.

O usuário cria um projeto escolhendo qual fluxo executar (A: Software, B: Landing Page, C: Marketing), e a partir daí o sistema orquestra automaticamente as fases do pipeline. O dashboard principal oferece uma visão consolidada de todos os projetos do usuário, com status em tempo real de cada fase, KPIs de produtividade e acesso rápido às ações mais comuns.

Esta fase também implementa o sistema de KanBan visual que materializa as tasks geradas pela Fase 4 do Fluxo A (Planejamento), além do painel de controle individual de cada projeto com visualização do progresso fase a fase.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ API REST de CRUD de projetos
- ✅ Sistema de estado de projeto com máquina de estados bem definida
- ✅ Dashboard principal com resumo de todos os projetos do usuário
- ✅ **Painel de status dos agentes em tempo real** (IDLE, RUNNING, PAUSED, QUEUED, ERROR, COMPLETED)
- ✅ Painel individual de projeto com progresso das fases
- ✅ KanBan visual de tasks (quando disponível da Fase 4)
- ✅ Vinculação de projeto para herança nos Fluxos B e C

---

## Funcionalidades Entregues

- **CRUD de Projetos:** Criar, listar, editar e arquivar projetos com metadados completos
- **Máquina de Estados:** Controle preciso do estado de cada projeto e fase
- **Dashboard:** Visão consolidada com métricas de progresso e acesso rápido
- **Painel de Agentes:** Status operacional em tempo real de cada agente (IDLE/RUNNING/PAUSED/QUEUED/ERROR/COMPLETED) via SSE
- **KanBan:** Visualização de tasks por fase com drag-and-drop de status
- **Herança de Projeto:** Linkagem de projeto base para Fluxos B e C

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

> 💡 **Dica:** Para esta phase, o modo recomendado é **phase completa**, pois o KanBan do projeto e o dashboard são funcionalidades visuais que ganham mais coerência quando desenvolvidas em conjunto.

---

## Tasks

### TASK-04-001 — Modelagem da Entidade Project no Domínio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir a entidade Project com estados, fases e metadados de rastreabilidade |

**Descrição:**
Definir a struct `Project` com: `ID`, `Name`, `Description`, `FlowType` (enum: `SOFTWARE`, `LANDING_PAGE`, `MARKETING`), `Status` (enum: `DRAFT`, `IN_PROGRESS`, `PAUSED`, `COMPLETED`, `ARCHIVED`), `CurrentPhaseNumber` (int), `Phases` (slice de `PhaseExecution`), `LinkedProjectID` (para Fluxos B e C herdarem do Fluxo A), `OwnerUserID`, `DynamicModeEnabled` (boolean — Modo Dinâmico de seleção de agentes), `SpecMD` (conteúdo do manifesto de contexto RAG acumulado), `TotalTokensUsed`, `TotalCostUSD`, `CreatedAt`, `UpdatedAt`. Struct `PhaseExecution` com: `PhaseNumber`, `PhaseName`, `Status`, `Track` (FULL/FRONTEND/BACKEND), `FeedbackCount`, `FeedbackLimit`, `Artifacts` (urls/conteúdos), `AgentTriad` (Produtor, Revisor, Refinador utilizados), `StartedAt`, `CompletedAt`.

**Critério de aceite:** Entidades com todos os campos; estados bem tipados; relacionamentos corretos.

---

### TASK-04-002 — Máquina de Estados de Projeto e Fase

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar transições de estado válidas para projetos e fases, prevenindo estados inválidos |

**Descrição:**
Implementar o `ProjectStateMachine` em `src/backend/domain/project/state_machine.go` que valida e executa transições de estado. Transições de projeto: `DRAFT → IN_PROGRESS` (ao iniciar a primeira fase), `IN_PROGRESS → PAUSED` (pela ação do usuário), `PAUSED → IN_PROGRESS`, `IN_PROGRESS → COMPLETED` (quando todas as fases concluídas), `IN_PROGRESS → ARCHIVED` (ação do usuário). Transições de fase: `PENDING → IN_PROGRESS`, `IN_PROGRESS → REVIEW` (Tríade finalizada, aguardando feedback), `REVIEW → IN_PROGRESS` (novo feedback recebido), `REVIEW → COMPLETED` (usuário aprova), `IN_PROGRESS → REJECTED` (gatilho automático), `REJECTED → IN_PROGRESS` (após correção automática). Qualquer tentativa de transição inválida retorna erro descritivo.

**Critério de aceite:** Transições válidas executadas; transições inválidas bloqueadas com erro; histórico de transições registrado.

---

### TASK-04-003 — Repositório de Projetos com MongoDB

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar persistência de projetos com queries otimizadas para o dashboard |

**Descrição:**
Implementar `ProjectMongoRepository` com: `Create`, `FindByID`, `FindByOwner` (com paginação e filtros por status e FlowType), `Update`, `Archive`, `UpdatePhase` (atualiza apenas a fase corrente sem sobrescrever o documento inteiro), `UpdateSpecMD` (atualiza o manifesto de contexto RAG). Criar índices: `{owner_id, status}` para o dashboard do usuário, `{status, created_at}` para relatórios administrativos. Usar projeção para o dashboard (retornar apenas `id, name, status, currentPhase, updatedAt` — não carregar artefatos completos).

**Critério de aceite:** Queries com projeção eficiente; índices criados; `UpdatePhase` atômico sem race condition.

---

### TASK-04-004 — Handlers de CRUD de Projetos na API

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor endpoints REST para ciclo de vida completo dos projetos |

**Descrição:**
Implementar handlers em `src/backend/api/handler/project_handler.go`: `GET /api/v1/projects` (listagem dos projetos do usuário logado, com filtros e paginação), `GET /api/v1/projects/:id` (detalhe completo do projeto), `POST /api/v1/projects` (criação — valida linked project se FlowType for B ou C), `PUT /api/v1/projects/:id` (edição de metadados enquanto DRAFT), `POST /api/v1/projects/:id/pause` (pausa projeto em andamento), `POST /api/v1/projects/:id/resume` (retoma projeto pausado), `POST /api/v1/projects/:id/archive` (arquiva projeto). O usuário só acessa seus próprios projetos (user_id da sessão JWT).

**Critério de aceite:** CRUD funcional; isolamento por usuário; transições de estado via endpoints específicos.

---

### TASK-04-005 — Sistema de Vinculação de Projetos (Herança Fluxos B e C)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a lógica de herança de contexto entre projetos para acelerar Fluxos B e C |

**Descrição:**
Implementar o `ProjectInheritanceService` que, quando um projeto do Fluxo B ou C é criado com `LinkedProjectID`, extrai automaticamente do projeto raiz (Fluxo A): o `SpecMD` (contexto consolidado de todas as fases já completadas), o nome e descrição do produto, as tecnologias utilizadas, a paleta de cores e identidade visual (se documentada na Fase 8). Este contexto extraído é injetado como contexto inicial do novo projeto do Fluxo B ou C, eliminando a necessidade de re-informar o branding ao criar a landing page ou estratégia de marketing.

**Critério de aceite:** Contexto extraído corretamente do projeto base; injeção no contexto inicial do novo projeto; validação de que o projeto base está COMPLETED ou IN_PROGRESS.

---

### TASK-04-006 — Dashboard Principal no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar o dashboard principal com visão consolidada de todos os projetos e métricas do usuário |

**Descrição:**
Implementar a página `/dashboard` com:

**Seção de Projetos:**
Header com estatísticas globais (total de projetos, projetos em andamento, concluídos, tokens utilizados este mês); grid de cards de projeto (nome, status com badge colorido, fluxo tipo, fase atual, % de progresso, última atividade); barra de busca e filtros por status e tipo de fluxo; botão "Novo Projeto" proeminente que abre o wizard de criação; estado vazio encorajador para usuário sem projetos ainda; polling de atualização a cada 30 segundos.

**Painel de Status dos Agentes (`AgentStatusPanel`):**
Seção lateral ou em aba separada exibindo o status operacional em tempo real de todos os agentes do sistema:
- Indicadores visuais por status: `IDLE` ● cinza, `RUNNING` ● verde pulsante (animação), `PAUSED` ● amarelo, `QUEUED` ● azul com contador de mensagens na fila, `ERROR` ● vermelho com tooltip do erro, `COMPLETED` ● verde sólido
- Nome do agente, provider e modelo associado
- Task atual em execução (para agentes `RUNNING`)
- Número de mensagens no buffer do channel (para agentes `QUEUED`)
- Conectado ao SSE `GET /api/v1/agents/status/stream` com fallback para polling de 5s

**Critério de aceite:** Dashboard carrega projetos; polling funcional; painel de agentes atualiza em tempo real via SSE; indicadores visuais corretos para cada status; estado vazio amigável; filtros funcionais.

---

### TASK-04-007 — Wizard de Criação de Novo Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar um assistente guiado multi-etapa para criação de novos projetos |

**Descrição:**
Implementar o wizard de criação de projeto com 3 steps: **Step 1 — Tipo de Fluxo:** cards visuais para escolha entre Fluxo A (Software), Fluxo B (Landing Page) e Fluxo C (Marketing), com descrição e features de cada um. **Step 2 — Configurações Base:** nome do projeto, descrição inicial, toggle de Modo Dinâmico (com explicação do que é), se Fluxo B ou C, seletor de projeto existente para herança (opcional). **Step 3 — Confirmação:** resumo visual da configuração escolhida antes de criar. Barra de progresso dos steps, navegação entre steps com validação, botão "Criar Projeto" que redireciona para o painel do projeto recém-criado.

**Critério de aceite:** Wizard 3 steps funcional; validação por step; redirecionamento após criação; herança de projeto funcional nos Fluxos B e C.

---

### TASK-04-008 — Painel Individual de Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de detalhe de projeto com visão completa do progresso das fases |

**Descrição:**
Implementar a página `/projects/:id` com: header do projeto (nome, tipo de fluxo, status, ações: pausar/retomar/arquivar), timeline visual das fases (barras de progresso sequenciais com status por fase — pending, in_progress, review, completed, rejected), fase ativa em destaque com informações da Tríade de agentes em uso e contador de feedbacks restantes, seção de artefatos gerados por fase (links/previews), seção de custos (tokens e custo USD por fase), botão de ação principal contextual ("Iniciar Fase", "Enviar Feedback", "Aprovar Entrega").

**Critério de aceite:** Timeline de fases visual e precisa; fase ativa em destaque; artefatos navegáveis; custos exibidos.

---

### TASK-04-009 — KanBan Visual de Tasks

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Implementar o KanBan de tasks geradas pela Fase 4 (Planejamento) com drag-and-drop de status |

**Descrição:**
Implementar o componente KanBan na aba "KanBan" do painel de projeto. Colunas: `TODO`, `IN_PROGRESS`, `DONE`, `BLOCKED`. Cards de task com: título, descrição (tooltip), tipo (FRONTEND/BACKEND/INFRA/TEST/DOC — badge colorido), complexidade (LOW/MEDIUM/HIGH/CRITICAL — indicador visual), estimativa de horas. Drag-and-drop entre colunas para atualizar status. Filtros por tipo e complexidade. Exportação do KanBan como JSON. O KanBan é populado automaticamente quando a Fase 4 é concluída no Fluxo A.

**Critério de aceite:** KanBan com 4 colunas; drag-and-drop funcional; filtros por tipo/complexidade; populado automaticamente pela Fase 4.

---

### TASK-04-010 — API de Tasks do Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar endpoints para gestão das tasks geradas pelo roadmap |

**Descrição:**
Implementar a entidade `Task` com: `ID`, `ProjectID`, `PhaseID`, `EpicID`, `Title`, `Description`, `Type` (FRONTEND/BACKEND/INFRA/TEST/DOC), `Complexity` (LOW/MEDIUM/HIGH/CRITICAL), `EstimatedHours`, `Status` (TODO/IN_PROGRESS/DONE/BLOCKED), `AssignedAgentID` (para rastreabilidade), `CreatedAt`, `UpdatedAt`. Implementar handlers: `GET /api/v1/projects/:id/tasks` (listagem com filtros), `PUT /api/v1/projects/:id/tasks/:taskId/status` (atualiza status — acionado pelo drag-and-drop do KanBan), `POST /api/v1/projects/:id/tasks/bulk` (criação em massa a partir do JSON da Fase 4).

**Critério de aceite:** Tasks persistidas; atualização de status por drag-and-drop; bulk create a partir do JSON da Fase 4.

---

### TASK-04-011 — Sistema de Notificações In-App

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Notificar o usuário sobre eventos importantes do pipeline sem exigir polling manual |

**Descrição:**
Implementar Server-Sent Events (SSE) no backend (`GET /api/v1/projects/:id/events`) que emite eventos em tempo real para o frontend quando: uma fase muda de status, a Tríade completa a revisão e aguarda feedback do usuário, um gatilho de rejeição automática é acionado, uma fase é concluída com sucesso. No frontend, conectar ao SSE ao abrir o painel de projeto e atualizar a UI sem reload. Exibir toast notifications para cada evento recebido. Fallback para polling de 10 segundos se SSE não for suportado.

**Critério de aceite:** SSE emite eventos corretos; frontend atualiza sem reload; toast notifications exibidas; fallback para polling.

---

### TASK-04-012 — Testes de Integração da Gestão de Projetos

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Validar o fluxo completo de criação e progressão de projetos com testes de integração |

**Descrição:**
Implementar testes de integração usando environment Docker para: criação de projeto e validações (nome duplicado, FlowType inválido, linked project inexistente), transições de estado válidas e inválidas na máquina de estados, isolamento de projetos por usuário (usuário A não acessa projetos do usuário B), criação bulk de tasks a partir de JSON da Fase 4. Usar `testcontainers-go` para subir MongoDB e Redis localmente nos testes.

**Critério de aceite:** Testes de integração passando; isolamento de usuário validado; máquina de estados testada.

---

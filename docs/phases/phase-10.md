# Phase 10 — Fluxo A: Fase 5 — Desenvolvimento de Software (Frontend e Backend)

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-10 |
| **Título** | Fluxo A — Fase 5: Desenvolvimento de Software (Frontend e Backend) |
| **Tipo** | Backend |
| **Prioridade** | Crítica |
| **Pré-requisitos** | PHASE-09 concluída (Roadmap gerado e ingerido) |

---

## Descrição Detalhada

A **Fase 5 do Fluxo A** é o coração do produto gerado pela plataforma: o **Desenvolvimento de Software**. É aqui que os agentes programadores constroem, em código real e funcional, toda a solução projetada nas fases de Engenharia e Arquitetura.

Esta fase opera de forma ramificada nos dois trilhos (Frontend e Backend) simultaneamente. Os agentes são especializados por trilho: o Agente Dev Backend é especialista em Golang com GIN, Clean Architecture, MongoDB e RabbitMQ; o Agente Dev Frontend é especialista em React, TypeScript, TailwindCSS e Next.js.

O pipeline de desenvolvimento é orientado pelas tasks do KanBan geradas na Fase 4. Para cada task, a Tríade gera o código correspondente. O sistema mantém um contexto de código acumulado (similar ao RAG, mas para código) para que os agentes mantenham consistência entre as tasks.

Esta phase também implementa o **Gatilho de Rejeição Automática** que, quando as Fases 6 (Testes) ou 7 (Segurança) detectam falhas catastróficas, retorna o projeto para cá sem consumir feedbacks manuais do usuário.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Pipeline de desenvolvimento task-by-task orientado pelo KanBan
- ✅ Agentes Dev especializados por trilho (Frontend Golang/React)
- ✅ Sistema de contexto de código acumulado entre tasks
- ✅ Terminal SSE para exibição do código em tempo real
- ✅ Gatilho de Rejeição Automática implementado e funcional
- ✅ Repositório virtual com todos os arquivos gerados organizados

---

## Funcionalidades Entregues

- **Desenvolvimento Task-by-Task:** Cada task do KanBan vira um ciclo completo da Tríade
- **Contexto de Código:** Agentes mantêm consistência com o código já gerado nas tasks anteriores
- **Terminal ao Vivo:** O usuário vê o código sendo gerado em tempo real
- **Repositório Virtual:** Todos os arquivos gerados organizados em estrutura de diretórios

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa (Modo Automático)

Todas as tasks de desenvolvimento são executadas sequencialmente pela Tríade de Agentes, sem interrupções. O sistema avança automaticamente de uma task para a próxima, aprovando internamente cada entregável antes de avançar.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a phase |
| **Velocidade** | Mais rápido — ideal para projetos menores ou bem especificados |
| **Feedback** | Aplicado ao conjunto completo de código gerado |
| **Ideal para** | Projetos com requisitos bem definidos e baixa criticidade |

### 🎯 Executar uma Task Específica ⭐ **Recomendado**

O usuário seleciona **uma ou mais tasks individualmente** do KanBan. A Tríade desenvolve apenas a(s) task(s) selecionada(s) e aguarda aprovação explícita antes de continuar.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — máximo controle sobre o código gerado |
| **Velocidade** | Mais controlado — requer revisão entre tasks |
| **Feedback** | Específico para cada task — o contexto de código é atualizado para as próximas |
| **Ideal para** | Desenvolvimento de software onde cada task afeta diretamente as seguintes |

### 🔀 Modo Híbrido

É possível **combinar os dois modos**: inicie em automático para tasks simples (CRUD, modelos) e pause manualmente nas tasks críticas (autenticação, integrações, segurança) que exigem revisão cuidadosa do código gerado.

> ⭐ **Fortemente Recomendado:** Use **task a task** nesta phase. O código gerado em cada task compõe o contexto das tasks seguintes — revisar e aprovar task por task garante consistência arquitetural ao longo de todo o desenvolvimento. Você pode sempre selecionar múltiplas tasks para execução em lote quando forem de baixo risco.

---

## Tasks

### TASK-10-001 — Prompts Especializados para Agentes Dev Backend (Golang)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts de sistema que tornam o Agente Developer Backend especialista em Golang |

**Descrição:**
Criar system prompts abrangentes para o Agente Dev Backend cobrindo: (1) Estrutura de código Go com Clean Architecture (`domain/`, `usecase/`, `infra/`, `api/`), (2) Padrões de código Golang: interfaces para injeção de dependência, tratamento explícito de erros, uso de context.Context em funções que fazem I/O, (3) GIN para handlers HTTP com validação via `go-playground/validator`, (4) MongoDB com repositório pattern, (5) Godoc em todas as funções públicas, (6) Testes unitários com testify para cada função de negócio, (7) Nunca gerar código com hardcoded secrets. O Produtor gera a implementação completa de uma task. O Revisor verifica Clean Architecture, tratamento de erros, testabilidade. O Refinador produz código de nível sênior.

**Critério de aceite:** Prompts geram código Go de qualidade profissional; Clean Architecture respeitada; godoc em todas as funções.

---

### TASK-10-002 — Prompts Especializados para Agentes Dev Frontend (React/TypeScript)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts para o Agente Developer Frontend especialista em React/TypeScript/TailwindCSS |

**Descrição:**
Criar system prompts para o Agente Dev Frontend cobrindo: (1) React com TypeScript (tipagem forte, sem `any`), (2) Componentes funcionais com hooks, sem class components, (3) TailwindCSS para estilização, respeitando a paleta de cores do usuário, (4) Acessibilidade (ARIA labels, roles, keyboard navigation), (5) Performance: lazy loading, memoização com useCallback/useMemo onde justificado, (6) Testes com React Testing Library, (7) Internacionalização básica estruturada, (8) Responsividade mobile-first. O Revisor verifica UX, acessibilidade, tipos TypeScript corretos, ausência de prop drilling excessivo. O Refinador entrega componentes de produção.

**Critério de aceite:** Prompts geram React/TypeScript de qualidade; acessibilidade incluída; responsividade garantida.

---

### TASK-10-003 — Execução de Desenvolvimento por Task (Task-Driven Development)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o fluxo de desenvolvimento onde cada task do KanBan gera um ciclo completo da Tríade |

**Descrição:**
Implementar o `TaskDevelopmentOrchestrator` que: recebe uma Task do KanBan (ou um conjunto de tasks relacionadas de um épico), compõe o prompt de desenvolvimento com: contexto de código já gerado (via CodeContextAccumulator), detalhes da task (título, descrição, tipo, complexidade), arquitetura definida na Fase 3, requisitos da Fase 2, prompts do usuário, executa o ciclo completo da Tríade para a task, ao final, extrai os arquivos de código do output do Refinador, atualiza o `CodeContext` com os novos arquivos, atualiza o status da task no KanBan para `DONE`, avança para a próxima task. O usuário pode optar por execução automática sequencial (todas as tasks) ou manual (uma task por vez).

**Critério de aceite:** Execução task-by-task funcional; contexto de código atualizado após cada task; status KanBan atualizado.

---

### TASK-10-004 — Sistema de Contexto de Código Acumulado (CodeContext)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Manter um contexto de código acumulado que permite que os agentes gerem código consistente entre tasks |

**Descrição:**
Implementar o `CodeContextAccumulator` que mantém um manifesto do código gerado até o momento. Por ser inviável enviar todos os arquivos gerados para cada nova task (limite de tokens), o acumulador mantém: índice de arquivos gerados (caminho + descrição de propósito de cada arquivo), interfaces e structs públicas exportadas (para Backend), componentes e suas props (para Frontend), variáveis de ambiente necessárias, dependências adicionadas. Este índice comprimido é injetado no prompt de cada nova task para garantir consistência (importações corretas, nomes de funções corretos, etc.). Tamanho máximo do contexto: 3000 tokens.

**Critério de aceite:** Contexto comprimido gerado corretamente; injetado no prompt de cada task; máximo de 3000 tokens respeitado.

---

### TASK-10-005 — Repositório Virtual de Arquivos Gerados

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Organizar e persistir todos os arquivos gerados em uma estrutura de diretórios virtual acessível |

**Descrição:**
Implementar o `VirtualRepository` que armazena os arquivos gerados pela Fase 5. Cada arquivo tem: `path` (caminho relativo, ex: `src/backend/domain/user/user.go`), `content` (conteúdo do arquivo), `taskId` (task que gerou o arquivo), `language` (golang/typescript/etc.), `version` (timestamp da última versão), `phaseNumber`. Armazenar no MongoDB como coleção `code_files` com índice em `{project_id, path}`. Implementar endpoints: `GET /api/v1/projects/:id/files` (lista a árvore de arquivos), `GET /api/v1/projects/:id/files/:fileId` (conteúdo de um arquivo), `GET /api/v1/projects/:id/files/download` (download de todos os arquivos como ZIP).

**Critério de aceite:** Arquivos organizados em estrutura de diretórios; versionamento por task; download ZIP funcional.

---

### TASK-10-006 — Terminal de código em Tempo Real (SSE)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Exibir o código sendo gerado em tempo real via SSE, como um terminal de desenvolvimento ao vivo |

**Descrição:**
Adaptar o sistema SSE para incluir eventos de código em andamento: `CODE_BLOCK_START {language, filepath}`, `CODE_CONTENT {chunk}` (fragmentos de código em streaming), `CODE_BLOCK_END {filepath, linesCount}`. No frontend, implementar um componente "Terminal de Desenvolvimento" que exibe o código sendo gerado em tempo real com: syntax highlighting dinâmico por linguagem, indicação da task sendo desenvolvida, indicação do agente em uso (Produtor/Revisor/Refinador), contador de linhas de código gerado, histórico de arquivos gerados nesta sessão como lista navegável.

**Critério de aceite:** Streaming de código funcional; syntax highlighting em tempo real; histórico de arquivos navegável.

---

### TASK-10-007 — Explorador de Arquivos (File Explorer) no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar um explorador de arquivos estilo IDE para navegar pelo código gerado |

**Descrição:**
Implementar o componente `FileExplorer` na aba "Código" do painel de projeto com: árvore de arquivos expansível organizada por diretórios, ícones por extensão de arquivo (Go fish ícone para .go, React ícone para .tsx, etc.), ao clicar em um arquivo, exibir o conteúdo no painel lateral com syntax highlighting (Monaco Editor ou CodeMirror), breadcrumb de navegação, busca de arquivos por nome, botão de download de arquivo individual, opção de download de todos os arquivos como ZIP.

**Critério de aceite:** Árvore de arquivos funcional; editor de código com syntax highlighting; busca por nome; download individual e ZIP.

---

### TASK-10-008 — Gatilho de Rejeição Automática (Auto-Rejection Trigger)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o mecanismo que retorna o desenvolvimento para a Fase 5 quando falhas catastróficas são detectadas nas Fases 6 e 7 |

**Descrição:**
Implementar o `AutoRejectionTrigger` que é chamado pelas Fases 6 (Testes) e 7 (Segurança) quando detectam falhas catastróficas. Critérios de falha catastrófica: cobertura de testes abaixo de 50% (configurável), vulnerabilidade crítica CVSS ≥ 9.0, exposição de credentials no código, falha de compilação do código gerado. O trigger: cria uma nova execução de desenvolvimento na Fase 5 injetando o relatório de falhas como contexto mandatório, **não consome feedback manual do usuário** (FeedbackCount não é incrementado), notifica o usuário sobre a rejeição automática com uma mensagem clara explicando o motivo, executa a reexecução em background. O usuário é notificado quando o problema foi corrigido automaticamente.

**Critério de aceite:** Trigger não consome feedback manual; usuário notificado; relatório de falhas injetado no contexto de correção.

---

### TASK-10-009 — Configuração de Execução Automática vs Manual

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir que o usuário escolha entre execução automática sequencial ou manual task-a-task |

**Descrição:**
Implementar dois modos de execução para a Fase 5: **Modo Automático:** o sistema executa todas as tasks em sequência automaticamente, aprovando cada uma com base em critérios internos (Tríade concluída sem erros graves), e o usuário apenas acompanha o progresso. **Modo Manual:** após cada task, o sistema aguarda aprovação explícita do usuário antes de avançar para a próxima task. O usuário pode enviar feedback específico para uma task antes de aprovar. Implementar `POST /api/v1/projects/:id/phases/5/mode` para alternar entre modos. No frontend, toggle visível no painel da Fase 5 com explicação das diferenças.

**Critério de aceite:** Modo automático processa todas as tasks sem interrupção; modo manual aguarda aprovação por task; alternância dinâmica entre modos.

---

### TASK-10-010 — Validação de Compilação do Código Gerado

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar automaticamente se o código Go gerado compila corretamente antes de concluir a task |

**Descrição:**
Implementar o `CodeCompilationValidator` que, após o Refinador gerar o código Go de uma task, executa `go build ./...` em um sandbox isolado (container Docker efêmero) para verificar se o código compila. Se houver erros de compilação, eles são retornados ao Refinador como feedback automático para correção (sem consumir feedback manual do usuário). Máximo de 2 tentativas de autocorreção de compilação por task. Se após 2 tentativas ainda não compilar, marca a task como `BLOCKED` e notifica o usuário. Para o frontend TypeScript, executar `tsc --noEmit` de forma similar.

**Critério de aceite:** Compilação verificada em sandbox; erros retornados ao Refinador automaticamente; task bloqueada após 2 falhas.

---

### TASK-10-011 — Resumo e Métricas da Fase de Desenvolvimento

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Consolidar métricas da Fase 5 para acompanhamento de progresso e qualidade |

**Descrição:**
Implementar `GET /api/v1/projects/:id/phases/5/summary` retornando: total de tasks do KanBan, tasks concluídas (DONE), tasks em andamento (IN_PROGRESS), tasks bloqueadas (BLOCKED), número de arquivos gerados por trilho (Frontend/Backend), total de linhas de código gerado, tempo médio de desenvolvimento por task, número de rejeições automáticas acionadas, custo total em tokens da Fase 5. No frontend, exibir como dashboard de progresso no painel da Fase 5 com gráfico de progresso, lista de tasks por status e métricas de produtividade dos agentes.

**Critério de aceite:** Métricas calculadas corretamente; dashboard visual no frontend; tempo por task rastreado.

---

### TASK-10-012 — Integração com Sistema de Versionamento (Git)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar um repositório Git local para o código gerado, com commits automáticos por task concluída |

**Descrição:**
Implementar o `GitIntegration` que mantém um repositório Git local para cada projeto. A cada task concluída e aprovada: cria um commit com a mensagem padronizada `feat(task-id): título da task`, inclui todos os arquivos gerados/modificados pela task. Ao final da Fase 5: cria uma tag `v0.1.0-dev-complete`. Implementar `GET /api/v1/projects/:id/git/log` que retorna o histórico de commits. `GET /api/v1/projects/:id/git/download` gera um tarball do repositório completo com histórico Git para download.

**Critério de aceite:** Commits automáticos por task; mensagens padronizadas; tag ao final da fase; download do repositório com histórico.

---

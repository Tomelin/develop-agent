# Phase 05 — Gestão de Prompts do Usuário

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-05 |
| **Título** | Gestão de Prompts do Usuário |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-01, PHASE-02 concluídas |

---

## Descrição Detalhada

O sistema de prompts do usuário é um dos diferenciais mais poderosos da plataforma. Conforme definido no PROJECT.md e PROPOSAL.md, cada usuário pode cadastrar múltiplos prompts contextuais agrupados por fase do pipeline. Esses prompts são **agregados automaticamente** às instruções base de cada agente antes de cada execução, funcionando como regras de domínio persistentes e personalizadas.

Exemplos de uso: o usuário define uma vez que "todo backend deve ser em Golang com GIN" e que "o design deve seguir a paleta de cores #1A1A2E em dark mode" — a partir daí, todos os agentes de todas as fases receberão essas diretrizes automaticamente sem que o usuário precise repetir nas conversas.

Esta fase implementa o CRUD completo de prompts no backend, a interface de gestão no frontend, e o motor de agregação de prompts que compõe a instrução final enviada aos modelos de IA em cada execução.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ API REST para CRUD de prompts por grupo/fase
- ✅ Motor de agregação que compõe a instrução final com prompts do usuário
- ✅ Interface de gestão de prompts com editor rico
- ✅ Preview de como os prompts serão injetados em cada fase
- ✅ Sistema de templates de prompts pré-definidos

---

## Funcionalidades Entregues

- **Banco de Prompts por Grupo:** Prompts categorizados por fase e tipo de fluxo
- **Agregação Automática:** Composição automática do prompt final antes de cada execução
- **Templates Padrão:** Biblioteca de prompts pré-definidos para padrões comuns
- **Preview de Composição:** Visualização de como os prompts serão compostos antes de executar

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

> 💡 **Dica:** Para esta phase, o modo recomendado é **híbrido** — execute automaticamente os prompts globais/de grupo e revise individualmente as tasks de criação de prompts específicos antes de aplicá-los ao pipeline.

---

## Tasks

### TASK-05-001 — Modelagem da Entidade UserPrompt no Domínio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir a entidade UserPrompt com agrupamento por fase e ordem de injeção configurável |

**Descrição:**
Definir a struct `UserPrompt` com: `ID`, `UserID`, `Title` (nome descritivo do prompt, ex: "Stack Backend Golang"), `Content` (conteúdo do prompt — até 2000 chars), `Group` (enum: `GLOBAL`, `PROJECT_CREATION`, `ENGINEERING`, `ARCHITECTURE`, `PLANNING`, `DEVELOPMENT`, `TESTING`, `SECURITY`, `DOCUMENTATION`, `DEVOPS`, `LANDING_PAGE`, `MARKETING`), `Priority` (int — ordem de injeção, menor = primeiro), `Enabled` (boolean), `Tags` (strings para organização), `CreatedAt`, `UpdatedAt`. Prompts do grupo `GLOBAL` são injetados em TODAS as fases. Criar interface `UserPromptRepository`.

**Critério de aceite:** Struct com campo Group cobrindo todas as fases; prioridade de injeção configurável; GLOBAL injetado em todos.

---

### TASK-05-002 — Repositório de Prompts do Usuário

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Persistir e recuperar prompts do usuário com ordenação por prioridade |

**Descrição:**
Implementar `UserPromptMongoRepository` com: `Create`, `FindByID`, `FindByUserAndGroup` (retorna prompts de um usuário para um grupo específico, ordenados por `Priority ASC`, apenas `Enabled: true`), `FindAllByUser` (listagem completa para a interface de gestão, com suporte a filtros por group e enabled), `Update`, `Delete`, `Reorder` (atualiza a prioridade de múltiplos prompts em uma operação atômica). Criar índice composto `{user_id, group, priority}` para queries eficientes.

**Critério de aceite:** Query por grupo retorna prompts ordenados e habilitados; reorder atômico; índice composto criado.

---

### TASK-05-003 — Motor de Agregação de Prompts

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a lógica de composição do prompt final a partir de múltiplas fontes |

**Descrição:**
Implementar o `PromptAggregator` em `src/backend/domain/prompt/aggregator.go`. O aggregator compõe o prompt final na seguinte ordem: **(1)** System prompts do agente selecionado (identidade e persona base), **(2)** Prompts GLOBAL do usuário (diretrizes universais — ex: linguagem de programação, paleta de cores), **(3)** Prompts específicos do grupo da fase atual (ex: prompts do grupo `ARCHITECTURE` para a Fase 3), **(4)** Contexto RAG da fase anterior (conteúdo do `SPEC.md`), **(5)** Instrução específica da fase (o que exatamente o agente deve fazer nesta execução). Retornar a composição final em formato de `[]Message` compatível com qualquer provider de LLM.

**Critério de aceite:** Composição na ordem correta; GLOBAL inserido em todos os grupos; contexto RAG incluído; output compatível com providers.

---

### TASK-05-004 — Handlers de CRUD de Prompts na API

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor endpoints REST para gestão completa dos prompts do usuário |

**Descrição:**
Implementar handlers em `src/backend/api/handler/prompt_handler.go`: `GET /api/v1/prompts` (listagem de todos os prompts do usuário logado, com filtros por group e enabled), `GET /api/v1/prompts/:group` (prompts de um grupo específico), `POST /api/v1/prompts` (criação de novo prompt), `PUT /api/v1/prompts/:id` (atualização de prompt), `DELETE /api/v1/prompts/:id` (remoção), `PUT /api/v1/prompts/reorder` (reordena prioridades em batch), `GET /api/v1/prompts/preview/:group` (preview da composição final para um grupo — útil para validação antes de executar uma fase). Todos os endpoints isolados por `user_id` da sessão.

**Critério de aceite:** CRUD completo funcional; preview de composição funcional; isolamento por usuário.

---

### TASK-05-005 — Sistema de Templates de Prompts Padrão

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Fornecer uma biblioteca de templates de prompts pré-definidos para onboarding rápido |

**Descrição:**
Criar um endpoint `GET /api/v1/prompts/templates` que retorna uma lista de templates de prompts organizados por grupo e por caso de uso, como: Stack Backend Golang+GIN, Stack Frontend React+TailwindCSS, Dark Mode Design, PostgreSQL como banco de dados, Clean Architecture, SOLID Principles, TDD First, OWASP Security Priority, RESTful API Design, GraphQL API Design. Implementar o endpoint `POST /api/v1/prompts/from-template` que cria um prompt no banco do usuário a partir de um template selecionado. Templates são estáticos no código (não persistidos no banco), mas podem ser expandidos via arquivo de configuração YAML.

**Critério de aceite:** Lista de templates retornada por categoria; criação de prompt a partir de template funcional; templates úteis e bem escritos.

---

### TASK-05-006 — Interface de Gestão de Prompts no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de gestão de prompts com editor rico e organização por grupos/fases |

**Descrição:**
Implementar a página `/prompts` no frontend organizada como tabs por grupo: Global, Project Creation, Engineering, Architecture, Planning, Development, Testing, Security, Documentation, DevOps, Landing Page, Marketing. Em cada tab: lista de prompts do grupo com título, preview do conteúdo (truncado), prioridade, toggle de enabled/disabled, botões de editar e remover. Reordenação por drag-and-drop entre prompts do mesmo grupo. Botão "Adicionar Prompt" por tab. Botão "Importar do Template" que abre modal com biblioteca de templates.

**Critério de aceite:** Tabs por grupo funcional; drag-and-drop de reordenação; toggle de enable/disable; importação de templates.

---

### TASK-05-007 — Editor de Prompt Rico

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar um editor de prompt com sugestões, formatação e preview de composição |

**Descrição:**
Implementar o modal/drawer de criação e edição de prompt com: campo título (nome descritivo), editor de textarea com contador de caracteres (limite 2000), seletor de grupo (dropdown com todas as categorias), campo de tags (chips para organização), toggle enabled/disabled, seção de dicas de escrita de prompts efetivos (expandível), botão "Preview de Composição" que abre painel lateral mostrando como o prompt se encaixa no prompt final completo da fase selecionada.

**Critério de aceite:** Editor com contador de chars; preview de composição funcional; dicas de boas práticas.

---

### TASK-05-008 — Preview de Composição de Prompts

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Permitir que o usuário visualize exatamente como seus prompts serão compostos antes de executar uma fase |

**Descrição:**
Implementar a tela `/prompts/preview` com: seletor de fase/grupo (dropdown), seletor de agente (para visualizar os system prompts do agente escolhido), exibição visual da composição em blocos coloridos: bloco azul = system prompts do agente, bloco verde = prompts GLOBAL do usuário, bloco amarelo = prompts do grupo selecionado, bloco roxo = contexto RAG (representado por placeholder), bloco cinza = instrução da fase (representado por placeholder). Botão de copiar o prompt completo composto.

**Critério de aceite:** Preview com blocos coloridos por origem; seletor de agente e grupo; botão de copiar funcionando.

---

### TASK-05-009 — Validação e Limites de Prompts

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar validações e limites para garantir que o sistema não seja sobrecarregado por prompts excessivos |

**Descrição:**
Implementar validações no serviço de prompts: limite máximo de 50 prompts por grupo por usuário, limite de 2000 caracteres por prompt, estimativa de tokens total dos prompts (usando tokenizer simplificado) e aviso quando exceder 4000 tokens (pode impactar custo e performance), validação de que o conteúdo do prompt não contém instruções maliciosas que possam subverter o comportamento dos agentes (lista básica de padrões proibidos como "ignore all previous instructions"). Retornar erros descritivos para cada violação.

**Critério de aceite:** Limites de quantidade e tamanho respeitados; estimativa de tokens no response; detecção básica de prompt injection.

---

### TASK-05-010 — Exportação e Importação de Configuração de Prompts

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir que o usuário exporte e importe sua configuração completa de prompts para backup e portabilidade |

**Descrição:**
Implementar `GET /api/v1/prompts/export` que gera um arquivo JSON com todos os prompts do usuário organizados por grupo (útil para backup e compartilhamento de configurações). Implementar `POST /api/v1/prompts/import` que recebe o JSON e cria ou atualiza os prompts do usuário (com opção de merge ou substituição total). No frontend, botões "Exportar" e "Importar" na tela de prompts. Ao importar, exibir preview dos prompts que serão criados antes de confirmar.

**Critério de aceite:** Exportação gera JSON válido; importação com preview antes de confirmar; opção de merge ou substituição.

---

### TASK-05-011 — Testes Unitários do Motor de Agregação

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que a composição de prompts está correta em todos os cenários |

**Descrição:**
Testar o `PromptAggregator` nos seguintes cenários: composição com todos os tipos de prompts disponíveis (verificar ordem correta), composição sem prompts GLOBAL (deve funcionar sem), composição sem prompts do grupo (apenas GLOBAL e system prompts do agente), composição sem contexto RAG (primeira fase do projeto), usuário sem prompts cadastrados (apenas system prompts do agente e instrução da fase), limite máximo de prompts (verificar que o token estimate está correto).

**Critério de aceite:** Todos os cenários testados e passando; ordem de composição validada; edge cases cobertos.

---

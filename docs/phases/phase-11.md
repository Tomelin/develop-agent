# Phase 11 — Fluxo A: Fase 6 — Testes de Software

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-11 |
| **Título** | Fluxo A — Fase 6: Testes de Software |
| **Tipo** | Backend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-10 concluída (Desenvolvimento implementado) |

---

## Descrição Detalhada

A **Fase 6 do Fluxo A** é a de **Testes de Software**. Os agentes desta fase são especialistas em garantia de qualidade (QA) e revisão de código com foco em testabilidade. Eles analisam o código gerado na Fase 5 e produzem uma suite completa de testes que aumenta significativamente a confiança no software entregue.

A abordagem é guiada pelo TDD (Test-Driven Development) em espírito — os agentes de teste identificam os comportamentos esperados a partir dos requisitos da Fase 2 e verificam se o código da Fase 5 os implementa corretamente. Para o backend Go, são gerados testes unitários com `testify`, testes de integração com `testcontainers-go` e configuração de pipeline de CI via GitHub Actions. Para o frontend React, são gerados testes com React Testing Library e Cypress para testes E2E dos fluxos críticos.

O aspecto mais crítico desta fase é o **Gatilho de Rejeição Automática**: se os agentes de teste encontrarem falhas catastróficas (cobertura abaixo do mínimo, funcionalidades críticas sem comportamento esperado), o framework retorna automaticamente o código para a Fase 5 para correção, sem consumir feedbacks manuais do usuário.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Suite completa de testes unitários para Backend (Go) e Frontend (React)
- ✅ Testes de integração para os fluxos mais críticos
- ✅ Pipeline de CI configurado com GitHub Actions
- ✅ Relatório de cobertura de testes com threshold configurável
- ✅ Gatilho de Rejeição Automática baseado em cobertura mínima
- ✅ Testes E2E para os fluxos críticos do frontend (Cypress)

---

## Funcionalidades Entregues

- **Testes Unitários:** Cobertura de toda lógica de negócio Backend e Frontend
- **Testes de Integração:** Validação de endpoints críticos contra banco real
- **Pipeline CI:** GitHub Actions executando testes em cada push
- **Relatório de Cobertura:** Visualização de cobertura por arquivo e função

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa ⭐ **Recomendado**

Todos os testes são gerados sequencialmente pela Tríade de Agentes, sem interrupções. O sistema avança automaticamente de uma task para a próxima até concluir a suite completa de testes.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a phase |
| **Velocidade** | Mais rápido — testes gerados em bloco têm mais consistência |
| **Feedback** | Aplicado à suite de testes como um todo |
| **Ideal para** | Geração de testes que têm melhor coerência quando vistos em conjunto |

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. A Tríade desenvolve apenas a(s) task(s) escolhida(s) e aguarda aprovação explícita antes de prosseguir.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — o usuário controla o ritmo |
| **Velocidade** | Mais controlado — requer interação entre tasks |
| **Feedback** | Granular e específico para cada tipo de teste |
| **Ideal para** | Quando se deseja gerar apenas testes unitários (sem E2E), por exemplo |

### 🔀 Modo Híbrido

É possível **combinar os dois modos**: execute automaticamente os testes unitários e de integração, mas revise individualmente a configuração do CI/CD (TASK-11-004) antes de incluí-la no repositório.

> 💡 **Dica:** Para esta phase (Testes — Fase 6 do Fluxo A), o modo recomendado é **phase completa** — testes gerados de forma holística são mais consistentes e coerentes entre si.

---

## Tasks

### TASK-11-001 — Prompts do Agente QA Engineer para Testes Go

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts que tornam o Agente QA especialista em testes de Golang |

**Descrição:**
Criar system prompts para o Agente QA Backend com instruções de: (1) Análise do código da Fase 5 para identificar todos os casos de uso cobertos e não cobertos, (2) Testes unitários com `github.com/stretchr/testify/assert` e `testify/mock` para isolamento de dependências, (3) Naming convention: `Test[FunctionName]_[Scenario]_[ExpectedResult]` (ex: `TestCreateUser_WithDuplicateEmail_ReturnsConflictError`), (4) Table-driven tests para múltiplos cenários de uma função, (5) Testes de integração usando `testcontainers-go` para MongoDB, Redis, (6) Coverage mínima de 80% em packages de domain e usecase, (7) Mock de providers externos (LLM, email). Prompt do Revisor QA: verificar que os testes realmente testam o comportamento (não apenas que chamam funções), detectar testes que passam sempre sem testar nada, verificar coverage.

**Critério de aceite:** Prompts geram testes com naming correto; table-driven tests; mocks de dependências; verificação de comportamento.

---

### TASK-11-002 — Prompts do Agente QA para Testes Frontend (React Testing Library)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts para testes de componentes React com foco em comportamento do usuário |

**Descrição:**
Criar prompts para o Agente QA Frontend com instruções de: (1) React Testing Library com `@testing-library/react` — testar o que o usuário vê, não os internos do componente, (2) Testes de interação com `userEvent` (clicks, digitação, submissão de formulários), (3) Testes de estados: loading, erro, vazio, populado, (4) Mock de chamadas de API com `msw` (Mock Service Worker), (5) Assertions de acessibilidade com `jest-axe`, (6) Testes de rotas com `MemoryRouter`, (7) Snapshot tests para componentes estáticos de apresentação. Prompt do Revisor: verificar que testes não testam implementação (sem `.instance()`, sem acesso a estado interno), verificar scenarios de erro cobertos.

**Critério de aceite:** Testes com RTL testam behavior do usuário; MSW para mock de API; acessibilidade testada.

---

### TASK-11-003 — Análise de Cobertura de Testes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Medir e reportar a cobertura de testes do código gerado na Fase 5 |

**Descrição:**
Implementar o `CoverageAnalyzer` que: executa `go test ./... -cover -coverprofile=coverage.out` em sandbox, gera o relatório HTML com `go tool cover -html=coverage.out`, parseia o relatório para extrair: cobertura por package, cobertura por função, cobertura total do projeto. Armazena o relatório de cobertura como artefato da Fase 6. Compara a cobertura obtida com o threshold configurável por projeto (default: 80%). Se abaixo do threshold, aciona o `AutoRejectionTrigger` com relatório detalhado dos packages/funções não cobertos para guiar a Fase 5 na correção.

**Critério de aceite:** Relatório de cobertura gerado; comparação com threshold; aciona rejeição automática se abaixo do mínimo.

---

### TASK-11-004 — Configuração de Pipeline CI com GitHub Actions

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar a configuração de CI que executa todos os testes automaticamente a cada commit |

**Descrição:**
Implementar o agente de CI que gera o arquivo `.github/workflows/ci.yml` para o projeto do usuário com: **Backend workflow:** checkout, setup-go, go mod download, golangci-lint, go test com coverage, publish coverage report como comentário no PR. **Frontend workflow:** checkout, setup-node, npm install, tsc --noEmit, eslint, jest (com coverage), upload coverage para Codecov (opcional). **Triggers:** em push para main/develop e em pull requests. Armazenar o arquivo gerado no VirtualRepository do projeto. Validar sintaxe YAML do arquivo gerado antes de entregar.

**Critério de aceite:** CI YAML gerado com jobs Backend e Frontend; triggers corretos; sintaxe YAML válida.

---

### TASK-11-005 — Testes de Integração para Endpoints Críticos

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar testes de integração que validam os endpoints mais críticos contra infraestrutura real (em container) |

**Descrição:**
Implementar o agente de testes de integração que analisa os endpoints gerados na Fase 5 e seleciona os mais críticos (endpoints de autenticação, CRUD principal, fluxos de negócio) para gerar testes de integração usando `testcontainers-go`. Os testes: sobem containers de MongoDB e Redis, executam requests HTTP reais contra o servidor GIN, validam o estado do banco após execução, testam cenários de erro (404, 400, 401, 409), garantem idempotência quando necessário. Os testes de integração ficam em `src/backend/integration_test/` para separação clara.

**Critério de aceite:** Testes de integração executam contra containers reais; cenários de sucesso e erro cobertos; separação em pasta dedicada.

---

### TASK-11-006 — Gerador de Testes E2E com Cypress

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar testes end-to-end dos fluxos críticos do frontend usando Cypress |

**Descrição:**
Implementar o agente de testes E2E que, baseado nos fluxos de usuário identificados na documentação, gera testes Cypress para: fluxo de login (login válido, inválido, logout), fluxo de criação de projeto (wizard completo), navegação pelo KanBan (filtros, drag-and-drop básico), submissão de feedback de fase, visualização de artefato gerado. Os testes Cypress ficam em `src/frontend/cypress/e2e/`. Incluir configuração `cypress.config.ts` com baseUrl, viewports mobile e desktop, e configuração de fixtures para dados de teste.

**Critério de aceite:** Testes Cypress gerados para fluxos críticos; configuração correta; testes executam sem erros.

---

### TASK-11-007 — Dashboard de Resultados de Testes no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir os resultados de cobertura e testes de forma visual no painel do projeto |

**Descrição:**
Implementar a aba "Testes" no painel de projeto com: gráfico circular de cobertura total (verde > 80%, amarelo 60-80%, vermelho < 60%), tabela de cobertura por package/módulo com indicador colorido, lista de testes que falharam (se houver) com mensagem de erro expandível, comparação de cobertura entre versões (antes/depois de um ciclo da Tríade), badges de status do CI (passando/falhando). Carregar os dados do endpoint de artefatos da Fase 6.

**Critério de aceite:** Gráfico de cobertura com código de cores; tabela por package; comparação entre versões; badges de CI.

---

### TASK-11-008 — Relatório de Qualidade de Testes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar um relatório consolidado de qualidade dos testes para o usuário final |

**Descrição:**
Ao concluir a Fase 6, gerar o documento `QUALITY_REPORT.md` como artefato com: sumário executivo da qualidade do código, cobertura por camada (domain, usecase, infra, api), funções críticas sem cobertura (lista prioritizada), tipos de testes gerados (unitários, integração, E2E) com contagem, issues de testabilidade identificadas (código difícil de testar sugere refatoração), recomendações de melhoria futura. O documento é gerado pelo Refinador da Tríade de Testes com base nos dados coletados.

**Critério de aceite:** QUALITY_REPORT.md gerado com todas as seções; issues de testabilidade identificadas; recomendações úteis.

---

### TASK-11-009 — Verificação de Compilação e Execução dos Testes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que os testes gerados compilam e executam corretamente antes de entregar |

**Descrição:**
Implementar o `TestExecutionValidator` que, após o Refinador gerar os arquivos de teste, executa `go test ./... -run .` em sandbox para verificar: compilação dos testes, execução sem panics, ausência de race conditions (executar com flag `-race`). Para testes de frontend: `npm test -- --watchAll=false`. Se algum teste falhar por erro de implementação do próprio teste (não por bug no código do projeto), retornar ao Refinador para correção automática (sem consumir feedback do usuário). Se o teste revelar um bug real no código, acionar o `AutoRejectionTrigger`.

**Critério de aceite:** Testes compilam e executam; race conditions detectadas; distinção entre erro no teste vs bug no código.

---

### TASK-11-010 — Testes de Snapshot para Componentes de UI

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar testes de snapshot para componentes React de apresentação para detectar regressões visuais |

**Descrição:**
O agente QA Frontend gera testes de snapshot (com `@testing-library/react` + `jest`) para os componentes de apresentação pura (botões, cards, badges, formulários em estado inicial). Os snapshots são armazenados em `src/frontend/src/__snapshots__/`. Implementar o script de atualização de snapshots no package.json (`npm run test:update-snapshots`). Incluir no README do projeto gerado as instruções de quando e como atualizar snapshots.

**Critério de aceite:** Snapshots gerados para componentes de apresentação; script de atualização; README com instruções.

---

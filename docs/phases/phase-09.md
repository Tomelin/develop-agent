# Phase 09 — Fluxo A: Fase 4 — Planejamento e Roadmap (KanBan Output)

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-09 |
| **Título** | Fluxo A — Fase 4: Planejamento e Roadmap com Output KanBan |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-08 concluída (Fases 2 e 3 do projeto implementadas) |

---

## Descrição Detalhada

A **Fase 4 do Fluxo A** é o Planejamento e elaboração do Roadmap. Esta fase é única porque seu output principal não é um documento Markdown, mas sim um **JSON determinístico** que o backend Golang ingere e materializa em cards de KanBan visual para o usuário.

O Agente Planejador analisa toda a arquitetura definida na Fase 3 e a quebra em unidades gerenciáveis: fases de desenvolvimento, épicos e tasks. Cada task recebe metadados ricos: tipo (Frontend/Backend/Infra/Test/Doc), complexidade estimada (LOW/MEDIUM/HIGH/CRITICAL) e estimativa em horas.

Esta phase implementa a lógica de parsing e validação do JSON de output da Tríade de Planejamento, a ingestão das tasks no banco de dados, e a interface KanBan visual que o usuário utilizará para acompanhar o progresso do desenvolvimento real.

Um aspecto crítico desta fase é garantir que o JSON gerado pelo agente seja rigorosamente estruturado — o Agente Revisor tem como objetivo principal verificar a validade do schema JSON antes de aprovar o output.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Agentes de planejamento com prompt que gera JSON determinístico
- ✅ Parser e validador do schema JSON de roadmap
- ✅ Ingestão automática das tasks no banco de dados após aprovação
- ✅ KanBan visual completo com todas as tasks do roadmap
- ✅ Métricas de roadmap (total de tasks, por tipo, estimativa total de horas)
- ✅ Exportação do roadmap em JSON e CSV

---

## Funcionalidades Entregues

- **JSON de Roadmap:** Output estruturado com fases, épicos e tasks bem definidas
- **Ingestão Automática:** Tasks criadas no sistema ao aprovar a Fase 4
- **KanBan Dinâmico:** Visualização de tasks por fase, épico e status
- **Métricas de Esforço:** Estimativas totais e por categoria

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

> 💡 **Dica:** Para esta phase (Planejamento e Roadmap — Fase 4 do Fluxo A), o modo recomendado é **phase completa**. O JSON de roadmap e o KanBan são gerados de forma holística. Revise o output consolidado e use o feedback da phase para ajustar o roadmap inteiro antes de aprovar.

---

## Tasks

### TASK-09-001 — Prompts do Agente Planejador para Output JSON

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts que instruem o Agente Planejador a gerar um JSON determinístico e bem estruturado |

**Descrição:**
Criar os system prompts para o Agente Planejador (Produtor da Fase 4) que instruem: (1) analisar a arquitetura da Fase 3 e os requisitos da Fase 2 para decompor em tasks, (2) gerar **exclusivamente** um JSON válido sem texto antes ou depois (critical: sem markdown code fences), (3) seguir rigorosamente o schema definido (ver TASK-09-002), (4) cada task deve ser granular o suficiente para ser executada por um único agente em uma única sessão, (5) estimar complexidade baseando-se no tipo de trabalho (criação de CRUD = LOW, integração de sistema externo = HIGH). Prompt do Revisor: validar estrutura JSON, verificar que cada task tem contexto suficiente, checar completude de cobertura dos requisitos, verificar estimativas coerentes.

**Critério de aceite:** Agente gera JSON puro sem markdown; schema correto; Revisor valida estrutura.

---

### TASK-09-002 — Schema JSON de Roadmap e Validação

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir e implementar o schema JSON obrigatório para o roadmap gerado pelo agente |

**Descrição:**
Definir o schema JSON e implementar o `RoadmapSchemaValidator` que valida o output do agente:

```json
{
  "project_id": "string",
  "phases": [{
    "id": "string",
    "name": "string",
    "description": "string",
    "order": "int",
    "epics": [{
      "id": "string",
      "title": "string",
      "description": "string",
      "tasks": [{
        "id": "string",
        "title": "string",
        "description": "string",
        "type": "FRONTEND|BACKEND|INFRA|TEST|DOC",
        "complexity": "LOW|MEDIUM|HIGH|CRITICAL",
        "estimated_hours": "int",
        "track": "FRONTEND|BACKEND|FULL",
        "dependencies": ["task_id"]
      }]
    }]
  }]
}
```

O validador verifica: JSON parseável, campos obrigatórios presentes, enums válidos, `estimated_hours` >= 1 e <= 200, sem `task_id` de dependência referenciando task inexistente. Em caso de falha de validação, o validador retorna erros estruturados que são usados como prompt de correção para o Refinador.

**Critério de aceite:** Schema definido; validação completa; erros estruturados retornados ao Refinador para autocorreção.

---

### TASK-09-003 — Parser e Ingestão do JSON de Roadmap

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Parsear o JSON do roadmap e criar todas as tasks no banco de dados ao aprovar a Fase 4 |

**Descrição:**
Implementar o `RoadmapIngester` que é executado automaticamente quando o usuário aprova a Fase 4. O ingester: parseia o JSON do output do Refinador, valida o schema (via RoadmapSchemaValidator), cria registros de `Phase`, `Epic` e `Task` no banco de dados em operação bulk atômica, associa todas as tasks ao projeto, emite evento `ROADMAP_INGESTED {taskCount, phaseCount, epicCount}` via SSE, atualiza o status da Fase 4 para COMPLETED. Se o JSON for inválido na hora da aprovação (edge case), retorna erro e solicita novo ciclo da Tríade.

**Critério de aceite:** Ingestão atômica; todas as entidades criadas; evento emitido; validação antes de persistir.

---

### TASK-09-004 — Métricas e Resumo do Roadmap

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Calcular e expor métricas consolidadas do roadmap para planejamento do projeto |

**Descrição:**
Implementar `GET /api/v1/projects/:id/roadmap/summary` que retorna: total de tasks por tipo (FRONTEND/BACKEND/INFRA/TEST/DOC), total de tasks por complexidade (LOW/MEDIUM/HIGH/CRITICAL), total de horas estimadas por tipo e por phase de desenvolvimento, número de fases e épicos no roadmap, path crítico estimado (soma de horas das tasks CRITICAL). Implementar como query agregada no MongoDB com pipeline de aggregation para performance. Exibir no painel do projeto como cards de métricas.

**Critério de aceite:** Métricas calculadas via aggregation MongoDB; path crítico estimado; response < 100ms.

---

### TASK-09-005 — Interface KanBan Completa

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Implementar o KanBan visual rico com todas as features necessárias para acompanhar o desenvolvimento |

**Descrição:**
Implementar a página `/projects/:id/kanban` com: colunas fixas (TODO, IN_PROGRESS, REVIEW, DONE, BLOCKED), cards de tasks com título, badge de tipo colorido, badge de complexidade, indicador de trilha (Front/Back), estimativa de horas, opção de expandir para ver descrição completa, drag-and-drop de cards entre colunas (atualiza status via API), cabeçalho de coluna com contagem de tasks e soma de horas estimadas, filtro superior: por fase de dev, por épico, por tipo (FRONTEND/BACKEND), por complexidade, busca por título/descrição.

**Critério de aceite:** KanBan funcional com drag-and-drop; filtros operacionais; contagem e horas por coluna; cards expansíveis.

---

### TASK-09-006 — Visualização de Épicos e Dependências

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Visualizar o roadmap em formato de épicos com agrupamento e dependências entre tasks |

**Descrição:**
Implementar a aba "Épicos" na tela do roadmap com: accordion por fase de desenvolvimento, dentro de cada fase, lista de épicos como seções colapsáveis, dentro de cada épico, lista de tasks em modo lista (não KanBan), indicadores de dependências entre tasks (seta visual), destaque de tasks bloqueadas (dependências não concluídas), progresso por épico (tasks DONE / total tasks). Status de épico calculado automaticamente: PENDING (0% tasks done), IN_PROGRESS (>0%), COMPLETED (100%).

**Critério de aceite:** Accordion por fase e épico; dependências visuais; progresso por épico calculado automaticamente.

---

### TASK-09-007 — Exportação do Roadmap

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir exportação do roadmap em múltiplos formatos para uso em ferramentas externas |

**Descrição:**
Implementar `GET /api/v1/projects/:id/roadmap/export` com parâmetro `format` suportando: `json` (o JSON original gerado pelo agente), `csv` (uma linha por task com todos os campos), `markdown` (documento Markdown legível com hierarquia Fase → Épico → Task), `jira` (formato de importação do Jira em CSV). No frontend, botão "Exportar Roadmap" com dropdown de formatos. Download automático do arquivo no formato selecionado.

**Critério de aceite:** Exportação funcional nos 4 formatos; download automático no frontend; CSV válido para Jira.

---

### TASK-09-008 — Timeline de Gantt Simplificada

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Visualizar o roadmap em formato de Gantt para planejamento temporal |

**Descrição:**
Implementar a aba "Timeline" na tela do roadmap com um Gantt simplificado. As tasks são dispostas em ordem de dependências com estimativa de duração proporcional. Assumir que tasks do mesmo épico podem ser paralelas quando não há dependência explícita entre elas, tasks CRITICAL têm cor vermelha, HIGH têm laranja, MEDIUM têm amarelo, LOW têm verde. O Gantt é apenas visual e não impõe datas reais (usa "sprints/blocos de tempo" relativos). Botão para exportar o Gantt como imagem PNG.

**Critério de aceite:** Gantt visual com código de cores; paralelo quando sem dependência; exportação como PNG.

---

### TASK-09-009 — Atualização de Status de Tasks via Drag-and-Drop

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Persistir a atualização de status das tasks quando o usuário move cards no KanBan |

**Descrição:**
Implementar `PUT /api/v1/projects/:id/tasks/:taskId/status` que atualiza o status de uma task (`TODO→IN_PROGRESS→REVIEW→DONE`, ou qualquer status → `BLOCKED`). Validar: não é possível mover para `DONE` se há dependências ainda em `TODO` ou `IN_PROGRESS`, ao mover para `DONE`, verificar se todas as tasks do épico estão concluídas e atualizar status do épico. Registrar histórico de mudanças de status com timestamp para auditoria. No frontend, otimistic update (move o card imediatamente e reverte se a API retornar erro).

**Critério de aceite:** Atualização persistida; validação de dependências; histórico de mudanças; optimistic update no frontend.

---

### TASK-09-010 — Teste de Schema e Ingestão do Roadmap

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que o parser do JSON de roadmap é robusto a variações e erros do output do agente |

**Descrição:**
Testes unitários e de integração do `RoadmapSchemaValidator` e `RoadmapIngester` cobrindo: JSON perfeitamente válido (happy path), JSON com campo obrigatório faltando (deve retornar erro descritivo), JSON com enum inválido (tipo ou complexidade não reconhecido), JSON com referência de dependência inválida (task_id inexistente), JSON com `estimated_hours` inválido (0, negativo, acima de 200), ingestão atômica (simular falha no meio da criação bulk — nenhuma task deve ser criada parcialmente), JSON com aninhamento correto mas arrays vazios.

**Critério de aceite:** Todos os cenários de erro capturados; ingestão atômica validada; mensagens de erro descritivas.

---

# Phase 16 — Painel de Billing e Auditoria de Custos de LLM

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-16 |
| **Título** | Painel de Billing e Auditoria de Custos de LLM |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Média |
| **Pré-requisitos** | PHASE-06 concluída (Pipeline da Tríade com rastreamento de tokens) |

---

## Descrição Detalhada

O Painel de Billing é um dos diferenciais competitivos da plataforma, conforme destacado no PROJECT.md. Ele oferece **transparência total dos custos** gerados pelo uso dos modelos de IA em cada fase, projeto e agente. Isso é especialmente crítico no Modo Dinâmico, onde múltiplos modelos são sorteados aleatoriamente e os preços variam significativamente entre providers.

O sistema de billing não apenas coletará dados — ele os apresentará de forma que o usuário possa: entender exatamente em que está gastando, identificar quais modelos são mais eficientes em custo/qualidade, projetar o custo de futuros projetos baseado no histórico, e, para negócios que revendem o serviço, viabilizar o repasse de custos para clientes finais.

Esta fase implementa a coleta centralizada de dados de billing (já iniciada na Phase 06), os dashboards analíticos e os relatórios exportáveis.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Dashboard de billing com visualizações por projeto, fase, agente e período
- ✅ Tabela de preços atualizada por provider/modelo
- ✅ Projeção de custo antes de iniciar uma fase
- ✅ Alertas de budget por projeto
- ✅ Relatórios exportáveis em CSV e PDF
- ✅ API de billing para integração externa

---

## Funcionalidades Entregues

- **Dashboard Analítico:** Gráficos de custo por projeto, fase, modelo e período
- **Tabela de Preços Real:** Preços atualizados por provider/modelo
- **Projeção de Custo:** Estimativa antes de executar uma fase
- **Alertas de Budget:** Notificações quando custo ultrapassa threshold

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa ⭐ **Recomendado**

Todo o sistema de billing e auditoria de custos é implementado sequencialmente pela Tríade, desde o collector de tokens até o dashboard visual.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a implementação |
| **Velocidade** | Mais rápido — o sistema de billing é coeso e interdependente |
| **Feedback** | Aplicado ao painel de billing como um todo |
| **Ideal para** | Implementar o sistema de billing de ponta a ponta de uma só vez |

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Útil para implementar apenas o coletor de tokens sem o dashboard, ou apenas os alertas de budget.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task |
| **Velocidade** | Mais controlado |
| **Feedback** | Específico para o componente de billing selecionado |
| **Ideal para** | Refinar um componente específico do billing sem reprocessar os demais |

### 🔀 Modo Híbrido

Execute automaticamente as tasks de coleta e persistência de tokens, e revise individualmente o dashboard e os alertas de orçamento que têm impacto direto na experiência do usuário.

> 💡 **Dica:** Para esta phase (Billing e Auditoria de Custos), o modo recomendado é **phase completa** — o sistema de billing é mais robusto quando implementado de forma unificada, com o coletor e o visualizador desenvolvidos em conjunto.

---

## Tasks

### TASK-16-001 — Modelagem da Entidade BillingRecord

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir a entidade de registro de billing com granularidade suficiente para auditoria completa |

**Descrição:**
Definir a struct `BillingRecord` já prevista na PHASE-06 em maior detalhe: `ID`, `ProjectID`, `UserID`, `PhaseNumber`, `PhaseName`, `TriadRole` (PRODUCER/REVIEWER/REFINER), `AgentID`, `AgentName`, `Provider` (OPENAI/ANTHROPIC/GOOGLE/OLLAMA), `Model` (claude-3-opus, gpt-4o, etc.), `PromptTokens`, `CompletionTokens`, `TotalTokens`, `PricePerMillionPromptTokens` (snapshot do preço no momento da chamada), `PricePerMillionCompletionTokens`, `EstimatedCostUSD` (calculado: `(promptTokens/1M * pricePrompt) + (completionTokens/1M * priceCompletion)`), `DurationMs`, `IsAutoRejection` (boolean — se foi acionado por rejeição automática), `Timestamp`. Criar índices: `{user_id, project_id}`, `{user_id, timestamp}`, `{project_id, phase_number}`.

**Critério de aceite:** Struct com granularidade completa; snapshot de preço no momento da chamada; índices para queries de dashboard.

---

### TASK-16-002 — Tabela de Preços de Modelos LLM

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Manter uma tabela atualizada com os preços por 1M de tokens de cada provider/modelo suportado |

**Descrição:**
Implementar o `ModelPricingTable` como configuração YAML (atualizável sem redeploy): listagem de todos os modelos suportados com `prompt_price_per_million_tokens` e `completion_price_per_million_tokens` em USD. Implementar `GET /api/v1/billing/pricing` que retorna a tabela formatada. No frontend, exibir a tabela de preços na aba de configurações com data de última atualização. Implementar uma rotina de atualização de preços (inicialmente manual via YAML, futura via scraping dos sites dos providers). Garantir que o preço snapshottado no BillingRecord é sempre o vigente no momento da chamada.

**Critério de aceite:** Tabela YAML com preços por modelo; endpoint de pricing funcional; snapshot no momento da chamada.

---

### TASK-16-003 — API de Billing Analytics

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar os endpoints de analytics de billing com aggregation pipelines otimizados |

**Descrição:**
Implementar os endpoints de billing analytics via MongoDB aggregation: `GET /api/v1/billing/summary` (custo total do usuário no período, por projeto, por modelo — parâmetros: `from`, `to`), `GET /api/v1/projects/:id/billing` (custo por fase, por agente, por modelo, total do projeto), `GET /api/v1/billing/by-model` (comparativo de custo por modelo — útil para avaliar eficiência), `GET /api/v1/billing/by-phase` (qual fase do pipeline custa mais em média), `GET /api/v1/billing/top-projects` (projetos mais caros do usuário). Todos os endpoints com parâmetros de filtro de período (`from`, `to` em ISO 8601) e paginação.

**Critério de aceite:** Todos os endpoints com aggregation MongoDB correto; filtros de período funcionais; respostas com dados corretos.

---

### TASK-16-004 — Dashboard Principal de Billing no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar o painel de billing com visualizações ricas e interativas |

**Descrição:**
Implementar a página `/billing` com: **Cards de Resumo:** custo total do mês, custo total de todos os projetos, custo médio por projeto, modelo mais caro utilizado; **Gráfico de Linha:** Custo por dia no período selecionado (0-30 dias), com toggle para visualizar por projeto individual; **Gráfico de Pizza:** Distribuição de custo por modelo/provider; **Gráfico de Barras:** Custo por fase do pipeline (qual fase consome mais); **Tabela Detalhada:** Todos os registros de billing com: data, projeto, fase, agente, modelo, tokens (prompt/completion), custo, tipo (normal vs auto-rejection); seletor de período (esta semana, este mês, este ano, personalizado). Usar uma biblioteca de gráficos leve (Recharts ou Chart.js).

**Critério de aceite:** Dashboard com todos os gráficos; seletor de período; tabela com todos os detalhes; cálculos corretos.

---

### TASK-16-005 — Projeção de Custo Antes de Executar uma Fase

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Estimar o custo de uma fase antes de iniciá-la para que o usuário possa tomar uma decisão informada |

**Descrição:**
Implementar o `CostEstimator` que, antes de iniciar uma fase, estima o custo baseado em: tamanho médio do contexto acumulado (tokens estimados do SPEC.md + prompts do usuário + instrução da fase), modelo selecionado (ou modelos candidatos no Modo Dinâmico), complexidade da fase (fases de desenvolvimento custam mais que fases de documentação, baseado em histórico), benchmark de custo de fases similares de projetos anteriores do usuário. Retornar estimativa como range (mínimo-médio-máximo) em USD. No frontend, exibir a estimativa com destaque antes do botão "Iniciar Fase".

**Critério de aceite:** Estimativa como range; baseada em histórico do usuário quando disponível; exibida antes do botão de iniciar.

---

### TASK-16-006 — Sistema de Alertas de Budget por Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Notificar o usuário quando o custo de um projeto ultrapassa thresholds configurados |

**Descrição:**
Implementar o sistema de budget de projeto: o usuário pode configurar um budget máximo por projeto (em USD). O sistema verifica o custo acumulado a cada nova execução da Tríade. Se o custo atingir 80% do budget, envia notificação in-app de alerta "Projeto X atingiu 80% do budget configurado". Se atingir 100%, pausa automaticamente o projeto (com notificação explicativa) e requer confirmação do usuário para continuar. Implementar o endpoint `PUT /api/v1/projects/:id/budget` para configurar o budget. Exibir o gauge de budget visível no painel do projeto.

**Critério de aceite:** Alerta em 80% do budget; pausa automática em 100%; gauge visual no painel do projeto; configuração por projeto.

---

### TASK-16-007 — Exportação de Relatórios de Billing

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir exportação de relatórios de billing em formatos adequados para contabilidade e repasse |

**Descrição:**
Implementar `GET /api/v1/billing/export` com suporte a: **CSV:** uma linha por BillingRecord com todos os campos (para importação em planilhas ou sistemas de contabilidade), **PDF:** relatório formatado com sumário executivo, tabelas de custo por projeto e detalhe de registros (usando biblioteca de geração de PDF), **JSON:** dump completo de todos os registros para integração com sistemas externos. Filtros de export: por projeto, por período, por provider. No frontend, botão de exportar no dashboard de billing com seleção de formato e filtros.

**Critério de aceite:** Exportação em CSV, PDF e JSON; filtros por projeto e período; PDF com formatação profissional.

---

### TASK-16-008 — API de Billing para Integração Externa

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor endpoints de billing com autenticação por API key para integração com sistemas externos de cobrança |

**Descrição:**
Implementar uma API de billing dedicada para integrações externas: `GET /api/v1/external/billing` (protegida por API key ao invés de JWT — para uso em scripts/automações), retorna os mesmos dados de billing mas em formato otimizado para processamento programático, com campos adicionais de `billable_amount_usd` calculado. Implementar geração e gestão de API keys em `POST /api/v1/api-keys` (cria API key com nome descritivo e permissão restrita a billing), `GET /api/v1/api-keys` (lista api keys — sem mostrar o valor completo), `DELETE /api/v1/api-keys/:id` (revoga API key). Ideal para integrações com Stripe, QuickBooks, etc.

**Critério de aceite:** API de billing com auth por API key; gestão de API keys funcional; campo billable_amount_usd calculado.

---

### TASK-16-009 — Comparativo de Custo por Provider/Modelo

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Mostrar ao usuário qual modelo/provider oferece melhor custo-benefício baseado no histórico |

**Descrição:**
Implementar a aba "Análise de Modelos" na página de billing com: tabela comparativa de modelos utilizados com: nome do modelo, total de execuções, custo médio por execução, custo total acumulado, taxa de sucesso da Tríade (aprovação sem rejeição automática), nota de eficiência (custo÷qualidade estimada), gráfico de radar por modelo comparando: custo, velocidade (DurationMs médio), taxa de sucesso, gráfico de tendência do custo por modelo ao longo do tempo, recomendação de modelo mais eficiente baseada nos dados históricos reais do usuário.

**Critério de aceite:** Tabela comparativa com métricas por modelo; gráfico radar; tendência temporal; recomendação baseada em dados reais.

---

### TASK-16-010 — Testes do Sistema de Billing

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir precisão dos cálculos de custo e funcionamento correto dos alertas de budget |

**Descrição:**
Testes unitários e de integração cobrindo: cálculo correto de custo para cada provider (prompt tokens × preço + completion tokens × preço), snapshot de preço no momento da chamada (se preço mudar depois, o registro histórico não muda), alerta de 80% do budget disparado corretamente (mock do custo acumulado), pausa automática do projeto em 100% do budget, projeção de custo: verificar que a estimativa está dentro de 20% do custo real em casos históricos, exports: CSV com todos os campos corretos, PDF gerado sem erros, JSON com schema correto.

**Critério de aceite:** Cálculos de custo corretos por provider; alertas disparados nos thresholds corretos; snapshot de preço imutável; exports validados.

---

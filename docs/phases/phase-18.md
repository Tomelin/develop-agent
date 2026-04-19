# Phase 18 — Observabilidade, Monitoramento e Performance

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-18 |
| **Título** | Observabilidade, Monitoramento e Performance da Plataforma |
| **Tipo** | Backend + DevOps |
| **Prioridade** | Média |
| **Pré-requisitos** | PHASE-01, PHASE-06 concluídas |

---

## Descrição Detalhada

Esta fase implementa a **camada de observabilidade completa** da plataforma da Agência de IA. Uma plataforma de produção que executa pipelines complexos de IA, gerencia múltiplas conexões com providers externos de LLM e processa mensagens assíncronas via RabbitMQ precisa de observabilidade robusta para garantir disponibilidade, detectar gargalos e diagnosticar falhas rapidamente.

A observabilidade aqui é implementada em três pilares: **Métricas** (Prometheus + Grafana), **Logs** (ELK ou Loki) e **Traces** (OpenTelemetry). Além disso, esta fase implementa otimizações de performance identificadas durante o desenvolvimento das fases anteriores: caching estratégico, índices otimizados e compressão de respostas HTTP.

Para a perspectiva do usuário, esta fase também implementa um simples painel de status da plataforma mostrando a saúde dos serviços em tempo real.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Métricas Prometheus expostas para todas as métricas críticas do sistema
- ✅ Distributed tracing com OpenTelemetry para pipelines de IA
- ✅ Dashboard de status da plataforma acessível para usuários
- ✅ Alertas configurados para eventos críticos
- ✅ Otimizações de performance implementadas (caching, índices)
- ✅ Documentação de runbook operacional

---

## Funcionalidades Entregues

- **Métricas de Sistema:** Prometheus com Grafana para visualização
- **Tracing Distribuído:** OpenTelemetry rastreando cada execução da Tríade
- **Status Page:** Página pública de status da plataforma
- **Performance:** Cache Redis, compressão HTTP, índices MongoDB otimizados

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Todo o stack de observabilidade (Prometheus, Grafana, OpenTelemetry, Jaeger) é implementado sequencialmente pela Tríade, sem interrupções.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a implementação |
| **Velocidade** | Mais rápido — stack completo de observabilidade gerado em bloco |
| **Feedback** | Aplicado ao sistema de observabilidade como um todo |
| **Ideal para** | Quem quer configurar tudo de uma vez e revisar o stack completo |

### 🎯 Executar uma Task Específica ⭐ **Recomendado**

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Útil para implementar apenas o Prometheus sem o Grafana, ou apenas o tracing sem as métricas.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por componente de observabilidade |
| **Velocidade** | Mais controlado — cada ferramenta implementada e validada antes de continuar |
| **Feedback** | Específico para o componente selecionado |
| **Ideal para** | Implementar observabilidade de forma incremental, validando cada camada |

### 🔀 Modo Híbrido

Execute automaticamente as métricas básicas (Prometheus + health endpoints) e depois revise individualmente os dashboards Grafana, os traces distribuídos (Jaeger) e os alertas críticos.

> 💡 **Dica:** Para esta phase (Observabilidade), o modo recomendado é **task a task** ou **híbrido**. Cada ferramenta (Prometheus, Grafana, OTel, Jaeger) deve ser validada individualmente em ambiente de staging antes de adicionar a próxima camada — isso evita problemas de configuração em cascata.

---

## Tasks

### TASK-18-001 — Instrumentação com Prometheus

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor métricas Prometheus cobrindo todos os aspectos críticos da plataforma |

**Descrição:**
Instrumentar a aplicação com `github.com/prometheus/client_golang` expondo as métricas em `/metrics`: **Métricas de HTTP:** `http_requests_total{method, path, status}`, `http_request_duration_seconds{method, path}` (histogram), `http_requests_in_flight` (gauge). **Métricas de Pipeline:** `triad_executions_total{phase, status}`, `triad_execution_duration_seconds{phase}` (histogram), `feedback_cycles_total{phase}`, `auto_rejections_total{phase, reason}`. **Métricas de LLM:** `llm_requests_total{provider, model, status}`, `llm_tokens_prompt_total{provider, model}` (counter), `llm_tokens_completion_total{provider, model}` (counter), `llm_request_duration_seconds{provider, model}` (histogram). **Métricas de Filas:** `rabbitmq_messages_pending`, `rabbitmq_messages_processed_total`, `rabbitmq_messages_failed_total`. **Métricas de Banco:** `mongodb_operations_total{operation}`, `redis_cache_hits_total`, `redis_cache_misses_total`.

**Critério de aceite:** Todas as métricas expostas em `/metrics`; labels corretos; histograms com buckets adequados.

---

### TASK-18-002 — Distributed Tracing com OpenTelemetry

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar rastreamento distribuído para visualizar a execução completa de cada pipeline da Tríade |

**Descrição:**
Instrumentar a aplicação com OpenTelemetry SDK para Go (`go.opentelemetry.io/otel`). Criar spans para: cada request HTTP (middleware automático do GIN), cada execução do TriadOrchestrator com sub-spans para Produtor, Revisor e Refinador, cada chamada ao AgentSDK (span com atributos: provider, model, prompt_tokens, completion_tokens), cada operação no MongoDB e Redis, cada publicação/consumo no RabbitMQ. Exportar traces para Jaeger (ou OTLP para qualquer backend compatível). O trace completo de uma execução da Tríade deve mostrar a cascata de spans do request HTTP até cada chamada ao LLM.

**Critério de aceite:** Traces exportados para Jaeger; span completo de uma Tríade visível com sub-spans; atributos de tokens nos spans de LLM.

---

### TASK-18-003 — Dashboard Grafana para a Plataforma

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps |
| **Objetivo** | Criar dashboards Grafana pré-configurados para monitoramento da plataforma |

**Descrição:**
Criar as definições de dashboard Grafana como arquivos JSON (provisionados automaticamente via `docker-compose`): **Dashboard Overview:** RPS, latência P50/P95/P99, taxa de erro HTTP, projetos ativos, fases em execução. **Dashboard de Pipeline:** taxa de conclusão da Tríade, tempo médio por fase, taxa de rejeição automática, feedbacks por fase. **Dashboard de LLM:** custo por provider em tempo real, tokens por minuto por provider, taxa de erro por provider, latência por modelo. **Dashboard de Infraestrutura:** CPU/memória por container, tamanho da fila RabbitMQ, hit rate do Redis, operações MongoDB por segundo.

**Critério de aceite:** 4 dashboards criados; provisionados automaticamente via docker-compose; métricas corretas em cada painel.

---

### TASK-18-004 — Alertas e Notificações Operacionais

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps |
| **Objetivo** | Configurar alertas automáticos para condições operacionais críticas |

**Descrição:**
Configurar alertas Prometheus (via Alertmanager ou Grafana Alerts) para: **Crítico:** taxa de erro HTTP > 5% por 5 minutos, fila RabbitMQ > 100 mensagens por mais de 10 minutos, qualquer provider de LLM com 100% de falha por 5 minutos, uso de memória > 90%. **Aviso:** latência P95 > 2 segundos, tag de redis hit rate < 50%, custo de LLM no dia acima de threshold configurado. Notificações para: Slack webhook (se configurado), email (se configurado), e registro de alerta no banco de dados para visualização no painel de admin.

**Critério de aceite:** Alertas configurados com thresholds corretos; notificações por Slack/email; log de alertas no banco.

---

### TASK-18-005 — Otimização de Cache Redis para Contextos Frequentes

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar caching estratégico para reduzir latência e custo de operações frequentes |

**Descrição:**
Identificar e implementar caching Redis para: **Prompts do Usuário por Grupo:** cache de 5 minutos para a lista de prompts de um usuário por grupo (evita queries MongoDB a cada execução de fase), **Agentes por Skill:** cache de 10 minutos para lista de agentes habilitados por skill (evita queries no início de cada sorteio), **Pricing Table:** cache de 24 horas para a tabela de preços de modelos, **SPEC.md do Projeto:** cache de 15 minutos para o SPEC.md de cada fase (muito acessado na composição de prompts), **Profile do Usuário:** cache de 30 minutos. Implementar invalidação de cache: ao atualizar prompts/agentes, invalida o cache correspondente via evento Redis `DEL`.

**Critério de aceite:** Cache implementado para todos os recursos listados; invalidação funcional; hit rate measurável via metrics.

---

### TASK-18-006 — Revisão e Otimização de Índices MongoDB

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Revisar e otimizar todos os índices MongoDB para garantir performance em escala |

**Descrição:**
Executar análise de `explain()` nas principais queries do sistema e otimizar índices: **BillingRecords:** índice composto `{user_id, timestamp}` para queries de dashboard de billing com range de datas, **Projects:** índice composto `{owner_id, status, updated_at}` para dashboard do usuário com filtros, **Tasks:** índice composto `{project_id, status, type}` para o KanBan com filtros, **TriadExecution:** índice em `{project_id, phase_number, status}` para queries de acompanhamento, **BillingRecords agregation:** adicionar índice em `{provider, model, timestamp}` para o dashboard comparativo de modelos. Verificar que nenhum índice está duplicado (overhead de escrita desnecessário). Documentar todos os índices e suas justificativas.

**Critério de aceite:** explain() mostrando IXSCAN (não COLLSCAN) nas queries críticas; sem índices duplicados; documentação dos índices.

---

### TASK-18-007 — Compressão e Otimização de Respostas HTTP

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Reduzir o tamanho das respostas HTTP para melhorar performance para o usuário final |

**Descrição:**
Implementar middleware de compressão gzip/br no GIN para respostas > 1KB. Implementar paginação eficiente em todas as listagens (cursor-based pagination para coleções grandes, ao invés de offset que degrada com volume alto). Implementar projeção explícita em todas as queries MongoDB (nunca retornar campos desnecessários). Adicionar headers de cache HTTP corretos: `Cache-Control: no-cache` para dados dinâmicos, `Cache-Control: max-age=3600` para dados estáticos (tabela de preços, lista de templates). Implementar ETag para respostas de detalhe de projeto (evitar re-download se dados não mudaram).

**Critério de aceite:** Gzip funcionando (verificar via Accept-Encoding header); paginação cursor-based; ETags em endpoints de detalhe.

---

### TASK-18-008 — Página de Status da Plataforma (Status Page)

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar uma página pública de status da plataforma para comunicar disponibilidade dos serviços |

**Descrição:**
Implementar a página `/status` (pública, sem autenticação) com: status geral da plataforma (OPERATIONAL/DEGRADED/OUTAGE), status individual de cada componente: API Backend, Workers de IA, MongoDB, Redis, RabbitMQ, providers de LLM (OpenAI, Anthropic, Google, Ollama), histórico de incidentes dos últimos 7 dias com título, data, duração e resolução, uptime dos últimos 30 dias por componente. Atualização automática a cada 30 segundos via polling (ou SSE). Design minimalista e de carregamento muito rápido. Endpoint backend: `GET /api/v1/status` (público, sem auth) que agrega todos os health checks.

**Critério de aceite:** Status page pública sem auth; componentes individuais monitorados; histórico de incidentes; uptime percentual.

---

### TASK-18-009 — Runbook Operacional

| Campo | Valor |
|-------|-------|
| **Camada** | Documentação |
| **Objetivo** | Documentar os procedimentos operacionais para situações comuns de manutenção e incidente |

**Descrição:**
Criar o `docs/RUNBOOK.md` com procedimentos para: **Deploy:** procedimento de deploy zero-downtime, rollback de versão, atualização de variáveis de ambiente sem downtime, **Escalonamento:** como aumentar número de workers, como escalar MongoDB (replica set), como aumentar memória do Redis, **Incidentes Comuns:** fila RabbitMQ travada (procedimento de recovery), worker preso em execução (timeout manual), provider de LLM com rate limit (como redirecionar para outro provider temporariamente), **Manutenção:** procedimento de backup do MongoDB, procedimento de limpeza de billing records antigos (data retention), como atualizar a tabela de preços de modelos.

**Critério de aceite:** Runbook com todos os procedimentos listados; passos claros e testados; links para dashboards relevantes em cada procedimento.

---

### TASK-18-010 — Load Testing e Baseline de Performance

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps |
| **Objetivo** | Estabelecer o baseline de performance da plataforma e identificar gargalos com load testing |

**Descrição:**
Criar scripts de load testing com `k6` cobrindo: **Cenário 1 - Dashboard:** 100 usuários concorrentes navegando no dashboard por 5 minutos, **Cenário 2 - Iniciar Fase:** 20 usuários concorrentes iniciando fases simultâneas (satura os workers), **Cenário 3 - Billing Dashboard:** 50 usuários consultando o dashboard de billing com diferentes períodos. Métricas coletadas: P50, P95, P99 de latência, taxa de erro, throughput (req/s). Executar os testes em ambiente de staging e documentar o baseline. Identificar e documentar os gargalos encontrados (ex: "o endpoint /billing/summary degrada com > 50 usuários concorrentes — adicionar cache").

**Critério de aceite:** Scripts k6 criados para os 3 cenários; baseline documentado; gargalos identificados e documentados.

---

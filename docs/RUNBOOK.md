# RUNBOOK Operacional — Observabilidade e Incidentes

## 1) Deploy e rollback (zero-downtime)
1. Faça build da nova imagem backend.
2. Suba nova versão em paralelo (rolling update).
3. Valide `/health/ready` e `/metrics` da nova versão.
4. Direcione tráfego gradualmente para os novos pods/containers.
5. Em erro, execute rollback para a imagem anterior e invalide cache Redis das chaves críticas.

## 2) Escalonamento
### Workers IA
1. Aumente o número de réplicas do worker.
2. Monitore `rabbitmq_messages_pending` e latência P95.

### MongoDB
1. Escale em replica set.
2. Revalide índices e `explain()` das queries críticas.

### Redis
1. Aumente memória/limites da instância.
2. Monitore hit rate (`redis_cache_hits_total` / misses).

## 3) Incidentes comuns
### Fila RabbitMQ travada
1. Verifique conexão e consumidores ativos.
2. Inspecione DLQ e mensagens com falha.
3. Reinicie consumidores presos e reprocese mensagens idempotentes.

### Worker preso em execução
1. Localize `execution_id` no log estruturado.
2. Aplique timeout manual no job e marque como erro controlado.
3. Reenfileire execução se seguro.

### Provider LLM com falha/rate-limit
1. Valide `llm_requests_total{status="error"}` por provider.
2. Redirecione temporariamente para provider alternativo via configuração.
3. Reavalie custo/latência após estabilização.

## 4) Manutenção
### Backup MongoDB
1. Execute backup full diário.
2. Teste restore em ambiente de staging semanalmente.

### Retenção de billing
1. Arquive registros antigos em storage frio.
2. Remova dados fora da janela de retenção com job agendado.

### Atualização de tabela de preços
1. Atualize `src/backend/config/model_pricing.yaml`.
2. Recarregue serviço de billing.
3. Valide dashboard de custo por provider.

## 5) Dashboards mínimos para operação
- HTTP: RPS, erro (%), latência P50/P95/P99.
- LLM: erros por provider/modelo, tokens por minuto, custo diário.
- Filas: pendências, throughput, falhas.
- Banco/cache: operações Mongo, latência Redis, hit rate.

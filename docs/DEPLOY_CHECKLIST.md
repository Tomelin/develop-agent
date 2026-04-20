# Deploy Checklist — Produção

## Pré-deploy
- [ ] CI principal verde (testes unitários, integração, contrato e E2E).
- [ ] Validação manual em staging dos fluxos críticos (A, B e C).
- [ ] Backup do MongoDB de produção concluído e testado.
- [ ] Janela de manutenção comunicada (quando aplicável).

## Deploy
- [ ] `docker compose pull` executado sem falhas.
- [ ] `docker compose up -d` aplicado no stack de produção.
- [ ] Health checks (`/healthz` e `/api/v1/status`) retornando OK.
- [ ] Logs dos primeiros 5 minutos sem erros críticos.

## Pós-deploy
- [ ] Smoke test: login, criação de projeto e início de fase.
- [ ] Métricas de erro/latência sem spikes.
- [ ] Alertas operacionais não disparados.
- [ ] Comunicação de conclusão enviada para o time.

## Rollback (meta: < 5 minutos)
- [ ] Reverter imagem para tag estável anterior.
- [ ] Restaurar compose stack da versão anterior.
- [ ] Confirmar health checks e smoke test mínimo.
- [ ] Registrar incidente e análise pós-mortem.

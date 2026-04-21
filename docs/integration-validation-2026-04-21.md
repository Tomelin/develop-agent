# Validação de integração Backend x Frontend (2026-04-21)

Escopo validado:
- Backend: `src/backend/api/handler/*`
- Frontend: `src/frontend/src/services/*`
- Contratos de request/response usados pelos serviços.

## ✅ Paths alinhados
- Organização: `/org/members`, `/org/invite`, `/org/members/:userId/role`, `/org/members/:userId`.
- Projetos principais: `/projects`, `/projects/:id`, `/projects/:id/pause`, `/projects/:id/resume`, `/projects/:id/archive`.
- Fases implementadas: `5`, `6`, `7`, `13`, `14`, `15`, `19`.
- Billing core: `/billing/pricing`, `/billing/summary`, `/billing/by-model`, `/billing/by-phase`, `/billing/top-projects`, `/billing/export`.
- Interview core: `/projects/:id/interview`, `/projects/:id/interview/message`, `/projects/:id/interview/confirm`, `/projects/:id/interview/regenerate-vision`.

## ❌ Paths do frontend sem handler no backend (quebram hoje)

### Phase 8
- `GET /projects/:id/phases/:phaseNumber/artifacts`
- `GET /projects/:id/phases/:phaseNumber/triad-progress`
- `GET /projects/:id/phases/:phaseNumber/tracks/:track/feedbacks`
- `POST /projects/:id/phases/:phaseNumber/tracks/:track/feedback`
- `GET /notifications`
- `POST /notifications/:id/read`

### Phase 17 (admin/config dinâmica)
- `GET /projects/:id/triad-selections`
- `GET /projects/:id/selection-logs`
- `PUT /projects/:id/dynamic-mode`
- `GET /projects/:id/dynamic-mode/preview`
- `GET /projects/:id/diversity-metrics`
- `GET /projects/:id/agent-config/matrix`
- `PUT /projects/:id/agent-config/matrix`
- `POST /projects/:id/agent-config/cost-preview`
- `GET /admin/settings`
- `PUT /admin/settings`
- `GET /admin/feature-flags`
- `PUT /admin/feature-flags`
- `GET /feature-flags`

### Phase 20 (colaboração/marketplace/integrações/pricing público)
- `GET /projects/:id/collaborators`
- `POST /projects/:id/collaborators`
- `PUT /projects/:id/collaborators/:userId/role`
- `DELETE /projects/:id/collaborators/:userId`
- `GET /marketplace/templates`
- `POST /marketplace/templates`
- `POST /marketplace/templates/:templateId/use`
- `POST /marketplace/templates/:templateId/star`
- `GET /projects/:id/integrations`
- `GET /integrations/github/auth`
- `POST /integrations/jira`
- `POST /projects/:id/integrations/jira/sync`
- `POST /integrations/slack/webhook`
- `GET /pricing/plans`
- `POST /pricing/checkout`
- `GET /roadmap/public`
- `POST /roadmap/features/:featureId/vote`
- `POST /roadmap/features/suggestions`
- `PUT /admin/roadmap/features/:featureId/status`

## ⚠️ Inconsistências de contrato (payload/query/response)

1) **Paginação de projetos**
- Frontend envia `page` + `size`.
- Backend lê `page` + `limit`.
- Efeito: `size` é ignorado, backend usa default de `limit=20`.

2) **Resposta de update de status de task**
- Frontend espera `Task` no retorno de `PUT /projects/:id/tasks/:taskId/status`.
- Backend retorna `204 No Content`.
- Efeito: caller recebe `undefined` onde tipagem diz `Task`.

3) **Roadmap detalhado**
- Frontend chama `GET /projects/:id/roadmap`.
- Backend expõe `GET /projects/:id/roadmap/summary` e `GET /projects/:id/roadmap/export`, mas não `/roadmap`.

## Recomendação prática
1. Implementar handlers faltantes (ou ocultar features no frontend até backend existir).
2. Padronizar paginação (`size` vs `limit`) em ambos os lados.
3. Ajustar contrato de `updateTaskStatus`:
   - opção A: backend retornar task atualizada;
   - opção B: frontend tipar retorno como `void`.
4. Para roadmap, alinhar endpoint único (`/roadmap`) ou ajustar frontend para usar apenas os endpoints já existentes.

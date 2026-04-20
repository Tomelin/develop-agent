# Phase 20 — Status de Implementação (Incremental)

## Escopo implementado nesta entrega

Esta entrega cobre a base técnica da **TASK-20-001 (Multi-Tenancy por Organização)** para permitir isolamento por organização no backend:

1. **Novo domínio de organização**
   - Entidade `Organization` com `Plan` e metadados principais.
2. **Propagação de `organization_id`**
   - Campos adicionados em `User`, `Project` e `BillingRecord`.
3. **Contexto de organização no JWT**
   - Token de acesso passa a carregar `organization_id`.
4. **Middlewares de organização**
   - Resolução de organização por JWT (com override opcional por header `X-Organization-ID`).
5. **Filtragem por tenant**
   - Filtros de listagem de projetos e billing passam a exigir `organization_id`.

## Decisões de implementação

- Estratégia inicial de isolamento: **single database + `organization_id` em documentos**.
- Escopo incremental: priorização da infraestrutura de autenticação/consulta, sem concluir ainda gestão de membros, convites, RBAC completo por organização e integrações externas da Phase 20.

## Próximos passos recomendados

1. Implementar repositório e endpoints de `Organization`.
2. Implementar `OrganizationMember` e RBAC (OWNER/ADMIN/MEMBER/VIEWER).
3. Migrar consultas de leitura por ID para receber `organization_id` explicitamente no repositório (evitar validação posterior em handler).
4. Adicionar migração/backfill para usuários/projetos legados sem `organization_id`.
5. Completar TASK-20-002 até TASK-20-010 de forma incremental.

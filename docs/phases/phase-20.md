# Phase 20 — Evolução, Multi-Tenancy e Roadmap Futuro

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-20 |
| **Título** | Evolução da Plataforma: Multi-Tenancy, Escalabilidade e Roadmap Futuro |
| **Tipo** | Backend + Frontend + Estratégia |
| **Prioridade** | Baixa (Futura) |
| **Pré-requisitos** | PHASE-01 a PHASE-19 concluídas |

---

## Descrição Detalhada

Esta é a fase de **evolução estratégica** da plataforma para escala de mercado institucional. Com todas as funcionalidades principais implementadas e validadas nas fases anteriores, é hora de preparar a plataforma para crescimento: suporte a múltiplos usuários/organizações (multi-tenancy), recursos de colaboração em equipe, mercado de agentes e templates, e integrações com ferramentas populares do ecossistema de desenvolvimento.

Esta fase também consolida o **Roadmap de Produto** formal da plataforma — um documento vivo que guia a evolução da Agência de IA além do MVP inicial. Inclui pesquisa de demanda de usuários, priorização de features e planejamento de versões futuras.

Além das features de produto, esta fase implementa melhorias de escalabilidade arquitetural: migração de partes do sistema para microserviços quando necessário, implementação de multi-region, e suporte a deploys em Kubernetes com alta disponibilidade.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Suporte básico a múltiplas organizações (multi-tenant)
- ✅ Sistema de colaboração em equipe (múltiplos usuários por projeto)
- ✅ Marketplace de templates de prompts e agentes
- ✅ Integrações com GitHub, Jira e Slack
- ✅ Roadmap de produto formalizado e público
- ✅ Plano de precificação e monetização

---

## Funcionalidades Entregues

- **Multi-Tenancy:** Organizações com múltiplos usuários e projetos compartilhados
- **Collaboração:** Múltiplos usuários trabalhando no mesmo projeto
- **Marketplace:** Compartilhamento de configurações de agentes e templates
- **Integrações:** GitHub, Jira, Slack como saídas do pipeline

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Todas as features de evolução da plataforma são implementadas sequencialmente pela Tríade, cobrindo multi-tenancy, integrações e marketplace.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a fase de evolução |
| **Velocidade** | Mais rápido — features de evolução implementadas em bloco |
| **Feedback** | Aplicado ao conjunto de features de evolução como um todo |
| **Ideal para** | Times que querem executar toda a evolução da plataforma de uma vez |

### 🎯 Executar uma Task Específica ⭐ **Recomendado**

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Cada task desta phase representa uma grande feature independente (ex: apenas a integração com GitHub, sem o Jira).

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por feature — cada integração validada antes de continuar |
| **Velocidade** | Mais controlado — integrações complexas exigem validação cuidadosa |
| **Feedback** | Específico para a feature ou integração selecionada |
| **Ideal para** | Features complexas de integração (GitHub, Jira, Slack, Stripe) que devem ser validadas individualmente |

### 🔀 Modo Híbrido

Execute automaticamente o modelo de multi-tenancy e gestão de membros, e revise individualmente cada integração externa (GitHub OAuth, Jira API, Slack Webhooks, Stripe Checkout) antes de aprovar.

> ⭐ **Fortemente Recomendado:** Use **task a task** nesta phase. Cada task representa uma feature estratégica de alto impacto. Integrações externas (GitHub, Jira, Slack, Stripe) devem ser implementadas, testadas e validadas individualmente antes de avançar para a próxima — erros em integrações OAuth têm impacto imediato na experiência do usuário final.

---

## Tasks

### TASK-20-001 — Modelo de Multi-Tenancy por Organização

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a estrutura de dados e isolamento para múltiplas organizações na plataforma |

**Descrição:**
Implementar a entidade `Organization` com: `ID`, `Name`, `Slug` (único, usado em URLs), `Plan` (FREE/STARTER/PRO/ENTERPRISE), `MaxUsers`, `MaxProjectsPerMonth`, `MaxTokensPerMonth`, `BillingEmail`, `CreatedAt`. Adicionar `OrganizationID` em todas as entidades que requerem isolamento por tenant: User (membro de uma organização), Project, Agent catalog (cada organização tem seu próprio catálogo de agentes, além dos globais), BillingRecord. Atualizar todos os repositórios e handlers para filtrar por `organization_id` automaticamente. Implementar um middleware de resolução de organização que extrai a org do subdomínio ou do JWT do usuário.

**Critério de aceite:** Toda query filtrada por org_id automaticamente; isolamento completo entre organizações; middleware de resolução de org.

---

### TASK-20-002 — Gestão de Membros da Organização

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Implementar o sistema de membros e roles dentro de uma organização |

**Descrição:**
Implementar o modelo de membros com roles: `OWNER` (cria e gerencia a org, acesso total), `ADMIN` (gerencia usuários e configurações, acesso total exceto billing), `MEMBER` (cria projetos e usa a plataforma), `VIEWER` (apenas visualiza projetos e resultados). Endpoints: `POST /api/v1/org/invite` (convida usuário por email), `GET /api/v1/org/members` (lista membros), `PUT /api/v1/org/members/:userId/role` (altera role), `DELETE /api/v1/org/members/:userId` (remove da org). Atualizar o sistema de autorização para considerar a role na organização além da role no sistema.

**Critério de aceite:** RBAC por organização funcional; convite por email implementado; roles com permissões distintas validadas.

---

### TASK-20-003 — Colaboração em Projeto (Múltiplos Usuários)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir que múltiplos membros da organização colaborem no mesmo projeto |

**Descrição:**
Expandir o modelo de projeto para suportar colaboradores: `Project.Members` como lista de `{UserID, Role: OWNER|EDITOR|VIEWER}`. O dono do projeto pode adicionar/remover colaboradores. Permissões: OWNER (todas as ações), EDITOR (pode iniciar fases, enviar feedback, aprovar, mas não arquivar ou deletar), VIEWER (apenas visualiza artefatos e progresso). Implementar `POST /api/v1/projects/:id/collaborators` para adicionar colaborador. No frontend, seção "Equipe" no painel do projeto com lista de colaboradores, roles e botões de gestão.

**Critério de aceite:** Colaboradores com roles distintas; permissões validadas no backend; UI de gestão de equipe no painel do projeto.

---

### TASK-20-004 — Marketplace de Templates de Prompts

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Criar um marketplace onde usuários podem compartilhar e reutilizar templates de prompts |

**Descrição:**
Implementar o marketplace de templates públicos: **Backend:** entidade `PublicTemplate` com `ID`, `Title`, `Description`, `Category`, `Content`, `Group` (fase alvo), `Stars` (contador), `UsageCount` (quantas vezes foi usado), `CreatorID`, `Visibility` (PUBLIC/PRIVATE), `Tags`. Endpoints: `GET /api/v1/marketplace/templates` (listagem pública com filtros e busca), `POST /api/v1/marketplace/templates` (publicar template), `POST /api/v1/marketplace/templates/:id/use` (adicionar ao meu banco de prompts). **Frontend:** página `/marketplace` com grid de templates, filtros por categoria e fase alvo, preview do conteúdo, botão "Usar este Template", contador de uso e estrelas.

**Critério de aceite:** Marketplace público sem auth para visualização; uso de template cria prompt no perfil do usuário; contador de uso atualizado.

---

### TASK-20-005 — Integração com GitHub

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Integrar a plataforma com GitHub para push automático do código gerado |

**Descrição:**
Implementar a integração OAuth com GitHub: **OAuth Flow:** `GET /api/v1/integrations/github/auth` (redireciona para GitHub OAuth), `GET /api/v1/integrations/github/callback` (troca o code por access token e armazena). **Repositório:** ao completar a Fase 5 (ou qualquer fase que gera código), opção de "Push para GitHub" que cria/atualiza um repositório no GitHub do usuário com o código gerado. **PR Automático:** ao concluir a Fase 5, opção de criar um Pull Request no repositório com as mudanças geradas. Configurável por projeto: repositório de destino, branch. No frontend: seção "Integrações" no painel do projeto com card do GitHub.

**Critério de aceite:** OAuth com GitHub funcional; push de código automático; criação de PR válido no GitHub.

---

### TASK-20-006 — Integração com Jira

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Sincronizar as tasks do KanBan geradas na Fase 4 com projetos do Jira |

**Descrição:**
Implementar integração com Jira Cloud via API REST: configuração de credenciais (Jira URL, email, API token) por organização em `/api/v1/integrations/jira`. Ao concluir a Fase 4 (Planejamento), opção de "Sincronizar com Jira" que: cria Épicos no Jira para cada épico do roadmap, cria Stories no Jira para cada task do roadmap com todos os metadados (tipo, complexidade mapeada para story points, descrição), atualiza o status das tasks no Jira quando o status muda no KanBan da plataforma (bidirecional opcional). No frontend: toggle de sincronização com Jira no painel da Fase 4.

**Critério de aceite:** Epics e Stories criados no Jira com dados corretos; sincronização de status bidirecional (opcional); configuração por organização.

---

### TASK-20-007 — Integração com Slack (Notificações)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Enviar notificações de eventos do pipeline para canais Slack configurados |

**Descrição:**
Implementar integração com Slack via Incoming Webhooks: `POST /api/v1/integrations/slack/webhook` (configura o webhook URL e o canal padrão). Eventos que geram notificações Slack: fase concluída e aguardando feedback (com link para o painel), fase aprovada e próxima fase iniciada, projeto concluído (com resumo de métricas), rejeição automática disparada (com motivo), alerta de budget atingido. Formato das mensagens: Slack Block Kit com formatação rich — ícones, cores por tipo de evento (verde = sucesso, amarelo = atenção, vermelho = falha), botão de link para o painel.

**Critério de aceite:** Notificações com Block Kit por tipo de evento; link para o painel em cada mensagem; configuração por organização/projeto.

---

### TASK-20-008 — Plano de Precificação e Monetização

| Campo | Valor |
|-------|-------|
| **Camada** | Estratégia + Frontend |
| **Objetivo** | Definir e implementar os planos de precificação da plataforma para monetização |

**Descrição:**
Definir e implementar os planos de assinatura: **FREE:** 1 usuário, 2 projetos/mês, apenas Fluxo A (até Fase 4), sem Modo Dinâmico, 10.000 tokens/mês. **STARTER:** 1 usuário, 10 projetos/mês, todos os fluxos, Modo Dinâmico, 100.000 tokens/mês. **PRO:** 5 usuários, projetos ilimitados, todos os fluxos, Modo Dinâmico, 500.000 tokens/mês, suporte prioritário. **ENTERPRISE:** usuários ilimitados, tudo do PRO, SLA, suporte dedicado, deploy on-premises disponível. Implementar integração básica com Stripe para checkout de planos. No frontend: página `/pricing` pública com tabela comparativa de planos, botão de upgrade no painel quando o usuário atinge limites.

**Critério de aceite:** 4 planos definidos com limites claros; integração Stripe para checkout; limites enforçados no backend; página de pricing pública.

---

### TASK-20-009 — Roadmap Público de Produto

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend + Estratégia |
| **Objetivo** | Publicar o roadmap de produto da plataforma para transparência e engajamento da comunidade |

**Descrição:**
Criar a página pública `/roadmap` com: visão de produto de longo prazo (12-18 meses), features agrupadas por milestone (v1.0, v1.5, v2.0), status de cada feature (Planejado/Em Desenvolvimento/Concluído), opção de voto nas features (usuários logados podem votar nas features mais desejadas), formulário de sugestão de nova feature, changelog das versões já lançadas. Os dados do roadmap são gerenciados via admin panel (`/admin/roadmap`). Exibição pública sem necessidade de login para leitura, login necessário apenas para voto/sugestão.

**Critério de aceite:** Roadmap público com milestone; votação de features funcional; changelog de versões; admin panel para gestão.

---

### TASK-20-010 — Documentação de Arquitetura Multi-Tenant e Guia de Contribuição

| Campo | Valor |
|-------|-------|
| **Camada** | Documentação |
| **Objetivo** | Documentar as decisões de arquitetura multi-tenant e criar guia para contribuidores externos |

**Descrição:**
Criar documentação completa para a fase de evolução: **`docs/adr/ADR-MULTI-TENANCY.md`:** documentar as decisões arquiteturais tomadas para multi-tenancy (isolamento por `org_id` vs schemas separados, shared database vs database-per-tenant, motivos da escolha). **`CONTRIBUTING.md`:** guia completo para contribuidores externos: configuração do ambiente, convenções de código, processo de PR, como rodar testes, como adicionar um novo provider de LLM, como adicionar novos campos ao schema de agente. **`docs/ARCHITECTURE.md`:** diagrama C4 (Context, Container, Component) da plataforma completa, incluindo a evolução para multi-tenancy.

**Critério de aceite:** ADR de multi-tenancy documentado com tradeoffs; CONTRIBUTING.md com passo-a-passo de setup; diagrama C4 completo e atualizado.

---

# Phase 19 — Testes de Ponta a Ponta e Qualidade da Plataforma

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-19 |
| **Título** | Testes de Ponta a Ponta e Qualidade da Plataforma |
| **Tipo** | Backend + Frontend + DevOps |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-01 a PHASE-17 concluídas |

---

## Descrição Detalhada

Esta fase é dedicada exclusivamente à **qualidade e confiabilidade** da plataforma da Agência de IA como um todo. Enquanto as fases anteriores incluíam testes unitários e de integração específicos por módulo, aqui o foco é na validação end-to-end dos fluxos completos — desde a criação de um projeto até a entrega do software completo.

Além dos testes automatizados, esta fase inclui a definição formal dos critérios de aceite de cada fluxo, a criação de dados de teste representativos, e a implementação de um ambiente de staging completo que espelha a produção.

A qualidade de uma plataforma de IA é complexa de medir — o output dos agentes tem variabilidade inerente. Por isso, esta fase também implementa um sistema de avaliação de qualidade dos outputs da Tríade usando um "juiz" LLM que pontua a qualidade técnica do que foi gerado.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Suite de testes E2E cobrindo os 3 fluxos principais
- ✅ Ambiente de staging completo com dados representativos
- ✅ Sistema de avaliação de qualidade dos outputs da Tríade (LLM Judge)
- ✅ Relatório de qualidade da plataforma com métricas objetivas
- ✅ Checklist de deploy para produção
- ✅ Documentação de troubleshooting para os cenários mais comuns

---

## Funcionalidades Entregues

- **Testes E2E:** Cobertura completa dos 3 fluxos (A, B, C)
- **LLM Judge:** Avaliação de qualidade dos outputs usando IA
- **Staging Environment:** Ambiente idêntico a produção para validação
- **Deploy Checklist:** Procedimento seguro de go-live

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Toda a suite de testes E2E e o sistema de qualidade de plataforma são implementados sequencialmente pela Tríade, sem interrupções.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a phase |
| **Velocidade** | Mais rápido — suite de QA gerada de ponta a ponta |
| **Feedback** | Aplicado ao conjunto completo de testes e métricas |
| **Ideal para** | Quem quer cobrir toda a suite de qualidade de uma vez |

### 🎯 Executar uma Task Específica ⭐ **Recomendado**

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Útil para executar apenas os testes E2E do Fluxo A sem o Fluxo B, ou apenas o LLM Judge sem os testes de carga.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por tipo de teste/validação |
| **Velocidade** | Mais controlado — cada tipo de teste validado antes de continuar |
| **Feedback** | Específico para o fluxo ou componente sendo testado |
| **Ideal para** | Testar um fluxo específico ou componente isolado de qualidade |

### 🔀 Modo Híbrido

Execute automaticamente os testes de contrato e E2E, e revise individualmente o LLM Judge e o relatório de qualidade que exigem interpretação humana dos resultados.

> 💡 **Dica:** Para esta phase (Testes E2E e Qualidade de Plataforma), o modo **task a task** é recomendado para as tasks que envolvem análise de resultados (LLM Judge, Relatório de Qualidade), e **phase completa** para as tasks de setup de ambiente e fixtures de dados.

---

## Tasks

### TASK-19-001 — Suite de Testes E2E para o Fluxo A Completo

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Validar o fluxo completo de desenvolvimento de software do login até o download do código |

**Descrição:**
Implementar suite de testes E2E usando Playwright (ou Cypress) que simula um usuário real executando o Fluxo A completo. O teste usa mocks de LLM (não chamadas reais) para ser determinístico e rápido. Fluxo testado: login → dashboard → criar projeto Fluxo A → configurar Modo Dinâmico → iniciar Fase 1 (entrevista com 3 mensagens mock) → confirmar visão → aguardar Fase 2 (Tríade mock) → aprovar → aguardar Fase 3 → aprovar → ver KanBan gerado (Fase 4) → iniciar Fase 5 (desenvolvimento task-by-task com código mock) → ver resultado de testes (Fase 6 com cobertura mock) → ver resultado de segurança (Fase 7 mock) → baixar código final. Cada step verificado com assertions robustas.

**Critério de aceite:** Fluxo A completo testado com mocks; cada step com assertions; teste determinístico e executável em CI.

---

### TASK-19-002 — Suite de Testes E2E para Fluxos B e C

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Validar os fluxos de Landing Page e Marketing de ponta a ponta |

**Descrição:**
Implementar testes E2E para: **Fluxo B (Landing Page):** criar projeto Fluxo B → preencher brief manual → iniciar Tríade (mock) → ver preview da landing page no iframe → enviar feedback → aprovar → download do HTML ZIP → verificar score de conversão exibido. **Fluxo B com herança:** criar projeto → vincular ao projeto do Fluxo A existente → verificar que o brief foi auto-preenchido com dados do projeto base. **Fluxo C (Marketing):** criar projeto Fluxo C → preencher brief → iniciar Tríade → ver calendário editorial gerado → download do pack de conteúdo.

**Critério de aceite:** Fluxos B e C testados com mocks; herança de projeto verificada; downloads funcionais no teste.

---

### TASK-19-003 — Testes de Contrato da API (Contract Testing)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Garantir que o frontend e backend nunca ficam dessincronizados em relação ao schema da API |

**Descrição:**
Implementar contract testing com Pact.io: o frontend define os contratos (quais endpoints usa e com quais schemas de request/response), o backend valida que satisfaz todos os contratos do frontend. Contratos críticos para definir: schema de autenticação (login, refresh, me), schema de criação de projeto, schema de listagem de projetos com paginação, schema de eventos SSE do pipeline, schema de billing summary. Integrar a verificação de contratos no CI — se o backend mudar a API de forma breaking sem atualização do contrato, o CI falha.

**Critério de aceite:** Contratos definidos para endpoints críticos; verificação automática no CI; breaking changes detectadas automatically.

---

### TASK-19-004 — LLM Judge — Avaliação de Qualidade dos Outputs da Tríade

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar um sistema de avaliação automática da qualidade dos artefatos gerados pelos agentes |

**Descrição:**
Implementar o `LLMJudge` — um agente avaliador independente (usando um modelo diferente do usado na Tríade) que avalia a qualidade do output final do Refinador em cada fase. O Judge analisa e pontua (0-10) cada artefato por critérios específicos: **Código (Fase 5):** corretude lógica, aderência à arquitetura, style guide, documentação. **Testes (Fase 6):** cobertura real vs obtida, qualidade das assertions, cenários cobertos. **Documentação (Fase 8):** completude, precisão, clareza, exemplos funcionais. O score do Judge é armazenado e exibido ao usuário como indicativo de qualidade (com disclaimer de que é uma avaliação automatizada). Agregar scores por projeto e por fase para métricas de qualidade da plataforma.

**Critério de aceite:** Judge avalia artefatos de cada fase com critérios específicos; score armazenado e exibido; agregação por projeto.

---

### TASK-19-005 — Ambiente de Staging Completo

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps |
| **Objetivo** | Criar um ambiente de staging que espelha a produção para validação antes de cada release |

**Descrição:**
Configurar o ambiente de staging usando docker-compose em um servidor de staging com: todos os serviços de infraestrutura (MongoDB, Redis, RabbitMQ) em configuração similar à produção, backend e frontend com as imagens de produção (não de desenvolvimento), dados de exemplo pré-carregados (usuário admin, 3 projetos em diferentes estados, agentes configurados), variáveis de ambiente de staging com API keys de teste dos providers (quando disponível), scripts de reset do staging para limpar e recarregar dados entre validações. Deploy automático para staging quando CI passa na branch `develop`.

**Critério de aceite:** Staging funcional com dados de exemplo; deploy automático da branch `develop`; reset de dados funcional.

---

### TASK-19-006 — Testes de Carga para o Pipeline da Tríade

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps |
| **Objetivo** | Validar o comportamento do sistema sob carga com múltiplos projetos executando simultaneamente |

**Descrição:**
Criar scripts k6 específicos para simular carga no pipeline da Tríade usando mocks de LLM (para não gerar custos reais): **Cenário 1 — 10 projetos simultâneos:** 10 usuários diferentes iniciam fases ao mesmo tempo, validar que os workers processam sem deadlock e todos concluem corretamente. **Cenário 2 — Pico de feedback:** 50 usuários enviando feedback simultaneamente (stressando o endpoint SSE e o RabbitMQ). **Cenário 3 — Recovery após falha:** simular falha de 1 worker durante processamento, validar que mensagens são reprocessadas pelo outro worker. Documentar a capacidade máxima suportada pela configuração padrão.

**Critério de aceite:** 10 projetos simultâneos sem deadlock; recovery após falha de worker funcional; capacidade máxima documentada.

---

### TASK-19-007 — Dados de Teste Representativos (Test Fixtures)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar um conjunto abrangente de dados de teste para facilitar o desenvolvimento e CI |

**Descrição:**
Criar o pacote `src/backend/testdata/` com fixtures representativos: **Projetos:** um projeto em cada estado possível (DRAFT, IN_PROGRESS cada fase, PAUSED, COMPLETED, ARCHIVED), projetos com e sem Modo Dinâmico, projetos dos 3 tipos de fluxo (A, B, C), **Agentes:** um agente por provider com prompts de exemplo, agente desabilitado, **Prompts do usuário:** conjunto de prompts GLOBAL e por grupo, **Billing Records:** 6 meses de registros de billing para testar os dashboards analíticos. As fixtures podem ser carregadas via `make seed-test-data` em ambiente de desenvolvimento.

**Critério de aceite:** Fixtures cobrindo todos os estados e tipos; comando de seed disponível; usados pelos testes de integração.

---

### TASK-19-008 — Relatório de Qualidade da Plataforma

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Consolidar métricas objetivas da qualidade da plataforma para acompanhamento contínuo |

**Descrição:**
Implementar `GET /api/v1/admin/quality-report` (admin only) que agrega as métricas de qualidade da plataforma: cobertura de testes da plataforma (%), taxa de sucesso das Tríades (% sem auto-rejection), score médio do LLM Judge por fase, tempo médio de execução por fase, custo médio por tipo de projeto, taxa de uptime dos últimos 30 dias, número de projetos completados vs abandonados. No frontend, exibir como cards de métricas na área de admin. Útil para monitorar a saúde da plataforma ao longo do tempo.

**Critério de aceite:** Relatório com todas as métricas; endpoint protegido por role ADMIN; visualização em cards no frontend.

---

### TASK-19-009 — Checklist de Deploy para Produção

| Campo | Valor |
|-------|-------|
| **Camada** | DevOps + Documentação |
| **Objetivo** | Documentar o procedimento completo e seguro de deploy da plataforma para produção |

**Descrição:**
Criar `docs/DEPLOY_CHECKLIST.md` com o procedimento de go-live: **Pré-deploy:** validar que todos os testes passam no CI, validar que o staging foi testado manualmente com os fluxos críticos, criar backup do banco de dados de produção, comunicar janela de manutenção (se aplicável). **Deploy:** procedimento de deploy com docker-compose pull + up, verificar health checks de todos os serviços, verificar logs dos primeiros 5 minutos. **Pós-deploy:** smoke test manual (login, criar projeto, iniciar fase), verificar métricas Prometheus (sem spike de erros), verificar alertas não disparados, comunicar conclusão do deploy. **Rollback:** procedimento de rollback em menos de 5 minutos se problemas críticos.

**Critério de aceite:** Checklist completo e testado; procedures de rollback claros; smoke test definido com exatamente o que verificar.

---

### TASK-19-010 — Revisão de Segurança da Plataforma (Platform Security Review)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + DevOps |
| **Objetivo** | Realizar uma revisão completa de segurança da própria plataforma (não dos projetos gerados) |

**Descrição:**
Executar uma revisão de segurança focada na plataforma Agência de IA: **Autenticação:** JWT com RS256 configurado corretamente, refresh token rotation funcional, session invalidation no logout. **Autorização:** todas as rotas de ADMIN verificadas, isolamento de dados entre usuários validado (usuário A não acessa dados do usuário B). **API Security:** rate limiting em endpoints de autenticação, injection prevention em todos os inputs, sem exposição de stack traces em produção. **Infra Security:** MongoDB com auth habilitada, Redis com auth (requirepass), RabbitMQ com credentials, todos os serviços em rede interna (não expostos diretamente). **Dependências:** `govulncheck` e `npm audit` passando sem vulnerabilidades críticas. Gerar `PLATFORM_SECURITY.md` com resultado da revisão.

**Critério de aceite:** Todos os pontos verificados; PLATFORM_SECURITY.md gerado; sem vulnerabilidades críticas na plataforma.

---

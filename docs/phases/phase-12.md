# Phase 12 — Fluxo A: Fase 7 — Auditoria de Segurança

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-12 |
| **Título** | Fluxo A — Fase 7: Auditoria de Segurança (OWASP) |
| **Tipo** | Backend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-11 concluída (Testes gerados e aprovados) |

---

## Descrição Detalhada

A **Fase 7 do Fluxo A** é a **Auditoria de Segurança**. Os agentes desta fase assumem a persona de um **Security Engineer sênior paranoico**, realizando uma auditoria profunda do código gerado na Fase 5 contra as principais vulnerabilidades do OWASP Top 10 e melhores práticas de segurança modernas.

Esta fase é crítica porque software com vulnerabilidades de segurança pode causar danos severos ao negócio do cliente final. O Agente Revisor de Segurança é o mais rigoroso de todos — ele **nunca aprova** código que tenha vulnerabilidades críticas ou altas (CVSS ≥ 7.0), e o Refinador deve corrigi-las antes de qualquer entrega.

O **Gatilho de Rejeição Automática** também opera aqui: se vulnerabilidades de CVSS ≥ 9.0 forem encontradas, o sistema retorna o código para a Fase 5 automaticamente (sem consumir feedbacks manuais), pois indica um problema estrutural no código gerado que precisa de reescrita, não apenas correção pontual.

O output principal desta fase é o `SECURITY_AUDIT.md` — um relatório completo da auditoria com todas as vulnerabilidades encontradas, sua severidade, prova de conceito (quando aplicável) e as correções aplicadas.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Auditoria OWASP Top 10 completa do código gerado
- ✅ Relatório `SECURITY_AUDIT.md` com severidade, CVE e correções
- ✅ Análise estática de segurança com `gosec` para código Go
- ✅ Verificação de dependências vulneráveis (CVE database)
- ✅ Código corrigido e auditado entregue como artefato final
- ✅ Gatilho de Rejeição Automática para vulnerabilidades críticas

---

## Funcionalidades Entregues

- **Auditoria OWASP:** Verificação sistemática contra OWASP Top 10
- **Análise Estática:** `gosec` + `npm audit` para detecção automatizada
- **Relatório de Segurança:** Documento detalhado com criticidade e correções
- **Hardening Automático:** Correções de segurança aplicadas pelo Refinador

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Toda a auditoria de segurança é executada sequencialmente pela Tríade, sem interrupções. O Agente Security Engineer analisa o código completo e entrega um relatório consolidado ao final.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a auditoria |
| **Velocidade** | Mais rápido — análise completa em uma rodada |
| **Feedback** | Aplicado ao relatório de auditoria como um todo |
| **Ideal para** | Quando o usuário quer uma visão rápida e consolidada da postura de segurança |

### 🎯 Executar uma Task Específica ⭐ **Recomendado**

O usuário seleciona **uma ou mais tasks individualmente**. A Tríade executa apenas a(s) análise(s) escolhida(s) e aguarda aprovação explícita.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por análise — máximo controle sobre as correções de segurança |
| **Velocidade** | Mais controlado — cada issue revisado individualmente |
| **Feedback** | Granular e específico para cada vetor de ataque |
| **Ideal para** | Auditoria de segurança onde cada finding exige decisão humana antes de corrigir |

### 🔀 Modo Híbrido

Execute automaticamente as análises estáticas (`gosec`, `govulncheck`) e pause nas tasks de verificação de IDOR e controle de acesso — que requerem julgamento humano sobre o design de segurança do sistema.

> ⭐ **Fortemente Recomendado:** Use **task a task** nesta phase. Vulnerabilidades de segurança têm implicações diretas no negócio e exigem revisão humana cuidadosa antes de aprovar cada categoria de análise. O Gatilho de Rejeição Automática cobre vulnerabilidades CVSS ≥ 9.0, mas as demais requerem sua decisão.

---

## Tasks

### TASK-12-001 — Prompts do Agente Security Engineer

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts que transformam o Agente Produtor e Revisor em especialistas de segurança rigorosos |

**Descrição:**
Criar system prompts para o Agente Security Engineer com instruções para: **(Produtor)** Análise sistemática do código pela checklist OWASP Top 10, verificação de: SQL/NoSQL injection, broken authentication (tokens, sessões), exposição de dados sensíveis (PII, secrets, credentials), broken access control (rotas sem autenticação, escalada de privilégios), security misconfiguration (headers, CORS, HTTPS), XSS (para frontend), CSRF, insecure deserialization, components com CVEs conhecidos, insufficient logging. **(Revisor)** — persona ultra-crítica: nunca aprova com vulnerabilidade de CVSS ≥ 7.0. Deve gerar para cada issue: severidade (CRITICAL/HIGH/MEDIUM/LOW), CVSS score, descrição técnica, prova de conceito de exploração, arquivo e linha afetados, remediação sugerida. **(Refinador)** aplica todas as correções e justifica cada mudança.

**Critério de aceite:** Prompts cobrem OWASP Top 10 completo; Revisor rejeita com CVSS ≥ 7.0; relatório com POC de exploração.

---

### TASK-12-002 — Análise Estática com gosec

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Executar a ferramenta gosec automaticamente no código Go gerado e incluir os resultados na auditoria |

**Descrição:**
Implementar o `StaticSecurityAnalyzer` que executa `gosec -fmt json ./...` no código gerado em sandbox seguro. Parseia o output JSON do gosec e extrai: lista de issues com severidade (HIGH/MEDIUM/LOW), regra violada (ex: G101 - Hardcoded Credentials, G501 - Import of DES, etc.), arquivo e linha. Injeta os resultados do gosec como contexto adicional no prompt do Agente Produtor de Segurança para que ele analise também os issues detectados estaticamente. Os issues do gosec são incluídos no relatório final com indicação de "detectado por análise estática" vs "detectado por análise manual".

**Critério de aceite:** gosec executado em sandbox; resultados parseados; injetados no prompt do agente; incluídos no relatório.

---

### TASK-12-003 — Análise de Dependências Vulneráveis

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar se as dependências usadas no código gerado possuem CVEs conhecidos |

**Descrição:**
Implementar o `DependencyVulnerabilityChecker` que executa: **Go:** `govulncheck ./...` (ferramenta oficial Google) para detectar CVEs em dependências Go, **Node.js/Frontend:** `npm audit --json` para detectar CVEs em dependências npm. Parseia os outputs e extrai: nome do pacote, versão vulnerável, CVE ID, severidade, versão corrigida disponível. Inclui os resultados no relatório de auditoria. O Refinador gera um `go.mod`/`package.json` atualizado com as versões corrigidas das dependências vulneráveis.

**Critério de aceite:** govulncheck e npm audit executados; CVEs identificados; dependências atualizadas pelo Refinador.

---

### TASK-12-004 — Verificação de Secrets no Código

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Detectar e remover qualquer secret, credencial ou chave de API hardcoded no código gerado |

**Descrição:**
Implementar o `SecretScanner` que executa `trufflehog filesystem` (ou `gitleaks detect`) no repositório virtual do projeto para detectar: API keys hardcoded, senhas em código, tokens de autenticação, certificados privados em código, strings que se parecem com secrets (regex patterns). Se algum secret for encontrado, é um **Gatilho de Rejeição Automática imediata** independente do CVSS (qualquer exposure de credential é CRITICAL). O Refinador recebe a lista exata de arquivos e linhas afetados para correção. Após correção, o scanner é re-executado para confirmar a limpeza.

**Critério de aceite:** TruffleHog/Gitleaks executado; qualquer secret dispara rejeição automática imediata; re-scan após correção.

---

### TASK-12-005 — Verificação de Configurações de Segurança

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar e corrigir configurações de segurança críticas na aplicação gerada |

**Descrição:**
O Agente Security Engineer verifica na Fase 7: **Headers de Segurança HTTP:** Content-Security-Policy, X-Frame-Options, X-Content-Type-Options, Strict-Transport-Security, Referrer-Policy — todos devem estar configurados no middleware GIN. **CORS:** não permitir wildcard `*` em produção; origens explicitamente configuradas. **Rate Limiting:** endpoints de autenticação com rate limit adequado (configurado na Fase 2 de autenticação, mas verificar). **TLS:** verificar que a configuração do servidor força HTTPS em produção. **Banco de Dados:** verificar que queries não são vulneráveis a injection, que índices de TTL de sessão estão configurados, que usuário do banco tem permissões mínimas necessárias. Gera lista de configurações corrigidas/adicionadas.

**Critério de aceite:** Todos os headers de segurança presentes; CORS sem wildcard; rate limiting verificado; configurações documentadas.

---

### TASK-12-006 — Verificação de Controle de Acesso (IDOR e Broken Access Control)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar que nenhum endpoint permite acesso não autorizado a recursos de outros usuários |

**Descrição:**
O Agente Security Engineer analisa todos os endpoints gerados na Fase 5 verificando: IDOR (Insecure Direct Object Reference) — qualquer endpoint que recebe um ID deve verificar se o recurso pertence ao usuário autenticado (ex: `GET /projects/:id` deve validar que o projeto pertence ao `user_id` do JWT), Broken Access Control — verificar que rotas que requerem ADMIN não são acessíveis por usuários comuns, Mass Assignment — verificar que a deserialização de request body não aceita campos sensíveis que não deveriam ser atualizáveis (ex: não permitir que o usuário envie `role: ADMIN` no body de uma atualização de perfil). Lista todos os endpoints com análise de controle de acesso.

**Critério de aceite:** Todos os endpoints analisados; IDOR em cada endpoint verificado; Mass Assignment prevenido; lista documentada.

---

### TASK-12-007 — Geração do Relatório de Auditoria de Segurança

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Consolidar todos os findings em um relatório de auditoria profissional e detalhado |

**Descrição:**
O Refinador da Tríade de Segurança gera o `SECURITY_AUDIT.md` com estrutura profissional: **Executive Summary** (score de segurança 0-100, número de issues por severidade, status geral PASS/FAIL), **Findings Table** (findings ordenados por severidade com ID, título, CVSS, status: FIXED/WONT_FIX/ACCEPTED), **Detalhe por Finding** (descrição técnica, prova de conceito de exploração, impacto, remediação aplicada ou sugerida, arquivo e linha original e corrigida), **Análise de Dependências** (CVEs encontrados e versões atualizadas), **Análise Estática** (gosec findings e status), **Recomendações de Segurança Futura** (monitoramento, pen-test sugerido, atualizações de dependências).

**Critério de aceite:** Relatório profissional com executive summary; score de segurança quantificado; cada finding detalhado com POC e remediação.

---

### TASK-12-008 — Gatilho de Rejeição Automática para Vulnerabilidades Críticas

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a lógica de rejeição automática específica para a Fase 7 de Segurança |

**Descrição:**
Integrar o `AutoRejectionTrigger` à Fase 7 com critérios específicos de segurança: **Rejeição Automática Imediata:** qualquer vulnerabilidade CVSS ≥ 9.0 CRÍTICA, qualquer secret/credential exposto no código, falha de autenticação/autorização (IDOR ou broken access control estrutural). **O Refinador tenta correger primeiro:** vulnerabilidades de CVSS 7.0-8.9 (HIGH) — o Refinador aplica as correções e re-auditoria sem retornar para a Fase 5. Apenas se o Refinador não conseguir corrigir após 2 tentativas é que o Gatilho retorna para a Fase 5. Documentar no relatório de auditoria se o Gatilho foi acionado e quantas vezes, com as vulnerabilidades que o causaram.

**Critério de aceite:** Rejeição imediata para CRITICAL e secrets; HIGH tem até 2 tentativas de autocorreção; documentado no relatório.

---

### TASK-12-009 — Interface de Resultados de Segurança no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir os resultados da auditoria de segurança de forma visual e acessível no painel do projeto |

**Descrição:**
Implementar a aba "Segurança" no painel de projeto com: score de segurança como gauge circular colorido (verde ≥ 80, amarelo 60-79, vermelho < 60), resumo por severidade como cards coloridos (CRITICAL = vermelho, HIGH = laranja, MEDIUM = amarelo, LOW = azul), tabela de findings com filtros por severidade e status, ao clicar em um finding, expandir com detalhes técnicos completos, badge "CORRIGIDO" verde para issues resolvidos pelo Refinador, botão de download do relatório `SECURITY_AUDIT.md` completo.

**Critério de aceite:** Score visual; tabela de findings com filtros; detalhes expansíveis; badge de corrigido; download do relatório.

---

### TASK-12-010 — Testes do Pipeline de Segurança

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Validar que o pipeline de auditoria de segurança detecta corretamente as vulnerabilidades intencionais |

**Descrição:**
Implementar um conjunto de fixtures de código especialmente construídas com vulnerabilidades conhecidas (para fins de teste): código com SQL injection, código com hardcoded secret, endpoint sem verificação de autenticação, CORS com wildcard. Executar o pipeline de segurança completo contra essas fixtures e verificar que: todas as vulnerabilidades são detectadas, as severidades estão corretas, o Gatilho de Rejeição é acionado para CRITICAL/CVSS ≥ 9.0, o Refinador corrige de forma adequada. Garantir que o scanner de secrets detecta as patterns mais comuns (AWS keys, GitHub tokens, JWT secrets, etc.).

**Critério de aceite:** Todas as vulnerabilidades das fixtures detectadas; severidades corretas; Gatilho acionado corretamente; patterns comuns de secrets detectadas.

---

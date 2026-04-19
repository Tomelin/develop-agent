# Phase 17 — Modo Dinâmico Multi-Modelo e Controle de Configuração

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-17 |
| **Título** | Modo Dinâmico Multi-Modelo e Controle de Configuração |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Média |
| **Pré-requisitos** | PHASE-03, PHASE-06 concluídas |

---

## Descrição Detalhada

O **Modo Dinâmico** é uma das inovações mais distintas da plataforma, conforme descrito no PROJECT.md: em vez de usar sempre o mesmo modelo de IA em todas as fases de um projeto, o sistema sorteia aleatoriamente modelos diferentes para cada posição da Tríade (Produtor, Revisor, Refinador) a cada execução.

O objetivo é **eliminar os vieses cognitivos** que surgem quando um único modelo verifica seu próprio trabalho. Exemplos reais: o Produtor Claude gera código em um estilo, o Revisor GPT-4o identifica inconsistências que o Claude não veria no próprio output, e o Refinador Gemini aplica uma abordagem diferente para resolver as críticas.

Esta fase implementa o mecanismo de sorteio, os controles de configuração do Modo Dinâmico, o painel de visualização do sorteio realizado em cada execução, e as métricas de diversidade de modelos utilizados no projeto.

Também cobre as configurações globais da plataforma — onde o administrador pode controlar configurações de sistema, como tempo de timeout dos agentes, paralelismo de workers, e outras configurações operacionais.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Modo Dinâmico configurável por projeto (on/off)
- ✅ Algoritmo de sorteio com garantia de diversidade de providers
- ✅ Painel de visualização do sorteio atual por fase
- ✅ Histórico de quais modelos foram usados em cada execução
- ✅ Configurações globais da plataforma (admin only)
- ✅ Métricas de diversidade de modelos por projeto

---

## Funcionalidades Entregues

- **Sorteio Inteligente:** Diversidade de providers garantida na mesma Tríade
- **Configuração por Projeto:** Modo Dinâmico ativável independentemente por projeto
- **Auditoria de Sorteio:** Registro completo de qual modelo foi sorteado e por quê
- **Configurações Globais:** Painel administrativo de configurações da plataforma

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa ⭐ **Recomendado**

Todo o sistema de Modo Dinâmico e configurações é implementado em sequência pela Tríade, desde o algoritmo de sorteio até o painel de configuração do usuário.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a implementação |
| **Velocidade** | Mais rápido — sorteador, configurações e UI geradas em conjunto |
| **Feedback** | Aplicado ao sistema completo de Modo Dinâmico |
| **Ideal para** | Implementar o Modo Dinâmico de forma integrada e consistente |

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Útil para implementar apenas o algoritmo de sorteio sem a UI, ou apenas as feature flags sem o painel.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task |
| **Velocidade** | Mais controlado |
| **Feedback** | Específico para o componente selecionado |
| **Ideal para** | Refinar componentes específicos do Modo Dinâmico sem reprocessar os demais |

### 🔀 Modo Híbrido

Execute automaticamente o algoritmo de sorteio e a persistência de histórico, e revise individualmente as feature flags e as configurações de diversidade de provider que têm maior impacto no comportamento do sistema.

> 💡 **Dica:** Para esta phase (Modo Dinâmico), o modo recomendado é **phase completa** — o sorteador, as garantias de diversidade e o painel de configuração funcionam melhor quando desenvolvidos como um sistema coeso.

---

## Tasks

### TASK-17-001 — Algoritmo de Sorteio com Diversidade de Providers

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o algoritmo de seleção aleatória que garante diversidade máxima de providers na Tríade |

**Descrição:**
Implementar o `DynamicAgentSelector` em `src/backend/domain/agent/dynamic_selector.go` com o algoritmo de sorteio: **(1)** Busca todos os agentes habilitados para a skill da fase, **(2)** Agrupa por provider (OpenAI, Anthropic, Google, Ollama), **(3)** Tenta selecionar 3 providers diferentes (Produtor de provider A, Revisor de provider B, Refinador de provider C), **(4)** Dentro de cada provider, sorteia o agente específico (caso haja múltiplos do mesmo provider para a skill), **(5)** Se não houver 3 providers distintos disponíveis: aceita repetição de provider mas nunca do mesmo agente (mesmo Nome), **(6)** Se houver apenas 1 agente disponível: usa o mesmo 3 vezes com log de aviso. Registrar o resultado do sorteio como `TriadSelection{PhaseName, Producer, Reviewer, Refiner, SelectionTimestamp}` para auditoria.

**Critério de aceite:** Diversidade de providers garantida quando possível; fallback progressivo bem definido; sorteio registrado para auditoria.

---

### TASK-17-002 — Toggle de Modo Dinâmico por Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Implementar o controle de ativação do Modo Dinâmico por projeto |

**Descrição:**
Quando o Modo Dinâmico está **desativado**: o usuário define manualmente qual agente usar para cada papel da Tríade por fase (configuração fixa). Quando está **ativado**: o sistema sorteia automaticamente a cada execução. Implementar: `PUT /api/v1/projects/:id/dynamic-mode` (ativa/desativa o modo e configura os agentes fixos quando desativado), endpoint de `preview` que simula um sorteio sem executar (para o usuário ver quais modelos seriam sorteados antes de iniciar), validação de que há agentes suficientes para o Modo Dinâmico (mínimo 2 agentes habilitados para a skill). No frontend: toggle de Modo Dinâmico no painel do projeto com explicação, quando desativado: seletor de agente fixo por posição da Tríade.

**Critério de aceite:** Toggle funcional; preview de sorteio sem executar; validação de agentes mínimos; configuração de agentes fixos quando desativado.

---

### TASK-17-003 — Painel de Visualização do Sorteio da Tríade

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir de forma visual e clara quais agentes foram sorteados para a Tríade de cada fase |

**Descrição:**
Implementar o componente `TriadCompositionPanel` exibido no painel de projeto com: para cada fase executada (ou em execução), mostrar os 3 agentes sorteados para a Tríade (Produtor, Revisor, Refinador) com: avatar colorido do agente (cor do provider), nome do agente, modelo de IA utilizado, provider (badge colorido: GPT blue, Claude orange, Gemini green, Ollama purple), quando Modo Dinâmico ativo: indicador "Sorteado dinamicamente" com o ícone de dado, quando fixo: indicador "Configuração fixa". Exibir também o histórico completo de todos os sorteios do projeto em uma tabela expandível.

**Critério de aceite:** Tríade visualizada por fase; cores por provider; indicador de dinâmico vs fixo; histórico completo.

---

### TASK-17-004 — Análise de Diversidade do Projeto

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Métricas de quantos providers/modelos diferentes foram utilizados no projeto |

**Descrição:**
Implementar `GET /api/v1/projects/:id/diversity-metrics` que retorna: lista de todos os providers utilizados no projeto com % de uso, lista de todos os modelos utilizados, número de Tríades onde todos os 3 providers eram diferentes (diversidade máxima), número de Tríades com providers repetidos, distribuição de uso por papel (Produtor/Revisor/Refinador) — qual modelo foi mais usado como Revisor, etc. No frontend, exibir na aba "Configurações" do projeto um card de "Diversidade de IA" com gráfico de pizza dos providers utilizados e score de diversidade (0-100%).

**Critério de aceite:** Métricas calculadas corretamente; score de diversidade quantificado; gráfico de pizza por provider.

---

### TASK-17-005 — Configurações Globais da Plataforma (Admin Panel)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Criar o painel de configurações globais da plataforma acessível apenas para admins |

**Descrição:**
Implementar a página `/admin/settings` (apenas para role ADMIN): **Configurações de Workers:** número máximo de workers concorrentes, timeout de execução por agente (default: 5 min), timeout total da Tríade (default: 20 min). **Configurações de Modelo:** model padrão quando nenhum agente for selecionado, modelo secundário para geração de SPEC.md (leve e barato), **Limites de Sistema:** máximo de projetos por usuário, máximo de fases em execução simultânea por usuário, limite de tokens por SPEC.md. **Configurações de Retry:** número máximo de tentativas, tempo de backoff. Todas configuradas via API `PUT /api/v1/admin/settings` e persistidas no MongoDB.

**Critério de aceite:** Painel de admin com todas as configurações; apenas role ADMIN tem acesso; configurações persistidas e aplicadas imediatamente.

---

### TASK-17-006 — Configuração de Agentes por Fase (Modo Fixo por Projeto)

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a interface de configuração de agentes fixos por fase quando o Modo Dinâmico está desativado |

**Descrição:**
Implementar a tela de configuração de agentes do projeto (aba "Configuração de Agentes" no painel do projeto): matriz de seleção de agentes onde as linhas são as fases (1-9 + Fluxo B + Fluxo C) e as colunas são os papéis (Produtor, Revisor, Refinador). Em cada célula: dropdown de seleção de agente habilitado e compatível com a skill da fase, badge mostrando o provider do agente selecionado, opção de deixar "Dinâmico" para fases específicas enquanto outras são fixas. Botão de "Aplicar a Todas as Fases" (preenche o mesmo agente em todas as fases do mesmo papel). Preview do custo estimado com a configuração atual.

**Critério de aceite:** Matriz de configuração completa; seleção por fase e papel; opção de misto (algumas fases fixas, outras dinâmicas); preview de custo.

---

### TASK-17-007 — Logs de Auditoria do Sorteio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Manter logs completos de auditoria de todos os sorteios realizados para reprodutibilidade e análise |

**Descrição:**
Implementar a entidade `TriadSelectionLog` com: `ID`, `ProjectID`, `PhaseNumber`, `ExecutionID`, `Mode` (DYNAMIC/FIXED), `CandidateAgents` (lista de todos os agentes disponíveis para a skill no momento do sorteio), `SelectedTriad` {producer, reviewer, refiner}, `SelectionReason` (texto explicando por que cada agente foi selecionado — ex: "sorteado aleatoriamente de 3 candidatos do provider Anthropic"), `Timestamp`. Implementar `GET /api/v1/projects/:id/selection-logs` que retorna todos os logs de sorteio do projeto com filtros por fase. Useful para debugging e para o usuário entender as decisões do Modo Dinâmico.

**Critério de aceite:** Logs com candidatos e motivo de seleção; endpoint com filtros; útil para debugging de sorteios.

---

### TASK-17-008 — Testes do Algoritmo de Sorteio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que o algoritmo de sorteio cobre todos os cenários corretamente |

**Descrição:**
Testes unitários do `DynamicAgentSelector` cobrindo: sorteio com 4 providers disponíveis (todos diferentes na Tríade — cenário ideal), sorteio com 2 providers disponíveis (2 diferentes e 1 repetido — fallback nível 1), sorteio com 1 provider disponível (todos iguais mas agentes diferentes — fallback nível 2), sorteio com 1 agente total para a skill (mesmo agente nos 3 papéis com aviso), distribuição estatística (executar 1000 sorteios e verificar que nenhum agente é sistematicamente favorecido — distribuição uniforme esperada), reprodutibilidade (com seed fixo, o sorteio é determinístico — para testes), log de auditoria gerado corretamente.

**Critério de aceite:** Todos os cenários de fallback testados; distribuição uniforme estatisticamente verificada; log de auditoria gerado.

---

### TASK-17-009 — Feature Flag System para Funcionalidades em Desenvolvimento

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Implementar um sistema de feature flags para controle gradual de lançamento de funcionalidades |

**Descrição:**
Implementar o `FeatureFlagService` que permite ativar/desativar funcionalidades da plataforma sem redeploy. Feature flags gerenciadas pelo admin via `GET/PUT /api/v1/admin/feature-flags`. Flags iniciais: `DYNAMIC_MODE_ENABLED` (Modo Dinâmico disponível para usuários), `FLOW_B_ENABLED` (Fluxo de Landing Page disponível), `FLOW_C_ENABLED` (Fluxo de Marketing disponível), `DEVOPS_PHASE_ENABLED` (Fase 9 disponível), `BILLING_PANEL_ENABLED`, `AUTO_REJECTION_ENABLED`. No frontend, componente `FeatureGate` que verifica a flag antes de renderizar uma funcionalidade, ocultando-a graciosamente se desabilitada.

**Critério de aceite:** Flags gerenciadas via admin; FeatureGate oculta funcionalidades desabilitadas; sem redeploy necessário.

---

### TASK-17-010 — Documentação do Modo Dinâmico para o Usuário

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar documentação in-app explicando o Modo Dinâmico de forma acessível ao usuário |

**Descrição:**
Implementar a seção "Como Funciona o Modo Dinâmico" na documentação in-app com: explicação visual do problema (viés de modelo único — "o juiz não pode ser o mesmo que o réu"), como o sorteio funciona (diagrama mostrando múltiplos modelos sendo selecionados), benefícios esperados (qualidade vs custo), riscos e compensações (pode custar mais com modelos premium sendo sorteados), dicas de configuração (quando usar Modo Dinâmico vs fixo). Acessível via botão "?" ao lado do toggle de Modo Dinâmico no painel do projeto.

**Critério de aceite:** Documentação clara e visual; diagrama do sorteio; dicas de configuração práticas; acessível via contexto.

---

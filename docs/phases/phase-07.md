# Phase 07 — Fluxo A: Fase 1 — Criação do Projeto (Agente Entrevistador)

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-07 |
| **Título** | Fluxo A — Fase 1: Criação do Projeto com Agente Entrevistador |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-01 a PHASE-06 concluídas |

---

## Descrição Detalhada

Esta fase implementa a **Fase 1 do Fluxo A** do pipeline de desenvolvimento de software: a **Criação do Projeto com o Agente Entrevistador**. É a porta de entrada do pipeline completo e talvez a mais estratégica — pois a qualidade do entendimento do produto aqui determina a qualidade de todas as fases subsequentes.

O Agente Entrevistador opera de forma diferente das demais fases. Ele **não utiliza a Tríade** — age como um único agente conversacional que conduz uma entrevista estruturada com o usuário para transformar uma ideia vaga em uma visão clara e documentada de produto. A conversa é iterativa: o agente faz perguntas, escuta, reformula o que entendeu e pede confirmação antes de avançar.

O usuário pode enviar até **10 feedbacks/iterações** nesta fase (o dobro das demais) porque ela é a mais crítica para o alinhamento do produto. O agente não avança sem o aval explícito do usuário.

Ao final desta fase, o sistema gera automaticamente o `SPEC.md` inicial com a visão consolidada do produto, que servirá como base para todas as fases seguintes.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Interface de chat conversacional com o Agente Entrevistador
- ✅ Lógica de entrevista estruturada com perguntas guiadas por template
- ✅ Validação de que o agente não avança sem confirmação explícita do usuário
- ✅ Geração automática do documento de Visão do Produto ao finalizar
- ✅ SPEC.md inicial gerado e armazenado no projeto
- ✅ Limite de 10 feedbacks por sessão de entrevista

---

## Funcionalidades Entregues

- **Chat Conversacional:** Interface de chat em tempo real com o Agente Entrevistador
- **Entrevista Estruturada:** Template de perguntas cobrindo produto, usuários, diferencial, tecnologia
- **Documentação Automática:** Geração do documento de Visão do Produto ao finalizar
- **Controle de Iterações:** Limite de 10 feedbacks com contador visível

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Todas as tasks são executadas sequencialmente pela Tríade de Agentes (Produtor → Revisor → Refinador), sem interrupções entre elas. O sistema avança automaticamente de uma task para a próxima até concluir a phase inteira e aguarda uma única aprovação ao final.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a phase |
| **Velocidade** | Mais rápido — execução contínua sem pausas |
| **Feedback** | Aplicado à phase como um todo |
| **Ideal para** | Phases bem compreendidas onde o usuário confia na execução automática |

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. A Tríade desenvolve apenas a(s) task(s) escolhida(s) e aguarda aprovação explícita antes de prosseguir para a próxima.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — o usuário controla o ritmo |
| **Velocidade** | Mais controlado — requer interação entre tasks |
| **Feedback** | Granular e específico para cada task |
| **Ideal para** | Tasks críticas ou complexas que exigem revisão antes de avançar |

### 🔀 Modo Híbrido

É possível **combinar os dois modos**: inicie a phase automaticamente e pause manualmente em qualquer task que exija atenção especial. Após aprovar aquela task individualmente, a execução automática retoma a partir da próxima task.

> 💡 **Dica:** Esta phase é a **Fase 1 do Fluxo A** (Entrevistador), que por natureza é interativa e conversacional. O modo natural de operação é sempre task a task — o Agente Entrevistador conduz iterações com o usuário até chegar à Confirmação da Visão do Produto.

---

## Tasks

### TASK-07-001 — Implementação do Agente Entrevistador Especializado

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar a lógica do Agente Entrevistador com seu template de perguntas estruturadas |

**Descrição:**
Implementar o `InterviewerAgent` como um tipo especializado de execução de fase, diferente da Tríade padrão. O agente opera em modo conversacional: mantém o histórico completo da conversa e usa-o em cada mensagem subsequente. Implementar o template de perguntas estruturadas que o agente segue, cobrindo: (1) Qual é o problema que o produto resolve? (2) Quem são os usuários principais e seus perfis? (3) Qual é o diferencial competitivo? (4) Quais são as 3-5 funcionalidades mais críticas do MVP? (5) Há alguma preferência de tecnologia ou restrição técnica? (6) Qual é o prazo e nível de urgência? (7) Há alguma integração com sistemas existentes? (8) Como será o modelo de negócio/monetização? O agente conduz a entrevista organicamente, não como um formulário rígido.

**Critério de aceite:** Agente conduz entrevista organicamente; cobre todas as áreas; reformula o entendimento para validação.

---

### TASK-07-002 — Persistência da Sessão de Entrevista (Histórico de Conversa)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Persistir o histórico completo da conversa de entrevista para permitir continuação e auditoria |

**Descrição:**
Implementar a entidade `InterviewSession` com: `ID`, `ProjectID`, `Messages` (slice de `{Role: USER|ASSISTANT, Content, Timestamp}`), `Status` (ACTIVE/AWAITING_CONFIRMATION/COMPLETED/ABANDONED), `IterationCount`, `MaxIterations` (10), `CompletedAt`. Armazenar no MongoDB com a sessão inteira (histórico completo de mensagens). Quando a sessão é retomada (usuário fecha e reabre o navegador), carregar o histórico completo para o frontend. Implementar `GET /api/v1/projects/:id/interview` e `POST /api/v1/projects/:id/interview/message`.

**Critério de aceite:** Histórico persistido; sessão retomável; contador de iterações correto.

---

### TASK-07-003 — Endpoint de Mensagem do Chat de Entrevista

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Expor endpoint de envio de mensagem que mantém a conversa com o Agente Entrevistador |

**Descrição:**
Implementar `POST /api/v1/projects/:id/interview/message` (SSE para streaming de resposta em tempo real). O endpoint recebe `{content: "minha ideia é..."}`, adiciona a mensagem do usuário ao histórico, envia o histórico completo ao InterviewerAgent via AgentSDK com streaming habilitado, transmite a resposta token-a-token via SSE para uma experiência de chat fluida como o ChatGPT, salva a mensagem do assistente ao histórico, incrementa `IterationCount`, quando `IterationCount >= 10`, adiciona nota na resposta indicando o limite de iterações. Registrar tokens e custo no BillingTracker.

**Critério de aceite:** Streaming funcional; histórico atualizado; limite de 10 iterações; billing registrado.

---

### TASK-07-004 — Fluxo de Confirmação e Encerramento da Entrevista

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o mecanismo de confirmação explícita do usuário para encerrar a entrevista e gerar o documento de visão |

**Descrição:**
Implementar `POST /api/v1/projects/:id/interview/confirm` que o usuário chama quando está satisfeito com o entendimento gerado. O endpoint: muda o status da sessão para `COMPLETED`, aciona o `VisionDocumentGenerator` (ver próxima task), gera o `SPEC.md` inicial do projeto, muda a Fase 1 do projeto para `COMPLETED`, emite evento `PHASE_1_COMPLETED` via SSE. Implementar também `POST /api/v1/projects/:id/interview/regenerate-vision` que regenera o documento de visão sem encerrar a sessão (para ajustes finos antes da confirmação final).

**Critério de aceite:** Confirmação encerra sessão; documento de visão gerado; fase avançada; evento emitido.

---

### TASK-07-005 — Gerador de Documento de Visão do Produto

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar automaticamente um documento estruturado de Visão do Produto a partir da entrevista |

**Descrição:**
Implementar `VisionDocumentGenerator` que recebe o histórico completo da entrevista e usa um modelo LLM para gerar um documento estruturado em Markdown (`VISION.md`) contendo: Sumário Executivo (2-3 parágrafos), Problema a ser Resolvido, Público-Alvo (personas detalhadas), Proposta de Valor Única, Escopo do MVP (lista de funcionalidades in/out of scope), Requisitos Não-Funcionais identificados, Tecnologias sugeridas/restrições, Integrações necessárias, Modelo de Negócio, Próximos Passos claros. O documento gerado é armazenado no projeto como artefato da Fase 1 e como base do SPEC.md inicial.

**Critério de aceite:** Documento gerado com todas as seções; baseado no conteúdo real da entrevista; armazenado como artefato.

---

### TASK-07-006 — Interface de Chat de Entrevista no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a interface de chat com design premium para a entrevista com o Agente Entrevistador |

**Descrição:**
Implementar a tela `/projects/:id/interview` com: interface de chat inspirada em produtos como ChatGPT/Claude, mensagens do usuário à direita (fundo colorido com cor primária), mensagens do agente à esquerda (fundo neutro escuro), renderização de Markdown nas respostas do agente (negrito, listas, títulos), animação de "digitando..." enquanto o agente processa, streaming de resposta token-a-token (o texto aparece progressivamente como o ChatGPT), campo de input fixo na parte inferior com envio por Enter ou botão, contador de iterações restantes (`X de 10 iterações utilizadas`) no topo, botão "Confirmar e Avançar" somente disponível após pelo menos 3 iterações.

**Critério de aceite:** Chat com streaming funcional; Markdown renderizado; contador de iterações; botão de confirmação com regra de mínimo de 3 iterações.

---

### TASK-07-007 — Tela de Preview do Documento de Visão

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir o documento de Visão do Produto gerado para revisão antes de confirmar o encerramento da entrevista |

**Descrição:**
Implementar um drawer/modal slide-over que exibe o `VISION.md` gerado com: renderização Markdown completa com estilos bem formatados (títulos, listas, tabelas), indicador de "última geração: X minutos atrás", botão "Regenerar Documento" (chama a API de regeneração sem encerrar a entrevista), botão "Confirmar Visão e Avançar para Engenharia" (chama o endpoint de confirmação), botão de download do documento em Markdown e PDF. Layout fullscreen opcional para leitura confortável do documento longo.

**Critério de aceite:** Markdown renderizado com formatação completa; regeneração funcional; botão de confirmação com feedback visual; download em MD e PDF.

---

### TASK-07-008 — Indicador de Progresso da Entrevista

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Mostrar ao usuário o progresso da entrevista com as áreas já cobertas e as pendentes |

**Descrição:**
Implementar um sidebar ou panel no chat que exibe as áreas temáticas da entrevista com indicador de cobertura: ✅ Problema identificado / ⏳ Pendente, ✅ Público-alvo definido, ✅ Funcionalidades do MVP, ⏳ Tecnologias e restrições, etc. A cobertura é detectada pelo backend via análise do histórico de mensagens (usando regex ou LLM leve). Isso ajuda o usuário a saber o que já foi discutido e o que ainda falta para uma entrevista completa antes de finalizar.

**Critério de aceite:** Área de cobertura atualizada dinamicamente; análise funcional; design não invasivo ao chat.

---

### TASK-07-009 — Histórico de Entrevistas Passadas

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Permitir que o usuário acesse e revise entrevistas passadas de projetos anteriores |

**Descrição:**
Implementar a seção "Histórico da Entrevista" no painel do projeto com: linha do tempo das mensagens trocadas durante a entrevista (apenas leitura), possibilidade de copiar mensagens individuais, link para o documento de Visão gerado, data e duração da entrevista. Útil para auditoria e para o usuário revisar as decisões tomadas na descoberta do produto.

**Critério de aceite:** Histórico completo exibido; apenas leitura; link para documento de Visão.

---

### TASK-07-010 — Testes do Agente Entrevistador

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Validar o comportamento do Agente Entrevistador em diferentes cenários de entrada |

**Descrição:**
Implementar testes usando mock do AgentSDK para: entrevista completa simulada (histórico de 5 turnos, confirmação, geração de VISION.md), tentativa de enviar mensagem após atingir limite de 10 iterações (deve retornar erro), tentativa de confirmar sem mínimo de interações (deve retornar aviso), geração de VISION.md a partir de histórico com informações incompletas (deve gerar documento com seções marcadas como "A definir"), retomada de sessão (carregar histórico existente e continuar).

**Critério de aceite:** Todos os cenários testados; limite de iterações validado; VISION.md gerado mesmo com infos incompletas.

---

# Phase 08 — Fluxo A: Fases 2 e 3 — Engenharia e Arquitetura de Software

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-08 |
| **Título** | Fluxo A — Fases 2 e 3: Engenharia e Arquitetura de Software |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Alta |
| **Pré-requisitos** | PHASE-07 concluída (Fase 1 do projeto implementada) |

---

## Descrição Detalhada

Esta phase implementa as **Fases 2 e 3 do Fluxo A**: Engenharia de Software e Arquitetura de Software. Ambas utilizam a Tríade completa (Produtor → Revisor → Refinador) e são as fases onde o produto começa a ganhar forma técnica estruturada.

**Fase 2 — Engenharia de Software:** Aqui os agentes de engenharia de requisitos se debruçam sobre a Visão do Produto (gerada na Fase 1) para extrair e documentar todas as regras de negócio, requisitos funcionais (RF) e não-funcionais (RNF) da plataforma. O output é um documento de engenharia que será a base de toda a arquitetura e desenvolvimento.

**Fase 3 — Arquitetura de Software:** Com os requisitos bem definidos, os agentes arquitetos projetam a solução técnica: modelagem de dados (entidades e relacionamentos), escolha da stack tecnológica (validada contra os prompts do usuário), design patterns a serem aplicados, desenho de arquitetura de alto nível (microserviços, monolítico, etc.), e definição de infraestrutura.

A partir da Fase 2, o projeto se divide em dois trilhos paralelos: **Frontend** e **Backend**, e cada trilho tem sua própria Tríade. Esta fase implementa essa bifurcação no pipeline.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Pipeline de execução paralela para trilhos Frontend e Backend
- ✅ Prompts especializados para agentes de Engenharia e Arquitetura
- ✅ Interface de visualização dos artefatos das Fases 2 e 3
- ✅ Suporte a feedback separado por trilho (Frontend/Backend)
- ✅ Diagrama de arquitetura gerado como artefato visual

---

## Funcionalidades Entregues

- **Execução Paralela:** Frontend e Backend executam em paralelo a partir da Fase 2
- **Documentos de Engenharia:** RFs, RNFs e regras de negócio estruturadas
- **Documento de Arquitetura:** Stack, modelagem de dados, design patterns, infra
- **Visualizador de Artefatos:** Renderização de documentos Markdown com diagramas

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

> 💡 **Dica:** Para esta phase (Engenharia e Arquitetura — Fases 2 e 3 do Fluxo A), o modo recomendado é **híbrido** — os prompts e contextos RAG podem rodar automaticamente, mas a tarefa de validação de alinhamento entre fases (TASK-08-010) deve ser revisada individualmente antes de aprovar a Fase 3.

---

## Tasks

### TASK-08-001 — Suporte a Trilhos Paralelos no Pipeline (Frontend/Backend Split)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a bifurcação do pipeline em trilhos Frontend e Backend a partir da Fase 2 |

**Descrição:**
Adaptar o `TriadOrchestrator` e o `ProjectStateMachine` para suportar fases com dois trilhos paralelos. Quando `PhaseTrack == BOTH`, criar dois registros de `TriadExecution` simultâneos: um para FRONTEND e um para BACKEND. Publicar duas mensagens separadas na fila RabbitMQ. A fase só avança para a próxima quando **ambos** os trilhos estiverem `COMPLETED`. Implementar `GET /api/v1/projects/:id/phases/:phaseNumber/tracks` que retorna o status de cada trilho individualmente. A aprovação por trilho é separada: o usuário pode aprovar o trilho Frontend independentemente do Backend.

**Critério de aceite:** Dois TriadExecution criados para fases com split; fase avança somente quando ambos completos; aprovação independente por trilho.

---

### TASK-08-002 — Prompts Especializados para Agente de Engenharia de Requisitos

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir os prompts de sistema e de execução específicos para a Fase 2 de Engenharia |

**Descrição:**
Criar os prompts especializados para o agente de Engenharia de Requisitos (Produtor da Fase 2) que instruem a extrair e documentar: lista numerada de Requisitos Funcionais (RF001 a RFxxx), lista de Requisitos Não-Funcionais por categoria (Performance, Segurança, Escalabilidade, Usabilidade, Manutenibilidade), regras de negócio críticas com linguagem precisa, glossário de termos do domínio. Prompt do Revisor: verificar completude dos requisitos, ambiguidades, conflitos entre RFs, ausência de RNFs críticos. Prompt do Refinador: gerar documento `ENGINEERING.md` limpo e numerado. Armazenar os prompts no seed de agentes.

**Critério de aceite:** Prompts geram documentos de engenharia de qualidade profissional; Revisor detecta ambiguidades e gaps; output em formato padronizado.

---

### TASK-08-003 — Prompts Especializados para Agente de Arquitetura

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir os prompts específicos para os agentes de Arquitetura Frontend e Backend |

**Descrição:**
Criar prompts para os agentes de Arquitetura incluindo: **Backend Architect:** modelagem de entidades (com atributos e relacionamentos), escolha de banco de dados e justificativa, padrões arquiteturais (Clean Architecture, CQRS, Event Sourcing), definição de APIs (REST vs GraphQL vs gRPC), estratégia de cache e mensageria, considerações de segurança e autenticação. **Frontend Architect:** estrutura de componentes, gestão de estado (Redux, Zustand), padrões de design (Design System, Atomic Design), estratégia de roteamento, SSR vs SPA, responsividade e acessibilidade. O Revisor de arquitetura deve verificar consistência com os RFs e RNFs da Fase 2.

**Critério de aceite:** Dois conjuntos de prompts distintos para Front e Back; Revisor verifica alinhamento com requisitos da Fase 2.

---

### TASK-08-004 — Instrução de Fase com Contexto RAG

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que cada fase recebe o contexto completo das fases anteriores via RAG |

**Descrição:**
Implementar o `PhaseContextBuilder` que, ao iniciar qualquer fase (2 em diante), compõe o contexto RAG do projeto concatenando: o `SPEC.md` da Fase 1 (Visão do Produto), o `SPEC.md` da Fase 2 (para Fase 3 em diante), os prompts GLOBAL do usuário, os prompts específicos do grupo da fase. Injetar esse contexto como `context` no prompt do Produtor. Garantir que o contexto nunca exceda o limite de tokens do modelo utilizado (fallback: truncar o SPEC.md mais antigo se necessário, preservando o mais recente).

**Critério de aceite:** Contexto RAG composto corretamente; limite de tokens respeitado com truncamento inteligente; contexto verificável nos logs de execução.

---

### TASK-08-005 — Handler de Inicialização das Fases 2 e 3

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Adaptar o handler de início de fase para suportar as particularidades das Fases 2 e 3 |

**Descrição:**
Adaptar `POST /api/v1/projects/:id/phases/:phaseNumber/start` para tratar a bifurcação: para Fases 2+ que têm split Frontend/Backend, criar os dois TriadExecution simultaneamente, publicar duas mensagens na fila (uma por trilho), retornar IDs de ambas as execuções no response. Implementar validação de que a Fase 1 está `COMPLETED` antes de iniciar a Fase 2. Para a Fase 3, validar que a Fase 2 está `COMPLETED` em ambos os trilhos antes de prosseguir.

**Critério de aceite:** Validação de pré-requisitos; dois TriadExecution para fases com split; response com IDs de ambas as execuções.

---

### TASK-08-006 — Visualizador de Artefatos das Fases no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar componente rico para visualização dos documentos gerados pela Tríade |

**Descrição:**
Implementar o componente `ArtifactViewer` reutilizável que renderiza os artefatos gerados pelas fases. Funcionalidades: renderização Markdown completa com suporte a diagramas Mermaid (para diagramas de arquitetura), syntax highlighting para blocos de código, índice de navegação auto-gerado a partir dos headings do documento, botão de copiar por bloco de código, modo fullscreen para leitura imersiva, botão de download em Markdown e PDF, diff visual entre versões quando há feedbacks (mostra o que mudou após cada iteração da Tríade).

**Critério de aceite:** Markdown com Mermaid renderizado; syntax highlighting; índice de navegação; diff entre versões.

---

### TASK-08-007 — Interface de Acompanhamento da Tríade em Tempo Real

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de acompanhamento da execução da Tríade com visualização em tempo real de cada step |

**Descrição:**
Implementar o componente `TriadProgressPanel` que exibe durante a execução: os três agentes da Tríade (Produtor, Revisor, Refinador) como cards com avatar gerado dinamicamente (iniciais do nome + cor do provider), indicador animado de "processando" para o agente ativo, output parcial em streaming para o agente ativo (os primeiros 500 chars como preview), após conclusão de cada step, exibir indicador de sucesso e estatísticas (tokens, duração, modelo). Para fases com split Frontend/Backend, exibir as duas Tríades lado a lado (quando tela permitir) ou como tabs.

**Critério de aceite:** Animação de processamento; streaming de preview; estatísticas por step; layout de split Frontend/Backend.

---

### TASK-08-008 — Interface de Feedback por Trilho

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Permitir que o usuário envie feedbacks separados para os trilhos Frontend e Backend |

**Descrição:**
Implementar o formulário de feedback com: tabs separadas para "Feedback Frontend" e "Feedback Backend", textarea de feedback rico com suporte a listas e formatação básica, contador de feedbacks restantes bem visível (`3 de 5 feedbacks utilizados`), histórico dos feedbacks anteriores enviados (expandível), preview em tempo real do que o Refinador receberá (feedback + artefato atual), botão "Enviar Feedback e Refinar" que aciona novo ciclo da Tríade apenas do trilho selecionado, botão "Aprovar e Avançar" separado por trilho.

**Critério de aceite:** Feedback separado por trilho; contador visível; histórico de feedbacks; preview do que será enviado ao agente.

---

### TASK-08-009 — Extrator de Diagrama de Arquitetura

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Extrair o diagrama de arquitetura Mermaid do documento gerado pela Fase 3 e disponibilizá-lo separadamente |

**Descrição:**
Implementar o `DiagramExtractor` que parseia o `ARCHITECTURE.md` gerado pela Fase 3 e extrai blocos de código Mermaid. Armazena os diagramas como artefatos separados do projeto com tipo `DIAGRAM_ARCHITECTURE`. Implementar `GET /api/v1/projects/:id/diagram` que retorna a lista de diagramas extraídos. O frontend usa essa API para renderizar os diagramas em um painel dedicado de "Visualização de Arquitetura" com zoom e pan. Se a Fase 3 não gerou diagramas Mermaid, o backend instrui o agente a regenerar com diagramas incluídos.

**Critério de aceite:** Diagramas Mermaid extraídos; endpoint dedicado; renderização com zoom e pan no frontend.

---

### TASK-08-010 — Validação de Alinhamento entre Fases 2 e 3

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar automaticamente se a arquitetura gerada na Fase 3 está alinhada com os requisitos da Fase 2 |

**Descrição:**
Implementar o `AlignmentValidator` que executa após a conclusão da Fase 3 (como pós-processamento). Usa um LLM leve para verificar: todos os RFs da Fase 2 têm componente arquitetural correspondente na Fase 3, todos os RNFs de performance e segurança têm estratégias definidas, não há tecnologias na arquitetura que conflitem com as restrições definidas na Fase 2. Gera um relatório de alinhamento `ALIGNMENT_REPORT.md` com score de cobertura (%) e lista de gaps encontrados. O relatório é exibido ao usuário como informação (não bloqueia o avanço).

**Critério de aceite:** Relatório de alinhamento gerado com score; gaps identificados; informativo mas não bloqueante.

---

### TASK-08-011 — Notificação de Conclusão das Fases 2 e 3

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Notificar o usuário quando as Tríades das Fases 2 e 3 concluem e aguardam seu feedback |

**Descrição:**
Implementar o sistema de notificação in-app e (opcionalmente) por email quando a Tríade de uma fase conclui e aguarda feedback do usuário. No backend: criar entidade `Notification` com `UserID`, `ProjectID`, `PhaseNumber`, `Type`, `Message`, `Read`, `CreatedAt`. Endpoint `GET /api/v1/notifications` retorna as notificações não lidas. No frontend: badge numérico no ícone de sino no header com contagem de não lidas, dropdown de notificações ao clicar no sino, clique na notificação navega para o painel de feedback da fase correspondente.

**Critério de aceite:** Notificações criadas ao concluir Tríade; badge no header; navegação ao clicar; marcar como lida.

---

### TASK-08-012 — Testes de Validação das Fases 2 e 3

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Validar o fluxo completo das Fases 2 e 3 incluindo o split de trilhos |

**Descrição:**
Testes de integração cobrindo: execução paralela (ambos os trilhos executam e precisam ser aprovados independentemente), validação de pré-requisitos (Fase 2 não inicia sem Fase 1 completa), contexto RAG correto injetado (verificar que SPEC.md da Fase 1 está no prompt do Produtor da Fase 2), feedback em trilho único (feedbacks no trilho Frontend não afetam o trilho Backend), aprovação parcial (Fase só é marcada COMPLETED quando ambos os trilhos estão aprovados).

**Critério de aceite:** Split de trilhos funcional; pré-requisitos validados; contexto RAG verificado; aprovação por trilho independente.

---

# Phase 15 — Fluxo C: Estratégia de Marketing

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-15 |
| **Título** | Fluxo C: Estratégia de Marketing e Campanhas Multi-Canal |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Média |
| **Pré-requisitos** | PHASE-04, PHASE-05, PHASE-06 concluídas |

---

## Descrição Detalhada

O **Fluxo C** é o serviço de criação de **Estratégias de Marketing e Campanhas** para múltiplos canais digitais. Similar ao Fluxo B, pode ser executado de forma independente ou herdando o contexto de um projeto do Fluxo A.

Os agentes desta fase são especialistas em marketing digital, copywriting persuasivo e estratégia de canais pagos e orgânicos. A Tríade opera com: o Produtor cria a estratégia e o material criativo, o Revisor analisa com foco em engajamento, ROI esperado e compliance das plataformas, e o Refinador entrega as campanhas otimizadas.

Os canais suportados incluem: **LinkedIn** (conteúdo orgânico, LinkedIn Ads, artigos), **Instagram** (feed, stories, reels, Instagram Ads), **Google Ads** (Search, Display, YouTube), e outros canais relevantes. Para cada canal, o sistema gera conteúdo adaptado às especificidades de formato, tom e audiência de cada plataforma.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Agentes especializados em Marketing por canal
- ✅ Sistema de herança de contexto do Fluxo A
- ✅ Pipeline de criação de campanhas via Tríade
- ✅ Calendário editorial gerado para cada canal
- ✅ Pack de ativos de campanha (copies, hashtags, CTAs)
- ✅ Score de engajamento esperado por campanha

---

## Funcionalidades Entregues

- **Estratégia Multi-Canal:** Plano de marketing cobrindo LinkedIn, Instagram, Google Ads
- **Calendário Editorial:** Cronograma de publicações por canal com frequência recomendada
- **Pack de Conteúdo:** Copies completos, hashtags, CTAs e sugestões visuais
- **Análise de ROI:** Estimativas de alcance, impressões e conversões esperadas

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Toda a estratégia de marketing e pack de conteúdo são gerados sequencialmente pela Tríade para todos os canais, sem interrupções.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a estratégia multi-canal |
| **Velocidade** | Mais rápido — estratégia gerada de ponta a ponta |
| **Feedback** | Aplicado à estratégia como um todo |
| **Ideal para** | Quando o usuário quer uma visão completa da estratégia antes de revisar |

### 🎯 Executar uma Task Específica ⭐ **Recomendado para Fluxo C**

O usuário seleciona **um ou mais canais individualmente** (LinkedIn, Instagram, Google Ads) para geração separada da estratégia e conteúdo.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por canal — revisão antes de avançar para o próximo |
| **Velocidade** | Mais controlado — cada canal revisado individualmente |
| **Feedback** | Específico para o canal selecionado (tom, formato, segmentação) |
| **Ideal para** | Quando os canais têm audiências e tons muito distintos que merecem revisão individual |

### 🔀 Modo Híbrido

Execute automaticamente a estratégia geral (TASK-15-003) e depois revise task a task o conteúdo gerado para cada canal, aprovando o LinkedIn antes de prosseguir para Instagram e Google Ads.

> 💡 **Dica:** Para o **Fluxo C (Marketing)**, o modo **task a task por canal** é muito recomendado. Cada plataforma tem sua própria voz, formato e audiência — aprovação granular garante que o conteúdo de cada canal seja revisado com o contexto correto. Comece pela estratégia geral e depois revise canal por canal.

---

## Tasks

### TASK-15-001 — Prompts do Agente Marketing Strategist

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts que tornam o Agente Marketing um especialista em estratégia digital e copywriting |

**Descrição:**
Criar system prompts para: **(Produtor — Marketing Strategist):** instruções para criar, dado o produto e o público-alvo, uma estratégia de marketing completa incluindo: defineção de objetivos SMART (Specific, Measurable, Achievable, Relevant, Time-bound), buyer personas detalhadas, mensagens-chave por persona, canais prioritários com justificativa (onde a audiência está), conteúdo adaptado por canal (LinkedIn formal/insights, Instagram visual/stories, Google Search com intenção de busca), plano de conteúdo com frequência, budget sugerido e alocação por canal, KPIs e métricas de sucesso. **(Revisor — CRO & Marketing Specialist):** analisar alinhamento da mensagem com o público, coerência entre canais, compliance das plataformas (políticas de anúncios), realismo das estimativas de ROI, identificar gaps de funil (topo, meio, fundo). **(Refinador)** otimiza e entrega a estratégia completa.

**Critério de aceite:** Prompts geram estratégia com SMART goals, personas e plano de conteúdo; Revisor verifica compliance de plataformas.

---

### TASK-15-002 — Extrator de Contexto do Projeto para Marketing

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Extrair informações relevantes do Fluxo A para o contexto de marketing |

**Descrição:**
Implementar o `MarketingContextExtractor` análogo ao extrator do Fluxo B, mas focado em informações de marketing: nome e tagline do produto, problema resolvido e público que tem esse problema, benefícios principais (não features técnicas), diferenciais competitivos, modelo de negócio e precificação (se definido), mercado-alvo (B2B/B2C, segmento), qualquer menção a concorrentes na Visão do Produto, tom de comunicação desejado (inferido dos prompts do usuário). Compõe um `MarketingBrief` estruturado que é o ponto de partida para o Marketing Strategist.

**Critério de aceite:** Extração focada em informações de marketing; brief com problema/público/diferenciais/modelo de negócio.

---

### TASK-15-003 — Geração de Estratégia de Marketing Completa

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Estruturar o output da Tríade de Marketing em um documento de estratégia profissional |

**Descrição:**
Definir a estrutura do `MARKETING_STRATEGY.md` que o Refinador deve gerar: **Executive Summary** (objetivos, budget sugerido, principais canais), **Análise de Audiência** (buyer personas com dores, objetivos, comportamento online), **Estratégia por Canal** (seção para cada canal com: público específico naquele canal, tipo de conteúdo, frequência, formato), **Messaging Framework** (mensagem central, variações por audiência, tom de voz), **Calendário de Conteúdo** (30 dias de planejamento em formato tabela), **Budget Allocation** (% recomendada por canal), **KPIs por Canal** (impressões, CTR, CAC, ROAS esperados), **Plano de 90 dias** (fases: Brand Awareness → Consideração → Conversão).

**Critério de aceite:** Documento com todas as seções; persona com dores reais; plano de 90 dias estratégico; KPIs quantificados.

---

### TASK-15-004 — Geração de Conteúdo por Canal: LinkedIn

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar um pack completo de conteúdo para LinkedIn (organic e Ads) |

**Descrição:**
O Agente Marketing gera para LinkedIn: **Conteúdo Orgânico:** 10 posts para feed (mix de: thought leadership, caso de uso de produto, behind-the-scenes, dados/estatísticas do setor, anúncio de produto sem ser salesy), 3 artigos longos (LinkedIn Articles) sobre temas relevantes para o ICP (Ideal Customer Profile), **LinkedIn Ads:** 5 variações de copy para Sponsored Content (headline 150 chars, intro text 600 chars, CTA), 3 variações de Message Ads (InMail), targeting sugerido (cargo, indústria, tamanho de empresa), bid strategy recomendada. Cada post com: copy completo, hashtags otimizadas (3-5 por post), melhor horário de publicação, sugestão visual (o que mostrar na imagem — instruções para designer).

**Critério de aceite:** 10 posts orgânicos completos; 5 variações de Ads; targeting definido; hashtags incluídas; sugestão visual para cada post.

---

### TASK-15-005 — Geração de Conteúdo por Canal: Instagram

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar pack de conteúdo para Instagram adaptado a feed, stories e reels |

**Descrição:**
O Agente Marketing gera para Instagram: **Feed:** 12 posts com copy de legenda (máximo 2200 chars, mas idealmente 150-200 para engajamento), hashtags estratégicas (20-25 por post, mix de nicho e trending), **Stories:** 8 sequências de stories (3-5 slides cada) para: apresentar o produto, case de sucesso, FAQ interativo, enquete, **Reels:** scripts para 5 reels de 30-60 segundos com: hook nos primeiros 3 segundos, estrutura do conteúdo, texto na tela, CTA final, **Instagram Ads:** 4 variações de creative copy (headline, primary text, description) para diferentes objetivos (Awareness, Consideration, Conversion). Sugestão de paleta visual por formato.

**Critério de aceite:** Conteúdo adaptado ao formato de cada tipo de post; scripts de Reels com hook; hashtags estratégicas por quantidade correta.

---

### TASK-15-006 — Geração de Campanhas Google Ads

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar estrutura de campanha e copies para Google Ads Search e Display |

**Descrição:**
O Agente Marketing gera para Google Ads: **Search Campaigns:** lista de palavras-chave por intenção (informacional, navegacional, transacional/comercial), agrupadas em Ad Groups temáticos, para cada Ad Group: 3 RSAs (Responsive Search Ads) com 15 headlines (30 chars) e 4 descriptions (90 chars), extensões de anúncio (Sitelinks, Callouts, Structured Snippets), estratégia de negativa de palavras-chave, **Display Campaign:** 5 textos de banner (headline 30 chars, description 90 chars), **YouTube (Se aplicável):** script de anúncio TrueView (5s + 30s version), **Budget Escalonado:** sugestão de budget por fase (teste, escala, otimização) com estimativa de cliques e conversões por fase de maturidade.

**Critério de aceite:** Keywords com intenções definidas; Ad Groups temáticos; RSAs dentro dos limites de caracteres; budget escalonado.

---

### TASK-15-007 — Calendário Editorial Interativo no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir o calendário de conteúdo de forma interativa e exportável |

**Descrição:**
Implementar a aba "Calendário" no painel de projeto do Fluxo C com: visualização em grade de calendário (mês atual + próximos 2 meses), cada dia com os posts planejados para aquele dia, filtro por canal (LinkedIn / Instagram / Google), ao clicar em um dia/post: drawer com copy completo, hashtags, sugestão visual, melhor horário, código de cores por canal (LinkedIn = azul, Instagram = rosê, Google = multicolor), botão de exportar como CSV (para ferramentas de scheduling como Buffer, Hootsuite), botão de exportar como ICS (calendário padrão para Google Calendar/Outlook).

**Critério de aceite:** Calendário visual com filtro por canal; copy completo ao clicar; exportação em CSV e ICS.

---

### TASK-15-008 — Pack de Download de Conteúdo por Canal

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Gerar arquivos organizados por canal prontos para download e uso imediato |

**Descrição:**
Implementar `GET /api/v1/projects/:id/marketing/export` que gera um ZIP organizado por canal: `linkedin/organic/posts.md` (todos os posts), `linkedin/ads/campaigns.md`, `instagram/feed/posts.md`, `instagram/stories/sequences.md`, `instagram/reels/scripts.md`, `instagram/ads/creatives.md`, `google-ads/keywords.csv`, `google-ads/ads.md`, `strategy/MARKETING_STRATEGY.md`, `strategy/CALENDAR.csv`. No frontend, botão de download do pack completo com indicador do total de peças de conteúdo geradas (ex: "57 peças de conteúdo prontas para uso"). Possibilidade de download por canal individual.

**Critério de aceite:** ZIP organizado por canal; CSV de keywords válido; download por canal individual; contador de peças de conteúdo.

---

### TASK-15-009 — Score e Análise de Performance Esperada

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar estimativas de performance para cada canal e campanha gerada |

**Descrição:**
Após conclusão da Tríade, gerar o `PERFORMANCE_FORECAST.md` com estimativas baseadas em benchmarks médios do setor: LinkedIn Ads: CTR esperado (0,3-0,5%), CPL estimado, **Instagram Ads:** CPM, CTR, custo por seguidor, **Google Search:** CPC médio por keyword (baseado no CPC médio do setor identificado na pesquisa de keywords), **Orgânico LinkedIn:** alcance estimado por post (baseado no tamanho médio de rede do ICP). Disclaimer claro de que são estimativas baseadas em benchmarks setoriais, não garantias. Exibir no painel como cards de métricas por canal com faixas (mínimo-médio-otimista).

**Critério de aceite:** Estimativas com faixas (min-médio-otimista); benchmarks setoriais indicados; disclaimer claro; exibição por canal.

---

### TASK-15-010 — Integração com Ferramentas de Marketing (Webhooks)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Permitir integração com ferramentas de automação de marketing via webhook |

**Descrição:**
Implementar um sistema de webhooks que permite que o usuário configure endpoints externos para receber o conteúdo de marketing ao ser gerado: `POST /api/v1/projects/:id/marketing/webhooks` (cadastra URL de webhook), quando o Fluxo C é concluído, o sistema faz POST para a URL configurada com o payload JSON contendo todos os posts e metadados. Implementar validação da URL de webhook (deve responder com 200 em teste de saúde), retry de 3 tentativas em caso de falha, log de entregas (success/failure, timestamp, response status). Casos de uso: enviar automaticamente para um Make.com/Zapier scenario que publica nas redes sociais.

**Critério de aceite:** Webhook configurado e disparado ao concluir; retry em falhas; log de entregas; validação da URL.

---

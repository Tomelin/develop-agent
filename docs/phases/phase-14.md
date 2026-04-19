# Phase 14 — Fluxo B: Criação de Landing Page

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-14 |
| **Título** | Fluxo B: Criação de Landing Page de Alta Conversão |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Média |
| **Pré-requisitos** | PHASE-04, PHASE-05, PHASE-06 concluídas (gestão de projetos e pipeline da Tríade) |

---

## Descrição Detalhada

O **Fluxo B** é o serviço de criação de **Landing Pages de Alta Conversão**. É um fluxo autônomo que pode ser executado de forma independente (sem vínculo com um projeto de Software) ou herdando o contexto rico de um projeto do Fluxo A para acelerar e alinhar a criação da página.

A landing page é desenvolvida pelos agentes especialistas em Design, Copywriting e Desenvolvimento Web. O agente Produtor cria o conceito (estrutura, copy, design visual), o agente Revisor analisa com foco em conversão, SEO e experiência do usuário, e o agente Refinador entrega a versão final impecável.

Quando o Fluxo B herda um projeto do Fluxo A, o sistema extrai automaticamente: proposta de valor, branding, cores, público-alvo e nome do produto para criar uma landing page perfeitamente alinhada ao produto desenvolvido. Isso elimina horas de briefing e garante consistência de marca.

O output final é código HTML/CSS/JavaScript (ou Next.js, dependendo da preferência do usuário) funcional, pronto para ser hospedado.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Agentes especializados em Landing Page (Designer, Copywriter, Dev Web)
- ✅ Sistema de herança de contexto do Fluxo A
- ✅ Pipeline de desenvolvimento da landing page via Tríade
- ✅ Preview em tempo real da landing page sendo gerada
- ✅ Output final em HTML/CSS/JS funcional
- ✅ Análise de métricas de conversão do resultado

---

## Funcionalidades Entregues

- **Agentes Especializados:** Tríade com foco em conversão, UX e código Web
- **Herança de Contexto:** Aproveitamento automático do branding do Fluxo A
- **Preview em Tempo Real:** Visualização da landing page enquanto é gerada
- **Análise de Conversão:** Score de potencial de conversão da página gerada

---

## Modo de Execução

> O usuário tem **controle total da granularidade de execução** desta phase. Ao visualizar a lista de tasks abaixo, escolha como deseja prosseguir:

### 🚀 Executar a Phase Completa

Todo o fluxo de criação da Landing Page é executado sequencialmente pela Tríade, desde o brief até os exports. O usuário aprova o resultado final consolidado.

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Única — ao final de toda a criação |
| **Velocidade** | Mais rápido — landing page gerada de ponta a ponta |
| **Feedback** | Aplicado à landing page como um todo |
| **Ideal para** | Usuários que querem uma primeira versão rápida para iterar depois |

### 🎯 Executar uma Task Específica ⭐ **Recomendado para Fluxo B**

O usuário seleciona **uma ou mais tasks individualmente** da lista abaixo. Útil para revisar seções específicas da landing page (ex: só o hero, só o CTA, só o SEO).

| Aspecto | Detalhe |
|---------|---------|
| **Aprovação** | Individual por task — controle granular sobre cada seção |
| **Velocidade** | Mais controlado — cada seção revisada individualmente |
| **Feedback** | Específico para a seção ou aspecto selecionado |
| **Ideal para** | Ajuste fino de copy, design ou SEO antes de publicar |

### 🔀 Modo Híbrido ⭐ **Mais Usado no Fluxo B**

Execute automaticamente a geração inicial da landing page e depois pause task a task para revisar: o copy do hero, os depoimentos, o CTA final e o score de conversão antes de confirmar o export.

> 💡 **Dica:** Para o **Fluxo B (Landing Page)**, o modo hybrid é o mais natural. Gere a página completa automaticamente e depois use tasks individuais para refinar seções específicas, especialmente headline, CTA e social proof que impactam diretamente a taxa de conversão.

---

## Tasks

### TASK-14-001 — Prompts do Agente Landing Page Designer

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar os prompts especializados para os agentes de criação de landing page |

**Descrição:**
Criar system prompts para: **(Produtor — Landing Page Creator):** instruções para criar landing pages de alta conversão seguindo as seções fundamentais: Hero (headline poderosa, subheadline, CTA principal), Social Proof (depoimentos, logos de clientes, métricas), Features/Benefits, Como Funciona (3 steps), FAQ, CTA Final. Instruções de copywriting: headline com benefício principal, voz ativa, linguagem do público-alvo, urgência sem ser spam, clareza sobre o que o usuário ganha. Instruções de design: paleta do usuário/projeto herdado, tipografia profissional, espaçamento generoso, mobile-first. **(Revisor — CRO Specialist):** analisar o potencial de conversão (headline clara, CTA visível, prova social, sem fricção), SEO on-page (H1 único, meta description, alt texts), performance (imagens otimizadas, CSS mínimo), WCAG 2.1 básico.

**Critério de aceite:** Prompts cobrem todas as seções de landing page eficaz; Revisor foca em métricas de conversão; WCAG básico verificado.

---

### TASK-14-002 — Extrator de Contexto do Projeto Base (Herança)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Extrair automaticamente o contexto relevante de um projeto do Fluxo A para alimentar o Fluxo B |

**Descrição:**
Implementar o `LandingPageContextExtractor` que, quando um projeto do Fluxo B é vinculado a um projeto do Fluxo A, extrai automaticamente: **Do VISION.md (Fase 1):** nome do produto, tagline, problema resolvido, público-alvo, proposta de valor única, diferenciais, **Do SPEC.md consolidado:** as 3-5 funcionalidades mais impactantes para destacar, o modelo de negócio (freemium, assinatura, etc.), integrações relevantes para o usuário final, **Dos prompts do usuário:** paleta de cores e tema (dark/light), tipografia preferida, restrições de design. Compõe um `LandingPageBrief` estruturado que é injetado como primeiro contexto do Produtor da landing page.

**Critério de aceite:** Extração funcional quando vinculado a projeto do Fluxo A; brief estruturado com todas as informações relevantes; funciona sem vínculo (brief manual).

---

### TASK-14-003 — Formulário de Brief Manual (Sem Herança)

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar o formulário de brief para usuários que criam landing pages sem vínculo com um projeto de software |

**Descrição:**
Implementar a tela de criação de projeto do Fluxo B com: toggle "Vincular a projeto existente" vs "Criar do zero", quando "Criar do zero": formulário de brief com campos: nome do produto (obrigatório), problema que resolve (obrigatório), público-alvo (obrigatório), proposta de valor única (obrigatório), 3-5 funcionalidades/benefícios principais (lista dinâmica), seção de estilo: paleta de cores (color pickers), tema claro/escuro, tom de comunicação (Profissional/Moderno/Descontraído/Inspirador), idioma da página (PT-BR, EN-US, ES), tipo de output preferido: HTML/CSS/JS vanilla vs Next.js. Validações obrigatórias antes de iniciar a Tríade.

**Critério de aceite:** Formulário com todos os campos; validações funcionais; toggle entre herança e criação do zero; seleção de stack de output.

---

### TASK-14-004 — Execução da Tríade de Landing Page

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Adaptar a execução da Tríade para o contexto específico de criação de landing page |

**Descrição:**
Adaptar o `TriadOrchestrator` para o Fluxo B: **(Produtor)** gera a landing page completa em código (HTML/CSS/JS ou Next.js com TailwindCSS), com todas as seções, copy e estilização. O output é código funcional completo — não um rascunho. **(Revisor)** analisa: taxa de conversão esperada (10 critérios de CRO), SEO on-page, acessibilidade WCAG 2.1 básico, responsividade (executa mentalmente a visualização em mobile), performance (identifica imagens não otimizadas, CSS pesado). Gera lista de melhorias priorizadas. **(Refinador)** aplica todas as melhorias e entrega código de produção. Limite de **5 feedbacks** por sessão. Armazenar o HTML final como artefato do projeto.

**Critério de aceite:** Tríade gera código funcional completo; Revisor analisa CRO e acessibilidade; output é HTML/CSS/JS válido.

---

### TASK-14-005 — Preview em Tempo Real da Landing Page

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Exibir um preview renderizado da landing page enquanto ela está sendo gerada pelo agente |

**Descrição:**
Implementar o componente `LandingPagePreview` que: ao receber o código HTML do SSE de streaming, renderiza progressivamente em um `<iframe>` sandboxed, exibe o preview em modo "splitscreen" — terminal à esquerda com código sendo gerado, preview à direita mostrando o resultado visual em tempo real, toggles de resolução: Desktop (1440px), Tablet (768px), Mobile (375px), ao finalizar a geração, exibe botão "Abrir em aba completa" para visualização em fullscreen, ao aprovar o output da Tríade, habilita o botão "Download da Landing Page" (ZIP com HTML/CSS/JS/imagens).

**Critério de aceite:** Preview renderizado em iframe sandboxed; split view código/preview; simulação de resolução; download ZIP.

---

### TASK-14-006 — Score de Análise de Conversão

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Gerar um score quantificado do potencial de conversão da landing page entregue |

**Descrição:**
Após a conclusão da Tríade, o sistema usa um agente especializado (leve) para analisar a landing page final e gerar um `CONVERSION_REPORT.md` com score de 0-100 baseado em: headline clara e focada no benefício (0-15 pts), subheadline complementar e persuasiva (0-10 pts), CTA visível acima da dobra (0-15 pts), prova social presente (0-10 pts), proposta de valor clara em menos de 5 segundos de leitura (0-15 pts), ausência de distrações ou opções paralelas (0-10 pts), responsividade mobile adequada (0-10 pts), velocidade de carregamento estimada (0-10 pts), SEO básico (H1, meta description, alt texts) (0-5 pts). Exibir o score como headline no painel do projeto com detalhamento por critério.

**Critério de aceite:** Score quantificado com detalhamento por critério; relatório MD gerado; exibido no painel do projeto.

---

### TASK-14-007 — Geração de Múltiplas Variantes (A/B Test Ready)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Gerar variantes da landing page para testes A/B, variando headline e CTA |

**Descrição:**
Implementar o `VariantGenerator` que, após a landing page principal ser aprovada, pode gerar (sob demanda do usuário) até 3 variantes alternativas com: variação de headline (3 abordagens diferentes — benefit-focused, problem-focused, curiosity-driven), variação de CTA (3 textos alternativos de botão), variação de hero image/illustration concept (3 estilos). Cada variante é uma landing page completa, não apenas snippets. No frontend, seção "Variantes A/B" com: lista de variantes geradas, preview individual, botão de download por variante, comparação side-by-side dos scores de conversão.

**Critério de aceite:** Até 3 variantes geradas por demanda; cada variante é código completo; comparação de scores entre variantes.

---

### TASK-14-008 — Otimização SEO On-Page

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Verificar e otimizar automaticamente todos os elementos de SEO on-page da landing page |

**Descrição:**
Implementar um passo de pós-processamento de SEO que: verifica H1 único e contendo a keyword principal, meta title entre 50-60 caracteres, meta description entre 150-160 caracteres e orientada à conversão, OG tags presentes (og:title, og:description, og:image), Twitter Card tags, Schema.org markup para o tipo de produto (Product, Organization ou Software), alt text em todas as imagens, links internos com texto descritivo, canonical tag configurada. Gera o `SEO_CHECKLIST.md` com status de cada item. O Refinador aplica os ajustes de SEO faltantes antes de finalizar.

**Critério de aceite:** Checklist SEO completo; H1, meta tags e OG tags presentes; Schema.org markup correto; SEO_CHECKLIST.md gerado.

---

### TASK-14-009 — Export e Deploy da Landing Page

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Permitir que o usuário exporte a landing page pronta para hospedagem imediata |

**Descrição:**
Implementar `GET /api/v1/projects/:id/landing-page/export` com opções: **HTML Bundle:** ZIP com index.html, style.css, script.js, imagens (se geradas) — prontos para hospedar em qualquer servidor estático, **Vercel Deploy:** gera o `vercel.json` de configuração e instrução de `vercel deploy`, **Netlify:** gera o `netlify.toml` e instrução de deploy, **GitHub Pages:** gera o workflow `.github/workflows/pages.yml` para deploy automático. No frontend, dropdown "Exportar para" com as 4 opções, cada uma com instruções passo a passo de como hospedar.

**Critério de aceite:** ZIP da landing page funcional; configurações de Vercel, Netlify e GitHub Pages gerados; instruções de deploy claras.

---

### TASK-14-010 — Histórico e Versionamento de Landing Pages

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Frontend |
| **Objetivo** | Manter histórico de versões da landing page para permitir rollback e comparação entre versões |

**Descrição:**
Implementar o versionamento de landing pages: cada ciclo completo da Tríade (ou iteração de feedback) gera uma nova versão numerada da landing page. Armazenar cada versão com: número de versão, código HTML completo, score de conversão, data de geração, status (DRAFT/APPROVED). Implementar `GET /api/v1/projects/:id/landing-page/versions` que lista todas as versões. No frontend: lista de versões com preview thumbnail, score de conversão por versão, botão de comparar versão A vs versão B lado a lado, botão de "Usar esta versão como base para próxima iteração".

**Critério de aceite:** Versionamento por iteração; preview thumbnail por versão; comparação lado a lado; rollback para versão anterior.

---

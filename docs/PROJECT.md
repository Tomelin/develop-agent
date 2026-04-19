# 📝 Estruturação Oficial do Projeto: Agência de IA

## Visão Geral
Uma **esteira de produção de software inteligente e automatizada**, onde cada "profissional" é um agente especialista de IA. O sistema garante rastreabilidade total da ideia à entrega, aplicando ciclos de revisão obrigatórios para assegurar consistência, segurança e qualidade de alto nível no ciclo de desenvolvimento (SDLC).

---

## 1. O Motor Principal: A Tríade de Agentes
Em todas as etapas e fluxos de desenvolvimento (Software, Landing Page ou Marketing), o sistema opera em um conceito de **três agentes trabalhando em conjunto** para garantir a qualidade máxima de cada entrega:

- **Agente Produtor:** Executa a tarefa inicial, desenvolvendo o artefato raiz (código estrutural, texto da funcionalidade, rascunho de arquitetura).
- **Agente Revisor:** Analisa o trabalho do Produtor contra as regras de negócio definidas, boas práticas e requisitos da aplicação. Ele **nunca altera o artefato**. O revisor apenas emite críticas e correções (Feedbacks). Assume uma persona paranoica e crítica (Engenheiro de Qualidade Sênior cético), não aprovando nada que não esteja perfeito.
- **Agente Refinador (Entregador):** Pega o artefato original e interage ativamente com as críticas do Revisor, aplicando os ajustes necessários para produzir a versão final impecável da entrega.

---

## 2. Os Três Grandes Fluxos de Serviço

### Fluxo A: Desenvolvimento de Software (End-to-End)
Dividido em 9 fases estruturais. A partir da Fase 2 (Engenharia), a trilha de desenvolvimento divide as execuções em dois caminhos paralelos para maior velocidade: **Frontend** e **Backend**.

*   **Fase 1: Criação do Projeto (Agente Entrevistador)**
    *   **Objetivo:** Transformar uma ideia bruta em uma visão clara de produto.
    *   **Comportamento:** Age de forma estritamente conversacional. Faz perguntas, escuta e reformula o que foi recebido para validar o entendimento do projeto. Não avança para a próxima etapa sem consolidação sólida e o aval do usuário.
    *   **Limites de Interação:** O usuário pode enviar até **10 feedbacks** de refinamento para o projeto base.
*   **Fase 2: Engenharia de Software**
    *   Definição abrangente de todas as regras de negócio, e levantamento dos requisitos funcionais e não-funcionais da plataforma.
*   **Fase 3: Arquitetura de Software**
    *   Modelagem de dados robusta, definição da stack tecnológica (linguagens/ferramentas), de infraestrutura e aplicação de design patterns de projeto.
*   **Fase 4: Planejamento (Roadmap)**
    *   Divisão técnica da arquitetura construída em fases menores, épicos, e tarefas assinaláveis.
    *   **Outputs Estruturados:** O agente é instruído a gerar um formato `JSON` determinístico contendo as Tasks e estimativas de sua complexidade para o core em Golang ingerir os objetos e materializá-los em cards de um KanBan para o visual do usuário.
*   **Fase 5: Desenvolvimento**
    *   Mão no código. A atuação acontece com os agentes programadores focando diretamente e de forma ramificada entre construção Front e Back.
*   **Fase 6: Testes de Software**
    *   Testes unitários, abordagem guiada por regras de TDD e automação de integração (CI).
*   **Fase 7: Segurança**
    *   Auditoria severa contra vulnerabilidades do código gerado (OWASP), garantindo blindagem aos dados da camada principal.
*   **Fase 8: Documentação**
    *   Geração exata e atual de manuais operacionais, tutoriais de utilização, referências da documentação de APIs, e READMEs atrativos para repositório.
*   **Fase 9: DevOps e Deploy (Infra como Código) - Bônus Evolutivo**
    *   Agente DevOps varre arquitetura (Fase 3) e o Source desenvolvido (Fase 5) para arquitetar a infraestrutura de modo automático gerando e disponibilizando os Dockerfiles necessários, testes de docker-compose, github-actions (.pipelines) e manifestos de subida no Kubernetes (`.yaml`).

> 📋 **Regra de Ciclos (Fase 2 à 8):** Diferente da primeira fase, nas Fases seguintes o cliente/usuário possui o teto de **até 5 feedbacks** limitados para as requisições de melhorias da rodada.

### Fluxo B: Criação de Landing Page
Serviço autônomo focado na rápida geração de páginas captadoras e de alta conversão do mercado publicatório.
*   Pode efetuar a herança autônoma do projeto gerado na Fase 1 (do Fluxo A). Ele usa a densa base de branding, propostas de valor informadas para acelerar, baratear e alinhar o perfil as instruções da nova landing page desejada.
*   Usa obrigatoriamente a premissa da avaliação da **Tríade de Agentes**.

### Fluxo C: Estratégia de Marketing
Serviço voltado para criação de campanhas inteligentes dedicadas às redes e mecanismos de tráfego, como LinkedIn, Instagram e Google Ads.
*   Idêntico às Landing pages, herda sem atritos os dados do escopo e arquitetura do projeto raiz.
*   Submete campanhas ao circuito crítico da **Tríade de Agentes** focada no engajamento para liberação do conteúdo.

---

## 3. Customização Avançada, Prompts Sistêmicos e Controle de Contextos

### Gestão do Perfil do Usuário e Contextos
*   **Painéis de Orientação:** A plataforma inicia com um modelo padrão (`admin`) oferecendo acesso aos bancos de prompts agrupados. 
*   **Regras de Domínio Agregadas:** Se o usuário determina preferências base de criação, elas passam a ser vinculadas a todas as requisições de sistema - injetando o prompt global na base de diretrizes enviadas aos Modelos - impedindo distorções arquiteturais (*ex: "Desenvolver utilizando Golang Back-end, design no NextJS; Dark Mode e Paleta em Tons Escuros e Amarelo"*). Estilização corporativa sempre honrada.

### Personalização Completa: O Catálogo de Agentes
Cada Agente AI listado para participação na Agência pode ser cadastrado do zero ou configurado, preenchendo as seguintes propriedades:
1.  **Nome e Descrição:** (Ex: "Bob, Arquiteto Sênior e Especialista Cloud")
2.  **Modelo de AI:** (Qual motor de linguagem impulsiona seu processamento na API, Ex: Claude-3-Opus, GPT-4o, Gemini-1.5-Pro).
3.  **Lista de System Prompts:** Comandando a identidade base da persona do agente.
4.  **Lista de Habilidades (Skills/Tarefas):** Arestas do workflow suportadas.
5.  **Toggle Habilitador:** Switcher Boolean `Enabled` determinando o status operacional na plataforma.

**🧠 Visão Multi-Modelo (Modo Dinâmico):** Opção inovadora da Agência para contornar vieses dos modelos. Se habilitado na fase de processamento, o backend sorteia agentes diferentes para a Tríade da fase (Ex: O Claude checa o código que o ChatGPT cometeu, que é então refinado pelo Gemini), assegurando visões cognitivas plurais.

---

## 4. Otimizações de Nível Profissional (Melhorias Contínuas)

Para escalar o produto visando mercado institucional consolidado e uso constante, temos o seguinte pilar de tecnologias agregadas:

*   **⚡ Handoff e Gestor de Memória (RAG) entre Fases:** Transitar o escopo da Fase 1 para instruir as últimas não causa sobrecargas de preço das APIs graças à consolidação. Um modelo secundário extremamento ágil avalia a fase terminada e empacota toda a vivência em um manifesto resumido (`SPEC.md`). A fase superior assimila essa sinopse que mantém vivo perfeitamente o seu histórico essencial, otimizando dezenas milhares de Prompt Tokens.
*   **🔂 Gatilhos de Rejeição Automática Anti-Despesas:** Se Segurança ou Testes encontrarem falhas catastróficas nas entregas, o framework desarma o Agente com poder de devolver em Background todo o repositório de volta para Fase 5 (Dev) instruindo as correções vitais *sem* descontar dos 5 feedbacks manuais preciosos do cliente, resolvendo problemas estruturais internamente.
*   **📊 Painel Visível e Auditoria de Finanças/Billing de LLMs:** Controle transparente dos custos da esteira dinâmica de sorteios aleatórios de Agentes. Centralizando log detalhado das informações cobradas por *Prompt Tokens* e *Completion Tokens*, mapeando centavos em requisições de alta densidade por ID de Projeto, viabilizando repasses monetários coesos para usuários.

---

## 5. Granularidade de Execução: Phase Completa ou Task Individual

O usuário tem **controle total sobre o nível de granularidade** com que deseja acionar a esteira de agentes. Ao visualizar o painel de uma phase, ele pode escolher entre dois modos de execução:

### 🚀 Executar a Phase Completa

Todas as tasks da phase são desenvolvidas sequencialmente pela Tríade de Agentes, sem interrupções. O sistema avança automaticamente de task em task até concluir toda a phase e aguarda a aprovação do usuário uma única vez ao final.

*   **Ideal para:** Phases bem compreendidas onde o usuário confia na execução automática e prioriza velocidade.
*   **Aprovação:** Uma única ao final de toda a phase, dentro do limite de feedbacks daquela fase.
*   **Feedback:** Aplicado à phase como um todo.

### 🎯 Executar uma Task Específica

O usuário seleciona **uma ou mais tasks individualmente** da lista da phase. A Tríade desenvolve apenas a(s) task(s) selecionada(s) e aguarda aprovação explícita antes de prosseguir.

*   **Ideal para:** Tasks críticas, complexas ou que exigem revisão minuciosa antes de impactar as tarefas seguintes.
*   **Aprovação:** Individual por task — o usuário controla o ritmo do pipeline.
*   **Feedback:** Granular e específico para cada task.

### 🔀 Modo Híbrido

Os dois modos podem ser **combinados livremente** dentro da mesma phase: inicie a phase em execução automática completa e pause manualmente em uma task específica que exige atenção. Após aprovar aquela task individualmente, o modo automático retoma a partir da próxima task.

> 💡 **Recomendação:** Use o modo automático (phase completa) para phases de Documentação (Fase 8) e DevOps (Fase 9), e prefira o modo granular por task nas phases de Desenvolvimento (Fase 5) e Segurança (Fase 7), onde cada entregável impacta diretamente as tasks seguintes.

---

## 6. AgentSDK — Contrato de Interface para Conexão com Modelos de IA

O **core de negócio da API** que conecta a plataforma com os modelos de IA **obrigatoriamente implementa o padrão de interface** (`Provider`). Nenhum código de domínio ou orquestração pode depender diretamente de um SDK específico (OpenAI, Anthropic, Gemini, Ollama). A dependência é sempre feita em cima da **abstração**, nunca da implementação concreta.

### Por que Interface?

*   **Desacoplamento:** Trocar um provider de LLM não requer mudança no core da plataforma — apenas o provider concreto é substituído.
*   **Testabilidade:** Testes unitários e de integração usam um `MockProvider` determinístico, sem custo de API.
*   **Modo Dinâmico:** O sorteio de agentes na Tríade funciona porque todos os providers satisfazem o mesmo contrato — o orquestrador não sabe (nem precisa saber) qual provider ele está chamando.
*   **Observabilidade:** O `BillingTracker` e o `RetryWrapper` decoram a interface sem alterar as implementações concretas (*Decorator Pattern*).

### Estrutura de Pacotes

```
src/backend/pkg/agentsdk/
├── provider.go          # Contrato (interface Provider + todos os tipos)
├── openai/
│   └── provider.go      # Implementação concreta OpenAI
├── anthropic/
│   └── provider.go      # Implementação concreta Anthropic
├── gemini/
│   └── provider.go      # Implementação concreta Google Gemini
└── ollama/
    └── provider.go      # Implementação concreta Ollama (modelos locais)
```

### Regra de Uso Obrigatória

```go
// ✅ CORRETO — o domínio depende da abstração
var p agentsdk.Provider = anthropic.New()
p.Initialize(ctx, cfg)
resp, _ := p.Complete(ctx, req)

// ❌ PROIBIDO — o domínio não pode depender da implementação concreta
import "github.com/anthropics/anthropic-sdk-go"
client := sdk.NewClient(...) // acoplamento direto ao SDK
```

> ⚠️ **Regra de Ouro:** Todo o código em `src/backend/domain/` e `src/backend/usecase/` deve referenciar exclusivamente `agentsdk.Provider`. As implementações concretas (`openai`, `anthropic`, `gemini`, `ollama`) só são instanciadas na camada de injeção de dependências (`src/backend/infra/` ou `main.go`).

---

## 7. Comunicação entre Agentes — Protocolo de Mensagens

Toda troca de informação entre os agentes da plataforma (Produtor → Revisor → Refinador, bem como interações com o orquestrador) é feita através de **Go Channels tipados** transportando um **envelope JSON padronizado**. Nenhum agente acessa o output de outro diretamente — todo acesso é mediado pelo channel.

### Envelope de Mensagem (`AgentMessage`)

```go
// AgentMessage é o envelope padrão para toda comunicação entre agentes.
// Transportado exclusivamente via channels — nunca via chamada direta.
type AgentMessage struct {
    ID        string    `json:"id"`         // UUID único da mensagem
    From      string    `json:"from"`       // Identificador do agente de origem (e.g. "producer-abc123")
    To        string    `json:"to"`         // Identificador do agente de destino (e.g. "reviewer-xyz456")
    Message   string    `json:"message"`    // Payload — conteúdo da mensagem (artefato, crítica, refinamento)
    Status    string    `json:"status"`     // Estado atual: "pending" | "processing" | "done" | "error"
    Timestamp time.Time `json:"timestamp"`  // Momento de criação da mensagem
    Meta      map[string]any `json:"meta,omitempty"` // Metadados opcionais (phase, token_usage, etc.)
}
```

### Transporte via Go Channels

```go
// AgentChannel é o canal de comunicação de um agente.
// Criado automaticamente na instanciação do agente;
// destruído automaticamente quando o agente é encerrado.
type AgentChannel struct {
    In  <-chan AgentMessage // Canal de entrada — mensagens recebidas pelo agente
    Out chan<- AgentMessage // Canal de saída — mensagens enviadas pelo agente
}
```

**Regras de ciclo de vida:**

| Evento | Comportamento |
|--------|---------------|
| Agente criado | Channel criado automaticamente com buffer padrão (`agent.channel_buffer_size = 10`) |
| Agente destruído | Channel fechado e removido automaticamente |
| Buffer cheio | Mensagem enfileirada — o remetente bloqueia até haver espaço (comportamento configurável) |
| Agente em erro | Channel pausado — mensagens pendentes preservadas até o agente se recuperar |

### Configuração

```yaml
# src/backend/config/config.yaml
agent:
  channel_buffer_size: 10   # tamanho do buffer do channel (padrão: 10, configurável)
  channel_drain_timeout: 30  # segundos para drenar mensagens pendentes antes de destruir
```

### Status do Agente (visível no Dashboard)

| Status | Descrição |
|--------|-----------|
| `IDLE` | Agente ativo e aguardando mensagens |
| `RUNNING` | Processando uma mensagem |
| `PAUSED` | Temporariamente suspenso pelo orquestrador |
| `QUEUED` | Há mensagens no buffer aguardando processamento |
| `ERROR` | Falha no processamento — requer atenção |
| `COMPLETED` | Tarefa concluída — agente aguarda encerramento |

> ⚠️ **Regra de Ouro:** A comunicação direta entre agentes é **proibida**. Todo dado transferido entre agentes DEVE passar pelo `AgentMessage` enviado via channel. Isso garante rastreabilidade total, possibilidade de replay e observabilidade da fila de cada agente.

---

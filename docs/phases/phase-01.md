# Phase 01 — Fundação e Infraestrutura da Plataforma

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-01 |
| **Título** | Fundação e Infraestrutura da Plataforma |
| **Tipo** | Backend + DevOps |
| **Prioridade** | Crítica |
| **Pré-requisitos** | Nenhum — esta é a fase inicial |

---

## Descrição Detalhada

Esta é a fase zero da construção da plataforma. Aqui estabelecemos toda a fundação técnica sobre a qual o sistema da Agência de IA será construído. O objetivo é ter um ambiente de desenvolvimento padronizado, um servidor HTTP funcional com roteamento básico, conexões com todos os serviços de infraestrutura (banco de dados, cache, mensageria) e o padrão arquitetural definido e documentado para que todas as fases subsequentes sigam o mesmo caminho.

A fase garante que a equipe de desenvolvimento (humana ou por agentes de IA) tenha um ponto de partida sólido, sem débitos técnicos arquiteturais, respeitando as escolhas de stack definidas no PROJECT.md: **Golang com GIN** para o backend, **MongoDB** para persistência, **Redis** para cache e **RabbitMQ** para mensageria assíncrona entre agentes.

Nesta fase também é configurado o ambiente de desenvolvimento local via Docker Compose, permitindo que qualquer desenvolvedor ou agente suba o ambiente com um único comando.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Repositório inicializado com estrutura de diretórios padronizada (`src/backend/`, `src/frontend/`)
- ✅ Servidor HTTP GIN funcional com health check respondendo
- ✅ Conexão ativa com MongoDB, Redis e RabbitMQ
- ✅ Sistema de configuração via YAML + variáveis de ambiente (Viper)
- ✅ Middleware de logging estruturado, CORS e recovery configurados
- ✅ Docker Compose para ambiente de desenvolvimento local
- ✅ Makefile com comandos de desenvolvimento padronizados
- ✅ Documentação de arquitetura inicial (ADR - Architecture Decision Record)

---

## Funcionalidades Entregues

- **API Gateway:** Ponto de entrada único via GIN com roteamento versionado (`/api/v1/`)
- **Health Check:** Endpoint `/health` que valida conectividade com todos os serviços
- **Config Management:** Carregamento de configurações via Viper com fallback para env vars
- **Infrastructure Adapters:** Wrappers para MongoDB, Redis e RabbitMQ com injeção de dependência
- **Observabilidade Básica:** Logging estruturado em JSON com nível configurável

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

> 💡 **Dica:** Para esta phase, o modo recomendado é **task a task**, pois cada entregável de infraestrutura influencia diretamente as configurações das tasks seguintes.

---

## Tasks

### TASK-01-001 — Inicialização do Repositório e Estrutura de Diretórios

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar a estrutura de diretórios canônica do projeto conforme definido no PLAYBOOK.md |

**Descrição:**
Inicializar o repositório Git e criar a estrutura de diretórios completa seguindo a convenção obrigatória do projeto. Todo código backend deve residir em `src/backend/` e todo código frontend em `src/frontend/`. A estrutura interna do backend deve seguir Clean Architecture com separação clara entre `domain`, `usecase`, `infra` e `api`. Criar o `go.mod` com o module path correto e o `.gitignore` adequado para projetos Go.

**Critério de aceite:** `go mod init` executado, estrutura de pastas criada conforme spec, `.gitignore` configurado.

---

### TASK-01-002 — Configuração do Sistema de Configuração com Viper

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar carregamento de configurações via arquivo YAML com suporte a override por variáveis de ambiente |

**Descrição:**
Implementar o pacote `src/backend/config/` utilizando `github.com/spf13/viper`. O sistema deve carregar um arquivo `config.yaml` na raiz do projeto e permitir que qualquer configuração seja sobrescrita via variável de ambiente com prefixo `APP_`. Estruturar as configurações em seções lógicas: `server`, `database`, `redis`, `rabbitmq`, `llm`. Incluir valores default para todas as configurações não-críticas. Implementar validação das configurações obrigatórias na inicialização, falhando com erro descritivo se estiverem ausentes.

**Critério de aceite:** `config.Load()` retorna struct populado; override via env var funciona; erro claro se configuração crítica ausente.

---

### TASK-01-003 — Implementação do Servidor HTTP com GIN

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Configurar o framework GIN com middlewares essenciais e roteamento versionado |

**Descrição:**
Implementar o servidor HTTP utilizando `github.com/gin-gonic/gin`. Configurar os middlewares obrigatórios: `gin.Recovery()` para capturar panics sem derrubar o servidor, `CORS` configurável via config (origens permitidas, métodos, headers), `RequestID` para rastreabilidade de cada requisição com UUID único injetado nos headers e no contexto. Criar o grupo de rotas `/api/v1/` para versionamento da API. O servidor deve inicializar com graceful shutdown ao receber SIGTERM/SIGINT, aguardando até 30 segundos para requisições em andamento completarem.

**Critério de aceite:** Servidor sobe na porta configurada; middlewares aplicados; graceful shutdown funcional.

---

### TASK-01-004 — Adapter de Conexão com MongoDB

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o adapter de MongoDB com connection pooling e operações CRUD genéricas |

**Descrição:**
Implementar o pacote `src/backend/infra/database/mongodb/` utilizando o driver oficial `go.mongodb.org/mongo-driver`. O adapter deve gerenciar um connection pool configurável (min/max connections, timeout de conexão, timeout de operação). Implementar uma interface `Repository` genérica com métodos `FindOne`, `FindMany`, `InsertOne`, `UpdateOne`, `DeleteOne` e `Aggregate`. O adapter deve reconectar automaticamente em caso de perda de conexão. Implementar health check que valida se o MongoDB está acessível e com latência aceitável.

**Critério de aceite:** Conexão estabelecida com retry automático; CRUD básico funcional; health check respondendo.

---

### TASK-01-005 — Adapter de Conexão com Redis

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar o adapter de Redis para cache, sessões e controle de estado |

**Descrição:**
Implementar o pacote `src/backend/infra/cache/redis/` utilizando `github.com/redis/go-redis/v9`. O adapter deve suportar operações básicas de cache com TTL configurável: `Set`, `Get`, `Delete`, `Exists`, `SetNX` (para locks distribuídos). Implementar serialização/deserialização JSON automática para structs Go. O adapter deve suportar pipeline para operações em lote quando performance for crítica. Incluir health check com verificação de latência.

**Critério de aceite:** Operações de cache funcionais; TTL respeitado; pipeline funcional; health check OK.

---

### TASK-01-006 — Adapter de Mensageria com RabbitMQ

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar producer e consumer de mensagens para comunicação assíncrona entre fases do pipeline |

**Descrição:**
Implementar o pacote `src/backend/infra/messaging/rabbitmq/` utilizando `github.com/rabbitmq/amqp091-go`. Criar abstração de `Publisher` (para envio de mensagens a filas e exchanges) e `Consumer` (para processamento de mensagens com ack/nack explícito). Implementar reconexão automática com backoff exponencial em caso de queda do broker. Criar a topologia de exchanges e filas do sistema: uma exchange principal `agency.events` do tipo `topic`, com filas para cada fase do pipeline (`phase.1`, `phase.2`, ... `phase.9`). Implementar dead-letter queue para mensagens que falharem após 3 tentativas.

**Critério de aceite:** Publisher envia mensagem; Consumer processa e dá ack; reconexão automática funcional; DLQ configurada.

---

### TASK-01-007 — Implementação do Health Check Endpoint

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Criar endpoint `/health` que valida o estado de todos os serviços de infraestrutura |

**Descrição:**
Implementar o handler `GET /health` que verifica e reporta o status de conectividade com MongoDB, Redis e RabbitMQ. A resposta deve seguir o padrão de health check com status `healthy`/`degraded`/`unhealthy` para cada componente, latência de cada verificação em milissegundos e status geral do sistema. Retornar HTTP 200 quando todos os serviços estão saudáveis, HTTP 503 quando qualquer serviço está indisponível. Implementar também `/health/live` (apenas verifica se o processo está vivo) e `/health/ready` (verifica se o processo está pronto para receber tráfego).

**Critério de aceite:** `/health` reporta status de todos os serviços; `/health/live` sempre retorna 200 se o processo está vivo; `/health/ready` retorna 503 se algum serviço está down.

---

### TASK-01-008 — Configuração do Docker Compose para Desenvolvimento

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + DevOps |
| **Objetivo** | Criar ambiente de desenvolvimento completo com um único comando via Docker Compose |

**Descrição:**
Criar o arquivo `docker-compose.yml` na raiz do projeto com todos os serviços necessários para desenvolvimento local: MongoDB (com autenticação habilitada e volume persistente), Redis (com persistência RDB), RabbitMQ (com management plugin habilitado na porta 15672). Criar um serviço `backend` que compila e executa a aplicação Go com hot-reload via `air`. Configurar networks isoladas para cada grupo de serviços. Adicionar `docker-compose.override.yml` para configurações específicas de desenvolvimento (ports expostas, volumes de source code). Criar `.env.example` documentando todas as variáveis de ambiente necessárias.

**Critério de aceite:** `docker-compose up` sobe todos os serviços sem erros; backend conecta a todos os serviços; hot-reload funcional.

---

### TASK-01-009 — Makefile com Comandos de Desenvolvimento

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + DevOps |
| **Objetivo** | Padronizar os comandos de desenvolvimento, build e teste com um Makefile |

**Descrição:**
Criar `Makefile` na raiz do projeto com os alvos essenciais: `make dev` (sobe o ambiente Docker Compose), `make build` (compila o binário para produção), `make test` (executa todos os testes unitários com cobertura), `make test-integration` (testes de integração com ambiente docker), `make lint` (executa golangci-lint), `make migrate` (executa migrações do banco), `make docs` (gera documentação Swagger), `make clean` (remove artefatos de build). Documentar cada alvo com comentários `##` para que `make help` exiba a lista de comandos disponíveis.

**Critério de aceite:** Todos os alvos do Makefile executam sem erro; `make help` lista todos os comandos com descrições.

---

### TASK-01-010 — Configuração do Logger Estruturado

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar logging estruturado em JSON com contexto de request rastreável |

**Descrição:**
Implementar o pacote de logging em `src/backend/pkg/logger/` utilizando `go.uber.org/zap`. O logger deve emitir logs em formato JSON (produção) ou human-readable colorido (desenvolvimento), com nível configurável via config. Criar middleware GIN que injeta o logger no contexto de cada request com o `request_id` pré-preenchido, para que todos os logs de handlers carreguem o identificador da requisição automaticamente. Implementar campos padrão em todos os logs: `timestamp`, `level`, `service`, `request_id`, `phase_id` (quando aplicável), `agent_id` (quando aplicável).

**Critério de aceite:** Logs em JSON em produção; `request_id` presente em todos os logs de uma requisição; nível configurável.

---

### TASK-01-011 — Implementação da Estrutura de Erros Padronizada

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir um sistema de erros de domínio e de API que padronize respostas de erro para toda a plataforma |

**Descrição:**
Implementar o pacote `src/backend/pkg/errors/` com tipos de erros de domínio (`NotFoundError`, `ValidationError`, `ConflictError`, `UnauthorizedError`, `InternalError`) e um handler de erros global para o GIN que converte erros de domínio em respostas HTTP padronizadas. A resposta de erro deve seguir o schema: `{"error": {"code": "NOT_FOUND", "message": "...", "details": {...}}}`. Garantir que erros internos nunca exponham stack traces ou detalhes sensíveis em produção.

**Critério de aceite:** Erros de domínio mapeados para HTTP correto; schema de resposta consistente; stack traces não expostos em produção.

---

### TASK-01-012 — Configuração do Linter e Qualidade de Código

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Configurar análise estática de código para garantir qualidade desde o primeiro commit |

**Descrição:**
Configurar `golangci-lint` com um arquivo `.golangci.yml` habilitando os linters mais importantes para o projeto: `errcheck` (erros não verificados), `govet` (análise do compilador), `staticcheck` (bugs e problemas de performance), `gofmt` (formatação), `goimports` (organização de imports), `gosec` (segurança), `exhaustive` (switch exhaustiveness). Integrar a execução do linter no Makefile e criar um GitHub Actions workflow que executa o lint em todo PR.

**Critério de aceite:** `make lint` passa sem erros no código base inicial; GitHub Actions configurado.

---

### TASK-01-013 — Criação do Architecture Decision Record (ADR)

| Campo | Valor |
|-------|-------|
| **Camada** | Backend + Documentação |
| **Objetivo** | Documentar as decisões arquiteturais tomadas nesta fase para rastreabilidade histórica |

**Descrição:**
Criar o diretório `docs/adr/` com os primeiros Architecture Decision Records da plataforma. Criar ADRs para: escolha do GIN como framework HTTP (vs Echo, Fiber), escolha do MongoDB como banco principal (vs PostgreSQL), escolha do RabbitMQ para mensageria (vs Kafka), adoção da Clean Architecture no backend, adoção do Viper para configuração. Cada ADR deve seguir o formato padrão: Status, Contexto, Decisão, Consequências.

**Critério de aceite:** ADRs criados com formato padronizado; decisões justificadas com prós e contras.

---

# Phase 02 — Autenticação e Gestão de Usuários

## Visão Geral

| Campo | Valor |
|-------|-------|
| **ID** | PHASE-02 |
| **Título** | Autenticação e Gestão de Usuários |
| **Tipo** | Backend + Frontend |
| **Prioridade** | Crítica |
| **Pré-requisitos** | PHASE-01 concluída (infraestrutura base) |

---

## Descrição Detalhada

Esta fase implementa o sistema de identidade da plataforma. Conforme definido no PROPOSAL.md e PROJECT.md, a plataforma inicia com um modelo de usuário único chamado `admin`, mas deve ser arquitetada de forma que suporte múltiplos usuários no futuro sem refatoração estrutural.

O sistema de autenticação utiliza JWT (JSON Web Tokens) para sessões stateless, com tokens de acesso de curta duração e refresh tokens de longa duração armazenados de forma segura. O perfil de usuário é o coração da personalização da plataforma, pois é nele que ficam os prompts contextuais que serão injetados em todas as fases do pipeline de agentes.

O frontend desta fase consiste na tela de login, tela de perfil do usuário e gerenciamento básico de conta. A integração entre frontend e backend via axios deve ser configurada aqui, com interceptors para refresh automático de token.

---

## Delivery

Ao final desta fase, a plataforma deverá ter:

- ✅ Endpoint de login com validação de credenciais e emissão de JWT
- ✅ Endpoint de refresh de token
- ✅ Middleware de autenticação protegendo todas as rotas privadas
- ✅ CRUD completo de perfil de usuário
- ✅ Tela de login funcional no frontend
- ✅ Tela de perfil do usuário no frontend
- ✅ Seed de dados com usuário `admin` criado automaticamente

---

## Funcionalidades Entregues

- **Autenticação JWT:** Login com email/senha, emissão de access token (15min) e refresh token (7 dias)
- **Proteção de Rotas:** Middleware verificando JWT em todas as rotas `/api/v1/` (exceto `/auth/login` e `/health`)
- **Perfil do Usuário:** Visualização e edição de nome, email e configurações gerais
- **Seed Automático:** Criação do usuário `admin` na primeira inicialização do sistema

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

> 💡 **Dica:** Para esta phase, o modo recomendado é **phase completa**, pois autenticação e gestão de usuários formam um bloco coeso e interdependente que é melhor revisado em conjunto.

---

## Tasks

### TASK-02-001 — Modelagem da Entidade User no Domínio

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Definir a entidade User com todos os campos necessários para o perfil e autenticação |

**Descrição:**
Definir a struct `User` no pacote `src/backend/domain/user/` com os campos: `ID` (ObjectID do MongoDB), `Name`, `Email` (único, com índice), `PasswordHash` (bcrypt), `Role` (enum: `ADMIN`, `USER`), `Prompts` (mapa de prompts por grupo de agente — vazio inicialmente, usado nas fases de gestão de prompts), `Enabled` (boolean), `CreatedAt`, `UpdatedAt`. Criar interfaces `UserRepository` e `UserService` com os métodos necessários. Implementar validações de domínio na criação do usuário (email válido, senha com mínimo de 8 caracteres, nome obrigatório).

**Critério de aceite:** Struct definida com validações; interfaces criadas; testes unitários de validação passando.

---

### TASK-02-002 — Repositório de Usuários com MongoDB

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar a camada de persistência de usuários no MongoDB |

**Descrição:**
Implementar `UserMongoRepository` em `src/backend/infra/database/mongodb/user_repository.go` satisfazendo a interface `UserRepository`. Implementar os métodos: `Create`, `FindByID`, `FindByEmail`, `Update`, `Delete`, `List`. Criar índices no MongoDB: índice único em `email`, índice em `role` para filtros futuros. Garantir que `PasswordHash` nunca seja retornado em queries de listagem (projeção explícita). Implementar soft delete (campo `DeletedAt`) para manter histórico de usuários.

**Critério de aceite:** CRUD funcional no MongoDB; índice único de email prevenindo duplicatas; projeção de senha funcional.

---

### TASK-02-003 — Implementação do Serviço de Autenticação JWT

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar geração, validação e refresh de tokens JWT |

**Descrição:**
Implementar o pacote `src/backend/pkg/auth/` com: geração de access token JWT com claims (`user_id`, `email`, `role`, `exp`), geração de refresh token (UUID armazenado no Redis com TTL de 7 dias apontando para o user_id), validação de access token (verificação de assinatura e expiração), endpoint de refresh que troca o refresh token por um novo par de tokens (rotação de refresh tokens). Usar `github.com/golang-jwt/jwt/v5` com algoritmo RS256 (par de chaves RSA). O secret/chave privada deve vir da configuração.

**Critério de aceite:** Tokens gerados com claims corretos; validação rejeita tokens expirados/inválidos; refresh token rotation funcional.

---

### TASK-02-004 — Handler de Login e Autenticação

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar os endpoints de autenticação da API |

**Descrição:**
Implementar os handlers em `src/backend/api/handler/auth_handler.go`: `POST /api/v1/auth/login` (recebe email/senha, valida, retorna access + refresh token), `POST /api/v1/auth/refresh` (recebe refresh token, valida no Redis, retorna novo par), `POST /api/v1/auth/logout` (invalida o refresh token no Redis), `GET /api/v1/auth/me` (retorna dados do usuário autenticado via JWT). Validar inputs com `github.com/go-playground/validator/v10`. Rate limiting no endpoint de login (máximo 5 tentativas por IP em 1 minuto) para prevenir brute force.

**Critério de aceite:** Login retorna tokens; logout invalida refresh token; rate limiting no login funcional.

---

### TASK-02-005 — Middleware de Autenticação JWT

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Proteger todas as rotas privadas da API com verificação de JWT |

**Descrição:**
Implementar o middleware GIN `AuthMiddleware` que extrai o JWT do header `Authorization: Bearer <token>`, valida assinatura e expiração, injeta o `UserContext` (user_id, email, role) no contexto GIN para que handlers downstream possam acessar sem re-parsear o token. Implementar `RoleMiddleware` para proteção de rotas que requerem papel específico (ex: apenas ADMIN). O middleware deve retornar `401 Unauthorized` para token ausente/inválido e `403 Forbidden` para papel insuficiente.

**Critério de aceite:** Rotas protegidas rejeitam requests sem JWT válido; `UserContext` disponível em handlers; RBAC funcional.

---

### TASK-02-006 — Seed de Dados: Usuário Admin Padrão

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir que o sistema sempre tenha um usuário admin operacional na primeira inicialização |

**Descrição:**
Implementar o pacote `src/backend/infra/seed/` com um `Seeder` que é executado na inicialização da aplicação (antes do servidor HTTP subir). O seeder deve verificar se o usuário `admin` já existe no banco — se sim, não faz nada (idempotente). Se não existe, cria com: email `admin@agency.ai`, senha gerada aleatoriamente (impressa no log de startup com aviso de segurança), nome `Administrator`, role `ADMIN`. Implementar flag de configuração `FORCE_ADMIN_RESET` que recria o admin se necessário.

**Critério de aceite:** Seed idempotente; admin criado apenas uma vez; senha logada no startup para configuração inicial.

---

### TASK-02-007 — Handler de Gestão de Perfil do Usuário

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Implementar endpoints para visualização e edição do perfil do usuário autenticado |

**Descrição:**
Implementar os handlers em `src/backend/api/handler/user_handler.go`: `GET /api/v1/users/me` (retorna perfil completo do usuário autenticado, exceto PasswordHash), `PUT /api/v1/users/me` (atualiza nome e configurações gerais), `PUT /api/v1/users/me/password` (troca de senha exigindo senha atual), `GET /api/v1/users` (listagem de usuários — apenas ADMIN). O usuário só pode editar o próprio perfil; o ADMIN pode editar qualquer perfil.

**Critério de aceite:** Perfil retornado sem PasswordHash; autorização respeitada; validação de senha atual na troca.

---

### TASK-02-008 — Tela de Login no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de login com design premium e integração com a API de autenticação |

**Descrição:**
Implementar a página `/login` no frontend com: formulário de email e senha com validação client-side, feedback visual claro de erros de autenticação (credenciais inválidas, conta desabilitada), loading state durante a chamada de API, redirecionamento automático para o dashboard após login bem-sucedido, persistência do access token no `localStorage` e refresh token em cookie `HttpOnly` (via backend). Layout com identidade visual premium da plataforma (dark mode, tipografia moderna, animações suaves de entrada).

**Critério de aceite:** Login funcional; tokens armazenados corretamente; UX de erro clara; redirecionamento após login.

---

### TASK-02-009 — Configuração do Cliente HTTP (Axios) no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Configurar o cliente HTTP com interceptors para autenticação automática e refresh de token |

**Descrição:**
Criar o módulo `src/frontend/src/services/api.ts` com instância do axios configurada com: `baseURL` apontando para o backend, interceptor de request que injeta o `Authorization: Bearer <token>` em todos os requests, interceptor de response que detecta `401 Unauthorized` e tenta refrescar o token automaticamente antes de retornar o erro para o componente. Implementar fila de requests pendentes durante o refresh para evitar múltiplas chamadas de refresh simultâneas. Exportar funções de serviço tipadas para cada domínio da API.

**Critério de aceite:** Token injetado automaticamente; refresh automático funcional; fila de requests durante refresh sem duplicatas.

---

### TASK-02-010 — Tela de Perfil do Usuário no Frontend

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Criar a tela de perfil para que o usuário visualize e edite suas informações |

**Descrição:**
Implementar a página `/profile` no frontend com: visualização dos dados do usuário (nome, email, role, data de criação), formulário de edição de nome, formulário de troca de senha (com campos: senha atual, nova senha, confirmação), feedback visual de sucesso/erro para cada operação, seção de resumo de uso (projetos criados, fases executadas — dados placeholders para ser preenchidos em fases futuras). Design consistente com o sistema de design da plataforma.

**Critério de aceite:** Dados do usuário exibidos corretamente; edição de nome funcional; troca de senha funcional com validação.

---

### TASK-02-011 — Proteção de Rotas no Frontend (Route Guards)

| Campo | Valor |
|-------|-------|
| **Camada** | Frontend |
| **Objetivo** | Garantir que rotas privadas exijam autenticação e redirecionar usuários não autenticados |

**Descrição:**
Implementar o componente `PrivateRoute` que verifica a presença e validade do token antes de renderizar a rota. Se o token estiver ausente ou expirado (e o refresh falhar), redirecionar para `/login` mantendo a URL de destino como parâmetro de query (`?redirect=/dashboard`) para retorno após autenticação. Implementar também verificação de role para rotas que exigem perfil específico (ex: rotas de admin).

**Critério de aceite:** Rotas privadas inacessíveis sem autenticação; redirecionamento com URL de retorno funcional; RBAC no frontend.

---

### TASK-02-012 — Testes Unitários da Camada de Autenticação

| Campo | Valor |
|-------|-------|
| **Camada** | Backend |
| **Objetivo** | Garantir cobertura de testes na lógica crítica de autenticação |

**Descrição:**
Implementar testes unitários utilizando `github.com/stretchr/testify` para: validação de geração e parsing de JWT, lógica de refresh token (incluindo token expirado, token não encontrado no Redis, token já revogado), validações de domínio da entidade User (email inválido, senha fraca), use cases de login (credenciais inválidas, usuário desabilitado, sucesso). Utilizar mocks para Redis e MongoDB com `github.com/stretchr/testify/mock`. Cobertura mínima de 80% nos pacotes de autenticação.

**Critério de aceite:** Testes passando com cobertura ≥ 80%; casos de erro cobertos; mocks sem dependências externas.

---

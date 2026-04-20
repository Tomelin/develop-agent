package project

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	domain "github.com/develop-agent/backend/internal/domain/project"
	"go.mongodb.org/mongo-driver/v2/bson"
	"gopkg.in/yaml.v3"
)

type Phase13Service struct {
	projects domain.ProjectRepository
	files    domain.CodeFileRepository
}

func NewPhase13Service(projects domain.ProjectRepository, files domain.CodeFileRepository) *Phase13Service {
	return &Phase13Service{projects: projects, files: files}
}

func (s *Phase13Service) Run(ctx context.Context, projectID, ownerID string, in domain.Phase13RunInput) (*domain.Phase13DeliveryReport, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, fmt.Errorf("project not found")
	}

	if !in.IncludeDevOps {
		// Fase 9 pode ser opcional; quando não usada, concluímos no final da documentação.
		in.IncludeDevOps = false
	}

	pid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, err
	}

	artifacts := make([]string, 0, 24)
	warnings := make([]string, 0)
	backendURL := strings.TrimSpace(in.BackendBaseURL)
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}
	frontendURL := strings.TrimSpace(in.FrontendURL)
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	openAPI, apiReference := buildOpenAPIAndReference(backendURL)
	if err := validateOpenAPI(openAPI); err != nil {
		warnings = append(warnings, "OpenAPI validation fallback used: "+err.Error())
	}

	files := []domain.CodeFile{
		mkFile(pid, "docs/prompts/TECH_WRITER.md", "TASK-13-001", "markdown", 8, phase13TechWriterPrompt()),
		mkFile(pid, "README.md", "TASK-13-002", "markdown", 8, buildReadme(frontendURL)),
		mkFile(pid, "docs/api/openapi.yaml", "TASK-13-003", "yaml", 8, openAPI),
		mkFile(pid, "docs/API_REFERENCE.md", "TASK-13-003", "markdown", 8, apiReference),
		mkFile(pid, "docs/OPERATIONS.md", "TASK-13-002", "markdown", 8, operationsManual()),
		mkFile(pid, "CONTRIBUTING.md", "TASK-13-002", "markdown", 8, contributingGuide()),
		mkFile(pid, "docs/prompts/DEVOPS.md", "TASK-13-004", "markdown", 9, phase13DevOpsPrompt()),
	}

	if in.IncludeDevOps {
		files = append(files,
			mkFile(pid, "infra/Dockerfile.backend", "TASK-13-005", "dockerfile", 9, dockerfileBackend()),
			mkFile(pid, "infra/Dockerfile.frontend", "TASK-13-005", "dockerfile", 9, dockerfileFrontend()),
			mkFile(pid, "infra/.dockerignore", "TASK-13-005", "plaintext", 9, dockerIgnore()),
			mkFile(pid, "infra/docker-compose.yml", "TASK-13-006", "yaml", 9, dockerCompose()),
			mkFile(pid, ".github/workflows/ci.yml", "TASK-13-007", "yaml", 9, githubActionsCI()),
			mkFile(pid, ".github/workflows/cd.yml", "TASK-13-007", "yaml", 9, githubActionsCD()),
			mkFile(pid, "k8s/backend-deployment.yaml", "TASK-13-008", "yaml", 9, k8sBackendDeployment()),
			mkFile(pid, "k8s/frontend-deployment.yaml", "TASK-13-008", "yaml", 9, k8sFrontendDeployment()),
			mkFile(pid, "k8s/services.yaml", "TASK-13-008", "yaml", 9, k8sServices()),
			mkFile(pid, "k8s/hpa.yaml", "TASK-13-008", "yaml", 9, k8sHPA()),
			mkFile(pid, "k8s/pdb.yaml", "TASK-13-008", "yaml", 9, k8sPDB()),
			mkFile(pid, "k8s/configmap.yaml", "TASK-13-008", "yaml", 9, k8sConfigMap()),
			mkFile(pid, "k8s/secret.template.yaml", "TASK-13-008", "yaml", 9, k8sSecretTemplate()),
			mkFile(pid, "k8s/ingress.yaml", "TASK-13-008", "yaml", 9, k8sIngress()),
		)
	}

	for i := range files {
		if err := s.files.Upsert(ctx, &files[i]); err != nil {
			return nil, err
		}
		artifacts = append(artifacts, files[i].Path)
	}

	summary := buildProjectSummary(p, artifacts)
	if err := s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   pid,
		Path:        "PROJECT_SUMMARY.md",
		TaskID:      "TASK-13-010",
		Language:    "markdown",
		PhaseNumber: 9,
		Content:     summary,
	}); err != nil {
		return nil, err
	}
	artifacts = append(artifacts, "PROJECT_SUMMARY.md")
	sort.Strings(artifacts)

	p.Status = domain.ProjectCompleted
	p.UpdatedAt = time.Now().UTC()
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	return &domain.Phase13DeliveryReport{
		GeneratedAt:   time.Now().UTC(),
		ProjectID:     projectID,
		ProjectStatus: string(p.Status),
		IncludeDevOps: in.IncludeDevOps,
		Artifacts:     artifacts,
		Warnings:      warnings,
	}, nil
}

func mkFile(projectID bson.ObjectID, path, task, lang string, phase int, content string) domain.CodeFile {
	return domain.CodeFile{ProjectID: projectID, Path: filepath.ToSlash(path), TaskID: task, Language: lang, PhaseNumber: phase, Content: content}
}

func validateOpenAPI(content string) error {
	var payload map[string]any
	if err := yaml.Unmarshal([]byte(content), &payload); err != nil {
		return err
	}
	if _, ok := payload["openapi"]; !ok {
		return fmt.Errorf("missing openapi field")
	}
	if _, ok := payload["paths"]; !ok {
		return fmt.Errorf("missing paths field")
	}
	return nil
}

func buildOpenAPIAndReference(baseURL string) (string, string) {
	op := fmt.Sprintf(`openapi: 3.0.3
info:
  title: Agência de IA API
  version: 1.0.0
  description: API principal do backend
servers:
  - url: %s/api/v1
security:
  - bearerAuth: []
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  schemas:
    ErrorResponse:
      type: object
      properties:
        error:
          type: string
paths:
  /projects:
    get:
      summary: Lista projetos
      responses:
        '200':
          description: OK
    post:
      summary: Cria projeto
      responses:
        '201':
          description: Criado
  /projects/{id}:
    get:
      summary: Busca projeto por ID
      parameters:
        - in: path
          name: id
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
        '404':
          description: Não encontrado
  /projects/{id}/tasks:
    get:
      summary: Lista tasks do projeto
      parameters:
        - in: path
          name: id
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
  /projects/{id}/phases/6/analyze-coverage:
    post:
      summary: Analisa cobertura de testes
      responses:
        '200':
          description: OK
        '400':
          description: Erro de validação
  /projects/{id}/phases/7/run-audit:
    post:
      summary: Executa auditoria de segurança
      responses:
        '200':
          description: OK
  /projects/{id}/phases/13/run:
    post:
      summary: Gera artefatos de documentação e DevOps
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  project_status:
                    type: string
                  artifacts:
                    type: array
                    items:
                      type: string
`, baseURL)

	ref := `# API_REFERENCE

## Autenticação

| Mecanismo | Valor |
|---|---|
| Tipo | Bearer JWT |
| Header | Authorization: Bearer <token> |

## Endpoints

### GET /api/v1/projects
- Lista projetos do usuário autenticado.
- Resposta: **200 OK**.

### POST /api/v1/projects
- Cria um novo projeto.
- Resposta: **201 Created**.

### GET /api/v1/projects/{id}
- Busca detalhes de um projeto específico.
- Path params: ` + "`id`" + ` (string)
- Respostas: **200**, **404**.

### GET /api/v1/projects/{id}/tasks
- Lista tasks por projeto.
- Path params: ` + "`id`" + `.
- Resposta: **200**.

### POST /api/v1/projects/{id}/phases/6/analyze-coverage
- Executa análise de cobertura da Fase 6.
- Body exemplo:
` + "```json" + `
{"backend_dir":"src/backend","threshold":80}
` + "```" + `
- Curl:
` + "```bash" + `
curl -X POST "$API/projects/{id}/phases/6/analyze-coverage" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"backend_dir":"src/backend"}'
` + "```" + `

### POST /api/v1/projects/{id}/phases/7/run-audit
- Executa auditoria de segurança.
- Resposta: relatório consolidado com findings.

### POST /api/v1/projects/{id}/phases/13/run
- Gera artefatos da Fase 13 (Documentação + DevOps).
- Body exemplo:
` + "```json" + `
{"backend_base_url":"http://localhost:8080","frontend_url":"http://localhost:3000","include_devops":true}
` + "```" + `
`

	return op, ref
}

func buildReadme(frontendURL string) string {
	return fmt.Sprintf(`# 🚀 Agência de IA

[![CI](https://img.shields.io/badge/CI-GitHub_Actions-blue)](#)
[![Coverage](https://img.shields.io/badge/Coverage-Phase6-green)](#)
[![License](https://img.shields.io/badge/License-MIT-yellow)](./LICENSE)

Plataforma de esteira inteligente para desenvolvimento orientado por Tríade de Agentes (Produtor, Revisor e Refinador).

## ✨ Features
- 🤖 Tríade de agentes com revisão obrigatória
- 🧭 Execução por phase completa ou task individual
- 🔐 Segurança e testes com gatilhos de auto-rejeição
- 📚 Geração de documentação e artefatos DevOps

## ⚡ Quickstart (<= 5 comandos)
`+"```bash"+`
cp .env.example .env
docker compose -f infra/docker-compose.yml up -d
# aguarde os healthchecks
open %s
`+"```"+`

## 🏗️ Arquitetura
`+"```mermaid"+`
flowchart LR
    UI[Frontend] --> API[Backend Gin]
    API --> DB[(MongoDB)]
    API --> MQ[(RabbitMQ)]
    API --> CACHE[(Redis)]
`+"```"+`

## 📖 Documentação
- [API Reference](docs/API_REFERENCE.md)
- [OpenAPI](docs/api/openapi.yaml)
- [Operations](docs/OPERATIONS.md)
- [Contributing](CONTRIBUTING.md)

## 🤝 Contribuindo
Veja [CONTRIBUTING.md](CONTRIBUTING.md).

## 📄 Licença
MIT.
`, frontendURL)
}

func operationsManual() string {
	return `# OPERATIONS

## Deploy
- Build de imagens via GitHub Actions CD.
- Aplicar manifests com kubectl apply -f k8s/.

## Escalabilidade
- HPA configurado por CPU e memória.
- PDB para disponibilidade durante manutenção.

## Monitoramento
- Expor /health no backend e frontend.
- Integrar logs de aplicação com stack centralizada.

## Troubleshooting
1. Verificar pods (` + "`kubectl get pods -n agency`" + `).
2. Verificar logs (` + "`kubectl logs <pod>`" + `).
3. Validar secrets e configmaps aplicados.
`
}

func contributingGuide() string {
	return `# CONTRIBUTING

## Setup local
1. Instale Docker e Docker Compose.
2. Configure arquivo .env a partir de .env.example.
3. Suba ambiente com ` + "`docker compose -f infra/docker-compose.yml up`" + `.

## Workflow de PR
- Crie branch por feature: ` + "`feat/<tema>`" + `.
- Abra PR com contexto, riscos e checklist de testes.

## Style guide
- Go: gofmt obrigatório e testes com ` + "`go test ./...`" + `.
- Nomes explícitos e funções pequenas.

## Processo de review
- Revisão obrigatória antes de merge.
- Não mergear com CI em falha.
`
}

func phase13TechWriterPrompt() string {
	return `# Prompt — Tech Writer (TASK-13-001)

- Não invente funcionalidades: documente apenas o que existe no código/arquitetura.
- README deve conter hero, badges, pré-requisitos, quickstart (máximo 5 comandos), features, diagrama Mermaid, contribuição e licença.
- API Reference: um endpoint por seção, com método, URL, autenticação, parâmetros, status codes, schema e exemplo curl.
- OPERATIONS.md: deploy, escalabilidade, monitoramento e troubleshooting.
- CONTRIBUTING.md: setup local, workflow de PR, style guide e processo de review.
- Revisor deve validar completude, precisão e quickstart funcional com docker-compose.
`
}

func phase13DevOpsPrompt() string {
	return `# Prompt — DevOps (TASK-13-004)

- Inferir portas, variáveis de ambiente e dependências do código.
- Gerar Dockerfiles multi-stage, docker-compose dev, workflows CI/CD e manifests Kubernetes.
- Otimizar imagens para produção e executar como non-root.
- Pipeline deve conter lint, test, build, push e deploy.
- Kubernetes com securityContext, liveness/readiness, HPA e PDB.
`
}

func dockerfileBackend() string {
	return `FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY src/backend/go.mod src/backend/go.sum ./
RUN go mod download
COPY src/backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
USER nonroot:nonroot
WORKDIR /app
COPY --from=builder /out/api /app/api
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD ["/app/api", "health"]
ENTRYPOINT ["/app/api"]
`
}

func dockerfileFrontend() string {
	return `FROM node:22-alpine AS builder
WORKDIR /app
COPY src/frontend/package*.json ./
RUN npm ci
COPY src/frontend/ .
RUN npm run build

FROM nginx:1.27-alpine
RUN adduser -D -H -u 10001 app && chown -R app:app /var/cache/nginx /var/run /var/log/nginx /etc/nginx/conf.d
USER app
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost/ || exit 1
`
}

func dockerIgnore() string {
	return ".git\nnode_modules\ncoverage.out\n*.log\n.tmp\n"
}

func dockerCompose() string {
	return `version: "3.9"
services:
  mongo:
    image: mongo:8
    ports: ["27017:27017"]
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 3s
      retries: 10
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
  rabbitmq:
    image: rabbitmq:3-management
    ports: ["5672:5672", "15672:15672"]
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
  backend:
    build:
      context: .
      dockerfile: infra/Dockerfile.backend
    depends_on:
      mongo: { condition: service_healthy }
      redis: { condition: service_healthy }
      rabbitmq: { condition: service_healthy }
    environment:
      APP_ENV: development
    ports: ["8080:8080"]
`
}

func githubActionsCI() string {
	return `name: CI
on: [push, pull_request]
jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.25' }
      - run: go test ./...
        working-directory: src/backend
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22', cache: 'npm', cache-dependency-path: 'src/frontend/package-lock.json' }
      - run: npm ci
        working-directory: src/frontend
      - run: npm run lint
        working-directory: src/frontend
`
}

func githubActionsCD() string {
	return `name: CD
on:
  push:
    branches: [main]
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - run: echo "build and push images"
  deploy-staging:
    runs-on: ubuntu-latest
    needs: [build-and-push]
    steps:
      - run: echo "deploy staging"
  deploy-production:
    runs-on: ubuntu-latest
    needs: [deploy-staging]
    environment: production
    steps:
      - run: echo "deploy production with manual approval"
`
}

func k8sBackendDeployment() string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: agency-backend
spec:
  replicas: 2
  selector: { matchLabels: { app: agency-backend } }
  template:
    metadata: { labels: { app: agency-backend } }
    spec:
      containers:
        - name: backend
          image: ghcr.io/org/agency-backend:latest
          ports: [{containerPort: 8080}]
          resources:
            requests: { cpu: "100m", memory: "128Mi" }
            limits: { cpu: "500m", memory: "512Mi" }
          livenessProbe: { httpGet: { path: /health, port: 8080 } }
          readinessProbe: { httpGet: { path: /health, port: 8080 } }
          securityContext:
            runAsNonRoot: true
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
`
}

func k8sFrontendDeployment() string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: agency-frontend
spec:
  replicas: 2
  selector: { matchLabels: { app: agency-frontend } }
  template:
    metadata: { labels: { app: agency-frontend } }
    spec:
      containers:
        - name: frontend
          image: ghcr.io/org/agency-frontend:latest
          ports: [{containerPort: 80}]
          livenessProbe: { httpGet: { path: /health, port: 80 } }
          readinessProbe: { httpGet: { path: /health, port: 80 } }
          securityContext:
            runAsNonRoot: true
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
`
}

func k8sServices() string {
	return `apiVersion: v1
kind: Service
metadata:
  name: agency-backend
spec:
  selector: { app: agency-backend }
  ports: [{ port: 8080, targetPort: 8080 }]
---
apiVersion: v1
kind: Service
metadata:
  name: agency-frontend
spec:
  type: LoadBalancer
  selector: { app: agency-frontend }
  ports: [{ port: 80, targetPort: 80 }]
`
}

func k8sHPA() string {
	return `apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: agency-backend
spec:
  scaleTargetRef: { apiVersion: apps/v1, kind: Deployment, name: agency-backend }
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource: { name: cpu, target: { type: Utilization, averageUtilization: 70 } }
    - type: Resource
      resource: { name: memory, target: { type: Utilization, averageUtilization: 75 } }
`
}

func k8sPDB() string {
	return `apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: agency-backend-pdb
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: agency-backend
`
}

func k8sConfigMap() string {
	return `apiVersion: v1
kind: ConfigMap
metadata:
  name: agency-config
data:
  APP_ENV: production
  LOG_LEVEL: info
`
}

func k8sSecretTemplate() string {
	return `apiVersion: v1
kind: Secret
metadata:
  name: agency-secrets
type: Opaque
stringData:
  JWT_PRIVATE_KEY_B64: "<replace-me>"
  MONGO_URI: "<replace-me>"
`
}

func k8sIngress() string {
	return `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agency-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
    - hosts: ["agency.example.com"]
      secretName: agency-tls
  rules:
    - host: agency.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: agency-frontend
                port:
                  number: 80
`
}

func buildProjectSummary(p *domain.Project, artifacts []string) string {
	completed := 0
	for _, ph := range p.Phases {
		if ph.Status == domain.PhaseCompleted {
			completed++
		}
	}
	return fmt.Sprintf(`# PROJECT_SUMMARY

- Projeto: **%s**
- Status final: **COMPLETED**
- Fases completas: **%d/%d**
- Total de artefatos desta entrega: **%d**
- Tokens consumidos: **%d**
- Custo total estimado (USD): **%.4f**

## Artefatos gerados
%s
`, p.Name, completed, len(p.Phases), len(artifacts), p.TotalTokensUsed, p.TotalCostUSD, bullets(artifacts))
}

func bullets(items []string) string {
	if len(items) == 0 {
		return "- (nenhum)"
	}
	sort.Strings(items)
	b := strings.Builder{}
	for _, item := range items {
		b.WriteString("- `" + item + "`\n")
	}
	return b.String()
}

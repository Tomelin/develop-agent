package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"

	"github.com/develop-agent/backend/api/handler"
	"github.com/develop-agent/backend/api/health"
	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/api/server"
	"github.com/develop-agent/backend/config"
	"github.com/develop-agent/backend/internal/domain/agent"
	"github.com/develop-agent/backend/internal/domain/billing"
	"github.com/develop-agent/backend/internal/domain/interview"
	"github.com/develop-agent/backend/internal/domain/organization"
	"github.com/develop-agent/backend/internal/domain/project"
	"github.com/develop-agent/backend/internal/domain/prompt"
	"github.com/develop-agent/backend/internal/domain/user"
	"github.com/develop-agent/backend/internal/infra/cache/redis"
	"github.com/develop-agent/backend/internal/infra/database/mongodb"
	"github.com/develop-agent/backend/internal/infra/messaging/rabbitmq"
	"github.com/develop-agent/backend/internal/infra/seed"
	usecaseauth "github.com/develop-agent/backend/internal/usecase/auth"
	usecasebilling "github.com/develop-agent/backend/internal/usecase/billing"
	usecaseinterview "github.com/develop-agent/backend/internal/usecase/interview"
	usecaseorganization "github.com/develop-agent/backend/internal/usecase/organization"
	usecaseproject "github.com/develop-agent/backend/internal/usecase/project"
	usecaseprompt "github.com/develop-agent/backend/internal/usecase/prompt"
	usecasetriad "github.com/develop-agent/backend/internal/usecase/triad"
	"github.com/develop-agent/backend/pkg/agentsdk"
	"github.com/develop-agent/backend/pkg/agentsdk/anthropic"
	"github.com/develop-agent/backend/pkg/agentsdk/gemini"
	"github.com/develop-agent/backend/pkg/agentsdk/mock"
	"github.com/develop-agent/backend/pkg/agentsdk/ollama"
	"github.com/develop-agent/backend/pkg/agentsdk/openai"
	pkgauth "github.com/develop-agent/backend/pkg/auth"
	"github.com/develop-agent/backend/pkg/logger"
	"github.com/develop-agent/backend/pkg/observability"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if logErr := logger.Setup(cfg.App.Env); logErr != nil {
		log.Fatalf("Failed to setup logger: %v", logErr)
	}
	defer func() { _ = logger.Global().Sync() }()

	mongoClient, err := mongodb.NewAdapter(cfg.Mongo.URI)
	if err != nil {
		logger.Global().Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() { _ = mongoClient.Close(context.Background()) }()

	redisClient, err := redis.NewAdapter(cfg.Redis.Addr, cfg.Redis.Password)
	if err != nil {
		logger.Global().Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer func() { _ = redisClient.Close() }()

	rmqClient, err := rabbitmq.NewAdapter(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Global().Error("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		defer func() { _ = rmqClient.Close() }()
	}

	userRepo := mongodb.NewUserRepository(mongoClient, cfg.Mongo.DBName)
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Mongo indexes", zap.Error(err))
	}
	agentRepo := mongodb.NewAgentRepository(mongoClient, cfg.Mongo.DBName)
	if err := agentRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Agent Mongo indexes", zap.Error(err))
	}
	projectRepo := mongodb.NewProjectRepository(mongoClient, cfg.Mongo.DBName)
	if err := projectRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Project Mongo indexes", zap.Error(err))
	}
	taskRepo := mongodb.NewTaskRepository(mongoClient, cfg.Mongo.DBName)
	if err := taskRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Task Mongo indexes", zap.Error(err))
	}
	codeFileRepo := mongodb.NewCodeFileRepository(mongoClient, cfg.Mongo.DBName)
	if err := codeFileRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure CodeFile Mongo indexes", zap.Error(err))
	}
	promptRepo := mongodb.NewUserPromptRepository(mongoClient, cfg.Mongo.DBName)
	if err := promptRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Prompt Mongo indexes", zap.Error(err))
	}
	interviewRepo := mongodb.NewInterviewRepository(mongoClient, cfg.Mongo.DBName)
	if err := interviewRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Interview Mongo indexes", zap.Error(err))
	}
	billingRepo := mongodb.NewBillingRepository(mongoClient, cfg.Mongo.DBName)
	if err := billingRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Billing Mongo indexes", zap.Error(err))
	}
	orgRepo := mongodb.NewOrganizationRepository(mongoClient, cfg.Mongo.DBName)
	if err := orgRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Organization Mongo indexes", zap.Error(err))
	}

	if err := seed.NewAdminSeeder(userRepo).Run(context.Background(), cfg.Seed.ForceAdminReset); err != nil {
		logger.Global().Fatal("Failed to seed admin user", zap.Error(err))
	}
	if err := seed.NewAgentsSeeder(agentRepo).Run(context.Background()); err != nil {
		logger.Global().Fatal("Failed to seed agents catalog", zap.Error(err))
	}

	tokenManager, err := pkgauth.NewTokenManager(
		cfg.Auth.JWTPrivateKeyB64,
		cfg.Auth.JWTIssuer,
		cfg.Auth.JWTAudience,
		cfg.Auth.AccessTTLMinutes,
		cfg.Auth.RefreshTTLDays,
	)
	if err != nil {
		logger.Global().Fatal("Failed to initialize token manager", zap.Error(err))
	}

	authService := usecaseauth.NewService(userRepo, tokenManager, pkgauth.NewRedisRefreshStore(redisClient))
	projectService := usecaseproject.NewService(projectRepo, taskRepo)
	agentSelector := agent.NewSelectorService(agentRepo, nil, true, nil)
	if rmqClient != nil {
		projectService = projectService.WithPublisher(rabbitmq.NewPublisher(rmqClient))
	}
	developmentService := usecaseproject.NewDevelopmentService(projectRepo, taskRepo, codeFileRepo)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	agentHandler := handler.NewAgentHandler(agentRepo)
	organizationHandler := handler.NewOrganizationHandler(usecaseorganization.NewService(orgRepo, userRepo))
	projectHandler := handler.NewProjectHandler(projectRepo, projectService)
	taskHandler := handler.NewTaskHandler(taskRepo, projectRepo)
	phase5Handler := handler.NewPhase5Handler(projectRepo, developmentService)
	phase6Service := usecaseproject.NewPhase6Service(codeFileRepo, developmentService.TriggerAutoRejection)
	phase6Handler := handler.NewPhase6Handler(projectRepo, phase6Service)
	phase7Service := usecaseproject.NewPhase7Service(codeFileRepo, developmentService.TriggerAutoRejection)
	phase7Handler := handler.NewPhase7Handler(projectRepo, phase7Service)
	phase13Service := usecaseproject.NewPhase13Service(projectRepo, codeFileRepo)
	phase13Handler := handler.NewPhase13Handler(projectRepo, phase13Service)
	phase14Service := usecaseproject.NewPhase14Service(projectRepo, codeFileRepo)
	phase14Handler := handler.NewPhase14Handler(projectRepo, phase14Service)
	phase15Service := usecaseproject.NewPhase15Service(projectRepo, codeFileRepo)
	phase15Handler := handler.NewPhase15Handler(projectRepo, phase15Service)
	promptHandler := handler.NewPromptHandler(promptRepo, usecaseprompt.NewService(promptRepo))
	interviewProvider := &agentsdk.MetricsProvider{Base: mock.New()}
	interviewService := usecaseinterview.NewService(interviewRepo, projectRepo, interviewProvider, nil)
	interviewHandler := handler.NewInterviewHandler(interviewService)
	billingService, err := usecasebilling.NewService(billingRepo, cfg.Billing.PricingFile)
	if err != nil {
		logger.Global().Fatal("Failed to load billing pricing table", zap.Error(err))
	}
	billingHandler := handler.NewBillingHandler(billingService)
	phase19Service := usecaseproject.NewPhase19Service(projectRepo, codeFileRepo)
	phase19Handler := handler.NewPhase19Handler(projectRepo, phase19Service)
	adminQualityHandler := handler.NewAdminQualityHandler(usecaseproject.NewAdminQualityReportService(projectRepo, codeFileRepo))
	integrationCompatHandler := handler.NewIntegrationCompatHandler(projectRepo, agentSelector)

	srv := server.New(cfg)
	v1 := srv.Router().Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
		authHandler.Register(v1)

		private := v1.Group("")
		private.Use(middleware.AuthMiddleware(authService))
		private.Use(middleware.OrganizationMiddleware())
		userHandler.Register(private)
		agentHandler.Register(private)
		organizationHandler.Register(private)
		projectHandler.Register(private)
		taskHandler.Register(private)
		phase5Handler.Register(private)
		phase6Handler.Register(private)
		phase7Handler.Register(private)
		phase13Handler.Register(private)
		phase14Handler.Register(private)
		phase15Handler.Register(private)
		phase19Handler.Register(private)
		promptHandler.Register(private)
		interviewHandler.Register(private)
		billingHandler.Register(private)
		integrationCompatHandler.Register(private)

		adminOnly := private.Group("")
		adminOnly.Use(middleware.RoleMiddleware("ADMIN"))
		adminQualityHandler.Register(adminOnly)
	}

	healthHandler := health.NewHandler(mongoClient, redisClient, rmqClient)
	healthHandler.Register(srv.Router())
	statusHandler := handler.NewStatusHandler(healthHandler)
	statusHandler.Register(v1)
	srv.Router().GET("/metrics", observability.Handler())

	if rmqClient != nil {
		bootstrapPhaseWorkers(context.Background(), rmqClient, agentSelector, projectService, codeFileRepo)
	}

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Global().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Global().Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logger.Global().Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Global().Info("Server exiting")
}

var _ user.Repository = (*mongodb.UserRepository)(nil)
var _ agent.Repository = (*mongodb.AgentRepository)(nil)
var _ project.ProjectRepository = (*mongodb.ProjectRepository)(nil)
var _ project.TaskRepository = (*mongodb.TaskRepository)(nil)
var _ project.CodeFileRepository = (*mongodb.CodeFileRepository)(nil)
var _ prompt.UserPromptRepository = (*mongodb.UserPromptRepository)(nil)
var _ interview.Repository = (*mongodb.InterviewRepository)(nil)
var _ billing.Repository = (*mongodb.BillingRepository)(nil)
var _ organization.Repository = (*mongodb.OrganizationRepository)(nil)

func bootstrapPhaseWorkers(ctx context.Context, rmqClient *rabbitmq.Adapter, selector *agent.Service, projectService *usecaseproject.Service, codeFiles project.CodeFileRepository) {
	for phase := 1; phase <= 9; phase++ {
		skill, err := skillForPhase(phase)
		if err != nil {
			logger.Global().Error("Failed to map phase skill", zap.Int("phase", phase), zap.Error(err))
			continue
		}
		triadSelection, err := selector.SelectTriad(ctx, skill)
		if err != nil {
			logger.Global().Error("Failed to select triad for phase worker", zap.Int("phase", phase), zap.Error(err))
			continue
		}
		producer, err := buildProviderForAgent(ctx, triadSelection.Producer)
		if err != nil {
			logger.Global().Error("Failed to build producer provider", zap.Int("phase", phase), zap.Error(err))
			continue
		}
		reviewer, err := buildProviderForAgent(ctx, triadSelection.Reviewer)
		if err != nil {
			logger.Global().Error("Failed to build reviewer provider", zap.Int("phase", phase), zap.Error(err))
			continue
		}
		refiner, err := buildProviderForAgent(ctx, triadSelection.Refiner)
		if err != nil {
			logger.Global().Error("Failed to build refiner provider", zap.Int("phase", phase), zap.Error(err))
			continue
		}

		worker := &usecasetriad.Worker{
			QueueName: fmt.Sprintf("phase.%d", phase),
			Consumer:  rabbitmq.NewConsumer(rmqClient),
			Orchestrator: &usecasetriad.Orchestrator{
				Producer: producer,
				Reviewer: reviewer,
				Refiner:  refiner,
				Events:   usecasetriad.NewBroker(),
			},
			OnSuccess: func(ctx context.Context, job usecasetriad.Job, refined string) error {
				pid, err := bson.ObjectIDFromHex(job.ProjectID)
				if err != nil {
					return err
				}
				if err := codeFiles.Upsert(ctx, &project.CodeFile{
					ProjectID:   pid,
					Path:        fmt.Sprintf("artifacts/phase_%d/output.md", job.PhaseNumber),
					Content:     refined,
					TaskID:      fmt.Sprintf("PHASE-%d-TRIAD", job.PhaseNumber),
					Language:    "markdown",
					Version:     time.Now().UTC(),
					PhaseNumber: job.PhaseNumber,
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
				}); err != nil {
					return err
				}
				track := project.Track(strings.ToUpper(strings.TrimSpace(job.Track)))
				if track == "" {
					track = project.TrackFull
				}
				_, err = projectService.ApprovePhaseTrack(ctx, job.ProjectID, job.OwnerUserID, job.PhaseNumber, track)
				return err
			},
		}
		go func(w *usecasetriad.Worker, phaseNumber int) {
			if err := w.Start(); err != nil {
				logger.Global().Error("Failed to start phase worker", zap.Int("phase", phaseNumber), zap.Error(err))
			}
		}(worker, phase)
	}
}

func skillForPhase(phase int) (agent.Skill, error) {
	switch phase {
	case 1:
		return agent.SkillProjectCreation, nil
	case 2:
		return agent.SkillEngineering, nil
	case 3:
		return agent.SkillArchitecture, nil
	case 4:
		return agent.SkillPlanning, nil
	case 5:
		return agent.SkillDevelopmentBackend, nil
	case 6:
		return agent.SkillTesting, nil
	case 7:
		return agent.SkillSecurity, nil
	case 8:
		return agent.SkillDocumentation, nil
	case 9:
		return agent.SkillDevOps, nil
	default:
		return "", fmt.Errorf("phase %d not supported", phase)
	}
}

func buildProviderForAgent(ctx context.Context, ag agent.Agent) (agentsdk.Provider, error) {
	var provider agentsdk.Provider
	switch ag.Provider {
	case agent.ProviderOpenAI:
		provider = openai.New()
	case agent.ProviderAnthropic:
		provider = anthropic.New()
	case agent.ProviderGoogle:
		provider = gemini.New()
	case agent.ProviderOllama:
		provider = ollama.New()
	default:
		return nil, fmt.Errorf("unsupported provider: %s", ag.Provider)
	}
	_ = provider.Initialize(ctx, agentsdk.Config{
		Token: os.Getenv(strings.TrimSpace(ag.ApiKeyRef)),
		Model: ag.Model,
	})
	return provider, nil
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/develop-agent/backend/api/handler"
	"github.com/develop-agent/backend/api/health"
	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/api/server"
	"github.com/develop-agent/backend/config"
	"github.com/develop-agent/backend/internal/domain/user"
	"github.com/develop-agent/backend/internal/infra/cache/redis"
	"github.com/develop-agent/backend/internal/infra/database/mongodb"
	"github.com/develop-agent/backend/internal/infra/messaging/rabbitmq"
	"github.com/develop-agent/backend/internal/infra/seed"
	usecaseauth "github.com/develop-agent/backend/internal/usecase/auth"
	pkgauth "github.com/develop-agent/backend/pkg/auth"
	"github.com/develop-agent/backend/pkg/logger"
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
	defer mongoClient.Close(context.Background())

	redisClient, err := redis.NewAdapter(cfg.Redis.Addr, cfg.Redis.Password)
	if err != nil {
		logger.Global().Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	rmqClient, err := rabbitmq.NewAdapter(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Global().Error("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		defer rmqClient.Close()
	}

	userRepo := mongodb.NewUserRepository(mongoClient, cfg.Mongo.DBName)
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Global().Fatal("Failed to ensure Mongo indexes", zap.Error(err))
	}

	if err := seed.NewAdminSeeder(userRepo).Run(context.Background(), cfg.Seed.ForceAdminReset); err != nil {
		logger.Global().Fatal("Failed to seed admin user", zap.Error(err))
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
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)

	srv := server.New(cfg)
	v1 := srv.Router().Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
		authHandler.Register(v1)

		private := v1.Group("")
		private.Use(middleware.AuthMiddleware(authService))
		userHandler.Register(private)
	}

	health.NewHandler(mongoClient, redisClient, rmqClient).Register(srv.Router())

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

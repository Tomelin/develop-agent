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

	"github.com/develop-agent/backend/api/health"
	"github.com/develop-agent/backend/api/server"
	"github.com/develop-agent/backend/config"
	"github.com/develop-agent/backend/internal/infra/cache/redis"
	"github.com/develop-agent/backend/internal/infra/database/mongodb"
	"github.com/develop-agent/backend/internal/infra/messaging/rabbitmq"
	"github.com/develop-agent/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init logger
	if logErr := logger.Setup(cfg.App.Env); logErr != nil {
		log.Fatalf("Failed to setup logger: %v", logErr)
	}
	defer func() {
		_ = logger.Global().Sync()
	}()

	// Init adapters
	mongoClient, err := mongodb.NewAdapter(cfg.Mongo.URI)
	if err != nil {
		logger.Global().Error("Failed to connect to MongoDB", zap.Error(err))
	} else {
		defer mongoClient.Close(context.Background())
		logger.Global().Info("Connected to MongoDB")
	}

	redisClient, err := redis.NewAdapter(cfg.Redis.Addr, cfg.Redis.Password)
	if err != nil {
		logger.Global().Error("Failed to connect to Redis", zap.Error(err))
	} else {
		defer redisClient.Close()
		logger.Global().Info("Connected to Redis")
	}

	rmqClient, err := rabbitmq.NewAdapter(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Global().Error("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		defer rmqClient.Close()
		logger.Global().Info("Connected to RabbitMQ")
	}

	// Init server
	srv := server.New(cfg)

	// Register routes
	v1 := srv.Router().Group("/api/v1")
	{
		// Demo route
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	// Register Health
	healthHandler := health.NewHandler(mongoClient, redisClient, rmqClient)
	healthHandler.Register(srv.Router())

	// Graceful shutdown
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

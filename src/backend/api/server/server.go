package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/develop-agent/backend/config"
	"github.com/develop-agent/backend/pkg/logger"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
	config *config.Config
}

func New(cfg *config.Config) *Server {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middlewares
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())

	return &Server{
		router: r,
		config: cfg,
	}
}

// Router returns the inner gin router to register routes.
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:              ":" + s.config.App.Port,
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Global().Info("Server starting", zap.String("port", s.config.App.Port))
	return s.srv.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	logger.Global().Info("Shutting down server gracefully...")
	return s.srv.Shutdown(ctx)
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID, ok := c.Get("request_id")
		var rID string
		if ok {
			if id, valid := reqID.(string); valid {
				rID = id
			}
		}

		// Assuming we want to export context key from logger pkg eventually, using string for now.
		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, rID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.Global().Info("HTTP request",
			zap.String("request_id", rID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Modify based on config if needed
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

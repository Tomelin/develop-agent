package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var globalLogger *zap.Logger

// Setup initializes the global logger based on environment.
func Setup(env string) error {
	var err error
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	globalLogger, err = cfg.Build()
	if err != nil {
		return err
	}

	// Make sure the global logger is replaced
	zap.ReplaceGlobals(globalLogger)
	return nil
}

// WithContext returns a logger with context values injected.
func WithContext(ctx context.Context) *zap.Logger {
	if globalLogger == nil {
		return zap.NewNop()
	}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return globalLogger.With(zap.String("request_id", reqID))
	}

	return globalLogger
}

// Global returns the global logger.
func Global() *zap.Logger {
	if globalLogger == nil {
		return zap.NewNop()
	}
	return globalLogger
}

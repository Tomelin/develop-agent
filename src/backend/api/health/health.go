package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/internal/infra/cache/redis"
	"github.com/develop-agent/backend/internal/infra/database/mongodb"
	"github.com/develop-agent/backend/internal/infra/messaging/rabbitmq"
)

type ComponentStatus string

const (
	Healthy   ComponentStatus = "healthy"
	Degraded  ComponentStatus = "degraded"
	Unhealthy ComponentStatus = "unhealthy"
)

type CheckResult struct {
	Status  ComponentStatus `json:"status"`
	Latency int64           `json:"latency_ms,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type Response struct {
	Status     ComponentStatus        `json:"status"`
	Components map[string]CheckResult `json:"components,omitempty"`
}

type Handler struct {
	mongoClient *mongodb.Adapter
	redisClient *redis.Adapter
	rmqClient   *rabbitmq.Adapter
}

func NewHandler(mongo *mongodb.Adapter, redis *redis.Adapter, rmq *rabbitmq.Adapter) *Handler {
	return &Handler{
		mongoClient: mongo,
		redisClient: redis,
		rmqClient:   rmq,
	}
}

// Register registers health routes.
func (h *Handler) Register(router *gin.Engine) {
	router.GET("/health", h.CheckAll)
	router.GET("/health/live", h.Live)
	router.GET("/health/ready", h.Ready)
}

func (h *Handler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (h *Handler) Ready(c *gin.Context) {
	resp := h.performChecks(c.Request.Context())
	if resp.Status != Healthy {
		c.JSON(http.StatusServiceUnavailable, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CheckAll(c *gin.Context) {
	resp := h.performChecks(c.Request.Context())
	if resp.Status != Healthy {
		c.JSON(http.StatusServiceUnavailable, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) performChecks(ctx context.Context) Response {
	components := make(map[string]CheckResult)
	overallStatus := Healthy

	// Mongo Check
	if h.mongoClient != nil {
		lat, err := h.mongoClient.Ping(ctx)
		components["mongodb"] = formatResult(lat, err)
		if err != nil {
			overallStatus = Unhealthy
		}
	} else {
		components["mongodb"] = CheckResult{Status: Unhealthy, Error: "not configured"}
		overallStatus = Unhealthy
	}

	// Redis Check
	if h.redisClient != nil {
		lat, err := h.redisClient.Ping(ctx)
		components["redis"] = formatResult(lat, err)
		if err != nil {
			overallStatus = Unhealthy
		}
	} else {
		components["redis"] = CheckResult{Status: Unhealthy, Error: "not configured"}
		overallStatus = Unhealthy
	}

	// RabbitMQ Check
	if h.rmqClient != nil {
		lat, err := h.rmqClient.Ping()
		components["rabbitmq"] = formatResult(lat, err)
		if err != nil {
			overallStatus = Unhealthy
		}
	} else {
		components["rabbitmq"] = CheckResult{Status: Unhealthy, Error: "not configured"}
		overallStatus = Unhealthy
	}

	return Response{
		Status:     overallStatus,
		Components: components,
	}
}

func formatResult(latency int64, err error) CheckResult {
	if err != nil {
		return CheckResult{
			Status:  Unhealthy,
			Latency: latency,
			Error:   err.Error(),
		}
	}
	status := Healthy
	if latency > 1000 {
		status = Degraded
	}
	return CheckResult{
		Status:  status,
		Latency: latency,
	}
}

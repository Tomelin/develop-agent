package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/api/health"
)

type statusLevel string

const (
	statusOperational statusLevel = "OPERATIONAL"
	statusDegraded    statusLevel = "DEGRADED"
	statusOutage      statusLevel = "OUTAGE"
)

type incident struct {
	Title      string    `json:"title"`
	OccurredAt time.Time `json:"occurred_at"`
	Duration   string    `json:"duration"`
	Resolution string    `json:"resolution"`
}

type StatusHandler struct {
	health *health.Handler
}

func NewStatusHandler(healthHandler *health.Handler) *StatusHandler {
	return &StatusHandler{health: healthHandler}
}

func (h *StatusHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/status", h.Get)
}

func (h *StatusHandler) Get(c *gin.Context) {
	snap := h.health.Snapshot(c.Request.Context())

	providers := map[string]string{
		"openai":    "unknown",
		"anthropic": "unknown",
		"google":    "unknown",
		"ollama":    "unknown",
	}

	status := statusOperational
	switch snap.Status {
	case health.Degraded:
		status = statusDegraded
	case health.Unhealthy:
		status = statusOutage
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     status,
		"checked_at": time.Now().UTC(),
		"components": snap.Components,
		"providers":  providers,
		"incidents":  []incident{},
	})
}

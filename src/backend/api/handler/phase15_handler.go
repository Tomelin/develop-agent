package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase15Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase15Service
}

type runPhase15Request struct {
	UseLinkedProject bool                         `json:"use_linked_project"`
	Channels         []string                     `json:"channels"`
	MonthlyBudgetUSD float64                      `json:"monthly_budget_usd"`
	ManualBrief      project.MarketingManualBrief `json:"manual_brief"`
}

type configureMarketingWebhookRequest struct {
	URL string `json:"url"`
}

func NewPhase15Handler(projects project.ProjectRepository, service *projectuc.Phase15Service) *Phase15Handler {
	return &Phase15Handler{projects: projects, service: service}
}

func (h *Phase15Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/15/run", h.Run)
	projects.GET("/:id/marketing/export", h.Export)
	projects.POST("/:id/marketing/webhooks", h.ConfigureWebhook)
}

func (h *Phase15Handler) Run(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req runPhase15Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.service.Run(c.Request.Context(), c.Param("id"), mustUserID(c), project.Phase15RunInput{
		UseLinkedProject: req.UseLinkedProject,
		ManualBrief:      req.ManualBrief,
		Channels:         req.Channels,
		MonthlyBudgetUSD: req.MonthlyBudgetUSD,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Phase15Handler) Export(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	channels := parseChannels(c.Query("channels"))
	raw, filename, pieces, err := h.service.ExportPack(c.Request.Context(), c.Param("id"), mustUserID(c), channels)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Header("X-Marketing-Pieces", toString(pieces))
	c.Data(http.StatusOK, "application/zip", raw)
}

func (h *Phase15Handler) ConfigureWebhook(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req configureMarketingWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	res, err := h.service.ConfigureWebhook(c.Request.Context(), c.Param("id"), mustUserID(c), project.MarketingWebhookInput{URL: req.URL})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Phase15Handler) canAccessProject(c *gin.Context) bool {
	p, err := h.projects.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return false
	}
	if p.OwnerUserID.Hex() != mustUserID(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return false
	}
	if p.OrganizationID.Hex() != mustOrganizationID(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return false
	}
	return true
}

func parseChannels(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func toString(n int) string {
	return strconv.Itoa(n)
}

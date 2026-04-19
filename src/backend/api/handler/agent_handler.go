package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/internal/domain/agent"
	"github.com/develop-agent/backend/internal/domain/user"
)

type AgentHandler struct {
	repo     agent.Repository
	validate *validator.Validate
}

type createAgentRequest struct {
	Name          string         `json:"name" validate:"required,min=2"`
	Description   string         `json:"description"`
	Provider      agent.Provider `json:"provider" validate:"required"`
	Model         string         `json:"model" validate:"required"`
	SystemPrompts []string       `json:"system_prompts"`
	Skills        []agent.Skill  `json:"skills" validate:"required,min=1"`
	Enabled       *bool          `json:"enabled"`
	ApiKeyRef     string         `json:"api_key_ref"`
}

type updateAgentRequest = createAgentRequest

func NewAgentHandler(repo agent.Repository) *AgentHandler {
	return &AgentHandler{repo: repo, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *AgentHandler) Register(rg *gin.RouterGroup) {
	agents := rg.Group("/agents")
	agents.GET("", h.List)
	agents.GET("/:id", h.GetByID)
	agents.POST("", middleware.RoleMiddleware(string(user.RoleAdmin)), h.Create)
	agents.PUT("/:id", middleware.RoleMiddleware(string(user.RoleAdmin)), h.Update)
	agents.DELETE("/:id", middleware.RoleMiddleware(string(user.RoleAdmin)), h.Delete)
	agents.POST("/:id/test", middleware.RoleMiddleware(string(user.RoleAdmin)), h.TestConnection)
}

func (h *AgentHandler) List(c *gin.Context) {
	filter := agent.ListFilter{}
	if enabledParam := strings.TrimSpace(c.Query("enabled")); enabledParam != "" {
		enabled := enabledParam == "true"
		filter.Enabled = &enabled
	}
	if provider := agent.Provider(strings.ToUpper(strings.TrimSpace(c.Query("provider")))); provider != "" {
		filter.Provider = provider
	}
	if skill := agent.Skill(strings.ToUpper(strings.TrimSpace(c.Query("skill")))); skill != "" {
		filter.Skill = skill
	}

	items, err := h.repo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list agents"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *AgentHandler) GetByID(c *gin.Context) {
	item, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req createAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !req.Provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}
	for _, sk := range req.Skills {
		if !sk.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill"})
			return
		}
	}
	if _, err := h.repo.FindByName(c.Request.Context(), req.Name); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "agent name already exists"})
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	a, err := agent.New(req.Name, req.Description, req.Provider, req.Model, req.SystemPrompts, req.Skills, enabled, req.ApiKeyRef)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(c.Request.Context(), a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create agent"})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AgentHandler) Update(c *gin.Context) {
	var req updateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	item, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	if !req.Provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}
	for _, sk := range req.Skills {
		if !sk.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill"})
			return
		}
	}

	item.Name = strings.TrimSpace(req.Name)
	item.Description = strings.TrimSpace(req.Description)
	item.Provider = req.Provider
	item.Model = strings.TrimSpace(req.Model)
	item.SystemPrompts = req.SystemPrompts
	item.Skills = req.Skills
	item.ApiKeyRef = strings.TrimSpace(req.ApiKeyRef)
	if req.Enabled != nil {
		item.Enabled = *req.Enabled
	}
	if err := h.repo.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update agent"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *AgentHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AgentHandler) TestConnection(c *gin.Context) {
	item, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	if !item.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent is disabled"})
		return
	}
	if strings.TrimSpace(item.ApiKeyRef) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent has no api key reference configured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":       true,
		"message":  "connectivity check passed",
		"provider": item.Provider,
		"model":    item.Model,
	})
}

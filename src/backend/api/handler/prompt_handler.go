package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "github.com/develop-agent/backend/internal/domain/prompt"
	promptuc "github.com/develop-agent/backend/internal/usecase/prompt"
)

type PromptHandler struct {
	repo     domain.UserPromptRepository
	service  *promptuc.Service
	validate *validator.Validate
}

type createPromptRequest struct {
	Title    string       `json:"title" validate:"required,min=2"`
	Content  string       `json:"content" validate:"required"`
	Group    domain.Group `json:"group" validate:"required"`
	Priority int          `json:"priority"`
	Enabled  *bool        `json:"enabled"`
	Tags     []string     `json:"tags"`
}

type updatePromptRequest = createPromptRequest

type reorderPromptsRequest struct {
	Items []domain.ReorderItem `json:"items" validate:"required,min=1"`
}

type previewRequest struct {
	AgentSystemPrompts []string `json:"agent_system_prompts"`
	RAGContext         string   `json:"rag_context"`
	PhaseInstruction   string   `json:"phase_instruction"`
	UserInstruction    string   `json:"user_instruction"`
}

type fromTemplateRequest struct {
	TemplateID string `json:"template_id" validate:"required"`
	Priority   int    `json:"priority"`
}

func NewPromptHandler(repo domain.UserPromptRepository, service *promptuc.Service) *PromptHandler {
	return &PromptHandler{repo: repo, service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *PromptHandler) Register(rg *gin.RouterGroup) {
	prompts := rg.Group("/prompts")
	prompts.GET("", h.List)
	prompts.GET("/templates", h.Templates)
	prompts.POST("/from-template", h.FromTemplate)
	prompts.PUT("/reorder", h.Reorder)
	prompts.GET("/preview/:group", h.Preview)
	prompts.GET("/export", h.Export)
	prompts.POST("/import", h.Import)
	prompts.GET("/:group", h.ListByGroup)
	prompts.POST("", h.Create)
	prompts.PUT("/:id", h.Update)
	prompts.DELETE("/:id", h.Delete)
}

func (h *PromptHandler) List(c *gin.Context) {
	filter := domain.ListFilter{UserID: mustUserID(c)}
	if group := domain.Group(strings.ToUpper(strings.TrimSpace(c.Query("group")))); group != "" {
		if !group.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group"})
			return
		}
		filter.Group = group
	}
	if enabledRaw := strings.TrimSpace(c.Query("enabled")); enabledRaw != "" {
		enabled := enabledRaw == "true"
		filter.Enabled = &enabled
	}
	items, err := h.repo.FindAllByUser(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *PromptHandler) ListByGroup(c *gin.Context) {
	group := domain.Group(strings.ToUpper(strings.TrimSpace(c.Param("group"))))
	if !group.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group"})
		return
	}
	items, err := h.repo.FindByUserAndGroup(c.Request.Context(), mustUserID(c), group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *PromptHandler) Create(c *gin.Context) {
	var req createPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil || !req.Group.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	created, err := h.service.Create(c.Request.Context(), mustUserID(c), req.Title, req.Content, req.Group, req.Priority, enabled, req.Tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": created, "token_estimate": domain.EstimateTokens(created.Content), "token_warning": domain.EstimateTokens(created.Content) > domain.WarnTokens})
}

func (h *PromptHandler) Update(c *gin.Context) {
	var req updatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil || !req.Group.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	updated, err := h.service.Update(c.Request.Context(), mustUserID(c), c.Param("id"), req.Title, req.Content, req.Group, req.Priority, enabled, req.Tags)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": updated, "token_estimate": domain.EstimateTokens(updated.Content), "token_warning": domain.EstimateTokens(updated.Content) > domain.WarnTokens})
}

func (h *PromptHandler) Delete(c *gin.Context) {
	if err := h.repo.Delete(c.Request.Context(), mustUserID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prompt not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PromptHandler) Reorder(c *gin.Context) {
	var req reorderPromptsRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.repo.Reorder(c.Request.Context(), mustUserID(c), req.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PromptHandler) Preview(c *gin.Context) {
	group := domain.Group(strings.ToUpper(strings.TrimSpace(c.Param("group"))))
	if !group.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group"})
		return
	}
	in := previewRequest{AgentSystemPrompts: []string{}}
	_ = c.ShouldBindJSON(&in)
	messages, tokenEstimate := h.service.Preview(c.Request.Context(), mustUserID(c), group, in.AgentSystemPrompts, in.RAGContext, in.PhaseInstruction, in.UserInstruction)
	c.JSON(http.StatusOK, gin.H{"group": group, "messages": messages, "token_estimate": tokenEstimate, "token_warning": tokenEstimate > domain.WarnTokens})
}

func (h *PromptHandler) Templates(c *gin.Context) {
	byGroup := map[domain.Group][]promptuc.Template{}
	for _, t := range h.service.Templates() {
		byGroup[t.Group] = append(byGroup[t.Group], t)
	}
	c.JSON(http.StatusOK, gin.H{"items": byGroup})
}

func (h *PromptHandler) FromTemplate(c *gin.Context) {
	var req fromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	item, err := h.service.CreateFromTemplate(c.Request.Context(), mustUserID(c), req.TemplateID, req.Priority)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h *PromptHandler) Export(c *gin.Context) {
	raw, err := h.service.Export(c.Request.Context(), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export prompts"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=prompts-export.json")
	c.Data(http.StatusOK, "application/json", raw)
}

func (h *PromptHandler) Import(c *gin.Context) {
	replace := c.Query("mode") == "replace"
	if mode := strings.TrimSpace(c.Query("mode")); mode != "" {
		replace = mode == "replace"
	}
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	created, err := h.service.Import(c.Request.Context(), mustUserID(c), raw, replace)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"created": created, "mode": map[bool]string{true: "replace", false: "merge"}[replace]})
}

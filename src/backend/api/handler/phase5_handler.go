package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase5Handler struct {
	projects project.ProjectRepository
	service  *projectuc.DevelopmentService
	validate *validator.Validate
}

type setPhase5ModeRequest struct {
	Mode project.ExecutionMode `json:"mode" validate:"required"`
}

type autoRejectionRequest struct {
	CoveragePercent     float64  `json:"coverage_percent"`
	MaxCVSS             float64  `json:"max_cvss"`
	CredentialsExposed  bool     `json:"credentials_exposed"`
	CompilationFailed   bool     `json:"compilation_failed"`
	FailureDescriptions []string `json:"failure_descriptions"`
	SourcePhase         int      `json:"source_phase"`
}

func NewPhase5Handler(projects project.ProjectRepository, service *projectuc.DevelopmentService) *Phase5Handler {
	return &Phase5Handler{projects: projects, service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *Phase5Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/5/mode", h.SetMode)
	projects.POST("/:id/phases/5/execute", h.ExecuteAll)
	projects.POST("/:id/phases/5/tasks/:taskId/execute", h.ExecuteTask)
	projects.GET("/:id/phases/5/summary", h.Summary)
	projects.GET("/:id/phases/5/code-context", h.CodeContext)
	projects.POST("/:id/phases/5/auto-rejection", h.AutoRejection)
	projects.GET("/:id/files", h.ListFiles)
	projects.GET("/:id/files/download", h.DownloadFiles)
	projects.GET("/:id/files/:fileId", h.GetFile)
}

func (h *Phase5Handler) SetMode(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req setPhase5ModeRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil || !req.Mode.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	mode, err := h.service.SetExecutionMode(c.Request.Context(), c.Param("id"), mustUserID(c), req.Mode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mode": mode})
}

func (h *Phase5Handler) ExecuteAll(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	total, err := h.service.ExecuteAllPending(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"executed_tasks": total})
}

func (h *Phase5Handler) ExecuteTask(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	if err := h.service.ExecuteTask(c.Request.Context(), c.Param("id"), mustUserID(c), c.Param("taskId")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Phase5Handler) Summary(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	summary, err := h.service.Summary(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *Phase5Handler) CodeContext(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	manifest, err := h.service.BuildCodeContext(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, manifest)
}

func (h *Phase5Handler) AutoRejection(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req autoRejectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	report := project.CatastrophicFailureReport{
		CoveragePercent:      req.CoveragePercent,
		MaxCVSS:              req.MaxCVSS,
		CredentialsExposed:   req.CredentialsExposed,
		CompilationFailed:    req.CompilationFailed,
		FailureDescriptions:  req.FailureDescriptions,
		SourcePhase:          req.SourcePhase,
		RequestedBySystemTag: "quality-gate",
	}
	if err := h.service.TriggerAutoRejection(c.Request.Context(), c.Param("id"), mustUserID(c), report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "auto-rejection-triggered"})
}

func (h *Phase5Handler) ListFiles(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	files, err := h.service.ListFiles(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": files})
}

func (h *Phase5Handler) GetFile(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	file, err := h.service.GetFile(c.Request.Context(), c.Param("id"), c.Param("fileId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

func (h *Phase5Handler) DownloadFiles(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	zipBytes, err := h.service.DownloadZIP(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=phase5-files.zip")
	c.Data(http.StatusOK, "application/zip", zipBytes)
}

func (h *Phase5Handler) canAccessProject(c *gin.Context) bool {
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

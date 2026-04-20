package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase7Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase7Service
	validate *validator.Validate
}

type runSecurityAuditRequest struct {
	BackendDir     string `json:"backend_dir"`
	FrontendDir    string `json:"frontend_dir"`
	ProjectRootDir string `json:"project_root_dir"`
	HighRetryCount int    `json:"high_retry_count"`
}

func NewPhase7Handler(projects project.ProjectRepository, service *projectuc.Phase7Service) *Phase7Handler {
	return &Phase7Handler{projects: projects, service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *Phase7Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/7/run-audit", h.RunAudit)
}

func (h *Phase7Handler) RunAudit(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req runSecurityAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	report, err := h.service.RunAudit(c.Request.Context(), c.Param("id"), mustUserID(c), projectuc.Phase7AuditInput{
		BackendDir:     req.BackendDir,
		FrontendDir:    req.FrontendDir,
		ProjectRootDir: req.ProjectRootDir,
		HighRetryCount: req.HighRetryCount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Phase7Handler) canAccessProject(c *gin.Context) bool {
	p, err := h.projects.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return false
	}
	if p.OwnerUserID.Hex() != mustUserID(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return false
	}
	return true
}

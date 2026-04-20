package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase6Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase6Service
	validate *validator.Validate
}

type analyzeCoverageRequest struct {
	BackendDir string  `json:"backend_dir"`
	Threshold  float64 `json:"threshold"`
}

type validateTestsRequest struct {
	BackendDir  string `json:"backend_dir"`
	FrontendDir string `json:"frontend_dir"`
}

func NewPhase6Handler(projects project.ProjectRepository, service *projectuc.Phase6Service) *Phase6Handler {
	return &Phase6Handler{projects: projects, service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *Phase6Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/6/analyze-coverage", h.AnalyzeCoverage)
	projects.POST("/:id/phases/6/validate-tests", h.ValidateTests)
}

func (h *Phase6Handler) AnalyzeCoverage(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req analyzeCoverageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	report, below, err := h.service.AnalyzeCoverage(c.Request.Context(), c.Param("id"), mustUserID(c), req.BackendDir, req.Threshold)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quality := h.service.BuildQualityReport(report, nil)
	if err := h.service.PersistQualityReport(c.Request.Context(), c.Param("id"), quality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"below_threshold": below, "report": report})
}

func (h *Phase6Handler) ValidateTests(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req validateTestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	res, err := h.service.ValidateTests(c.Request.Context(), req.BackendDir, req.FrontendDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Phase6Handler) canAccessProject(c *gin.Context) bool {
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

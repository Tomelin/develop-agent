package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase14Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase14Service
}

type runPhase14Request struct {
	UseLinkedProject bool                           `json:"use_linked_project"`
	GenerateVariants bool                           `json:"generate_variants"`
	VariantCount     int                            `json:"variant_count"`
	ManualBrief      project.LandingPageManualBrief `json:"manual_brief"`
}

func NewPhase14Handler(projects project.ProjectRepository, service *projectuc.Phase14Service) *Phase14Handler {
	return &Phase14Handler{projects: projects, service: service}
}

func (h *Phase14Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/14/run", h.Run)
}

func (h *Phase14Handler) Run(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req runPhase14Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	result, err := h.service.Run(c.Request.Context(), c.Param("id"), mustUserID(c), project.Phase14RunInput{
		UseLinkedProject: req.UseLinkedProject,
		ManualBrief:      req.ManualBrief,
		GenerateVariants: req.GenerateVariants,
		VariantCount:     req.VariantCount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Phase14Handler) canAccessProject(c *gin.Context) bool {
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

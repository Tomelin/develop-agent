package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase13Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase13Service
}

type runPhase13Request struct {
	BackendBaseURL string `json:"backend_base_url"`
	FrontendURL    string `json:"frontend_url"`
	IncludeDevOps  *bool  `json:"include_devops"`
}

func NewPhase13Handler(projects project.ProjectRepository, service *projectuc.Phase13Service) *Phase13Handler {
	return &Phase13Handler{projects: projects, service: service}
}

func (h *Phase13Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/13/run", h.Run)
}

func (h *Phase13Handler) Run(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req runPhase13Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	includeDevOps := true
	if req.IncludeDevOps != nil {
		includeDevOps = *req.IncludeDevOps
	}
	result, err := h.service.Run(c.Request.Context(), c.Param("id"), mustUserID(c), project.Phase13RunInput{
		BackendBaseURL: req.BackendBaseURL,
		FrontendURL:    req.FrontendURL,
		IncludeDevOps:  includeDevOps,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Phase13Handler) canAccessProject(c *gin.Context) bool {
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

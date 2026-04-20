package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type Phase19Handler struct {
	projects project.ProjectRepository
	service  *projectuc.Phase19Service
}

type runPhase19Request struct {
	IncludeFlowA bool `json:"include_flow_a"`
	IncludeFlowB bool `json:"include_flow_b"`
	IncludeFlowC bool `json:"include_flow_c"`
}

func NewPhase19Handler(projects project.ProjectRepository, service *projectuc.Phase19Service) *Phase19Handler {
	return &Phase19Handler{projects: projects, service: service}
}

func (h *Phase19Handler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.POST("/:id/phases/19/run", h.Run)
}

func (h *Phase19Handler) Run(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req runPhase19Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.service.Run(c.Request.Context(), c.Param("id"), mustUserID(c), project.Phase19RunInput{
		IncludeFlowA: req.IncludeFlowA,
		IncludeFlowB: req.IncludeFlowB,
		IncludeFlowC: req.IncludeFlowC,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Phase19Handler) canAccessProject(c *gin.Context) bool {
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

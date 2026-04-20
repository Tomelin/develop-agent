package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/internal/domain/project"
	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type ProjectHandler struct {
	repo     project.ProjectRepository
	service  *projectuc.Service
	validate *validator.Validate
}

type createProjectRequest struct {
	Name               string           `json:"name" validate:"required,min=2"`
	Description        string           `json:"description"`
	FlowType           project.FlowType `json:"flow_type" validate:"required"`
	LinkedProjectID    string           `json:"linked_project_id"`
	DynamicModeEnabled bool             `json:"dynamic_mode_enabled"`
}

type updateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
}

func NewProjectHandler(repo project.ProjectRepository, service *projectuc.Service) *ProjectHandler {
	return &ProjectHandler{repo: repo, service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *ProjectHandler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	projects.GET("", h.List)
	projects.GET("/:id", h.GetByID)
	projects.GET("/:id/phases/:phaseNumber/tracks", h.GetPhaseTracks)
	projects.POST("", h.Create)
	projects.PUT("/:id", h.Update)
	projects.POST("/:id/pause", h.Pause)
	projects.POST("/:id/resume", h.Resume)
	projects.POST("/:id/archive", h.Archive)
	projects.POST("/:id/phases/:phaseNumber/start", h.StartPhase)
	projects.POST("/:id/phases/:phaseNumber/tracks/:track/approve", h.ApprovePhaseTrack)
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID := mustUserID(c)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	filter := project.ProjectListFilter{
		OwnerID:  userID,
		Status:   project.ProjectStatus(strings.ToUpper(strings.TrimSpace(c.Query("status")))),
		FlowType: project.FlowType(strings.ToUpper(strings.TrimSpace(c.Query("flow_type")))),
		Page:     page,
		Limit:    limit,
	}
	items, total, err := h.repo.FindDashboardByOwner(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total, "page": page, "limit": limit})
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID := mustUserID(c)
	p, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if p.OwnerUserID.Hex() != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if !req.FlowType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flow type"})
		return
	}

	created, err := h.service.CreateProject(c.Request.Context(), projectuc.CreateProjectInput{
		Name:               req.Name,
		Description:        req.Description,
		FlowType:           req.FlowType,
		OwnerUserID:        mustUserID(c),
		LinkedProjectID:    req.LinkedProjectID,
		DynamicModeEnabled: req.DynamicModeEnabled,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := mustUserID(c)
	p, err := h.repo.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	if p.OwnerUserID.Hex() != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	if p.Status != project.ProjectDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only DRAFT project can be edited"})
		return
	}
	p.Name = strings.TrimSpace(req.Name)
	p.Description = strings.TrimSpace(req.Description)
	if err := h.repo.Update(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Pause(c *gin.Context) {
	p, err := h.service.Pause(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Resume(c *gin.Context) {
	p, err := h.service.Resume(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Archive(c *gin.Context) {
	if err := h.service.Archive(c.Request.Context(), c.Param("id"), mustUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) StartPhase(c *gin.Context) {
	phaseNumber, err := strconv.Atoi(c.Param("phaseNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phase number"})
		return
	}
	phase, err := h.service.StartPhase(c.Request.Context(), c.Param("id"), mustUserID(c), phaseNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, phase)
}

func (h *ProjectHandler) GetPhaseTracks(c *gin.Context) {
	phaseNumber, err := strconv.Atoi(c.Param("phaseNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phase number"})
		return
	}
	tracks, err := h.service.GetPhaseTracks(c.Request.Context(), c.Param("id"), mustUserID(c), phaseNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func (h *ProjectHandler) ApprovePhaseTrack(c *gin.Context) {
	phaseNumber, err := strconv.Atoi(c.Param("phaseNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phase number"})
		return
	}
	track := project.Track(strings.ToUpper(strings.TrimSpace(c.Param("track"))))
	if track != project.TrackFrontend && track != project.TrackBackend && track != project.TrackFull {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track"})
		return
	}
	phase, err := h.service.ApprovePhaseTrack(c.Request.Context(), c.Param("id"), mustUserID(c), phaseNumber, track)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, phase)
}

func mustUserID(c *gin.Context) string {
	ctxRaw, _ := c.Get(middleware.UserContextKey)
	ctx := ctxRaw.(gin.H)
	return ctx["user_id"].(string)
}

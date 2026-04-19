package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/develop-agent/backend/internal/domain/project"
)

type TaskHandler struct {
	tasks    project.TaskRepository
	projects project.ProjectRepository
	validate *validator.Validate
}

type createTaskRequest struct {
	PhaseID         string                 `json:"phase_id"`
	EpicID          string                 `json:"epic_id"`
	Title           string                 `json:"title" validate:"required,min=2"`
	Description     string                 `json:"description"`
	Type            project.TaskType       `json:"type" validate:"required"`
	Complexity      project.TaskComplexity `json:"complexity" validate:"required"`
	EstimatedHours  float64                `json:"estimated_hours"`
	AssignedAgentID string                 `json:"assigned_agent_id"`
}

type updateTaskStatusRequest struct {
	Status project.TaskStatus `json:"status" validate:"required"`
}

func NewTaskHandler(tasks project.TaskRepository, projects project.ProjectRepository) *TaskHandler {
	return &TaskHandler{tasks: tasks, projects: projects, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *TaskHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/projects/:id/tasks", h.List)
	rg.POST("/projects/:id/tasks/bulk", h.BulkCreate)
	rg.PUT("/projects/:id/tasks/:taskId/status", h.UpdateStatus)
}

func (h *TaskHandler) List(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	items, err := h.tasks.ListByProject(c.Request.Context(), project.TaskListFilter{
		ProjectID:  c.Param("id"),
		Type:       project.TaskType(strings.ToUpper(strings.TrimSpace(c.Query("type")))),
		Complexity: project.TaskComplexity(strings.ToUpper(strings.TrimSpace(c.Query("complexity")))),
		Status:     project.TaskStatus(strings.ToUpper(strings.TrimSpace(c.Query("status")))),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *TaskHandler) BulkCreate(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req []createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	projectID, _ := bson.ObjectIDFromHex(c.Param("id"))
	tasks := make([]*project.Task, 0, len(req))
	for _, item := range req {
		if err := h.validate.Struct(item); err != nil || !item.Type.IsValid() || !item.Complexity.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task payload"})
			return
		}
		tasks = append(tasks, &project.Task{
			ProjectID:       projectID,
			PhaseID:         strings.TrimSpace(item.PhaseID),
			EpicID:          strings.TrimSpace(item.EpicID),
			Title:           strings.TrimSpace(item.Title),
			Description:     strings.TrimSpace(item.Description),
			Type:            item.Type,
			Complexity:      item.Complexity,
			EstimatedHours:  item.EstimatedHours,
			Status:          project.TaskTodo,
			AssignedAgentID: strings.TrimSpace(item.AssignedAgentID),
		})
	}
	if err := h.tasks.BulkCreate(c.Request.Context(), tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tasks"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"created": len(tasks)})
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	if !h.canAccessProject(c) {
		return
	}
	var req updateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil || !req.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.tasks.UpdateStatus(c.Request.Context(), c.Param("id"), c.Param("taskId"), req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) canAccessProject(c *gin.Context) bool {
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

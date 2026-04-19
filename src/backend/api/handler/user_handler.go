package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/internal/domain/user"
)

type UserHandler struct {
	repo     user.Repository
	validate *validator.Validate
}

type updateMeRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type updatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

func NewUserHandler(repo user.Repository) *UserHandler {
	return &UserHandler{repo: repo, validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *UserHandler) Register(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/me", h.GetMe)
	users.PUT("/me", h.UpdateMe)
	users.PUT("/me/password", h.UpdatePassword)
	users.GET("", middleware.RoleMiddleware(string(user.RoleAdmin)), h.List)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	ctxRaw, _ := c.Get(middleware.UserContextKey)
	ctx := ctxRaw.(gin.H)
	u, err := h.repo.FindByID(c.Request.Context(), ctx["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u.Sanitize())
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	ctxRaw, _ := c.Get(middleware.UserContextKey)
	ctx := ctxRaw.(gin.H)
	u, err := h.repo.FindByID(c.Request.Context(), ctx["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	u.Name = req.Name
	if err := h.repo.Update(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, u.Sanitize())
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req updatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := user.ValidatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctxRaw, _ := c.Get(middleware.UserContextKey)
	ctx := ctxRaw.(gin.H)
	u, err := h.repo.FindByID(c.Request.Context(), ctx["user_id"].(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err := u.CheckPassword(req.CurrentPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
		return
	}
	updated, err := user.New(u.Name, u.Email, req.NewPassword, u.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.PasswordHash = updated.PasswordHash
	if err := h.repo.Update(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/api/middleware"
	"github.com/develop-agent/backend/internal/domain/user"
	orguc "github.com/develop-agent/backend/internal/usecase/organization"
)

type OrganizationHandler struct {
	service  *orguc.Service
	validate *validator.Validate
}

type inviteMemberRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=OWNER ADMIN MEMBER VIEWER"`
}

type updateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=OWNER ADMIN MEMBER VIEWER"`
}

func NewOrganizationHandler(service *orguc.Service) *OrganizationHandler {
	return &OrganizationHandler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *OrganizationHandler) Register(rg *gin.RouterGroup) {
	org := rg.Group("/org")
	org.GET("", h.Get)
	org.GET("/members", h.ListMembers)
	org.POST("/invite", middleware.OrganizationRoleMiddleware(string(user.OrganizationRoleOwner), string(user.OrganizationRoleAdmin)), h.InviteMember)
	org.PUT("/members/:userId/role", middleware.OrganizationRoleMiddleware(string(user.OrganizationRoleOwner), string(user.OrganizationRoleAdmin)), h.UpdateMemberRole)
	org.DELETE("/members/:userId", middleware.OrganizationRoleMiddleware(string(user.OrganizationRoleOwner)), h.RemoveMember)
}

func (h *OrganizationHandler) Get(c *gin.Context) {
	org, err := h.service.GetOrganization(c.Request.Context(), mustOrganizationID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) ListMembers(c *gin.Context) {
	items, err := h.service.ListMembers(c.Request.Context(), mustOrganizationID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (h *OrganizationHandler) InviteMember(c *gin.Context) {
	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.service.InviteMember(c.Request.Context(), orguc.InviteInput{
		OrganizationID:   mustOrganizationID(c),
		Name:             req.Name,
		Email:            req.Email,
		OrganizationRole: user.OrganizationRole(strings.ToUpper(strings.TrimSpace(req.Role))),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *OrganizationHandler) UpdateMemberRole(c *gin.Context) {
	var req updateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.service.UpdateMemberRole(
		c.Request.Context(),
		mustOrganizationID(c),
		c.Param("userId"),
		user.OrganizationRole(strings.ToUpper(strings.TrimSpace(req.Role))),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	if err := h.service.RemoveMember(c.Request.Context(), mustOrganizationID(c), c.Param("userId")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

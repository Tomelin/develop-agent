package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	domain "github.com/develop-agent/backend/internal/domain/billing"
	billinguc "github.com/develop-agent/backend/internal/usecase/billing"
)

type BillingHandler struct {
	service *billinguc.Service
}

func NewBillingHandler(service *billinguc.Service) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) Register(rg *gin.RouterGroup) {
	billing := rg.Group("/billing")
	billing.GET("/pricing", h.Pricing)
	billing.GET("/summary", h.Summary)
	billing.GET("/by-model", h.ByModel)
	billing.GET("/by-phase", h.ByPhase)
	billing.GET("/top-projects", h.TopProjects)
	billing.GET("/export", h.Export)

	rg.GET("/projects/:id/billing", h.ProjectBilling)
}

func (h *BillingHandler) Pricing(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Pricing(c.Request.Context()))
}

func (h *BillingHandler) Summary(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.Summary(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *BillingHandler) ProjectBilling(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter.ProjectID = c.Param("id")
	data, err := h.service.ProjectDetails(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *BillingHandler) ByModel(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.ByModel(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data})
}

func (h *BillingHandler) ByPhase(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.ByPhase(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data})
}

func (h *BillingHandler) TopProjects(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.TopProjects(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": data})
}

func (h *BillingHandler) Export(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter.ProjectID = strings.TrimSpace(c.Query("project_id"))
	filter.Provider = strings.ToUpper(strings.TrimSpace(c.Query("provider")))
	format := c.DefaultQuery("format", "csv")
	payload, contentType, err := h.service.Export(c.Request.Context(), filter, format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=billing-export."+strings.ToLower(format))
	c.Data(http.StatusOK, contentType, payload)
}

func parseFilter(c *gin.Context) (domain.QueryFilter, error) {
	from, to, err := billinguc.ParseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		return domain.QueryFilter{}, err
	}
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	return domain.QueryFilter{UserID: mustUserID(c), From: from, To: to, Page: page, Limit: limit}, nil
}

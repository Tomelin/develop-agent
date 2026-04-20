package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	projectuc "github.com/develop-agent/backend/internal/usecase/project"
)

type AdminQualityHandler struct {
	service *projectuc.AdminQualityReportService
}

func NewAdminQualityHandler(service *projectuc.AdminQualityReportService) *AdminQualityHandler {
	return &AdminQualityHandler{service: service}
}

func (h *AdminQualityHandler) Register(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	admin.GET("/quality-report", h.Report)
}

func (h *AdminQualityHandler) Report(c *gin.Context) {
	report, err := h.service.Build(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build report"})
		return
	}
	c.JSON(http.StatusOK, report)
}

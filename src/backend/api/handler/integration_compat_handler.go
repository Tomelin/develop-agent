package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IntegrationCompatHandler expõe endpoints ainda não implementados em profundidade,
// mas necessários para manter o contrato de integração com o frontend.
type IntegrationCompatHandler struct{}

func NewIntegrationCompatHandler() *IntegrationCompatHandler { return &IntegrationCompatHandler{} }

func (h *IntegrationCompatHandler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		// Phase 8
		projects.GET("/:id/phases/:phaseNumber/artifacts", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.GET("/:id/phases/:phaseNumber/triad-progress", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.GET("/:id/phases/:phaseNumber/tracks/:track/feedbacks", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.POST("/:id/phases/:phaseNumber/tracks/:track/feedback", func(c *gin.Context) { c.Status(http.StatusNoContent) })

		// Roadmap compat
		projects.GET("/:id/roadmap", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"project_id": c.Param("id"), "phases": []any{}})
		})

		// Phase 17
		projects.GET("/:id/triad-selections", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.GET("/:id/selection-logs", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.PUT("/:id/dynamic-mode", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"enabled": true}) })
		projects.GET("/:id/dynamic-mode/preview", func(c *gin.Context) {
			emptyTriad := gin.H{
				"producer": gin.H{"id": "", "name": "", "provider": "OPENAI", "model": ""},
				"reviewer": gin.H{"id": "", "name": "", "provider": "OPENAI", "model": ""},
				"refiner":  gin.H{"id": "", "name": "", "provider": "OPENAI", "model": ""},
			}
			c.JSON(http.StatusOK, gin.H{"eligible_agents": 0, "triad": emptyTriad, "notes": []string{}})
		})
		projects.GET("/:id/diversity-metrics", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"project_id":               c.Param("id"),
				"diversity_score":          0,
				"providers":                []any{},
				"models":                   []string{},
				"full_diversity_triads":    0,
				"repeated_provider_triads": 0,
				"role_distribution":        []any{},
			})
		})
		projects.GET("/:id/agent-config/matrix", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"rows": []any{}}) })
		projects.PUT("/:id/agent-config/matrix", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"rows": []any{}}) })
		projects.POST("/:id/agent-config/cost-preview", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"monthly_estimated_usd": 0, "note": "preview compat"})
		})

		// Phase 20 - colaboração / integrações
		projects.GET("/:id/collaborators", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.POST("/:id/collaborators", func(c *gin.Context) { c.Status(http.StatusNoContent) })
		projects.PUT("/:id/collaborators/:userId/role", func(c *gin.Context) { c.Status(http.StatusNoContent) })
		projects.DELETE("/:id/collaborators/:userId", func(c *gin.Context) { c.Status(http.StatusNoContent) })
		projects.GET("/:id/integrations", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
		projects.POST("/:id/integrations/jira/sync", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	}

	rg.GET("/notifications", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
	rg.POST("/notifications/:id/read", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Phase 17 - admin/flags
	rg.GET("/admin/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"workers": gin.H{"max_concurrency": 1, "agent_timeout_seconds": 60, "triad_timeout_seconds": 120},
			"models":  gin.H{"default_model": "", "spec_generation_model": ""},
			"limits":  gin.H{"max_projects_per_user": 1, "max_parallel_phases_per_user": 1, "max_spec_tokens": 4000},
			"retry":   gin.H{"max_attempts": 1, "backoff_seconds": 1},
		})
	})
	rg.PUT("/admin/settings", func(c *gin.Context) {
		var body map[string]any
		_ = c.ShouldBindJSON(&body)
		c.JSON(http.StatusOK, body)
	})
	rg.GET("/admin/feature-flags", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
	rg.PUT("/admin/feature-flags", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
	rg.GET("/feature-flags", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })

	// Phase 20 - marketplace/integrations/pricing/roadmap público
	rg.GET("/marketplace/templates", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": []any{}, "total": 0, "page": 1, "size": 20})
	})
	rg.POST("/marketplace/templates", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.POST("/marketplace/templates/:templateId/use", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.POST("/marketplace/templates/:templateId/star", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.GET("/integrations/github/auth", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"auth_url": ""}) })
	rg.POST("/integrations/jira", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.POST("/integrations/slack/webhook", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.GET("/pricing/plans", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
	rg.POST("/pricing/checkout", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"checkout_url": ""}) })
	rg.GET("/roadmap/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"vision": "", "features": []any{}, "changelog": []any{}})
	})
	rg.POST("/roadmap/features/:featureId/vote", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.POST("/roadmap/features/suggestions", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	rg.PUT("/admin/roadmap/features/:featureId/status", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Compat com frontend legado
	rg.POST("/agents/test-config", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
}

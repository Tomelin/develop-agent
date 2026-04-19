package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecaseauth "github.com/develop-agent/backend/internal/usecase/auth"
)

const UserContextKey = "user_context"

func AuthMiddleware(authService *usecaseauth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(UserContextKey, gin.H{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"role":    claims.Role,
		})
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxRaw, ok := c.Get(UserContextKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
			return
		}
		userCtx, ok := ctxRaw.(gin.H)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}
		role, _ := userCtx["role"].(string)
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}

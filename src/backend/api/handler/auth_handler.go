package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/develop-agent/backend/api/middleware"
	usecaseauth "github.com/develop-agent/backend/internal/usecase/auth"
)

const refreshTokenCookieName = "refresh_token"

type AuthHandler struct {
	authService *usecaseauth.Service
	validate    *validator.Validate
	limiter     *loginRateLimiter
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginRateLimiter struct {
	mu      sync.Mutex
	attempt map[string][]time.Time
}

func NewAuthHandler(authService *usecaseauth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
		limiter:     &loginRateLimiter{attempt: map[string][]time.Time{}},
	}
}

func (h *AuthHandler) Register(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/logout", h.Logout)
	auth.GET("/me", h.Me)
}

func (h *AuthHandler) Login(c *gin.Context) {
	ip := c.ClientIP()
	if !h.limiter.Allow(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts"})
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, resp.RefreshToken, resp.RefreshExpiresAt)
	resp.RefreshToken = ""
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	resp, err := h.authService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, resp.RefreshToken, resp.RefreshExpiresAt)
	resp.RefreshToken = ""
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshTokenCookieName)
	if refreshToken != "" {
		if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
			return
		}
	}
	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	ctx, ok := c.Get(middleware.UserContextKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}
	c.JSON(http.StatusOK, ctx)
}

func (l *loginRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UTC()
	cutoff := now.Add(-1 * time.Minute)
	entries := l.attempt[ip]
	filtered := entries[:0]
	for _, t := range entries {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= 5 {
		l.attempt[ip] = filtered
		return false
	}
	filtered = append(filtered, now)
	l.attempt[ip] = filtered
	return true
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, refreshToken string, refreshExpiresAt string) {
	exp, err := time.Parse("2006-01-02T15:04:05Z", refreshExpiresAt)
	if err != nil {
		exp = time.Now().UTC().Add(7 * 24 * time.Hour)
	}
	maxAge := int(time.Until(exp).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		refreshToken,
		maxAge,
		"/api/v1/auth",
		"",
		c.Request.TLS != nil,
		true,
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		"/api/v1/auth",
		"",
		c.Request.TLS != nil,
		true,
	)
}

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	interviewuc "github.com/develop-agent/backend/internal/usecase/interview"
)

type InterviewHandler struct {
	service *interviewuc.Service
}

type sendMessageRequest struct {
	Content string `json:"content"`
}

func NewInterviewHandler(service *interviewuc.Service) *InterviewHandler {
	return &InterviewHandler{service: service}
}

func (h *InterviewHandler) Register(rg *gin.RouterGroup) {
	group := rg.Group("/projects/:id/interview")
	group.GET("", h.GetSession)
	group.POST("/message", h.SendMessage)
	group.POST("/confirm", h.Confirm)
	group.POST("/regenerate-vision", h.RegenerateVision)
	group.GET("/events", h.Events)
}

func (h *InterviewHandler) GetSession(c *gin.Context) {
	session, err := h.service.GetSession(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *InterviewHandler) SendMessage(c *gin.Context) {
	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	session, _, err := h.service.StreamMessage(c.Request.Context(), c.Param("id"), mustUserID(c), req.Content, func(delta string) error {
		if _, wErr := fmt.Fprintf(c.Writer, "event: token\ndata: %s\n\n", jsonEscape(delta)); wErr != nil {
			return wErr
		}
		flusher.Flush()
		return nil
	})
	if err != nil {
		_, _ = fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", jsonEscape(err.Error()))
		flusher.Flush()
		return
	}

	payload, _ := json.Marshal(session)
	_, _ = fmt.Fprintf(c.Writer, "event: done\ndata: %s\n\n", string(payload))
	flusher.Flush()
}

func (h *InterviewHandler) Confirm(c *gin.Context) {
	session, err := h.service.Confirm(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *InterviewHandler) RegenerateVision(c *gin.Context) {
	session, err := h.service.RegenerateVision(c.Request.Context(), c.Param("id"), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *InterviewHandler) Events(c *gin.Context) {
	ch := h.service.Broker().Subscribe(c.Param("id"))
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case ev := <-ch:
			c.SSEvent(ev.Type, ev)
			return true
		}
	})
}

func jsonEscape(v string) string {
	b, _ := json.Marshal(v)
	return string(b)
}

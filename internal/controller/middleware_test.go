package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-blog/config"
	"mini-blog/internal/service"
	"mini-blog/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUser_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthenticateUser)
	r.GET("/", func(c *gin.Context) {
		userID, exists := c.Get(UserIdKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not set"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Инициализация конфига для JWT
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL:   12,
		},
	})

	token, _ := service.GenerateToken(nil, 123)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestAuthenticateUser_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthenticateUser)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthenticateUser_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthenticateUser)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestRequestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.NewTestLogger()
	r := gin.New()
	r.Use(RequestLoggingMiddleware(log))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping?foo=bar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

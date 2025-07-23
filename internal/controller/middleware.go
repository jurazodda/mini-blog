package controller

import (
	"mini-blog/internal/service"
	"net/http"
	"strings"
	"time"

	"mini-blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

var (
	UserIdKey = "user_id"
)

func AuthenticateUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	claims, err := service.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userIDFloat, ok := claims[UserIdKey].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	c.Set(UserIdKey, int(userIDFloat))
}

func RequestLoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build log entry
		logEntry := logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent())

		if raw != "" {
			logEntry = logEntry.Str("query", raw)
		}

		if len(c.Errors) > 0 {
			logEntry = logEntry.Str("errors", c.Errors.String())
		}

		logEntry.Msg("HTTP request processed")
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

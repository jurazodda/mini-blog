package controller

import (
	"errors"
	"mini-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
	case errors.Is(err, service.ErrUserIdRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
	case errors.Is(err, service.ErrPostIdRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
	case errors.Is(err, service.ErrCommentIdRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment ID is required"})
	case errors.Is(err, service.ErrTitleRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
	case errors.Is(err, service.ErrContentRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
	case errors.Is(err, service.ErrCommentRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment is required"})
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden access"})
	case errors.Is(err, service.ErrAccessDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	case errors.Is(err, service.ErrUsernameRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
	case errors.Is(err, service.ErrEmailRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
	case errors.Is(err, service.ErrPasswordRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
	case errors.Is(err, service.ErrInvalidEmailFormat):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
	case errors.Is(err, service.ErrInvalidPassword):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

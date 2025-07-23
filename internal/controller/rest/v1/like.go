package v1

import (
	"mini-blog/internal/controller"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *RestHandlerv1) CreateLikePost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "CreateLikePost").Logger()
	userID := c.GetInt("user_id")

	in := struct {
		PostID int `json:"post_id"`
	}{}
	err := c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.CreateLikePost(c, in.PostID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to toggle like post")
		controller.HandleError(c, err)
		return
	}

	// Возвращаем разные статусы в зависимости от операции
	if result.IsLiked {
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *RestHandlerv1) CreateLikeComment(c *gin.Context) {
	logger := h.Logger.With().Str("method", "CreateLikeComment").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert comment_id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.CreateLikeComment(c, commentID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to toggle like comment")
		controller.HandleError(c, err)
		return
	}

	statusCode := http.StatusCreated
	if result.Message == "like removed" {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, gin.H{"result": result})
}

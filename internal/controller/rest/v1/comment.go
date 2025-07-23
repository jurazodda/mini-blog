package v1

import (
	"mini-blog/internal/controller"
	"mini-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *RestHandlerv1) CreateComment(c *gin.Context) {
	logger := h.Logger.With().Str("method", "CreateComment").Logger()
	userID := c.GetInt("user_id")

	var in service.CreateCommentIn
	err := c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCommment, err := h.Service.CreateComment(c, in, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create comment")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": createdCommment})
}

func (h *RestHandlerv1) GetCommentByID(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetCommentByID").Logger()
	userID := c.GetInt("user_id")

	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.Service.GetCommentByID(c, userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment by id")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": comment})
}

func (h *RestHandlerv1) UpdateComment(c *gin.Context) {
	logger := h.Logger.With().Str("method", "UpdateComment").Logger()
	userID := c.GetInt("user_id")

	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	in := struct {
		Comment string `json:"comment"`
	}{}
	err = c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedComment, err := h.Service.UpdateComment(c, in.Comment, userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": updatedComment})
}

func (h *RestHandlerv1) DeleteComment(c *gin.Context) {
	logger := h.Logger.With().Str("method", "DeleteComment").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeleteComment(c, userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

func (h *RestHandlerv1) GetCommentsByPostID(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetCommentsByPostID").Logger()

	idStr := c.Param("post_id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert post_id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comments, err := h.Service.GetCommentsByPostID(c, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by post id")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

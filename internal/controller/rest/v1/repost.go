package v1

import (
	"mini-blog/internal/controller"
	"mini-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *RestHandlerv1) CreateRepost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "CreateRepost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	var in service.CreateRepostIn
	err := c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdRepost, err := h.Service.CreateRepost(c, in, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create repost")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"repost": createdRepost})
}

func (h *RestHandlerv1) GetReposts(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetReposts").Logger()

	reposts, err := h.Service.GetReposts(c)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reposts": reposts})
}

func (h *RestHandlerv1) GetRepostByID(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetRepostByID").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	repostID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repost, err := h.Service.GetRepostByID(c, userID, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get repost by id")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"repost": repost})
}

/*
func (h *RestHandlerv1) UpdateRepost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "UpdateRepost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	repostID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var in service.UpdateRepostIn
	err = c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRepost, err := h.Service.UpdateRepost(c, in, userID, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update repost")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"repost": updatedRepost})
}
*/

func (h *RestHandlerv1) DeleteRepost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "DeleteRepost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	repostID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeleteRepost(c, userID, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete repost")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "repost deleted"})
}

func (h *RestHandlerv1) GetRepostsByUserID(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetRepostsByUserID").Logger()

	idStr := c.Param("user_id")
	targetUserID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert user_id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reposts, err := h.Service.GetRepostsByUserID(c, targetUserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts by user id")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"reposts": reposts})
}

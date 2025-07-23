package v1

import (
	"mini-blog/internal/controller"
	"mini-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *RestHandlerv1) CreatePost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "CreatePost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	var in service.CreatePostIn
	err := c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPost, err := h.Service.CreatePost(c, in, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create post")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": createdPost})
}

func (h *RestHandlerv1) GetPosts(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetPosts").Logger()
	userID := c.GetInt(controller.UserIdKey)

	postList, err := h.Service.GetPosts(c, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post_list": postList})
}

func (h *RestHandlerv1) SearchPosts(c *gin.Context) {
	logger := h.Logger.With().Str("method", "SearchPosts").Logger()
	userID := c.GetInt(controller.UserIdKey)

	var searchQuery service.SearchQuery
	err := c.ShouldBindQuery(&searchQuery)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.SearchPosts(c, userID, searchQuery)
	if err != nil {
		logger.Error().Err(err).Msg("failed to search posts")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RestHandlerv1) GetPostByID(c *gin.Context) {
	logger := h.Logger.With().Str("method", "GetPostByID").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.Service.GetPostByID(c, userID, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get post by id")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *RestHandlerv1) UpdatePost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "UpdatePost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var in service.UpdatePostIn
	err = c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost, err := h.Service.UpdatePost(c, in, userID, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update post")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": updatedPost})
}

func (h *RestHandlerv1) DeletePost(c *gin.Context) {
	logger := h.Logger.With().Str("method", "DeletePost").Logger()
	userID := c.GetInt(controller.UserIdKey)

	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error().Err(err).Msg("failed to convert id to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeletePost(c, userID, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete post")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}


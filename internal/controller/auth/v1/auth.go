package v1

import (
	"mini-blog/internal/controller"
	"mini-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandlerV1) SignUp(c *gin.Context) {
	var in service.SignUpIn

	err := c.BindJSON(&in)
	if err != nil {
		h.Logger.Error().Str("method", "SignUp").Err(err).Msg("failed bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.SignUp(c, in)
	if err != nil {
		h.Logger.Error().Str("method", "SignUp").Err(err).Msg("failed signup user")
		controller.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *AuthHandlerV1) SignIn(c *gin.Context) {
	var in service.SignInIn
	err := c.BindJSON(&in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.SignIn(c, in)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed signin user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := service.GenerateToken(c, user.ID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed generate token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

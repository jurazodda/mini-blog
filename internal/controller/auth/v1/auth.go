package v1

import (
	"mini-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandlerV1) SignUp(c *gin.Context) {
	in := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed signup"})
		return
	}

	err = h.Service.SignUp(c, service.SignUpIn{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *AuthHandlerV1) SignIn(c *gin.Context) {
	logger :=  h.Logger.With().Str("method", "SignIn").Logger()

	in := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(&in)
	if err != nil {
		logger.Error().Err(err).Msg("failed bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.SignIn(c, service.SignInIn{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed signin user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
	}

	token, err := h.GenerateToken(c, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed generate token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

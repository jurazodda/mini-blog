package v1

import (
	"mini-blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthHandlerV1 struct {
	Service service.Service
	Logger  zerolog.Logger
}

func (h *AuthHandlerV1) InitRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())

	user := r.Group("/auth/v1")
	{
		user.POST("/sign-up", h.SignUp)
		user.POST("/sign-in", h.SignIn)
	}
}

package v1

import (
	"mini-blog/internal/service"
	"mini-blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthHandlerV1 struct {
	Service service.ServiceI
	Logger  *logger.Logger
}

func (h *AuthHandlerV1) InitRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())

	user := r.Group("/auth/v1")
	{
		user.POST("/sign-up", h.SignUp)
		user.POST("/sign-in", h.SignIn)
	}
}

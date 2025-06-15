package v1

import (
	"mini-blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type RestHandlerv1 struct {
	Service service.Service
	Logger  zerolog.Logger
}

func (h *RestHandlerv1) InitRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())

	// routeGroup
}

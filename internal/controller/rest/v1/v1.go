package v1

import (
	"mini-blog/internal/controller"
	"mini-blog/internal/service"
	"mini-blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RestHandlerv1 struct {
	Service service.ServiceI
	Logger  *logger.Logger
}

func (h *RestHandlerv1) InitRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(controller.RequestLoggingMiddleware(h.Logger))
	// Upload route for images

	post := r.Group("/post", controller.AuthenticateUser)
	{
		post.POST("", h.CreatePost)
		post.GET("", h.GetPosts)
		post.GET("/search", h.SearchPosts)
		post.GET("/:id", h.GetPostByID)
		post.PUT("/:id", h.UpdatePost)
		post.DELETE("/:id", h.DeletePost)
	}

	comment := r.Group("/comment", controller.AuthenticateUser)
	{
		comment.POST("", h.CreateComment)
		comment.GET("/:id", h.GetCommentByID)
		comment.GET("/post/:post_id", h.GetCommentsByPostID)
		comment.PUT("/:id", h.UpdateComment)
		comment.DELETE("/:id", h.DeleteComment)
	}
	
	like := r.Group("/like", controller.AuthenticateUser)
	{
		like.POST("/post", h.CreateLikePost)
		like.POST("/comment/:comment_id", h.CreateLikeComment)
	}

	repost := r.Group("/repost", controller.AuthenticateUser)
	{
		repost.POST("", h.CreateRepost)
		repost.GET("", h.GetReposts)
		repost.GET("/:id", h.GetRepostByID)
		repost.GET("/user/:user_id", h.GetRepostsByUserID)
		// repost.PUT("/:id", h.UpdateRepost)
		repost.DELETE("/:id", h.DeleteRepost)
	}

}

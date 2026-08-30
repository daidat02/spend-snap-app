package post

import (
	"spendsnap-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPostRoutes(router *gin.Engine, createPostHandler *CreatePostHandler) {
	group := router.Group("/api/v1")
	protected := group.Group("")
	protected.Use(middleware.VerifyToken())

	{
		protected.POST("/posts", createPostHandler.CreatePost)
	}

}

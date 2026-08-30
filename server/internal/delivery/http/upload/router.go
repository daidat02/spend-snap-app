package upload

import (
	"spendsnap-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(router *gin.Engine, uploadHandler *UploadHandler) {
	group := router.Group("/api/v1")
	protected := group.Group("")
	protected.Use(middleware.VerifyToken())
	{
		protected.POST("/upload/:filename", uploadHandler.UploadFile)
	}
}

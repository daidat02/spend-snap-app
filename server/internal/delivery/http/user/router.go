package user

import (
	"spendsnap-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h *UserHandler) {
	group := r.Group("/api/v1")
	group.POST("/users/register", h.RegisterAccount)
	group.POST("/users/login", h.Login)

	protected := group.Group("")
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/users/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}	
package transaction

import (
	"spendsnap-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(router *gin.Engine, createHandler *CreateTransactionHandler) {
	transactionGroup := router.Group("/transactions")
	transactionGroup.Use(middleware.VerifyToken())
	{
		transactionGroup.POST("/", createHandler.CreateTransaction)
	}
}

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"spendsnap-backend/internal/config"
	postHttp "spendsnap-backend/internal/delivery/http/post"
	"spendsnap-backend/internal/delivery/http/transaction"
	uploadHttp "spendsnap-backend/internal/delivery/http/upload"
	userHttp "spendsnap-backend/internal/delivery/http/user"
	r2storage "spendsnap-backend/pkg/storage"
)

func newRouter(cfg *config.Config, dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	allowedOrigins := []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002"}

	if cfg.App.Env == "production" {
		allowedOrigins = []string{"https://spendsnap.com", "https://www.spendsnap.com"}
	}
	// cấu hình CORS để cho phép các yêu cầu từ các nguồn được chỉ định
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}));

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"}) 
	})

	// DI qua Dependencies per-module — wire 
	userDeps := userHttp.NewDependencies(dbPool)
	userHttp.RegisterUserRoutes(r, userDeps.Handler)

	storageProvider, err := r2storage.NewR2Provider(cfg.R2)
	if err != nil {
		panic(err)
	}
	uploadDeps := uploadHttp.NewDependencies(storageProvider)
	uploadHttp.RegisterUploadRoutes(r, uploadDeps.Handler)

	transactionDeps := transaction.NewDependencies(dbPool)
	transaction.RegisterTransactionRoutes(r, transactionDeps.Handler)
	
	postDeps := postHttp.NewDependencies(dbPool, uploadDeps.Usecase, transactionDeps.Usecase)
	postHttp.RegisterPostRoutes(r, postDeps.Handler)
	return r
}
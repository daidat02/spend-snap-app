package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	// Khai báo alias gọn gàng, không trùng lặp
	"spendsnap-backend/internal/config"
	uploadHttp "spendsnap-backend/internal/delivery/http/upload"
	userHttp "spendsnap-backend/internal/delivery/http/user"
	userRepo "spendsnap-backend/internal/repository/postgres"
	uploadUsecase "spendsnap-backend/internal/usecase/upload"
	userUsecase "spendsnap-backend/internal/usecase/user"
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
		c.JSON(200, gin.H{"status": "ok"}) // Bỏ dấu chấm phẩy
	})

	// Khởi tạo các dependency với tên biến ngắn gọn, tránh trùng lặp alias
	repo := userRepo.NewPostgresUserRepository(dbPool)
	uc := userUsecase.NewCreateUserUsecase(repo)

	handler := userHttp.NewUserHandler(uc)	
	// Dùng chung userHttp cho cả Handler và Router
	userHttp.RegisterUserRoutes(r, handler)

	storage, err := r2storage.NewR2Provider(cfg.R2)
	if err != nil {
		panic(err)
	}
	uploadUsecase := uploadUsecase.NewCreateUploadUsecase(storage)
	uploadHandler := uploadHttp.NewUploadHandler(uploadUsecase)

	// Dùng chung uploadHttp cho cả Handler và Router
	uploadHttp.RegisterUploadRoutes(r, uploadHandler)
	return r
}
package main

import (
	"fmt"
	"log"
	"spendsnap-backend/internal/config"
	"spendsnap-backend/pkg/database"
	"spendsnap-backend/pkg/utils"
)

func main() {
	// Load Config 
	cfg,err := config.LoadConfig(); 
	fmt.Printf("🌱 SpendSnap Server đang khởi động [Môi trường: %s]...\n", cfg.App.Env)
	if err != nil{
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	dsn := cfg.Database.DB_url

	if dsn == ""{
		dsn = "postgres://postgres:secretpassword@localhost:5432/spendsnap_db?sslmode=disable"
	}
	dbPool := database.NewPostgresDatabase(dsn)

	// Khởi tạo JWT với secret và thời gian hết hạn từ cấu hình
	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.AccessTokenExpire, cfg.JWT.RefreshTokenExpire)

	// Khởi tạo router và các route
	r := newRouter(cfg, dbPool)

	// Server sẽ lắng nghe trên cổng được chỉ định trong cấu hình
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
	log.Printf("SpendSnap server Đang Chạy trên cổng :%s", cfg.App.Port)

}

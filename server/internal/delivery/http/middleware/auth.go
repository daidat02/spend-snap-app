package middleware

import (
	"net/http"
	"spendsnap-backend/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy chuỗi Authorization từ Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Yêu cầu cung cấp Authorization Header",
			})
			c.Abort() // Dừng request ngay lập tức, không cho đi tiếp vào Handler
			return
		}

		// 2. Kiểm tra định dạng "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Định dạng token không hợp lệ (phải là: Bearer <token>)",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Giải mã & Validate Token bằng hàm utils.ValidateToken đã viết
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token không hợp lệ hoặc đã hết hạn",
			})
			c.Abort()
			return
		}

		// 4. Lưu thông tin User vào gin.Context để các Handler đằng sau lấy dùng
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		// 5. Cho phép request đi tiếp vào Handler
		c.Next()
	}
}

// VertifyToken giữ backward-compat với code cũ (sai chính tả)
func VertifyToken() gin.HandlerFunc { return VerifyToken() }
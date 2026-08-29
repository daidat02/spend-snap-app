package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userdomain "spendsnap-backend/internal/domain/user"
	usecase "spendsnap-backend/internal/usecase/user"
)

// DTO (Data Transfer Object) để nhận dữ liệu từ yêu cầu tạo người dùng



type UserHandler struct{
	usecase *usecase.CreateUserUsecase
}

func NewUserHandler(usecase *usecase.CreateUserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}


func (h *UserHandler) RegisterAccount(c *gin.Context){
	var req userdomain.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &userdomain.User{
		Id: req.Id,
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		Status: req.Status,
	}

	userCreated, err := h.usecase.RegisterAccount(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Đăng Ký thành công", "data": userCreated})
}		


func (h *UserHandler) Login(c *gin.Context){
	var req userdomain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, accessToken, err := h.usecase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công", "data": userResponse, "access_token": accessToken})
}
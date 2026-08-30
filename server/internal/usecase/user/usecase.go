package user

import (
	"context"
	"errors"
	"log"
	"time"

	// import domain user
	userdomain "spendsnap-backend/internal/domain/user"
	utils "spendsnap-backend/pkg/utils"
)
type CreateUserUsecase struct{
	userRepo userdomain.UserRepository
}

func NewCreateUserUsecase(userRepo userdomain.UserRepository) *CreateUserUsecase{
	return &CreateUserUsecase{
		userRepo:userRepo,
	}
}

func (uc *CreateUserUsecase) RegisterAccount(ctx context.Context, user *userdomain.User ) (*userdomain.UserResponse, error){
	
	
	// Khởi tạo các giá trị mặc định bắt buộc trước khi Validate
	user.Status = userdomain.StatusActive
	
	// Validate kiểm tra tính hợp lệ của thực thể người dùng
	if err := user.Validate();err != nil{
		log.Println("lỗi dữ liệu user không hợp lệ:", err)
		return nil, err
	}

	hashPassword,err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("lỗi khi hash password:", err)
		return nil, err
	}
	user.Password = hashPassword

	if user.ID == "" {
		user.ID = utils.NewID()
	}

	userCreated, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		log.Println("lỗi khi tạo người dùng:", err)
		return nil, err
	}

	response := &userdomain.UserResponse{
		ID:          userCreated.ID,
		Firstname:   userCreated.Firstname,
		Lastname:    userCreated.Lastname,
		Email:       userCreated.Email,
		Username:    userCreated.Username,
		PhoneNumber: userCreated.PhoneNumber,
		Status:      userCreated.Status,
	}
	return response, nil
}

func (uc *CreateUserUsecase) Login(ctx context.Context, req *userdomain.LoginRequest) (*userdomain.UserResponse, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
    // 1. Validate kiểm tra tính hợp lệ của dữ liệu đầu vào trước
	log.Println("Bắt đầu xử lý đăng nhập cho email:", req.Email)
    if err := req.Validate(); err != nil {
        log.Println("lỗi dữ liệu đăng nhập không hợp lệ:", err)
        return nil, "", err
    }

    // 2. Tìm user trong DB theo Email
    userFound, err := uc.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        log.Println("lỗi khi truy vấn user theo email:", err)
        return nil, "", err
    }
    if userFound == nil {
        return nil, "", errors.New("email hoặc mật khẩu không chính xác")
    }

    // 3. Kiểm tra trạng thái tài khoản (chỉ cho phép user đang Active đăng nhập)
    if userFound.Status != userdomain.StatusActive {
        return nil, "", errors.New("tài khoản của bạn đã bị khóa hoặc chưa kích hoạt")
    }

    // 4. So sánh password hash
    if !utils.CheckPasswordHash(req.Password, userFound.Password){
		log.Println("Mật khẩu không chính xác cho email:", req.Email)
		return nil, "", errors.New("email hoặc mật khẩu không chính xác")
	}
	accessToken, err := utils.GenerateToken(userFound.ID, userFound.Email)
	if err != nil {
		log.Println("lỗi khi tạo token:", err)
		return nil, "", err
	}

    response := &userdomain.UserResponse{
        ID:          userFound.ID,
        Firstname:   userFound.Firstname,
        Lastname:    userFound.Lastname,
        Email:       userFound.Email,
		Username:    userFound.Username,
        AvatarURL:   userFound.AvatarURL,
        PhoneNumber: userFound.PhoneNumber,
        Bio:         userFound.Bio,
        Status:      userFound.Status,
        CreatedAt:   userFound.CreatedAt,
    }

    return response, accessToken, nil
}
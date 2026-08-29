package user

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)
type CreateUserRequest struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	Status string `json:"status"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserResponse struct {
    // Tên trường trong Go (Viết hoa) | Tag chỉ định key JSON (Viết thường)
    Id          string `json:"id"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    AvatarURL   string `json:"avatar_url,omitempty"`   // Nếu rỗng thì ẩn
    PhoneNumber string `json:"phone_number,omitempty"` // Nếu rỗng thì ẩn
    Status      string `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}



func (req *LoginRequest) Validate() error {
    if strings.TrimSpace(req.Email) == "" {
        return errors.New("lỗi: email không được để trống")
    }
    if _, err := mail.ParseAddress(req.Email); err != nil {
        return errors.New("lỗi: định dạng email không hợp lệ")
    }
    if req.Password == "" {
        return errors.New("lỗi: mật khẩu không được để trống")
    }
    return nil
}
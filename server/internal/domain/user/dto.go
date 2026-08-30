package user

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)
type CreateUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
    Username string `json:"user_name"`
	Firstname string `json:"first_name"`
	Lastname string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	Bio string `json:"bio"`
	Status string `json:"status"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserResponse struct {
    ID          string `json:"id"`
    Email       string `json:"email"`
    Username    string `json:"user_name"` // Nếu rỗng thì ẩn
    Firstname   string `json:"first_name,omitempty"`
    Lastname    string `json:"last_name,omitempty"`
    AvatarURL   *string `json:"avatar_url,omitempty"`   // Nếu rỗng thì ẩn
    PhoneNumber *string `json:"phone_number,omitempty"` // Nếu rỗng thì ẩn
    Bio         *string `json:"bio,omitempty"`
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
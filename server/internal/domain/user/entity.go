package user

import (
	"errors"
	"time"
)

// khai báo các trạng thái hợp lệ
const (
	StatusActive = "active"
	StatusInactive = "inactive"
)

// thực thể người dùng
type User struct{
	Id string
	Email string
	Password string
	Username string
	Firstname string
	Lastname string
	AvatarURL *string
	PhoneNumber *string
	Bio *string
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}



// Validate kiểm tra tính hợp lệ của thực thể người dùng
func (user *User) Validate() error {
	if user.Email == "" {
		return errors.New("Lỗi: Email không được để trống")
	}
	if user.Firstname == "" {
		return errors.New("Lỗi: Tên không được để trống")
	}
	if user.Lastname == "" {
		return errors.New("Lỗi: Họ không được để trống")
	}
	if user.Password == "" {
		return errors.New("Lỗi: Mật khẩu không được để trống")
	}
	if status := user.Status; status != StatusActive && status != StatusInactive {
		return errors.New("Lỗi: Trạng thái không hợp lệ")
	}
	return nil
}



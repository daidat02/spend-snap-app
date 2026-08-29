package user

import "context"

// UserRepository là giao diện định nghĩa các phương thức để thao tác với thực thể người dùng trong cơ sở dữ liệu
type UserRepository interface{
	Create(ctx context.Context, user *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByID(ctx context.Context, id string) (*User, error)
	// Update(ctx context.Context, user *User) error
	// Delete(ctx context.Context, id string) error
}
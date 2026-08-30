package post

import (
	"errors"
	"time"
)


type Post struct{
	ID string
	UserID string
	ImageUrl string
	Caption *string
	Visibility string
	LocationName *string
	TransactionID *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (post *Post) Validate() error {
	if post.UserID == "" {
		return errors.New("Lỗi: UserID không được để trống")
	}
	if post.ImageUrl == "" {
		return errors.New("Lỗi: ImageUrl không được để trống")
	}
	if post.Visibility != "public" && post.Visibility != "private" && post.Visibility != "friends" {
		return errors.New("Lỗi: Visibility không hợp lệ")
	}
	return nil
}

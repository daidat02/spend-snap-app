package transaction

import (
	"errors"
	"time"
)


type Transaction struct {
	ID              string
	UserID          string
	CategoryID      string
	Amount          int64
	Type            string
	Note            *string
	Source          string
	TransactionDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func (t *Transaction) Validate() error {
	if t.UserID == "" {
		return errors.New("Lỗi: UserID không được để trống")
	}
	if t.CategoryID == "" {
		return errors.New("Lỗi: CategoryID không được để trống")
	}
	if t.Amount <= 0 {
		return errors.New("Lỗi: Amount phải lớn hơn 0")
	}
	if t.Type != "income" && t.Type != "expense" {
		return errors.New("Lỗi: Type không hợp lệ")
	}
	if t.Source == "" {
		return errors.New("Lỗi: Source không được để trống")
	}
	if t.TransactionDate.IsZero() {
		return errors.New("Lỗi: TransactionDate không được để trống")
	}
	return nil
}
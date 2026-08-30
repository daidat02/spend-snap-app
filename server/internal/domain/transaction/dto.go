package transaction

import "time"

type CreateTransactionRequest struct {
	CategoryID      string  `json:"category_id" binding:"required"`
	Amount          int64   `json:"amount" binding:"required,gt=0"`
	Type            string  `json:"type" binding:"omitempty,oneof=expense income"`
	Note            *string `json:"note"`
	TransactionDate *string `json:"transaction_date"`
}

type TransactionResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	CategoryID      string    `json:"category_id"`
	Amount          int64     `json:"amount"`
	Type            string    `json:"type"`
	Note            *string   `json:"note,omitempty"`
	Source          string    `json:"source"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

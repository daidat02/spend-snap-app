package transaction

import "context"

import "github.com/jackc/pgx/v5"

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) (*Transaction, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, transaction *Transaction) (*Transaction, error)
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) (*TransactionResponse, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, transaction *Transaction) (*TransactionResponse, error)
}

package postgres

import (
	"context"
	domain "spendsnap-backend/internal/domain/transaction"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTransactionRepository struct {
	dbPool *pgxpool.Pool
}

func NewPostgresTransactionRepository(dbPool *pgxpool.Pool) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		dbPool: dbPool,
	}
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	err := r.dbPool.QueryRow(ctx,
		CreateTransactionQuery,
		transaction.ID,
		transaction.UserID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Type,
		transaction.Note,
		transaction.Source,
		transaction.TransactionDate,
	).Scan(&transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *PostgresTransactionRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, transaction *domain.Transaction) (*domain.Transaction, error) {
	err := tx.QueryRow(ctx,
		CreateTransactionQuery,
		transaction.ID,
		transaction.UserID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Type,
		transaction.Note,
		transaction.Source,
		transaction.TransactionDate,
	).Scan(&transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

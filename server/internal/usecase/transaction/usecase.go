package transaction

import (
	"context"
	"time"

	domain "spendsnap-backend/internal/domain/transaction"
	"spendsnap-backend/pkg/utils"

	"github.com/jackc/pgx/v5"
)

type CreateTransactionUsecase struct {
	repo domain.TransactionRepository
}

func NewCreateTransactionUsecase(repo domain.TransactionRepository) *CreateTransactionUsecase {
	return &CreateTransactionUsecase{
		repo: repo,
	}
}

func (uc *CreateTransactionUsecase) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.TransactionResponse, error) {
	if transaction.ID == "" {
		transaction.ID = utils.NewID()
	}
	if transaction.Type == "" {
		transaction.Type = "expense"
	}
	if transaction.Source == "" {
		transaction.Source = "manual"
	}
	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	createdTransaction, err := uc.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return toResponse(createdTransaction), nil
}

func (uc *CreateTransactionUsecase) CreateWithTx(ctx context.Context, tx pgx.Tx, transaction *domain.Transaction) (*domain.TransactionResponse, error) {
	if transaction.ID == "" {
		transaction.ID = utils.NewID()
	}
	if transaction.Type == "" {
		transaction.Type = "expense"
	}
	if transaction.Source == "" {
		transaction.Source = "manual"
	}
	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	createdTransaction, err := uc.repo.CreateWithTx(ctx, tx, transaction)
	if err != nil {
		return nil, err
	}
	return toResponse(createdTransaction), nil
}

func toResponse(t *domain.Transaction) *domain.TransactionResponse {
	return &domain.TransactionResponse{
		ID:              t.ID,
		UserID:          t.UserID,
		CategoryID:      t.CategoryID,
		Amount:          t.Amount,
		Type:            t.Type,
		Note:            t.Note,
		Source:          t.Source,
		TransactionDate: t.TransactionDate,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

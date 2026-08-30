package transaction

import (
	postgres "spendsnap-backend/internal/repository/postgres"
	usecase "spendsnap-backend/internal/usecase/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionDependency struct {
	Repo   *postgres.PostgresTransactionRepository
	Usecase *usecase.CreateTransactionUsecase
	Handler *CreateTransactionHandler

}

func NewDependencies(dbPool *pgxpool.Pool) *TransactionDependency {
	repo := postgres.NewPostgresTransactionRepository(dbPool)
	uc := usecase.NewCreateTransactionUsecase(repo)
	handler := NewCreateTransactionHandler(uc)

	return &TransactionDependency{
		Repo:    repo,
		Usecase: uc,
		Handler: handler,
	}
}

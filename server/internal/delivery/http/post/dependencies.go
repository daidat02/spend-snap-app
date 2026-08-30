package post

import (
	"spendsnap-backend/internal/domain/storage"
	transactionDomain "spendsnap-backend/internal/domain/transaction"
	"spendsnap-backend/internal/repository/postgres"
	postUsecase "spendsnap-backend/internal/usecase/post"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Repo    *postgres.PostgresPostRepository
	Usecase *postUsecase.CreatePostUsecase
	Handler *CreatePostHandler
}

func NewDependencies(dbPool *pgxpool.Pool, uploadUC storage.UploadUsecase, transactionUC transactionDomain.TransactionUsecase) *Dependencies {
	repo := postgres.NewPostgresPostRepository(dbPool)
	uc := postUsecase.NewCreatePostUsecase(dbPool, repo, uploadUC, transactionUC)
	handler := NewCreatePostHandler(uc)

	return &Dependencies{
		Repo:    repo,
		Usecase: uc,
		Handler: handler,
	}
}

package user

import (
	"spendsnap-backend/internal/repository/postgres"
	"spendsnap-backend/internal/usecase/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Repo    *postgres.PostgresUserRepository
	Usecase *user.CreateUserUsecase
	Handler *UserHandler
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := postgres.NewPostgresUserRepository(dbPool)
	uc := user.NewCreateUserUsecase(repo)
	handler := NewUserHandler(uc)

	return &Dependencies{
		Repo:    repo,
		Usecase: uc,
		Handler: handler,
	}
}

package postgres

import (
	"context"

	domain "spendsnap-backend/internal/domain/post"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPostRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostRepository(db *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{
		db: db,
	}
}

func (repo *PostgresPostRepository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	_, err := repo.db.Exec(
		ctx,
		queryCreatePost,
		post.ID,
		post.UserID,
		post.ImageUrl,
		post.Caption,
		post.Visibility,
		post.LocationName,
		post.TransactionID,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostgresPostRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, post *domain.Post) (*domain.Post, error) {
	_, err := tx.Exec(
		ctx,
		queryCreatePost,
		post.ID,
		post.UserID,
		post.ImageUrl,
		post.Caption,
		post.Visibility,
		post.LocationName,
		post.TransactionID,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

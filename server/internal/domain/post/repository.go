package post

import "context"

import "github.com/jackc/pgx/v5"

type PostRepository interface {
	Create(ctx context.Context, post *Post) (*Post ,error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, post *Post) (*Post ,error)
	// GetByID(id string) (*Post, error)
	// Update(post *Post) error
	// Delete(id string) error
}

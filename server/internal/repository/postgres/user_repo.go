package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	userdomain "spendsnap-backend/internal/domain/user"
)

type PostgresUserRepository struct{
	db *pgxpool.Pool
}


func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository{
	return &PostgresUserRepository{
		db: db,
	}
}


func(repo *PostgresUserRepository)Create(ctx context.Context,user *userdomain.User) (*userdomain.User, error){
	log.Println("Creating user in Postgres repository:", user)
	_, err := repo.db.Exec(
		ctx,
		queryCreateUser,
		user.ID,
		user.Email,
		user.Password,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.PhoneNumber,
		user.Status,
	)
	if err != nil{
		log.Println("Lỗi khi tạo người dùng:", err)
		return nil, err
	}
	return user, nil
}


func (repo *PostgresUserRepository) FindByEmail(ctx context.Context, email string) ( *userdomain.User, error){
	log.Println("Finding user by email in Postgres repository:", email)
	u := &userdomain.User{}
	err := repo.db.QueryRow(ctx, queryGetByEmail, email).Scan(
		&u.ID,
		&u.Email,
		&u.Firstname,
		&u.Lastname,
		&u.Password, // Nhận về chuỗi password_hash đã lưu trong DB
		&u.Username,
		&u.AvatarURL,
		&u.PhoneNumber,
		&u.Bio,
		&u.Status,
		&u.CreatedAt,
	)

	if err != nil{
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Println("Lỗi khi tìm người dùng:", err)
		return nil, err
	}
	return  u, nil
}

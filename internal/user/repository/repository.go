package repository

import (
	"context"
	"database/sql"

	"github.com/ProNinjaDev/GoUserApi/internal/user"
)

type Repository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id int64) (*user.User, error)
	GetByFilter(ctx context.Context, name string, status string) ([]user.User, error)
	Update(ctx context.Context, id int64, u user.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	query := "INSERT INTO users (name, status) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRowContext(ctx, query, u.Name, u.Status).Scan(&u.Id)
}

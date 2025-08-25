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

func (r *userRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	query := "SELECT id, name, status FROM users WHERE id = $1"
	var u user.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.Id, &u.Name, &u.Status)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

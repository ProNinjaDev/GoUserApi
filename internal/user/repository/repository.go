package repository

import (
	"context"
	"database/sql"
	"strconv"

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

func (r *userRepository) GetByFilter(ctx context.Context, name, statusString string) ([]user.User, error) {
	query := "SELECT id, name, status FROM users WHERE 1=1"
	args := []any{}
	argCnt := 1

	if name != "" {
		query += " AND name LIKE $" + strconv.Itoa(argCnt)
		args = append(args, name)
		argCnt++
	}

	if statusString != "" {
		status, err := strconv.ParseBool(statusString)
		if err != nil {
			return nil, err
		}

		query += " AND status = $" + strconv.Itoa(argCnt)
		args = append(args, status)
		argCnt++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User

		if err := rows.Scan(&u.Id, &u.Name, &u.Status); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *userRepository) Update(ctx context.Context, id int64, u user.User) error {
	query := "UPDATE users SET name = $1, status = $2 WHERE id = $3"
	result, err := r.db.ExecContext(ctx, query, u.Name, u.Status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

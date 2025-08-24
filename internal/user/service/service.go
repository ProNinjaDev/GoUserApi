package service

import (
	"context"

	"github.com/ProNinjaDev/GoUserApi/internal/user"
	"github.com/ProNinjaDev/GoUserApi/internal/user/repository"
)

type Service interface {
	Create(ctx context.Context, u *user.User) error
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, u *user.User) error {
	return s.repo.Create(ctx, u)
}

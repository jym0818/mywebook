package service

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewuserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Signup(ctx context.Context, user domain.User) error {
	return u.repo.Create(ctx, user)
}

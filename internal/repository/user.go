package repository

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/dao"
	"time"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type userRepository struct {
	dao dao.UserDAO
}

func NewuserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *userRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
		Utime:    time.UnixMilli(user.Utime),
		CTime:    time.UnixMilli(user.Ctime),
	}

}

func (repo *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
		Utime:    user.Utime.UnixMilli(),
		Ctime:    user.CTime.UnixMilli(),
	}

}

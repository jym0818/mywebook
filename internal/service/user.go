package service

import (
	"context"
	"errors"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (u *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//先去查找
	user, err := u.repo.FindByPhone(ctx, phone)
	//err不为找不到
	if err != repository.ErrUserNotExists {
		return user, err
	}
	user = domain.User{
		Phone: phone,
	}
	err = u.repo.Create(ctx, user)
	if err != repository.ErrUserDuplicateEmail {
		return user, err
	}
	return u.repo.FindByPhone(ctx, phone)

}

func (u *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	//先去查找
	user, err := u.repo.FindByWechat(ctx, info.OpenID)
	//err不为找不到
	if err != repository.ErrUserNotExists {
		return user, err
	}
	user = domain.User{
		WechatInfo: domain.WechatInfo{OpenID: info.OpenID, UnionID: info.UnionID},
	}
	err = u.repo.Create(ctx, user)
	if err != repository.ErrUserDuplicateEmail {
		return user, err
	}
	return u.repo.FindByWechat(ctx, info.OpenID)

}

func (u *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotExists {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func NewuserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Signup(ctx context.Context, user domain.User) error {
	//密码加密
	str, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(str)
	return u.repo.Create(ctx, user)
}

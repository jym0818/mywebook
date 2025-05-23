package repository

import (
	"context"
	"database/sql"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
	"time"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicate
var ErrUserNotExists = dao.ErrUserNotExists

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repo *userRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := repo.dao.FindByWechat(ctx, openID)
	return repo.toDomain(u), err
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	return repo.toDomain(u), err
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	return repo.toDomain(u), err
}
func NewuserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	user, err := repo.cache.Get(ctx, id)
	if err == nil {
		return user, nil
	}
	//找不到或者系统错误
	ue, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	//回写
	u := repo.toDomain(ue)
	err = repo.cache.Set(ctx, u)
	if err != nil {
		//记录日志就可以
	}
	return u, nil
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *userRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Utime:    time.UnixMilli(user.Utime),
		CTime:    time.UnixMilli(user.Ctime),
		WechatInfo: domain.WechatInfo{
			OpenID:  user.WechatOpenID.String,
			UnionID: user.WechatUnionID.String,
		},
	}

}

func (repo *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:            user.Id,
		Email:         sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:         sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password:      user.Password,
		Utime:         user.Utime.UnixMilli(),
		Ctime:         user.CTime.UnixMilli(),
		WechatOpenID:  sql.NullString{String: user.WechatInfo.OpenID, Valid: user.WechatInfo.OpenID != ""},
		WechatUnionID: sql.NullString{String: user.WechatInfo.UnionID, Valid: user.WechatInfo.UnionID != ""},
	}

}

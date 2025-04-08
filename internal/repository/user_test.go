package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/cache"
	cachemocks "github.com/jym/mywebook/internal/repository/cache/mocks"
	"github.com/jym/mywebook/internal/repository/dao"
	daomocks "github.com/jym/mywebook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		id int64

		wantRes domain.User
		wantErr error
	}{
		{
			name: "回写",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {

				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExists)

				userDAO := daomocks.NewMockUserDAO(ctrl)
				userDAO.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "1139022499@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Password: "jy123456@",
					Ctime:    now.UnixMilli(),
					Utime:    now.UnixMilli(),
				}, nil)

				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "1139022499@qq.com",
					Password: "jy123456@",
					CTime:    now,
					Phone:    "15212345678",
					Utime:    now,
				}).Return(nil)

				return userDAO, userCache
			},
			id: 123,
			wantRes: domain.User{

				Id:       123,
				Email:    "1139022499@qq.com",
				Password: "jy123456@",
				CTime:    now,
				Phone:    "15212345678",
				Utime:    now,
			},
			wantErr: nil,
		},
		{
			name: "查找成功缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {

				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "1139022499@qq.com",
						Password: "jy123456@",
						CTime:    now,
						Phone:    "15212345678",
						Utime:    now,
					}, nil)

				userDAO := daomocks.NewMockUserDAO(ctrl)

				return userDAO, userCache
			},
			id: 123,
			wantRes: domain.User{

				Id:       123,
				Email:    "1139022499@qq.com",
				Password: "jy123456@",
				CTime:    now,
				Phone:    "15212345678",
				Utime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {

				userCache := cachemocks.NewMockUserCache(ctrl)
				userCache.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExists)

				userDAO := daomocks.NewMockUserDAO(ctrl)
				userDAO.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, errors.New("系统错误"))
				
				return userDAO, userCache
			},
			id:      123,
			wantRes: domain.User{},
			wantErr: errors.New("系统错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d, c := tc.mock(ctrl)
			repo := NewuserRepository(d, c)
			res, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

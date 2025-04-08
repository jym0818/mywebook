package service

import (
	"context"
	"errors"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
	repomocks "github.com/jym/mywebook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		email    string
		password string
		wantErr  error
		wantRes  domain.User
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1139022499@qq.com").Return(domain.User{
					Id:       1,
					Email:    "1139022499@qq.com",
					Password: "$2a$10$99cKH8rOxvb77wzUZHUM6em2YnLkgeNURNcrqZsVJoD3fLV.hWoiG",
					CTime:    now,
				}, nil)
				return repo
			},
			email:    "1139022499@qq.com",
			password: "jy123456@",
			wantRes: domain.User{
				Id:       1,
				Email:    "1139022499@qq.com",
				Password: "$2a$10$99cKH8rOxvb77wzUZHUM6em2YnLkgeNURNcrqZsVJoD3fLV.hWoiG",
				CTime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1139022499@qq.com").Return(domain.User{}, repository.ErrUserNotExists)
				return repo
			},
			email:    "1139022499@qq.com",
			password: "jy123456@",
			wantRes:  domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1139022499@qq.com").Return(domain.User{
					Id:       1,
					Email:    "1139022499@qq.com",
					Password: "$2a$10$99cKH8rOxvb77wzUZHUM6em2YnLkgeNURNcrqZsVJoD3fLV.hWoiG",
					CTime:    now,
				}, nil)
				return repo
			},
			email:    "1139022499@qq.com",
			password: "jy1234561@",
			wantRes:  domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1139022499@qq.com").Return(domain.User{}, errors.New("系统错误"))
				return repo
			},
			email:    "1139022499@qq.com",
			password: "jy123456@",
			wantRes:  domain.User{},
			wantErr:  errors.New("系统错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := NewuserService(tc.mock(ctrl))
			res, err := svc.Login(context.Background(), tc.email, tc.password)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "jy123456@"
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	t.Log(string(p))
}

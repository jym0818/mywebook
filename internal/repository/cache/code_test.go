package cache

import (
	"context"
	"errors"
	"github.com/jym/mywebook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15904922108"}, "123456").
					Return(res)
				return cmd
			},
			phone:   "15904922108",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "发送验证码太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15904922108"}, "123456").
					Return(res)
				return cmd
			},
			phone:   "15904922108",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "发送验证码太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-2))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15904922108"}, "123456").
					Return(res)
				return cmd
			},
			phone:   "15904922108",
			code:    "123456",
			wantErr: ErrUnknownForCode,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis err"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:15904922108"}, "123456").
					Return(res)
				return cmd
			},
			phone:   "15904922108",
			code:    "123456",
			wantErr: errors.New("mock redis err"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := NewRedisCodeCache(tc.mock(ctrl))
			err := uc.Set(context.Background(), "login", tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

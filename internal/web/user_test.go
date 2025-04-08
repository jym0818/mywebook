package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/service"
	svcmocks "github.com/jym/mywebook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody Result
	}{
		{
			name: "成功注册",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1139022499@qq.com",
					Password: "jy123456@",
				}).Return(nil)
				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456@"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "注册成功",
			},
		},
		{
			name: "参数绑定错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456@"，
}
`,
			wantCode: http.StatusBadRequest,
			wantBody: Result{
				Msg: "系统错误",
			},
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			reqBody: `
{
	"email":"1139022499qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456@"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "邮箱格式错误",
			},
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456",
	"confirmPassword":"jy123456"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "密码格式错误",
			},
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456!"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "两次密码不一致",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1139022499@qq.com",
					Password: "jy123456@",
				}).Return(service.ErrUserDuplicateEmail)
				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456@"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "邮箱已注册",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				svc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "1139022499@qq.com",
					Password: "jy123456@",
				}).Return(errors.New("系统错误"))
				return svc
			},
			reqBody: `
{
	"email":"1139022499@qq.com",
	"password":"jy123456@",
	"confirmPassword":"jy123456@"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Msg: "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			user := NewUserHandler(tc.mock(ctrl), nil)
			s := gin.Default()
			user.RegisterRouters(s)

			res, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			assert.NoError(t, err)
			res.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.ServeHTTP(resp, res)
			var result Result
			err = json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, result)
		})
	}
}

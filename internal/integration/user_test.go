package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jym/mywebook/internal/web"
	"github.com/jym/mywebook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_SendCode(t *testing.T) {
	s := InitWeb()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name    string
		resBody string

		wantCode int
		wantBody web.Result

		after  func(*testing.T)
		before func(*testing.T)
	}{
		{
			name:   "发送成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				val, err := rdb.GetDel(ctx, "phone_code:login:15904922108").Result()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			resBody: `
{
	"phone":"15904922108"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				err := rdb.Set(ctx, "phone_code:login:15904922108", "123456", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				val, err := rdb.GetDel(ctx, "phone_code:login:15904922108").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是6位
				assert.Equal(t, val, "123456")
			},
			resBody: `
{
	"phone":"15904922108"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			res, err := http.NewRequest(http.MethodPost, "/user/sms/send_code", bytes.NewBuffer([]byte(tc.resBody)))
			require.NoError(t, err)
			res.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			s.ServeHTTP(resp, res)
			var result web.Result
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, result)
			tc.after(t)
		})
	}
}

func TestRedisDB(t *testing.T) {
	rdb := ioc.InitRedis()
	rdb.FlushAll(context.Background())
}

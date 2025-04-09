package integration

import (
	"context"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCodeCache_Set(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr: "47.120.51.5:6379",
	})
	rd := cache.NewRedisCodeCache(cmd).(*cache.RedisCodeCache)
	testCases := []struct {
		name   string
		after  func(t *testing.T)
		before func(t *testing.T)

		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name:   "发送成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := cmd.Get(ctx, "phone_code:login:15904922108").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "123465")
				ttl, err := cmd.TTL(ctx, "phone_code:login:15904922108").Result()
				assert.NoError(t, err)
				assert.True(t, ttl > time.Minute*9+time.Second*30)
				err = cmd.Del(ctx, "phone_code:login:15904922108").Err()
				assert.NoError(t, err)
			},
			biz:   "login",
			phone: "15904922108",
			code:  "123465",
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				err := cmd.Set(ctx, "phone_code:login:15904922108", "123456", time.Minute*10).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := cmd.Get(ctx, "phone_code:login:15904922108").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "123456")
				ttl, err := cmd.TTL(ctx, "phone_code:login:15904922108").Result()
				assert.NoError(t, err)
				assert.True(t, ttl > time.Minute*9)
				err = cmd.Del(ctx, "phone_code:login:15904922108").Err()
				assert.NoError(t, err)
			},
			biz:     "login",
			phone:   "15904922108",
			code:    "78910",
			wantErr: cache.ErrCodeSendTooMany,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := rd.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}

func TestRedisCodeCache_Get(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr: "47.120.51.5:6379",
	})
	cmd.FlushAll(context.Background())
}

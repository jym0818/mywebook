package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var LuaSlideWindow string

type RedisSlideWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	rate     int
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}
func (l *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	b, err := l.cmd.Eval(ctx, LuaSlideWindow, []string{key}, l.interval.Milliseconds(), l.rate, time.Now().UnixMilli()).Bool()
	fmt.Println(b)
	fmt.Println(err)
	return b, err
}

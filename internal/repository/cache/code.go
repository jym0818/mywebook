package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
type RedisCodeCache struct {
	cmd redis.Cmdable
}

func (r *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := r.cmd.Eval(ctx, luaVerifyCode, []string{r.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooManyTimes
	default:
		return false, nil
	}
}

func (r *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := r.cmd.Eval(ctx, luaSetCode, []string{r.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -1:
		return ErrCodeSendTooMany
	case -2:
		return ErrUnknownForCode
	default:
		return nil
	}

}
func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

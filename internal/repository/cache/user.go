package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/mywebook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExists = redis.Nil

type UserCache interface {
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, uid int64) (domain.User, error)
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (r *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	data, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, r.key(u.Id), data, r.expiration).Err()
}

func (r *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (r *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	data, err := r.cmd.Get(ctx, r.key(uid)).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	return user, err
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 10,
	}
}

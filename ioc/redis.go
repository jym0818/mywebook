package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "47.120.51.5:6379",
		Password: "",
		DB:       1,
	})
	return cmd
}

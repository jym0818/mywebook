package startup

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var redisClient redis.Cmdable

func InitRedis() redis.Cmdable {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "47.120.51.5:6379",
		})

		for err := redisClient.Ping(context.Background()).Err(); err != nil; {
			panic(err)
		}
	}
	return redisClient
}

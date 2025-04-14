package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string
	}
	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	cmd := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: "",
		DB:       1,
	})
	return cmd
}

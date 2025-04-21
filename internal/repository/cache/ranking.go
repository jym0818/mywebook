package cache

import (
	"context"
	"encoding/json"
	"github.com/jym/mywebook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	cmd redis.Cmdable
	key string
}

func (cache *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	//节省内存
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	val, err := json.Marshal(&arts)
	if err != nil {
		return err
	}
	//这个过期时间要大于计算时间
	//可以考虑永不过期
	return cache.cmd.Set(ctx, cache.key, val, time.Hour*24).Err()
}

func (cache *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := cache.cmd.Get(ctx, cache.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(data, &res)
	return res, err
}

func NewRedisRankingCache(cmd redis.Cmdable) RankingCache {
	return &RedisRankingCache{
		cmd: cmd,
		key: "ranking",
	}

}

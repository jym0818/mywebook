package repository

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type rankingRepository struct {
	cache cache.RankingCache
	local *cache.RankingLocalCache
}

func (repo *rankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := repo.local.Get(ctx)
	if err == nil {
		return arts, nil
	}
	arts, err = repo.cache.Get(ctx)
	if err == nil {
		_ = repo.local.Set(ctx, arts)
	} else {
		//redis出现问题了 ，那我们直接从本地缓存读取，不去考虑它是否过期
		arts, err = repo.local.ForceGet(ctx)
	}
	return arts, err
}

func (repo *rankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	//1.先操作本地缓存 因为本地缓存几乎不可能失败
	_ = repo.local.Set(ctx, arts)
	return repo.cache.Set(ctx, arts)

}

func NewrankingRepository(cache cache.RankingCache, local cache.RankingLocalCache) RankingRepository {
	return &rankingRepository{
		cache: cache,
		local: &local,
	}
}

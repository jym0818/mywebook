package repository

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

type rankingRepository struct {
	cache cache.RankingCache
}

func (repo *rankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return repo.cache.Set(ctx, arts)
}

func NewrankingRepository(cache cache.RankingCache) RankingRepository {
	return &rankingRepository{
		cache: cache,
	}
}

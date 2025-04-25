package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/jym/mywebook/internal/domain"
	"time"
)

type RankingLocalCache struct {
	//可以考虑用uber的缓存
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		topN: atomicx.NewValue[[]domain.Article](),
		ddl:  atomicx.NewValueOf(time.Now()),
		//永不过期 或者非常长
		expiration: time.Minute * 3,
	}
}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.topN.Store(arts)
	ddl := time.Now().Add(r.expiration)
	r.ddl.Store(ddl)
	return nil
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := r.ddl.Load()
	arts := r.topN.Load()
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("已经过期了")
	}
	return arts, nil
}

func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	return arts, nil
}

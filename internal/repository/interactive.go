package repository

import (
	"context"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (repo *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {
	err := repo.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}
	//写入缓存
	return repo.cache.IncrReadCntIfPresent(ctx, biz, id)
}

func NewinteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
	}
}

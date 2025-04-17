package repository

import (
	"context"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (repo *interactiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func (repo *interactiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := repo.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	//缓存
	return repo.cache.IncrLikeCntIfPresent(ctx, biz, id)
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

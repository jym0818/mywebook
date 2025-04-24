package repository

import (
	"context"
	"github.com/jym/mywebook/interactive/domain"
	"github.com/jym/mywebook/interactive/repository/cache"
	"github.com/jym/mywebook/interactive/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (repo *interactiveRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	return repo.dao.BatchIncrReadCnt(ctx, bizs, ids)
}

func (repo *interactiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	//先去缓存中查找
	intr, err := repo.cache.Get(ctx, biz, id)
	if err == nil {
		return intr, nil
	}
	//去数据库查找
	ie, err := repo.dao.Get(ctx, biz, id)
	if err == nil {
		res := domain.Interactive{
			LikeCnt:    ie.LikeCnt,
			ReadCnt:    ie.ReadCnt,
			CollectCnt: ie.CollectCnt,
		}
		if er := repo.cache.Set(ctx, biz, id, res); er != nil {
			//缓存设置失败，记录日志
		}
		return res, nil
	}
	return domain.Interactive{}, err
}

func (repo *interactiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	err := repo.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		Cid:   cid,
		BizId: id,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return repo.cache.IncrCollectCntIfPresent(ctx, biz, id)
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

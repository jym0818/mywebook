package service

import (
	"context"
	"github.com/jym/mywebook/interactive/domain"
	"github.com/jym/mywebook/interactive/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
}
type interactiveService struct {
	repo repository.InteractiveRepository
}

func (svc *interactiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *interactiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	intr, err := svc.repo.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	if uid > 0 {
		//说明登录了，需要判断是否点赞
		var eg errgroup.Group
		eg.Go(func() error {
			intr.Liked, err = svc.repo.Liked(ctx, biz, id, uid)
			return err
		})
		eg.Go(func() error {
			intr.Collected, err = svc.repo.Collected(ctx, biz, id, uid)
			return err
		})

		err := eg.Wait()
		if err != nil {
			//记录日志就可以
		}

	}
	return intr, nil
}

func (svc *interactiveService) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return svc.repo.AddCollectionItem(ctx, biz, id, cid, uid)
}

func NewinteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (svc *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.IncrLike(ctx, biz, id, uid)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.DecrLike(ctx, biz, id, uid)
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

package service

import (
	"context"
	"github.com/jym/mywebook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
}
type interactiveService struct {
	repo repository.InteractiveRepository
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

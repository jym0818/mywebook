package service

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewarticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (svc *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	return svc.repo.Save(ctx, article)
}

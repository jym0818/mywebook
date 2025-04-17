package service

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(c context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func (svc *articleService) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetPublishedById(ctx, id)
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *articleService) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	return svc.repo.List(ctx, uid, limit, offset)
}

func (svc *articleService) Withdraw(c context.Context, article domain.Article) error {
	return svc.repo.SyncStatus(c, article.Id, article.Author.Id, domain.ArticleStatusPrivate)
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func NewarticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (svc *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnPublished
	if article.Id > 0 {
		err := svc.repo.Update(ctx, article)
		return article.Id, err
	}
	return svc.repo.Create(ctx, article)
}

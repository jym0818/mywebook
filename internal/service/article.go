package service

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/events/article"
	"github.com/jym/mywebook/internal/repository"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(c context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64, uid int64) (domain.Article, error)

	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	producer article.Producer
}

func (svc *articleService) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.repo.ListPub(ctx, start, offset, limit)
}

func NewarticleService(repo repository.ArticleRepository, producer article.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
	}
}

func (svc *articleService) GetPublishedById(ctx context.Context, id int64, uid int64) (domain.Article, error) {

	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		//发送事件了也就是某人读了某篇文章
		go func() {
			er := svc.producer.ProduceReadEvent(ctx, article.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				//记录日志
			}
		}()
	}
	return art, err

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

func (svc *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnPublished
	if article.Id > 0 {
		err := svc.repo.Update(ctx, article)
		return article.Id, err
	}
	return svc.repo.Create(ctx, article)
}

package repository

import (
	"context"
	"fmt"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
}

type articleRepository struct {
	dao dao.ArticleDAO
}

func (repo *articleRepository) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	res, err := repo.dao.GetByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	arts := []domain.Article{}
	for _, v := range res {
		arts = append(arts, repo.toDomain(v))
	}
	fmt.Println(arts)
	return arts, nil
}

func (repo *articleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, uid, uint8(status))
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Sync(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Update(ctx context.Context, article domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(article))
}

func (repo *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(article))
}

func NewarticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleRepository{dao: dao}
}

func (repo *articleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
		Ctime:    article.Ctime.UnixMilli(),
		Utime:    article.Utime.UnixMilli(),
	}
}
func (repo *articleRepository) toDomain(article dao.Article) domain.Article {
	return domain.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			Id: article.AuthorId,
		},
		Status: domain.ArticleStatus(article.Status),
		Ctime:  time.UnixMilli(article.Ctime),
		Utime:  time.UnixMilli(article.Utime),
	}
}

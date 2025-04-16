package repository

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/dao"
)

type ArticleRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleRepository struct {
	dao dao.ArticleDAO
}

func (repo *articleRepository) Save(ctx context.Context, article domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(article))
}

func NewarticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleRepository{dao: dao}
}

func (repo *articleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	}
}
func (repo *articleRepository) toDomain(article dao.Article) domain.Article {
	return domain.Article{
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			Id: article.AuthorId,
		},
	}
}

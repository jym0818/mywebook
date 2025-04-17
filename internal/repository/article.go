package repository

import (
	"context"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository/cache"
	"github.com/jym/mywebook/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type articleRepository struct {
	dao      dao.ArticleDAO
	cache    cache.ArticleCache
	userRepo UserRepository
}

func (repo *articleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		},
	}
	return res, nil
}

func (repo *articleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return repo.toDomain(art), nil
}

func (repo *articleRepository) List(ctx context.Context, uid int64, limit int, offset int) ([]domain.Article, error) {
	//先查找缓存
	if offset == 0 && limit <= 100 {
		data, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				repo.preCache(ctx, data)
			}()
			return data[:limit], nil
		}
	}
	res, err := repo.dao.GetByAuthor(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	arts := []domain.Article{}
	for _, v := range res {
		arts = append(arts, repo.toDomain(v))
	}
	//回写缓存
	go func() {
		repo.preCache(ctx, arts)
		err = repo.cache.SetFirstPage(ctx, uid, arts)
		if err != nil {
			//记录日志
		}
	}()
	return arts, nil
}

func (repo *articleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	go func() {
		err := repo.cache.DelFirstPage(ctx, uid)
		if err != nil {
			//记录日志
		}
	}()
	return repo.dao.SyncStatus(ctx, id, uid, uint8(status))
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	go func() {
		err := repo.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			//记录日志
		}
	}()
	return repo.dao.Sync(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Update(ctx context.Context, article domain.Article) error {
	go func() {
		err := repo.cache.DelFirstPage(ctx, article.Author.Id)
		if err != nil {
			//记录日志
		}
	}()
	return repo.dao.UpdateById(ctx, repo.toEntity(article))
}

func (repo *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	go func() {
		err := repo.cache.DelFirstPage(ctx, article.Author.Id)
		if err != nil {
			//记录日志
		}
	}()
	return repo.dao.Insert(ctx, repo.toEntity(article))
}

func (repo *articleRepository) preCache(ctx context.Context,
	arts []domain.Article) {
	// 1MB
	const contentSizeThreshold = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		// 你也可以记录日志
		if err := repo.cache.Set(ctx, arts[0]); err != nil {

		}
	}
}

func NewarticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository) ArticleRepository {
	return &articleRepository{dao: dao, cache: cache, userRepo: userRepo}
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

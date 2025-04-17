package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/mywebook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	DelFirstPage(ctx context.Context, uid int64) error
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	Set(ctx context.Context, article domain.Article) error
}

type RedisArticle struct {
	cmd redis.Cmdable
}

func (r *RedisArticle) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, r.authorArtKey(article.Id), data, time.Minute).Err()
}

func (r *RedisArticle) DelFirstPage(ctx context.Context, uid int64) error {
	return r.cmd.Del(ctx, r.firstPageKey(uid)).Err()
}

func (r *RedisArticle) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := range arts {
		// 只缓存摘要部分
		arts[i].Content = arts[i].Abstract()
	}
	res, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, r.firstPageKey(uid), string(res), time.Minute*10).Err()
}

func (r *RedisArticle) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	bs, err := r.cmd.Get(ctx, r.firstPageKey(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(bs, &arts)
	return arts, err
}
func (r *RedisArticle) firstPageKey(author int64) string {
	return fmt.Sprintf("article:first_page:%d", author)
}
func (r *RedisArticle) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

func NewRedisArticle(cmd redis.Cmdable) ArticleCache {
	return &RedisArticle{cmd: cmd}
}

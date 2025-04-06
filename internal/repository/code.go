package repository

import (
	"context"
	"github.com/jym/mywebook/internal/repository/cache"
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verfiy(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeRepository struct {
	cache cache.CodeCache
}

func (c *codeRepository) Verfiy(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}

func (c *codeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func NewcodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}

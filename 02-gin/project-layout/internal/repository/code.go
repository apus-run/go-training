package repository

import (
	"context"

	"project-layout/internal/repository/cache/code"
)

var (
	ErrCodeVerifyTooManyTimes = code.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = code.ErrCodeSendTooMany
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error

	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeRepository struct {
	cache code.CodeCache
}

func NewCodeRepository(cache code.CodeCache) CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}

func (c codeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c codeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}

package usecases

import (
	"context"
	"time"

	"github.com/jdluques/url-shortener/internal/application/ports"
	"github.com/jdluques/url-shortener/internal/domain/url"
)

type RedirectUseCase struct {
	repo  url.URLRepository
	cache ports.Cache
	now   func() time.Time
}

func NewRedirectUseCase(
	repo url.URLRepository,
	cache ports.Cache,
) *RedirectUseCase {
	return &RedirectUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (usecase *RedirectUseCase) Execute(
	ctx context.Context,
	code string,
) (string, error) {
	cache_hit, err := usecase.cache.Get(ctx, code)
	if err == nil {
		return cache_hit, nil
	}

	original_url, err := usecase.repo.FindByCode(ctx, code)
	if err != nil {
		return "", err
	}

	if err := usecase.cache.Set(ctx, code, original_url.OriginalURL(), 24*time.Hour); err != nil {
		return "", err
	}

	return original_url.OriginalURL(), nil
}

package usecases

import (
	"context"
	"time"

	"github.com/jdluques/url-shortener/internal/application/ports"
	"github.com/jdluques/url-shortener/internal/domain/url"
)

type ShortenURLUseCase struct {
	repo         url.URLRepository
	cache        ports.Cache
	idGen        url.IDGenerator
	shortCodeGen url.ShortCodeGenerator
	now          func() time.Time
}

func NewShortenURLUseCase(
	repo url.URLRepository,
	cache ports.Cache,
	idGen url.IDGenerator,
	shortCodeGen url.ShortCodeGenerator,
) *ShortenURLUseCase {
	return &ShortenURLUseCase{
		repo:         repo,
		cache:        cache,
		idGen:        idGen,
		shortCodeGen: shortCodeGen,
	}
}

func (service *ShortenURLUseCase) Execute(
	ctx context.Context,
	originalURL string,
	expiresIn time.Duration,
) (*url.URL, error) {
	id, err := service.idGen.NextID()
	if err != nil {
		return nil, err
	}

	shortCode, err := service.shortCodeGen.Generate(id)
	if err != nil {
		return nil, err
	}

	url, err := url.NewURL(
		id,
		originalURL,
		shortCode,
		service.now(),
		expiresIn,
	)
	if err != nil {
		return nil, err
	}

	if err := service.repo.Save(ctx, url); err != nil {
		return nil, err
	}

	key := url.ShortCode()
	value := url.OriginalURL()
	ttl := 24 * time.Hour
	if err := service.cache.Set(ctx, key, value, ttl); err != nil {
		return nil, err
	}

	return url, nil
}

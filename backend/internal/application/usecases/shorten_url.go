package usecases

import (
	"context"
	"time"

	"github.com/jdluques/url-shortener/internal/domain/url"
)

type ShortenURLUseCase struct {
	repo         url.URLRepository
	idGen        url.IDGenerator
	shortCodeGen url.ShortCodeGenerator
	now          func() time.Time
}

func NewShortenURLUseCase(
	repo url.URLRepository,
	idGen url.IDGenerator,
	shortCodeGen url.ShortCodeGenerator,
) *ShortenURLUseCase {
	return &ShortenURLUseCase{
		repo:         repo,
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

	return url, nil
}

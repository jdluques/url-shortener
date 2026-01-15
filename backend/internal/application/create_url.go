package application

import (
	"context"
	"time"

	"github.com/jdluques/url-shortener/internal/domain/url"
)

type CreateURLService struct {
	repo         url.URLRepository
	idGen        url.IDGenerator
	shortCodeGen url.ShortCodeGenerator
	now          func() time.Time
}

func NewCreateURLService(
	repo url.URLRepository,
	idGen url.IDGenerator,
	shortCodeGen url.ShortCodeGenerator,
) *CreateURLService {
	return &CreateURLService{
		repo:         repo,
		idGen:        idGen,
		shortCodeGen: shortCodeGen,
	}
}

func (service *CreateURLService) Create(
	ctx context.Context,
	originalURL string,
	expiresAt time.Time,
) (*url.URL, error) {
	id, err := service.idGen.NextID()
	if err != nil {
		return nil, err
	}

	shortCode, err := service.shortCodeGen.Generate(id)
	if err != nil {
		return nil, err
	}

	expiresIn := 30 * 24 * time.Hour

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

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jdluques/url-shortener/internal/domain/url"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (repo *URLRepository) Save(ctx context.Context, url *url.URL) error {
	_, err := repo.db.ExecContext(
		ctx,
		`INSERT INTO urls (id, original_url, short_code, clicks, expires_at)
         VALUES ($1, $2, $3, $4`,
		url.ID(),
		url.OriginalURL(),
		url.ShortCode(),
		url.Clicks(),
		url.ExpiresAt(),
	)
	return err
}

func (repo *URLRepository) FindByCode(ctx context.Context, code string) (*url.URL, error) {
	row := repo.db.QueryRowContext(
		ctx,
		`SELECT id, original_url, short_code, clicks, expires_at, created_at
         FROM urls
         WHERE short_code=$1`,
		code,
	)

	var id, clicks int64
	var original_url, short_code string
	var expires_at, created_at time.Time
	err := row.Scan(&id, &original_url, &short_code, &clicks, &expires_at, &created_at)
	if err != nil {
		return nil, err
	}

	return url.HydrateURL(id, clicks, original_url, short_code, created_at, expires_at)
}

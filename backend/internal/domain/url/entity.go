package url

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidId         = errors.New("invalid url id")
	ErrInvalidURLString  = errors.New("invalid url string")
	ErrInvalidaShortCode = errors.New("invalid url short code")
)

type URL struct {
	id int64

	originalURL string
	shortCode   string
	clicks      int64

	createdAt time.Time
	expiresAt time.Time
}

// Creation

func NewURL(id int64, originalURL, shortCode string, now time.Time, expiresIn time.Duration) (*URL, error) {
	url := URL{
		id:          id,
		originalURL: originalURL,
		shortCode:   shortCode,
		clicks:      0,
		createdAt:   now,
		expiresAt:   now.Add(expiresIn),
	}

	if err := url.ValidateURL(); err != nil {
		return nil, err
	}

	return &url, nil
}

func HydrateURL(id, clicks int64, originalUrl, shortCode string, createdAt, expiresAt time.Time) (*URL, error) {
	url := URL{
		id:          id,
		originalURL: originalUrl,
		shortCode:   shortCode,
		clicks:      clicks,
		createdAt:   createdAt,
		expiresAt:   expiresAt,
	}

	if err := url.ValidateURL(); err != nil {
		return nil, err
	}

	return &url, nil
}

// Validation

func (url URL) ValidateURL() error {
	if err := validateId(url.id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}
	if err := validateURLString(url.originalURL); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}
	if err := validateShortCode(url.shortCode); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}
	if err := validateClicks(url.clicks); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	return nil
}

func validateId(id int64) error {
	if id == 0 {
		return ErrInvalidId
	}

	return nil
}

func validateURLString(rawURL string) error {
	if rawURL == "" {
		return ErrInvalidURLString
	}

	url, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURLString
	}

	if url.Scheme == "" || url.Host == "" {
		return ErrInvalidURLString
	}

	return nil
}

func validateShortCode(code string) error {
	if code == "" {
		return ErrInvalidaShortCode
	}

	return nil
}

func validateClicks(clicks int64) error {
	if clicks < 0 {
		return fmt.Errorf("%w: clicks must be non-negative integers", ErrInvalidURL)
	}

	return nil
}

// Getters

func (url URL) ID() int64 {
	return url.id
}

func (url URL) OriginalURL() string {
	return url.originalURL
}

func (url URL) ShortCode() string {
	return url.shortCode
}

func (url URL) Clicks() int64 {
	return url.clicks
}

func (url URL) CreatedAt() time.Time {
	return url.createdAt
}

func (url URL) ExpiresAt() time.Time {
	return url.expiresAt
}

// Behaviours

func (url *URL) Use() {
	url.clicks++
}


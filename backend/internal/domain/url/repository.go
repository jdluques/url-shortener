package url

import "context"

type URLRepository interface {
	Save(ctx context.Context, url *URL) error
	FindByCode(ctx context.Context, code string) (*URL, error)
}

package url

import (
	"context"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type Repository interface {
	Create(ctx context.Context, url *urlDomain.Url) (int64, error)
	GetByShortCode(ctx context.Context, id string) (*urlDomain.Url, error)
	IncrementClicks(ctx context.Context, id string) error
	UpdateShortCode(ctx context.Context, id int64, shortCode string) error
}

type Usecase interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

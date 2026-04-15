package url

import (
	"context"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type Repository interface {
	Create(ctx context.Context, url *urlDomain.Url) error
	GetByID(ctx context.Context, id string) (*urlDomain.Url, error)
	IncrementClicks(ctx context.Context, id string) error
}

type Usecase interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

type CreateInput struct {
}

type UpdateInput struct {
}

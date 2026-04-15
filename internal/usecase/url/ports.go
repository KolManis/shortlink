package url

import (
	"context"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	Create(ctx context.Context, tx pgx.Tx, url *urlDomain.Url) (int64, error)
	UpdateShortCode(ctx context.Context, tx pgx.Tx, id int64, shortCode string) error
	GetByShortCode(ctx context.Context, id string) (*urlDomain.Url, error)
	IncrementClicks(ctx context.Context, id string) error
}

type Usecase interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

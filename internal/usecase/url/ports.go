package url

import (
	"context"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/jackc/pgx/v5"
)

// mockery --name=Repository --dir=internal/usecase/url --output=internal/mocks
// mockery --name=Cache --dir=internal/usecase/url --output=internal/mocks
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	//Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
}

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	Create(ctx context.Context, tx pgx.Tx, url *urlDomain.Url) (int64, error)
	UpdateShortCode(ctx context.Context, tx pgx.Tx, id int64, shortCode string) error
	GetByShortCode(ctx context.Context, id string) (*urlDomain.Url, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*urlDomain.Url, error)
	IncrementClicks(ctx context.Context, id string) error
}

type Usecase interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
	// GetByShortCode(ctx context.Context, shortCode string)
	// IncrementClicks(ctx context.Context, shortCode string)
}

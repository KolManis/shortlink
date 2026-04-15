package postgres

import (
	"context"
	"errors"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

type urlScanner interface {
	Scan(dest ...any) error
}

func (r *Repository) Create(ctx context.Context, url *urlDomain.Url) error {
	const query = `
        INSERT INTO links (id, original_url, created_at, clicks)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.pool.Exec(ctx, query,
		url.ID,
		url.OriginalURL,
		url.CreatedAt,
		url.Clicks,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*urlDomain.Url, error) {
	const query = `
		SELECT id, original_url, created_at, clicks
		FROM links
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanUrl(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, urlDomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) IncrementClicks(ctx context.Context, id string) error {
	query := `
		UPDATE links
		SET clicks = clicks + 1
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func scanUrl(scanner urlScanner) (*urlDomain.Url, error) {
	var (
		url urlDomain.Url
	)

	if err := scanner.Scan(
		&url.ID,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.Clicks,
	); err != nil {
		return nil, err
	}

	return &url, nil
}

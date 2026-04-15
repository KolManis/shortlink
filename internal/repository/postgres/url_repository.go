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

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *Repository) Create(ctx context.Context, tx pgx.Tx, url *urlDomain.Url) (int64, error) {
	const query = `
        INSERT INTO links (original_url, created_at, clicks)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64
	var err error
	if tx != nil {
		err = tx.QueryRow(ctx, query, url.OriginalURL, url.CreatedAt, url.Clicks).Scan(&id)
	} else {
		err = r.pool.QueryRow(ctx, query, url.OriginalURL, url.CreatedAt, url.Clicks).Scan(&id)
	}
	return id, err
}

func (r *Repository) GetByShortCode(ctx context.Context, shortCode string) (*urlDomain.Url, error) {
	const query = `
        SELECT id, short_code, original_url, created_at, clicks
        FROM links
        WHERE short_code = $1
    `

	var url urlDomain.Url
	err := r.pool.QueryRow(ctx, query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.Clicks,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, urlDomain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *Repository) IncrementClicks(ctx context.Context, shortCode string) error {
	query := `
        UPDATE links
        SET clicks = clicks + 1
        WHERE short_code = $1
    `
	_, err := r.pool.Exec(ctx, query, shortCode)
	return err
}

func (r *Repository) UpdateShortCode(ctx context.Context, tx pgx.Tx, id int64, shortCode string) error {
	const query = `
        UPDATE links
        SET short_code = $1
        WHERE id = $2
    `
	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, shortCode, id)
	} else {
		_, err = r.pool.Exec(ctx, query, shortCode, id)
	}
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

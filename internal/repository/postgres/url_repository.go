package postgres

import (
	"context"
	"errors"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db DB
}

func New(db DB) *Repository {
	return &Repository{db: db}
}

type urlScanner interface {
	Scan(dest ...any) error
}

func (r *Repository) Create(ctx context.Context, url *urlDomain.Url) (int64, error) {
	const query = `
        INSERT INTO links (original_url, created_at, clicks)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int64

	err := r.db.QueryRow(
		ctx,
		query,
		url.OriginalURL,
		url.CreatedAt,
		url.Clicks,
	).Scan(&id)

	return id, err
}

func (r *Repository) GetByOriginalURL(ctx context.Context, originalURL string) (*urlDomain.Url, error) {
	const query = `
        SELECT id, short_code, original_url, created_at, clicks
        FROM links
        WHERE original_url = $1
    `

	var url urlDomain.Url
	err := r.db.QueryRow(ctx, query, originalURL).Scan(
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

func (r *Repository) GetByShortCode(ctx context.Context, shortCode string) (*urlDomain.Url, error) {
	const query = `
        SELECT id, short_code, original_url, created_at, clicks
        FROM links
        WHERE short_code = $1
    `

	var url urlDomain.Url
	err := r.db.QueryRow(ctx, query, shortCode).Scan(
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
	_, err := r.db.Exec(ctx, query, shortCode)
	return err
}

func (r *Repository) UpdateShortCode(ctx context.Context, id int64, shortCode string) error {
	const query = `
        UPDATE links
        SET short_code = $1
        WHERE id = $2
    `

	_, err := r.db.Exec(ctx, query, shortCode, id)

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

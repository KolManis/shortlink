package postgres

import (
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

// func scanTask(scanner urlScanner) (..., error) {
// 	var (
// 		status string
// 	)

// 	if err := scanner.Scan(
// 		&....ID,
// 		&.....CreatedAt,
// 		&.....UpdatedAt,
// 	); err != nil {
// 		return nil, err
// 	}

// 	.....Status = .....Status(status)

// 	return &...., nil
// }

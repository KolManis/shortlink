package uow

import (
	"context"

	"github.com/KolManis/shortlink/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(r *postgres.Repository) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}

	repo := postgres.New(tx)
	err = fn(repo)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Adapter struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Adapter, error) {
	a := &Adapter{}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}
	a.Pool = pool

	return a, nil
}

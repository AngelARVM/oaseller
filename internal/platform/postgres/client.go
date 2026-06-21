package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL is required.")
	}

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

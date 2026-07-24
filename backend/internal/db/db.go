package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to reach database: %w", err)
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	schema, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		return fmt.Errorf("unable to read schema file: %w", err)
	}
	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		return fmt.Errorf("unable to execute schema: %w", err)
	}
	fmt.Println("Database schema applied successfully.")
	return nil
}

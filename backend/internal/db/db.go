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
	schema, err := os.ReadFile("internal/db/learning.sql")
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

// EnsureDefaultUser inserts a placeholder user row if it doesn't already
// exist. It is a temporary stand-in until Supabase Auth is wired in and
// requests carry a real authenticated user id.
func EnsureDefaultUser(ctx context.Context, pool *pgxpool.Pool, id, username, email string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash)
		VALUES ($1, $2, $3, 'not-a-real-password-hash')
		ON CONFLICT (id) DO NOTHING
	`, id, username, email)
	if err != nil {
		return fmt.Errorf("unable to ensure default user: %w", err)
	}
	return nil
}

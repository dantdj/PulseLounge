package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func initDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := waitForDB(db, 10, 2*time.Second); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("wait for database: %w; close database: %v", err, closeErr)
		}
		return nil, fmt.Errorf("wait for database: %w", err)
	}

	if err := ensureSchema(db); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("ensure schema: %w; close database: %v", err, closeErr)
		}
		return nil, fmt.Errorf("ensure schema: %w", err)
	}

	return db, nil
}

func buildDatabaseURLFromEnv() string {
	user := getenvOrDefault("POSTGRES_USER", "postgres")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := getenvOrDefault("POSTGRES_HOST", "localhost")
	port := getenvOrDefault("POSTGRES_PORT", "5432")
	name := getenvOrDefault("POSTGRES_DB", "pulselounge")
	sslMode := getenvOrDefault("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		name,
		url.QueryEscape(sslMode),
	)
}

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func waitForDB(db *sql.DB, attempts int, delay time.Duration) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := db.PingContext(ctx)
		cancel()
		if err == nil {
			return nil
		}

		lastErr = err
		time.Sleep(delay)
	}

	return lastErr
}

func ensureSchema(db *sql.DB) error {
	const createMessagesTable = `
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    body TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, createMessagesTable); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("schema creation timed out: %w", err)
		}
		return err
	}

	return nil
}

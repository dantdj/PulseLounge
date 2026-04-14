package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"pulselounge/internal/postgres"
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

	return db, nil
}

func buildDatabaseURLFromEnv() string {
	return postgres.BuildURLFromEnv()
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

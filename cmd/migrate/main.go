package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"pulselounge/internal/postgres"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate <up|status>")
	}

	command := os.Args[1]
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = postgres.BuildURLFromEnv()
	}

	db, err := openDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database: %v", closeErr)
		}
	}()

	log.SetFlags(0)
	goose.SetLogger(&stdoutLogger{})
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to configure goose dialect: %v", err)
	}

	dir := defaultDir()

	switch command {
	case "up":
		if err := goose.Up(db, dir); err != nil {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	case "status":
		if err := goose.Status(db, dir); err != nil {
			log.Fatalf("failed to read migration status: %v", err)
		}
	default:
		log.Fatalf("unknown migration command %q; expected up or status", command)
	}
}

func defaultDir() string {
	if dir := os.Getenv("MIGRATIONS_DIR"); dir != "" {
		return dir
	}
	return "db/migrations"
}

func openDB(databaseURL string) (*sql.DB, error) {
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

func waitForDB(db *sql.DB, attempts int, delay time.Duration) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}

		lastErr = err
		time.Sleep(delay)
	}

	return lastErr
}

type stdoutLogger struct{}

func (*stdoutLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func (*stdoutLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

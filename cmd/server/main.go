package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//go:embed web/dist/*
var embeddedUI embed.FS

type healthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type message struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type app struct {
	db *sql.DB
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = buildDatabaseURLFromEnv()
	}

	db, err := initDB(databaseURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	uiFS, err := fs.Sub(embeddedUI, "web/dist")
	if err != nil {
		log.Fatalf("failed to load embedded UI: %v", err)
	}

	application := app{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/messages", application.listMessagesHandler)
	mux.HandleFunc("/", spaHandler(uiFS))

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (a app) listMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := a.db.QueryContext(r.Context(), "SELECT id, body, created_at FROM messages ORDER BY id")
	if err != nil {
		http.Error(w, "failed to query messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := make([]message, 0)
	for rows.Next() {
		var m message
		if err := rows.Scan(&m.ID, &m.Body, &m.CreatedAt); err != nil {
			http.Error(w, "failed to decode messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "message iteration failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messages)
}

func initDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := waitForDB(db, 10, 2*time.Second); err != nil {
		db.Close()
		return nil, fmt.Errorf("wait for database: %w", err)
	}

	if err := ensureSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("ensure schema: %w", err)
	}

	return db, nil
}

func buildDatabaseURLFromEnv() string {
	user := getenvDefault("POSTGRES_USER", "postgres")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := getenvDefault("POSTGRES_HOST", "localhost")
	port := getenvDefault("POSTGRES_PORT", "5432")
	name := getenvDefault("POSTGRES_DB", "pulselounge")
	sslMode := getenvDefault("POSTGRES_SSLMODE", "disable")

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

func getenvDefault(key, fallback string) string {
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

func spaHandler(uiFS fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(uiFS))

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if cleanPath == "." {
			cleanPath = ""
		}

		if cleanPath != "" {
			if _, err := fs.Stat(uiFS, cleanPath); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		indexHTML, err := fs.ReadFile(uiFS, "index.html")
		if err != nil {
			http.Error(w, "index file not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	}
}

package messages

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresRepositoryCreateAndList(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to run Postgres integration tests")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close database: %v", err)
		}
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    body TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	repo := NewPostgresRepository(db)
	bodyPrefix := fmt.Sprintf("integration-test-%d", time.Now().UnixNano())
	cleanupPattern := bodyPrefix + "%"
	t.Cleanup(func() {
		if _, err := db.ExecContext(context.Background(), "DELETE FROM messages WHERE body LIKE $1", cleanupPattern); err != nil {
			t.Fatalf("cleanup test rows: %v", err)
		}
	})

	created, err := repo.Create(ctx, bodyPrefix)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected created message to have an id")
	}
	if created.Body != bodyPrefix {
		t.Fatalf("expected created body %q, got %q", bodyPrefix, created.Body)
	}
	if created.CreatedAt.IsZero() {
		t.Fatalf("expected created_at to be populated")
	}

	listed, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	var found Message
	var ok bool
	for _, message := range listed {
		if strings.HasPrefix(message.Body, bodyPrefix) {
			found = message
			ok = true
			break
		}
	}
	if !ok {
		t.Fatalf("expected created message to appear in list")
	}
	if found.ID != created.ID {
		t.Fatalf("expected listed id %d, got %d", created.ID, found.ID)
	}
}

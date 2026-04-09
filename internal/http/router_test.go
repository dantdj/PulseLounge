package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pulselounge/internal/messages"
)

func TestRouterWiresListMessagesEndpoint(t *testing.T) {
	repo := &fakeMessageRepo{listResult: []messages.Message{}}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected list to be called once, got %d", repo.listCalls)
	}
}

func TestRouterWiresCreateMessagesEndpoint(t *testing.T) {
	createdAt := time.Date(2026, time.January, 4, 12, 0, 0, 0, time.UTC)
	repo := &fakeMessageRepo{
		createFn: func(_ context.Context, body string) (messages.Message, error) {
			return messages.Message{
				ID:        10,
				Body:      body,
				CreatedAt: createdAt,
			}, nil
		},
	}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPost, "/api/messages", strings.NewReader(`{"body":"from router"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if repo.createCall != 1 {
		t.Fatalf("expected create to be called once, got %d", repo.createCall)
	}
}

func TestRouterWiresEditMessagesEndpoint(t *testing.T) {
	repo := &fakeMessageRepo{}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader(`{"id":1,"newBody":"updated from router"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
}

func TestRouterReturnsNotFoundForUnknownEndpoint(t *testing.T) {
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(&fakeMessageRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

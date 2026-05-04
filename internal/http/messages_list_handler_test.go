package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pulselounge/internal/messages"
)

func TestMessagesHandlerListSuccess(t *testing.T) {
	createdAt := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	repo := &fakeMessageRepo{
		listResult: []messages.Message{
			{ID: 1, Body: "hello", CreatedAt: createdAt},
			{ID: 2, Body: "world", CreatedAt: createdAt.Add(time.Minute)},
		},
	}
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/channels/7/messages", nil)
	rr := httptest.NewRecorder()
	handler.List(rr, req, 7)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected list to be called once, got %d", repo.listCalls)
	}
	if repo.listChannelID != 7 {
		t.Fatalf("expected channel id 7, got %d", repo.listChannelID)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %q", contentType)
	}

	var got []messages.Message
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(got))
	}
	if got[0].ID != 1 || got[0].Body != "hello" {
		t.Fatalf("unexpected first message: %+v", got[0])
	}
}

func TestMessagesHandlerListMethodNotAllowed(t *testing.T) {
	repo := &fakeMessageRepo{}
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", nil)
	rr := httptest.NewRecorder()
	handler.List(rr, req, 7)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	if repo.listCalls != 0 {
		t.Fatalf("expected list not to be called, got %d", repo.listCalls)
	}
	assertJSONError(t, rr, "method not allowed")
}

func TestMessagesHandlerListServiceError(t *testing.T) {
	repo := &fakeMessageRepo{listErr: errors.New("boom")}
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/channels/7/messages", nil)
	rr := httptest.NewRecorder()
	handler.List(rr, req, 7)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected list to be called once, got %d", repo.listCalls)
	}
	assertJSONError(t, rr, "failed to query messages")
}

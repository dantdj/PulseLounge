package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pulselounge/internal/messages"
)

func TestMessagesHandlerCreateSuccess(t *testing.T) {
	createdAt := time.Date(2026, time.January, 3, 12, 0, 0, 0, time.UTC)
	repo := &fakeMessageRepo{
		createFn: func(_ context.Context, channelID int64, authorID int64, body string, imageKey string) (messages.Message, error) {
			return messages.Message{
				ID:        9,
				AuthorID:  authorID,
				ChannelID: channelID,
				Body:      body,
				CreatedAt: createdAt,
			}, nil
		},
	}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", strings.NewReader(`{"body":"new msg"}`))
	rr := httptest.NewRecorder()
	handler.Create(rr, req, 7)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %q", contentType)
	}
	if repo.createChannelID != 7 {
		t.Fatalf("expected channel id 7, got %d", repo.createChannelID)
	}
	if repo.createAuthorID != devAuthorID {
		t.Fatalf("expected author id %d, got %d", devAuthorID, repo.createAuthorID)
	}
}

func TestMessagesHandlerCreateMethodNotAllowed(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodGet, "/api/channels/7/messages", nil)
	rr := httptest.NewRecorder()
	handler.Create(rr, req, 7)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	assertJSONError(t, rr, "method not allowed")
}

func TestMessagesHandlerCreateInvalidBody(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", strings.NewReader("{"))
	rr := httptest.NewRecorder()
	handler.Create(rr, req, 7)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if repo.createCall != 0 {
		t.Fatalf("expected create not to be called, got %d", repo.createCall)
	}
	assertJSONError(t, rr, "invalid request body")
}

func TestMessagesHandlerCreateRejectsBlankBody(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", strings.NewReader(`{"body":"   "}`))
	rr := httptest.NewRecorder()
	handler.Create(rr, req, 7)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if repo.createCall != 0 {
		t.Fatalf("expected create not to be called, got %d", repo.createCall)
	}
	assertJSONError(t, rr, "message body cannot be empty")
}

func TestMessagesHandlerCreateServiceError(t *testing.T) {
	repo := &fakeMessageRepo{
		createErr: errors.New("boom"),
	}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", strings.NewReader(`{"body":"new msg"}`))
	rr := httptest.NewRecorder()
	handler.Create(rr, req, 7)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if repo.createCall != 1 {
		t.Fatalf("expected create to be called once, got %d", repo.createCall)
	}
	assertJSONError(t, rr, "failed to create message")
}

func assertJSONError(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %q", contentType)
	}

	var got struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if got.Error != want {
		t.Fatalf("expected error %q, got %q", want, got.Error)
	}
}

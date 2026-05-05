package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"pulselounge/internal/messages"
	"strings"
	"testing"
)

func TestMessagesHandlerEditMethodNotAllowed(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodGet, "/api/messages/1", nil)
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 1)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	if repo.editCalls != 0 {
		t.Fatalf("expected edit not to be called, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "method not allowed")
}

func TestMessagesHandlerEditInvalidBody(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPut, "/api/messages/1", strings.NewReader("{"))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 1)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if repo.editCalls != 0 {
		t.Fatalf("expected edit not to be called, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "invalid request body")
}

func TestMessagesHandlerEditRejectsBlankBody(t *testing.T) {
	repo := &fakeMessageRepo{}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPut, "/api/messages/1", strings.NewReader(`{"newBody":"   "}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 1)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if repo.editCalls != 0 {
		t.Fatalf("expected edit not to be called, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "message body cannot be empty")
}

func TestMessagesHandlerEditMessageNotFound(t *testing.T) {
	repo := &fakeMessageRepo{
		editFn: func(ctx context.Context, id int64, body string) error {
			return messages.ErrMessageNotFound
		},
	}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPut, "/api/messages/999", strings.NewReader(`{"newBody":"updated msg"}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 999)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "message not found")
}

func TestMessagesHandlerEditServiceError(t *testing.T) {
	repo := &fakeMessageRepo{
		editErr: errors.New("boom"),
	}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)

	req := httptest.NewRequest(http.MethodPut, "/api/messages/1", strings.NewReader(`{"newBody":"updated msg"}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 1)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "failed to edit message")
}

func TestMessagesHandlerEditSuccess(t *testing.T) {
	var capturedID int64
	var capturedBody string

	repo := &fakeMessageRepo{
		editFn: func(_ context.Context, id int64, body string) error {
			capturedID = id
			capturedBody = body
			return nil
		},
	}
	store := &fakeMediaStore{}
	handler := NewMessagesHandler(messages.NewService(repo), store)
	newBody := "updated msg"
	reqBody := `{"newBody":"` + newBody + `"}`

	req := httptest.NewRequest(http.MethodPut, "/api/messages/42", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req, 42)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
	if capturedID != 42 {
		t.Fatalf("expected edit to be called with id 42, got %d", capturedID)
	}
	if capturedBody != newBody {
		t.Fatalf("expected edit to be called with body %q, got %q", newBody, capturedBody)
	}
}

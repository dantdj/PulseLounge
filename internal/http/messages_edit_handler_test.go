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
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

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
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader("{"))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

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
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader(`{"id":1,"newBody":"   "}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

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
		editFn: func(ctx context.Context, id int, body string) error {
			return messages.ErrMessageNotFound
		},
	}
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader(`{"id":999,"newBody":"updated msg"}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

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
	handler := NewMessagesHandler(messages.NewService(repo))

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader(`{"id":1,"newBody":"updated msg"}`))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
	assertJSONError(t, rr, "failed to edit message")
}

func TestMessagesHandlerEditSuccess(t *testing.T) {
	var capturedID int
	var capturedBody string

	repo := &fakeMessageRepo{
		editFn: func(_ context.Context, id int, body string) error {
			capturedID = id
			capturedBody = body
			return nil
		},
	}
	handler := NewMessagesHandler(messages.NewService(repo))
	newBody := "updated msg"
	reqBody := `{"id":42,"newBody":"` + newBody + `"}`

	req := httptest.NewRequest(http.MethodPut, "/api/messages", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()
	handler.Edit(rr, req)

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

package httpapi

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggingPreservesRequestID(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	handler := withRequestLogging(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := RequestIDFromContext(r.Context()); got != "test-request-id" {
			t.Fatalf("expected request id in context, got %q", got)
		}
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/messages", nil)
	req.Header.Set(requestIDHeader, "test-request-id")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get(requestIDHeader); got != "test-request-id" {
		t.Fatalf("expected response request id header, got %q", got)
	}
	for _, want := range []string{
		"msg=\"http request completed\"",
		"request_id=test-request-id",
		"method=POST",
		"path=/api/messages",
		"status=201",
	} {
		if !strings.Contains(logs.String(), want) {
			t.Fatalf("expected log to contain %q, got %q", want, logs.String())
		}
	}
}

func TestPanicRecoveryLogsAndReturnsInternalServerError(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	handler := withRequestLogging(logger, withPanicRecovery(logger, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))

	req := httptest.NewRequest(http.MethodGet, "/api/panic", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	for _, want := range []string{
		"msg=\"panic while handling request\"",
		"panic=boom",
		"status=500",
	} {
		if !strings.Contains(logs.String(), want) {
			t.Fatalf("expected log to contain %q, got %q", want, logs.String())
		}
	}
}

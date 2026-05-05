package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"pulselounge/internal/channels"
	"pulselounge/internal/messages"
)

func TestRouterWiresListChannelMessagesEndpoint(t *testing.T) {
	repo := &fakeMessageRepo{listResult: []messages.Message{}}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo), testChannelService(), testMediaStore())

	req := httptest.NewRequest(http.MethodGet, "/api/channels/7/messages", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected list to be called once, got %d", repo.listCalls)
	}
	if repo.listChannelID != 7 {
		t.Fatalf("expected channel id 7, got %d", repo.listChannelID)
	}
}

func TestRouterWiresCreateChannelMessagesEndpoint(t *testing.T) {
	createdAt := time.Date(2026, time.January, 4, 12, 0, 0, 0, time.UTC)
	repo := &fakeMessageRepo{
		createFn: func(_ context.Context, channelID int64, authorID int64, body string, imageKey string) (messages.Message, error) {
			return messages.Message{
				ID:        10,
				AuthorID:  authorID,
				ChannelID: channelID,
				Body:      body,
				CreatedAt: createdAt,
			}, nil
		},
	}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo), testChannelService(), testMediaStore())

	req := httptest.NewRequest(http.MethodPost, "/api/channels/7/messages", strings.NewReader(`{"body":"from router"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if repo.createCall != 1 {
		t.Fatalf("expected create to be called once, got %d", repo.createCall)
	}
	if repo.createChannelID != 7 {
		t.Fatalf("expected channel id 7, got %d", repo.createChannelID)
	}
}

func TestRouterWiresEditMessagesEndpoint(t *testing.T) {
	repo := &fakeMessageRepo{}
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(repo), testChannelService(), testMediaStore())

	req := httptest.NewRequest(http.MethodPut, "/api/messages/1", strings.NewReader(`{"newBody":"updated from router"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if repo.editCalls != 1 {
		t.Fatalf("expected edit to be called once, got %d", repo.editCalls)
	}
}

func TestRouterWiresDeleteChannelsEndpoint(t *testing.T) {
	channelRepo := &fakeChannelRepo{}
	mux := NewRouter(
		http.NotFoundHandler(),
		messages.NewService(&fakeMessageRepo{}),
		channels.NewService(channelRepo),
		testMediaStore(),
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/channels/7", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if channelRepo.deleteCall != 1 {
		t.Fatalf("expected delete to be called once, got %d", channelRepo.deleteCall)
	}
	if channelRepo.deleteID != 7 {
		t.Fatalf("expected channel id 7, got %d", channelRepo.deleteID)
	}
}

func TestRouterReturnsNotFoundForUnknownEndpoint(t *testing.T) {
	mux := NewRouter(http.NotFoundHandler(), messages.NewService(&fakeMessageRepo{}), testChannelService(), testMediaStore())

	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

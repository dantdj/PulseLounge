package messages

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	listResult  []Message
	listErr     error
	createFn    func(context.Context, string) (Message, error)
	createCalls int
}

func (f *fakeRepository) List(context.Context) ([]Message, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listResult, nil
}

func (f *fakeRepository) Create(ctx context.Context, body string) (Message, error) {
	f.createCalls++
	if f.createFn != nil {
		return f.createFn(ctx, body)
	}
	return Message{}, nil
}

func TestServiceCreateRejectsEmptyBody(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)

	_, err := service.Create(context.Background(), "   \n\t ")
	if !errors.Is(err, ErrEmptyBody) {
		t.Fatalf("expected ErrEmptyBody, got %v", err)
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected repository create not to be called, got %d", repo.createCalls)
	}
}

func TestServiceCreateTrimsBodyBeforePersisting(t *testing.T) {
	createdAt := time.Date(2026, time.March, 10, 10, 0, 0, 0, time.UTC)
	repo := &fakeRepository{
		createFn: func(_ context.Context, body string) (Message, error) {
			return Message{ID: 1, Body: body, CreatedAt: createdAt}, nil
		},
	}
	service := NewService(repo)

	got, err := service.Create(context.Background(), "  hello world  ")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got.Body != "hello world" {
		t.Fatalf("expected trimmed body, got %q", got.Body)
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected repository create to be called once, got %d", repo.createCalls)
	}
}

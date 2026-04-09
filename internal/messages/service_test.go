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
	editFn      (func(ctx context.Context, id int, body string) error)
	editCalls   int
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

func (f *fakeRepository) Edit(ctx context.Context, id int, body string) error {
	f.editCalls++
	if f.editFn != nil {
		return f.editFn(ctx, id, body)
	}
	return nil
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

func TestServiceEditRejectsEmptyBody(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)

	// Provide whitespace-only body to Edit, which should be rejected with ErrEmptyBody
	err := service.Edit(context.Background(), 1, "   \n\t ")
	if !errors.Is(err, ErrEmptyBody) {
		t.Fatalf("expected ErrEmptyBody, got %v", err)
	}
}

func TestServiceEditTrimsBodyBeforePersisting(t *testing.T) {
	var capturedID int
	var capturedBody string

	repo := &fakeRepository{
		editFn: func(_ context.Context, id int, body string) error {
			capturedID = id
			capturedBody = body
			return nil
		},
	}

	service := NewService(repo)

	err := service.Edit(context.Background(), 42, "  updated message  ")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if capturedID != 42 {
		t.Fatalf("expected ID 42, got %d", capturedID)
	}
	if capturedBody != "updated message" {
		t.Fatalf("expected trimmed body 'updated message', got %q", capturedBody)
	}
}

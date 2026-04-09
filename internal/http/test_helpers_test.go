package httpapi

import (
	"context"

	"pulselounge/internal/messages"
)

type fakeMessageRepo struct {
	listResult []messages.Message
	listErr    error
	listCalls  int
	editFn     func(ctx context.Context, id int, body string) error
	editErr    error
	editCalls  int
	createFn   func(ctx context.Context, body string) (messages.Message, error)
	createErr  error
	createCall int
}

func (f *fakeMessageRepo) List(context.Context) ([]messages.Message, error) {
	f.listCalls++
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listResult, nil
}

func (f *fakeMessageRepo) Create(ctx context.Context, body string) (messages.Message, error) {
	f.createCall++
	if f.createFn != nil {
		return f.createFn(ctx, body)
	}
	if f.createErr != nil {
		return messages.Message{}, f.createErr
	}
	return messages.Message{}, nil
}

func (f *fakeMessageRepo) Edit(ctx context.Context, id int, body string) error {
	f.editCalls++
	if f.editFn != nil {
		return f.editFn(ctx, id, body)
	}
	if f.editErr != nil {
		return f.editErr
	}

	return nil
}

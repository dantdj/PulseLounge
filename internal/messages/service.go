package messages

import (
	"context"
	"errors"
	"strings"
)

var ErrEmptyBody = errors.New("message body cannot be empty")

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) List(ctx context.Context) ([]Message, error) {
	return s.repo.List(ctx)
}

func (s Service) Create(ctx context.Context, body string) (Message, error) {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return Message{}, ErrEmptyBody
	}

	return s.repo.Create(ctx, trimmedBody)
}

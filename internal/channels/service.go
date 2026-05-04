package channels

import (
	"context"
	"errors"
	"strings"
)

var ErrEmptyName = errors.New("channel name cannot be empty")
var ErrChannelNotFound = errors.New("channel not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) List(ctx context.Context) ([]Channel, error) {
	return s.repo.List(ctx)
}

func (s Service) Create(ctx context.Context, name string) (Channel, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return Channel{}, ErrEmptyName
	}

	return s.repo.Create(ctx, trimmedName)
}

func (s Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

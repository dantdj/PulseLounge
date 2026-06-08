package channels

import (
	"context"
	"errors"
	"strings"

	"pulselounge/internal/logging"
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
	result, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	logging.FromContext(ctx).InfoContext(ctx, "listed channels", "channel_count", len(result))
	return result, nil
}

func (s Service) Create(ctx context.Context, name string) (Channel, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return Channel{}, ErrEmptyName
	}

	channel, err := s.repo.Create(ctx, trimmedName)
	if err != nil {
		return Channel{}, err
	}

	logging.FromContext(ctx).InfoContext(ctx, "created channel", "channel_id", channel.ID)
	return channel, nil
}

func (s Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	logging.FromContext(ctx).InfoContext(ctx, "deleted channel", "channel_id", id)
	return nil
}

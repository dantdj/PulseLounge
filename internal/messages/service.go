package messages

import (
	"context"
	"errors"
	"strings"

	"pulselounge/internal/logging"
)

var ErrEmptyBody = errors.New("message body cannot be empty")
var ErrMessageNotFound = errors.New("message not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) ListByChannel(ctx context.Context, channelID int64) ([]Message, error) {
	result, err := s.repo.ListByChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}

	logging.FromContext(ctx).InfoContext(ctx, "listed channel messages", "channel_id", channelID, "message_count", len(result))
	return result, nil
}

func (s Service) CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string, imageKey string) (Message, error) {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return Message{}, ErrEmptyBody
	}

	message, err := s.repo.CreateInChannel(ctx, channelID, authorID, trimmedBody, imageKey)
	if err != nil {
		return Message{}, err
	}

	logging.FromContext(ctx).InfoContext(
		ctx,
		"created message",
		"message_id", message.ID,
		"channel_id", message.ChannelID,
		"author_id", message.AuthorID,
		"has_image", imageKey != "",
	)
	return message, nil
}

func (s Service) Edit(ctx context.Context, id int64, body string) error {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return ErrEmptyBody
	}

	if err := s.repo.Edit(ctx, id, trimmedBody); err != nil {
		return err
	}

	logging.FromContext(ctx).InfoContext(ctx, "edited message", "message_id", id)
	return nil
}

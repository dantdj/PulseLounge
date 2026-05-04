package messages

import (
	"context"
	"errors"
	"strings"
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
	return s.repo.ListByChannel(ctx, channelID)
}

func (s Service) CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string) (Message, error) {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return Message{}, ErrEmptyBody
	}

	return s.repo.CreateInChannel(ctx, channelID, authorID, trimmedBody)
}

func (s Service) Edit(ctx context.Context, id int64, body string) error {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return ErrEmptyBody
	}

	return s.repo.Edit(ctx, id, trimmedBody)
}

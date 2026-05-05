package messages

import "context"

type Repository interface {
	ListByChannel(ctx context.Context, channelID int64) ([]Message, error)
	CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string, imageKey string) (Message, error)
	Edit(ctx context.Context, id int64, body string) error
}

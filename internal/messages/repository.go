package messages

import "context"

type Repository interface {
	List(ctx context.Context) ([]Message, error)
	Create(ctx context.Context, body string) (Message, error)
	Edit(ctx context.Context, id int, body string) error
}

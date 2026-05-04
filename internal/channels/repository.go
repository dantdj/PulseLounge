package channels

import "context"

type Repository interface {
	List(ctx context.Context) ([]Channel, error)
	Create(ctx context.Context, name string) (Channel, error)
	Delete(ctx context.Context, id int64) error
}

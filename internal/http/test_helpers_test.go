package httpapi

import (
	"context"

	"pulselounge/internal/channels"
	"pulselounge/internal/messages"
)

type fakeMessageRepo struct {
	listResult        []messages.Message
	listErr           error
	listCalls         int
	listChannelID     int64
	editFn            func(ctx context.Context, id int64, body string) error
	editErr           error
	editCalls         int
	createFn          func(ctx context.Context, channelID int64, authorID int64, body string) (messages.Message, error)
	createErr         error
	createCall        int
	createChannelID   int64
	createAuthorID    int64
	createRequestBody string
}

func (f *fakeMessageRepo) ListByChannel(_ context.Context, channelID int64) ([]messages.Message, error) {
	f.listCalls++
	f.listChannelID = channelID
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listResult, nil
}

func (f *fakeMessageRepo) CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string) (messages.Message, error) {
	f.createCall++
	f.createChannelID = channelID
	f.createAuthorID = authorID
	f.createRequestBody = body
	if f.createFn != nil {
		return f.createFn(ctx, channelID, authorID, body)
	}
	if f.createErr != nil {
		return messages.Message{}, f.createErr
	}
	return messages.Message{}, nil
}

func (f *fakeMessageRepo) Edit(ctx context.Context, id int64, body string) error {
	f.editCalls++
	if f.editFn != nil {
		return f.editFn(ctx, id, body)
	}
	if f.editErr != nil {
		return f.editErr
	}

	return nil
}

type fakeChannelRepo struct {
	listResult []channels.Channel
	createFn   func(ctx context.Context, name string) (channels.Channel, error)
	deleteFn   func(ctx context.Context, id int64) error
	deleteID   int64
	deleteCall int
}

func (f *fakeChannelRepo) List(context.Context) ([]channels.Channel, error) {
	return f.listResult, nil
}

func (f *fakeChannelRepo) Create(ctx context.Context, name string) (channels.Channel, error) {
	if f.createFn != nil {
		return f.createFn(ctx, name)
	}
	return channels.Channel{}, nil
}

func (f *fakeChannelRepo) Delete(ctx context.Context, id int64) error {
	f.deleteCall++
	f.deleteID = id
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func testChannelService() channels.Service {
	return channels.NewService(&fakeChannelRepo{})
}

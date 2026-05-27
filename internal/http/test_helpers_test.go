package httpapi

import (
	"context"

	"pulselounge/internal/channels"
	"pulselounge/internal/media"
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
	createFn          func(ctx context.Context, channelID int64, authorID int64, body string, imageKey string) (messages.Message, error)
	createErr         error
	createCall        int
	createChannelID   int64
	createAuthorID    int64
	createRequestBody string
	createImageKey    string
}

func (f *fakeMessageRepo) ListByChannel(_ context.Context, channelID int64) ([]messages.Message, error) {
	f.listCalls++
	f.listChannelID = channelID
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listResult, nil
}

func (f *fakeMessageRepo) CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string, imageKey string) (messages.Message, error) {
	f.createCall++
	f.createChannelID = channelID
	f.createAuthorID = authorID
	f.createRequestBody = body
	f.createImageKey = imageKey
	if f.createFn != nil {
		return f.createFn(ctx, channelID, authorID, body, imageKey)
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

func testMediaStore() media.Store {
	return &fakeMediaStore{}
}

type fakeMediaStore struct {
	saveFn func(id string, contentType string, data []byte) (string, error)
}

func (f *fakeMediaStore) Save(id string, contentType string, data []byte) (string, error) {
	if f.saveFn != nil {
		return f.saveFn(id, contentType, data)
	}
	return "https://example.com/" + id, nil
}

func (f *fakeMediaStore) PublicURL(id string) string {
	return "https://example.com/" + id
}

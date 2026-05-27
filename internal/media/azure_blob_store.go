package media

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// TODO: Should pull these from some kind of config at some point
const (
	thumbnailFolder = "thumbs/"
	retryAttempts   = 3
	retryBaseDelay  = 200 * time.Millisecond
)

func NewAzureBlobStore(connectionString string, containerName string, publicBaseURL string) (*AzureBlobStore, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	store := &AzureBlobStore{
		client:        client,
		containerName: strings.Trim(containerName, "/"),
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}

	return store, nil
}

type AzureBlobStore struct {
	client        *azblob.Client
	containerName string
	publicBaseURL string
	containerMu   sync.Mutex
	containerOK   bool
}

func (s *AzureBlobStore) Save(id string, contentType string, file []byte) (string, error) {
	ctx := context.Background()
	if err := s.ensureContainer(ctx); err != nil {
		return "", err
	}

	err := withRetry(ctx, func(ctx context.Context) error {
		_, err := s.client.UploadBuffer(ctx, s.containerName, id, file, &azblob.UploadBufferOptions{
			HTTPHeaders: &blob.HTTPHeaders{
				BlobContentType: to.Ptr(contentType),
			},
		})
		return err
	})
	if err != nil {
		return "", err
	}

	return s.PublicURL(id), nil
}

func (s *AzureBlobStore) ensureContainer(ctx context.Context) error {
	s.containerMu.Lock()
	defer s.containerMu.Unlock()

	if s.containerOK {
		return nil
	}

	access := container.PublicAccessTypeBlob
	_, err := s.client.CreateContainer(ctx, s.containerName, &azblob.CreateContainerOptions{
		Access: &access,
	})
	if err != nil {
		var responseErr *azcore.ResponseError
		if !errors.As(err, &responseErr) || responseErr.ErrorCode != "ContainerAlreadyExists" {
			return err
		}
	}

	s.containerOK = true
	return nil
}

func (s *AzureBlobStore) PublicURL(id string) string {
	return s.publicBaseURL + "/" + id
}

func ThumbnailKey(id string) string {
	return thumbnailFolder + id
}

func withRetry(ctx context.Context, fn func(context.Context) error) error {
	if retryAttempts <= 0 {
		return errors.New("attempts must be greater than 0")
	}

	delay := retryBaseDelay
	for attempt := 1; attempt <= retryAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		if !shouldRetryAzureError(err) {
			return err
		}

		// If this was the last attempt, return the error immediately instead of waiting for the delay
		if attempt == retryAttempts {
			return err
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		delay *= 2
	}

	return nil
}

func shouldRetryAzureError(err error) bool {
	var responseErr *azcore.ResponseError
	if errors.As(err, &responseErr) {
		switch responseErr.StatusCode {
		case http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		default:
			return false
		}
	}

	// Check for network errors that may indicate a retryable,
	// not-necessarily-Azure issue (e.g. DNS failure)
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}

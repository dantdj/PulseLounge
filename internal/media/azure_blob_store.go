package media

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// TODO: Should pull this from some kind of config at some point
var thumbnailFolder = "thumbs/"

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

	_, err := s.client.UploadBuffer(ctx, s.containerName, id, file, &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
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

func (s *AzureBlobStore) Exists(id string) (bool, error) {
	blobClient := s.client.ServiceClient().NewContainerClient(s.containerName).NewBlobClient(id)
	_, err := blobClient.GetProperties(context.Background(), nil)
	if err != nil {
		var responseErr *azcore.ResponseError
		if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *AzureBlobStore) PublicURL(id string) string {
	return s.publicBaseURL + "/" + id
}

func ThumbnailKey(id string) string {
	return thumbnailFolder + id
}

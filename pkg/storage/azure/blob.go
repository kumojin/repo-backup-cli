package azure

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/kumojin/repo-backup-cli/pkg/storage"
)

type defaultBlobRepository struct {
	cfg    *config.Config
	client *azblob.Client
}

func NewBlobRepository(cfg *config.Config, client *azblob.Client) storage.BlobRepository {
	return defaultBlobRepository{cfg: cfg, client: client}
}

func (r defaultBlobRepository) Upload(ctx context.Context, blobName string, in io.Reader) (string, error) {
	_, err := r.client.UploadStream(ctx, r.cfg.AzureStorageConfig.ContainerName, blobName, in, nil)
	if err != nil {
		return "", err
	}

	return r.getUrl(blobName)
}

func (r defaultBlobRepository) getUrl(blobName string) (string, error) {
	url, err := url.JoinPath(r.cfg.AzureStorageConfig.AccountUrl, r.cfg.AzureStorageConfig.ContainerName, blobName)
	if err != nil {
		return "", fmt.Errorf("failed to construct blob URL: %w", err)
	}

	return url, nil
}

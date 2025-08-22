package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/kumojin/repo-backup-cli/pkg/storage"
	"github.com/minio/minio-go/v7"
)

type defaultBlobRepository struct {
	cfg    *config.Config
	client *minio.Client
}

func NewBlobRepository(cfg *config.Config, client *minio.Client) storage.BlobRepository {
	return defaultBlobRepository{cfg: cfg, client: client}
}

func (r defaultBlobRepository) Upload(ctx context.Context, blobName string, in io.Reader) (string, error) {
	info, err := r.client.PutObject(ctx, r.cfg.ObjectStorageConfig.BucketName, blobName, in, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to object storage: %w", err)
	}

	return info.Location, nil
}

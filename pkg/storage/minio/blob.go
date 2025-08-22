package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

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
	size, err := getSize(in)
	if err != nil {
		return "", fmt.Errorf("failed to get size of input stream: %w", err)
	}

	_, err = r.client.PutObject(ctx, r.cfg.ObjectStorageConfig.BucketName, blobName, in, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to object storage: %w", err)
	}

	return r.getUrl(blobName)
}

func (r defaultBlobRepository) getUrl(blobName string) (string, error) {
	scheme := "http"
	if r.cfg.ObjectStorageConfig.UseSSL {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, r.cfg.ObjectStorageConfig.Endpoint)
	url, err := url.JoinPath(baseURL, r.cfg.ObjectStorageConfig.BucketName, blobName)
	if err != nil {
		return "", fmt.Errorf("failed to construct MinIO URL: %w", err)
	}

	return url, nil
}

func getSize(stream io.Reader) (int64, error) {
	seeker, ok := stream.(io.Seeker)
	if !ok {
		return 0, errors.New("cannot cast reader to seeker")
	}

	size, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return size, nil
}

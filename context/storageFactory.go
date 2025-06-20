package context

import (
	"fmt"

	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/kumojin/repo-backup-cli/pkg/storage"
	"github.com/kumojin/repo-backup-cli/pkg/storage/azure"
	"github.com/kumojin/repo-backup-cli/pkg/storage/minio"
)

func NewBlobRepository(cfg *config.Config) (storage.BlobRepository, error) {
	switch cfg.StorageBackend {
	case config.StorageBackendObject:
		minioClient, err := GetMinioClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to get MinIO client: %w", err)
		}
		return minio.NewBlobRepository(cfg, minioClient), nil

	case config.StorageBackendAzure:
		azureClient, err := GetAzureBlobClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to get Azure client: %w", err)
		}
		return azure.NewBlobRepository(cfg, azureClient), nil

	default:
		return nil, fmt.Errorf("unsupported storage backend: %s (supported: %s, %s)", cfg.StorageBackend, config.StorageBackendAzure, config.StorageBackendObject)
	}
}

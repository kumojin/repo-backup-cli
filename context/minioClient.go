package context

import (
	"fmt"

	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
)

func GetMinioClient(cfg *config.Config) (*minio.Client, error) {
	if minioClient == nil {
		client, err := minio.New(cfg.ObjectStorageConfig.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.ObjectStorageConfig.AccessKey, cfg.ObjectStorageConfig.SecretKey, ""),
			Secure: cfg.ObjectStorageConfig.UseSSL,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create MinIO client: %w", err)
		}

		minioClient = client
	}

	return minioClient, nil
}

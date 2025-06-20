package context

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

var (
	azureClient *azblob.Client
)

func GetAzureBlobClient(cfg *config.Config) (*azblob.Client, error) {
	if azureClient == nil {
		credentials, err := azblob.NewSharedKeyCredential(cfg.AzureStorageConfig.AccountName, cfg.AzureStorageConfig.ApiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure SharedKeyCredential: %w", err)
		}

		azureClient, err = azblob.NewClientWithSharedKeyCredential(cfg.AzureStorageConfig.AccountUrl, credentials, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
		}
	}

	return azureClient, nil
}

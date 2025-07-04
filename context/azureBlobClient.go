package context

import (
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
			return nil, err
		}

		azureClient, err = azblob.NewClientWithSharedKeyCredential(cfg.AzureStorageConfig.AccountUrl, credentials, nil)
		if err != nil {
			return nil, err
		}
	}

	return azureClient, nil
}

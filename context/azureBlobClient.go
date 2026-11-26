package context

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

var (
	azureClient *azblob.Client
)

func GetAzureBlobClient(cfg *config.Config) (*azblob.Client, error) {
	if azureClient == nil {
		client, err := newAzureBlobClient(cfg.AzureStorageConfig)
		if err != nil {
			return nil, err
		}
		azureClient = client
	}

	return azureClient, nil
}

func newAzureBlobClient(storageConfig config.AzureStorageConfig) (*azblob.Client, error) {
	if storageConfig.UsesSharedKey() {
		return newSharedKeyBlobClient(storageConfig)
	}

	return newAzureIdentityBlobClient(storageConfig)
}

func newSharedKeyBlobClient(storageConfig config.AzureStorageConfig) (*azblob.Client, error) {
	credential, err := azblob.NewSharedKeyCredential(storageConfig.AccountName, storageConfig.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure SharedKeyCredential: %w", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(storageConfig.AccountUrl, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	return client, nil
}

func newAzureIdentityBlobClient(storageConfig config.AzureStorageConfig) (*azblob.Client, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure DefaultAzureCredential: %w", err)
	}

	client, err := azblob.NewClient(storageConfig.AccountUrl, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	return client, nil
}

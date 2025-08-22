package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	StorageBackendAzure  = "azure"
	StorageBackendObject = "object"
)

type AzureStorageConfig struct {
	AccountName   string
	ApiKey        string
	AccountUrl    string
	ContainerName string
}

func newAzureStorageConfig() (AzureStorageConfig, error) {
	accountName := viper.GetString(azureStorageAccountNameKey)
	apiKey := viper.GetString(azureStorageApiKeyKey)
	accountUrl := viper.GetString(azureStorageAccountUrlKey)
	containerName := viper.GetString(azureStorageContainerNameKey)

	if accountName == "" || apiKey == "" || accountUrl == "" || containerName == "" {
		return AzureStorageConfig{}, fmt.Errorf("azure Storage configuration is incomplete")
	}

	return AzureStorageConfig{
		AccountName:   accountName,
		ApiKey:        apiKey,
		AccountUrl:    accountUrl,
		ContainerName: containerName,
	}, nil
}

type ObjectStorageConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

func newObjectStorageConfig() (ObjectStorageConfig, error) {
	endpoint := viper.GetString(objectStorageEndpointKey)
	accessKey := viper.GetString(objectStorageAccessKeyKey)
	secretKey := viper.GetString(objectStorageSecretKeyKey)
	bucketName := viper.GetString(objectStorageBucketNameKey)
	useSSL := viper.GetBool(objectStorageUseSSLKey)

	if endpoint == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		return ObjectStorageConfig{}, fmt.Errorf("object storage configuration is incomplete")
	}

	return ObjectStorageConfig{
		Endpoint:   endpoint,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
		UseSSL:     useSSL,
	}, nil
}

func createStorageConfigs(storageBackend string) (AzureStorageConfig, ObjectStorageConfig, error) {
	var azureStorageConfig AzureStorageConfig
	var objectStorageConfig ObjectStorageConfig
	var err error

	switch storageBackend {
	case StorageBackendAzure:
		azureStorageConfig, err = newAzureStorageConfig()
		if err != nil {
			return azureStorageConfig, objectStorageConfig, fmt.Errorf("failed to create Azure storage config: %w", err)
		}
	case StorageBackendObject:
		objectStorageConfig, err = newObjectStorageConfig()
		if err != nil {
			return azureStorageConfig, objectStorageConfig, fmt.Errorf("failed to create object storage config: %w", err)
		}
	default:
		return azureStorageConfig, objectStorageConfig, fmt.Errorf("unsupported storage backend: %s (supported: %s, %s)", storageBackend, StorageBackendAzure, StorageBackendObject)
	}

	return azureStorageConfig, objectStorageConfig, nil
}

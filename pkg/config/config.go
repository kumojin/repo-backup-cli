package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	azureStorageAccountNameKey   = "AZURE_STORAGE_ACCOUNT_NAME"
	azureStorageApiKeyKey        = "AZURE_STORAGE_API_KEY"
	azureStorageAccountUrlKey    = "AZURE_STORAGE_ACCOUNT_URL"
	azureStorageContainerNameKey = "AZURE_STORAGE_CONTAINER_NAME"
	objectStorageEndpointKey     = "OBJECT_STORAGE_ENDPOINT"
	objectStorageAccessKeyKey    = "OBJECT_STORAGE_ACCESS_KEY"
	objectStorageSecretKeyKey    = "OBJECT_STORAGE_SECRET_KEY"
	objectStorageBucketNameKey   = "OBJECT_STORAGE_BUCKET_NAME"
	objectStorageUseSSLKey       = "OBJECT_STORAGE_USE_SSL"
	storageBackendKey            = "STORAGE_BACKEND"
	githubTokenKey               = "CLI_GITHUB_TOKEN"
	sentryDsnKey                 = "SENTRY_DSN"
)

type SentryConfig struct {
	Dsn string
}

func NewSentryConfig() SentryConfig {
	return SentryConfig{
		Dsn: viper.GetString(sentryDsnKey),
	}
}

type Config struct {
	AzureStorageConfig  AzureStorageConfig
	ObjectStorageConfig ObjectStorageConfig
	SentryConfig        SentryConfig
	GitHubToken         string
	Organization        string
	StorageBackend      string
}

func New(filepath string) (*Config, error) {
	viper.SetConfigName(filepath)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	token := viper.GetString(githubTokenKey)
	if token == "" {
		return nil, fmt.Errorf("github token is not set in the configuration file")
	}

	storageBackend := viper.GetString(storageBackendKey)
	if storageBackend == "" {
		storageBackend = StorageBackendAzure // Defaults to Azure blob storage
	}

	azureStorageConfig, objectStorageConfig, err := createStorageConfigs(storageBackend)
	if err != nil {
		return nil, err
	}

	return &Config{
		AzureStorageConfig:  azureStorageConfig,
		ObjectStorageConfig: objectStorageConfig,
		GitHubToken:         token,
		SentryConfig:        NewSentryConfig(),
		StorageBackend:      storageBackend,
	}, nil
}

func (c *Config) WithOrganization(organization string) *Config {
	c.Organization = organization

	return c
}

func (c *Config) GetSentryConfig() SentryConfig {
	return c.SentryConfig
}

func (c *Config) IsSentryEnabled() bool {
	return c.SentryConfig.Dsn != ""
}

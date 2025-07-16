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
	githubTokenKey               = "GITHUB_TOKEN"
)

type AzureStorageConfig struct {
	AccountName   string
	ApiKey        string
	AccountUrl    string
	ContainerName string
}

func NewAzureStorageConfig() (AzureStorageConfig, error) {
	accountName := viper.GetString(azureStorageAccountNameKey)
	apiKey := viper.GetString(azureStorageApiKeyKey)
	accountUrl := viper.GetString(azureStorageAccountUrlKey)
	containerName := viper.GetString(azureStorageContainerNameKey)

	if accountName == "" || apiKey == "" || accountUrl == "" || containerName == "" {
		return AzureStorageConfig{}, fmt.Errorf("Azure Storage configuration is incomplete")
	}

	return AzureStorageConfig{
		AccountName:   accountName,
		ApiKey:        apiKey,
		AccountUrl:    accountUrl,
		ContainerName: containerName,
	}, nil
}

type Config struct {
	AzureStorageConfig AzureStorageConfig
	GitHubToken        string
	Organization       string
}

func New(filepath string) (*Config, error) {
	viper.SetConfigName(filepath)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file: %w", err)
	}

	token := viper.GetString(githubTokenKey)
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set in the configuration file")
	}

	azureStorageConfig, err := NewAzureStorageConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AzureStorageConfig: azureStorageConfig,
		GitHubToken:        token,
	}, nil
}

func (c *Config) WithOrganization(organization string) *Config {
	updatedConfig := *c

	updatedConfig.Organization = organization

	return &updatedConfig
}

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
	githubTokenKey               = "CLI_GITHUB_TOKEN"
	sentryDsnKey                 = "SENTRY_DSN"
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
		return AzureStorageConfig{}, fmt.Errorf("azure Storage configuration is incomplete")
	}

	return AzureStorageConfig{
		AccountName:   accountName,
		ApiKey:        apiKey,
		AccountUrl:    accountUrl,
		ContainerName: containerName,
	}, nil
}

type SentryConfig struct {
	Dsn string
}

func NewSentryConfig() (SentryConfig, error) {
	dsn := viper.GetString(sentryDsnKey)

	if dsn == "" {
		return SentryConfig{}, fmt.Errorf("sentry config is incomplete")
	}

	return SentryConfig{
		Dsn: dsn,
	}, nil
}

type Config struct {
	AzureStorageConfig AzureStorageConfig
	SentryConfig       SentryConfig
	GitHubToken        string
	Organization       string
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

	azureStorageConfig, err := NewAzureStorageConfig()
	if err != nil {
		return nil, err
	}

	sentryConfig, err := NewSentryConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AzureStorageConfig: azureStorageConfig,
		GitHubToken:        token,
		SentryConfig:       sentryConfig,
	}, nil
}

func (c *Config) WithOrganization(organization string) *Config {
	c.Organization = organization

	return c
}

func (c *Config) GetSentryConfig() SentryConfig {
	return c.SentryConfig
}

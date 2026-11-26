package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setAzureStorageEnv(t *testing.T, apiKey string) {
	t.Helper()
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set(azureStorageAccountNameKey, "account")
	viper.Set(azureStorageAccountUrlKey, "https://account.blob.core.windows.net")
	viper.Set(azureStorageContainerNameKey, "backup")
	viper.Set(azureStorageApiKeyKey, apiKey)
}

func TestAzureStorageConfigWithApiKeyUsesSharedKey(t *testing.T) {
	setAzureStorageEnv(t, "key")

	cfg, err := newAzureStorageConfig()

	require.NoError(t, err)
	assert.True(t, cfg.UsesSharedKey())
}

func TestAzureStorageConfigWithoutApiKeyUsesAzureIdentity(t *testing.T) {
	setAzureStorageEnv(t, "")

	cfg, err := newAzureStorageConfig()

	require.NoError(t, err)
	assert.False(t, cfg.UsesSharedKey())
}

func TestAzureStorageConfigRequiresAccountUrl(t *testing.T) {
	setAzureStorageEnv(t, "")
	viper.Set(azureStorageAccountUrlKey, "")

	_, err := newAzureStorageConfig()

	assert.Error(t, err)
}

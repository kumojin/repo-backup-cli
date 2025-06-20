package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const githubTokenKey = "github_token"

type Config struct {
	GitHubToken  string
	Organization string
}

func New(filepath string) (*Config, error) {
	viper.SetConfigName(filepath)
	viper.SetConfigType("env")

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file: %w", err)
	}

	token := viper.GetString(githubTokenKey)
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set in the configuration file")
	}

	return &Config{
		GitHubToken: token,
	}, nil
}

func (c *Config) WithOrganization(organization string) *Config {
	updatedConfig := *c

	updatedConfig.Organization = organization

	return &updatedConfig
}

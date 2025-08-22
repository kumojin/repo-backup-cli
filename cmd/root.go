package cmd

import (
	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	rootConfig     *config.Config
	organization   string
	configFilepath string
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rbk",
		Short:         "CLI tool to backup private repositories from one a github organization to a remote file storage service",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVarP(&configFilepath, "config", "c", ".env", "Path to environment configuration file")
	cmd.PersistentFlags().StringVarP(&organization, "organization", "o", "Kumojin", "GitHub organization to use")

	cmd.AddCommand(ReposCommand())
	cmd.AddCommand(BackupCommand())

	return cmd
}

func GetConfig() (*config.Config, error) {
	if rootConfig != nil {
		return rootConfig, nil
	}

	var err error
	rootConfig, err = config.New(configFilepath)
	if err != nil {
		return nil, err
	}

	return rootConfig.WithOrganization(organization), nil
}

package cmd

import (
	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	rootConfig     *config.Config
	organization   string
	configFilepath string
)

func RootCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:               "rbk",
		Short:             "CLI tool to backup private repositories from one a github organization to a remote file storage service",
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: preRun,
	}

	cmd.PersistentFlags().StringVarP(&configFilepath, "config", "c", ".env", "Path to environment configuration file")
	cmd.PersistentFlags().StringVarP(&organization, "organization", "o", "", "GitHub organization to use")

	err := cmd.MarkPersistentFlagRequired("organization")
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(ReposCommand())
	cmd.AddCommand(BackupCommand())

	return cmd, nil
}

func preRun(_ *cobra.Command, _ []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	if !cfg.IsSentryEnabled() {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.GetSentryConfig().Dsn,
		EnableLogs:       true,
		SendDefaultPII:   true,
		AttachStacktrace: true,
	})
}

func getConfig() (*config.Config, error) {
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

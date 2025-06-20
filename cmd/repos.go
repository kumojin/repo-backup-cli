package cmd

import (
	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/uc"

	"github.com/spf13/cobra"
)

func ReposCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List all repositories in an organization (defaults to `Kumojin` organization)",
		RunE:  runReposCommand,
	}

	return cmd
}

func runReposCommand(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	client := appContext.GetGitHubClient(cfg)

	usecase := uc.NewListPrivateReposUseCase(client)

	repos, err := usecase.Do(cfg.Organization)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		println(*repo.Name)
	}

	return nil
}

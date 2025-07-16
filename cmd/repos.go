package cmd

import (
	"context"

	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/github"
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

	githubClient := github.NewClient(appContext.GetGithubClient(cfg))

	usecase := uc.NewListPrivateReposUseCase(githubClient)

	ctx := context.Background()

	repos, err := usecase.Do(ctx, cfg.Organization)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		println(*repo.Name)
	}

	return nil
}

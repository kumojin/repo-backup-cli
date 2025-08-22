package cmd

import (
	"context"
	"log/slog"

	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/kumojin/repo-backup-cli/pkg/logging"
	"github.com/kumojin/repo-backup-cli/pkg/uc"

	"github.com/spf13/cobra"
)

func ReposCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List all repositories in an organization (defaults to `Kumojin` organization)",
		Run:   runReposCommand,
	}

	return cmd
}

func runReposCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	logger := logging.NewLogger(ctx)

	cfg, err := GetConfig()
	if err != nil {
		logger.Error("could not get config", slog.Any("error", err))
		return
	}

	logger = logger.With(
		slog.String("organization", cfg.Organization),
	)

	githubClient := github.NewClient(appContext.GetGithubClient(cfg))

	usecase := uc.NewListPrivateReposUseCase(githubClient)

	repos, err := usecase.Do(ctx, cfg.Organization)
	if err != nil {
		logger.Error("could not list repositories", slog.Any("error", err))
		return
	}

	for _, repo := range repos {
		println(*repo.Name)
	}
}

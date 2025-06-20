package cmd

import (
	"context"

	"github.com/google/go-github/v72/github"
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

	client := github.NewClient(nil).WithAuthToken(cfg.GitHubToken)

	repos, _, err := client.Repositories.ListByOrg(context.TODO(), cfg.Organization, &github.RepositoryListByOrgOptions{
		Type: "private",
	})

	if err != nil {
		return err
	}

	for _, repo := range repos {
		println(*repo.Name)
	}

	return nil
}

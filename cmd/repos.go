package cmd

import (
	"context"

	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

var (
	reposOrganization string
)

func ReposCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List all repositories in an organization (defaults to `Kumojin` organization)",
		RunE:  runReposCommand,
	}

	cmd.Flags().StringVarP(&reposOrganization, "organization", "o", "Kumojin", "GitHub organization to list repositories from")

	return cmd
}

func runReposCommand(cmd *cobra.Command, args []string) error {
	client := github.NewClient(nil).WithAuthToken("your_token")

	repos, _, err := client.Repositories.ListByOrg(context.TODO(), reposOrganization, &github.RepositoryListByOrgOptions{
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

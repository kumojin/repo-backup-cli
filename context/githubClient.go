package context

import (
	"github.com/google/go-github/v89/github"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

var (
	githubClient *github.Client
)

func GetGithubClient(cfg *config.Config) (*github.Client, error) {
	if githubClient == nil {
		client, err := github.NewClient(github.WithAuthToken(cfg.GitHubToken))
		if err != nil {
			return nil, err
		}
		githubClient = client
	}

	return githubClient, nil
}

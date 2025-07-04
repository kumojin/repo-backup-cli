package context

import (
	"github.com/google/go-github/v73/github"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

var (
	githubClient *github.Client
)

func GetGithubClient(cfg *config.Config) *github.Client {
	if githubClient == nil {
		githubClient = github.NewClient(nil).WithAuthToken(cfg.GitHubToken)
	}

	return githubClient
}

package uc

import (
	"context"

	gh "github.com/google/go-github/v78/github"
	"github.com/kumojin/repo-backup-cli/pkg/github"
)

type ListPrivateReposUseCase interface {
	Do(ctx context.Context, organization string) ([]gh.Repository, error)
}

type listPrivateReposUseCase struct {
	githubClient github.Client
}

func NewListPrivateReposUseCase(client github.Client) ListPrivateReposUseCase {
	return &listPrivateReposUseCase{
		githubClient: client,
	}
}

func (uc *listPrivateReposUseCase) Do(ctx context.Context, organization string) ([]gh.Repository, error) {
	var repos []*gh.Repository

	repos, err := uc.githubClient.ListOrgRepos(ctx, organization, "private")
	if err != nil {
		return nil, err
	}

	var filteredRepos []gh.Repository
	for _, repo := range repos {
		if !repo.GetArchived() {
			filteredRepos = append(filteredRepos, *repo)
		}
	}

	return filteredRepos, nil
}

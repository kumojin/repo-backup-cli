package uc

import (
	"github.com/google/go-github/v72/github"
)

type ListPrivateReposUseCase interface {
	Do(organization string) ([]github.Repository, error)
}

type listPrivateReposUseCase struct {
	githubClient *github.Client
}

func NewListPrivateReposUseCase(client *github.Client) ListPrivateReposUseCase {
	return &listPrivateReposUseCase{
		githubClient: client,
	}
}

func (uc *listPrivateReposUseCase) Do(organization string) ([]github.Repository, error) {
	repos, _, err := uc.githubClient.Repositories.ListByOrg(
		nil,
		organization,
		&github.RepositoryListByOrgOptions{
			Type: "private",
		},
	)
	if err != nil {
		return nil, err
	}

	var filteredRepos []github.Repository
	for _, repo := range repos {
		if !repo.GetArchived() {
			filteredRepos = append(filteredRepos, *repo)
		}
	}

	return filteredRepos, nil
}

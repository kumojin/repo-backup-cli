package uc

import (
	"context"

	"github.com/google/go-github/v72/github"
)

const maxPerPage = 100

type ListPrivateReposUseCase interface {
	Do(ctx context.Context, organization string) ([]github.Repository, error)
}

type listPrivateReposUseCase struct {
	githubClient *github.Client
}

func NewListPrivateReposUseCase(client *github.Client) ListPrivateReposUseCase {
	return &listPrivateReposUseCase{
		githubClient: client,
	}
}

func (uc *listPrivateReposUseCase) Do(ctx context.Context, organization string) ([]github.Repository, error) {
	var repos []*github.Repository
	hasMore := true
	page := 1

	for hasMore {
		pageRepos, res, err := uc.githubClient.Repositories.ListByOrg(
			ctx,
			organization,
			&github.RepositoryListByOrgOptions{
				Type: "private",
				ListOptions: github.ListOptions{
					PerPage: maxPerPage,
					Page:    page,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		repos = append(repos, pageRepos...)

		if res.NextPage > page {
			page = res.NextPage
		} else {
			hasMore = false
		}
	}

	var filteredRepos []github.Repository
	for _, repo := range repos {
		if !repo.GetArchived() {
			filteredRepos = append(filteredRepos, *repo)
		}
	}

	return filteredRepos, nil
}

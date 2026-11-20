package uc

import (
	"context"
	"errors"
	"testing"

	gh "github.com/google/go-github/v93/github"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListPrivateReposUseCase_SuccessfullyListNonArchivedRepos(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return(
			[]*gh.Repository{
				{
					Name:     new("repo1"),
					Archived: new(false),
					Private:  new(true),
				},
				{
					Name:     new("repo2"),
					Archived: new(true), // Archived repo, should be filtered out
					Private:  new(true),
				},
				{
					Name:     new("repo3"),
					Archived: new(false),
					Private:  new(true),
				},
			},
			nil,
		)

	useCase := NewListPrivateReposUseCase(mockClient)

	// When
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Then
	assert.NoError(t, err)

	expectedRepos := []gh.Repository{
		{
			Name:     new("repo1"),
			Archived: new(false),
			Private:  new(true),
		},
		{
			Name:     new("repo3"),
			Archived: new(false),
			Private:  new(true),
		},
	}
	assert.Equal(t, expectedRepos, repos)
}

func TestListPrivateReposUseCase_ErrorFromGitHubClient(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	githubApiError := errors.New("github API error")
	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return(nil, githubApiError)

	useCase := NewListPrivateReposUseCase(mockClient)

	// When
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Then
	assert.ErrorIs(t, err, githubApiError)
	assert.Nil(t, repos)
}

func TestListPrivateReposUseCase_NoRepositoriesFound(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return([]*gh.Repository{}, nil)

	useCase := NewListPrivateReposUseCase(mockClient)

	// When
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Then
	assert.NoError(t, err)
	assert.Empty(t, repos)
}

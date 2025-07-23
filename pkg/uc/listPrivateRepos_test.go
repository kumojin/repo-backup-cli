package uc

import (
	"context"
	"errors"
	"testing"

	gh "github.com/google/go-github/v73/github"
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
					Name:     gh.Ptr("repo1"),
					Archived: gh.Ptr(false),
					Private:  gh.Ptr(true),
				},
				{
					Name:     gh.Ptr("repo2"),
					Archived: gh.Ptr(true), // Archived repo, should be filtered out
					Private:  gh.Ptr(true),
				},
				{
					Name:     gh.Ptr("repo3"),
					Archived: gh.Ptr(false),
					Private:  gh.Ptr(true),
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
			Name:     gh.Ptr("repo1"),
			Archived: gh.Ptr(false),
			Private:  gh.Ptr(true),
		},
		{
			Name:     gh.Ptr("repo3"),
			Archived: gh.Ptr(false),
			Private:  gh.Ptr(true),
		},
	}
	assert.Equal(t, expectedRepos, repos)
}

func TestListPrivateReposUseCase_ErrorFromGitHubClient(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return(nil, errors.New("github API error"))

	useCase := NewListPrivateReposUseCase(mockClient)

	// When
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Then
	assert.Error(t, err)
	assert.Equal(t, "github API error", err.Error())
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

package uc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v73/github"
	"github.com/kumojin/repo-backup-cli/pkg/github/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListPrivateReposUseCase_SuccessfullyListNonArchivedRepos(t *testing.T) {
	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Setup mock expectations
	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return(
			[]*github.Repository{
				{
					Name:     github.Ptr("repo1"),
					Archived: github.Ptr(false),
					Private:  github.Ptr(true),
				},
				{
					Name:     github.Ptr("repo2"),
					Archived: github.Ptr(true), // Archived repo, should be filtered out
					Private:  github.Ptr(true),
				},
				{
					Name:     github.Ptr("repo3"),
					Archived: github.Ptr(false),
					Private:  github.Ptr(true),
				},
			},
			nil,
		)

	// Create use case with mock client
	useCase := NewListPrivateReposUseCase(mockClient)

	// Execute the use case
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Assertions
	assert.NoError(t, err)

	expectedRepos := []github.Repository{
		{
			Name:     github.Ptr("repo1"),
			Archived: github.Ptr(false),
			Private:  github.Ptr(true),
		},
		{
			Name:     github.Ptr("repo3"),
			Archived: github.Ptr(false),
			Private:  github.Ptr(true),
		},
	}
	assert.Equal(t, expectedRepos, repos)
}

func TestListPrivateReposUseCase_ErrorFromGitHubClient(t *testing.T) {
	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Setup mock expectations
	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return(nil, errors.New("github API error"))

	// Create use case with mock client
	useCase := NewListPrivateReposUseCase(mockClient)

	// Execute the use case
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "github API error", err.Error())
	assert.Nil(t, repos)
}

func TestListPrivateReposUseCase_NoRepositoriesFound(t *testing.T) {
	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Setup mock expectations
	mockClient.EXPECT().
		ListOrgRepos(mock.Anything, "kumojin", "private").
		Return([]*github.Repository{}, nil)

	// Create use case with mock client
	useCase := NewListPrivateReposUseCase(mockClient)

	// Execute the use case
	repos, err := useCase.Do(context.Background(), "kumojin")

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, repos)
}

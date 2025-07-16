package uc

import (
	"context"
	"errors"
	"testing"
	"time"

	gh "github.com/google/go-github/v73/github"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSaveBackupFunc is a mock implementation of SaveBackupFunc
type MockSaveBackupFunc struct {
	mock.Mock
}

func (m *MockSaveBackupFunc) Do(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

// createBackupTestMocks contains all the mocks used in tests
type createBackupTestMocks struct {
	githubClient              *github.MockClient
	listPrivateRepos          *MockListPrivateReposUseCase
	getOrganizationArchiveUrl *MockGetOrganizationArchiveUrlUseCase
	saveBackupFunc            func(url string) (string, error)
	saveBackupMock            *MockSaveBackupFunc
}

// newCreateBackupTestMocks creates and returns all the mocks needed for testing
func newCreateBackupTestMocks(t *testing.T) *createBackupTestMocks {
	mockGithubClient := github.NewMockClient(t)
	mockListPrivateRepos := NewMockListPrivateReposUseCase(t)
	mockGetArchiveUrl := NewMockGetOrganizationArchiveUrlUseCase(t)

	// Create mock save backup function
	mockSaveBackup := new(MockSaveBackupFunc)
	saveBackupFunc := func(url string) (string, error) {
		return mockSaveBackup.Do(url)
	}

	return &createBackupTestMocks{
		githubClient:              mockGithubClient,
		listPrivateRepos:          mockListPrivateRepos,
		getOrganizationArchiveUrl: mockGetArchiveUrl,
		saveBackupFunc:            saveBackupFunc,
		saveBackupMock:            mockSaveBackup,
	}
}

// createUseCase creates a new use case instance with the provided mocks
func (m *createBackupTestMocks) createUseCase() CreateBackupUseCase {
	return NewCreateBackupUseCase(m.githubClient, m.listPrivateRepos, m.getOrganizationArchiveUrl).
		WithPollingInterval(1) // Use 1ns for faster tests
}

func TestCreateBackupUseCase_Success(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("exported"),
	}
	archiveURL := "https://api.github.com/archive/kumojin/12345.zip"
	savePath := "/tmp/backup.zip"

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	// Setup expectations for getting archive URL
	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(archiveURL, nil)

	// Setup expectations for saving backup
	mocks.saveBackupMock.On("Do", archiveURL).Return(savePath, nil)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, savePath, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_SuccessOnSecondCall(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	pendingMigration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("pending"),
	}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("exported"),
	}
	archiveURL := "https://api.github.com/archive/kumojin/12345.zip"
	savePath := "/tmp/backup.zip"

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(pendingMigration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(pendingMigration, nil).
		Once()

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil).
		Once()

	// Setup expectations for getting archive URL
	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(archiveURL, nil)

	// Setup expectations for saving backup
	mocks.saveBackupMock.On("Do", archiveURL).Return(savePath, nil)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, savePath, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_ListRepositoriesError(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	expectedError := errors.New("failed to list repositories")

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return([]gh.Repository{}, expectedError)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list private repositories")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_StartMigrationError(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	expectedError := errors.New("failed to start migration")

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(nil, expectedError)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start migration")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_FailedMigration(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("pending"),
	}
	failedMigration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("failed"),
	}

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(failedMigration, nil)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "migration failed")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_GetMigrationStatusError(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("pending"),
	}
	expectedError := errors.New("failed to get migration status")

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(nil, expectedError)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get migration status")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_GetArchiveURLError(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("exported"),
	}
	expectedError := errors.New("failed to get archive URL")

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	// Setup expectations for getting archive URL
	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return("", expectedError)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get migration archive URL")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_SaveBackupError(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("exported"),
	}
	archiveURL := "https://api.github.com/archive/kumojin/12345.zip"
	expectedError := errors.New("failed to save backup")

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	// Setup expectations for getting archive URL
	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(archiveURL, nil)

	// Setup expectations for saving backup
	mocks.saveBackupMock.On("Do", archiveURL).Return("", expectedError)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_ContextCancellation(t *testing.T) {
	// Setup mocks
	mocks := newCreateBackupTestMocks(t)

	// Setup test data
	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	migration := &gh.Migration{
		ID:    gh.Ptr(int64(12345)),
		State: gh.Ptr("pending"),
	}

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Setup expectations for listing repositories
	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	// Setup expectations for starting migration
	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	// Setup expectations for getting migration status - never called because we'll cancel the context
	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Run(func(ctx context.Context, org string, id int64) {
			// Cancel the context before returning
			cancel()
			// Sleep to ensure the context cancellation is processed
			time.Sleep(1)
		}).
		Return(nil, context.Canceled)

	// Create use case with mocks
	useCase := mocks.createUseCase()

	// Execute the use case
	result, err := useCase.Do(ctx, organization, mocks.saveBackupFunc)

	// Assertions
	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, result)
}

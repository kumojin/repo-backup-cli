package uc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gh "github.com/google/go-github/v78/github"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSaveBackupFunc is a mock implementation of SaveBackupFunc
type MockSaveBackupFunc struct {
	mock.Mock
}

func (m *MockSaveBackupFunc) Do(reader io.Reader) (string, error) {
	// Read the content to pass it to the mock for assertions
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	args := m.Called(string(content))
	return args.String(0), args.Error(1)
}

// createBackupTestMocks contains all the mocks used in tests
type createBackupTestMocks struct {
	githubClient              *github.MockClient
	listPrivateRepos          *MockListPrivateReposUseCase
	getOrganizationArchiveUrl *MockGetOrganizationArchiveUrlUseCase
	saveBackupFunc            func(reader io.Reader) (string, error)
	saveBackupMock            *MockSaveBackupFunc
}

// newCreateBackupTestMocks creates and returns all the mocks needed for testing
func newCreateBackupTestMocks(t *testing.T) *createBackupTestMocks {
	mockGithubClient := github.NewMockClient(t)
	mockListPrivateRepos := NewMockListPrivateReposUseCase(t)
	mockGetArchiveUrl := NewMockGetOrganizationArchiveUrlUseCase(t)

	mockSaveBackup := new(MockSaveBackupFunc)
	saveBackupFunc := func(reader io.Reader) (string, error) {
		return mockSaveBackup.Do(reader)
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
	// Given
	mocks := newCreateBackupTestMocks(t)
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
	archiveContent := "mock archive content"
	savePath := "/tmp/backup.zip"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(server.URL, nil)

	mocks.saveBackupMock.On("Do", archiveContent).Return(savePath, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, savePath, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_SuccessOnSecondCall(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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
	archiveContent := "mock archive content"
	savePath := "/tmp/backup.zip"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(pendingMigration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(pendingMigration, nil).
		Once()

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil).
		Once()

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(server.URL, nil)

	mocks.saveBackupMock.On("Do", archiveContent).Return(savePath, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, savePath, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_ListRepositoriesError(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

	organization := "kumojin"
	expectedError := errors.New("failed to list repositories")

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return([]gh.Repository{}, expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_StartMigrationError(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

	organization := "kumojin"
	repos := []gh.Repository{
		{Name: gh.Ptr("repo1")},
		{Name: gh.Ptr("repo2")},
	}
	repoNames := []string{"repo1", "repo2"}
	expectedError := errors.New("failed to start migration")

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(nil, expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_FailedMigration(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(failedMigration, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.ErrorIs(t, err, ErrMigrationFailed)
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_GetMigrationStatusError(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(nil, expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.ErrorIs(t, err, expectedError)
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_GetArchiveURLError(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return("", expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get migration archive URL")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_SaveBackupError(t *testing.T) {
	// Given
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
	archiveContent := "mock archive content"
	expectedError := errors.New("failed to save backup")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(server.URL, nil)

	mocks.saveBackupMock.On("Do", archiveContent).Return("", expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.Error(t, err)
	assert.Empty(t, result)
	mocks.saveBackupMock.AssertExpectations(t)
}

func TestCreateBackupUseCase_HTTPStatusNotOK(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(server.URL, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download archive, got status: 500 Internal Server Error")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_HTTPGetError(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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
	// Use an invalid URL scheme to force an HTTP error
	invalidURL := "invalid://not-a-valid-url"

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil)

	mocks.getOrganizationArchiveUrl.EXPECT().Do(mock.Anything, organization, int64(12345)).Return(invalidURL, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, mocks.saveBackupFunc)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download archive")
	assert.Empty(t, result)
}

func TestCreateBackupUseCase_ContextCancellation(t *testing.T) {
	// Given
	mocks := newCreateBackupTestMocks(t)

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

	ctx, cancel := context.WithCancel(context.Background())

	mocks.listPrivateRepos.EXPECT().Do(mock.Anything, organization).Return(repos, nil)

	mocks.githubClient.EXPECT().
		StartMigration(mock.Anything, organization, repoNames).
		Return(migration, nil)

	mocks.githubClient.EXPECT().
		GetMigrationStatus(mock.Anything, organization, int64(12345)).
		Return(migration, nil).
		Maybe()

	useCase := mocks.createUseCase()

	cancel()

	// When
	result, err := useCase.Do(ctx, organization, mocks.saveBackupFunc)

	// Then
	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, result)
}

package uc

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createRemoteBackupTestMocks contains all the mocks used in tests
type createRemoteBackupTestMocks struct {
	blobRepository      *storage.MockBlobRepository
	createBackupUseCase *MockCreateBackupUseCase
}

// newCreateRemoteBackupTestMocks creates and returns all the mocks needed for testing
func newCreateRemoteBackupTestMocks(t *testing.T) *createRemoteBackupTestMocks {
	mockBlobRepository := storage.NewMockBlobRepository(t)
	mockCreateBackupUseCase := NewMockCreateBackupUseCase(t)

	getCurrentTime = func() time.Time { return time.Date(2025, 7, 23, 0, 0, 0, 0, time.UTC) }

	return &createRemoteBackupTestMocks{
		blobRepository:      mockBlobRepository,
		createBackupUseCase: mockCreateBackupUseCase,
	}
}

// createUseCase creates a CreateRemoteBackupUseCase with the mocks
func (m *createRemoteBackupTestMocks) createUseCase() CreateRemoteBackupUseCase {
	return NewCreateRemoteBackupUseCase(
		m.blobRepository,
		m.createBackupUseCase,
	)
}

func TestCreateRemoteBackupUseCase_Success(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"
	expectedBlobURL := "https://storage.azure.com/blob/2025-07-23-org-migration.tar.gz"

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			reader := strings.NewReader(archiveContent)
			saveFunc(reader)
		}).
		Return(expectedBlobURL, nil)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(blobName string) bool {
			return strings.Contains(blobName, "kumojin-migration.tar.gz")
		}), mock.AnythingOfType("*strings.Reader")).
		Run(func(ctx context.Context, blobName string, reader io.Reader) {
			content, err := io.ReadAll(reader)
			assert.NoError(t, err)
			assert.Equal(t, archiveContent, string(content))
		}).
		Return(expectedBlobURL, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedBlobURL, result)
}

func TestCreateRemoteBackupUseCase_CreateBackupError(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	expectedError := errors.New("failed to create backup")

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Return("", expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create backup")
	assert.Empty(t, result)
}

func TestCreateRemoteBackupUseCase_BlobUploadError(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"
	uploadError := errors.New("failed to upload blob")

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			reader := strings.NewReader(archiveContent)
			_, err := saveFunc(reader)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to upload blob")
		}).
		Return("", uploadError)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(blobName string) bool {
			return strings.Contains(blobName, "kumojin-migration.tar.gz")
		}), mock.AnythingOfType("*strings.Reader")).
		Return("", uploadError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload blob")
	assert.Empty(t, result)
}

func TestCreateRemoteBackupUseCase_BlobNameFormat(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"
	expectedBlobURL := "https://storage.azure.com/blob/test-blob.tar.gz"

	var capturedBlobName string

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			reader := strings.NewReader(archiveContent)
			saveFunc(reader)
		}).
		Return(expectedBlobURL, nil)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*strings.Reader")).
		Run(func(ctx context.Context, blobName string, reader io.Reader) {
			capturedBlobName = blobName
		}).
		Return(expectedBlobURL, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedBlobURL, result)
	assert.Contains(t, capturedBlobName, "-kumojin-migration.tar.gz")
	assert.Contains(t, capturedBlobName, "2025-07-23")
}

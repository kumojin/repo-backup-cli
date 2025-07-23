package uc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			// Call the saveBackupFunc with our test server URL
			saveFunc(server.URL)
		}).
		Return(expectedBlobURL, nil)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(blobName string) bool {
			return strings.Contains(blobName, "org-migration.tar.gz")
		}), mock.AnythingOfType("*http.bodyEOFSignal")).
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

func TestCreateRemoteBackupUseCase_HTTPError(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			// Call the saveBackupFunc with our test server URL that returns 404
			_, err := saveFunc(server.URL)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "could not download archive, got status: 404 Not Found")
		}).
		Return("", errors.New("could not download archive, got status: 404 Not Found"))

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not download archive, got status: 404 Not Found")
	assert.Empty(t, result)
}

func TestCreateRemoteBackupUseCase_DownloadError(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	connectionError := errors.New("connection error")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Force connection close by hijacking the connection
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, err := hj.Hijack()
			if err == nil {
				conn.Close() // Abruptly close the connection
			}
		}
	}))
	defer server.Close()

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			_, err := saveFunc(server.URL)
			assert.Error(t, err)
		}).
		Return("", connectionError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization)

	// Then
	assert.Error(t, err)
	assert.ErrorIs(t, err, connectionError)
	assert.Empty(t, result)
}

func TestCreateRemoteBackupUseCase_BlobUploadError(t *testing.T) {
	// Given
	mocks := newCreateRemoteBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"
	uploadError := errors.New("failed to upload blob")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			_, err := saveFunc(server.URL)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to upload blob")
		}).
		Return("", uploadError)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.MatchedBy(func(blobName string) bool {
			return strings.Contains(blobName, "org-migration.tar.gz")
		}), mock.AnythingOfType("*http.bodyEOFSignal")).
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(archiveContent))
	}))
	defer server.Close()

	var capturedBlobName string

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			saveFunc(server.URL)
		}).
		Return(expectedBlobURL, nil)

	mocks.blobRepository.EXPECT().
		Upload(mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*http.bodyEOFSignal")).
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
	assert.Contains(t, capturedBlobName, "-org-migration.tar.gz")
	assert.Contains(t, capturedBlobName, "2025-07-23") // Current date based on context
}

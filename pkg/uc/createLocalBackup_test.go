package uc

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createLocalBackupTestMocks contains all the mocks used in tests
type createLocalBackupTestMocks struct {
	createBackupUseCase *MockCreateBackupUseCase
}

// newCreateLocalBackupTestMocks creates and returns all the mocks needed for testing
func newCreateLocalBackupTestMocks(t *testing.T) *createLocalBackupTestMocks {
	mockCreateBackupUseCase := NewMockCreateBackupUseCase(t)

	return &createLocalBackupTestMocks{
		createBackupUseCase: mockCreateBackupUseCase,
	}
}

// createUseCase creates a CreateLocalBackupUseCase with the mocks
func (m *createLocalBackupTestMocks) createUseCase() CreateLocalBackupUseCase {
	return NewCreateLocalBackupUseCase(m.createBackupUseCase)
}

func TestCreateLocalBackupUseCase_Success(t *testing.T) {
	// Given
	mocks := newCreateLocalBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"

	tempDir := t.TempDir()
	backupPath := filepath.Join(tempDir, "backup.tar.gz")

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			reader := strings.NewReader(archiveContent)
			result, err := saveFunc(reader)
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
		}).
		Return(backupPath, nil)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, backupPath)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, backupPath, result)

	content, err := os.ReadFile(backupPath)
	assert.NoError(t, err)
	assert.Equal(t, archiveContent, string(content))

	absPath, err := filepath.Abs(backupPath)
	assert.NoError(t, err)
	assert.Equal(t, absPath, result)
}

func TestCreateLocalBackupUseCase_CreateBackupError(t *testing.T) {
	// Given
	mocks := newCreateLocalBackupTestMocks(t)
	organization := "kumojin"
	backupPath := "/tmp/backup.tar.gz"
	expectedError := errors.New("failed to create backup")

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Return("", expectedError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, backupPath)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create backup")
	assert.Empty(t, result)
}

func TestCreateLocalBackupUseCase_FileCreationError(t *testing.T) {
	// Given
	mocks := newCreateLocalBackupTestMocks(t)
	organization := "kumojin"
	archiveContent := "mock archive content"

	// Use an invalid path that will cause os.Create to fail
	invalidPath := "/root/nonexistent/backup.tar.gz"

	var capturedError error

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			reader := strings.NewReader(archiveContent)
			result, err := saveFunc(reader)
			capturedError = err
			assert.Error(t, err)
			assert.Empty(t, result)
		}).
		Return("", errors.New("permission denied"))

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, invalidPath)

	// Then
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, capturedError)
}

func TestCreateLocalBackupUseCase_IOCopyError(t *testing.T) {
	// Given
	mocks := newCreateLocalBackupTestMocks(t)
	organization := "kumojin"

	// Create a temporary file path
	tempDir := t.TempDir()
	backupPath := filepath.Join(tempDir, "backup.tar.gz")

	// Create a reader that will cause an error during copy
	errorReader := &errorReader{err: errors.New("read error")}
	readError := errors.New("read error")

	mocks.createBackupUseCase.EXPECT().
		Do(mock.Anything, organization, mock.AnythingOfType("uc.SaveBackupFunc")).
		Run(func(ctx context.Context, org string, saveFunc SaveBackupFunc) {
			result, err := saveFunc(errorReader)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "read error")
			assert.Empty(t, result)
		}).
		Return("", readError)

	useCase := mocks.createUseCase()

	// When
	result, err := useCase.Do(context.Background(), organization, backupPath)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read error")
	assert.Empty(t, result)
}

// errorReader is a helper struct that implements io.Reader and always returns an error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

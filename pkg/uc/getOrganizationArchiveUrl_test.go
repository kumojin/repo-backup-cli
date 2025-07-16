package uc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrganizationArchiveUrlUseCase_SuccessfullyGetURL(t *testing.T) {
	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Setup mock expectations
	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Return("https://api.github.com/archive/kumojin/12345.zip", nil)

	// Create use case with mock client with short durations for testing
	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(DefaultTimeoutDuration, 1)

	// Execute the use case
	url, err := useCase.Do(context.Background(), "kumojin", int64(12345))

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/archive/kumojin/12345.zip", url)
}

func TestGetOrganizationArchiveUrlUseCase_ErrorFirstAttemptThenSuccess(t *testing.T) {
	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Set up a counter to track call attempts
	callCount := 0

	// Setup mock expectations - first call returns error, second call succeeds
	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Run(func(ctx context.Context, org string, orgID int64) {
			callCount++
		}).
		Return("", errors.New("not ready yet")).
		Once()

	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Run(func(ctx context.Context, org string, orgID int64) {
			callCount++
		}).
		Return("https://api.github.com/archive/kumojin/12345.zip", nil).
		Once()

	// Create use case with mock client with short durations for testing
	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(DefaultTimeoutDuration, 1)

	// Execute the use case
	url, err := useCase.Do(context.Background(), "kumojin", int64(12345))

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/archive/kumojin/12345.zip", url)
	assert.Equal(t, 2, callCount, "The API should have been called exactly twice")
}

func TestGetOrganizationArchiveUrlUseCase_ContextTimeout(t *testing.T) {
	// Use a regular context without timeout as the timeout is now controlled by the use case options
	ctx := context.Background()

	// Create mock client
	mockClient := mocks.NewMockClient(t)

	// Setup mock expectations - always return error to force timeout
	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Return("", errors.New("not ready yet")).
		Maybe() // Use Maybe since we don't know exactly how many times it will be called before timing out

	// Create use case with mock client with very short timeout for testing
	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(1, 5*time.Millisecond)

	// Execute the use case
	url, err := useCase.Do(ctx, "kumojin", int64(12345))

	// Assertions
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Empty(t, url)
}

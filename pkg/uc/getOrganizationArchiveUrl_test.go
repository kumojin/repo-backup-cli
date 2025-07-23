package uc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrganizationArchiveUrlUseCase_SuccessfullyGetURL(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Return("https://api.github.com/archive/kumojin/12345.zip", nil)

	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(DefaultTimeoutDuration, 1)

	// When
	url, err := useCase.Do(context.Background(), "kumojin", int64(12345))

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/archive/kumojin/12345.zip", url)
}

func TestGetOrganizationArchiveUrlUseCase_ErrorFirstAttemptThenSuccess(t *testing.T) {
	// Given
	mockClient := github.NewMockClient(t)

	callCount := 0

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

	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(DefaultTimeoutDuration, 1)

	// When
	url, err := useCase.Do(context.Background(), "kumojin", int64(12345))

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "https://api.github.com/archive/kumojin/12345.zip", url)
	assert.Equal(t, 2, callCount, "The API should have been called exactly twice")
}

func TestGetOrganizationArchiveUrlUseCase_ContextTimeout(t *testing.T) {
	// Given
	ctx := context.Background()

	mockClient := github.NewMockClient(t)

	mockClient.EXPECT().
		GetMigrationArchiveURL(mock.Anything, "kumojin", int64(12345)).
		Return("", errors.New("not ready yet")).
		Maybe() // Use Maybe since we don't know exactly how many times it will be called before timing out

	useCase := NewGetOrganizationArchiveUrlUseCase(mockClient).WithDurationOptions(1, 5*time.Millisecond)

	// When
	url, err := useCase.Do(ctx, "kumojin", int64(12345))

	// Then
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Empty(t, url)
}

package uc

import (
	"context"
	"fmt"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github"
)

const (
	// DefaultTimeoutDuration is the default time to wait for the archive URL to be available
	DefaultTimeoutDuration = time.Second * 30
	// DefaultTickerDuration is the default interval to check for the archive URL
	DefaultTickerDuration = time.Second * 5
)

type GetOrganizationArchiveUrlUseCase interface {
	Do(ctx context.Context, organization string, organizationID int64) (string, error)
	WithDurationOptions(timeout, ticker time.Duration) GetOrganizationArchiveUrlUseCase
}
type getOrganizationArchiveUrlUseCase struct {
	gitHubClient    github.Client
	timeoutDuration time.Duration
	tickerDuration  time.Duration
}

func NewGetOrganizationArchiveUrlUseCase(client github.Client) GetOrganizationArchiveUrlUseCase {
	return &getOrganizationArchiveUrlUseCase{
		gitHubClient:    client,
		timeoutDuration: DefaultTimeoutDuration,
		tickerDuration:  DefaultTickerDuration,
	}
}

// WithDurationOptions sets custom durations for timeout and ticker and returns the modified use case
func (uc *getOrganizationArchiveUrlUseCase) WithDurationOptions(timeout, ticker time.Duration) GetOrganizationArchiveUrlUseCase {
	uc.timeoutDuration = timeout
	uc.tickerDuration = ticker

	return uc
}

func (uc *getOrganizationArchiveUrlUseCase) Do(ctx context.Context, organization string, organizationID int64) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, uc.timeoutDuration)
	defer cancel()

	ticker := time.NewTicker(uc.tickerDuration)
	defer ticker.Stop()

	var err error
	var archiveURL string

	fmt.Println("Trying to get migration archive URL...")

	for {
		select {
		case <-ticker.C:
			archiveURL, err = uc.gitHubClient.GetMigrationArchiveURL(ctx, organization, organizationID)
			if err == nil {
				fmt.Println("Migration archive URL retrieved successfully.")
				return archiveURL, nil
			}
		case <-ctxTimeout.Done():
			return "", fmt.Errorf("context timed out while getting migration archive URL: %w", ctxTimeout.Err())
		}
	}
}

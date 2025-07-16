package uc

import (
	"context"
	"fmt"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github"
)

type GetOrganizationArchiveUrlUseCase interface {
	Do(ctx context.Context, organization string, organizationID int64) (string, error)
}
type getOrganizationArchiveUrlUseCase struct {
	gitHubClient github.Client
}

func NewGetOrganizationArchiveUrlUseCase(client github.Client) GetOrganizationArchiveUrlUseCase {
	return &getOrganizationArchiveUrlUseCase{
		gitHubClient: client,
	}
}
func (uc *getOrganizationArchiveUrlUseCase) Do(ctx context.Context, organization string, organizationID int64) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)
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
			return "", fmt.Errorf("context timed out while getting migration archive URL: %w", err)
		}
	}
}

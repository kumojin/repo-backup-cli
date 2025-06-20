package uc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v72/github"
)

type GetOrganizationArchiveUrlUseCase interface {
	Do(ctx context.Context, organization string, organizationID int64) (string, error)
}
type getOrganizationArchiveUrlUseCase struct {
	gitHubClient *github.Client
}

func NewGetOrganizationArchiveUrlUseCase(client *github.Client) GetOrganizationArchiveUrlUseCase {
	return &getOrganizationArchiveUrlUseCase{
		gitHubClient: client,
	}
}
func (uc *getOrganizationArchiveUrlUseCase) Do(ctx context.Context, organization string, organizationID int64) (string, error) {
	attempt := 1

	// TODO: update this mechanism and use context timeout instead with ticker instead
	for {
		archiveURL, err := uc.gitHubClient.Migrations.MigrationArchiveURL(ctx, organization, organizationID)
		if err != nil {
			if attempt < 3 {
				attempt++
				time.Sleep(time.Millisecond * 6000)

				fmt.Println("attempt ", attempt)

				continue
			}

			return "", fmt.Errorf("error getting migration archive URL: %w", err)
		}

		return archiveURL, nil
	}
}

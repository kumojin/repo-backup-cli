package uc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github"
)

const defaultPollingInterval = 5 * time.Second

type SaveBackupFunc func(reader io.Reader) (string, error)

type CreateBackupUseCase interface {
	Do(ctx context.Context, organization string, saveBackupFunc SaveBackupFunc) (string, error)
	WithPollingInterval(interval time.Duration) CreateBackupUseCase
}

type createBackupUseCase struct {
	githubClient                     github.Client
	listPrivateReposUseCase          ListPrivateReposUseCase
	getOrganizationArchiveUrlUseCase GetOrganizationArchiveUrlUseCase
	pollingInterval                  time.Duration
}

func NewCreateBackupUseCase(
	client github.Client,
	listPrivateRepoUseCase ListPrivateReposUseCase,
	getOrganizationArchiveUrlUseCase GetOrganizationArchiveUrlUseCase,
) CreateBackupUseCase {
	return &createBackupUseCase{
		githubClient:                     client,
		listPrivateReposUseCase:          listPrivateRepoUseCase,
		getOrganizationArchiveUrlUseCase: getOrganizationArchiveUrlUseCase,
		pollingInterval:                  defaultPollingInterval,
	}
}

func (uc *createBackupUseCase) WithPollingInterval(interval time.Duration) CreateBackupUseCase {
	uc.pollingInterval = interval
	return uc
}

func (uc *createBackupUseCase) Do(ctx context.Context, organization string, saveBackupFunc SaveBackupFunc) (string, error) {
	repos, err := uc.listPrivateReposUseCase.Do(ctx, organization)
	if err != nil {
		return "", fmt.Errorf("failed to list private repositories: %w", err)
	}

	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = *repo.Name
	}

	migration, err := uc.githubClient.StartMigration(ctx, organization, repoNames)
	if err != nil {
		return "", fmt.Errorf("failed to start migration: %w", err)
	}

	ticker := time.NewTicker(uc.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			migration, err = uc.githubClient.GetMigrationStatus(ctx, organization, migration.GetID())
			if err != nil {
				return "", fmt.Errorf("failed to get migration status: %w", err)
			}

			if migration.GetState() == "failed" {
				return "", errors.New("migration failed")
			}

			if migration.GetState() != "exported" {
				fmt.Println("Migration in progress, waiting for completion...")
				break
			}

			url, err := uc.getOrganizationArchiveUrlUseCase.Do(ctx, organization, migration.GetID())
			if err != nil {
				return "", fmt.Errorf("failed to get migration archive URL: %w", err)
			}

			resp, err := http.Get(url)
			if err != nil {
				return "", fmt.Errorf("failed to download archive: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf("failed to download archive, got status: %s", resp.Status)
			}

			return saveBackupFunc(resp.Body)
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

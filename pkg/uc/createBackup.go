package uc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/github"
)

const pollingInterval = 5 * time.Second

type SaveBackupFunc func(url string) (string, error)

type CreateBackupUseCase interface {
	Do(ctx context.Context, organization string, saveBackupFunc SaveBackupFunc) (string, error)
}

type createBackupUseCase struct {
	githubClient                     github.Client
	listPrivateReposUseCase          ListPrivateReposUseCase
	getOrganizationArchiveUrlUseCase GetOrganizationArchiveUrlUseCase
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
	}
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

	ticker := time.NewTicker(pollingInterval)
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
			}

			url, err := uc.getOrganizationArchiveUrlUseCase.Do(ctx, organization, migration.GetID())
			if err != nil {
				return "", fmt.Errorf("failed to get migration archive URL: %w", err)
			}

			return saveBackupFunc(url)
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

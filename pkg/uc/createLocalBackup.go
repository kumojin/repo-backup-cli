package uc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v72/github"
)

const pollingInterval = 5 * time.Second

type CreateLocalBackupUseCase interface {
	Do(ctx context.Context, organization string, backupPath string) (string, error)
}

type createLocalBackupUseCase struct {
	gitHubClient                     *github.Client
	listPrivateReposUseCase          ListPrivateReposUseCase
	getOrganizationArchiveUrlUseCase GetOrganizationArchiveUrlUseCase
}

func NewCreateLocalBackupUseCase(client *github.Client) CreateLocalBackupUseCase {
	return &createLocalBackupUseCase{
		gitHubClient:                     client,
		listPrivateReposUseCase:          NewListPrivateReposUseCase(client),
		getOrganizationArchiveUrlUseCase: NewGetOrganizationArchiveUrlUseCase(client),
	}
}

func (uc *createLocalBackupUseCase) Do(ctx context.Context, organization string, backupPath string) (string, error) {
	repos, err := uc.listPrivateReposUseCase.Do(ctx, organization)
	if err != nil {
		return "", fmt.Errorf("failed to list private repositories: %w", err)
	}

	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = *repo.Name
	}

	migration, _, err := uc.gitHubClient.Migrations.StartMigration(ctx, organization, repoNames, &github.MigrationOptions{
		ExcludeAttachments: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to start migration: %w", err)
	}

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			migration, _, err = uc.gitHubClient.Migrations.MigrationStatus(ctx, organization, migration.GetID())
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

			return saveMigrationArchive(url, backupPath)
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func saveMigrationArchive(url string, backupPath string) (string, error) {
	out, err := os.Create(backupPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not download archive, got status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	archivePath, err := filepath.Abs(out.Name())
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return archivePath, nil
}

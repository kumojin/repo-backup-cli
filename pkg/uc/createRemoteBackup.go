package uc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/kumojin/repo-backup-cli/pkg/storage"
)

type CreateRemoteBackupUseCase interface {
	Do(ctx context.Context, organization string) (string, error)
}

type createRemoteBackupUseCase struct {
	blobRepository      storage.BlobRepository
	gitHubClient        *github.Client
	createBackupUseCase CreateBackupUseCase
}

func NewCreateRemoteBackupUseCase(
	blobRepository storage.BlobRepository,
	githubClient *github.Client,
	createBackupUseCase CreateBackupUseCase,
) CreateRemoteBackupUseCase {
	return &createRemoteBackupUseCase{
		blobRepository:      blobRepository,
		gitHubClient:        githubClient,
		createBackupUseCase: createBackupUseCase,
	}
}

func (uc *createRemoteBackupUseCase) Do(ctx context.Context, organization string) (string, error) {
	saveMigrationArchive := func(url string) (string, error) {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("could not download archive, got status: %s", resp.Status)
		}

		blobName := fmt.Sprintf("%s-org-migration.tar.gz", time.Now().Format(time.DateOnly))

		blobUrl, err := uc.blobRepository.Upload(ctx, blobName, resp.Body)

		return blobUrl, nil
	}

	return uc.createBackupUseCase.Do(ctx, organization, saveMigrationArchive)
}

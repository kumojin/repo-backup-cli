package uc

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/kumojin/repo-backup-cli/pkg/storage"
)

type CreateRemoteBackupUseCase interface {
	Do(ctx context.Context, organization string) (string, error)
}

type createRemoteBackupUseCase struct {
	blobRepository      storage.BlobRepository
	createBackupUseCase CreateBackupUseCase
}

func NewCreateRemoteBackupUseCase(
	blobRepository storage.BlobRepository,
	createBackupUseCase CreateBackupUseCase,
) CreateRemoteBackupUseCase {
	return &createRemoteBackupUseCase{
		blobRepository:      blobRepository,
		createBackupUseCase: createBackupUseCase,
	}
}

func (uc *createRemoteBackupUseCase) Do(ctx context.Context, organization string) (string, error) {
	saveMigrationArchive := func(reader io.Reader) (string, error) {
		blobName := fmt.Sprintf("%s-%s-migration.tar.gz", time.Now().Format(time.DateOnly), organization)
		return uc.blobRepository.Upload(ctx, blobName, reader)
	}

	return uc.createBackupUseCase.Do(ctx, organization, saveMigrationArchive)
}

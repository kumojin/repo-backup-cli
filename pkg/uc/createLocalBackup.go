package uc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CreateLocalBackupUseCase interface {
	Do(ctx context.Context, organization string, backupPath string) (string, error)
}

type createLocalBackupUseCase struct {
	createBackupUseCase CreateBackupUseCase
}

func NewCreateLocalBackupUseCase(createBackupUseCase CreateBackupUseCase) CreateLocalBackupUseCase {
	return &createLocalBackupUseCase{
		createBackupUseCase: createBackupUseCase,
	}
}

func (uc *createLocalBackupUseCase) Do(ctx context.Context, organization string, backupPath string) (string, error) {
	saveMigrationArchive := func(reader io.Reader) (string, error) {
		out, err := os.Create(backupPath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		_, err = io.Copy(out, reader)
		if err != nil {
			return "", err
		}

		archivePath, err := filepath.Abs(out.Name())
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}

		return archivePath, nil
	}

	return uc.createBackupUseCase.Do(ctx, organization, saveMigrationArchive)
}

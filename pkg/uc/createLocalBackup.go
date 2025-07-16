package uc

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	saveMigrationArchive := func(url string) (string, error) {
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

	return uc.createBackupUseCase.Do(ctx, organization, saveMigrationArchive)
}

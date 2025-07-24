package cmd

import (
	"context"
	"log/slog"
	"os"

	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/kumojin/repo-backup-cli/pkg/storage/azure"
	"github.com/kumojin/repo-backup-cli/pkg/uc"

	"github.com/spf13/cobra"
)

func BackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Commands to backup repositories from an organization",
		RunE:  runReposCommand,
	}

	cmd.AddCommand(LocalBackupCommand())
	cmd.AddCommand(RemoteBackupCommand())

	return cmd
}

func LocalBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Backup repositories to local storage",
		RunE:  runLocalBackupCommand,
	}

	return cmd
}

func RemoteBackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Backup repositories to remote storage",
		RunE:  runRemoteBackupCommand,
	}

	return cmd
}

func runLocalBackupCommand(_ *cobra.Command, _ []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateLocalBackupUseCase(createBackupUseCase)

	ctx := context.Background()

	archivePath, err := usecase.Do(ctx, cfg.Organization, "archive.tar.gz")
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("organization", cfg.Organization),
		slog.String("backupURL", archivePath),
		slog.String("backupType", "local"),
	)

	logger.Info("backup completed successfully")

	return nil
}

func runRemoteBackupCommand(_ *cobra.Command, _ []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	azClient, err := appContext.GetAzureBlobClient(cfg)
	if err != nil {
		return err
	}
	blobRepository := azure.NewBlobRepository(cfg, azClient)

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateRemoteBackupUseCase(blobRepository, createBackupUseCase)

	remoteUrl, err := usecase.Do(context.Background(), cfg.Organization)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("organization", cfg.Organization),
		slog.String("backupURL", remoteUrl),
		slog.String("backupType", "remote"),
	)

	logger.Info("backup completed successfully")

	return nil
}

func getCreateBackupUseCase(cfg *config.Config) uc.CreateBackupUseCase {
	githubClient := github.NewClient(appContext.GetGithubClient(cfg))

	return uc.NewCreateBackupUseCase(
		githubClient,
		uc.NewListPrivateReposUseCase(githubClient),
		uc.NewGetOrganizationArchiveUrlUseCase(githubClient),
	)
}

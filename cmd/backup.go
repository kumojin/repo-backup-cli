package cmd

import (
	"context"
	"log/slog"

	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/config"
	"github.com/kumojin/repo-backup-cli/pkg/github"
	"github.com/kumojin/repo-backup-cli/pkg/logging"
	"github.com/kumojin/repo-backup-cli/pkg/uc"

	"github.com/spf13/cobra"
)

func BackupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Commands to backup repositories from an organization",
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
	ctx := context.Background()
	logger := logging.NewLogger(ctx).With(
		slog.String("backupType", "local"),
	)

	cfg, err := getConfig()
	if err != nil {
		logger.Error("could not get config", slog.Any("error", err))
		return err
	}

	logger = logger.With(slog.String("organization", cfg.Organization))

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateLocalBackupUseCase(createBackupUseCase)

	archivePath, err := usecase.Do(ctx, cfg.Organization, "archive.tar.gz")
	if err != nil {
		logger.Error("could not create local backup", slog.Any("error", err))
		return err
	}

	logger.With(slog.String("backupURL", archivePath)).Info("backup completed successfully")

	return nil
}

func runRemoteBackupCommand(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	logger := logging.NewLogger(ctx).With(
		slog.String("backupType", "remote"),
	)

	cfg, err := getConfig()
	if err != nil {
		logger.Error("could not get config", slog.Any("error", err))
		return err
	}

	logger = logger.With(slog.String("organization", cfg.Organization))

	blobRepository, err := appContext.NewBlobRepository(cfg)
	if err != nil {
		logger.Error("could not get blob repository", slog.Any("error", err))
		return err
	}

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateRemoteBackupUseCase(blobRepository, createBackupUseCase)

	remoteUrl, err := usecase.Do(ctx, cfg.Organization)
	if err != nil {
		logger.Error("could not create remote backup", slog.Any("error", err))
		return err
	}

	logger.With(
		slog.String("backupURL", remoteUrl),
	).Info("backup completed successfully")

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

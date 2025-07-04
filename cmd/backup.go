package cmd

import (
	"context"
	"fmt"

	appContext "github.com/kumojin/repo-backup-cli/context"
	"github.com/kumojin/repo-backup-cli/pkg/config"
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

	githubClient := appContext.GetGithubClient(cfg)

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateLocalBackupUseCase(githubClient, createBackupUseCase)

	ctx := context.Background()

	archivePath, err := usecase.Do(ctx, cfg.Organization, "archive.tar.gz")
	if err != nil {
		return err
	}

	fmt.Printf("Backup completed successfully! File saved at %s\n", archivePath)

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

	githubClient := appContext.GetGithubClient(cfg)

	createBackupUseCase := getCreateBackupUseCase(cfg)

	usecase := uc.NewCreateRemoteBackupUseCase(blobRepository, githubClient, createBackupUseCase)

	remoteUrl, err := usecase.Do(context.Background(), cfg.Organization)
	if err != nil {
		return err
	}

	fmt.Printf("Backup completed successfully! File saved remotely at %s\n", remoteUrl)

	return nil
}

func getCreateBackupUseCase(cfg *config.Config) uc.CreateBackupUseCase {
	githubClient := appContext.GetGithubClient(cfg)

	return uc.NewCreateBackupUseCase(
		githubClient,
		uc.NewListPrivateReposUseCase(githubClient),
		uc.NewGetOrganizationArchiveUrlUseCase(githubClient),
	)
}

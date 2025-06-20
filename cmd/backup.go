package cmd

import (
	"context"
	"fmt"

	appContext "github.com/kumojin/repo-backup-cli/context"
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

func runLocalBackupCommand(_ *cobra.Command, _ []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	client := appContext.GetGitHubClient(cfg)

	usecase := uc.NewCreateLocalBackupUseCase(client)

	ctx := context.Background()

	archivePath, err := usecase.Do(ctx, cfg.Organization, "archive.tar.gz")
	if err != nil {
		return err
	}

	fmt.Printf("Backup completed successfully! File saved at %s\n", archivePath)

	return nil
}

package cmd

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbk",
		Short: "CLI tool to backup private repositories from one a github organization to a remote file storage service",
	}

	cmd.AddCommand(ReposCommand())

	return cmd
}

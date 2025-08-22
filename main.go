package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/cmd"
)

func main() {
	defer sentry.Flush(2 * time.Second)

	rootCmd, err := cmd.RootCommand()
	if err != nil {
		log.Fatalf("could not create root command: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		sentry.CaptureException(fmt.Errorf("could not execute root cmd: %w", err))

		log.Fatalf("could not execute root command: %v", err)
	}
}

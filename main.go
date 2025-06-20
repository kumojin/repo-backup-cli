package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/cmd"
)

const flushTimeout = 2 * time.Second

func main() {
	defer sentry.Flush(flushTimeout)

	rootCmd, err := cmd.RootCommand()
	if err != nil {
		log.Fatalf("could not create root command: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		sentry.Flush(flushTimeout)

		log.Fatal(err)
	}
}

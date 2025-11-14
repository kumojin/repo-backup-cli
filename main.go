package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/fang"
	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/cmd"
	"github.com/kumojin/repo-backup-cli/internal/version"
)

const flushTimeout = 2 * time.Second

func main() {
	defer sentry.Flush(flushTimeout)

	rootCmd, err := cmd.RootCommand()
	if err != nil {
		log.Fatalf("could not create root command: %v", err)
	}

	opts := []fang.Option{
		fang.WithVersion(version.Tag),
		fang.WithCommit(version.Commit),
	}

	fmt.Println(version.Tag, version.Commit)

	if err := fang.Execute(context.Background(), rootCmd, opts...); err != nil {
		sentry.Flush(flushTimeout)

		log.Fatal(err)
	}
}

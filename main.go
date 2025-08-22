package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/cmd"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

func initSentry(cfg config.SentryConfig) (func(), error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Dsn,
	})
	if err != nil {
		return nil, err
	}

	return func() {
		sentry.Flush(2 * time.Second)
	}, nil
}

func main() {
	rootCmd := cmd.RootCommand()

	err := rootCmd.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("could not parse flags: %v", err)
	}

	cfg, err := cmd.GetConfig()
	if err != nil {
		log.Fatalf("could not get config: %v", err)
	}

	flush, err := initSentry(cfg.GetSentryConfig())
	if err != nil {
		log.Fatalf("could not init sentry: %v", err)
	}
	defer flush()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

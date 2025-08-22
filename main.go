package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kumojin/repo-backup-cli/cmd"
	"github.com/kumojin/repo-backup-cli/pkg/config"
)

func initSentry(cfg config.SentryConfig) (func(), error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Dsn,
		EnableLogs:       true,
		SendDefaultPII:   true,
		AttachStacktrace: true,
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
	if err != nil && err != flag.ErrHelp { // Ignore flag parsing errors if it's just help, seems like cobra does not handle those
		log.Fatalf("could not parse flags: %v", err)
	}

	cfg, err := cmd.GetConfig()
	if err != nil {
		log.Fatalf("could not get config: %v", err)
	}

	if cfg.IsSentryEnabled() {
		flush, err := initSentry(cfg.GetSentryConfig())
		if err != nil {
			log.Fatalf("could not init sentry: %v", err)
		}
		defer flush()
	}

	if err := rootCmd.Execute(); err != nil {
		sentry.CaptureException(fmt.Errorf("could not execute root cmd: %w", err))

		log.Fatalf("could not execute root command: %v", err)
	}
}

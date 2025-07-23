package main

import (
	"log"

	"github.com/kumojin/repo-backup-cli/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

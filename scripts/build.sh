#!/bin/sh
set -e

GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")

go build \
  -ldflags="-X 'github.com/kumojin/repo-backup-cli/internal/version.Tag=${GIT_TAG}' \
            -X 'github.com/kumojin/repo-backup-cli/internal/version.Commit=${GIT_COMMIT}'" \
  -o rbk \
  .
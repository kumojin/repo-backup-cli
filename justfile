# Available recipes for repo-backup-cli

# List available recipes
default:
    @just --list

# Copy .env.template to .env if it doesn't exist
setup-env:
    #!/bin/bash
    if [ ! -f .env ]; then
        echo "Creating .env file from template..."
        cp .env.template .env
        echo "Done! Please edit .env with your credentials."
    else
        echo ".env file already exists."
    fi

# Generate mocks using mockery
generate-mocks:
    #!/bin/bash
    @echo "Checking if mockery is installed..."
    if ! command -v mockery &> /dev/null; then
        echo "Error: mockery is not installed. Please install it using homebrew."
        exit 1
    fi
    
    @echo "Generating mocks..."
    mockery

# Clean generated mocks
clean-mocks:
    @echo "Cleaning mocks..."
    find . -path "*/mocks/mock_*.go" -delete

# Build the CLI
build:
    @echo "Building repo-backup-cli..."
    go build -o rbk .

# Run local backup
backup-local: build
    ./rbk backup local

# Run remote backup
backup-remote: build
    ./rbk backup remote

# List repositories
list-repos: build
    ./rbk repos

# Run all tests
test:
    go test -v ./...

# Setup the project (install dependencies and create .env)
setup:
    go mod download
    @just setup-env

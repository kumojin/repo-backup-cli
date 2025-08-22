# Repo Backup CLI

A command-line tool for backing up private GitHub repositories from an organization to local storage or remote object storage.

## Overview

Repo Backup CLI (rbk) provides functionality to:

- List all private, non-archived repositories in a GitHub organization
- Create local backups of repositories as archive files
- Create remote backups to object storage (Azure Blob Storage or S3-compatible storage)

The backup feature of this CLI leverages GitHub's Migration API to create a migration archive of all private non-archived repositories from a specified organization and then downloads or uploads this archive to your desired storage location.

## Storage Backends

The CLI supports two different storage backends for remote backups:

### Azure Blob Storage

Azure Blob Storage is Microsoft's object storage solution for the cloud. To use Azure Blob Storage as your backend:

- Set `STORAGE_BACKEND=azure` in your configuration
- Configure the following environment variables:
  - `AZURE_STORAGE_ACCOUNT_NAME` - Your Azure Storage account name
  - `AZURE_STORAGE_API_KEY` - Your Azure Storage API key
  - `AZURE_STORAGE_ACCOUNT_URL` - Your Azure Storage account URL
  - `AZURE_STORAGE_CONTAINER_NAME` - Your Azure Storage container name

### S3-Compatible Object Storage

The CLI also supports S3-compatible object storage services (such as MinIO, AWS S3, DigitalOcean Spaces, etc.). To use S3-compatible storage:

- Set `STORAGE_BACKEND=object` in your configuration
- Configure the following environment variables:
  - `OBJECT_STORAGE_ENDPOINT` - The endpoint URL of your S3-compatible service
  - `OBJECT_STORAGE_ACCESS_KEY` - Your access key
  - `OBJECT_STORAGE_SECRET_KEY` - Your secret key
  - `OBJECT_STORAGE_BUCKET_NAME` - The bucket name where backups will be stored
  - `OBJECT_STORAGE_USE_SSL` - Whether to use SSL (true/false, defaults to true)

## Development Setup

### Install Dependencies

This project requires several development dependencies which can be installed using Homebrew:

```bash
# Install all dependencies using the included Brewfile
brew bundle
```

### Configuration

You can use the included `justfile` to set up the environment:

```bash
just setup
```

The above command will create a `.env` file in the root directory from the `.env.template` file. You should replace the variables in it with the correct values.

## Building from Source

### Prerequisites

- [Go 1.25+](https://golang.org/doc/install) - The project requires Go version 1.25 or later
- [Just](https://github.com/casey/just) - **(Optional)** Command runner for development tasks (installable via `brew install just`)

### Clone the Repository

```bash
git clone https://github.com/kumojin/repo-backup-cli.git
cd repo-backup-cli
```

### Install Dependencies

Install all development dependencies using the included Brewfile:

```bash
brew bundle
```

This will install Go, Just, and other required tools.

### Build the Binary

You can build the project using either Go directly or the included Justfile:

#### Using Go

```bash
go build -o rbk .
```

#### Using Just

```bash
just build
```

Both commands will create an executable binary named `rbk` in the current directory.

### Install Globally (Optional)

To install the binary globally so you can run `rbk` from anywhere:

```bash
go install .
```

This will install the binary to your `$GOPATH/bin` directory (make sure it's in your `$PATH`).

### Verify Installation

Test that the binary works correctly:

```bash
./rbk --help
```

Or if you installed globally:

```bash
rbk --help
```

## Usage

### Basic Command Structure

```
rbk [command] [flags]
```

### Global Flags

- `-c, --config` - Path to environment configuration file (default: ".env")
- `-o, --organization` - GitHub organization to use

### Available Commands

#### List Repositories

List all private, non-archived repositories in the specified organization. Used primarily for debugging.

```bash
rbk repos
```

#### Backup Repositories

Create a backup of repositories from an organization:

```bash
rbk backup [local|remote]
```

##### Local Backup

Save the backup archive to local storage:

```bash
rbk backup local
```

This will save the archive as `archive.tar.gz` in the current directory.

##### Remote Backup

Upload the backup archive to your configured object storage backend:

```bash
rbk backup remote
```

This will create a blob/object with the name format `YYYY-MM-DD-org-migration.tar.gz` and upload it to your configured storage container/bucket (Azure Blob Storage or S3-compatible storage).

## Example

```bash
# List all private repos in the organization
rbk repos --organization myorg

# Create a local backup of repositories
rbk backup local --organization myorg

# Create a remote backup to object storage
rbk backup remote --organization myorg --config custom.env
```

## Development

### Debug Configuration

VS Code launch configurations are provided in the `.vscode/launch.json` file for debugging for all operations above.

### Using the Justfile

This project includes a `justfile` with common development tasks. Use the following command to list them all:

```bash
just
```

## Automatic backup via Github action

A GitHub Actions workflow is included that automatically runs a remote backup every day at midnight UTC.

The workflow can also be triggered manually from the "Actions" tab in your GitHub repository.

The workflow file is located at `.github/workflows/daily-backup.yml`.

### Setup

Add the following secrets to your GitHub repository:

- `CLI_GITHUB_TOKEN` - A GitHub personal access token with the necessary permissions
- `STORAGE_BACKEND` - The storage backend to use (`azure` or `object`)

**For Azure Blob Storage (`STORAGE_BACKEND=azure`):**

- `AZURE_STORAGE_ACCOUNT_NAME` - Your Azure Storage account name
- `AZURE_STORAGE_API_KEY` - Your Azure Storage API key
- `AZURE_STORAGE_ACCOUNT_URL` - Your Azure Storage account URL
- `AZURE_STORAGE_CONTAINER_NAME` - Your Azure Storage container name

**For S3-Compatible Storage (`STORAGE_BACKEND=object`):**

- `OBJECT_STORAGE_ENDPOINT` - The endpoint URL of your S3-compatible service
- `OBJECT_STORAGE_ACCESS_KEY` - Your access key
- `OBJECT_STORAGE_SECRET_KEY` - Your secret key
- `OBJECT_STORAGE_BUCKET_NAME` - The bucket name where backups will be stored
- `OBJECT_STORAGE_USE_SSL` - Whether to use SSL (true/false)

#### GitHub Token Requirements

The GitHub token must be a **classic personal access token** (not a fine-grained token) with the following permissions:

- `repo` - Full control of private repositories
- `admin:org` - Full control of orgs and teams, read and write org projects

To create a classic token follow these [instructions](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic).

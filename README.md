# Repo Backup CLI

A command-line tool for backing up private GitHub repositories from an organization to local storage or remote Azure Blob Storage.

## Overview

Repo Backup CLI (rbk) provides functionality to:

- List all private, non-archived repositories in a GitHub organization
- Create local backups of repositories as archive files
- Create remote backups to Azure Blob Storage

The backup feature of this CLI leverages GitHub's Migration API to create a migration archive of all private non-archived repositories from a specified organization and then downloads or uploads this archive to your desired storage location.

## Configuration

Create a `.env` file in the root directory from the `.env.template` file and replace the variables with the correct values.

## Usage

### Basic Command Structure

```
rbk [command] [flags]
```

### Global Flags

- `-c, --config` - Path to environment configuration file (default: ".env")
- `-o, --organization` - GitHub organization to use (default: "Kumojin")

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

Upload the backup archive to Azure Blob Storage:

```bash
rbk backup remote
```

This will create a blob with the name format `YYYY-MM-DD-org-migration.tar.gz` and upload it to your configured Azure container.

## Example

```bash
# List all private repos in the organization
rbk repos --organization myorg

# Create a local backup of repositories
rbk backup local --organization myorg

# Create a remote backup to Azure Blob Storage
rbk backup remote --organization myorg --config custom.env
```

## Development

### Debug Configuration

VS Code launch configurations are provided in the `.vscode/launch.json` file for debugging for all operations above.

# gitea-mirror sync

Sync Gitea mirrors

## Usage

```
gitea-mirror sync [<repository> ...] [flags]
```

## Description

The `sync` command triggers an immediate synchronization of the specified mirror repositories in the target Gitea instance. This is useful when you want to update a mirror immediately rather than waiting for the scheduled synchronization.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be synchronized.

## Examples

```bash
# Sync all mirrors defined in the configuration
gitea-mirror sync

# Sync specific mirrors
gitea-mirror sync repo1 repo2

# Sync a mirror with specific owner
gitea-mirror sync owner/repo
```

## Options

```
-h, --help   help for sync
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

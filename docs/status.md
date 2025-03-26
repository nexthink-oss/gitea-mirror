# gitea-mirror status

Print the status of the mirrors

## Usage

```
gitea-mirror status [<repository> ...] [flags]
```

## Description

The `status` command displays the current status of mirror repositories in the target Gitea instance. For each repository, it shows the last synchronization time in UTC.

If no specific repositories are provided as arguments, the status of all repositories defined in the configuration file will be displayed.

## Examples

```bash
# Show status of all mirrors defined in the configuration
gitea-mirror status

# Show status of specific mirrors
gitea-mirror status repo1 repo2

# Show status of a mirror with specific owner
gitea-mirror status owner/repo
```

## Options

```
-h, --help   help for status
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

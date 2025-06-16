# gitea-mirror delete

Delete Gitea mirrors

## Usage

```
gitea-mirror delete [<repository> ...] [flags]
```

## Description

The `delete` command removes mirror repositories from the target Gitea instance. This operation is irreversible and will delete the repository and all its data from the target system.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be deleted.

## Examples

```bash
# Delete all mirrors defined in the configuration
gitea-mirror delete

# Delete specific mirrors
gitea-mirror delete repo1 repo2

# Delete a mirror with non-default owner
gitea-mirror delete owner/repo
```

## Options

```
-h, --help   help for delete
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

# gitea-mirror recreate

Recreate Gitea mirrors

## Usage

```
gitea-mirror recreate [<repository> ...] [flags]
```

## Description

The `recreate` command deletes and then recreates mirrors in the target Gitea instance. This is particularly useful when you need to reset the sync token associated with a mirror, as the Gitea SDK doesn't support updating the token directly.

The process:

1. Deletes the existing mirror repository from the target Gitea instance
2. Creates a new mirror with the same configuration but with updated credentials

**Note**: Missing target organizations are automatically created as needed with the visibility of the first repository mirrored to that organization.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be recreated.

## Examples

```bash
# Recreate all mirrors defined in the configuration
gitea-mirror recreate

# Recreate specific mirrors
gitea-mirror recreate repo1 repo2

# Recreate a mirror with non-default owner
gitea-mirror recreate owner/repo
```

## Options

```
-h, --help   help for recreate
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

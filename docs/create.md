# gitea-mirror create

Create Gitea mirrors

## Usage

```
gitea-mirror create [<repository> ...] [flags]
```

## Description

The `create` command creates mirrors in the target Gitea instance for repositories defined in your configuration. It will:

1. Connect to the source server (GitHub or Gitea)
2. Retrieve repository information
3. Create a mirror in the target Gitea server
4. Set up the mirroring configuration

**Note**: Missing target organizations are automatically created as needed with the visibility of the first repository mirrored to that organization.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be processed.

## Examples

```bash
# Create all mirrors defined in the configuration
gitea-mirror create

# Create specific mirrors
gitea-mirror create repo1 repo2

# Create a mirror with specific owner
gitea-mirror create owner/repo
```

## Options

```
-h, --help   help for create
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

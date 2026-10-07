# gitea-mirror update

Update Gitea mirrors

## Usage

```
gitea-mirror update [<repository> ...] [flags]
```

## Description

The `update` command modifies the configuration of existing mirror repositories in the target Gitea instance. It can update properties such as:

- Mirror synchronization interval
- Public/private status
- Source token, for private sources (Gitea >= 1.27)

The values used for the update are taken from the current configuration.

On Gitea < 1.27 the source token cannot be updated in place; a warning is printed and [`recreate`](recreate.md) should be used instead. Use `--skip-credentials` (or `GM_SKIP_CREDENTIALS=true`) to leave mirror source credentials untouched, in which case no source token is required.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be updated.

## Examples

```bash
# Update all mirrors defined in the configuration
gitea-mirror update

# Update specific mirrors
gitea-mirror update repo1 repo2

# Update a mirror with specific owner
gitea-mirror update owner/repo

# Update settings without rotating the source token
gitea-mirror update --skip-credentials
```

## Options

```
-h, --help               help for update
    --skip-credentials   do not update mirror source credentials
```

## Options inherited from parent commands

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
```

# gitea-mirror

Manage Gitea mirrors

## Usage

```
gitea-mirror [command]
```

## Description

`gitea-mirror` is a command-line tool for managing repository mirrors in Gitea. It allows you to create, update, sync, check status, and delete mirror repositories. It supports mirroring repositories from both GitHub and Gitea sources.

## Available Commands

- [`gitea-mirror config`](config.md): Print the resolved configuration
- [`gitea-mirror create`](create.md): Create Gitea mirrors
- [`gitea-mirror recreate`](recreate.md): Re-create Gitea mirrors following source-token update
- [`gitea-mirror update`](update.md): Update Gitea mirrors
- [`gitea-mirror delete`](delete.md): Delete Gitea mirrors
- [`gitea-mirror status`](status.md): Print the status of the mirrors
- [`gitea-mirror sync`](sync.md): Sync Gitea mirrors


## Global Options

```
-c, --config-file file     configuration files (default [gitea-mirror.yaml])
-l, --labels stringSlice   filter repositories by label
-o, --owner string         default owner
-S, --source.token token   source API token
-T, --target.token token   target API token
-h, --help                 help for gitea-mirror
    --version              version for gitea-mirror
```

## Environment Variables

Configuration can also be provided via environment variables:

```
GM_CONFIG_FILE        - Configuration file name
GM_OWNER              - Default repository owner
GM_LABELS             - Comma-separated list of labels to filter repositories
GM_SOURCE_TOKEN       - Source API token
GM_TARGET_TOKEN       - Target API token
```

## Configuration File

`gitea-mirror` uses a YAML configuration file. By default, it looks for a file named `gitea-mirror.yaml` in the current directory. You can specify a different files with the `--config-file` option.

Example configuration:

```yaml
source:
  type: github

target:
  url: http://localhost:3000

defaults:
  owner: myorg
  interval: 1h
  public-source: false
  public-target: false
  labels: ["default-label", "team-a"]

repositories:
  - name: repo1
    # Inherits labels from defaults: ["default-label", "team-a"]
  - name: repo2
    owner: otherorg
    interval: 30m
    public-target: true
    labels: ["critical", "team-b"]
  - name: repo3
    labels: [] # Explicitly no labels, does not inherit default
```

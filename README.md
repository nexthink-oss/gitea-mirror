# gitea-mirror

A simple utility to manage collections of Gitea repository mirrors, supporting either Gitea or GitHub upstream source.

## Configuration

Configuration via `gitea-mirror.yaml` file. An example configuration is given in the `gitea-mirror.yaml.example` file.

### `source`

The upstream GitHub or Gitea instance to mirror from.

`type`: set to `github` to use `github.com`.

`url`: set to source from an upstream Gitea instance.

### `target`

The downstream Gitea instance to mirror to.

`url` is the URL of the Gitea instance.

### `defaults`

Default values for repositories.

`owner`: the default owner of repositories to mirror (default: unset).

`interval`: default source sync interval (default: `0s` = disabled)

`public`: (default: `false`).

### `repositories`

The list of repositories to mirror, including default overrides.

`name`: (required) the name of the repository to mirror

`owner`: (optional, if default set) the owner of the repository

`public`: (optional) whether the mirrored repository should be public (default: `false`)

`interval`: (optional) source sync interval

## Authentication

Required API tokens can be passed via `SOURCE_TOKEN` and `TARGET_TOKEN` environment variables, passed by argument, hardcoded via `token` key under `source` or `target`, or, if unset, will be interactively requested at runtime.

## Usage

### Create repository mirrors

```shell
gitea-mirror create
```

### Synchronise repository mirrors

```shell
gitea-mirror sync
```

### Check repository mirrors

```shell
gitea-mirror status
```

### Show configuration

```shell
gitea-mirror config
```

Note: this *will* expose API tokens in the output.

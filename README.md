# gitea-mirror

A simple utility to manage collections of Gitea repository mirrors, supporting either Gitea or GitHub upstream source.

## Configuration

Configuration via `gitea-mirror.yaml` file. An example configuration is given in the `gitea-mirror.yaml.example` file.

`source`: the upstream GitHub or Gitea instance to mirror from. `type: github` to use `github.com`; or just `url` for a Gitea instance.

`target`: the downstream Gitea instance to mirror to. `url` is the URL of the Gitea instance.

`repositories`: a list of repositories to mirror, each identified by `owner` and `name`.

Required API tokens can be passed via `SOURCE_TOKEN` and `TARGET_TOKEN` environment variables, hardcoded via `token` key under `source` or `target`, or interactively prompted for at runtime.

## Usage

### Create repository mirrors

```shell
gitea-mirror create
```

### Update repository mirrors

```shell
gitea-mirror update
```

### Check repository mirrors

```shell
gitea-mirror status
```

### Debug configuration

```shell
gitea-mirror debug
```

Note: this *will* expose API tokens in the output.

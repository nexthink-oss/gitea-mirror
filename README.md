# gitea-mirror

[![Go Report Card](https://goreportcard.com/badge/github.com/nexthink-oss/gitea-mirror)](https://goreportcard.com/report/github.com/nexthink-oss/gitea-mirror)
[![GoDoc](https://godoc.org/github.com/nexthink-oss/gitea-mirror?status.svg)](https://godoc.org/github.com/nexthink-oss/gitea-mirror)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

A command-line tool to manage collections of Gitea repository mirrors, supporting either Gitea or GitHub upstream sources.

## Features

- Mirror repositories from GitHub or Gitea to your Gitea instance
- Configure mirroring interval and repository visibility
- Batch-manage multiple repository mirrors
- Trigger manual synchronization
- Check status of mirrors

## Installation

```bash
go install github.com/nexthink-oss/gitea-mirror@latest
```

Or download from [releases](https://github.com/nexthink-oss/gitea-mirror/releases).

## Quick Start

1. Create a configuration file:

```bash
# Create a basic configuration
cat > gitea-mirror.yaml << EOF
source:
  type: github

target:
  url: https://gitea.example.com

defaults:
  owner: myorg
  interval: 8h
  public: false

repositories:
  - name: repo1
  - name: repo2
    interval: 1h
EOF
```

2. Run the create command:

```bash
# Set API tokens via environment variables
export SOURCE_TOKEN=your_github_token
export TARGET_TOKEN=your_gitea_token

# Create all mirrors defined in config
gitea-mirror create
```

## Documentation

For detailed usage and configuration information, see the documentation:

- [`gitea-mirror`](docs/gitea-mirror.md)
- [`gitea-mirror config`](docs/config.md)
- [`gitea-mirror create`](docs/create.md)
- [`gitea-mirror update`](docs/update.md)
- [`gitea-mirror delete`](docs/delete.md)
- [`gitea-mirror status`](docs/status.md)
- [`gitea-mirror sync`](docs/sync.md)

## Configuration

`gitea-mirror` reads its configuration from one or more YAML or TOML configuration files. By default, it looks for `gitea-mirror.yaml` in the current directory. If multiple configuration files are specified, they are merged in order, with later files taking precedence.

An example configuration file is provided in `gitea-mirror.example.yaml`.

### Configuration Sections

#### Source

```yaml
source:
  type: github  # Use "github" for GitHub, "gitea" for Gitea (default)
  url: http://gitea.upstream.example.com  # Required for Gitea source
  alt-url: https://gitea.example.com  # Optional, defaults to source.url
  token: token  # Optional, can be set via environment or command line
```

- `type`: the type of source instance, either `github` or `gitea`.
- `url`: the address of the source Gitea instance from the context within which `gitea-mirror` is run.
- `mirror-url`: the address of the source Gitea instance from the context of the the target Gitea instance, if not the mirror shouldn't use the source instance's configured `server.DOMAIN`.

#### Target

```yaml
target:
  url: https://gitea.example.com  # Required
  token: token  # Optional, can be set via environment or command line
```

#### Defaults

```yaml
defaults:
  owner: myorg  # Default repository owner
  interval: 8h  # Default sync interval (0s to disable)
  public: false  # Default visibility setting
```

#### Repositories

```yaml
repositories:
  - name: repo1  # Uses defaults
  - name: repo2
    owner: otherorg  # Override default owner
    interval: 1h  # Override default interval
    public: true  # Override default visibility
```

## Authentication

API tokens can be provided in several ways (in order of precedence):

1. Command line arguments (`-S/--source.token`, `-T/--target.token`)
2. Environment variables (`SOURCE_TOKEN`, `TARGET_TOKEN`)
3. Configuration file (`token` under `source` or `target`)
4. Interactive prompt (if none of the above are provided)

## Examples

```bash
# Create all mirrors
gitea-mirror create

# Synchronize specific repositories
gitea-mirror sync repo1 repo2

# Check status of mirrors
gitea-mirror status

# Update mirror configuration
gitea-mirror update

# Delete specific mirrors
gitea-mirror delete repo1

# Display current configuration
gitea-mirror config
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

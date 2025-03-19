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

The values used for the update are taken from the current configuration.

If no specific repositories are provided as arguments, all repositories defined in the configuration file will be updated.

## Examples

```bash
# Update all mirrors defined in the configuration
gitea-mirror update

# Update specific mirrors
gitea-mirror update repo1 repo2

# Update a mirror with specific owner
gitea-mirror update owner/repo
```

## Options

```
-h, --help   help for update
```

## Options inherited from parent commands

```
-C, --config-name string        configuration file name (without extension) (default "gitea-mirror")
-P, --config-path stringArray   configuration file path (default [.])
-o, --owner string              default owner
-S, --source.token string       source API token
-T, --target.token string       target API token
```

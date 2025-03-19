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

# Delete a mirror with specific owner
gitea-mirror delete owner/repo
```

## Options

```
-h, --help   help for delete
```

## Options inherited from parent commands

```
-C, --config-name string        configuration file name (without extension) (default "gitea-mirror")
-P, --config-path stringArray   configuration file path (default [.])
-o, --owner string              default owner
-S, --source.token string       source API token
-T, --target.token string       target API token
```

# gitea-mirror config

Print the resolved configuration

## Usage

```
gitea-mirror config [flags]
```

## Description

The `config` command displays the current configuration in YAML format. This includes all settings from configuration files, environment variables, and command line flags, with sensitive information like tokens omitted.

This command is useful for:
- Verifying configuration settings
- Debugging configuration issues
- Generating a template for a new configuration file

## Examples

```bash
# Print the current configuration
gitea-mirror config

# Save the configuration to a file
gitea-mirror config > config.yaml
```

## Options

```
-h, --help   help for config
```

## Options inherited from parent commands

```
-C, --config-name string        configuration file name (without extension) (default "gitea-mirror")
-P, --config-path stringArray   configuration file path (default [.])
-o, --owner string              default owner
-S, --source.token string       source API token
-T, --target.token string       target API token
```

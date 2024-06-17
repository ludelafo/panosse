# Configuration

Configuration can be set using environment variables, flags or a configuration
file.

The order of precedence is:

1. [Flags](#flags)
2. [Environment variables](#environment-variables)
3. [Configuration file](#configuration-file)

Display the current configuration with `panosse config`.

## Flags

Check the available flags for each command with `panosse help <command>` or
`panosse <command> --help`.

### Examples

- `panosse --config-file="path/to/config.yaml"`
- `panosse --dry-run=true`

## Environment variables

The environment variables can be set by prefexing the [flag names](#flags) with
`PANOSSE_` and converting them to uppercase.

### Examples

- `PANOSSE_CONFIG_FILE="path/to/config.yaml"`
- `PANOSSE_DRY_RUN=true`

## Configuration file

The order of precedence for the configuration file is:

1. `--config-file` flag allows to specify any configuration file
2. `config.yaml` in the current directory
3. `$HOME/.panosse/config.yaml` on Linux and `%USER%\.panosse\config.yaml` on
   Windows

If no configuration file is found, the default values are used from the
[flags](#flags) section.

### Examples

```yaml
# custom config.yaml, config.yaml in the current directory or ~/.panosse/config.yaml
dry-run: true
```

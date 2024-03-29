# panosse

[![Latest release](https://img.shields.io/github/v/release/ludelafo/panosse?include_prereleases)](https://github.com/ludelafo/panosse/releases)
[![License](https://img.shields.io/github/license/ludelafo/panosse)](https://github.com/ludelafo/panosse) [![Issues](https://img.shields.io/github/issues/ludelafo/panosse)](https://github.com/ludelafo/panosse/issues) [![Pull requests](https://img.shields.io/github/issues-pr/ludelafo/panosse)](https://github.com/ludelafo/panosse/pulls) [![Discussions](https://img.shields.io/github/discussions/ludelafo/panosse)](https://github.com/ludelafo/panosse/discussions)

panosse is a CLI tool to clean, encode, normalize and verify your FLAC music library.

It is a wrapper around [flac](https://xiph.org/flac/documentation_tools_flac.html) and [metaflac](https://xiph.org/flac/documentation_tools_metaflac.html) and uses [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) under the hood.

**Examples**

```sh
# Clean FLAC files
panosse clean

# Encode FLAC files
panosse encode

# Normalize FLAC files
panosse normalize

# Verify FLAC files
panosse verify
```

## Configuration

Configuration can be set using environment variables, flags or a configuration file.

The order of precedence is:

1. [Environment variables](#environment-variables)
2. [Flags](#flags)
3. [Configuration file](#configuration-file)

### Environment variables

The environment variables can be set by prefexing the [flag names](#flags) with `PANOSSE_` and converting them to uppercase.

**Examples**

- `PANOSSE_CONFIG_FILE="path/to/config.yml"`
- `PANOSSE_INPUT="path/to/your/music/library"`

### Flags

You can check the available flags for each command with `panosse help <command>` or `panosse <command> --help`.

**Examples**

- `panosse --config-file="path/to/config.yml"`
- `panosse --input="path/to/your/music/library"`

### Configuration file

An commented version of the example file is available at [config.yml](./config.yml).

The order of precedence for the configuration file is:

1. `config.yml` in the current directory
2. `$HOME/.panosse/config.yml` on Linux and `%USER%\.panosse\config.yml` on Windows

If no configuration file is found, the default values are used from the [flags](#flags) section.

**Examples**

```yml
# config.yml or $HOME/.panosse/config.yml
input: path/to/your/music/library
```

## Development

### Build panosse

You must have [Go](https://go.dev/) installed and configured to build panosse.

Once Go is installed, build panosse with the following command:

```sh
# Build panosse
go build
```

### Run panosse

Once panosse is built, you can run it with the following command:

```sh
# Run panosse
./panosse
```

## What does panosse mean?

panosse (`/pa.nɔs/`) is a Swiss-French word meaning mop. The idea is that a mop cleans a floor, panosse cleans FLAC files.

## License

panosse is licensed under the [GNU Affero General Public License (GNU AGPL-3.0)](./COPYING).

# panosse

[![Latest release](https://img.shields.io/github/v/release/ludelafo/panosse?include_prereleases)](https://github.com/ludelafo/panosse/releases)
[![License](https://img.shields.io/github/license/ludelafo/panosse)](https://github.com/ludelafo/panosse/blob/main/COPYING.md)
[![Issues](https://img.shields.io/github/issues/ludelafo/panosse)](https://github.com/ludelafo/panosse/issues)
[![Pull requests](https://img.shields.io/github/issues-pr/ludelafo/panosse)](https://github.com/ludelafo/panosse/pulls)
[![Discussions](https://img.shields.io/github/discussions/ludelafo/panosse)](https://github.com/ludelafo/panosse/discussions)

## What is panosse?

panosse is a CLI tool to clean, encode, normalize, and verify your FLAC music
library.

It is merely a wrapper around
[flac](https://xiph.org/flac/documentation_tools_flac.html) and
[metaflac](https://xiph.org/flac/documentation_tools_metaflac.html) and uses
[Cobra](https://github.com/spf13/cobra) and
[Viper](https://github.com/spf13/viper) under the hood.

```text
Usage:
  panosse [command]

Available Commands:
  clean       Clean FLAC files from blocks and tags
  config      Display panosse configuration
  encode      Encode FLAC files
  help        Help about any command
  normalize   Normalize FLAC files with ReplayGain
  verify      Verify FLAC files
```

> [!NOTE]
>
> This is my first Go project. The code may not be idiomatic. I am open to
> suggestions and improvements. Criticism is welcome!

## What panosse is not

panosse is not a music player, tag editor, or a music library manager. panosse
is focused on cleaning, encoding, normalizing, and verifying FLAC files.

Other tools can be used to manage your music library, such as
[beets](https://beets.io/),
[MusicBrainz Picard](https://picard.musicbrainz.org/), or
[foobar2000](https://www.foobar2000.org/).

As already mentioned, panosse is only a wrapper around flac and metaflac. It
does not provide much more functionality. It was developed to automate and set
sane defaults for my music library maintenance.

panosse tries to stay close to the UNIX philosophy of doing one thing and doing
it well. For example, panosse only proccesses one file at a time (except for
normalization), so you can easily parallelize the process using `find` and
`xargs` or similar tools.

## Usage

panosse can be used as a standalone binary or with Docker.

For detailed information, see the dedicated [Usage](./docs/01_USAGE.md)
documentation.

## Commands and flags

Every panosse's commands have a `help` command to describe the command's usage.

You can use `panosse [command] --help` or `panosse help [command]` to display
the help.

For detailed information, see the dedicated
[Commands and flags](./docs/02_COMMANDS_AND_FLAGS.md) documentation.

## Configuration

Configuration can be set using environment variables, flags or a configuration
file.

The order of precedence is:

1. Flags
2. Environment variables
3. Configuration file

Display the current configuration with `panosse config`.

For a commented version of the example file, check the
[`config.yaml`](./config.yaml) file.

For detailed information, see the dedicated
[Configuration](./docs/03_CONFIGURATION.md) documentation.

## Development

To build panosse, [Go](https://go.dev/) must be installed and configured .

Once Go is installed, build panosse with the following command:

```sh
# Build panosse
go build
```

Once panosse is built, run it with the following command:

```sh
# Run panosse
./panosse
```

## What does panosse mean?

panosse (`/pa.nɔs/`) is a Swiss-French word meaning mop. The idea is that a mop
cleans a floor, panosse cleans FLAC files.

## Contributing

If you have interested in contributing to panosse, check the
[Contributing](https://github.com/ludelafo/panosse/blob/main/CONTRIBUTING.md)
guide.

Thank you in advance!

## License

panosse is licensed under the
[GNU Affero General Public License (GNU AGPL-3.0)](https://github.com/ludelafo/panosse/blob/main/LICENCE.md).

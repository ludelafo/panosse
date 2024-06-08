# panosse

[![Latest release](https://img.shields.io/github/v/release/ludelafo/panosse?include_prereleases)](https://github.com/ludelafo/panosse/releases)
[![License](https://img.shields.io/github/license/ludelafo/panosse)](https://github.com/ludelafo/panosse/blob/main/LICENCE.md)
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

For usage and configuration, see the [Usage](#usage) section and the
[Configuration](#configuration) section. Check the
[Concrete example](#concrete-example) for a real-world example.

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

For usage and configuration, see the [Usage](#usage) section and the
[Configuration](#configuration) section. Check the
[Concrete example](#concrete-example) for a real-world example.

## Usage

panosse can be used as a standalone binary or with Docker.

To use panosse as a standalone binary, download the latest release from the
GitHub Releases page: <https://github.com/ludelafo/panosse/releases>.

> [!IMPORTANT]
>
> flac and metaflac must be installed on your computer in order to use panosse.
> You can install them using your package manager or download them from the
> [Xiph website](https://xiph.org/flac/download.html).
>
> panosse was tested with flac version 1.4.2 and metaflac version 1.4.2.

```sh
# Run panosse as a standalone binary
./panosse --help
```

To use panosse with Docker, pull the Docker image from GitHub Container Registry
and run it as a container: <https://ghcr.io/ludelafo/panosse>.

> [!IMPORTANT]
>
> flac and metaflac are already installed in the Docker image. No need to
> install them on your computer.

```sh
# Run panosse as a Docker container
docker run --rm ghcr.io/ludelafo/panosse --help

# Run panosse as a Docker container with a volume
docker run --rm \
  --volume "$(pwd)/custom-config.yaml:/config/custom-config.yaml" \
  --volume "$(pwd):/flac-files" \
  ghcr.io/ludelafo/panosse --config-file /config/custom-config.yaml /files/file.flac
```

Every panosse's commands have a `help` command to describe the command's usage.

You can use `panosse [command] --help` or `panosse help [command]` to display
the help.

```text
$ ./panosse --help
Usage:
  panosse [command]

Available Commands:
  clean       Clean FLAC files from blocks and tags
  config      Display panosse configuration
  encode      Encode FLAC files
  help        Help about any command
  normalize   Normalize FLAC files with ReplayGain
  verify      Verify FLAC files

Flags:
  -C, --config-file string             config file to use (optional - will use "config.yaml" or "~/.panosse/config.yaml" if available)
  -D, --dry-run                        perform a trial run with no changes made
  -F, --flac-command-path string       path to the flac command (checks in $PATH as well) (default "flac")
  -X, --force                          force processing even if no processing is needed
  -h, --help                           help for panosse
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
  -v, --version                        version for panosse

Use "panosse [command] --help" for more information about a command.
```

> [!TIP]
>
> The recommended order to execute panosse is:
>
> 1. [`verify`](#verify)
> 2. [`encode`](#encode)
> 3. [`normalize`](#normalize)
> 4. [`clean`](#clean)
>
> This will ensure all the files are correct, encoded with the latest FLAC
> version, normalized with ReplayGain, and cleaned from unnecessary blocks and
> tags after all other operations.

### Verify

```text
$ ./panosse verify --help
Check the integrity of the FLAC files.

It calls metaflac to verify the FLAC files.

Usage:
  panosse verify <file> [flags]

Examples:
  ## Verify a single FLAC file
  $ panosse verify file.flac

  ## Verify all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n 1 -P $(nproc) panosse verify

  ## Verify all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 panosse verify

Flags:
  -a, --verify-arguments strings   arguments passed to flac to verify the files (default [--test,--silent])
```

### Encode

```text
$ ./panosse encode --help
Encode FLAC files.

It calls flac to encode the FLAC files.

Usage:
  panosse encode <file> [flags]

Examples:
  ## Encode a single FLAC file
  $ panosse encode file.flac

  ## Encode all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n 1 -P $(nproc) panosse encode

  ## Encode all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 panosse encode

Flags:
  -a, --encode-arguments strings                   arguments passed to flac to encode the file (default [--compression-level-8,--delete-input-file,--no-padding,--force,--verify,--warnings-as-errors,--silent])
      --encode-if-encode-argument-tags-mismatch    encode if encode argument tags mismatch (missing or different) (default true)
      --encode-if-flac-versions-mismatch           encode if flac versions mismatch between host's flac version and file's flac version (default true)
      --save-encode-arguments-in-tag               save encode arguments in tag (default true)
      --save-encode-arguments-in-tag-name string   encode arguments tag name (default "FLAC_ARGUMENTS")
```

### Normalize

```text
$ ./panosse normalize --help
Normalize FLAC files with ReplayGain.

It calls metaflac to calculate and add the ReplayGain tags to the FLAC files.

Usage:
  panosse normalize <file 1> [<file 2>]... [flags]

Examples:
  ## Normalize some FLAC files
  $ panosse normalize file1.flac file2.flac

  ## Normalize all FLAC files in each sub-directory for a depth of 1 in parallel
  # This allows to consider the nested directories as one album for the normalization
  $ find . -mindepth 1 -maxdepth 1 -type d -print0 | xargs -0 -n 1 -P $(nproc) bash -c '
    dir="$1"
    flac_files=()

    # Find all FLAC files in the current directory and store them in an array
    while IFS= read -r -d "" file; do
      flac_files+=("$file")
    done < <(find "$dir" -type f -name "*.flac" -print0)

    # Check if there are any FLAC files found
    if [ ${#flac_files[@]} -ne 0 ]; then
      # Pass the .flac files to the panosse normalize command
      panosse normalize "${flac_files[@]}"
    fi
  ' {}

  ## Normalize all FLAC files in each sub-directory for a depth of 1 in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -mindepth 1 -maxdepth 1 -type d -print0 | sort -z | xargs -0 -n 1 bash -c '
    dir="$1"
    flac_files=()

    # Find all FLAC files in the current directory and store them in an array
    while IFS= read -r -d "" file; do
      flac_files+=("$file")
    done < <(find "$dir" -type f -name "*.flac" -print0)

    # Check if there are any FLAC files found
    if [ ${#flac_files[@]} -ne 0 ]; then
      # Pass the .flac files to the panosse normalize command
      panosse normalize "${flac_files[@]}"
    fi
  ' {}

Flags:
  -a, --normalize-arguments strings                     arguments passed to flac to normalize the files (default [--add-replay-gain])
      --normalize-if-any-replaygain-tags-are-missing    normalize if any ReplayGain tags are missing (default true)
      --normalize-if-normalize-argument-tags-mismatch   normalize if normalize arguments tags mismatch (missing or different) (default true)
  -t, --replaygain-tags strings                         ReplayGain tags (default [REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK])
      --save-normalize-arguments-in-tag                 save normalize arguments in tag (default true)
      --save-normalize-arguments-in-tag-name string     normalize arguments tag name (default "METAFLAC_ARGUMENTS")
```

### Clean

```text
$ ./panosse clean --help
Clean FLAC files from blocks and tags.

It calls metaflac to clean the FLAC files.

Usage:
  panosse clean <file> [flags]

Examples:
  ## Clean a single FLAC file
  $ panosse clean file.flac

  ## Clean all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n 1 -P $(nproc) panosse clean

  ## Clean all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 panosse clean

Flags:
  -a, --clean-arguments strings   arguments passed to metaflac to clean the file (default [--remove,--dont-use-padding,--block-type=APPLICATION,--block-type=CUESHEET,--block-type=PADDING,--block-type=PICTURE,--block-type=SEEKTABLE])
  -t, --tags-to-keep strings      tags to keep in the file (default [ALBUM,ALBUMARTIST,ARTIST,COMMENT,DISCNUMBER,FLAC_ARGUMENTS,GENRE,METAFLAC_ARGUMENTS,REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,TITLE,TRACKNUMBER,TOTALDISCS,TOTALTRACKS,YEAR])
```

## Configuration

Configuration can be set using environment variables, flags or a configuration
file.

The order of precedence is:

1. [Flags](#flags)
2. [Environment variables](#environment-variables)
3. [Configuration file](#configuration-file)

Display the current configuration with `panosse config`.

### Flags

Check the available flags for each command with `panosse help <command>` or
`panosse <command> --help`.

#### Examples

- `panosse --config-file="path/to/config.yaml"`
- `panosse --dry-run=true`

### Environment variables

The environment variables can be set by prefexing the [flag names](#flags) with
`PANOSSE_` and converting them to uppercase.

#### Examples

- `PANOSSE_CONFIG_FILE="path/to/config.yaml"`
- `PANOSSE_DRY_RUN=true`

### Configuration file

A commented version of the example file is available at
[config.yaml](./config.yaml).

The order of precedence for the configuration file is:

1. `--config-file` flag allows to specify a configuration file
2. `config.yaml` in the current directory
3. `$HOME/.panosse/config.yaml` on Linux and `%USER%\.panosse\config.yaml` on
   Windows

If no configuration file is found, the default values are used from the
[flags](#flags) section.

#### Examples

```yaml
# custom config.yaml, config.yaml in the current directory or ~/.panosse/config.yaml
dry-run: true
```

## Development

### Build panosse

To build panosse, [Go](https://go.dev/) must be installed and configured .

Once Go is installed, build panosse with the following command:

```sh
# Build panosse
go build
```

### Run panosse

Once panosse is built, run it with the following command:

```sh
# Run panosse
./panosse
```

### Build panosse Docker image

> [!TIP]
>
> Official Docker images for many platforms are already available on GitHub
> Container Registry: <https://ghcr.io/ludelafo/panosse>.

To build panosse Docker image for your current platform:

```sh
# Build panosse Docker image
docker build --tag ghcr.io/ludelafo/panosse .
```

### Run panosse Docker image

To run panosse Docker image:

```sh
# Run panosse Docker image
docker run --rm ghcr.io/ludelafo/panosse
```

You will probably need to map a volume to access the files on your host machine:

```sh
# Run panosse Docker image with a volume
docker run --rm --volume "$(pwd):/files" ghcr.io/ludelafo/panosse --config-file=/files/config.yaml verify /files/file.flac
```

## What does panosse mean?

panosse (`/pa.nɔs/`) is a Swiss-French word meaning mop. The idea is that a mop
cleans a floor, panosse cleans FLAC files.

## Concrete example

This section provides a real-world example of how to use panosse to clean,
encode, normalize, and verify a FLAC music library.

This section will use the panosse's Docker image to process the files.

### Explore and analyze the files

For this concrete example, all my FLAC files will be in a directory named
`files`.

Let's explore the current structure of the files:

```sh
# Display the files' structure
tree ./files
```

```text
./files/
├── processed
│   └── <empty for the moment>
└── raw
    ├── 2Pac - All Eyez on Me (1996) [CD FLAC]
    │   ├── CD 1
    │   │   ├── 01 - Ambitionz Az a Ridah.flac
    │   │   ├── 02 - All Bout U.flac
    │   │   ├── 03 - Skandalouz.flac
    │   │   ├── 04 - Got My Mind Made Up.flac
    │   │   ├── 05 - How Do U Want It.flac
    │   │   ├── 06 - 2 of Amerikaz Most Wanted.flac
    │   │   ├── 07 - No More Pain.flac
    │   │   ├── 08 - Heartz of Men.flac
    │   │   ├── 09 - Life Goes On.flac
    │   │   ├── 10 - Only God Can Judge Me.flac
    │   │   ├── 11 - Tradin War Stories.flac
    │   │   ├── 12 - California Love (Remix).flac
    │   │   ├── 13 - I Ain't Mad at Cha.flac
    │   │   ├── 14 - What'z Ya Phone #.flac
    │   │   └── folder.jpg
    │   └── CD 2
    │       ├── 01 - Can't C Me.flac
    │       ├── 02 - Shorty Wanna Be a Thug.flac
    │       ├── 03 - Holla at Me.flac
    │       ├── 04 - Wonda Why They Call U Bitch.flac
    │       ├── 05 - When We Ride.flac
    │       ├── 06 - Thug Passion.flac
    │       ├── 07 - Picture Me Rollin'.flac
    │       ├── 08 - Check Out Time.flac
    │       ├── 09 - Ratha Be Ya Nigga.flac
    │       ├── 10 - All Eyez on Me.flac
    │       ├── 11 - Run Tha Streetz.flac
    │       ├── 12 - Ain't Hard 2 Find.flac
    │       ├── 13 - Heaven Ain't Hard 2 Find.flac
    │       └── folder.jpg
    ├── Eminem - The Eminem Show - 2002 (2023 2xCD - FLAC); Expanded Edition
    │   ├── 01-01 - Curtains Up (skit).flac
    │   ├── 01-02 - White America.flac
    │   ├── 01-03 - Business.flac
    │   ├── 01-04 - Cleanin Out My Closet.flac
    │   ├── 01-05 - Square Dance.flac
    │   ├── 01-06 - The Kiss (skit).flac
    │   ├── 01-07 - Soldier.flac
    │   ├── 01-08 - Say Goodbye Hollywood.flac
    │   ├── 01-09 - Drips (feat. Obie Trice).flac
    │   ├── 01-10 - Without Me.flac
    │   ├── 01-11 - Paul Rosenberg (skit).flac
    │   ├── 01-12 - Sing for the Moment.flac
    │   ├── 01-13 - Superman (feat. Dina Rae).flac
    │   ├── 01-14 - Hailie’s Song.flac
    │   ├── 01-15 - Steve Berman (skit).flac
    │   ├── 01-16 - When the Music Stops (feat. D12).flac
    │   ├── 01-17 - Say What You Say (feat. Dr. Dre).flac
    │   ├── 01-18 - ’Till I Collapse (feat. Nate Dogg).flac
    │   ├── 01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac
    │   ├── 01-20 - Curtains Close (skit).flac
    │   ├── 02-01 - Stimulate.flac
    │   ├── 02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac
    │   ├── 02-03 - Bump Heads (DJ Green Lantern Version).flac
    │   ├── 02-04 - Jimmy, Brian and Mike.flac
    │   ├── 02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac
    │   ├── 02-06 - Brain Damage (Live from Tramps, New York , 1999).flac
    │   ├── 02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac
    │   ├── 02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac
    │   ├── 02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac
    │   ├── 02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac
    │   ├── 02-11 - Business (Instrumental).flac
    │   ├── 02-12 - Cleanin' Out My Closet (Instrumental).flac
    │   ├── 02-13 - Square Dance (Instrumental).flac
    │   ├── 02-14 - Without Me (Instrumental).flac
    │   ├── 02-15 - Sing for the Moment (Instrumental).flac
    │   ├── 02-16 - Superman (Instrumental).flac
    │   ├── 02-17 - Say What You Say (Instrumental).flac
    │   ├── 02-18 - Till I Collapse (Instrumental).flac
    │   └── folder.jpg
    └── The Notorious B.I.G. (1994) Ready To Die (2006 Remastered) [FLAC]
        ├── 01 Intro.flac
        ├── 02 Things Done Changed.flac
        ├── 03 Gimme The Loot.flac
        ├── 04 Machine Gun Funk.flac
        ├── 05 Warning.flac
        ├── 06 Ready To Die.flac
        ├── 07 One More Chance.flac
        ├── 08 Fuck Me (Interlude).flac
        ├── 09 The What.flac
        ├── 10 Juicy.flac
        ├── 11 Everyday Struggles.flac
        ├── 12 Me & My Bitch.flac
        ├── 13 Big Poppa.flac
        ├── 14 Respect.flac
        ├── 15 Friend Of Mine.flac
        ├── 16 Unbelievable.flac
        ├── 17 Suicidal Thoughts.flac
        ├── 18 Who Shot Ya.flac
        ├── 19 Just Playing (Dreams).flac
        └── folder.jpg
```

The first thing to notice is the difference in structure:

1. _2Pac - All Eyez on Me (1996)_ is split into two CDs, making it a nested
   structure
2. _Eminem - The Eminem Show (2002)_ has two CDs in a flat structure
3. _The Notorious B.I.G. - Ready To Die (1994)_ has a flat structure

The second thing to notice is the difference in naming:

1. _2Pac - All Eyez on Me (1996)_ misses all the featuring in the filenames
   (` (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down)` in the
   `02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac` track)
2. _Eminem - The Eminem Show (2002)_ has tracks named with a `` ` `` instead of
   `'`.
3. _The Notorious B.I.G. - Ready To Die (1994)_ tracks are named
   `[Track number] [Track title].flac` instead of
   `[Track number] - [Track title].flac` as the others albums
4. Each album directory has a different naming convention

Let's have a look at the files:

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/raw/2Pac - All Eyez on Me (1996) [CD FLAC]/CD 1/01 - Ambitionz Az a Ridah.flac"
```

<details>
<summary>Expand the output</summary>

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 14 bytes
  maximum framesize: 12367 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 12288024
  MD5 signature: e63d3931b934b23765bea0e754c27420
METADATA block #1
  type: 3 (SEEKTABLE)
  is last: false
  length: 504
  seek points: 28
    point 0: sample_number=0, stream_offset=0, frame_samples=4096
    point 1: sample_number=438272, stream_offset=859010, frame_samples=4096
    point 2: sample_number=880640, stream_offset=1945774, frame_samples=4096
    point 3: sample_number=1318912, stream_offset=3034368, frame_samples=4096
    point 4: sample_number=1761280, stream_offset=4142698, frame_samples=4096
    point 5: sample_number=2203648, stream_offset=5216326, frame_samples=4096
    point 6: sample_number=2641920, stream_offset=6289235, frame_samples=4096
    point 7: sample_number=3084288, stream_offset=7380840, frame_samples=4096
    point 8: sample_number=3526656, stream_offset=8468452, frame_samples=4096
    point 9: sample_number=3964928, stream_offset=9544547, frame_samples=4096
    point 10: sample_number=4407296, stream_offset=10633110, frame_samples=4096
    point 11: sample_number=4849664, stream_offset=11717855, frame_samples=4096
    point 12: sample_number=5287936, stream_offset=12802331, frame_samples=4096
    point 13: sample_number=5730304, stream_offset=13918381, frame_samples=4096
    point 14: sample_number=6172672, stream_offset=15003813, frame_samples=4096
    point 15: sample_number=6610944, stream_offset=16126360, frame_samples=4096
    point 16: sample_number=7053312, stream_offset=17221544, frame_samples=4096
    point 17: sample_number=7495680, stream_offset=18301284, frame_samples=4096
    point 18: sample_number=7933952, stream_offset=19407848, frame_samples=4096
    point 19: sample_number=8376320, stream_offset=20499114, frame_samples=4096
    point 20: sample_number=8818688, stream_offset=21591306, frame_samples=4096
    point 21: sample_number=9256960, stream_offset=22663085, frame_samples=4096
    point 22: sample_number=9699328, stream_offset=23733993, frame_samples=4096
    point 23: sample_number=10141696, stream_offset=24843150, frame_samples=4096
    point 24: sample_number=10579968, stream_offset=25917352, frame_samples=4096
    point 25: sample_number=11022336, stream_offset=27024462, frame_samples=4096
    point 26: sample_number=11464704, stream_offset=28145844, frame_samples=4096
    point 27: sample_number=11902976, stream_offset=29204061, frame_samples=4096
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 225
  vendor string: reference libFLAC 1.3.2 20170101
  comments: 10
    comment[0]: ALBUM=All Eyez on Me
    comment[1]: TITLE=Ambitionz Az a Ridah
    comment[2]: ARTIST=2Pac
    comment[3]: Date=1996
    comment[4]: GENRE=Hip-Hop
    comment[5]: ALBUMARTIST=2Pac
    comment[6]: TRACKNUMBER=01
    comment[7]: TRACKTOTAL=14
    comment[8]: DISCNUMBER=1
    comment[9]: DISCTOTAL=2
METADATA block #3
  type: 1 (PADDING)
  is last: true
  length: 8200
```

</details>

The FLAC file has four blocks:

1. STREAMINFO
2. SEEKTABLE
3. VORBIS_COMMENT
4. PADDING

The STREAMINFO block is mandatory and contains the audio stream information.

The SEEKTABLE block is optional and contains seek points for the audio stream.
It is used to quickly seek to a specific sample in the audio stream. Nowadays,
most players can read the audio stream without the seek table, so it is not
necessary to keep it.

The VORBIS_COMMENT block is optional and contains metadata about the audio
stream. It is used to store tags like the album, title, artist, etc.

The PADDING block is optional and is used to pad the FLAC file to a specific
size.

Let's have a closer look at the VORBIS_COMMENT block:

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/2Pac - All Eyez on Me (1996) [CD FLAC]/CD 1/01 - Ambitionz Az a Ridah.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 225
  vendor string: reference libFLAC 1.3.2 20170101
  comments: 10
    comment[0]: ALBUM=All Eyez on Me
    comment[1]: TITLE=Ambitionz Az a Ridah
    comment[2]: ARTIST=2Pac
    comment[3]: Date=1996
    comment[4]: GENRE=Hip-Hop
    comment[5]: ALBUMARTIST=2Pac
    comment[6]: TRACKNUMBER=01
    comment[7]: TRACKTOTAL=14
    comment[8]: DISCNUMBER=1
    comment[9]: DISCTOTAL=2
```

Elements to notice:

1. The FLAC file was encoded with FLAC version 1.3.2 (from the vendor string
   `libFLAC 1.3.2`)
2. The FLAC file has most of the tags I want to keep (see the [Clean](#clean)
   section) - I can assume that the tags are correct and the same for all the
   files
3. However, the FLAC file misses the ReplayGain tags - more on this later

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/raw/Eminem - The Eminem Show - 2002 (2023 2xCD - FLAC); Expanded Edition/01-03 - Business.flac"
```

<details>
<summary>Expand the output</summary>

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 910 bytes
  maximum framesize: 13244 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 11100852
  MD5 signature: 956c9ded0e0dad45dead1f55e650266f
METADATA block #1
  type: 3 (SEEKTABLE)
  is last: false
  length: 468
  seek points: 26
    point 0: sample_number=0, stream_offset=0, frame_samples=4096
    point 1: sample_number=438272, stream_offset=695700, frame_samples=4096
    point 2: sample_number=880640, stream_offset=1716555, frame_samples=4096
    point 3: sample_number=1318912, stream_offset=2780751, frame_samples=4096
    point 4: sample_number=1761280, stream_offset=3864149, frame_samples=4096
    point 5: sample_number=2203648, stream_offset=4886139, frame_samples=4096
    point 6: sample_number=2641920, stream_offset=5897313, frame_samples=4096
    point 7: sample_number=3084288, stream_offset=7072681, frame_samples=4096
    point 8: sample_number=3526656, stream_offset=8245812, frame_samples=4096
    point 9: sample_number=3964928, stream_offset=9402685, frame_samples=4096
    point 10: sample_number=4407296, stream_offset=10417171, frame_samples=4096
    point 11: sample_number=4849664, stream_offset=11434177, frame_samples=4096
    point 12: sample_number=5287936, stream_offset=12475506, frame_samples=4096
    point 13: sample_number=5730304, stream_offset=13542384, frame_samples=4096
    point 14: sample_number=6172672, stream_offset=14717808, frame_samples=4096
    point 15: sample_number=6610944, stream_offset=15877479, frame_samples=4096
    point 16: sample_number=7053312, stream_offset=16888577, frame_samples=4096
    point 17: sample_number=7495680, stream_offset=17919215, frame_samples=4096
    point 18: sample_number=7933952, stream_offset=18947842, frame_samples=4096
    point 19: sample_number=8376320, stream_offset=19968207, frame_samples=4096
    point 20: sample_number=8818688, stream_offset=21042817, frame_samples=4096
    point 21: sample_number=9256960, stream_offset=22205366, frame_samples=4096
    point 22: sample_number=9699328, stream_offset=23353935, frame_samples=4096
    point 23: sample_number=10141696, stream_offset=24371713, frame_samples=4096
    point 24: sample_number=10579968, stream_offset=25382529, frame_samples=4096
    point 25: sample_number=11022336, stream_offset=26372576, frame_samples=4096
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 448
  vendor string: reference libFLAC 1.4.1 20220922
  comments: 15
    comment[0]: replaygain_album_gain=-9.54 dB
    comment[1]: replaygain_album_peak=1
    comment[2]: TITLE=Business
    comment[3]: ARTIST=Eminem
    comment[4]: ALBUM=The Eminem Show
    comment[5]: GENRE=Hip-Hop
    comment[6]: COMMENT=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[7]: Date=2002
    comment[8]: replaygain_track_gain=-9.49 dB
    comment[9]: replaygain_track_peak=0.997559
    comment[10]: ALBUMARTIST=Eminem
    comment[11]: TRACKNUMBER=03
    comment[12]: TRACKTOTAL=20
    comment[13]: DISCNUMBER=1
    comment[14]: DISCTOTAL=2
METADATA block #3
  type: 6 (PICTURE)
  is last: false
  length: 503256
  type: 3 (Cover (front))
  MIME type: image/jpeg
  description:
  width: 1000
  height: 1000
  depth: 0
  colors: 0 (unindexed)
  data length: 503214
  data:
    00000000: FF D8 FF E0 00 10 4A 46 49 46 00 01 01 01 00 60 ......JFIF.....`
    [...]
    0007ADA0: C2 FF 00 C2 C2 EC BA 62 7F F9 AB C8 FF D9 00 00 .......b......
METADATA block #4
  type: 1 (PADDING)
  is last: true
  length: 8131
```

</details>

The FLAC file has five blocks:

1. STREAMINFO
2. SEEKTABLE
3. VORBIS_COMMENT
4. PICTURE
5. PADDING

The PICTURE block is optional and contains the cover art for the album (cut for
brevity). As the directory already contains a `folder.jpg` file, I can remove
the PICTURE block and keep the `folder.jpg` file.

Let's have a closer look at the VORBIS_COMMENT block:

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/Eminem - The Eminem Show - 2002 (2023 2xCD - FLAC); Expanded Edition/01-03 - Business.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 448
  vendor string: reference libFLAC 1.4.1 20220922
  comments: 15
    comment[0]: replaygain_album_gain=-9.54 dB
    comment[1]: replaygain_album_peak=1
    comment[2]: TITLE=Business
    comment[3]: ARTIST=Eminem
    comment[4]: ALBUM=The Eminem Show
    comment[5]: GENRE=Hip-Hop
    comment[6]: COMMENT=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[7]: Date=2002
    comment[8]: replaygain_track_gain=-9.49 dB
    comment[9]: replaygain_track_peak=0.997559
    comment[10]: ALBUMARTIST=Eminem
    comment[11]: TRACKNUMBER=03
    comment[12]: TRACKTOTAL=20
    comment[13]: DISCNUMBER=1
    comment[14]: DISCTOTAL=2
```

Elements to notice:

1. The FLAC file was encoded with FLAC version 1.4.1
2. The FLAC file has tags that I don't want to keep (see the [Clean](#clean)
   section), such as DATE
3. The FLAC file has the ReplayGain tags

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/raw/The Notorious B.I.G. (1994) Ready To Die (2006 Remastered) [FLAC]/03 - Gimme The Loot.flac"
```

<details>
<summary>Expand the output</summary>

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 1722 bytes
  maximum framesize: 13377 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 13412868
  MD5 signature: 3335fc652d261db991a51aef346cc2a9
METADATA block #1
  type: 3 (SEEKTABLE)
  is last: false
  length: 558
  seek points: 31
    point 0: sample_number=0, stream_offset=0, frame_samples=4096
    point 1: sample_number=438272, stream_offset=1085954, frame_samples=4096
    point 2: sample_number=880640, stream_offset=2188125, frame_samples=4096
    point 3: sample_number=1318912, stream_offset=3359773, frame_samples=4096
    point 4: sample_number=1761280, stream_offset=4521772, frame_samples=4096
    point 5: sample_number=2203648, stream_offset=5701990, frame_samples=4096
    point 6: sample_number=2641920, stream_offset=6861160, frame_samples=4096
    point 7: sample_number=3084288, stream_offset=8062077, frame_samples=4096
    point 8: sample_number=3526656, stream_offset=9258131, frame_samples=4096
    point 9: sample_number=3964928, stream_offset=10409462, frame_samples=4096
    point 10: sample_number=4407296, stream_offset=11547545, frame_samples=4096
    point 11: sample_number=4849664, stream_offset=12608362, frame_samples=4096
    point 12: sample_number=5287936, stream_offset=13822969, frame_samples=4096
    point 13: sample_number=5730304, stream_offset=15054283, frame_samples=4096
    point 14: sample_number=6172672, stream_offset=16258583, frame_samples=4096
    point 15: sample_number=6610944, stream_offset=17468552, frame_samples=4096
    point 16: sample_number=7053312, stream_offset=18604864, frame_samples=4096
    point 17: sample_number=7495680, stream_offset=19774694, frame_samples=4096
    point 18: sample_number=7933952, stream_offset=20941658, frame_samples=4096
    point 19: sample_number=8376320, stream_offset=22081831, frame_samples=4096
    point 20: sample_number=8818688, stream_offset=23230866, frame_samples=4096
    point 21: sample_number=9256960, stream_offset=24359738, frame_samples=4096
    point 22: sample_number=9699328, stream_offset=25582571, frame_samples=4096
    point 23: sample_number=10141696, stream_offset=26800669, frame_samples=4096
    point 24: sample_number=10579968, stream_offset=27922183, frame_samples=4096
    point 25: sample_number=11022336, stream_offset=29042231, frame_samples=4096
    point 26: sample_number=11464704, stream_offset=30210935, frame_samples=4096
    point 27: sample_number=11902976, stream_offset=31341158, frame_samples=4096
    point 28: sample_number=12345344, stream_offset=32433728, frame_samples=4096
    point 29: sample_number=12787712, stream_offset=33466084, frame_samples=4096
    point 30: sample_number=13225984, stream_offset=34426127, frame_samples=4096
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 193
  vendor string: reference libFLAC 1.2.1 20070917
  comments: 7
    comment[0]: artist=The Notorious B.I.G.
    comment[1]: title=Gimme The Loot
    comment[2]: album=Ready To Die (The Remaster)
    comment[3]: genre=Hip-Hop
    comment[4]: Comment=.
    comment[5]: DATE=2004
    comment[6]: TRACKNUMBER=03
METADATA block #3
  type: 1 (PADDING)
  is last: true
  length: 8179
```

</details>

The FLAC file has four blocks:

1. STREAMINFO
2. SEEKTABLE
3. VORBIS_COMMENT
4. PADDING

Let's have a closer look at the VORBIS_COMMENT block:

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/The Notorious B.I.G. (1994) Ready To Die (2006 Remastered) [FLAC]/03 - Gimme The Loot.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 193
  vendor string: reference libFLAC 1.2.1 20070917
  comments: 7
    comment[0]: artist=The Notorious B.I.G.
    comment[1]: title=Gimme The Loot
    comment[2]: album=Ready To Die (The Remaster)
    comment[3]: genre=Hip-Hop
    comment[4]: Comment=.
    comment[5]: DATE=2004
    comment[6]: TRACKNUMBER=03
```

Elements to notice:

1. The FLAC file was encoded with FLAC version 1.2.1
2. The FLAC file has tags that I don't want to keep (see the [Clean](#clean)
   section), such as DATE
3. The FLAC file misses the TRACKTOTAL
4. The FLAC file misses the ReplayGain tags as well

### Considerations

1. The FLAC files are encoded with different versions of FLAC
2. The FLAC files have different tags, some of which are missing and some of
   which I want to keep and some I want to remove
3. The FLAC files have different blocks, some of which are unnecessary and can
   be removed
4. ReplayGain tags are missing from some of the FLAC files and need to be added

The next natural step is to correctly tag all the files.

### Tag the files

> [!NOTE]
>
> The purpose of this step is to add all missing tags and rename files with
> their right name.
>
> If you want to use another tool, feel free to do so.

Many options exist to tag FLAC files but I will use [beets](https://beets.io/)
(v1.6.0) to tag the files with the following configuration:

<details>
<summary>Expand the configuration file</summary>

```yaml
## https://beets.readthedocs.io/en/stable/index.html

# Where to keep the database
library: /config/beets.db

# Where to import the music
directory: /files/processed

# Import options
import:
  # Do not edit the path. The filename can be changed.
  log: /config/beets.log
  write: yes
  hardlink: no
  copy: yes
  move: no
  incremental: no
  resume: ask
  detail: yes
  autotag: yes
  duplicate_action: ask
  none_rec_action: skip
  bell: yes
  timid: yes

match:
  max_rec:
    missing_tracks: strong
    unmatched_tracks: strong

item_fields:
  my_multidisc: 1 if disctotal > 1 else 0

  my_discnumber: u'%i' % disc

  my_source: |
    if media == 'Digital Media':
      return 'WEB'
    elif media == 'Enhanced CD':
      return 'CD'
    elif media == '':
      return False
    else:
      return media

  my_samplerate: |
    if format == "FLAC":
      return round(int(samplerate) / 1000, 1)

  my_bitdepth: |
    if format == "FLAC":
      return bitdepth

  my_year: |
    if year == original_year:
      return year
    else:
      return '{}, Reissue {}'.format(original_year,year)

album_fields:
  my_format: |
    average_bitrate = sum([item.bitrate for item in items]) / len(items)

    if average_bitrate > 480:
      return 'FLAC'
    elif average_bitrate < 480 and average_bitrate >= 320:
      return 'MP3 320'
    elif average_bitrate < 320 and average_bitrate >= 220:
      return 'MP3 V0'
    elif average_bitrate < 215 and average_bitrate >= 170 and average_bitrate != 192:
      return 'MP3 V2'
    elif average_bitrate == 192:
      return 'MP3 192'
    elif average_bitrate < 170:
      return 'MP3 %i' % average_bitrate

paths:
  default:
    $albumartist/$albumartist - $album ($my_year) [$my_format%if{$my_bitdepth &&
    $my_samplerate, $my_bitdepth-$my_samplerate}%if{$my_source,
    $my_source}]%if{$catalognum, {$catalognum$}}/%if{$my_multidisc,CD
    $my_discnumber/}$track - $title

format_item:
  $albumartist - $album (%if{$original_year,$original_year,$year}) - $track -
  $title

format_album: $albumartist - $album (%if{$original_year,$original_year,$year})

sort_item: $albumartist+ $album+ $albumdisambig+ $year+ $disc+ $track+ $title+

sort_album: $albumartist+ $album+ $albumdisambig+ $year+ $disc+

sort_case_insensitive: yes

per_disc_numbering: yes

# If the process is threaded
threaded: yes

# The musicbrainz host
musicbrainz:
  host: musicbrainz.org
  ratelimit: 1
  searchlimit: 3

# UI customization
ui:
  color: yes
  colors:
    text_success: green
    text_warning: yellow
    text_error: red
    text_highlight: red
    text_highlight_minor: lightgray
    action_default: turquoise
    action: blue

art_filename: cover

aunique:
  keys: albumartist year album
  disambiguators:
    albumtype year label catalognum albumdisambig releasegroupdisambig
  bracket: "[]"

ignore: [
    ## Linux

    "*~",

    ## macOS

    # General
    ".DS_Store",

    # Thumbnails
    "._*",

    ## Windows

    # Windows thumbnail cache files
    "Thumbs.db",
    "Thumbs.db:encryptable",
    "ehthumbs.db",
    "ehthumbs_vista.db",

    # Folder config file
    "[Dd]esktop.ini",

    # Windows shortcuts
    "*.lnk",
  ]

ignore_hidden: yes

clutter: [".DS_Store", "[Dd]esktop.ini", "Thumbs.DB"]

asciify_paths: false

replace:
  '[\\/]': "_"
  '^\.': "_"
  '[\x00-\x1f]': "_"
  '[<>:"\?\*\|]': "_"
  '\.$': "_"
  '\s+$': "_"
  '^\s+': "_"
  "^-": "_"
  "’": "'"
  "[“”]": '"'
  '[\xE8-\xEB]': e
  '[\xEC-\xEF]': i
  '[\xE2-\xE6]': a
  '[\xF2-\xF6]': o
  '[\xF8]': o

max_filename_length: 180

va_name: "Various Artists"

# Plugins
plugins: [
    # Autotagger
    chroma,
    discogs,
    spotify,
    deezer,
    fromfilename,
    # Metadata
    # absubmit,
    # acousticbrainz,
    edit,
    fetchart,
    ftintitle,
    lastgenre,
    # Path Formats
    inline,
    # Miscellaneous
    ihate,
    web,
  ]

discogs:
  user_token: "CHANGE_ME"

edit:
  itemfields: [track, title, artist, album]
  albumfields: [album, albumartist]

fetchart:
  auto: yes
  cautious: yes
  minwidth: 1000
  maxwidth: 1000
  enforce_ratio: yes
  sources: [filesystem, coverart, itunes, amazon, albumart, wikipedia]
  quality: 75

ftintitle:
  auto: yes
  drop: no
  format: "(feat. {0})"

hook:
  hooks:
    - event: write
      command: echo 'Processing "{item.path}"...'

    - event: album_imported
      command: echo 'Album "{album.path}" imported.'

ihate:
  warn: ["albumtype:live"]
  skip: []

lastgenre:
  auto: yes
  count: 1
  fallback: ""
  force: no
  prefer_specific: no
  souce: album
  title_case: yes

web:
  host: 0.0.0.0
  port: 8337
```

</details>

```sh
# Run beets with Docker
$ docker run --rm -it \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/Zurich \
  -v ./config:/config \
  -v ./files:/files \
  lscr.io/linuxserver/beets:1.6.0 beet --config /config/config.yaml import /files/raw
```

The exact editions for the albums are as follow:

- _2Pac - All Eyez on Me (1996)_, catalog number
  [314-524 204-2](https://musicbrainz.org/release/f5e7ddad-e38e-4621-9173-6bad2f126c33)
- _Eminem - The Eminem Show (2002)_, catalog number
  [00602445964222](https://www.discogs.com/release/25950724-Eminem-The-Eminem-Show)
- _The Notorious B.I.G. - Ready To Die (1994)_, catalog number
  [94567-2](https://musicbrainz.org/release/f42fe7d8-fa5e-3ee5-9a83-456c8c663ed5)

### Analyze the newly tagged files

Let's explore the current structure of the files again:

```sh
# Display the files' structure
tree ./files
```

```text
./files
├── processed
│   ├── 2Pac
│   │   └── 2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}
│   │       ├── CD 1
│   │       │   ├── 01 - Ambitionz az a Ridah.flac
│   │       │   ├── 02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac
│   │       │   ├── 03 - Skandalouz (feat. Nate Dogg).flac
│   │       │   ├── 04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac
│   │       │   ├── 05 - How Do U Want It (feat. K‐Ci & JoJo).flac
│   │       │   ├── 06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac
│   │       │   ├── 07 - No More Pain.flac
│   │       │   ├── 08 - Heartz of Men.flac
│   │       │   ├── 09 - Life Goes On.flac
│   │       │   ├── 10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac
│   │       │   ├── 11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac
│   │       │   ├── 12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac
│   │       │   ├── 13 - I Ain't Mad at Cha (feat. Danny Boy).flac
│   │       │   └── 14 - What'z Ya Phone # (feat. Danny Boy).flac
│   │       ├── CD 2
│   │       │   ├── 01 - Can't C Me (feat. George Clinton).flac
│   │       │   ├── 02 - Shorty Wanna Be a Thug.flac
│   │       │   ├── 03 - Holla at Me.flac
│   │       │   ├── 04 - Wonda Why They Call U Bytch.flac
│   │       │   ├── 05 - When We Ride (feat. Outlawz).flac
│   │       │   ├── 06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac
│   │       │   ├── 07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac
│   │       │   ├── 08 - Check Out Time (feat. Kurupt & Big Syke).flac
│   │       │   ├── 09 - Ratha Be Ya Nigga (feat. Richie Rich).flac
│   │       │   ├── 10 - All Eyez on Me (feat. Big Syke).flac
│   │       │   ├── 11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac
│   │       │   ├── 12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac
│   │       │   └── 13 - Heaven Ain't Hard 2 Find.flac
│   │       └── cover.jpg
│   ├── Eminem
│   │   └── Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}
│   │       ├── CD 1
│   │       │   ├── 01 - Curtains Up (skit).flac
│   │       │   ├── 02 - White America.flac
│   │       │   ├── 03 - Business.flac
│   │       │   ├── 04 - Cleanin' Out My Closet.flac
│   │       │   ├── 05 - Square Dance.flac
│   │       │   ├── 06 - The Kiss (skit).flac
│   │       │   ├── 07 - Soldier.flac
│   │       │   ├── 08 - Say Goodbye Hollywood.flac
│   │       │   ├── 09 - Drips.flac
│   │       │   ├── 10 - Without Me.flac
│   │       │   ├── 11 - Paul Rosenberg (skit).flac
│   │       │   ├── 12 - Sing for the Moment.flac
│   │       │   ├── 13 - Superman.flac
│   │       │   ├── 14 - Hailie's Song.flac
│   │       │   ├── 15 - Steve Berman (skit).flac
│   │       │   ├── 16 - When the Music Stops.flac
│   │       │   ├── 17 - Say What You Say.flac
│   │       │   ├── 18 - Till I Collapse.flac
│   │       │   ├── 19 - My Dad's Gone Crazy.flac
│   │       │   └── 20 - Curtains Close.flac
│   │       ├── CD 2
│   │       │   ├── 01 - Stimulate.flac
│   │       │   ├── 02 - The Conspiracy Freestyle.flac
│   │       │   ├── 03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac
│   │       │   ├── 04 - Jimmy, Brian and Mike.flac
│   │       │   ├── 05 - Freestyle (#1) (live at Tramps, New York, 1999).flac
│   │       │   ├── 06 - Brain Damage (live at Tramps, New York, 1999).flac
│   │       │   ├── 07 - Freestyle (#2) (live at Tramps, New York, 1999).flac
│   │       │   ├── 08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac
│   │       │   ├── 09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac
│   │       │   ├── 10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac
│   │       │   ├── 11 - Business (instrumental).flac
│   │       │   ├── 12 - Cleanin' Out My Closet (instrumental).flac
│   │       │   ├── 13 - Square Dance (instrumental).flac
│   │       │   ├── 14 - Without Me (instrumental).flac
│   │       │   ├── 15 - Sing for the Moment (instrumental).flac
│   │       │   ├── 16 - Superman (instrumental).flac
│   │       │   ├── 17 - Say What You Say (instrumental).flac
│   │       │   └── 18 - Till I Collapse (instrumental).flac
│   │       └── cover.jpg
│   └── The Notorious B.I.G.
│       └── The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}
│           ├── 01 - Intro.flac
│           ├── 02 - Things Done Changed.flac
│           ├── 03 - Gimme the Loot.flac
│           ├── 04 - Machine Gun Funk.flac
│           ├── 05 - Warning.flac
│           ├── 06 - Ready to Die.flac
│           ├── 07 - One More Chance.flac
│           ├── 08 - Fuck Me (interlude).flac
│           ├── 09 - The What.flac
│           ├── 10 - Juicy.flac
│           ├── 11 - Everyday Struggle.flac
│           ├── 12 - Me & My Bitch.flac
│           ├── 13 - Big Poppa.flac
│           ├── 14 - Respect.flac
│           ├── 15 - Friend of Mine.flac
│           ├── 16 - Unbelievable.flac
│           ├── 17 - Suicidal Thoughts.flac
│           ├── 18 - Who Shot Ya.flac
│           ├── 19 - Just Playing (Dreams).flac
│           └── cover.jpg
└── raw
    └── ...
```

Way better! Let's have a look at the metadata.

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/2Pac - All Eyez on Me (1996) [CD FLAC]/CD 1/01 - Ambitionz Az a Ridah.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 5658
  vendor string: reference libFLAC 1.3.2 20170101
  comments: 68
    comment[0]: ACOUSTID_FINGERPRINT=AQADtEr3pAx0rsWXSEWv40fJzDjuUHCShcR5uB-JM8mDSx2SxXkResd5_AnRJLmDymR27NEkRbhyoTlKIjqF8wKfK0FT56j2QdahMjom6tqRH8cVkriTw0sQ6T_OPSH0oDl31MSXKBY-fji0dUdumMzRX8g1JceZQzwuJsqRE9X0w1WRO_rhBzqXo6mON0pwcwv8CFdX4s7RnASzXsFzIUxE5mA07keYCf1u_IRY7bge5IsevD6aqkSeH9ahk0dTccLDDPcSg8nDoM9mXETzoTIdvKvQSznCRN3BSMnReA_yHE3Fo-6hU8eT5ahV-Klg-odONHpKnFvw8Hhj2I2MazvRHNWN_wl6HWEiZjkqZUeYJscjG9qaxWh-nDriJotDVNTxHbmPSxD5o1uNhyJ-g9ojEd8ORh-aSxTeVQl6IY2Y5cGVHU0T5FFjdM1iaPnxI66icERFHf6Rp4d1QeeORtOMh3jzgNrzCY-m480FVyTU9ArUH80Yqvgegpr0odeFj0flQ0tyfC36G0134vGhD9fxz2A25bDh_MKZC81ZPPJE3Nl4_KtwHpxyNJF29PLwEaKqI3_wRVGOR0NTacV_-EVuHlp-5FNYPDwOO8ksnD--5WgeHeEeTmDyHFkYbkWpo_GRKzUaTQs06rgTPcjTo6IKPx9SHbp0mGPRTSMeGr_hl8SlUGiOyjT-B8x1ZMm6PCiDMM3xpdA2ZUeYH5e0HE941Dp8SsM35O2h7whNhcG1Cwf1EN8VgeN89FI2NLk4IcvuItO6qHheNN2RTzY-JYaoHhcv5FmO5-iRytHh59C5HI0mPFFy3F9C2MYtDtfhmMeDPxfCrDzKHE2P8MrxjRC_I9ZxZTkeoak0oxqtoS8Syj_CjMclpTgTPLtghzGupJpxpJdxpxOa_ajGrUclHY314FEMURPxPjilI25CFxVVIn3g_5BFotHEEl-ON6GI_iKahxJxqUfzo_Lx7glKSjPCZFxUVDqaRsidQnt2hJmPLw-eJBnaSPDJSPiRZiS0VNKMPExxZRd-eBONHz_SPDt0VcGlK-joo4kSQW98oXfgZxu-I4wfPPxxoUlYCvWW41Q84z3CC3-gWTo-6UE7KcdlcDkecGfw0dgqaI5ynD-aZUwohFee4KnRdDnCk8cT9MwxJZl7VIUnNcLRhkrxBXmyFdePJieJ7kH2XEKYiD2ylN_RRDp65dAnGVaP3MejLA6eo6la3BHSK90hqszR49GDN1kKJkc1ZSJx5WiUX-iNP0OYSw8YMsuPZkHoRIeeotHEI-Q9fMqe4NFQS-Ph70f6HDqPRvpxTsGZI7zRvA8uKSfCHF1p_FXQS0eYkXmMSuHR9MgVpaio7Giie9Cp46KDNjp8OPlxeK-hUQrxQ08hfjgqyRkalST-BJxlDVkyhheuHI2GXJGNO0rgQ_eDLwjVI1dx5tB1fEcT9ngWfCn8gz8u_WjOo-uDPxeaMTy-H9TEo8-D6zrE6hQePcKZREFoHSbSjNnxDzqroKm4GU-Grw0atJJYnEcz8kTXB9xzoTlzIwvDdQgz8egv_II46cFf5Jrm4NKHptRxGmneydCzoWGkCw9FxOYldDeaZ5FxaTmaPEf4cbieoDJzI8zCHeHCowl5nKqhpUethwlyJjv-oqIy-EP6HxcJ8XxQa-nwLA5u-DmurUXzo6fwJ-iuGWnELMezH02z46kRuoGWE9eiB3FN1Dp-5MwPHZZ3lDqeLPiRik3xMwh9NPEu4tIY_B7CcFm44zpS7UJfE7WkQEvO4nwQU4mCjocrDdF-3D8-xoKihMaX40HzFGV33EdzBtUfkNWFMGNvMOUnNBMRXjduCa4OVeKDfNkjvEb9w3LeIv2FR-mg5XrxiQ--LBXRSEf48LgqjmgQ6nFwHWGex6iUZUdIoRGHt4bY4g--KEcuHbWKVA58qNePpm3xnLiO3misHc-P8kQTH_-DXke6LqPg5ahz5IvCo3l2dPkDXTwuZTzq4cMP2Rn8o9463IyP6_CJH32WEE03SvgfhNFzVOOC83Cb4Ig_iFV4fJJJ5MJF9BlSfXAuByK5o9-PizFu-MGP1gkazcaTg6mDMM9ykH3R8wjTZAl0OLzxRTny40x49IfzID108miq41qMI_3x6sMj52iO6g9-oTlzlOF-hIl09Dr0ZYd_5H9wJTseHk1VPDpSHd_CG1qKa5lxNE9RVTse-Wg-onpM_E_QM0cYhdnRRPkROse5C00ObX_wLkfM7GijIo_wH7J5NDGPKyPxMMYFKsepo5SPJi_xhEpQ_kaaaFkOTrTQNMhfozq0KD_O5cGZiUfYiEU-_OugxjvMH7XC4T0OOz3OBDUlo0lePFTw6wgzZnmMK0Naodds6IkCaz2uB3mUtHhu9DXyHDp53EQz-XgX8LBd_Mcjo3kr9IsT3Ajz8AGbHdaO8KGNMxTEicd95IviEJeOWkX6qPAPneTRqMeV5biPJq_QH5cyHY0-MI-Y4BdKMiuyROJ-VBIaz8KpaUeYo9ShR4lxfkcbFT6sc8cPLW2OZ_nwjgHz1GiefTgXotmYoLdynJ9QpYwML8s6pDt6Ucelos-NZsb3FKESKdAzF_xp5DTSlcR_nNHhP8lR-ZDzhGA6Ht5-9MWpTQp6HR6zHOP2oZk__PgFRTlqZ8MTIZcofCmyJM0-dPkEfccZCb_Q-EJ_CtfR52hUGt-DH2Eiho9xZUdT4TVySTmsP7gOPQoVog78ED-ckYd-ounxBe8Dyk3x-LhcNHnQP_jRjMwRdll3MIlEVCEd5IfLQ3_xbDlyHbV6pHkO6zr07GiIT5xG3Fko9IyFLImWUEZ29UJ_NCqpFD-aS3kQsruKMIyOcjmeQkslHf2FXDyeZBwqSXWIB_nx89CSaGxxRalw5jjSZ8W1I9TR6Cd-OWD4rNDCZUdI6WDSBJ_QVF9QHp-SJDk-SUYuIWHM4DmeBFdm4UoBwohABjFgBCAEIUKwAQIIRAhAQggEOACAARCJAUgAAYEgSBAiDDFCCIkAUAYAAYwQRAAEEAACAiEAEMIIIA0AiBkghBAACaSEEIIDiYQBRiFBgBJAAUqEMkIAQwAmzhCCCATAOkAYMAIJAYQBADNDCEUKMCUEYcAI4YgCQgjDCCBIMYmAEoQwAAUyQjABEAACMAIFQERIIBQAwgEGGFCAEYAQEAAgwY1zQCAkgFiCAQMIFgIYAYwDxACigBEAMEUUEIAYJAAgVCCkCFEAUQaNEEIBABwBAgwABERACEAoAooJwpBgAFEjiABWOESIEIYAJgBAgCAkjADMCGCYMEKAhABySDEJmEICAIIEQowYARQxyBIHCQEEgmcBEYIAJIQQAwjREAHCQASEABYIogQyiggBhDUIYIQIURIAAYCDAAiggBDAACWAAAwoIABQgAgBhEUAEKKkAsYqCqAwEwDgEcDAOGoMdcARwrQCRgGAgDBAAGEUcIAAAIBBggDJlDIACiMAEA4QQBAywAjliDJMEECUAYIJIYRSDAigFDEECEAYRYIooCQAgABAKSRCCkMMkMQqCpwAgBEBAJOAEEWFEQgYB0QiSAAChBCACGQIAAABAxwAAhggBCFMICAYkQgAKggADBilhAFCGCGAIR4ZAAQwRBgAJUEECIMEpAAAIggjSiJBEGGEEAEQRQCAAIQQRiwECACKAWCGIAosAYQgiAAogBAAAUEQMAIQBJRyggChRCLCKAKEEEIghBBxghDgGDDCEaRAEkogKAgwQABiABAOACiEAmAQI4EQAAAkiGMICEYFAIQIAJQAggCABWDGGCOAQEIYhgBASBACDIJCACAwVAIAYgAghAGlAGFGAKKFccIAhwRxBhAChAAQICIsMBAwIAxwhjnHlABESCIIAgowAA0ARBkkBIEKQOMQcU4BoIhBAAoBKDJAGUAIAlQARAQAAAghAQBCEgoQMUQ5KAEQiCGQCEQCACEEtoICICSzIBk
    comment[1]: ACOUSTID_ID=bdf0a406-30df-4521-9161-0aab63c5baee
    comment[2]: ALBUM=All Eyez on Me
    comment[3]: ALBUM ARTIST=2Pac
    comment[4]: ALBUM_ARTIST=2Pac
    comment[5]: ALBUMARTIST=2Pac
    comment[6]: ALBUMARTIST_CREDIT=2Pac
    comment[7]: ALBUMARTISTSORT=2Pac
    comment[8]: MUSICBRAINZ_ALBUMCOMMENT=explicit
    comment[9]: RELEASESTATUS=Official
    comment[10]: MUSICBRAINZ_ALBUMSTATUS=Official
    comment[11]: RELEASETYPE=a
    comment[12]: RELEASETYPE=l
    comment[13]: RELEASETYPE=b
    comment[14]: RELEASETYPE=u
    comment[15]: RELEASETYPE=m
    comment[16]: MUSICBRAINZ_ALBUMTYPE=a
    comment[17]: MUSICBRAINZ_ALBUMTYPE=l
    comment[18]: MUSICBRAINZ_ALBUMTYPE=b
    comment[19]: MUSICBRAINZ_ALBUMTYPE=u
    comment[20]: MUSICBRAINZ_ALBUMTYPE=m
    comment[21]: ARRANGER=
    comment[22]: ARTIST=2Pac
    comment[23]: ARTIST_CREDIT=2Pac
    comment[24]: ARTISTSORT=2Pac
    comment[25]: ASIN=B00000163G
    comment[26]: BPM=0
    comment[27]: CATALOGNUMBER=314-524 204-2
    comment[28]: DESCRIPTION=
    comment[29]: COMMENT=
    comment[30]: COMPILATION=0
    comment[31]: COMPOSER=
    comment[32]: COMPOSERSORT=
    comment[33]: RELEASECOUNTRY=US
    comment[34]: DATE=1996-02-13
    comment[35]: YEAR=1996
    comment[36]: DISC=1
    comment[37]: DISCNUMBER=1
    comment[38]: DISCSUBTITLE=Book 1
    comment[39]: DISCTOTAL=2
    comment[40]: DISCC=2
    comment[41]: TOTALDISCS=2
    comment[42]: ENCODEDBY=
    comment[43]: ENCODER=
    comment[44]: GENRE=Gangsta Rap
    comment[45]: GROUPING=
    comment[46]: ISRC=USKO10403591
    comment[47]: LABEL=Death Row Records
    comment[48]: PUBLISHER=Death Row Records
    comment[49]: LANGUAGE=eng
    comment[50]: LYRICIST=
    comment[51]: LYRICS=
    comment[52]: MUSICBRAINZ_ALBUMARTISTID=382f1005-e9ab-4684-afd4-0bdae4ee37f2
    comment[53]: MUSICBRAINZ_ALBUMID=f5e7ddad-e38e-4621-9173-6bad2f126c33
    comment[54]: MUSICBRAINZ_ARTISTID=382f1005-e9ab-4684-afd4-0bdae4ee37f2
    comment[55]: MUSICBRAINZ_RELEASEGROUPID=e2621417-9236-36b4-9f9e-376c416dc4b0
    comment[56]: MUSICBRAINZ_RELEASETRACKID=5ad65fe1-fc9a-36ac-9b9c-e8173edcd197
    comment[57]: MUSICBRAINZ_TRACKID=ffd277e5-bd7f-434b-ba95-4b278013e81a
    comment[58]: MUSICBRAINZ_WORKID=b540da92-e52b-4fd8-b8e9-6c17ceaba2f7
    comment[59]: MEDIA=CD
    comment[60]: ORIGINALDATE=1996-02-13
    comment[61]: SCRIPT=Latn
    comment[62]: TITLE=Ambitionz az a Ridah
    comment[63]: TRACK=1
    comment[64]: TRACKNUMBER=1
    comment[65]: TRACKTOTAL=14
    comment[66]: TRACKC=14
    comment[67]: TOTALTRACKS=14
```

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/Eminem - The Eminem Show - 2002 (2023 2xCD - FLAC); Expanded Edition/01-03 - Business.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 5561
  vendor string: reference libFLAC 1.4.1 20220922
  comments: 72
    comment[0]: ACOUSTID_FINGERPRINT=AQADtJISJZKSUQl-I3yMZAkzPRK-w0qWHNeJL8SD9Cl0JcijHG2Np5ODStSDH3mOfjhD_BD3QJcR-sdbPHmwJw92KSgZbL_QJ8UtPQGzLxwR5mMg-UWcR2jqHtOb40uKf3jOHZF7NFU2QXOOO0_wd-CVQD32HpeGM4PW6EezU8iTC2eDvjmu59h2whLOTHCWlMOVg68KX8dvmBEP5h7iP6gPUZs-tHLxC7UyNM-xE7qMrhOcHj9a9XhwDj3R6MGeHeeR86iLFPrRKWnwo88yND9oPtBTHrWO45BseIpoQ6P44YV55Dmeo07qoNl15B905C9e1NPRvEfO4z6058iPw_3RqseH8_BTnMeNPviRQxeFMBp-3MwD88ajWbh-9I4Qe8em_IHWGJeSHLkAP8ez4jy-oznJoBaDPIceEb3g_tBkHvtw5bjWI_yFdHph7RApo96Oo3nQecQDbWSIS_Dx7nh4aFOOk2iiB_3x3CjX4ypUPCfcKeiHXDSa_Kgc7BfOGQ2Vo64m5OWh28JT_IMr5thzfMR1aCMuHD_cLRY67WB-mFF44aGLB3mJH4-kI62PLtSh9QeTnRs-NMZ544eo6_B-1JFS4RGiHj6PZ3mGD6_RZ2g8BfEf_IGOa4dHBk9w7sV-vIYmD907PNWRT0fzoxsVhOdx58HV6PhqqJGEo8eF53B_fJCngP8RPoZ_nIepHVR3tGuRwcKl4mmHdofDHf0HlRuuSDe6iRmax8ehfuB5lDxS80Ksoz90nBf6LGjq48HHQxM14SR-2Ic2LQ90lD-a9XhavEf3B_hxCcrUQ-OOP2ijLEMa4ce_wOKOSi_yQx2jDI98eAv240pWBr-gWTms44X_wi7a8cZzNM_ByDiPci9-Oci0HuHxKcTl4TFEpcrx48KrF1rE48d7tDpunOhDNMfz4g139HB1oV0O_cR9VIt2NM_xQE-Oi0PIH_aD_Ogk5sSl4yyah5BCHo-gk3nAbMcvaAmHJ9F25D9-ZIqYg1ePK8sFfZ6RS0dnHu5Rbm6RKyxx9IaW1_hJaI0Gu7D84IczjkUP3QmsZsItnDt0nejowCFyhNzxCJX-gMkf4WHwaCHq6GCPV7DyEW3SC9dn3IKWM0dJQl_64PrQMRQc4rg-qWiVHZeN7cf4SMJD7ojU7Gg2JsetPCFyQRceZSFcZTGOPcfTo0xkQT8-8bB-oceTFNfh5wv4Dj9RJh1BZoHz4w-hHiGTpEf5C81zkOVRi4G35MicQi5Z5KgrY4-OSyX6A36OL1k0wVWPHvF-6D8mrggzHfkVfDx45uAFbdlxHcePxr7wGCf8LAe9DDeJjofDDw_0HFV93OiaPWh-nLsg3Tz6ID_8a8iLyxj3C9_h16gc5KLQp2hOBvuDxMOXF-HDFM2Pf6CJ0E0AMWHwwwd1olV6XCSa5-iz4SsqNjyeHdcrfE-gH1ocHn_QpEf440d7PJPx40f6VRhjiGMnVMuN3ER_-EX4SCtO_I7xH-_RUNqDqWNwfTgzdB_ECu8J7_jxEk574Vrh53hXvLCkEKeC73h2aMd3NFV49Dh0oTzy1PiJ8zgixaKDCw2mHz1ctULHQ7mPpzl6NGJz_Nh3fNCoCukvPOuRDyf6B3pRqxca5ccPdXGPPoiv4xKRSvgO-bhSwy9x7sDVQh8e_Yj8GN-DxloY5BFxPEF_uJI8nB_cGz-PZ9DiLDneK9iPh5nQXbggnriRjuQw57DDbOh3IUzFGyUnos2EG5aPv8JFaD7O1EGtHFOy4SSLK4P2CTaPdC-OSvOIyT42NUf-ZoJJQj_OHy40_fihq-hzqJJC-Ee7EEepCnl25PjRDE-PCx7D4zye6egZ5ER6HX8hH88hPXtg_qiT5MhKDfXxwz1RzivyGIqL_PBj4cyMBx8jhHmYQZES6fiPnwi7oUl7_MUffEVd5XCqI8cPrVEvvPjhfBdIWdBfXBGNHo_wQGzRSR_4sUMoHw3XCbmiCH0hhtLwIz6DavGNfLh2CqIq5YjFE4-u4eJRUT3ya9Ck7kaeOIzwEtfRo-kDfpsk4aWC50fSzEf4NPh_lNqF7wwuTcSeoE8eTBk1nMeFplp-lOSRFxdK9hDDCLkctMSPlw6oUmVwDRdWoTyaJ_qIl4nxCw-FMNGSRUzQRjL0o7GPPsJIXgl6PNoIHX5xCn9Q8YWu8PiEPlvQ5MH54HzRMTvuB-Lx5dAeHpdjaFKP5saXHE-O2kd6WMcfuBGH0I1x3mhiGloXHn_wHJp7wbyQ-_hD4zq-4_nhH_-KC6l4iMiHcy_COIvxDDW0SNzxT2Ce44d0-Bd6VqDv4FmOhtBU9TiFUD_y-oWfQZ3KgvkQUoJvoSSDLyFCEfYNRfsRU6TQP2h6BruPHk31oR8cbtXwd0geCvnxw8yHjjXySEcTGlqJ3Ebf4OXxa0hraM3xy-i1GI2m4yw0jcFtOD1EFX-O6FOHW8dr5GGHWjS0pzAjHf3xTMePRxnKBeJRi4-wH9-JcsNd-A-hH_fR2Hi0o-_RnCUqZx08_aiaHT-69dDS4yH-w8yOX3gOLVWNZzhh_cjPBe509CmOdgZAlWCEAOHIBYQCg7BRkBgokRIACKCscAgAIpwCACEoCACUAiQIAgJAAoBRBACHEHAMESeMFAIAAwAFCACgiCQiCAKQMgoQKIJARhGC4BkMMiE4UEYJA5QAglAEgAUCAQEQI8gQBIQkwxkgGCIAMUYEACAgAixhACgBACAAIQCEMhIwZIgiBAlDBEGGCSQoMEgYgxQgRDkADDUIgbKQUcIoAwhBwEEACEKAKSEIMBwIBgASwBgIFCAAKaQMEcIAKhwwBBgKFjMAGAGQAgoCYxwCSgkDgCABIGGIIkIJIZgBDlhnJABIGQQIdQYhAoBEBiGDBBDOGWEEMAYYIBACDBDJAaHECEQIcYgAi5hAxgGkESJSKQGQS4gAISgjBFHgjCEAEkMEYQwpJYQAhgAJCAFEGIAAYcYAaIghAiBhjLBOECMAEIAYVIAjQhgBABDAAOIOEsooQQRghCgAnAMGSOKAEMAQAARhxgAiABeAAIKYAIoAQAggCQLDgUCQkWEAAkoYAQgAyBBtAEAIGICQAQoggBxAUAIEjHJEASEcwIYJogiiBCpikCFAAAQECIgIIwRAhkBAoIEEIAQIskgYBYATygGBEALAJEKYE8IAAhAEyAEkBDNAASGEowIQ4JUxCBmBjAAICEY8E1IwIAwxTBABhUBScUKEQkgIJwBBQigjjBBECASQMBoAxIAAjhqnBAKMGQG0EYoxAwgAxACGDFEEEUOQogAgIBA0zgkjEgQKMSCQIEiYApQBgjpgFQCEAUgAYEAopIQSiAEBADAAECMIEcAAIADRAghvCEAUAIGUJIpRYIVRgEgigBDGGGAIAghAIAhCDhgEgCWQKIGYEoAQg4AyDBggBAEOIW8AVQgAJAFEBgAmDCUCACGRYcIYAgggxghhjKbEKGEEEII5KIChSCGgBBEKCIKIUFQ4EpASRigApHFCCgCAEgYapwQBCGCGFDCGAAMQ
    comment[1]: ACOUSTID_ID=e8ae92e6-3b32-4dd2-bb50-e3cfc6989167
    comment[2]: ALBUM=The Eminem Show (Expanded Edition)
    comment[3]: ALBUM ARTIST=Eminem
    comment[4]: ALBUM_ARTIST=Eminem
    comment[5]: ALBUMARTIST=Eminem
    comment[6]: ALBUMARTIST_CREDIT=Eminem
    comment[7]: ALBUMARTISTSORT=Eminem
    comment[8]: MUSICBRAINZ_ALBUMCOMMENT=
    comment[9]: RELEASESTATUS=Official
    comment[10]: MUSICBRAINZ_ALBUMSTATUS=Official
    comment[11]: RELEASETYPE=a
    comment[12]: RELEASETYPE=l
    comment[13]: RELEASETYPE=b
    comment[14]: RELEASETYPE=u
    comment[15]: RELEASETYPE=m
    comment[16]: MUSICBRAINZ_ALBUMTYPE=a
    comment[17]: MUSICBRAINZ_ALBUMTYPE=l
    comment[18]: MUSICBRAINZ_ALBUMTYPE=b
    comment[19]: MUSICBRAINZ_ALBUMTYPE=u
    comment[20]: MUSICBRAINZ_ALBUMTYPE=m
    comment[21]: ARRANGER=
    comment[22]: ARTIST=Eminem
    comment[23]: ARTIST_CREDIT=Eminem
    comment[24]: ARTISTSORT=Eminem
    comment[25]: ASIN=B0BLTJZ3LB
    comment[26]: BPM=0
    comment[27]: CATALOGNUMBER=B0035988-02
    comment[28]: DESCRIPTION=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[29]: COMMENT=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[30]: COMPILATION=0
    comment[31]: COMPOSER=
    comment[32]: COMPOSERSORT=
    comment[33]: RELEASECOUNTRY=XW
    comment[34]: DATE=2023-01-27
    comment[35]: YEAR=2023
    comment[36]: DISC=1
    comment[37]: DISCNUMBER=1
    comment[38]: DISCSUBTITLE=
    comment[39]: DISCTOTAL=2
    comment[40]: DISCC=2
    comment[41]: TOTALDISCS=2
    comment[42]: ENCODEDBY=
    comment[43]: ENCODER=
    comment[44]: GENRE=Hip Hop
    comment[45]: GROUPING=
    comment[46]: ISRC=USIR10211053
    comment[47]: LABEL=Aftermath Entertainment
    comment[48]: PUBLISHER=Aftermath Entertainment
    comment[49]: LANGUAGE=eng
    comment[50]: LYRICIST=
    comment[51]: LYRICS=
    comment[52]: MUSICBRAINZ_ALBUMARTISTID=b95ce3ff-3d05-4e87-9e01-c97b66af13d4
    comment[53]: MUSICBRAINZ_ALBUMID=8cbe48f7-a526-46ce-b10a-da3e9526c86a
    comment[54]: MUSICBRAINZ_ARTISTID=b95ce3ff-3d05-4e87-9e01-c97b66af13d4
    comment[55]: MUSICBRAINZ_RELEASEGROUPID=e9585ed4-d148-3711-bbee-55a97b58325a
    comment[56]: MUSICBRAINZ_RELEASETRACKID=6b7929b2-e47c-483d-bdf0-2f448318a618
    comment[57]: MUSICBRAINZ_TRACKID=df417b58-3b8a-45ee-9cdc-0b131538b147
    comment[58]: MUSICBRAINZ_WORKID=1989329c-a912-30da-a223-6aa77e9ed83a
    comment[59]: MEDIA=CD
    comment[60]: ORIGINALDATE=2002-05-27
    comment[61]: REPLAYGAIN_ALBUM_GAIN=-9.54 dB
    comment[62]: REPLAYGAIN_ALBUM_PEAK=1.000000
    comment[63]: REPLAYGAIN_TRACK_GAIN=-9.49 dB
    comment[64]: REPLAYGAIN_TRACK_PEAK=0.997559
    comment[65]: SCRIPT=Latn
    comment[66]: TITLE=Business
    comment[67]: TRACK=3
    comment[68]: TRACKNUMBER=3
    comment[69]: TRACKTOTAL=20
    comment[70]: TRACKC=20
    comment[71]: TOTALTRACKS=20
```

```sh
# List all tags
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list --block-type=VORBIS_COMMENT "/files/raw/The Notorious B.I.G. (1994) Ready To Die (2006 Remastered) [FLAC]/03 - Gimme The Loot.flac"
```

```text
METADATA block #2
  type: 4 (VORBIS_COMMENT)
  is last: false
  length: 5461
  vendor string: reference libFLAC 1.2.1 20070917
  comments: 68
    comment[0]: ACOUSTID_FINGERPRINT=AQADtImWKInEKMN11CKmtcYdjghTlM-Og_PRUUmCJlyE7IGe5Uh3Bk-Nb8GIHmGOS08hS5KQ64Jz40d0-DVKoXGDHzmPJmQSdKwcxD80Jjq8E_nxJth-uEd0oZkf6Emm4IFm8bg5VM-hWRkPJluyGDl-VJQEJlqUJPhxPgizPNB8PAkRGxOM_keuJXJwKUJ9NHGPPdAOnzEipTeu4rlwHdVzdJIu5PIETU_RhNuQ_7getOGDPFowPcGzkOgF1xdhmSOq49lBcVkG_cjxRkKjNcGf40j3QGNi48mCUyiMXiF0KTzOKGjYoaICHj5ehJmMPsaPC06e5HhVMUROeMkX6E4xvgifIK2OR5jOGBeTo9VYpN0CnxM6L8d_NJOYIcsLHbHRiBu456j3IOyg5cFfXKiwoT66kNCzJMGt48Z3TOeDfEF5NN2CH8cfIRwfaIl65JmloB_84xWM_IKOdLqOVsvRf7BGJaic0eDR9AvKhcdx6mjEJSlOisYPlg-8HHWuDKfCY_cRQsvx41GUqAmh60jDHyemHH0r9DuaaUuRHDmaPDa6S5KMSw-mRImOH_GVTLhoQmbyHN-Rewj5xFCV-HiSfEQaStrwEWPMFNdhNxPyQxURbosCh1ES4kbzB88eXBdC1XCHHcV5PIwq5FSOWowNxzz8nGjXoiyabmnwfkh-hD8qnsPF45GGO4UfHneO3IeiKAsvtHmIP_jhKsW3JMcpiXjzoZL0IIxz4cGNqiE-JkGP_grSiIeqJIerInyKHs2OR3yCR0lSPNNh3RidxcgvYjt0EWGWo1yd4see4G2O-0EkiSa0LLmCe8ezQFPzw70Q8rgycsngSkqFc8SVGKHcg4lCVB9qkRu-I8ePPutxwuXR7gGTPkSf8Mg1JIP744eX0QifJUKVOUgYLWTAK0dz5cGP3ELzZQPjz9BPhF0YND8mNHZuHPfRZw_sE6dUpPUJfVGEprowLxSFcMqOJurwox-PpmLw4VaCRnkefFKNSw9ywYfmHTsq4ifMQdeSPOg4BTfqQ9l68BZ0Bn2QyzFuNN7RC89RMhWuPIbTD3nhp5AQ_gzeFZw-hO0RZpOS45SKJj-6MkcoHo2bsTiOJg2Dbr7wgwl3BlUu4T72aMEPPaeE_MG1BHWSCA1-wQpz_F6SoWF09MelLDiLHlNaPCzyHGEmH5WyI3wO6zh15EzxbSmSaQ9y0rgW2M1xLpbx5EGlB2H6QfoRxNbwDXwOLRGvoD3yHGEmB5M-QweOWmGNKQ7kcagu3AdzPviWoyGXQ-OK3oGvHFyOq2IKPdGFf8HWJtBWIk-E8zl-CbW4B1OiNUezZEcviTmavEdXHqeSE6FTEs0DHWqqBw4jCSehPwjDow_OJKxw6UT3Gc6PH8-O8Dn0K8jxRBku-Di94MeVF06WNEj6I7oK_8cv42sOR1KCk0doytCXIW8OL0clLUqI89iqRGRwMxXeoIkelDJ1HA-qNrWQaoZ4PYgVH-WJMxN2ZUEO7ciJMhFjlEtmwXqCmsrxH9GPjymaH7VTHE9xXfihPAl5hJ_Q7PiyHF94NEctRPoGLZ2HWJtQjUqCZuJV3ONEOE_BSjlxF94RqOZBckqIenpShImP-3CUnMgTNriPY8SZCiOVF1q-Dl3GIOaCCx9z_GjqoJ-OFx_i7C485lKMR8eJ5sW05kf4Hzmh0UZEJVqi4tSH5nnQNyGefMXVAdZ3ISaSD1-MUbnxKcPbYMp5hL-gN7ho9Lgk1FPSQY-Oo4mEdiF-9IcZH81-5PuRBz9-qGodRI7ywI_w4N-RHzrC9EqKiSd-NBOPZ0W9_Mh74ebRTMoyhNyBWxI6Kkl6HH4Y4VkuWJFyBZOYEM-CHv5ROuvxqA4eJfqQVDyRHs-xP8HlTME5OXhyDpEXnVDdE88s9IGW4wlySseDytIRTYySGRqPeChpCXapIA9uHVUaHXjqI90u3LhNNGJzXBaRR3mgH2GPelnR7MJ7YUeeB4rZMIGr5MiPuEbT5fjRcFmOLPzRlEpxDS84PcMTKs0QU48gbke8GH9muHtwOsMRSz_0w9ORPcllnFrwKxVOnQivE82y4zpqHU3yQxdOHmHGPWjD40-Df8iO_hBKMsaZWZCfI813_MT340Se-xAnBzk-_Hpw50eYKDM07gpyb0flLPCPJKKOPMKlRtglNGeLUsqFRy_yKcsK1Uf4oKlI4lfxRsFZXHlxjQ12Hjm0VET4hLiSf_jFoPkxh1HgSctx4nRK9NFxb3h6_EFT6qjT4NeJP42QWP5wKg5SN8ctHh_ToMmPL0rwUsIjorLkI49wJUuFSvYELZqOu0f2oXE6lMF9vFlSfIsM7oF_Ir8GzUeYRsR4HMcTGc0mHuEm4Q-2R8FfaOEYD6WGxqJR2_jR_CmSEWclXMmXg0pGJbh5MCy8jOnwHwXjDU3CprgSsPaR58KnRmh6Ej9EhVkqfNKI_0W24_jxHf-QHp8ChsfRKQ98-ASfH75G5D3Oof0j9LA8tLnBSIetCc-N5j3yfdD3IXY59PDEH2k2akj2ZMJZolVyNH2O6oSVHvEPJpOp48ORS6hCJg9yQn3FwOOFPFfgpkHFDUei6AhzXGo3fEqSB025Ex1-NDhTpD1CPvjxo3l0g1HPIG-KKw80XoGPuDr4wz1CfUrBSnuCVklINP2higcTNSOySO8RpjfKEz-eH1pCDutzpH9w0cSc47RRRwuQKEeaFyVbPEpyVMyUHj7G-Bh8H5lehJRz_MeTo7lyxIdSqkSe6MM_7ISf4zrS5oOmbwxyjgNlgCDKcIAMQo4SZBBy0gMKjABIIAMAoMQS4YgwgICkgCWKKSQMEAILowEwEEBiADAGIcAIYMgAgIwxSAAAEHMKWAAEEU0pBAEAghADiBZEIAAMA1wQRBFQSgCjkBKACYUgEAwJqBQwBgQiiGLGIICAQ4oRoAxizAkkADJaEGC8AYoA4gQSQjAtgAKWKKOwQMQQIwk1ADlAFCFKASGsQgYJYYgiQBJHABAECUKEJgQh4ZwBiBiiDCYMSgWAQAYhRgRxQCEECPIMEaGMUoogQogXgBBAkRTIIECABEYAQIRCwEliFBAAEKMVIMQAAABhCAuKEACACQUAYEgABACAIiyHBEJAQAKIkAgoAoACkDCjFAFGIIEcQMwhAAQTgCmBiHEGAS8CEUIgoBRRQgFDEEICCGYAEYAQARDiAiEFLUBICKEUYwog4xBSBgSkBGSEAKKMEQYCApQAQgDFgCEGCGUtEABAJwwCigAgKCESGUeEZUIYoCAgilCFmEHMEHAYEcAhKABBWgkjCABIAQIIAIYYJogAhBiDTDNOMCSYEAAJ4QAS1DgKgMLIGGOEUgIABIQCADlCABNEAAOIUgwo4YAAQjCLlDJMSAmMIQSJAAARTjAGiBBSACEIBIIoJIQhRBAFARAMCwAEAAQAQ5gxyABDBBEAYKAIgEI4RJBgBACvAAHCAGIIIEggRAiBAAtogACCGGIBAIgIYJQABFAqmDFECGAAIAgYQ5CiABoDACACCAsEIUqEUQAggkhoFDAECMAoAEAoxATAgAiAnDDOKIIRpQwJAYABAAEgGVEKCAMIBEIgRIBxyghAwFIGGIHgwgQAAABjQikoCCAAMKSgJU4AwCQRQCiljABEIOCAw0YIRBEAQCEACWFWCIIggMgaIQA
    comment[1]: ACOUSTID_ID=2311f60b-1471-4008-b99a-4032f91c043b
    comment[2]: ALBUM=Ready to Die
    comment[3]: ALBUM ARTIST=The Notorious B.I.G.
    comment[4]: ALBUM_ARTIST=The Notorious B.I.G.
    comment[5]: ALBUMARTIST=The Notorious B.I.G.
    comment[6]: ALBUMARTIST_CREDIT=The Notorious B.I.G.
    comment[7]: ALBUMARTISTSORT=Notorious B.I.G., The
    comment[8]: MUSICBRAINZ_ALBUMCOMMENT=
    comment[9]: RELEASESTATUS=Official
    comment[10]: MUSICBRAINZ_ALBUMSTATUS=Official
    comment[11]: RELEASETYPE=a
    comment[12]: RELEASETYPE=l
    comment[13]: RELEASETYPE=b
    comment[14]: RELEASETYPE=u
    comment[15]: RELEASETYPE=m
    comment[16]: MUSICBRAINZ_ALBUMTYPE=a
    comment[17]: MUSICBRAINZ_ALBUMTYPE=l
    comment[18]: MUSICBRAINZ_ALBUMTYPE=b
    comment[19]: MUSICBRAINZ_ALBUMTYPE=u
    comment[20]: MUSICBRAINZ_ALBUMTYPE=m
    comment[21]: ARRANGER=
    comment[22]: ARTIST=The Notorious B.I.G.
    comment[23]: ARTIST_CREDIT=The Notorious B.I.G.
    comment[24]: ARTISTSORT=Notorious B.I.G., The
    comment[25]: ASIN=
    comment[26]: BPM=0
    comment[27]: CATALOGNUMBER=0249-86280-1
    comment[28]: DESCRIPTION=.
    comment[29]: COMMENT=.
    comment[30]: COMPILATION=0
    comment[31]: COMPOSER=
    comment[32]: COMPOSERSORT=
    comment[33]: RELEASECOUNTRY=DE
    comment[34]: DATE=2004-07-12
    comment[35]: YEAR=2004
    comment[36]: DISC=1
    comment[37]: DISCNUMBER=1
    comment[38]: DISCSUBTITLE=
    comment[39]: DISCTOTAL=1
    comment[40]: DISCC=1
    comment[41]: TOTALDISCS=1
    comment[42]: ENCODEDBY=
    comment[43]: ENCODER=
    comment[44]: GENRE=Rap
    comment[45]: GROUPING=
    comment[46]: ISRC=USBB40580807
    comment[47]: LABEL=Bad Boy Records
    comment[48]: PUBLISHER=Bad Boy Records
    comment[49]: LANGUAGE=eng
    comment[50]: LYRICIST=
    comment[51]: LYRICS=
    comment[52]: MUSICBRAINZ_ALBUMARTISTID=d5d97b2b-b83b-4976-814a-056d9076c8c3
    comment[53]: MUSICBRAINZ_ALBUMID=eca5a076-25f0-3d45-9e92-1d131321472a
    comment[54]: MUSICBRAINZ_ARTISTID=d5d97b2b-b83b-4976-814a-056d9076c8c3
    comment[55]: MUSICBRAINZ_RELEASEGROUPID=5afcfeac-118a-35e6-af0d-35ec9003354d
    comment[56]: MUSICBRAINZ_RELEASETRACKID=d506b393-1420-3a28-b05a-d34227fc85e9
    comment[57]: MUSICBRAINZ_TRACKID=e86ee6e1-ccd4-46fe-9781-cea01aa5db92
    comment[58]: MUSICBRAINZ_WORKID=45d3b990-7e44-3369-98f2-a79906450680
    comment[59]: MEDIA=CD
    comment[60]: ORIGINALDATE=1994-09-13
    comment[61]: SCRIPT=Latn
    comment[62]: TITLE=Gimme the Loot
    comment[63]: TRACK=3
    comment[64]: TRACKNUMBER=3
    comment[65]: TRACKTOTAL=19
    comment[66]: TRACKC=19
    comment[67]: TOTALTRACKS=19
```

Well, at least I have all the tags I need and even more!

### Use panosse to verify, encode, normalize and clean the files

Now that the files are correctly tagged, the final steps are:

- [Verify](#verify) the integrity of the FLAC files
- [Encode](#encode) the FLAC files using the latest FLAC version
- [Normalize](#normalize) the FLAC files using ReplayGain
- [Clean](#clean) the FLAC files

The following commands can be used to achieve the final steps.

**Verify**

```sh
# Verify the integrity of the FLAC files
$ find ./files/processed -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 docker run --rm --volume "$(pwd)/files:/files" ghcr.io/ludelafo/panosse verify --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac" verified
[panosse::verify] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac" verified
[panosse::verify] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac" verified
[panosse::verify] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac" verified
```

</details>

**Encode**

```sh
# Encode the FLAC files
$ find ./files/processed -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 docker run --rm --volume "$(pwd)/files:/files" ghcr.io/ludelafo/panosse encode --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac" encoded
[panosse::encode] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac" encoded
[panosse::encode] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac" encoded
[panosse::encode] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac" FLAC_ARGUMENTS tag added
```

</details>

If you reexecute the command, the execution will be much faster as the files are
already encoded.

**Normalize**

> [!NOTE]
>
> ReplayGain is a technique to normalize the volume of audio files. It is used
> to prevent the volume from being too loud or too quiet when playing different
> audio files. ReplayGain can be calculated for the album and the track. This is
> where the file structure comes into play:
>
> - _2Pac - All Eyez on Me (1996)_ and _Eminem - The Eminem Show (2002)_ are
>   split into multiple CDs. But they are still considered as one album, so the
>   ReplayGain should be calculated on all the files in the subdirectories
> - _The Notorious B.I.G. - Ready To Die (1994)_ has a flat structure, so the
>   ReplayGain is calculated on all the files its only directory.

```sh
# Normalize the FLAC files by directory
$ find ./files/processed -mindepth 2 -maxdepth 2 -type d -print0 | sort -z | xargs -0 -n 1 bash -c '
    dir="$1"
    flac_files=()

    # Find all FLAC files in the current directory and store them in an array
    while IFS= read -r -d "" file; do
      flac_files+=("$file")
    done < <(find "$dir" -type f -name "*.flac" -print0)

    # Check if there are any FLAC files found
    if [ ${#flac_files[@]} -ne 0 ]; then
      # Pass the .flac files to the panosse normalize command
      docker run --rm --volume "$(pwd)/files:/files" ghcr.io/ludelafo/panosse normalize --verbose "${flac_files[@]}"
    fi
  ' {}
```

<details>
<summary>Expand the output</summary>

```text
[panosse::normalize] "[./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac]" normalized
[panosse::normalize] "[./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac ./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac]" METAFLAC_ARGUMENTS tag added
[panosse::normalize] "[./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac]" normalized
[panosse::normalize] "[./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac ./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac]" METAFLAC_ARGUMENTS tag added
[panosse::normalize] "[./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac]" normalized
[panosse::normalize] "[./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac ./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac]" METAFLAC_ARGUMENTS tag added
```

</details>

If you reexecute the command, the execution will be much faster as the files are
already normalized.

**Clean**

```sh
# Clean the FLAC files
$ find ./files/processed -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 docker run --rm --volume "$(pwd)/files:/files" ghcr.io/ludelafo/panosse clean --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/02 - All Bout U (feat. Snoop Dogg, Nate Dogg & Dru Down).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/04 - Got My Mind Made Up (feat. Daz Dillinger, Kurupt, Redman & Method Man).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/05 - How Do U Want It (feat. K‐Ci & JoJo).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Dogg).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/07 - No More Pain.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/08 - Heartz of Men.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/09 - Life Goes On.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4‐Tay).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/11 - Tradin' War Stories (feat. Outlawz, C‐Bo & Storm).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/12 - California Love (remix) (feat. Dr. Dre & Roger Troutman).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/01 - Can't C Me (feat. George Clinton).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/02 - Shorty Wanna Be a Thug.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/03 - Holla at Me.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/04 - Wonda Why They Call U Bytch.flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/05 - When We Ride (feat. Outlawz).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/06 - Thug Passion (feat. Jewell, Outlawz & Storm).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/07 - Picture Me Rollin' (feat. Danny Boy, Big Syke & CPO Boss Hogg).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/08 - Check Out Time (feat. Kurupt & Big Syke).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/11 - Run tha Streetz (feat. Michel'le, Napoleon & Storm).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/12 - Ain't Hard 2 Find (feat. E‐40, B‐Legit, C‐Bo & Richie Rich).flac" cleaned
[panosse::clean] "./files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 2/13 - Heaven Ain't Hard 2 Find.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/01 - Curtains Up (skit).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/02 - White America.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/04 - Cleanin' Out My Closet.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/05 - Square Dance.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/06 - The Kiss (skit).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/07 - Soldier.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/08 - Say Goodbye Hollywood.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/09 - Drips.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/10 - Without Me.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/11 - Paul Rosenberg (skit).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/12 - Sing for the Moment.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/13 - Superman.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/14 - Hailie's Song.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/15 - Steve Berman (skit).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/16 - When the Music Stops.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/17 - Say What You Say.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/18 - Till I Collapse.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/19 - My Dad's Gone Crazy.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/20 - Curtains Close.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/01 - Stimulate.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/02 - The Conspiracy Freestyle.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/03 - Bump Heads (feat. 50 Cent, Tony Yayo, and Lloyd Banks).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/04 - Jimmy, Brian and Mike.flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/05 - Freestyle (#1) (live at Tramps, New York, 1999).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/06 - Brain Damage (live at Tramps, New York, 1999).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/07 - Freestyle (#2) (live at Tramps, New York, 1999).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/08 - Just Don't Give a Fuck (live at Tramps, New York, 1999).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/09 - The Way I Am (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/10 - The Real Slim Shady (live at the Fuji Rock Festival, 2001) (feat. Proof).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/11 - Business (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/12 - Cleanin' Out My Closet (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/13 - Square Dance (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/14 - Without Me (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/15 - Sing for the Moment (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/16 - Superman (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/17 - Say What You Say (instrumental).flac" cleaned
[panosse::clean] "./files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 2/18 - Till I Collapse (instrumental).flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/01 - Intro.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/02 - Things Done Changed.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/04 - Machine Gun Funk.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/05 - Warning.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/06 - Ready to Die.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/07 - One More Chance.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/08 - Fuck Me (interlude).flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/09 - The What.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/10 - Juicy.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/11 - Everyday Struggle.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/12 - Me & My Bitch.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/13 - Big Poppa.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/14 - Respect.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/15 - Friend of Mine.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/16 - Unbelievable.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/17 - Suicidal Thoughts.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/18 - Who Shot Ya.flac" cleaned
[panosse::clean] "./files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/19 - Just Playing (Dreams).flac" cleaned
```

</details>

### Check the results

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/processed/2Pac/2Pac - All Eyez on Me (1996) [FLAC 16-44.1 CD] {314-524 204-2}/CD 1/01 - Ambitionz az a Ridah.flac"
```

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 14 bytes
  maximum framesize: 12388 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 12288024
  MD5 signature: e63d3931b934b23765bea0e754c27420
METADATA block #1
  type: 4 (VORBIS_COMMENT)
  is last: true
  length: 571
  vendor string: reference libFLAC 1.4.3 20230623
  comments: 17
    comment[0]: ALBUM=All Eyez on Me
    comment[1]: ALBUMARTIST=2Pac
    comment[2]: ARTIST=2Pac
    comment[3]: DISCNUMBER=1
    comment[4]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[5]: GENRE=Gangsta Rap
    comment[6]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[7]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[8]: REPLAYGAIN_ALBUM_GAIN=-8.98 dB
    comment[9]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[10]: REPLAYGAIN_TRACK_GAIN=-7.86 dB
    comment[11]: REPLAYGAIN_TRACK_PEAK=1.00000000
    comment[12]: TITLE=Ambitionz az a Ridah
    comment[13]: TRACKNUMBER=1
    comment[14]: TOTALDISCS=2
    comment[15]: TOTALTRACKS=14
    comment[16]: YEAR=1996
```

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/processed/Eminem/Eminem - The Eminem Show (Expanded Edition) (2002, Reissue 2023) [FLAC 16-44.1 CD] {B0035988-02}/CD 1/03 - Business.flac"
```

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 908 bytes
  maximum framesize: 13244 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 11100852
  MD5 signature: 956c9ded0e0dad45dead1f55e650266f
METADATA block #1
  type: 4 (VORBIS_COMMENT)
  is last: true
  length: 680
  vendor string: reference libFLAC 1.4.3 20230623
  comments: 18
    comment[0]: ALBUM=The Eminem Show (Expanded Edition)
    comment[1]: ALBUMARTIST=Eminem
    comment[2]: ARTIST=Eminem
    comment[3]: COMMENT=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[4]: DISCNUMBER=1
    comment[5]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[6]: GENRE=Hip Hop
    comment[7]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[8]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[9]: REPLAYGAIN_ALBUM_GAIN=-9.39 dB
    comment[10]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[11]: REPLAYGAIN_TRACK_GAIN=-9.57 dB
    comment[12]: REPLAYGAIN_TRACK_PEAK=0.99755859
    comment[13]: TITLE=Business
    comment[14]: TRACKNUMBER=3
    comment[15]: TOTALDISCS=2
    comment[16]: TOTALTRACKS=20
    comment[17]: YEAR=2023
```

```sh
# List all available blocks
$ docker run --rm \
  --entrypoint "" \
  --volume "$(pwd)/files:/files" \
  ghcr.io/ludelafo/panosse \
    metaflac --list "/files/processed/The Notorious B.I.G./The Notorious B.I.G. - Ready to Die (1994, Reissue 2004) [FLAC 16-44.1 CD] {0249-86280-1}/03 - Gimme the Loot.flac"
```

```text
METADATA block #0
  type: 0 (STREAMINFO)
  is last: false
  length: 34
  minimum blocksize: 4096 samples
  maximum blocksize: 4096 samples
  minimum framesize: 1731 bytes
  maximum framesize: 13377 bytes
  sample_rate: 44100 Hz
  channels: 2
  bits-per-sample: 16
  total samples: 13412868
  MD5 signature: 3335fc652d261db991a51aef346cc2a9
METADATA block #1
  type: 4 (VORBIS_COMMENT)
  is last: true
  length: 600
  vendor string: reference libFLAC 1.4.3 20230623
  comments: 18
    comment[0]: ALBUM=Ready to Die
    comment[1]: ALBUMARTIST=The Notorious B.I.G.
    comment[2]: ARTIST=The Notorious B.I.G.
    comment[3]: COMMENT=.
    comment[4]: DISCNUMBER=1
    comment[5]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[6]: GENRE=Rap
    comment[7]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[8]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[9]: REPLAYGAIN_ALBUM_GAIN=-8.41 dB
    comment[10]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[11]: REPLAYGAIN_TRACK_GAIN=-7.69 dB
    comment[12]: REPLAYGAIN_TRACK_PEAK=0.98852539
    comment[13]: TITLE=Gimme the Loot
    comment[14]: TRACKNUMBER=3
    comment[15]: TOTALDISCS=1
    comment[16]: TOTALTRACKS=19
    comment[17]: YEAR=2004
```

Thanks to panosse, all files have been cleaned from their unwanted blocks and
tags!

## Contributing

If you have interested in contributing to panosse, check the
[Contributing](https://github.com/ludelafo/panosse/blob/main/CONTRIBUTING.md)
guide.

Thank you in advance!

## License

panosse is licensed under the
[GNU Affero General Public License (GNU AGPL-3.0)](https://github.com/ludelafo/panosse/blob/main/LICENCE.md).

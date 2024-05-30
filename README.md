# panosse

[![Latest release](https://img.shields.io/github/v/release/ludelafo/panosse?include_prereleases)](https://github.com/ludelafo/panosse/releases)
[![License](https://img.shields.io/github/license/ludelafo/panosse)](https://github.com/ludelafo/panosse)
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

As already mentioned, panosse is only a wrapper around flac and metaflac. It does
not provide much more functionality. It was developed to automate and set
sane defaults for my music library maintenance.

panosse tries to stay close to the UNIX philosophy of doing one thing and doing
it well. For example, panosse only proccesses one file at a time (except for
normalization), so you can easily parallelize the process using `find` and
`xargs` or similar tools.

For usage and configuration, see the [Usage](#usage) section and the
[Configuration](#configuration) section. Check the
[Concrete example](#concrete-example) for a real-world example.

## Usage

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
  -h, --help                           help for panosse
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
  -v, --version                        version for panosse

Use "panosse [command] --help" for more information about a command.
```

### Verify

```text
$ ./panosse verify --help
Check the integrity of the FLAC files.

It calls metaflac to verify the FLAC files.

Usage:
  panosse verify <file> [flags]

Examples:
  # Verify a single FLAC file
  $ panosse verify file.flac

  # Verify all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 -P$(nproc) panosse verify

Flags:
  -h, --help                       help for verify
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
  # Encode a single FLAC file
  $ panosse encode file.flac

  # Encode all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 -P$(nproc) panosse encode

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

It calls metaflac to add the ReplayGain tags to the FLAC files.

Usage:
  panosse normalize <file 1> [<file 2>]... [flags]

Examples:
  # Normalize some FLAC files
  $ panosse normalize file1.flac file2.flac

  # Normalize all FLAC files in all directories in parallel for a depth of 1
  # This allows to consider the nested directories as one album for the normalization
  $ find . -mindepth 1 -maxdepth 1 -type d -print0 | sort -z | while IFS= read -r -d '' dir; do
    mapfile -d '' -t flac_files < <(find "$dir" -type f -name "*.flac" -print0)
  
    if [ ${#flac_files[@]} -ne 0 ]; then
      panosse normalize --verbose "${flac_files[@]}"
    fi
  done

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
  # Clean a single FLAC file
  $ panosse clean file.flac

  # Clean all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 -P$(nproc) panosse clean

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

You can display the current configuration with `panosse config`.

### Flags

You can check the available flags for each command with `panosse help <command>`
or `panosse <command> --help`.

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

1. `--config-file` flag, allowing you to specify the configuration file
2. `config.yaml` in the current directory
3. `$HOME/.panosse/config.yaml` on Linux and `%USER%\.panosse\config.yaml` on
   Windows

If no configuration file is found, the default values are used from the
[flags](#flags) section.

#### Examples

```yaml
# custon config.yaml, config.yaml in the current directory or ~/.panosse/config.yaml
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

Once panosse is built, you can run it with the following command:

```sh
# Run panosse
./panosse
```

## What does panosse mean?

panosse (`/pa.nɔs/`) is a Swiss-French word meaning mop. The idea is that a mop
cleans a floor, panosse cleans FLAC files.

## Concrete example

This section provides a real-world example of how to use panosse to clean,
encode, normalize, and verify a FLAC music library.

### Structure and file contents

Let's say you have a music library with the following structure:

```text
.
├── 2Pac - All Eyez on Me (1996)
│   ├── CD 1
│   │   ├── 01 - Ambitionz Az a Ridah.flac
│   │   ├── 02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac
│   │   ├── 03 - Skandalouz (feat. Nate Dogg).flac
│   │   ├── 04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac
│   │   ├── 05 - How Do U Want It (feat. KC & JoJo).flac
│   │   ├── 06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac
│   │   ├── 07 - No More Pain.flac
│   │   ├── 08 - Heartz of Men.flac
│   │   ├── 09 - Life Goes On.flac
│   │   ├── 10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac
│   │   ├── 11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac
│   │   ├── 12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac
│   │   ├── 13 - I Ain't Mad at Cha (feat. Danny Boy).flac
│   │   ├── 14 - What'z Ya Phone # (feat. Danny Boy).flac
│   │   └── folder.jpg
│   └── CD 2
│       ├── 01 - Can't C Me (feat. George Clinton).flac
│       ├── 02 - Shorty Wanna Be a Thug.flac
│       ├── 03 - Holla at Me.flac
│       ├── 04 - Wonda Why They Call U Bitch.flac
│       ├── 05 - When We Ride (feat. Outlaw Immortalz).flac
│       ├── 06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac
│       ├── 07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac
│       ├── 08 - Check Out Time (feat. Big Syke & Kurupt).flac
│       ├── 09 - Ratha Be Ya Nigga (feat. Richie Rich).flac
│       ├── 10 - All Eyez on Me (feat. Big Syke).flac
│       ├── 11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac
│       ├── 12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac
│       ├── 13 - Heaven Ain't Hard 2 Find.flac
│       └── folder.jpg
├── Eminem - The Eminem Show (2002)
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
│   ├── 01-14 - Hailie's Song.flac
│   ├── 01-15 - Steve Berman (skit).flac
│   ├── 01-16 - When the Music Stops (feat. D12).flac
│   ├── 01-17 - Say What You Say (feat. Dr. Dre).flac
│   ├── 01-18 - 'Till I Collapse (feat. Nate Dogg).flac
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
│   ├── 02-18 - 'Till I Collapse (Instrumental).flac
│   └── folder.jpg
└── The Notorious B.I.G. - Ready To Die (1994)
    ├── 01 - Intro.flac
    ├── 02 - Things Done Changed.flac
    ├── 03 - Gimme The Loot.flac
    ├── 04 - Machine Gun Funk.flac
    ├── 05 - Warning.flac
    ├── 06 - Ready To Die.flac
    ├── 07 - One More Chance.flac
    ├── 08 - Fuck Me (Interlude).flac
    ├── 09 - The What.flac
    ├── 10 - Juicy.flac
    ├── 11 - Everyday Struggles.flac
    ├── 12 - Me & My Bitch.flac
    ├── 13 - Big Poppa.flac
    ├── 14 - Respect.flac
    ├── 15 - Friend Of Mine.flac
    ├── 16 - Unbelievable.flac
    ├── 17 - Suicidal Thoughts.flac
    ├── 18 - Who Shot Ya.flac
    ├── 19 - Just Playing (Dreams).flac
    └── folder.jpg
```

The first thing to note is the difference in structure:

1. _2Pac - All Eyez on Me (1996)_ is split into two CDs, making it a nested
   structure
2. _Eminem - The Eminem Show (2002)_ has two CDs in a flat structure
3. _The Notorious B.I.G. - Ready To Die (1994)_ has a flat structure

Let's have a look at the files:

```sh
# List all available blocks
$ metaflac --list "2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac"
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
$ metaflac --list --block-type=VORBIS_COMMENT "2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac"
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
$ metaflac --list "Eminem - The Eminem Show (2002)/01-03 - Business.flac"
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
$ metaflac --list --block-type=VORBIS_COMMENT "Eminem - The Eminem Show (2002)/01-03 - Business.flac"
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
$ metaflac --list "The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac"
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
$ metaflac --list --block-type=VORBIS_COMMENT "The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac"
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

ReplayGain is a technique to normalize the volume of audio files. It is used to
prevent the volume from being too loud or too quiet when playing different audio
files. ReplayGain can be calculated for the album and the track. This is where
the file structure comes into play:

- _2Pac - All Eyez on Me (1996)_ is split into two CDs, so the ReplayGain should
  be calculated on all the files in the subdirectories
- The other albums (_Eminem - The Eminem Show (2002)_ and _The Notorious
  B.I.G. - Ready To Die (1994)_) have a flat structure, so the ReplayGain is
  calculated on all the files in each directory

### Final steps

With the considerations in mind, the final steps are:

- [Verify](#verify) the integrity of the FLAC files
- [Encode](#encode) the FLAC files using the latest FLAC version
- [Normalize](#normalize) the FLAC files using ReplayGain
- [Clean](#clean) the FLAC files

The following commands can be used to achieve the final steps:

```sh
# Verify the integrity of the FLAC files
find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 ./panosse verify --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac" verified
[panosse::verify] "./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-02 - White America.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-03 - Business.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac" verified
[panosse::verify] "./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac" verified
[panosse::verify] "./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac" verified
```

</details>

```sh
# Encode the FLAC files
find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 ./panosse encode --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac" encoded
[panosse::encode] "./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-02 - White America.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-02 - White America.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-03 - Business.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-03 - Business.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac" encoded
[panosse::encode] "./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac" FLAC_ARGUMENTS tag added
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac" encoded
[panosse::encode] "./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac" FLAC_ARGUMENTS tag added
```

</details>

If you reexecute the command, the execution will be much faster as the files are
already encoded.

```sh
# Normalize the FLAC files by directory
find . -mindepth 1 -maxdepth 1 -type d -print0 | sort -z | while IFS= read -r -d '' dir; do
  # For each directory, find all FLAC files and store them in an array
  mapfile -d '' -t flac_files < <(find "$dir" -type f -name "*.flac" -print0)
  
  # Check if the array is not empty
  if [ ${#flac_files[@]} -ne 0 ]; then
    # Run the panosse normalize command with the found .flac files
    ./panosse normalize --verbose "${flac_files[@]}"
  fi
done
```

<details>
<summary>Expand the output</summary>

```text
[panosse::normalize] "[./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac ./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac ./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac ./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac ./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac ./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac ./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac ./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac ./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac ./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac ./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac ./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac ./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac ./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac ./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac ./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac ./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac ./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac ./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac ./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac ./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac]" normalized
[panosse::normalize] "[./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac ./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac ./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac ./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac ./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac ./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac ./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac ./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac ./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac ./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac ./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac ./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac ./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac ./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac ./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac ./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac ./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac ./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac ./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac ./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac ./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac ./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac ./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac]" METAFLAC_ARGUMENTS tag added
[panosse::normalize] "[./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac ./Eminem - The Eminem Show (2002)/01-02 - White America.flac ./Eminem - The Eminem Show (2002)/01-03 - Business.flac ./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac ./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac ./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac ./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac ./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac ./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac ./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac ./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac ./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac ./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac ./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac ./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac ./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac ./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac ./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac ./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac ./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac ./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac ./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac ./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac ./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac ./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac ./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac ./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac]" normalized
[panosse::normalize] "[./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac ./Eminem - The Eminem Show (2002)/01-02 - White America.flac ./Eminem - The Eminem Show (2002)/01-03 - Business.flac ./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac ./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac ./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac ./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac ./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac ./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac ./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac ./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac ./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac ./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac ./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac ./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac ./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac ./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac ./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac ./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac ./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac ./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac ./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac ./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac ./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac ./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac ./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac ./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac ./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac ./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac]" METAFLAC_ARGUMENTS tag added
[panosse::normalize] "[./Geto Boys - The Resurrection (1996)/01 - Ghetto Prisoner.flac ./Geto Boys - The Resurrection (1996)/02 - Still.flac ./Geto Boys - The Resurrection (1996)/03 - The World Is A Ghetto (feat. Flaj).flac ./Geto Boys - The Resurrection (1996)/04 - Open Minded (feat. DMG).flac ./Geto Boys - The Resurrection (1996)/05 - Killer For Scratch.flac ./Geto Boys - The Resurrection (1996)/06 - Hold It Down (feat. Facemob).flac ./Geto Boys - The Resurrection (1996)/07 - Blind Leading The Blind (feat. Menace Clan).flac ./Geto Boys - The Resurrection (1996)/08 - First Light Of The Day.flac ./Geto Boys - The Resurrection (1996)/09 - Time Taker.flac ./Geto Boys - The Resurrection (1996)/10 - Geto Boys And Girls.flac ./Geto Boys - The Resurrection (1996)/11 - Geto Fantasy.flac ./Geto Boys - The Resurrection (1996)/12 - I Just Wanna Die.flac ./Geto Boys - The Resurrection (1996)/13 - Niggas And Flies.flac ./Geto Boys - The Resurrection (1996)/14 - A Visit With Larry Hoover.flac ./Geto Boys - The Resurrection (1996)/15 - Point Of No Return.flac]" normalized
[panosse::normalize] "[./Geto Boys - The Resurrection (1996)/01 - Ghetto Prisoner.flac ./Geto Boys - The Resurrection (1996)/02 - Still.flac ./Geto Boys - The Resurrection (1996)/03 - The World Is A Ghetto (feat. Flaj).flac ./Geto Boys - The Resurrection (1996)/04 - Open Minded (feat. DMG).flac ./Geto Boys - The Resurrection (1996)/05 - Killer For Scratch.flac ./Geto Boys - The Resurrection (1996)/06 - Hold It Down (feat. Facemob).flac ./Geto Boys - The Resurrection (1996)/07 - Blind Leading The Blind (feat. Menace Clan).flac ./Geto Boys - The Resurrection (1996)/08 - First Light Of The Day.flac ./Geto Boys - The Resurrection (1996)/09 - Time Taker.flac ./Geto Boys - The Resurrection (1996)/10 - Geto Boys And Girls.flac ./Geto Boys - The Resurrection (1996)/11 - Geto Fantasy.flac ./Geto Boys - The Resurrection (1996)/12 - I Just Wanna Die.flac ./Geto Boys - The Resurrection (1996)/13 - Niggas And Flies.flac ./Geto Boys - The Resurrection (1996)/14 - A Visit With Larry Hoover.flac ./Geto Boys - The Resurrection (1996)/15 - Point Of No Return.flac]" METAFLAC_ARGUMENTS tag added
[panosse::normalize] "[./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac ./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac ./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac ./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac ./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac ./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac ./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac ./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac ./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac ./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac ./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac ./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac ./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac ./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac ./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac ./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac ./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac ./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac ./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac]" normalized
[panosse::normalize] "[./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac ./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac ./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac ./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac ./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac ./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac ./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac ./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac ./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac ./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac ./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac ./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac ./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac ./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac ./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac ./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac ./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac ./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac ./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac]" METAFLAC_ARGUMENTS tag added
```

</details>

If you reexecute the command, the execution will be much faster as the files are
already normalized.

```sh
# Clean the FLAC files
find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 ./panosse clean --verbose
```

<details>
<summary>Expand the output</summary>

```text
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/02 - All Bout U (feat. Snoop Doggy Dogg, Nate Dogg & Dru Down).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/03 - Skandalouz (feat. Nate Dogg).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/04 - Got My Mind Made Up (feat. Dat Niggaz Daz, Kurupt, Redman & Method Man).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/05 - How Do U Want It (feat. KC & JoJo).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/06 - 2 of Amerikaz Most Wanted (feat. Snoop Doggy Dogg).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/07 - No More Pain.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/08 - Heartz of Men.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/09 - Life Goes On.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/10 - Only God Can Judge Me (feat. Rappin' 4-Tay).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/11 - Tradin War Stories (feat. Dramacydal, C-Bo & Storm).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/12 - California Love (Remix) (feat. Dr. Dre & Roger Troutman).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/13 - I Ain't Mad at Cha (feat. Danny Boy).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 1/14 - What'z Ya Phone # (feat. Danny Boy).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/01 - Can't C Me (feat. George Clinton).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/02 - Shorty Wanna Be a Thug.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/03 - Holla at Me.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/04 - Wonda Why They Call U Bitch.flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/05 - When We Ride (feat. Outlaw Immortalz).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/06 - Thug Passion (feat. Dramarydal, Jewell & Storm).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/07 - Picture Me Rollin' (feat. Big Syke, CPO, Danny Boy).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/08 - Check Out Time (feat. Big Syke & Kurupt).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/09 - Ratha Be Ya Nigga (feat. Richie Rich).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/10 - All Eyez on Me (feat. Big Syke).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/11 - Run Tha Streetz (feat. Michel'le, Mutah & Storm).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/12 - Ain't Hard 2 Find (feat. B-Legit, C-Bo, E-40 & Richie Rich).flac" cleaned
[panosse::clean] file "./2Pac - All Eyez on Me (1996)/CD 2/13 - Heaven Ain't Hard 2 Find.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-01 - Curtains Up (skit).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-02 - White America.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-03 - Business.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-04 - Cleanin Out My Closet.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-05 - Square Dance.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-06 - The Kiss (skit).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-07 - Soldier.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-08 - Say Goodbye Hollywood.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-09 - Drips (feat. Obie Trice).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-10 - Without Me.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-11 - Paul Rosenberg (skit).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-12 - Sing for the Moment.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-13 - Superman (feat. Dina Rae).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-14 - Hailie's Song.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-15 - Steve Berman (skit).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-16 - When the Music Stops (feat. D12).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-17 - Say What You Say (feat. Dr. Dre).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-18 - 'Till I Collapse (feat. Nate Dogg).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-19 - My Dad’s Gone Crazy (feat. Hailie Jade).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/01-20 - Curtains Close (skit).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-01 - Stimulate.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-02 - The Conspiracy Freestyle (DJ Green Lantern Version).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-03 - Bump Heads (DJ Green Lantern Version).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-04 - Jimmy, Brian and Mike.flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-05 - Freestyle #1 (Live from Tramps, New York , 1999).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-06 - Brain Damage (Live from Tramps, New York , 1999).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-07 - Freestyle #2 (Live from Tramps, New York , 1999).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-08 - Just Don't Give a Fuck (Live from Tramps, New York , 1999).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-09 - The Way I Am (Live from Fuji Rock Festival, Japan , 2001).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-10 - The Real Slim Shady (Live from Fuji Rock Festival, Japan , 2001).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-11 - Business (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-12 - Cleanin' Out My Closet (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-13 - Square Dance (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-14 - Without Me (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-15 - Sing for the Moment (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-16 - Superman (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-17 - Say What You Say (Instrumental).flac" cleaned
[panosse::clean] file "./Eminem - The Eminem Show (2002)/02-18 - 'Till I Collapse (Instrumental).flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/01 - Intro.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/02 - Things Done Changed.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/04 - Machine Gun Funk.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/05 - Warning.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/06 - Ready To Die.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/07 - One More Chance.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/08 - Fuck Me (Interlude).flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/09 - The What.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/10 - Juicy.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/11 - Everyday Struggles.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/12 - Me & My Bitch.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/13 - Big Poppa.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/14 - Respect.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/15 - Friend Of Mine.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/16 - Unbelievable.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/17 - Suicidal Thoughts.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/18 - Who Shot Ya.flac" cleaned
[panosse::clean] file "./The Notorious B.I.G. - Ready To Die (1994)/19 - Just Playing (Dreams).flac" cleaned
```

</details>

### Check the results

Let's have a look at the files. Cleaner now!

In this real world example, the files
were not pre-processed using tools such as [beets](https://beets.io/) or
[MusicBrainz Picard](https://picard.musicbrainz.org/), thus missing some of the
tags I would like to have. However, the files are now clean and ready for further
processing.

```sh
# List all available blocks
$ metaflac --list "2Pac - All Eyez on Me (1996)/CD 1/01 - Ambitionz Az a Ridah.flac"
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
  length: 521
  vendor string: reference libFLAC 1.4.2 20221022
  comments: 14
    comment[0]: ALBUM=All Eyez on Me
    comment[1]: ALBUMARTIST=2Pac
    comment[2]: ARTIST=2Pac
    comment[3]: DISCNUMBER=1
    comment[4]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[5]: GENRE=Hip-Hop
    comment[6]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[7]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[8]: REPLAYGAIN_ALBUM_GAIN=-9.07 dB
    comment[9]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[10]: REPLAYGAIN_TRACK_GAIN=-7.86 dB
    comment[11]: REPLAYGAIN_TRACK_PEAK=1.00000000
    comment[12]: TITLE=Ambitionz Az a Ridah
    comment[13]: TRACKNUMBER=01
```

```sh
# List all available blocks
$ metaflac --list "Eminem - The Eminem Show (2002)/01-03 - Business.flac"
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
  length: 615
  vendor string: reference libFLAC 1.4.2 20221022
  comments: 15
    comment[0]: ALBUM=The Eminem Show
    comment[1]: ALBUMARTIST=Eminem
    comment[2]: ARTIST=Eminem
    comment[3]: COMMENT=− 2023 - Shady Records / Aftermath Records / Interscope Records / Expanded Edition / CD
    comment[4]: DISCNUMBER=1
    comment[5]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[6]: GENRE=Hip-Hop
    comment[7]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[8]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[9]: REPLAYGAIN_ALBUM_GAIN=-9.07 dB
    comment[10]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[11]: REPLAYGAIN_TRACK_GAIN=-9.57 dB
    comment[12]: REPLAYGAIN_TRACK_PEAK=0.99755859
    comment[13]: TITLE=Business
    comment[14]: TRACKNUMBER=03
```

```sh
# List all available blocks
$ metaflac --list "The Notorious B.I.G. - Ready To Die (1994)/03 - Gimme The Loot.flac"
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
  length: 521
  vendor string: reference libFLAC 1.4.2 20221022
  comments: 13
    comment[0]: ALBUM=Ready To Die (The Remaster)
    comment[1]: ARTIST=The Notorious B.I.G.
    comment[2]: COMMENT=.
    comment[3]: FLAC_ARGUMENTS=--compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent
    comment[4]: GENRE=Hip-Hop
    comment[5]: METAFLAC_ARGUMENTS=--add-replay-gain
    comment[6]: REPLAYGAIN_REFERENCE_LOUDNESS=89.0 dB
    comment[7]: REPLAYGAIN_ALBUM_GAIN=-9.07 dB
    comment[8]: REPLAYGAIN_ALBUM_PEAK=1.00000000
    comment[9]: REPLAYGAIN_TRACK_GAIN=-7.69 dB
    comment[10]: REPLAYGAIN_TRACK_PEAK=0.98852539
    comment[11]: TITLE=Gimme The Loot
    comment[12]: TRACKNUMBER=03
```

## License

panosse is licensed under the
[GNU Affero General Public License (GNU AGPL-3.0)](./COPYING).

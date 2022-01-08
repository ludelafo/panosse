`audio²`
========

A collection of Node.js scripts to process audio files

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![License](https://img.shields.io/npm/l/audiotwo.svg)](https://github.com/ludelafo/audiotwo/blob/main/package.json)

<!-- toc -->
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->
```sh
$ npm install
$ ./bin/run COMMAND
running command...
$ ./bin/run (-v|--version|version)
audiotwo/0.0.0 darwin-x64 node-v14.17.0
$ audiotwo --help [COMMAND]
USAGE
  $ audiotwo COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`audiotwo flactwoflac INPUT`](#audiotwo-flactwoflac-input)
* [`audiotwo help [COMMAND]`](#audiotwo-help-command)

## `audiotwo flactwoflac INPUT`

Process FLAC files to FLAC files

```
USAGE
  $ audiotwo flactwoflac INPUT

ARGUMENTS
  INPUT  input folder

OPTIONS
  -h, --help                               show CLI help
  --addReplayGain                          add ReplayGain to FLAC files
  --addReplayGainIfTagMismatches           add ReplayGain if tag mismatches
  --addReplayGainIfTagsAreNotPresent       add ReplayGain if tags are not present
  --cleanTags                              clean tags from FLAC files
  --compress                               compress FLAC files
  --compressIfTagIsNotPresent              compress if tag is not present
  --compressIfTagMismatches                compress if tag mismatches
  --compressIfVersionMismatches            compress if version mismatches
  --dryRun                                 run without actually doing anything
  --encoderSettings=encoderSettings        [default: flac --compression-level-8 --delete-input-file --no-padding --force --verify --warnings-as-errors --silent] encoder settings
  --encoderTag=encoderTag                  [default: ENCODER_SETTINGS] encoder tag name
  --fillMissingTags                        Fill missings tags with a "[No <tag>]" in file
  --force                                  force command
  --removeApplicationBlocks                remove application blocks
  --removeCuesheetBlocks                   remove cuesheet blocks
  --removeEmbeddedBlocks                   remove embedded blocks from FLAC files
  --removePaddingBlocks                    remove padding blocks
  --removePictureBlocks                    remove picture blocks
  --removeSeektableBlocks                  remove seek table blocks
  --replayGainSettings=replayGainSettings  [default: metaflac --add-replay-gain] ReplayGain settings
  --replayGainTag=replayGainTag            [default: REPLAYGAIN_SETTINGS] ReplayGain tag name
  --replayGainTags=replayGainTags          [default: REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK]
  --saveEncoderSettingsInTag               add encoder tag with encoder settings
  --saveReplayGainSettingsInTag            add ReplayGain tag with ReplayGain settings
  --tagsToKeep=tagsToKeep                  [default: ALBUM,ALBUMARTIST,ARTIST,COMMENT,DISCNUMBER,ENCODER_SETTINGS,GENRE,REPLAYGAIN_SETTINGS,REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,TITLE,TRACKNUMBER,TOTALDISCS,TOTALTRACKS,YEAR]
  --threadCount=threadCount                [default: 12] number of threads to use for processing
```

_See code: [src/commands/flactwoflac/index.ts](https://github.com/ludelafo/audiotwo/blob/v0.0.0/src/commands/flactwoflac/index.ts)_

## `audiotwo help [COMMAND]`

Display help for audiotwo

```
USAGE
  $ ./bin/run help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all  see all commands in CLI
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v3.2.6/src/commands/help.ts)_
<!-- commandsstop -->

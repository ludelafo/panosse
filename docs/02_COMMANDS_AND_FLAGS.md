<!--
Generate this file with the following command:

```sh
echo "# Commands and flags" > COMMANDS_AND_FLAGS.md
echo "" >> COMMANDS_AND_FLAGS.md

echo "> [!TIP]" >> COMMANDS_AND_FLAGS.md
echo ">" >> COMMANDS_AND_FLAGS.md
echo "> The recommended order to execute panosse is:" >> COMMANDS_AND_FLAGS.md
echo ">" >> COMMANDS_AND_FLAGS.md
echo "> 1. \`verify\`" >> COMMANDS_AND_FLAGS.md
echo "> 2. \`encode\`" >> COMMANDS_AND_FLAGS.md
echo "> 3. \`normalize\`" >> COMMANDS_AND_FLAGS.md
echo "> 4. \`clean\`" >> COMMANDS_AND_FLAGS.md
echo ">" >> COMMANDS_AND_FLAGS.md
echo "> This will ensure all the files are correct, encoded with the latest FLAC" >> COMMANDS_AND_FLAGS.md
echo "> version, normalized with ReplayGain, and cleaned from unnecessary blocks and" >> COMMANDS_AND_FLAGS.md
echo "> tags after all other operations." >> COMMANDS_AND_FLAGS.md
echo "" >> COMMANDS_AND_FLAGS.md

echo "Every panosse's commands have a \`help\` command to describe the command's usage." >> COMMANDS_AND_FLAGS.md
echo "" >> COMMANDS_AND_FLAGS.md
echo "You can use \`panosse [command] --help\` or \`panosse help [command]\` to display" >> COMMANDS_AND_FLAGS.md
echo "the help." >> COMMANDS_AND_FLAGS.md
echo "" >> COMMANDS_AND_FLAGS.md

echo "\`\`\`text" >> COMMANDS_AND_FLAGS.md
echo "$ ./panosse --help" >> COMMANDS_AND_FLAGS.md
./panosse --help >> COMMANDS_AND_FLAGS.md
echo "\`\`\`" >> COMMANDS_AND_FLAGS.md

for command in clean encode normalize verify; do
  echo "" >> COMMANDS_AND_FLAGS.md
  echo "## ${command^}" >> COMMANDS_AND_FLAGS.md
  echo "" >> COMMANDS_AND_FLAGS.md
  echo "\`\`\`text" >> COMMANDS_AND_FLAGS.md
  echo "$ ./panosse ${command} --help" >> COMMANDS_AND_FLAGS.md
  ./panosse ${command} --help >> COMMANDS_AND_FLAGS.md
  echo "\`\`\`" >> COMMANDS_AND_FLAGS.md
done
```
-->

# Commands and flags

> [!TIP]
>
> The recommended order to execute panosse is:
>
> 1. `verify`
> 2. `encode`
> 3. `normalize`
> 4. `clean`
>
> This will ensure all the files are correct, encoded with the latest FLAC
> version, normalized with ReplayGain, and cleaned from unnecessary blocks and
> tags after all other operations.

Every panosse's commands have a `help` command to describe the command's usage.

You can use `panosse [command] --help` or `panosse help [command]` to display
the help.

```text
$ ./panosse --help
panosse is a CLI tool to clean, encode, normalize, and verify your FLAC music library.

panosse is merely a wrapper around flac and metaflac and uses Cobra and Viper under the hood.

panosse is licensed under the GNU Affero General Public License (GNU AGPL-3.0).

For more information, see https://github.com/ludelafo/panosse.

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

## Clean

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
  -h, --help                      help for clean
  -t, --tags-to-keep strings      tags to keep in the file (default [ALBUM,ALBUMARTIST,ARTIST,COMMENT,COMPOSER,DISCNUMBER,FLAC_ARGUMENTS,GENRE,METAFLAC_ARGUMENTS,PERFORMER,REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,TITLE,TOTALDISCS,TOTALTRACKS,TRACKNUMBER,YEAR])

Global Flags:
  -C, --config-file string             config file to use (optional - will use "config.yaml" or "~/.panosse/config.yaml" if available)
  -D, --dry-run                        perform a trial run with no changes made
  -F, --flac-command-path string       path to the flac command (checks in $PATH as well) (default "flac")
  -X, --force                          force processing even if no processing is needed
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
```

## Encode

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
  -a, --encode-arguments strings                   arguments passed to flac to encode the file (default [--compression-level-8,--delete-input-file,--force,--no-padding,--silent,--verify,--warnings-as-errors])
      --encode-if-encode-argument-tags-mismatch    encode if encode argument tags mismatch (missing or different) (default true)
      --encode-if-flac-versions-mismatch           encode if flac versions mismatch between host's flac version and file's flac version (default true)
  -h, --help                                       help for encode
      --save-encode-arguments-in-tag               save encode arguments in tag (default true)
      --save-encode-arguments-in-tag-name string   encode arguments tag name (default "FLAC_ARGUMENTS")

Global Flags:
  -C, --config-file string             config file to use (optional - will use "config.yaml" or "~/.panosse/config.yaml" if available)
  -D, --dry-run                        perform a trial run with no changes made
  -F, --flac-command-path string       path to the flac command (checks in $PATH as well) (default "flac")
  -X, --force                          force processing even if no processing is needed
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
```

## Normalize

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
  -h, --help                                            help for normalize
  -a, --normalize-arguments strings                     arguments passed to flac to normalize the files (default [--add-replay-gain])
      --normalize-if-any-replaygain-tags-are-missing    normalize if any ReplayGain tags are missing (default true)
      --normalize-if-normalize-argument-tags-mismatch   normalize if normalize arguments tags mismatch (missing or different) (default true)
  -t, --replaygain-tags strings                         ReplayGain tags (default [REPLAYGAIN_REFERENCE_LOUDNESS,REPLAYGAIN_TRACK_GAIN,REPLAYGAIN_TRACK_PEAK,REPLAYGAIN_ALBUM_GAIN,REPLAYGAIN_ALBUM_PEAK])
      --save-normalize-arguments-in-tag                 save normalize arguments in tag (default true)
      --save-normalize-arguments-in-tag-name string     normalize arguments tag name (default "METAFLAC_ARGUMENTS")

Global Flags:
  -C, --config-file string             config file to use (optional - will use "config.yaml" or "~/.panosse/config.yaml" if available)
  -D, --dry-run                        perform a trial run with no changes made
  -F, --flac-command-path string       path to the flac command (checks in $PATH as well) (default "flac")
  -X, --force                          force processing even if no processing is needed
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
```

## Verify

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
  -h, --help                       help for verify
  -a, --verify-arguments strings   arguments passed to flac to verify the files (default [--test,--silent])

Global Flags:
  -C, --config-file string             config file to use (optional - will use "config.yaml" or "~/.panosse/config.yaml" if available)
  -D, --dry-run                        perform a trial run with no changes made
  -F, --flac-command-path string       path to the flac command (checks in $PATH as well) (default "flac")
  -X, --force                          force processing even if no processing is needed
  -M, --metaflac-command-path string   path to the metaflac command (checks in $PATH as well) (default "metaflac")
  -V, --verbose                        enable verbose output
```

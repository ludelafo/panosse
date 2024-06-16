/*
Copyright © 2023 Ludovic Delafontaine <@ludelafo>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"panosse/utils"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
type NormalizeCmdArgs struct {
	NormalizeArguments                       []string `mapstructure:"normalize-arguments"`
	NormalizeIfNormalizeArgumentTagsMismatch bool     `mapstructure:"normalize-if-normalize-argument-tags-mismatch"`
	NormalizeIfAnyReplayGainTagsAreMissing   bool     `mapstructure:"normalize-if-any-replaygain-tags-are-missing"`
	ReplaygainTags                           []string `mapstructure:"replaygain-tags"`
	SaveNormalizeArgumentsInTag              bool     `mapstructure:"save-normalize-arguments-in-tag"`
	SaveNormalizeArgumentsInTagName          string   `mapstructure:"save-normalize-arguments-in-tag-name"`
}

var normalizeCmdArgs NormalizeCmdArgs

var normalizeCmd = &cobra.Command{
	Use:   "normalize <file 1> [<file 2>]...",
	Short: "Normalize FLAC files with ReplayGain",
	Long: `Normalize FLAC files with ReplayGain.

It calls metaflac to calculate and add the ReplayGain tags to the FLAC files.`,
	Example: `  ## Normalize some FLAC files
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
  ' {}`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::normalize] ")

		// Get command line arguments from Viper
		viper.Unmarshal(&normalizeCmdArgs)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFiles := args

		// Command action
		normalizeArgumentsAsString := strings.Join(normalizeCmdArgs.NormalizeArguments, " ")

		needToNormalize := false

		for _, flacFile := range flacFiles {
			normalizeArgumentsTagContent, err := utils.GetTag(
				rootCmdArgs.MetaflacCommandPath,
				normalizeCmdArgs.SaveNormalizeArgumentsInTagName,
				flacFile,
			)

			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					resultCode := exitError.ExitCode()

					log.Fatalf(
						"ERROR - cannot get tag from file \"%s\" (exit code %d)",
						flacFile,
						resultCode,
					)
				}

				os.Exit(1)
			}

			if normalizeCmdArgs.NormalizeIfNormalizeArgumentTagsMismatch && normalizeArgumentsTagContent == "" {
				needToNormalize = true
			}

			if normalizeCmdArgs.NormalizeIfAnyReplayGainTagsAreMissing {
				for _, replaygainTag := range normalizeCmdArgs.ReplaygainTags {
					replaygainTagContent, _ := utils.GetTag(
						rootCmdArgs.MetaflacCommandPath,
						replaygainTag,
						flacFile,
					)

					if replaygainTagContent == "" {
						needToNormalize = true
						break
					}
				}
			}
		}

		if needToNormalize || rootCmdArgs.Force {
			if !rootCmdArgs.DryRun {
				err := utils.Normalize(
					rootCmdArgs.MetaflacCommandPath,
					normalizeCmdArgs.NormalizeArguments,
					flacFiles,
				)

				if err != nil {
					if exitError, ok := err.(*exec.ExitError); ok {
						resultCode := exitError.ExitCode()

						log.Fatalf(
							"ERROR - cannot normalize files \"%s\" (exit code %d)",
							flacFiles,
							resultCode,
						)
					}

					os.Exit(1)
				}
			}
		}

		if rootCmdArgs.Verbose {
			log.Printf("\"%s\" normalized\n", flacFiles)
		}

		for _, flacFile := range flacFiles {
			if normalizeCmdArgs.SaveNormalizeArgumentsInTag {
				if !rootCmdArgs.DryRun {
					utils.RemoveTag(
						rootCmdArgs.MetaflacCommandPath,
						normalizeCmdArgs.SaveNormalizeArgumentsInTagName,
						flacFile,
					)
					utils.SetTag(
						rootCmdArgs.MetaflacCommandPath,
						normalizeCmdArgs.SaveNormalizeArgumentsInTagName,
						normalizeArgumentsAsString,
						flacFile,
					)
				}
			}
		}

		if rootCmdArgs.Verbose {
			log.Printf(
				"\"%s\" %s tag added\n",
				flacFiles,
				normalizeCmdArgs.SaveNormalizeArgumentsInTagName,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	normalizeCmd.PersistentFlags().StringSliceVarP(
		&normalizeCmdArgs.NormalizeArguments,
		"normalize-arguments", "a", []string{
			"--add-replay-gain",
		}, "arguments passed to flac to normalize the files",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeCmdArgs.NormalizeIfNormalizeArgumentTagsMismatch,
		"normalize-if-normalize-argument-tags-mismatch",
		true,
		"normalize if normalize arguments tags mismatch (missing or different)",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeCmdArgs.NormalizeIfAnyReplayGainTagsAreMissing,
		"normalize-if-any-replaygain-tags-are-missing",
		true,
		"normalize if any ReplayGain tags are missing",
	)
	normalizeCmd.PersistentFlags().StringSliceVarP(
		&normalizeCmdArgs.ReplaygainTags,
		"replaygain-tags",
		"t",
		[]string{
			"REPLAYGAIN_REFERENCE_LOUDNESS",
			"REPLAYGAIN_TRACK_GAIN",
			"REPLAYGAIN_TRACK_PEAK",
			"REPLAYGAIN_ALBUM_GAIN",
			"REPLAYGAIN_ALBUM_PEAK",
		},
		"ReplayGain tags",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeCmdArgs.SaveNormalizeArgumentsInTag,
		"save-normalize-arguments-in-tag",
		true,
		"save normalize arguments in tag",
	)
	normalizeCmd.PersistentFlags().StringVar(
		&normalizeCmdArgs.SaveNormalizeArgumentsInTagName,
		"save-normalize-arguments-in-tag-name",
		"METAFLAC_ARGUMENTS",
		"normalize arguments tag name",
	)

	viper.BindPFlags(normalizeCmd.PersistentFlags())
}

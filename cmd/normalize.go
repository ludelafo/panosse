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

var (
	normalizeArguments                       []string
	normalizeIfNormalizeArgumentTagsMismatch bool
	normalizeIfAnyReplayGainTagsAreMissing   bool
	replaygainTags                           []string
	saveNormalizeArgumentsInTag              bool
	saveNormalizeArgumentsInTagName          string
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize <file 1> [<file 2>]...",
	Short: "Normalize FLAC files with ReplayGain",
	Long: `Normalize FLAC files with ReplayGain.

It calls metaflac to add the ReplayGain tags to the FLAC files.`,
	Example: `  # Normalize some FLAC files
  $ panosse normalize file1.flac file2.flac

  # Normalize all FLAC files in all directories in parallel for a depth of 1
  # This allows to consider the nested directories as one album for the normalization
  $ find . -mindepth 1 -maxdepth 1 -type d -print0 | sort -z | while IFS= read -r -d '' dir; do
    mapfile -d '' -t flac_files < <(find "$dir" -type f -name "*.flac" -print0)
  
    if [ ${#flac_files[@]} -ne 0 ]; then
      panosse normalize --verbose "${flac_files[@]}"
    fi
  done`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::normalize] ")

		// Get command line arguments from Viper
		normalizeArguments = viper.GetStringSlice("normalize-arguments")
		normalizeIfNormalizeArgumentTagsMismatch =
			viper.GetBool("normalize-if-normalize-argument-tags-mismatch")
		normalizeIfAnyReplayGainTagsAreMissing =
			viper.GetBool("normalize-if-any-replaygain-tags-are-missing")
		replaygainTags = viper.GetStringSlice("replaygain-tags")
		saveNormalizeArgumentsInTag =
			viper.GetBool("save-normalize-arguments-in-tag")
		saveNormalizeArgumentsInTagName =
			viper.GetString("save-normalize-arguments-in-tag-name")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFiles := args

		// Command action
		normalizeArgumentsAsString := strings.Join(normalizeArguments, " ")

		needToNormalize := false

		for _, flacFile := range flacFiles {
			normalizeArgumentsTagContent, err := utils.GetTag(
				metaflacCommandPath,
				saveNormalizeArgumentsInTagName,
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

			if normalizeIfNormalizeArgumentTagsMismatch && normalizeArgumentsTagContent == "" {
				needToNormalize = true
			}

			if normalizeIfAnyReplayGainTagsAreMissing {
				for _, replaygainTag := range replaygainTags {
					replaygainTagContent, _ := utils.GetTag(
						metaflacCommandPath,
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

		if needToNormalize {
			if !dryRun {
				err := utils.Normalize(
					metaflacCommandPath,
					normalizeArguments,
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

		if verbose {
			log.Printf("\"%s\" normalized\n", flacFiles)
		}

		for _, flacFile := range flacFiles {
			if saveNormalizeArgumentsInTag {
				if !dryRun {
					utils.RemoveTag(
						metaflacCommandPath,
						saveNormalizeArgumentsInTagName,
						flacFile,
					)
					utils.SetTag(
						metaflacCommandPath,
						saveNormalizeArgumentsInTagName,
						normalizeArgumentsAsString,
						flacFile,
					)
				}
			}
		}

		if verbose {
			log.Printf(
				"\"%s\" %s tag added\n",
				flacFiles,
				saveNormalizeArgumentsInTagName,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	normalizeCmd.PersistentFlags().StringSliceVarP(
		&normalizeArguments,
		"normalize-arguments", "a", []string{
			"--add-replay-gain",
		}, "arguments passed to flac to normalize the files",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeIfNormalizeArgumentTagsMismatch,
		"normalize-if-normalize-argument-tags-mismatch",
		true,
		"normalize if normalize arguments tags mismatch (missing or different)",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeIfAnyReplayGainTagsAreMissing,
		"normalize-if-any-replaygain-tags-are-missing",
		true,
		"normalize if any ReplayGain tags are missing",
	)
	normalizeCmd.PersistentFlags().StringSliceVarP(
		&replaygainTags,
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
		&saveNormalizeArgumentsInTag,
		"save-normalize-arguments-in-tag",
		true,
		"save normalize arguments in tag",
	)
	normalizeCmd.PersistentFlags().StringVar(
		&saveNormalizeArgumentsInTagName,
		"save-normalize-arguments-in-tag-name",
		"METAFLAC_ARGUMENTS",
		"normalize arguments tag name",
	)

	viper.BindPFlags(normalizeCmd.PersistentFlags())
}

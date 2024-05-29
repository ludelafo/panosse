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
	normalizeIfReplayGainTagsAreMissing      bool
	replaygainTags                           []string
	saveNormalizeArgumentsInTag              bool
	saveNormalizeArgumentsInTagName          string
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize <file 1> [<file 2>]...",
	Short: "Normalize FLAC files with ReplayGain",
	Long:  "Normalize FLAC files with ReplayGain.",
	Example: `  Normalize some FLAC files
    panosse normalize file1.flac file2.flac

  Normalize all FLAC files in a directory
    find . -name "*.flac" -exec panosse normalize {} \;

  Normalize all FLAC files in a directory and subdirectories, by subdirectories
    find . -type d -exec sh -c 'find "$1" -name "*.flac" -exec panosse normalize {} \;' _ {} \;`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		normalizeArguments = viper.GetStringSlice("normalize-arguments")
		normalizeIfNormalizeArgumentTagsMismatch =
			viper.GetBool("normalize-if-normalize-argument-tags-mismatch")
		normalizeIfReplayGainTagsAreMissing =
			viper.GetBool("normalize-if-replaygain-tags-are-missing")
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
				metaflacCommand,
				saveNormalizeArgumentsInTagName,
				flacFile,
			)

			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					resultCode := exitError.ExitCode()

					if verbose {
						log.Fatalf(
							"cannot get tag from file '%s' (exit code %d)",
							flacFile,
							resultCode,
						)
					}
				}

				os.Exit(1)
			}

			if normalizeIfNormalizeArgumentTagsMismatch && normalizeArgumentsTagContent == "" {
				needToNormalize = true
			}

			if normalizeIfReplayGainTagsAreMissing {
				for _, replaygainTag := range replaygainTags {
					replaygainTagContent, _ := utils.GetTag(
						metaflacCommand,
						replaygainTag,
						flacFile,
					)

					if replaygainTagContent == "" {
						needToNormalize = true
						break
					}
				}
			}

			if saveNormalizeArgumentsInTag {
				if !dryRun {
					utils.RemoveTag(
						metaflacCommand,
						saveNormalizeArgumentsInTagName,
						flacFile,
					)
					utils.SetTag(
						metaflacCommand,
						saveNormalizeArgumentsInTagName,
						normalizeArgumentsAsString,
						flacFile,
					)
				}

				if verbose {
					log.Printf(
						"file '%s' %s tag added\n",
						flacFile,
						saveNormalizeArgumentsInTagName,
					)
				}
			}
		}

		if needToNormalize {
			if !dryRun {
				err := utils.Normalize(
					flacCommand,
					normalizeArguments,
					flacFiles,
				)

				if err != nil {
					if exitError, ok := err.(*exec.ExitError); ok {
						resultCode := exitError.ExitCode()

						if verbose {
							log.Fatalf(
								"error normalizing files '%s' (exit code %d)",
								flacFiles,
								resultCode,
							)
						}
					}

					os.Exit(1)
				}
			}
		}

		if verbose {
			log.Printf("files '%s' normalized\n", flacFiles)
		}
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	normalizeCmd.PersistentFlags().StringSliceVarP(
		&normalizeArguments,
		"normalize-arguments", "a", []string{
			"--add-replay-gain",
		}, "normalize arguments",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeIfNormalizeArgumentTagsMismatch,
		"normalize-if-normalize-argument-tags-mismatch",
		true,
		"normalize if normalize arguments tag is missing",
	)
	normalizeCmd.PersistentFlags().BoolVar(
		&normalizeIfReplayGainTagsAreMissing,
		"normalize-if-replaygain-tags-are-missing",
		true,
		"normalize if ReplayGain tags are missing",
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
		"save ReplayGain settings in tag",
	)
	normalizeCmd.PersistentFlags().StringVar(
		&saveNormalizeArgumentsInTagName,
		"save-normalize-arguments-in-tag-name",
		"METAFLAC_ARGUMENTS",
		"normalize arguments tag name",
	)

	viper.BindPFlags(normalizeCmd.PersistentFlags())
}

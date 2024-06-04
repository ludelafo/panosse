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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
type CleanCmdArgs struct {
	CleanArguments []string `mapstructure:"clean-arguments"`
	TagsToKeep     []string `mapstructure:"tags-to-keep"`
}

var cleanCmdArgs CleanCmdArgs

var cleanCmd = &cobra.Command{
	Use:   "clean <file>",
	Short: "Clean FLAC files from blocks and tags",
	Long: `Clean FLAC files from blocks and tags.

It calls metaflac to clean the FLAC files.`,
	Example: `  ## Clean a single FLAC file
  $ panosse clean file.flac

  ## Clean all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n1 -P$(nproc) panosse clean

  ## Clean all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 panosse clean`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::clean] ")

		// Get command line arguments from Viper
		viper.Unmarshal(&cleanCmdArgs)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		tagsToKeepMap := map[string]string{}

		for _, tagToKeep := range cleanCmdArgs.TagsToKeep {
			tagContent, err := utils.GetTag(
				rootCmdArgs.MetaflacCommandPath,
				tagToKeep,
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

			if tagContent != "" {
				tagsToKeepMap[tagToKeep] = tagContent
			}
		}

		if !rootCmdArgs.DryRun {
			utils.RemoveAllTags(rootCmdArgs.MetaflacCommandPath, flacFile)
		}

		for _, tagToKeep := range cleanCmdArgs.TagsToKeep {
			tagContent, ok := tagsToKeepMap[tagToKeep]

			if !rootCmdArgs.DryRun && ok {
				utils.SetTag(
					rootCmdArgs.MetaflacCommandPath,
					tagToKeep,
					tagContent,
					flacFile,
				)
			}
		}

		if !rootCmdArgs.DryRun {
			err := utils.Clean(
				rootCmdArgs.MetaflacCommandPath,
				cleanCmdArgs.CleanArguments,
				flacFile,
			)

			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					resultCode := exitError.ExitCode()

					log.Fatalf(
						"ERROR - cannot clean file \"%s\" (exit code %d)",
						flacFile,
						resultCode,
					)
				}

				os.Exit(1)
			}
		}

		if rootCmdArgs.Verbose {
			log.Printf("\"%s\" cleaned\n", flacFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().StringSliceVarP(
		&cleanCmdArgs.CleanArguments,
		"clean-arguments",
		"a",
		[]string{
			"--remove",
			"--dont-use-padding",
			"--block-type=APPLICATION",
			"--block-type=CUESHEET",
			"--block-type=PADDING",
			"--block-type=PICTURE",
			"--block-type=SEEKTABLE",
		},
		"arguments passed to metaflac to clean the file",
	)
	cleanCmd.PersistentFlags().StringSliceVarP(
		&cleanCmdArgs.TagsToKeep,
		"tags-to-keep",
		"t",
		[]string{
			"ALBUM",
			"ALBUMARTIST",
			"ARTIST",
			"COMMENT",
			"DISCNUMBER",
			"FLAC_ARGUMENTS",
			"GENRE",
			"METAFLAC_ARGUMENTS",
			"REPLAYGAIN_REFERENCE_LOUDNESS",
			"REPLAYGAIN_ALBUM_GAIN",
			"REPLAYGAIN_ALBUM_PEAK",
			"REPLAYGAIN_TRACK_GAIN",
			"REPLAYGAIN_TRACK_PEAK",
			"TITLE",
			"TRACKNUMBER",
			"TOTALDISCS",
			"TOTALTRACKS",
			"YEAR",
		},
		"tags to keep in the file",
	)

	viper.BindPFlags(cleanCmd.PersistentFlags())
}

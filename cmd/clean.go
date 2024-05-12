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
	"fmt"
	"os"
	"panosse/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
var (
	fillMissingTags     bool
	fillMissingTagsWith string
	cleanArguments      []string
	tagsToKeep          []string
)

var cleanCmd = &cobra.Command{
	Use:   "clean <file>",
	Short: "Clean FLAC files from tags and blocks",
	Long:  `Clean FLAC files from tags and blocks.`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		fillMissingTags = viper.GetBool("fill-missing-tags")
		fillMissingTagsWith = viper.GetString("fill-missing-tags-with")
		cleanArguments = viper.GetStringSlice("clean-arguments")
		tagsToKeep = viper.GetStringSlice("tags-to-keep")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		tagsToKeepMap := map[string]string{}

		for _, tagToKeep := range tagsToKeep {
			tagContent := utils.GetTag(metaflacCommand, tagToKeep, flacFile, verbose)

			if tagContent == "" && fillMissingTags {
				tagsToKeepMap[tagToKeep] = fillMissingTagsWith
			} else {
				tagsToKeepMap[tagToKeep] = tagContent
			}
		}

		if !dryRun {
			utils.RemoveAllTags(metaflacCommand, flacFile, verbose)
		}

		for tagToKeep, tagContent := range tagsToKeepMap {
			if !dryRun {
				utils.SetTag(metaflacCommand, tagToKeep, tagContent, flacFile, verbose)
			}
		}

		if !dryRun {
			utils.Clean(metaflacCommand, cleanArguments, flacFile, verbose)
		}

		if verbose {
			fmt.Fprintf(os.Stdout, "file '%s' cleaned\n", flacFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVar(&fillMissingTags, "fill-missing-tags", true, "fill missing tags")
	cleanCmd.PersistentFlags().StringVar(&fillMissingTagsWith, "fill-missing-tags-content", "No content for this tag", "fill missing tags content")
	cleanCmd.PersistentFlags().StringSliceVarP(&cleanArguments, "clean-arguments", "a", []string{
		"--remove",
		"--dont-use-padding",
		"--block-type=APPLICATION",
		"--block-type=CUESHEET",
		"--block-type=PADDING",
		"--block-type=PICTURE",
		"--block-type=SEEKTABLE",
	}, "clean arguments")
	cleanCmd.PersistentFlags().StringSliceVarP(&tagsToKeep, "tags-to-keep", "t", []string{
		"ALBUM",
		"ALBUMARTIST",
		"ARTIST",
		"COMMENT",
		"DISCNUMBER",
		"ENCODER_SETTINGS",
		"GENRE",
		"REPLAYGAIN_SETTINGS",
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
	}, "tags to keep")

	viper.BindPFlags(cleanCmd.PersistentFlags())
}

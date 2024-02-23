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
	"panosse/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
var (
	fillMissingTags         bool
	fillMissingTagsWith     string
	removeApplicationBlocks bool
	removeCuesheetBlocks    bool
	removePaddingBlocks     bool
	removePictureBlocks     bool
	removeSeektableBlocks   bool
	tagsToKeep              []string
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean FLAC files from tags and blocks",
	Long:  `Clean FLAC files from tags and blocks.`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		fillMissingTags = viper.GetBool("fill-missing-tags")
		fillMissingTagsWith = viper.GetString("fill-missing-tags-with")
		removeApplicationBlocks = viper.GetBool("remove-application-blocks")
		removeCuesheetBlocks = viper.GetBool("remove-cuesheet-blocks")
		removePaddingBlocks = viper.GetBool("remove-padding-blocks")
		removePictureBlocks = viper.GetBool("remove-picture-blocks")
		removeSeektableBlocks = viper.GetBool("remove-seektable-blocks")
		tagsToKeep = viper.GetStringSlice("tags-to-keep")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		cleanCommand := []string{"--remove", "--dont-use-padding"}

		if removeApplicationBlocks {
			cleanCommand = append(cleanCommand, "--block-type=APPLICATION")
		}

		if removeCuesheetBlocks {
			cleanCommand = append(cleanCommand, "--block-type=CUESHEET")
		}

		if removePaddingBlocks {
			cleanCommand = append(cleanCommand, "--block-type=PADDING")
		}

		if removePictureBlocks {
			cleanCommand = append(cleanCommand, "--block-type=PICTURE")
		}

		if removeSeektableBlocks {
			cleanCommand = append(cleanCommand, "--block-type=SEEKTABLE")
		}

		if !dryRun {
			utils.Clean(metaflacCommand, cleanCommand, flacFile, verbose)
		}

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
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVar(&fillMissingTags, "fill-missing-tags", true, "fill missing tags")
	cleanCmd.PersistentFlags().StringVar(&fillMissingTagsWith, "fill-missing-tags-content", "No content for this tag", "fill missing tags content")
	cleanCmd.PersistentFlags().BoolVar(&removeApplicationBlocks, "remove-application-blocks", true, "remove application blocks")
	cleanCmd.PersistentFlags().BoolVar(&removeCuesheetBlocks, "remove-cuesheet-blocks", true, "remove cuesheet blocks")
	cleanCmd.PersistentFlags().BoolVar(&removePaddingBlocks, "remove-padding-blocks", true, "remove padding blocks")
	cleanCmd.PersistentFlags().BoolVar(&removePictureBlocks, "remove-picture-blocks", true, "remove picture blocks")
	cleanCmd.PersistentFlags().BoolVar(&removeSeektableBlocks, "remove-seektable-blocks", true, "remove seektable blocks")
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

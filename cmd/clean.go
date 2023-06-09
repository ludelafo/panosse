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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fillMissingTags bool
var tagsToKeep []string
var removeApplicationBlocks bool
var removeCuesheetBlocks bool
var removeEmbeddedBlocks bool
var removePaddingBlocks bool
var removePictureBlocks bool
var removeSeektableBlocks bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean FLAC files from tags and blocks",
	Long:  `Clean FLAC files from tags and blocks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clean called")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVar(&fillMissingTags, "fill-missing-tags", true, "fill missing tags with \"[No <tag>]\"")
	cleanCmd.PersistentFlags().BoolVar(&removeApplicationBlocks, "remove-application-blocks", true, "remove application blocks")
	cleanCmd.PersistentFlags().BoolVar(&removeCuesheetBlocks, "remove-cuesheet-blocks", true, "remove cuesheet blocks")
	cleanCmd.PersistentFlags().BoolVar(&removeEmbeddedBlocks, "remove-embedded-blocks", true, "remove embedded blocks")
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

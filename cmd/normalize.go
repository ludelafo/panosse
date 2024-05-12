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

var normalizeArguments []string
var normalizeArgumentsTagName string
var normalizeIfReplayGainArgumentsTagIsMissing bool
var normalizeIfReplayGainTagsAreMissing bool
var normalizeTags []string
var saveNormalizeArgumentsInTag bool

var normalizeCmd = &cobra.Command{
	Use:   "normalize <file 1> [<file 2> ... <file N>]",
	Short: "Normalize FLAC files with ReplayGain",
	Long:  `Normalize FLAC files with ReplayGain.`,
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		normalizeArguments = viper.GetStringSlice("normalize-arguments")
		normalizeArgumentsTagName = viper.GetString("normalize-arguments-tag-name")
		normalizeIfReplayGainArgumentsTagIsMissing = viper.GetBool("normalize-if-normalize-arguments-tag-is-missing")
		normalizeIfReplayGainTagsAreMissing = viper.GetBool("normalize-if-replaygain-tags-are-missing")
		normalizeTags = viper.GetStringSlice("replaygain-tags")
		saveNormalizeArgumentsInTag = viper.GetBool("save-normalize-arguments-in-tag")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("normalize called")
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	normalizeCmd.PersistentFlags().StringSliceVarP(&normalizeArguments, "normalize-arguments", "a", []string{
		"--add-replay-gain",
	}, "normalize arguments")
	normalizeCmd.PersistentFlags().StringVar(&normalizeArgumentsTagName, "normalize-arguments-tag-name", "NORMALIZE_ARGUMENTS", "normalize arguments tag name")
	normalizeCmd.PersistentFlags().BoolVar(&normalizeIfReplayGainArgumentsTagIsMissing, "normalize-if-normalize-arguments-tag-is-missing", true, "normalize if normalize arguments tag is missing")
	normalizeCmd.PersistentFlags().BoolVar(&normalizeIfReplayGainTagsAreMissing, "normalize-if-replaygain-tags-are-missing", true, "normalize if ReplayGain tags are missing")
	normalizeCmd.PersistentFlags().StringSliceVarP(&normalizeTags, "replaygain-tags", "t", []string{
		"REPLAYGAIN_REFERENCE_LOUDNESS",
		"REPLAYGAIN_TRACK_GAIN",
		"REPLAYGAIN_TRACK_PEAK",
		"REPLAYGAIN_ALBUM_GAIN",
		"REPLAYGAIN_ALBUM_PEAK",
	}, "ReplayGain tags")
	normalizeCmd.PersistentFlags().BoolVar(&saveNormalizeArgumentsInTag, "save-normalize-arguments-in-tag", true, "save ReplayGain settings in tag")

	viper.BindPFlags(normalizeCmd.PersistentFlags())
}

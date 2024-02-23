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

var normalizeIfReplayGainSettingsTagIsMissing bool
var normalizeIfReplayGainTagsAreMissing bool
var replayGainSettings []string
var replayGainSettingsTagName string
var replayGainTags []string
var saveReplayGainSettingsInTag bool

var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize FLAC files with ReplayGain",
	Long: `Normalize FLAC files by calculating and adding ReplayGain to them.
	
Each directory containing FLAC files will be used to calculate the normalization.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("normalize called")
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	normalizeCmd.PersistentFlags().BoolVar(&normalizeIfReplayGainSettingsTagIsMissing, "normalize-if-replaygain-settings-tag-is-missing", true, "normalize if ReplayGain settings tag is missing")
	normalizeCmd.PersistentFlags().BoolVar(&normalizeIfReplayGainTagsAreMissing, "normalize-if-replaygain-tags-are-missing", true, "normalize if ReplayGain tags are missing")
	normalizeCmd.PersistentFlags().StringSliceVarP(&replayGainSettings, "replaygain-settings", "r", []string{
		"--add-replay-gain",
	}, "ReplayGain settings")
	normalizeCmd.PersistentFlags().StringVar(&replayGainSettingsTagName, "replaygain-settings-tag-name", "REPLAYGAIN_SETTINGS", "ReplayGain settings tag name")
	normalizeCmd.PersistentFlags().StringSliceVarP(&replayGainTags, "replaygain-tags", "t", []string{
		"REPLAYGAIN_REFERENCE_LOUDNESS",
		"REPLAYGAIN_TRACK_GAIN",
		"REPLAYGAIN_TRACK_PEAK",
		"REPLAYGAIN_ALBUM_GAIN",
		"REPLAYGAIN_ALBUM_PEAK",
	}, "ReplayGain tags")
	normalizeCmd.PersistentFlags().BoolVar(&saveReplayGainSettingsInTag, "save-replaygain-settings-in-tag", true, "save ReplayGain settings in tag")

	viper.BindPFlags(normalizeCmd.PersistentFlags())
}

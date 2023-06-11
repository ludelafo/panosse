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

var encodeIfVersionsMismatch bool
var encodeIfEncoderSettingsTagIsMissing bool
var encoderSettingsTagName string
var encoderSettings []string
var extractVersionFromFlacVersionRegex string
var extractVersionFromMetaflacShowVendorTagRegex string
var saveEncoderSettingsInTag bool

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode FLAC files",
	Long: `Encode FLAC files.
	
Each FLAC file will be encoded in-place.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("encode called")
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.PersistentFlags().BoolVar(&encodeIfVersionsMismatch, "encode-if-versions-mismatch", true, "encode if versions mismatch between current FLAC version and file's FLAC version")
	encodeCmd.PersistentFlags().BoolVar(&encodeIfEncoderSettingsTagIsMissing, "encode-if-encoder-settings-tag-is-missing", true, "encode if encoder settings tag is missing")
	encodeCmd.PersistentFlags().StringVar(&encoderSettingsTagName, "encoder-settings-tag-name", "ENCODER_SETTINGS", "encoder settings tag name")
	encodeCmd.PersistentFlags().StringSliceVarP(&encoderSettings, "encoder-settings", "e", []string{
		"flac",
		"--compression-level-8",
		"--delete-input-file",
		"--no-padding",
		"--force",
		"--verify",
		"--warnings-as-errors",
		"--silent",
	}, "encoder settings")
	encodeCmd.PersistentFlags().StringVar(&extractVersionFromFlacVersionRegex, "extract-version-from-flac-version-regex", "flac ([\\d]+.[\\d]+.[\\d]+)", "extract version from `flac --version` regex")
	encodeCmd.PersistentFlags().StringVar(&extractVersionFromMetaflacShowVendorTagRegex, "extract-version-from-metaflac-show-vendor-tag-regex", "reference libFLAC ([\\d]+.[\\d]+.[\\d]+) [\\d]+", "extract version from `metaflac --show-vendor-tag` regex")
	encodeCmd.PersistentFlags().BoolVar(&saveEncoderSettingsInTag, "save-encoder-settings-in-tag", true, "save encoder settings in tag")

	viper.BindPFlags(encodeCmd.PersistentFlags())
}

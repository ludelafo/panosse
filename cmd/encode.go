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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
var (
	encodeIfFlacVersionsMismatch   bool
	encodeIfEncodeSettingsMismatch bool
	saveEncodeSettingsInTag        bool
	encodeSettingsTagName          string
	encodeSettings                 []string
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode FLAC file",
	Long:  `Encode FLAC file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		flacVersionFromFlacCommand := utils.GetFlacVersionFromFlacCommand(flacCommand, verbose)
		flacVersionFromFlacFile := utils.GetFlacVersionFromFlacFile(metaflacCommand, flacFile, verbose)
		encodeSettingsTagContent := utils.GetTag(metaflacCommand, encodeSettingsTagName, flacFile, verbose)
		encodeSettingsAsString := strings.Join(encodeSettings, " ")

		needToEncode := false

		if encodeIfFlacVersionsMismatch && flacVersionFromFlacCommand != flacVersionFromFlacFile {
			needToEncode = true
		}

		if encodeIfEncodeSettingsMismatch && encodeSettingsTagContent != encodeSettingsAsString {
			needToEncode = true
		}

		if needToEncode {
			if !dryRun {
				utils.Encode(flacCommand, encodeSettings, flacFile, verbose)
			}

			if verbose {
				fmt.Fprintf(os.Stdout, "file '%s' encoded\n", flacFile)
			}
		}

		if saveEncodeSettingsInTag {
			if !dryRun {
				utils.RemoveTag(metaflacCommand, encodeSettingsTagName, flacFile, verbose)
				utils.SetTag(metaflacCommand, encodeSettingsTagName, encodeSettingsAsString, flacFile, verbose)
			}

			if verbose {
				fmt.Fprintf(os.Stdout, "file '%s' tag added\n", flacFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.PersistentFlags().BoolVar(&encodeIfFlacVersionsMismatch, "encode-if-flac-versions-mismatch", true, "encode if flac versions mismatch between host's flac version and file's flac version")
	encodeCmd.PersistentFlags().BoolVar(&encodeIfEncodeSettingsMismatch, "encode-if-encode-settings-mismatch", true, "encode if encode settings mismatch")
	encodeCmd.PersistentFlags().BoolVar(&saveEncodeSettingsInTag, "save-encode-settings-in-tag", true, "save encode settings in tag")
	encodeCmd.PersistentFlags().StringVar(&encodeSettingsTagName, "encode-settings-tag-name", "FLAC_SETTINGS", "encode settings tag name")
	encodeCmd.PersistentFlags().StringSliceVar(&encodeSettings, "encode-settings", []string{
		"--compression-level-8",
		"--delete-input-file",
		"--no-padding",
		"--force",
		"--verify",
		"--warnings-as-errors",
		"--silent",
	}, "encode settings")

	viper.BindPFlags(encodeCmd.PersistentFlags())
}

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
	encodeIfFlacVersionsMismatch    bool
	encodeIfEncodeArgumentsMismatch bool
	saveEncodeArgumentsInTag        bool
	encodeArgumentsTagName          string
	encodeArguments                 []string
)

var encodeCmd = &cobra.Command{
	Use:   "encode <file>",
	Short: "Encode FLAC file",
	Long:  `Encode FLAC file.`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		encodeIfFlacVersionsMismatch = viper.GetBool("encode-if-flac-versions-mismatch")
		encodeIfEncodeArgumentsMismatch = viper.GetBool("encode-if-encode-arguments-mismatch")
		saveEncodeArgumentsInTag = viper.GetBool("save-encode-arguments-in-tag")
		encodeArgumentsTagName = viper.GetString("encode-arguments-tag-name")
		encodeArguments = viper.GetStringSlice("encode-arguments")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		flacVersionFromFlacCommand := utils.GetFlacVersionFromFlacCommand(flacCommand, verbose)
		flacVersionFromFlacFile := utils.GetFlacVersionFromFlacFile(metaflacCommand, flacFile, verbose)
		encodeArgumentsTagContent := utils.GetTag(metaflacCommand, encodeArgumentsTagName, flacFile, verbose)
		encodeArgumentsAsString := strings.Join(encodeArguments, " ")

		needToEncode := false

		if encodeIfFlacVersionsMismatch && flacVersionFromFlacCommand != flacVersionFromFlacFile {
			needToEncode = true
		}

		if encodeIfEncodeArgumentsMismatch && encodeArgumentsTagContent != encodeArgumentsAsString {
			needToEncode = true
		}

		if needToEncode {
			if !dryRun {
				utils.Encode(flacCommand, encodeArguments, flacFile, verbose)
			}

			if verbose {
				fmt.Fprintf(os.Stdout, "file '%s' encoded\n", flacFile)
			}
		} else {
			if verbose {
				fmt.Fprintf(os.Stdout, "file '%s' already encoded\n", flacFile)
			}
		}

		if saveEncodeArgumentsInTag {
			if !dryRun {
				utils.RemoveTag(metaflacCommand, encodeArgumentsTagName, flacFile, verbose)
				utils.SetTag(metaflacCommand, encodeArgumentsTagName, encodeArgumentsAsString, flacFile, verbose)
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
	encodeCmd.PersistentFlags().BoolVar(&encodeIfEncodeArgumentsMismatch, "encode-if-encode-arguments-mismatch", true, "encode if encode arguments mismatch")
	encodeCmd.PersistentFlags().BoolVar(&saveEncodeArgumentsInTag, "save-encode-arguments-in-tag", true, "save encode arguments in tag")
	encodeCmd.PersistentFlags().StringVar(&encodeArgumentsTagName, "encode-arguments-tag-name", "FLAC_SETTINGS", "encode arguments tag name")
	encodeCmd.PersistentFlags().StringSliceVarP(&encodeArguments, "encode-arguments", "a", []string{
		"--compression-level-8",
		"--delete-input-file",
		"--no-padding",
		"--force",
		"--verify",
		"--warnings-as-errors",
		"--silent",
	}, "encode arguments")

	viper.BindPFlags(encodeCmd.PersistentFlags())
}

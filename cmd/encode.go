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

// Command arguments
var (
	encodeArguments                    []string
	encodeIfEncodeArgumentTagsMismatch bool
	encodeIfFlacVersionsMismatch       bool
	saveEncodeArgumentsInTag           bool
	saveEncodeArgumentsInTagName       string
)

var encodeCmd = &cobra.Command{
	Use:   "encode <file>",
	Short: "Encode FLAC file",
	Long:  "Encode FLAC file.",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		encodeArguments = viper.GetStringSlice("encode-arguments")
		encodeIfEncodeArgumentTagsMismatch =
			viper.GetBool("encode-if-encode-argument-tags-mismatch")
		encodeIfFlacVersionsMismatch =
			viper.GetBool("encode-if-flac-versions-mismatch")
		saveEncodeArgumentsInTag =
			viper.GetBool("save-encode-arguments-in-tag")
		saveEncodeArgumentsInTagName =
			viper.GetString("save-encode-arguments-in-tag-name")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		flacVersionFromFlacCommand, err := utils.GetFlacVersionFromFlacCommand(
			flacCommand,
		)

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				resultCode := exitError.ExitCode()

				if verbose {
					log.Fatalf(
						"cannot get flac version from flac command (exit code %d)",
						resultCode,
					)
				}
			}

			os.Exit(1)
		}

		flacVersionFromFlacFile, err := utils.GetFlacVersionFromFlacFile(
			metaflacCommand,
			flacFile,
		)

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				resultCode := exitError.ExitCode()

				if verbose {
					log.Fatalf(
						"cannot get flac version from file '%s' (exit code %d)",
						flacFile,
						resultCode,
					)
				}
			}

			os.Exit(1)
		}

		encodeArgumentsTagContent, _ := utils.GetTag(
			metaflacCommand,
			saveEncodeArgumentsInTagName,
			flacFile,
		)

		encodeArgumentsAsString := strings.Join(encodeArguments, " ")

		needToEncode := false

		if encodeIfFlacVersionsMismatch && flacVersionFromFlacCommand != flacVersionFromFlacFile {
			needToEncode = true
		}

		if encodeIfEncodeArgumentTagsMismatch && encodeArgumentsTagContent != encodeArgumentsAsString {
			needToEncode = true
		}

		if needToEncode {
			if !dryRun {
				utils.Encode(flacCommand, encodeArguments, flacFile)
			}

			if verbose {
				log.Printf("file '%s' encoded\n", flacFile)
			}
		}

		if saveEncodeArgumentsInTag {
			if !dryRun {
				utils.RemoveTag(
					metaflacCommand,
					saveEncodeArgumentsInTagName,
					flacFile,
				)
				utils.SetTag(
					metaflacCommand,
					saveEncodeArgumentsInTagName,
					encodeArgumentsAsString,
					flacFile,
				)
			}

			if verbose {
				log.Printf(
					"file '%s' %s tag added\n",
					flacFile,
					saveEncodeArgumentsInTagName,
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.PersistentFlags().StringSliceVarP(
		&encodeArguments,
		"encode-arguments",
		"a",
		[]string{
			"--compression-level-8",
			"--delete-input-file",
			"--no-padding",
			"--force",
			"--verify",
			"--warnings-as-errors",
			"--silent",
		},
		"encode arguments",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeIfEncodeArgumentTagsMismatch,
		"encode-if-encode-argument-tags-mismatch",
		true,
		"encode if encode arguments mismatch",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeIfFlacVersionsMismatch,
		"encode-if-flac-versions-mismatch",
		true,
		"encode if flac versions mismatch between host's flac version and file's flac version",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&saveEncodeArgumentsInTag,
		"save-encode-arguments-in-tag",
		true,
		"save encode arguments in tag",
	)
	encodeCmd.PersistentFlags().StringVar(
		&saveEncodeArgumentsInTagName,
		"save-encode-arguments-in-tag-name",
		"FLAC_ARGUMENTS",
		"encode arguments tag name",
	)

	viper.BindPFlags(encodeCmd.PersistentFlags())
}

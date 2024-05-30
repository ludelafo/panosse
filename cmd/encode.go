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
	Short: "Encode FLAC files",
	Long: `Encode FLAC files.

It calls flac to encode the FLAC files.`,
	Example: `  # Encode a single FLAC file
  $ panosse encode file.flac

  # Encode all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n1 -P$(nproc) panosse encode`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::encode] ")

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
			flacCommandPath,
		)

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				resultCode := exitError.ExitCode()

				log.Fatalf(
					"ERROR - cannot get flac version from flac command (exit code %d)",
					resultCode,
				)
			}

			os.Exit(1)
		}

		flacVersionFromFlacFile, err := utils.GetFlacVersionFromFlacFile(
			metaflacCommandPath,
			flacFile,
		)

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				resultCode := exitError.ExitCode()

				log.Fatalf(
					"ERROR - cannot get flac version from file \"%s\" (exit code %d)",
					flacFile,
					resultCode,
				)
			}

			os.Exit(1)
		}

		encodeArgumentsTagContent, _ := utils.GetTag(
			metaflacCommandPath,
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
				utils.Encode(flacCommandPath, encodeArguments, flacFile)
			}

			if verbose {
				log.Printf("\"%s\" encoded\n", flacFile)
			}
		}

		if saveEncodeArgumentsInTag {
			if !dryRun {
				utils.RemoveTag(
					metaflacCommandPath,
					saveEncodeArgumentsInTagName,
					flacFile,
				)
				utils.SetTag(
					metaflacCommandPath,
					saveEncodeArgumentsInTagName,
					encodeArgumentsAsString,
					flacFile,
				)
			}

			if verbose {
				log.Printf(
					"\"%s\" %s tag added\n",
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
		"arguments passed to flac to encode the file",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeIfEncodeArgumentTagsMismatch,
		"encode-if-encode-argument-tags-mismatch",
		true,
		"encode if encode argument tags mismatch (missing or different)",
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

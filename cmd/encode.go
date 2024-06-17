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
type EncodeCmdArgs struct {
	EncodeArguments                    []string `mapstructure:"encode-arguments"`
	EncodeIfEncodeArgumentTagsMismatch bool     `mapstructure:"encode-if-encode-argument-tags-mismatch"`
	EncodeIfFlacVersionsMismatch       bool     `mapstructure:"encode-if-flac-versions-mismatch"`
	SaveEncodeArgumentsInTag           bool     `mapstructure:"save-encode-arguments-in-tag"`
	SaveEncodeArgumentsInTagName       string   `mapstructure:"save-encode-arguments-in-tag-name"`
}

var encodeCmdArgs EncodeCmdArgs

var encodeCmd = &cobra.Command{
	Use:   "encode <file>",
	Short: "Encode FLAC files",
	Long: `Encode FLAC files.

It calls flac to encode the FLAC files.`,
	Example: `  ## Encode a single FLAC file
  $ panosse encode file.flac

  ## Encode all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n 1 -P $(nproc) panosse encode

  ## Encode all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 panosse encode`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::encode] ")

		// Get command line arguments from Viper
		viper.Unmarshal(&encodeCmdArgs)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		// Command action
		flacVersionFromFlacCommand, err := utils.GetFlacVersionFromFlacCommand(
			rootCmdArgs.FlacCommandPath,
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
			rootCmdArgs.MetaflacCommandPath,
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
			rootCmdArgs.MetaflacCommandPath,
			encodeCmdArgs.SaveEncodeArgumentsInTagName,
			flacFile,
		)

		encodeArgumentsAsString := strings.Join(
			encodeCmdArgs.EncodeArguments,
			" ",
		)

		needToEncode := false

		if encodeCmdArgs.EncodeIfFlacVersionsMismatch && flacVersionFromFlacCommand != flacVersionFromFlacFile {
			needToEncode = true
		}

		if encodeCmdArgs.EncodeIfEncodeArgumentTagsMismatch && encodeArgumentsTagContent != encodeArgumentsAsString {
			needToEncode = true
		}

		if needToEncode || rootCmdArgs.Force {
			encodeCmdArgs.EncodeArguments = append(
				encodeCmdArgs.EncodeArguments,
				"--force",
				"--silent",
			)

			if !rootCmdArgs.DryRun {
				var output, err = utils.Encode(
					rootCmdArgs.FlacCommandPath,
					encodeCmdArgs.EncodeArguments,
					flacFile,
				)

				if err != nil {
					if exitError, ok := err.(*exec.ExitError); ok {
						resultCode := exitError.ExitCode()

						log.Fatalf(
							"ERROR - cannot encode file \"%s\" (exit code %d)",
							flacFile,
							resultCode,
						)

						if rootCmdArgs.Verbose {
							log.Fatalln(output)
						}
					}

					os.Exit(1)
				}
			}

			if rootCmdArgs.Verbose {
				log.Printf("\"%s\" encoded\n", flacFile)
			}
		}

		if encodeCmdArgs.SaveEncodeArgumentsInTag {
			if !rootCmdArgs.DryRun {
				utils.RemoveTag(
					rootCmdArgs.MetaflacCommandPath,
					encodeCmdArgs.SaveEncodeArgumentsInTagName,
					flacFile,
				)
				utils.SetTag(
					rootCmdArgs.MetaflacCommandPath,
					encodeCmdArgs.SaveEncodeArgumentsInTagName,
					encodeArgumentsAsString,
					flacFile,
				)
			}

			if rootCmdArgs.Verbose {
				log.Printf(
					"\"%s\" %s tag added\n",
					flacFile,
					encodeCmdArgs.SaveEncodeArgumentsInTagName,
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.PersistentFlags().StringSliceVarP(
		&encodeCmdArgs.EncodeArguments,
		"encode-arguments",
		"a",
		[]string{
			"--compression-level-8",
			"--exhaustive-model-search",
			"--no-padding",
			"--qlp-coeff-precision-search",
			"--verify",
			"--warnings-as-errors",
		},
		"arguments passed to flac to encode the file",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeCmdArgs.EncodeIfEncodeArgumentTagsMismatch,
		"encode-if-encode-argument-tags-mismatch",
		true,
		"encode if encode argument tags mismatch (missing or different)",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeCmdArgs.EncodeIfFlacVersionsMismatch,
		"encode-if-flac-versions-mismatch",
		true,
		"encode if flac versions mismatch between host's flac version and file's flac version",
	)
	encodeCmd.PersistentFlags().BoolVar(
		&encodeCmdArgs.SaveEncodeArgumentsInTag,
		"save-encode-arguments-in-tag",
		true,
		"save encode arguments in tag",
	)
	encodeCmd.PersistentFlags().StringVar(
		&encodeCmdArgs.SaveEncodeArgumentsInTagName,
		"save-encode-arguments-in-tag-name",
		"FLAC_ARGUMENTS",
		"encode arguments tag name",
	)

	viper.BindPFlags(encodeCmd.PersistentFlags())
}

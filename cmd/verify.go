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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
type VerifyCmdArgs struct {
	VerifyArguments []string `mapstructure:"verify-arguments"`
}

var verifyCmdArgs VerifyCmdArgs

var verifyCmd = &cobra.Command{
	Use:   "verify <file>",
	Short: "Verify FLAC files",
	Long: `Check the integrity of the FLAC files.

It calls metaflac to verify the FLAC files.`,
	Example: `  ## Verify a single FLAC file
  $ panosse verify file.flac

  ## Verify all FLAC files in the current directory recursively and in parallel
  $ find . -type f -name "*.flac" -print0 | xargs -0 -n 1 -P $(nproc) panosse verify

  ## Verify all FLAC files in the current directory recursively and in order
  # This approach is slower than the previous one but it can be useful to process
  # the files in a specific order (e.g., to follow the progression)
  $ find . -type f -name "*.flac" -print0 | sort -z | xargs -0 -n 1 panosse verify`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Set logger prefix for this file
		log.SetPrefix("[panosse::verify] ")

		// Get command line arguments from Viper
		viper.Unmarshal(&verifyCmdArgs)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		if !rootCmdArgs.DryRun {
			err := utils.Verify(
				rootCmdArgs.FlacCommandPath,
				verifyCmdArgs.VerifyArguments,
				flacFile,
			)

			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					resultCode := exitError.ExitCode()

					log.Fatalf(
						"ERROR - cannot verify file \"%s\" (exit code %d)",
						flacFile,
						resultCode,
					)
				}

				os.Exit(1)
			}
		}

		if rootCmdArgs.Verbose {
			log.Printf("\"%s\" verified\n", flacFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	cobra.OnInitialize()

	verifyCmd.PersistentFlags().StringSliceVarP(
		&verifyCmdArgs.VerifyArguments,
		"verify-arguments",
		"a",
		[]string{
			"--test",
			"--silent",
		},
		"arguments passed to flac to verify the files",
	)

	viper.BindPFlags(verifyCmd.PersistentFlags())
}

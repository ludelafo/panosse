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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
var (
	verifyArguments []string
)

var verifyCmd = &cobra.Command{
	Use:     "verify <file>",
	Version: rootCmd.Version,
	Short:   "Verify FLAC files",
	Long:    `Verify FLAC files.`,
	Args:    cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Get command line arguments from Viper
		verifyArguments = viper.GetStringSlice("verify-arguments")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get arguments for the command
		flacFile := args[0]

		if !dryRun {
			utils.Verify(flacCommand, verifyArguments, flacFile, verbose)
		}

		if verbose {
			fmt.Fprintf(os.Stdout, "file '%s' verified\n", flacFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	cobra.OnInitialize()

	verifyCmd.PersistentFlags().StringSliceVarP(&verifyArguments, "verify-arguments", "a", []string{
		"--test",
		"--silent",
	}, "verify arguments")

	viper.BindPFlags(verifyCmd.PersistentFlags())
}

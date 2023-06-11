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
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify FLAC files",
	Long:  `Verify FLAC files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("verify called")

		var stdout bytes.Buffer
		flacVerify := exec.Command("flac", "--test", "uncommon_11 - file starting with unparsable data")
		flacVerify.Stdout = &stdout

		err := flacVerify.Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				// Command failed with non-zero exit code
				resultCode := exitError.ExitCode()
				fmt.Fprintf(os.Stderr, "flac execution failed with result code: %d, %s", resultCode, stdout.String())
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}

		fmt.Println(stdout.String())
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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
)

var normalizeIfTagsAreNotPresent bool
var normalizeTags []string
var normalizerSettingsTagName string
var normalizerSettings []string
var saveNormalizerSettingsInTag bool

var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize FLAC files with ReplayGain",
	Long: `Normalize FLAC files by calculating and adding ReplayGain to them.
	
Each directory containing FLAC files will be used to calculate the normalization.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("normalize called")
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// normalizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// normalizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

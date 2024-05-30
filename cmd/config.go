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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Version: rootCmd.Version,
	Short:   "Display panosse configuration",
	Long:    "Display panosse configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		log.SetPrefix("")

		// Get command line arguments from Viper (clean)
		cleanArguments = viper.GetStringSlice("clean-arguments")
		tagsToKeep = viper.GetStringSlice("tags-to-keep")

		// Get command line arguments from Viper (encode)
		encodeArguments = viper.GetStringSlice("encode-arguments")
		encodeIfEncodeArgumentTagsMismatch =
			viper.GetBool("encode-if-encode-argument-tags-mismatch")
		encodeIfFlacVersionsMismatch =
			viper.GetBool("encode-if-flac-versions-mismatch")
		saveEncodeArgumentsInTag =
			viper.GetBool("save-encode-arguments-in-tag")
		saveEncodeArgumentsInTagName =
			viper.GetString("save-encode-arguments-in-tag-name")

		// Get command line arguments from Viper (normalize)
		normalizeArguments = viper.GetStringSlice("normalize-arguments")
		normalizeIfNormalizeArgumentTagsMismatch =
			viper.GetBool("normalize-if-normalize-argument-tags-mismatch")
		normalizeIfAnyReplayGainTagsAreMissing =
			viper.GetBool("normalize-if-replaygain-tags-are-missing")
		replaygainTags = viper.GetStringSlice("replaygain-tags")
		saveNormalizeArgumentsInTag =
			viper.GetBool("save-normalize-arguments-in-tag")
		saveNormalizeArgumentsInTagName =
			viper.GetString("save-normalize-arguments-in-tag-name")

		// Get command line arguments from Viper (verify)
		verifyArguments = viper.GetStringSlice("verify-arguments")

	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Common")
		log.Printf("  config-file: %s\n", viper.ConfigFileUsed())
		log.Printf("  dry-run: %t\n", dryRun)
		log.Printf("  flac-command-path: %s\n", flacCommandPath)
		log.Printf("  metaflac-command-path: %s\n", metaflacCommandPath)
		log.Printf("  verbose: %t\n", verbose)

		log.Println()
		log.Println("Clean")
		log.Printf("  clean-arguments: %s\n", cleanArguments)
		log.Printf("  tags-to-keep: %s\n", tagsToKeep)

		log.Println()
		log.Println("Encode")
		log.Printf("  encode-arguments: %s\n", encodeArguments)
		log.Printf("  encode-if-encode-argument-tags-mismatch: %t\n", encodeIfEncodeArgumentTagsMismatch)
		log.Printf("  encode-if-flac-versions-mismatch: %t\n", encodeIfFlacVersionsMismatch)
		log.Printf("  save-encode-arguments-in-tag: %t\n", saveEncodeArgumentsInTag)
		log.Printf("  save-encode-arguments-in-tag-name: %s\n", saveEncodeArgumentsInTagName)

		log.Println()
		log.Println("Normalize")
		log.Printf("  normalize-arguments: %s\n", normalizeArguments)
		log.Printf("  normalize-if-normalize-argument-tags-mismatch: %t\n", normalizeIfNormalizeArgumentTagsMismatch)
		log.Printf("  normalize-if-replaygain-tags-are-missing: %t\n", normalizeIfAnyReplayGainTagsAreMissing)
		log.Printf("  replaygain-tags: %s\n", replaygainTags)
		log.Printf("  save-normalize-arguments-in-tag: %t\n", saveEncodeArgumentsInTag)
		log.Printf("  save-normalize-arguments-in-tag-name: %s\n", saveNormalizeArgumentsInTagName)

		log.Println()
		log.Println("Verify")
		log.Printf("  verify-arguments: %s\n", verifyArguments)

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	cobra.OnInitialize()
}

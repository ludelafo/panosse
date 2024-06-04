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
	Use:   "config",
	Short: "Display panosse configuration",
	Long:  "Display panosse configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		log.SetPrefix("")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Common")
		log.Printf("  config-file: %s\n", viper.ConfigFileUsed())
		log.Printf("  dry-run: %t\n", RootCmdArgs.DryRun)
		log.Printf("  flac-command-path: %s\n", RootCmdArgs.FlacCommandPath)
		log.Printf("  metaflac-command-path: %s\n", RootCmdArgs.MetaflacCommandPath)
		log.Printf("  verbose: %t\n", RootCmdArgs.Verbose)

		log.Println()
		log.Println("Clean")
		log.Printf("  clean-arguments: %s\n", CleanCmdArgs.CleanArguments)
		log.Printf("  tags-to-keep: %s\n", CleanCmdArgs.TagsToKeep)

		log.Println()
		log.Println("Encode")
		log.Printf("  encode-arguments: %s\n", EncodeCmdArgs.EncodeArguments)
		log.Printf("  encode-if-encode-argument-tags-mismatch: %t\n", EncodeCmdArgs.EncodeIfEncodeArgumentTagsMismatch)
		log.Printf("  encode-if-flac-versions-mismatch: %t\n", EncodeCmdArgs.EncodeIfFlacVersionsMismatch)
		log.Printf("  save-encode-arguments-in-tag: %t\n", EncodeCmdArgs.SaveEncodeArgumentsInTag)
		log.Printf("  save-encode-arguments-in-tag-name: %s\n", EncodeCmdArgs.SaveEncodeArgumentsInTagName)

		log.Println()
		log.Println("Normalize")
		log.Printf("  normalize-arguments: %s\n", NormalizeCmdArgs.NormalizeArguments)
		log.Printf("  normalize-if-normalize-argument-tags-mismatch: %t\n", NormalizeCmdArgs.NormalizeIfNormalizeArgumentTagsMismatch)
		log.Printf("  normalize-if-replaygain-tags-are-missing: %t\n", NormalizeCmdArgs.NormalizeIfAnyReplayGainTagsAreMissing)
		log.Printf("  replaygain-tags: %s\n", NormalizeCmdArgs.ReplaygainTags)
		log.Printf("  save-normalize-arguments-in-tag: %t\n", NormalizeCmdArgs.SaveNormalizeArgumentsInTag)
		log.Printf("  save-normalize-arguments-in-tag-name: %s\n", NormalizeCmdArgs.SaveNormalizeArgumentsInTagName)

		log.Println()
		log.Println("Verify")
		log.Printf("  verify-arguments: %s\n", VerifyCmdArgs.VerifyArguments)

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	cobra.OnInitialize()
}

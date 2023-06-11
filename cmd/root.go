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
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var input string
var force bool
var dryRun bool
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "panosse",
	Short: "Clean, encode, normalize and verify your FLAC music library",
	Long: `panosse is a CLI tool to clean, encode, normalize and verify your FLAC music library.

It is a wrapper around flac and metaflac and uses Cobra and Viper under the hood.

Examples:

Clean FLAC files
panosse clean

Encode FLAC files
panosse encode

Normalize FLAC files
panosse normalize

Verify FLAC files
panosse verify`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "C", "", "config file (default is config.yaml or $HOME/.panosse/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&input, "input", "I", "", "where to look for music files (default is current directory)")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "F", false, "ignore warnings and errors")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "D", false, "dry run")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose output")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config.yaml in current directory or $HOME/.panosse
		viper.AddConfigPath(".")
		viper.AddConfigPath(path.Join(home, ".panosse"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("PANOSSE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err == nil {
		fmt.Fprintf(os.Stdout, "Using config file: %v\n", viper.ConfigFileUsed())
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

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
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command arguments
type RootCmdArgs struct {
	FlacCommandPath     string `mapstructure:"flac-command-path"`
	MetaflacCommandPath string `mapstructure:"metaflac-command-path"`
	ConfigFile          string
	DryRun              bool `mapstructure:"dry-run"`
	Verbose             bool `mapstructure:"verbose"`
}

var rootCmdArgs RootCmdArgs

var rootCmd = &cobra.Command{
	Use:     "panosse",
	Version: "0.1.0",
	Short:   "Clean, encode, normalize, and verify your FLAC music library",
	Long: `panosse is a CLI tool to clean, encode, normalize, and verify your FLAC music library.

panosse is merely a wrapper around flac and metaflac and uses Cobra and Viper under the hood.

panosse is licensed under the GNU Affero General Public License (GNU AGPL-3.0).

For more information, see https://github.com/ludelafo/panosse.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig, initViper)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SetVersionTemplate("panosse v{{.Version}}\n")

	rootCmd.PersistentFlags().StringVarP(
		&rootCmdArgs.ConfigFile,
		"config-file",
		"C",
		"",
		"config file to use (optional - will use \"config.yaml\" or \"~/.panosse/config.yaml\" if available)",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&rootCmdArgs.DryRun,
		"dry-run",
		"D",
		false,
		"perform a trial run with no changes made",
	)
	rootCmd.PersistentFlags().StringVarP(
		&rootCmdArgs.FlacCommandPath,
		"flac-command-path",
		"F",
		"flac",
		"path to the flac command (checks in $PATH as well)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&rootCmdArgs.MetaflacCommandPath,
		"metaflac-command-path",
		"M",
		"metaflac",
		"path to the metaflac command (checks in $PATH as well)",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&rootCmdArgs.Verbose,
		"verbose",
		"V",
		false,
		"enable verbose output",
	)

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initLogger() {
	log.SetFlags(0)
}

func initConfig() {
	if rootCmdArgs.ConfigFile != "" {
		viper.SetConfigFile(rootCmdArgs.ConfigFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config.yaml in current directory or ~/.panosse
		viper.AddConfigPath(".")
		viper.AddConfigPath(path.Join(home, ".panosse"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("panosse")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func initViper() {
	// Get command line arguments from Viper
	viper.Unmarshal(&rootCmdArgs)
}

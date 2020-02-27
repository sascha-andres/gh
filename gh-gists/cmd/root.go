/*
Copyright © 2020 Sascha Andres <sascha.andres@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"livingit.de/code/gh/helper"
	"livingit.de/code/gh/wrapper"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-gists",
	Short: "List gists for you or an organization",
	Long: `Returns gists as JSON stream.

When --organization is provided list gists for that organization/user`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := helper.Must("token")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading GitHub token: %s\n", err)
			os.Exit(1)
		}

		organization := viper.GetString("gists.organization")

		logger := logrus.
			WithField("package", "cmd").
			WithField("method", "gists::Run")

		logger.Infof("about to list gists in [%s]", organization)

		w, err := wrapper.NewGitHubWrapper(token)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		gists, err := w.GistsList(organization)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		for _, gist := range gists {
			data, err := json.Marshal(gist)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			fmt.Println(string(data))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-gists.yaml)")

	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub token")
	rootCmd.PersistentFlags().StringP("log-level", "l", "warn", "Set log level (defaulting to warn)\nmay break pipes as log messages appear within json stream")
	rootCmd.PersistentFlags().StringP("organization", "o", "", "Name of organization")

	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("gists.organization", rootCmd.PersistentFlags().Lookup("organization"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gh" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.ReadInConfig()

	logLevel := viper.GetString("log-level")
	if logLevel == "" {
		logLevel = "warn"
	}

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing loglevel: %s", err)
		os.Exit(1)
	}
	logrus.SetLevel(lvl)
}

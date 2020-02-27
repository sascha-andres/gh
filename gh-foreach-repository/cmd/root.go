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
	Use:   "gh-foreach-repository",
	Short: "Return GitHub repositories",
	Long:  `Return a list of GitHub repositories as a json stream`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.
			WithField("package", "cmd").
			WithField("method", "foreach-repository::Run")

		token, err := helper.Must("token")
		if err != nil {
			logger.Errorf("error reading GitHub token: %s", err)
			os.Exit(1)
		}
		affiliation := viper.GetString("foreach.repository.affiliation")
		visibility := viper.GetString("foreach.repository.visibility")

		w, err := wrapper.NewGitHubWrapper(token)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		repos, err := w.RepositoriesList(affiliation, visibility)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		for _, r := range repos {
			data, err := json.Marshal(r)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", data)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-foreach-repository.yaml)")

	rootCmd.PersistentFlags().StringP("visibility", "v", "", "Visibility of repositories to list. Can be one of all, public, or private.")
	rootCmd.PersistentFlags().StringP("affiliation", "a", "", "List repos of given affiliation[s]. Default: owner,collaborator,organization_member")
	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub token")

	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("foreach.repository.affiliation", rootCmd.PersistentFlags().Lookup("affiliation"))
	viper.BindPFlag("foreach.repository.visibility", rootCmd.PersistentFlags().Lookup("visibility"))
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

		viper.AddConfigPath(home)
		viper.SetConfigName(".gh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

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
	Use:   "gh-clone",
	Short: "Clone an organization",
	Long:  `Clones all repositories from an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := helper.Must("token")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading GitHub token: %s\n", err)
			os.Exit(1)
		}

		organization, err := helper.Must("clone.organization")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading organization: %s\n", err)
			os.Exit(1)
		}

		ssh := viper.GetBool("clone.ssh")
		user := viper.GetString("clone.user")
		email := viper.GetString("clone.email")

		logrus.
			WithField("package", "cmd").
			WithField("method", "clone::Run").
			Infof("about to clone [%s]", organization)

		w, err := wrapper.NewGitHubWrapper(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		repos, err := w.RepositoriesListByOrganization(organization)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, r := range repos {
			logrus.
				WithField("package", "cmd").
				WithField("method", "clone::Run").
				Infof("cloning [%s]", r.GetFullName())
			if ssh {
				_, _ = wrapper.Git("clone", r.GetSSHURL(), r.GetFullName())
			} else {
				_, _ = wrapper.Git("clone", r.GetHTMLURL(), r.GetFullName())
			}
			if user != "" {
				_, _ = wrapper.Git("-C", r.GetFullName(), "config", "user.name", user)
			}
			if email != "" {
				_, _ = wrapper.Git("-C", r.GetFullName(), "config", "user.email", email)
			}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-clone.yaml)")

	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub token")
	rootCmd.PersistentFlags().StringP("organization", "o", "", "Name of organization")
	rootCmd.Flags().StringP("user", "", "", "user.name for cloned repositories")
	rootCmd.Flags().StringP("email", "", "", "user.email for cloned repositories")
	rootCmd.Flags().BoolP("ssh", "s", true, "Use SSH to clone")

	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("clone.organization", rootCmd.PersistentFlags().Lookup("organization"))
	_ = viper.BindPFlag("clone.ssh", rootCmd.Flags().Lookup("ssh"))
	_ = viper.BindPFlag("clone.user", rootCmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("clone.email", rootCmd.Flags().Lookup("email"))
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

		// Search config in home directory with name ".gh-clone" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

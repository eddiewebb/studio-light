// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/eddiewebb/blync-studio-light/config"
)

// default in initConfig, unless passed as flag
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blync-studio-light",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.InfoLevel)
			log.Info("Verbose logging enabled")
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(configFromFlag)
	log.SetLevel(log.WarnLevel)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blync-studio-light.yaml)")
	rootCmd.PersistentFlags().IntP("device", "d", 0, "Device index for light to interface with")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Include info level logs")
	//nolint:errcheck
	viper.BindPFlag("device", rootCmd.PersistentFlags().Lookup("device"))
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configFromFlag() {
	C := config.Configuration{
		ConfigFile: cfgFile,
	}
	cfgFile = C.InitConfig()
}
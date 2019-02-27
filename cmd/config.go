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
    "bufio"
    "fmt"
    "os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/eddiewebb/blync-studio-light/config"
)

var C config.Configuration

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Viper Debug:")
		viper.Debug()
		if err := viper.Unmarshal(&C); err != nil{
			fmt.Println(err)
		}
		fmt.Println("Calendar Values:")
		fmt.Println(C.GoogleCalendar)
	},
}

var setConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Write new config based on user prompts",
	Long: `We'll ask some questions to generate an initial config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		newconfig := promptForCalValues()
		viper.Set("GoogleCalendar",newconfig)
		fmt.Println("ttempt to create new file"  + cfgFile)
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			fmt.Println(err)
		}
		fmt.Println("config set in " + viper.ConfigFileUsed())
	},
}

func promptForCalValues() config.GoogleCalendarConfiguration{
	gcalconfig := config.GoogleCalendarConfiguration{
		CalendarUri : prompt("What is the Calendar URI?"),
		ApiToken	: prompt("What is your API Token?"),
	}

	return gcalconfig
}

func prompt(message string) string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	return text
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setConfigCmd)
}

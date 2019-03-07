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
	"github.com/spf13/cobra"

	"github.com/eddiewebb/blync-studio-light/calendars"
	"github.com/eddiewebb/blync-studio-light/config"
	"github.com/spf13/viper"
)

var C config.Configuration

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setConfigCmd)
	configCmd.AddCommand(loginCmd)
	configCmd.AddCommand(setScheduleCmd)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or show configurations",
	Long: `Covers various configurations including:

- Default device index (0)
- Calendar information
- Rules`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Viper Debug:")
		viper.Debug()
		if err := viper.Unmarshal(&C); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Calendar Values:")
		fmt.Println(C.GoogleCalendar)
	},
}

var setConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Write new config based on user prompts",
	Long:  `We'll ask some questions to generate an initial config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		newconfig := promptForCalValues()
		viper.Set("GoogleCalendar", newconfig)
		log.Infoln("Attempt to create new file" + cfgFile)
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Fatal(err)
		}
		log.Infoln("config set in " + viper.ConfigFileUsed())
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Force new auth token with Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		calendarId := viper.GetString("googleCalendar.calendarId")
		email := viper.GetString("googleCalendar.email")
		calendars.RemoveExistingGoogleAuthToken()
		gcal,err := calendars.NewGoogleCalendar()
		if err != nil {
			log.Fatal(err)
		}
		color := gcal.GetColor(calendarId, email)
		log.Infof("The connected calendar shows status color or %s for user %s", color, email)
	},
}

var setScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Write new config based on user prompts",
	Long:  `We'll ask some questions to generate an initial config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		newconfig := promptForScheduleValues()
		viper.Set("Schedule", newconfig)
		fmt.Println("Attempt to update  file" + cfgFile)
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			fmt.Println(err)
		}
		fmt.Println("config set in " + viper.ConfigFileUsed())
	},
}

func promptForCalValues() config.GoogleCalendarConfiguration {
	gcalconfig := config.GoogleCalendarConfiguration{
		CalendarId: prompt("What is the Calendar ID (as seen in settings, usally your email)? "),
		Email:      prompt("What is the email of the attendee to base status on? "),
	}
	return gcalconfig
}

func promptForScheduleValues() (scheduleConfig config.StudioLightSchedule) {
	onTime := prompt("What ime (as HH:MM) should the light be allowed on? ")
	offTime := prompt("What time (as HH:MM) should the light always be off? ")
	daysOff := prompt("What days (as Saturday, Sunday, etc) should the light always be off? ")

	scheduleConfig, err := config.NewSchedule(onTime, offTime, daysOff)
	if err != nil {
		log.Fatal(err)
	}

	return scheduleConfig
}


func prompt(message string) string {
	fmt.Print(message)
	var input string
	fmt.Scanln(&input)
	return input
}

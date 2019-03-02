package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/eddiewebb/blync-studio-light/lights"
	"github.com/eddiewebb/blync-studio-light/calendars"
)


func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(updateCmd)
}

var calendarCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Interact with calendar (login)",
}

var updateCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Set light based on calendar events",
	Run: func(cmd *cobra.Command, args []string) {
		calendarId := viper.GetString("googleCalendar.calendarId")
		email := viper.GetString("googleCalendar.email")
		calendar := calendars.NewGoogleCalendarFromExistingToken()
		light.SetColor(calendar.GetColor(calendarId,email))
	},
		
}


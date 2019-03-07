package cmd

import (
	"github.com/eddiewebb/blync-studio-light/calendars"
	"github.com/eddiewebb/blync-studio-light/config"
	"github.com/eddiewebb/blync-studio-light/lights"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

func init() {
	rootCmd.AddCommand(refreshCmd)
	refreshCmd.AddCommand(calendarCmd)
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Interact with calendar (login)",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.InfoLevel)
			log.Info("Verbose logging enabled")
		}

		var schedule config.StudioLightSchedule
		if err := viper.UnmarshalKey("schedule", &schedule); err != nil {
			log.Info("no schedule in config")
		}

		today := time.Now().Weekday()
		offClock := minutesOfDay(schedule.OffHour, schedule.OffMinute)
		onClock := minutesOfDay(schedule.OnHour, schedule.OnMinute)
		nowHours, nowMinutes, _ := time.Now().Clock()
		now := minutesOfDay(nowHours, nowMinutes)
		log.Infof("Day: %s, TimeNow: %v:%v, onHour: %v:%02d, offHour: %v:%02d",today, nowHours, nowMinutes, schedule.OnHour, schedule.OnMinute, schedule.OffHour, schedule.OffMinute)
		if onClock <= now && now < offClock && ! schedule.DaysOffContains(today) {
			log.Info("Time is within schedule. Run `config schedule` to set/adjust off hours")
		} else {
			log.Warn("Time is outside configured schedule, lights off. Run `config schedule` to change off hours")
			light.Off()
			os.Exit(0)
		}
	},
}

func minutesOfDay(hour int, minutes int) int {
	return (hour * 60) + minutes
}


var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Set light based on calendar events",
	Run: func(cmd *cobra.Command, args []string) {
		calendarId := viper.GetString("googleCalendar.calendarId")
		email := viper.GetString("googleCalendar.email")

		log.Infof("Attempt to get calendar service", calendarId)
		if calendar,err := calendars.NewGoogleCalendar(); err != nil {
			log.Fatalf("Error accessing calendar %s: %v", calendarId,err)
		}else{			
			light.SetColor(calendar.GetColor(calendarId, email))
		}
	},
}

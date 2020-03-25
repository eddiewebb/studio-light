package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DaysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

type Configuration struct {
	GoogleCalendar GoogleCalendarConfiguration
	Schedule       StudioLightSchedule
	ConfigFile     string
}

// initConfig reads in config file and ENV variables if set.
func (c *Configuration) InitConfig() string {
	log.Infoln("Initalize Config")
	if c.ConfigFile != "" {
		log.Infoln("Config path provided as " + c.ConfigFile)
		// Use config file from the flag.
		viper.SetConfigFile(c.ConfigFile)
	} else {
		// Find home directory.
		log.Infoln("Loading default config file path")
		// Search config in home directory with name ".blync-studio-light" (without extension).
		viper.AddConfigPath(GetHomeDir())
		viper.SetConfigName(".blync-studio-light")
		c.ConfigFile = GetHomeDir() + "/.blync-studio-light.json"
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infoln("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Warnf("No configuration found at  %s, functionality will be limited until you run `config init\n", c.ConfigFile)
	}
	return c.ConfigFile
}

func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

type GoogleCalendarConfiguration struct {
	CalendarId string
	Email      string
}

type StudioLightSchedule struct {
	OnHour    int
	OnMinute  int
	OffHour   int
	OffMinute int
	DaysOff   []time.Weekday
}

func NewSchedule(on string, off string, days string) (s StudioLightSchedule, e error) {
	onHour, onMinutes, err := asHoursAndMinutes(on)
	if err != nil {
		e = err
	}
	offHour, offMinutes, err := asHoursAndMinutes(off)
	if err != nil {
		e = err
	}
	daysOff, err := parseWeekdays(days)
	if err != nil {
		e = err
	}

	s = StudioLightSchedule{
		OnHour:    onHour,
		OnMinute:  onMinutes,
		OffHour:   offHour,
		OffMinute: offMinutes,
		DaysOff:   daysOff,
	}
	return
}

func parseWeekdays(commaDays string) (days []time.Weekday, e error) {
	if commaDays == "" {
		return
	}
	for _, d := range strings.Split(commaDays, ",") {
		if day, ok := DaysOfWeek[strings.TrimSpace(d)]; ok {
			days = append(days, day)
		} else {
			return days, fmt.Errorf("The string %s is not a valid day from %v\n", d, DaysOfWeek)
		}
	}
	return
}

func asHoursAndMinutes(time string) (hours int, minutes int, e error) {
	if time == "" {
		return 0, 0, nil
	}
	parts := strings.Split(time, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 || hours > 23 {
		return 0, 0, fmt.Errorf("Not a valid number for time, please use HH:MM for hours with a 24 hour clock (0-23)")
	}
	minutes, merr := strconv.Atoi(parts[1])
	if merr != nil || minutes < 0 || minutes > 59 {
		return 0, 0, fmt.Errorf("Not a valid number for time, please use HH:MM for minutes as 0-59")
	}
	return
}

func (s *StudioLightSchedule) DaysOffContains(day time.Weekday) bool {
	for _, dayOff := range s.DaysOff {
		if dayOff == day {
			return true
		}
	}
	return false
}

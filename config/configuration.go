package config

import (

	"strconv"
	"strings"
)

type Configuration struct {
	GoogleCalendar GoogleCalendarConfiguration
	Schedule       StudioLightSchedule
}

type GoogleCalendarConfiguration struct {
	CalendarId string
	Email      string
}

type StudioLightSchedule struct {
	OffHour   int
	OffMinute int
	OnHour    int
	OnMinute  int
	DaysOff	  string
}

func (s *StudioLightSchedule) DaysOffContains(day int) bool{
    for _, d := range strings.Split(s.DaysOff,",") {
    	dayOff,_ := strconv.Atoi(d)
        if dayOff == day {
            return true
        }
    }
    return false
}
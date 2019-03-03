package config

import "testing"
import "gotest.tools/assert"
import "time"



func TestDayOfWeeksParsedCorrectly(t *testing.T){
	scheduleConfig,_ := NewSchedule("00:00","00:00","Tuesday,Wednesday")
	assert.Assert(t, scheduleConfig.DaysOffContains(DaysOfWeek["Tuesday"] ) )
	assert.Assert(t, scheduleConfig.DaysOffContains(time.Weekday(3)) )
}

func TestDayOfWeeksParsedCorrectlyWithSpaces(t *testing.T){
	scheduleConfig,_ := NewSchedule("00:00","00:00","Tuesday, Wednesday,Thursday")
	assert.Assert(t, scheduleConfig.DaysOffContains(DaysOfWeek["Tuesday"] ) )
	assert.Assert(t, scheduleConfig.DaysOffContains(time.Weekday(3)) )
	assert.Assert(t, scheduleConfig.DaysOffContains(time.Weekday(4)) )
}


func TestDayOfWeeksParsedCorrectlyNoMatch(t *testing.T){
	scheduleConfig,_ := NewSchedule("00:00","00:00","Tuesday,Wednesday")
	assert.Assert(t, ! scheduleConfig.DaysOffContains(DaysOfWeek["Monday"] ), "" )
}


func TestBadDayOfWeekThrowsError(t *testing.T){
	_,err := NewSchedule("00:00","00:00","MumpDay")
	assert.ErrorContains(t, err, "is not a valid day" )
}

func TestOnHoursParsedCorrectly(t *testing.T){
	scheduleConfig,_ := NewSchedule("18:15","","")
	assert.Equal(t, scheduleConfig.OnHour, 18 ) 
	assert.Equal(t, scheduleConfig.OnMinute, 15 ) 
	assert.Equal(t, scheduleConfig.OffHour, 0 ) 
	assert.Equal(t, scheduleConfig.OffMinute, 0 ) 
}

func TestOffHoursParsedCorrectly(t *testing.T){
	scheduleConfig,_ := NewSchedule("","23:41","")
	assert.Equal(t, scheduleConfig.OnHour, 0 ) 
	assert.Equal(t, scheduleConfig.OnMinute, 0 ) 
	assert.Equal(t, scheduleConfig.OffHour, 23 ) 
	assert.Equal(t, scheduleConfig.OffMinute, 41 ) 
}


func TestBadHourThrowsError(t *testing.T){
	_,err := NewSchedule("24:00","","")
	assert.ErrorContains(t, err, "Not a valid number for time, please use HH:MM for hours with a 24 hour clock (0-23)" )
}
func TestBadMinuteThrowsError(t *testing.T){
	_,err := NewSchedule("00:60","","")
	assert.ErrorContains(t, err, "Not a valid number for time, please use HH:MM for minutes as 0-59" )
}
func TestBadOffMinuteThrowsError(t *testing.T){
	_,err := NewSchedule("","00:60","")
	assert.ErrorContains(t, err, "Not a valid number for time, please use HH:MM for minutes as 0-59" )
}
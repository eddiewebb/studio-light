package calendars

import "testing"
import "gotest.tools/assert"
import 	"google.golang.org/api/calendar/v3"


//feel like there is a way to eliminate this for thestruct in the file under test, but getting type conflictc
type MockDAL struct{
	Items []*calendar.Event
}

func (m *MockDAL) ResolveEventList(calendarId string) []*calendar.Event {
	return m.Items
}


func getMockWithItems(items []*calendar.Event) (GoogleCalendar){
	dal := MockDAL{
		Items: items,
	}
	cal := GoogleCalendar{
		dal: &dal,
	}
	return cal
}


func TestGetColorIsGreenForEmptyEvents(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "green" )
}


func TestGetColorIsYellowForNewInvites(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{
				Email: userEmail,
				ResponseStatus: "needsAction",
			},
		},
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "yellow" )
}
func TestGetColorIsYellowForTentativeInvites(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{
				Email: userEmail,
				ResponseStatus: "tentative",
			},
		},
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "yellow" )
}

func TestGetColorIsRedForEventWithNoAttendees(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "red" )
}

func TestGetColorIsGreenForEventWithNoAttendeesButTransparent(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
		Transparency: "transparent",
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "green" )
}


func TestGetColorIsGreenForDeclinedInvites(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{
				Email: userEmail,
				ResponseStatus: "declined",
			},
		},
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "green" )
}

func TestGetColorIsGreenForTransparentEvents(t *testing.T){
	userEmail := "me@example.com"
	var mockItems []*calendar.Event
	event := calendar.Event{
		Status: "confirmed",
		Transparency: "transparent",
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{
				Email: userEmail,
				ResponseStatus: "needsAction",
			},
		},
	}
	mockItems = append(mockItems,&event)
	SUT := getMockWithItems(mockItems)
	color := SUT.GetColor("",userEmail)
	assert.Assert(t, color == "green" )
}

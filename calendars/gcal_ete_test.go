package calendars

import "testing"
import "gotest.tools/assert"
import 	"github.com/eddiewebb/blync-studio-light/config"

// This is ETE type test

//  Look for service mocking/stubs techniques to miic google API responses
func TestGetColorEte(t *testing.T){
	if testing.Short() {
		// -test.short
        t.Skip("skipping ete test in short mode.")
    }
	C := config.Configuration{
		ConfigFile: "./test/busy.json",
	}
	C.InitConfig()
	calendar,_ := NewGoogleCalendar()
	color := calendar.GetColor("n78adpbq0qvp59gqurtula39r4@group.calendar.google.com","ollitech@gmail.com")
	assert.Assert(t, color == "red" )
}


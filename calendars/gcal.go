package calendars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/golang/glog"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
)

const tokFile = "token.json"

type GoogleCalendar struct {
	service 		*gcal.Service
}


func NewGoogleCalendarFromExistingToken() (googleCalendar GoogleCalendar){
	googleCalendar.service = getCalendarService()
	return
}

func NewGoogleCalendarFromNewToken() (googleCalendar GoogleCalendar){
	generateAccessToken()
	googleCalendar.service = getCalendarService()
	return
}	


func (c *GoogleCalendar)  Verify(calendarId string) {
	fmt.Println("CalendarId: |" + calendarId + "|")
	cal, err := c.service.CalendarList.Get(calendarId).Do()
	if err != nil {
		glog.Fatalf("Unable to retrieve calendar: %v", err)
	}

	res, err := c.service.Events.List(cal.Id).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
	if err != nil {
		glog.Fatalf("Unable to retrieve calendar events list: %v", err)
	}
	for _, v := range res.Items {
		fmt.Printf("Calendar ID %q event: %v: %q\n", cal.Id, v.Updated, v.Summary)
	}
	fmt.Printf("Calendar cal.Id %q Summary: %v\n", cal.Id, res.Summary)
	fmt.Printf("Calendar cal.Id %q next page token: %v\n", cal.Id, res.NextPageToken)	
}

func (c *GoogleCalendar) GetColor(calendarId string, userEmail string) string{
	minTime := time.Now().Format(time.RFC3339)
	maxTime := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	events, err := c.service.Events.List(calendarId).ShowHiddenInvitations(true).TimeMin(minTime).TimeMax(maxTime).SingleEvents(true).Do()
	if err != nil {
		glog.Fatalf("Unable to retrieve list of events: %v", err)
	}

	color := "green"
	for _,event := range events.Items{
		//glog.Infof("%+v\n",event)
		glog.Infof("%s\n",event.Summary)
		// if calendar is only free/busy access, we can only use confirmed status. 
		// If we have full acces check addtional attributes
		if event.Transparency == "transparent" {
			//event is marked "free" in calendar, dont mark busy
			glog.Infoln("Event is marked free, ignore")
			continue
		}else{			
			status := event.Status	
			for _,attendee := range event.Attendees{
				if attendee.Email == userEmail{
					status = attendee.ResponseStatus
					fmt.Printf("Updated status to %s based on attendee %s\n",status, attendee.Email)
					break
				}
			}
			if status == "tentative" || status == "needsAction" {
				// event is worth noting, but could be one of many. Set to yellow and keep lookinh
				glog.Infof("Yellow\n")
				color = "yellow"
				continue
			}else if status == "confirmed" || status == "accepted" {
				// busy event, mark red and stop loooking
				color = "red"
				glog.Infof("Red\n")
				break
			}
		} 
	}	
	fmt.Println(color)
	return color
	
}

func generateAccessToken(){
	config := readConfig()
	tok := getTokenFromWeb(config)
	saveToken(tokFile, tok)
}

func getCalendarService() *gcal.Service {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		glog.Fatalln("No token.json found. Have you config init'd ?")
		os.Exit(1)
	}		
	config := readConfig()
	svc, err := gcal.New(config.Client(context.Background(), tok))
	if err != nil {
		glog.Fatalf("Unable to create Calendar service: %v", err)
	}
	return svc
}

func readConfig() *oauth2.Config{
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		glog.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gcal.CalendarReadonlyScope)
	if err != nil {
		glog.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		glog.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		glog.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		glog.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		glog.Fatalf("Encoding issue: %v", err)
	}
}

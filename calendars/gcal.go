package calendars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/eddiewebb/blync-studio-light/config"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
)

const configSuffix = "/.studio-light/gcal/" // where in users home directory which contains files below
const tokFile = "token.json"
const credsFile = "credentials.json"

type CalendarDAL interface {
	ResolveEventList(calendarId string) []*gcal.Event
}

// GoogleCalendarDAL implements interface CalendarDAL
type GoogleCalendarDAL struct {
	service *gcal.Service
}

type GoogleCalendar struct {
	dal CalendarDAL
}

func NewGoogleCalendar() (GoogleCalendar, error) {
	token, err := tokenFromFile(getFullPath(tokFile))
	if err != nil || token == nil {
		log.Warnf("No existing toke from %s, attempting to create new", tokFile)
		token = generateAccessToken()
		/*tok, err := tokenFromFile(token)
		if err != nil {
			log.Warn(err)
			log.Fatalln("No token.json found. Have you config init'd ?")
			os.Exit(1)
		}*/
	}
	log.Infof("token: %v", token)

	service, err := getCalendarService(token)
	if err != nil {
		log.Warnf("Failed ot get calendar")
		return GoogleCalendar{}, err
	}
	dal := GoogleCalendarDAL{
		service: service,
	}
	cal := GoogleCalendar{
		dal: &dal,
	}
	return cal, nil
}

func (c *GoogleCalendarDAL) ResolveEventList(calendarId string) []*gcal.Event {
	minTime := time.Now().Format(time.RFC3339)
	maxTime := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
	events, err := c.service.Events.List(calendarId).ShowHiddenInvitations(true).TimeMin(minTime).TimeMax(maxTime).SingleEvents(true).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of events: %v", err)
	}
	return events.Items
}

func (c *GoogleCalendar) GetColor(calendarId string, userEmail string) string {
	log.Infof("Checking Status")
	events := c.dal.ResolveEventList(calendarId)
	color := "green"
	if len(events) > 0 {
		for _, event := range events {
			//log.Infof("%+v\n",event)
			// if calendar is only free/busy access, we can only use confirmed status.
			// If we have full acces check addtional attributes
			if event.Transparency == "transparent" {
				//event is marked "free" in calendar, dont mark busy
				log.Infof("Event %s is marked free, ignore\n", event.Summary)
				continue
			} else {
				// Things we created have status confirmed, no attendees.
				status := event.Status
				// For invites sent to us, check out reply
				for _, attendee := range event.Attendees {
					if attendee.Email == userEmail {
						status = attendee.ResponseStatus
						log.Infof("Updated status to %s based on attendee %s\n", status, attendee.Email)
						break
					}
				}
				if status == "tentative" || status == "needsAction" {
					// event is worth noting, but could be one of many. Set to yellow and keep lookinh
					log.Infof("Event %s is marked %s\n", event.Summary, status)
					color = "yellow"
					continue
				} else if status == "confirmed" || status == "accepted" {
					// busy event, mark red and stop loooking
					log.Infof("Event %s is marked %s\n", event.Summary, status)
					color = "red"
					break
				}
			}
		}
	} else {
		log.Infof("No events currently scheduled.")
	}
	log.Infof("Setting light %s", color)
	return color
}

/*
*  Private Helpers
 */

func getCalendarService(token *oauth2.Token) (*gcal.Service, error) {
	config := readConfig()
	svc, err := gcal.New(config.Client(context.Background(), token))
	return svc, err
}

func generateAccessToken() *oauth2.Token {
	config := readConfig()
	tok := getTokenFromWeb(config)
	saveToken(getFullPath(tokFile), tok)
	return tok
}

func readConfig() *oauth2.Config {
	b, err := ioutil.ReadFile(getFullPath(credsFile))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gcal.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
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
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		log.Fatalf("Encoding issue: %v", err)
	}
}

func getFullPath(path string) string {
	return config.GetHomeDir() + configSuffix + path
}

func RemoveExistingGoogleAuthToken() {
	os.Remove(getFullPath(tokFile))
}

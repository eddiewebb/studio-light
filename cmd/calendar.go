package cmd

import (
	"encoding/json"
	"fmt"
	"flag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"github.com/golang/glog"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/eddiewebb/blync-studio-light/light"
)

var client *http.Client

func init() {
	calendarCmd.AddCommand(loginCmd)
	calendarCmd.AddCommand(logoutCmd)
	calendarCmd.AddCommand(verifyCmd)
	calendarCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(calendarCmd)
	flag.Parse()
}

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Interact with calendar (login)",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Do the oauth dance with Google",
	Long:  `Uses your credentials to get an access token for automated interactions`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			glog.Fatalf("Unable to read client secret file: %v", err)
		}
		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
		if err != nil {
			glog.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client = getClient(config)

	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout",
	Long:  `Forget you`,
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "COnfirm calendar status",
	Long:  `Uses your credentials to get an access token for automated interactions`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := getCalendarService()
		calendarId := viper.GetString("googleCalendar.calendarId")
		fmt.Println("CalendarId: |" + calendarId + "|")
		cal, err := svc.CalendarList.Get(calendarId).Do()
		if err != nil {
			glog.Fatalf("Unable to retrieve calendar: %v", err)
		}

		res, err := svc.Events.List(cal.Id).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
		if err != nil {
			glog.Fatalf("Unable to retrieve calendar events list: %v", err)
		}
		for _, v := range res.Items {
			fmt.Printf("Calendar ID %q event: %v: %q\n", cal.Id, v.Updated, v.Summary)
		}
		fmt.Printf("Calendar cal.Id %q Summary: %v\n", cal.Id, res.Summary)
		fmt.Printf("Calendar cal.Id %q next page token: %v\n", cal.Id, res.NextPageToken)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Set light based n calendar",
	Run: func(cmd *cobra.Command, args []string) {
		calendarId := viper.GetString("googleCalendar.calendarId")
		minTime := time.Now().Format(time.RFC3339)
		maxTime := time.Now().Add(time.Minute * 5).Format(time.RFC3339)
		svc := getCalendarService()
		events, err := svc.Events.List(calendarId).ShowHiddenInvitations(true).TimeMin(minTime).TimeMax(maxTime).SingleEvents(true).Do()
		if err != nil {
			glog.Fatalf("Unable to retrieve list of events: %v", err)
		}

		color := "green"
		for _,event := range events.Items{
			glog.Infof("%+v\n",event)

			// if calendar is only free/busy access, we can only use confirmed status. 
			// If we have full acces check addtional attributes
			if event.Transparency == "transparent" {
				//event is marked "free" in calendar, dont mark busy
				glog.Infoln("Event is marked free, ignore")
				continue
			}else{			
				status := event.Status	
				glog.Infof("Set status to %s based on event\n",status)
				for _,attendee := range event.Attendees{
					if attendee.Email == calendarId{
						status = attendee.ResponseStatus
						glog.Infof("Updated status to %s based on attendee %s\n",status, attendee.Email)
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

		light.SetColor(color)
		glog.Infoln("Updated!")
	},
}

func getCalendarService() *calendar.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		glog.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		glog.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	svc, err := calendar.New(client)
	if err != nil {
		glog.Fatalf("Unable to create Calendar service: %v", err)
	}
	return svc
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
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

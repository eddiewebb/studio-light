package cmd


import (
"fmt"
"encoding/json"
"io/ioutil"
        "log"
        "net/http"
        "os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	
        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/calendar/v3"
)

var client *http.Client

func init() {
	calendarCmd.AddCommand(loginCmd)
	calendarCmd.AddCommand(logoutCmd)
	calendarCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(calendarCmd)
}

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Interact with calendar (login)",
	
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Do the oauth dance with Google",
	Long: `Uses your credentials to get an access token for automated interactions`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client = getClient(config)

	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout",
	Long: `Forget you`,
	
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "COnfirm calendar status",
	Long: `Uses your credentials to get an access token for automated interactions`,
	Run: func(cmd *cobra.Command, args []string) {

        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }
		config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)
        svc, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to create Calendar service: %v", err)
		}
		listRes, err := svc.CalendarList.List().Fields("items/id").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve list of calendars: %v", err)
		}
		for _, v := range listRes.Items {
			log.Printf("Calendar ID: %v\n", v.Id)
		}

		calendarId := viper.GetString("googleCalendar.calendarId")
		fmt.Println("CalendarId: |" + calendarId +"|")
		cal,err := svc.CalendarList.Get(calendarId).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar: %v", err)
		}
		
		res, err := svc.Events.List(cal.Id).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar events list: %v", err)
		}
		for _, v := range res.Items {
			log.Printf("Calendar ID %q event: %v: %q\n", cal.Id, v.Updated, v.Summary)
		}
		log.Printf("Calendar cal.Id %q Summary: %v\n", cal.Id, res.Summary)
		log.Printf("Calendar cal.Id %q next page token: %v\n", cal.Id, res.NextPageToken)
	},
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
        json.NewEncoder(f).Encode(token)
}




package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"

	meet "github.com/ysugimoto/google-meet-api/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	tokenFile      = "token.json"
	credentialFile = "credentials.json"
)

// make sure you have created credential file.
// Or you can create via https://developers.google.com/calendar/quickstart/go quickly.
func main() {
	buf, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Failed to read credential file %v\n", err)
	}

	config, err := google.ConfigFromJSON(buf, meet.MeetScopes...)
	if err != nil {
		log.Fatalf("Failed to create oauth2 config %v\n", err)
	}

	token, err := getTokenFromFile(tokenFile)
	if err != nil {
		token, err = getTokenFromWeb(config, tokenFile)
		if err != nil {
			log.Fatalln("Failed to get oauth token", err)
		}
	}

	ctx := context.Background()
	client := config.Client(ctx, token)
	srv, err := meet.New(client)
	if err != nil {
		log.Fatalf("Failed to create meet service %v", err)
	}
	m, err := srv.Conference.Create().Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to create meet URL %v", err)
	}
	fmt.Printf("Google Meet URL has been generated: %s\n", m.URL)
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to open token file %v", err)
	}
	defer fp.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(fp).Decode(token); err != nil {
		return nil, fmt.Errorf("Failed to decode saved json file %v", err)
	}
	return token, nil
}

func getTokenFromWeb(config *oauth2.Config, file string) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-exchange", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web %v", err)
	}

	// save to local file
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open token file %v", err)
	}
	defer fp.Close()
	if err := json.NewEncoder(fp).Encode(token); err != nil {
		return nil, fmt.Errorf("Failed to encode token to json %v", err)
	}
	return token, nil
}

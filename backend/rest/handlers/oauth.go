package handlers

import (
	"context"
	"fmt"
	"net/http"

	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/boj/redistore.v1"
)

type Oauth2 struct {
	logger *log.Logger
	conf   *oauth2.Config
	Client *http.Client
	store  *redistore.RediStore
}

func NewOauth2(logger *log.Logger) *Oauth2 {
	oauth := new(Oauth2)
	oauth.logger = logger
	oauth.intializeConfiguration()
	return oauth
}

func (oauth *Oauth2) intializeConfiguration() {
	conf := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9090/callback-oauth",
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint:     google.Endpoint,
	}

	oauth.conf = conf
}

func (oauth *Oauth2) Authenticate(rw http.ResponseWriter, r *http.Request) {

	url := oauth.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	return

}

func (oauth *Oauth2) RedirectCallback(rw http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	oauth.logger.Printf("Code %v: \n", code)
	if code == "" {
		oauth.logger.Fatal("Code not found...")
		rw.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			rw.Write([]byte("User has denied Permission.."))
		}
	} else {
		ctx := context.Background()
		token, err := oauth.conf.Exchange(ctx, code)
		if err != nil {
			oauth.logger.Fatalf("Oauth Exchange Failed with %v", err)
		}

		client := oauth.conf.Client(ctx, token)
		oauth.Client = client

		// resp, err := client.Get("https://photoslibrary.googleapis.com/v1/albums")

		// if err != nil {
		// 	oauth.logger.Printf("Get: %v\n", err.Error())
		// 	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		// 	return
		// }
		// defer resp.Body.Close()

		// response, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	oauth.logger.Printf("ReadAll: %v\n", err.Error())
		// 	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		// 	return
		// }

		// oauth.logger.Printf("parseResponseBody: %v\n", string(response))

		// rw.Write([]byte("Hello, I'm protected\n"))
		// rw.Write([]byte(string(response)))
		// return
	}
}

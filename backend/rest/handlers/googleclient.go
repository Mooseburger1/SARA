package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/boj/redistore.v1"
)

type googleClientBuilder struct {
	logger *log.Logger
	conf   *oauth2.Config
	store  *redistore.RediStore
}

func NewGoogleClientBuilder() *googleClientBuilder {
	return &googleClientBuilder{
		logger: nil,
		conf:   nil,
		store:  nil,
	}
}

func (builder *googleClientBuilder) SetLogger(logger *log.Logger) *googleClientBuilder {
	builder.logger = logger
	return builder
}

func (builder *googleClientBuilder) SetStore(store *redistore.RediStore) *googleClientBuilder {
	builder.store = store
	return builder
}

func (builder *googleClientBuilder) SetConfig(conf *oauth2.Config) *googleClientBuilder {
	builder.conf = conf
	return builder
}

func (builder *googleClientBuilder) Build() *GoogleClient {
	gc := GoogleClient{logger: builder.logger,
		store: builder.store,
		conf:  builder.conf}
	return &gc
}

type GoogleClient struct {
	logger *log.Logger
	conf   *oauth2.Config
	store  *redistore.RediStore
}

func ConfigBuilder() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9090/oauth-callback",
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint:     google.Endpoint,
	}

	return conf
}

func (gc GoogleClient) GetConfiguration() *oauth2.Config {
	return gc.conf
}

// Authenticate routes user through Google's Oauth workflow. If the user has already
// Authenticated and Authorized the app, they will be redirected
// TODO IMPORVE THIS WORKFLOW - CHECKING FOR A SESSION IS NOT ENOUGH
// googleClient can get reinitialized with no authorizedClient but the
// redis DB can save the stored session causing of state of conflict. This should
// check for the presence of both an access token and the valid expiry time of the
// the token. Maybe experiment with rebuilding the authorizedClient with the persisted
// code. But I fear this code is no longer valid after expiry time. This should be fine though
// if expiry time has run out, need to redo the whole Oauth process.
func (gc GoogleClient) Authenticate(rw http.ResponseWriter, r *http.Request) {

	url := gc.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	templ := template.Must(template.New("Oauth_Redirect").Parse(`
	<h1>Visit the URL for the auth dialog: <a href={{.URL}}>Link</a></h1>`))

	err := templ.Execute(rw, struct{ URL string }{URL: url})

	if err != nil {
		http.Redirect(rw, r, "/", http.StatusBadRequest)
	}
	return

}

// RedirectCallback is the URL registered with Google API dashboard as the callback
// Handler after a user has performed OAuth. It will save all tokens from the OAuth
// Process within the session cookie and return to user for future use
func (gc GoogleClient) RedirectCallback(rw http.ResponseWriter, r *http.Request) {

	// Extract google code
	code := r.FormValue("code")
	gc.logger.Printf("Code %v: \n", code)
	if code == "" {
		gc.logger.Fatal("Code not found...")
		rw.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			rw.Write([]byte("User has denied Permission.."))
		}
	} else {

		// Utilize the code to generate an Acess Token
		ctx := context.Background()
		token, err := gc.conf.Exchange(ctx, code)
		if err != nil {
			gc.logger.Fatalf("Oauth Exchange Failed with %v\n", err)
		}

		// Generate a new session cookie
		session, err := gc.store.Get(r, "session-key")
		if err != nil {
			fmt.Printf("Error getting session: %v\n", err)
		}

		// save tokens in session cookie
		session.Values["access-token"] = token.AccessToken
		session.Values["token-type"] = token.TokenType
		session.Values["refresh-token"] = token.RefreshToken
		session.Values["oauth-code"] = code
		session.Values["expiry"] = token.Expiry.Format(time.RFC3339Nano)

		err = session.Save(r, rw)
		if err != nil {
			fmt.Printf("Error saving session & token: %v\n", err)
		}

		fmt.Print("session and token saved\n")

		return
	}
}

// ListAlbums utilizes photoslibrary googleapis to list all albums in the
// Google photos account.
func (gc GoogleClient) ListAlbums(rw http.ResponseWriter, r *http.Request, client *http.Client) {
	gc.logger.Printf("\n\nThe Client is: %v\n\n", client)
	req, err := http.NewRequest("GET", "https://photoslibrary.googleapis.com/v1/albums", nil)
	req.Header.Set("Accept", "application/json")

	// Use the client to make request ot Google Photos API for list albums
	resp, err := client.Do(req)
	if err != nil {
		gc.logger.Printf("Get: %v\n", err.Error())
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gc.logger.Printf("ReadAll: %v\n", err.Error())
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rw.Write([]byte(string(response)))
	return

}

func (gc GoogleClient) ListPicturesFromAlbum(rw http.ResponseWriter, r *http.Request, client *http.Client) {

	body := `{"albumId": "AAmUB7Udp11GyFjuG6dicqcuJoZpESzagCtF8mMbb9ekLoftkrmQ_yT2-wpMM2iUc93PS0VPozGd"}`
	req, err := http.NewRequest("POST", "https://photoslibrary.googleapis.com/v1/mediaItems:search", strings.NewReader(body))
	if err != nil {
		gc.logger.Fatalf("Failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	fmt.Printf("The response %v", resp)
	if err != nil {
		gc.logger.Printf("Get: %v\n", err.Error())
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gc.logger.Printf("ReadAll: %v\n", err.Error())
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rw.Write([]byte(string(response)))
	return

}

// OhNo is the default Redirect handler for when a user has done something stupid
func (gc GoogleClient) OhNo(rw http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.New("Oh-No").Parse(`
	<h1>OH NO</h1>`))

	_ = templ.Execute(rw, nil)

	return
}

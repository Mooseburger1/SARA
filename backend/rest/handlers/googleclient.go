package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/boj/redistore.v1"
)

type googleClientBuilder struct {
	logger *log.Logger
	conf   *oauth2.Config
	Client *http.Client
	store  *redistore.RediStore
}

func NewGoogleClientBuilder() *googleClientBuilder {
	return &googleClientBuilder{
		logger: nil,
		conf:   nil,
		Client: nil,
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

func (builder *googleClientBuilder) SetClient(client *http.Client) *googleClientBuilder {
	builder.Client = client
	return builder
}

func (builder *googleClientBuilder) Build() *googleClient {
	gc := googleClient{logger: builder.logger,
		store:  builder.store,
		conf:   builder.conf,
		Client: builder.Client}
	return &gc
}

type googleClient struct {
	logger *log.Logger
	conf   *oauth2.Config
	Client *http.Client
	store  *redistore.RediStore
}

func ConfigBuilder() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "682403741961-hrhr7ip4k16apjgcjonrs2ii8cofsgb6.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:9090/oauth-callback",
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint:     google.Endpoint,
	}

	return conf
}

func (gc *googleClient) Authenticate(rw http.ResponseWriter, r *http.Request) {

	url := gc.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	templ := template.Must(template.New("Oauth_Redirect").Parse(`
	<h1>Visit the URL for the auth dialog: <a href={{.URL}}>Link</a></h1>`))

	err := templ.Execute(rw, struct{ URL string }{URL: url})
	if err != nil {
		http.Redirect(rw, r, "/", http.StatusBadRequest)
	}
	return

}

func (gc *googleClient) RedirectCallback(rw http.ResponseWriter, r *http.Request) {

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
		ctx := context.Background()
		token, err := gc.conf.Exchange(ctx, code)
		if err != nil {
			gc.logger.Fatalf("Oauth Exchange Failed with %v\n", err)
		}

		session, err := gc.store.Get(r, "session-key")
		if err != nil {
			fmt.Printf("Error getting session: %v\n", err)
		}

		client := gc.conf.Client(ctx, token)
		gc.Client = client

		session.Values["access-token"] = token.AccessToken
		session.Values["token-type"] = token.TokenType
		session.Values["refresh-token"] = token.RefreshToken
		session.Values["expiry"] = token.Expiry.String()

		err = session.Save(r, rw)
		if err != nil {
			fmt.Printf("Error saving session & token: %v\n", err)
		}

		fmt.Print("session and token saved\n")

		return
	}
}

func (gc *googleClient) ListAlbums(rw http.ResponseWriter, r *http.Request) {

	resp, err := gc.Client.Get("https://photoslibrary.googleapis.com/v1/albums")
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

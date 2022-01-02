package middleware

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/boj/redistore.v1"
)

const (
	CLIENT_ID         = ""
	CLIENT_SECRET     = ""
	REDIRECT_URL      = "http://localhost:9090/oauth-callback"
	SESSION_KEY       = "session-key"
	ACCESS_TOKEN_KEY  = "access-token"
	REFRESH_TOKEN_KEY = "refresh-token"
	TOKEN_TYPE_KEY    = "token-type"
	EXPIRY_KEY        = "expiry"
	OAUTH_CODE_KEY    = "oauth-code"
)

var SCOPES = []string{"https://www.googleapis.com/auth/photoslibrary.readonly"}

type Middleware struct {
	logger *log.Logger
	store  *redistore.RediStore
	config *oauth2.Config
}

func NewMiddleWare(logger *log.Logger, store *redistore.RediStore) *Middleware {
	return &Middleware{
		logger: logger,
		store:  store,
		config: ConfigBuilder(),
	}
}

func ConfigBuilder() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		RedirectURL:  REDIRECT_URL,
		Scopes:       SCOPES,
		Endpoint:     google.Endpoint,
	}

	return conf
}

func (mw *Middleware) Authorized(handler func(http.ResponseWriter, *http.Request, *http.Client)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session, err := mw.store.Get(r, SESSION_KEY)
		if err != nil {
			mw.logger.Fatalf("Error retrieveing session cookie: %v", err)
		}

		accessToken := session.Values[ACCESS_TOKEN_KEY]

		if accessToken == nil {
			mw.Authenticate(rw, r)
			return
		}

		token := new(oauth2.Token)
		token.AccessToken = accessToken.(string)
		token.RefreshToken = session.Values[REFRESH_TOKEN_KEY].(string)
		token.TokenType = session.Values[TOKEN_TYPE_KEY].(string)
		token.Expiry, err = time.Parse(time.RFC3339Nano, session.Values[EXPIRY_KEY].(string))
		if err != nil {
			mw.logger.Fatalf("Error parsing time: %v", err)
			http.Redirect(rw, r, "/", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		client := mw.config.Client(ctx, token)

		handler(rw, r, client)

		return

	}
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
func (mw *Middleware) Authenticate(rw http.ResponseWriter, r *http.Request) {

	url := mw.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
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
func (mw *Middleware) RedirectCallback(rw http.ResponseWriter, r *http.Request) {

	// Extract google code
	code := r.FormValue("code")
	mw.logger.Printf("Code %v: \n", code)
	if code == "" {
		mw.logger.Fatal("Code not found...")
		rw.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			rw.Write([]byte("User has denied Permission.."))
		}
	} else {

		// Utilize the code to generate an Acess Token
		ctx := context.Background()
		token, err := mw.config.Exchange(ctx, code)
		if err != nil {
			mw.logger.Fatalf("Oauth Exchange Failed with %v\n", err)
		}

		// Generate a new session cookie
		session, err := mw.store.Get(r, "session-key")
		if err != nil {
			mw.logger.Fatalf("Error getting session: %v\n", err)
		}

		// save tokens in session cookie
		session.Values[ACCESS_TOKEN_KEY] = token.AccessToken
		session.Values[TOKEN_TYPE_KEY] = token.TokenType
		session.Values[REFRESH_TOKEN_KEY] = token.RefreshToken
		session.Values[OAUTH_CODE_KEY] = code
		session.Values[EXPIRY_KEY] = token.Expiry.Format(time.RFC3339Nano)

		err = session.Save(r, rw)
		if err != nil {
			mw.logger.Printf("Error saving session & token: %v\n", err)
		}
		return
	}
}

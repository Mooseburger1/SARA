package auth

import (
	backendConfig "backend/configuration"
	clientProto "backend/grpc/proto/api/client"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/boj/redistore.v1"
)

const (
	SESSION_KEY       = "session-key"
	ACCESS_TOKEN_KEY  = "access-token"
	REFRESH_TOKEN_KEY = "refresh-token"
	TOKEN_TYPE_KEY    = "token-type"
	EXPIRY_KEY        = "expiry"
	OAUTH_CODE_KEY    = "oauth-code"
)

//ClientHandlerFunc is an amended http.HandlerFunc that takes in the typical params of
//a http.HandlerFunc plus a *ClientInfo proto. It is the expected input handler for the
// Google Oauth Middleware.
type ClientHandlerFunc func(http.ResponseWriter, *http.Request, *clientProto.ClientInfo)

//AuthMiddleware manages the authentiation flow for any requests made to Google service
//endpoints. Any requests not made from a valid session will be redirected to Google's
//OAuth portal. Subsequent requests will be passed on to provided ClientHandlerFunc.
type AuthMiddleware struct {
	logger          *log.Logger
	store           *redistore.RediStore
	oauthconfig     *oauth2.Config
	backendconfig   *backendConfig.GOAuthConfig
	redirectHandler *ClientHandlerFunc
	responseWriter  *http.ResponseWriter
	request         *http.Request
}

//NewAuthMiddleware is a builder for the AuthMiddleware struct
func NewAuthMiddleware(logger *log.Logger, store *redistore.RediStore,
	internalconfig *backendConfig.GOAuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		logger:        logger,
		store:         store,
		oauthconfig:   ConfigBuilder(internalconfig),
		backendconfig: internalconfig,
	}
}

//ConfigBuilder receives server side configurations and builds expected Oauth
//proto needed for verified Google API services requests
func ConfigBuilder(internalConfig *backendConfig.GOAuthConfig) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     internalConfig.ClientID,
		ClientSecret: internalConfig.ClientSecret,
		RedirectURL:  internalConfig.RedirectUrl,
		Scopes:       internalConfig.Scopes,
		Endpoint:     google.Endpoint,
	}

	return conf
}

func (mw *AuthMiddleware) IsAuthorized(clientHandler ClientHandlerFunc) http.HandlerFunc {
	mw.clearRedirectHandler()
	return func(rw http.ResponseWriter, r *http.Request) {
		session, err := mw.store.Get(r, SESSION_KEY)
		if err != nil {
			mw.logger.Fatalf("Error retrieveing session cookie: %v", err)
		}

		accessToken := session.Values[ACCESS_TOKEN_KEY]

		mw.redirectHandler = &clientHandler
		mw.responseWriter = &rw
		mw.request = r

		if accessToken == nil {
			mw.Authenticate(rw, r)
			return
		}

		ts, _ := time.Parse(time.RFC3339Nano, session.Values[EXPIRY_KEY].(string))
		protoTime := timestamppb.New(ts)

		clientInfo := clientInfoBuilder(accessToken.(string),
			session.Values[REFRESH_TOKEN_KEY].(string),
			session.Values[TOKEN_TYPE_KEY].(string),
			/* expiry= */ protoTime,
			/* authMW= */ mw)

		clientHandler(rw, r, &clientInfo)
	}
}

func clientInfoBuilder(accessToken string, refreshToken string, tokenType string,
	expiry *timestamppb.Timestamp, authMW *AuthMiddleware) clientProto.ClientInfo {
	tokenInfo := clientProto.TokenInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		Expiry:       expiry}

	appCreds := clientProto.ApplicationCredentials{
		ClientId:     authMW.backendconfig.ClientID,
		ClientSecret: authMW.backendconfig.ClientSecret}

	scoping := clientProto.Scoping{
		Scopes: authMW.backendconfig.Scopes}

	url := clientProto.URL{
		RedirectUrl: authMW.backendconfig.RedirectUrl}

	return clientProto.ClientInfo{
		TokenInfo:      &tokenInfo,
		AppCredentials: &appCreds,
		AppScopes:      &scoping,
		Urls:           &url}

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
func (mw *AuthMiddleware) Authenticate(rw http.ResponseWriter, r *http.Request) {

	url := mw.oauthconfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	rw.WriteHeader(http.StatusUnauthorized)
	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(struct {
		Url string `json:"url"`
	}{Url: url})

	if err != nil {
		panic(err)
	}
}

// RedirectCallback is the URL registered with Google API dashboard as the callback
// Handler after a user has performed OAuth. It will save all tokens from the OAuth
// Process within the session cookie and return to user for future use
func (mw *AuthMiddleware) RedirectCallback(rw http.ResponseWriter, r *http.Request) {

	// Extract google code
	code := r.FormValue("code")

	if code == "" {
		mw.clearRedirectHandler()
		mw.logger.Fatal("Code not found...")
		rw.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			rw.Write([]byte("User has denied Permission.."))
		}
		return
	}
	// Utilize the code to generate an Acess Token
	ctx := context.Background()
	token, err := mw.oauthconfig.Exchange(ctx, code)
	if err != nil {
		mw.clearRedirectHandler()
		mw.logger.Fatalf("Oauth Exchange Failed with %v\n", err)
	}

	// Generate a new session cookie
	session, err := mw.store.Get(r, "session-key")
	if err != nil {
		mw.clearRedirectHandler()
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
		mw.clearRedirectHandler()
		mw.logger.Printf("Error saving session & token: %v\n", err)
	}

	mw.logger.Print("calling is Authorized")
	callback := mw.IsAuthorized(*mw.redirectHandler)
	callback(rw, r)
}

// clearRedirectHandler clears any stored clientHandler functions stored from previous
// isAuthorized invocations. This information is only needed in event a call was made
// from an unauthenticated client.
func (mw *AuthMiddleware) clearRedirectHandler() {
	mw.redirectHandler = nil
	mw.responseWriter = nil
	mw.request = nil
}
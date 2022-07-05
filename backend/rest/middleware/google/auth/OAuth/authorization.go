package auth

import (
	backendConfig "backend/configuration"
	config "backend/configuration"
	clientProto "backend/grpc/proto/api/client"
	common "backend/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

var INSTANCE *AuthMiddleware

//NewAuthMiddleware is a builder for the AuthMiddleware struct
func GetAuthMiddleware() *AuthMiddleware {
	logger := log.New(os.Stdout, "authorization-middleware", log.LstdFlags)

	if INSTANCE == nil {
		config := config.NewGOAuthConfig()
		INSTANCE = &AuthMiddleware{
			logger:        logger,
			store:         common.GetDefaultRedisInstance(),
			oauthconfig:   ConfigBuilder(config),
			backendconfig: config,
		}
		return INSTANCE
	}
	return INSTANCE
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

// IsAuthorized checks the cookie associated with current client call. If that session has a valid
// OAuth token, the caller is already OAuth verified and the request is propagated forward to the
// ClientHandlerFunc. Otherwise, the ClientHandlerFunc is cached, and the request is propageted through the OAuth flow.
// Upon successful completion of the OAuth flow, the cached ClientHandlerFunc will be inoked on Google OAuth redirect.
func (mw *AuthMiddleware) IsAuthorized(clientHandler ClientHandlerFunc) http.HandlerFunc {

	// clear any cached ClientHandlerFunc
	mw.clearRedirectHandler()

	// return new http.Handlerfunc
	return func(rw http.ResponseWriter, r *http.Request) {

		// Check redis cache for stored session from request cookie
		session, err := mw.store.Get(r, SESSION_KEY)
		if err != nil {
			mw.logger.Fatalf("Error retrieveing session cookie: %v", err)
		}

		accessToken := session.Values[ACCESS_TOKEN_KEY]

		// Cache the ClientHandlerFunc and the current requests
		// response writer & request
		mw.redirectHandler = &clientHandler
		mw.responseWriter = &rw
		mw.request = r

		// If there's no Token, redirect user through OAuth flow
		if accessToken == nil {

			mw.logger.Print("CALLING AUTHENTICATE!!!!!!!!!!")
			mw.Authenticate(rw, r)
			return
		}

		// Parse Client Info and pass request onto ClientHandlerFunc
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

// clientInfoBuilder takes in all the stored params of the authorized caller and
// constructs a ClientInfo proto needed to be propagated along to the ClientHandlerFunc
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

// Authenticate constructs the URL to route the caller through Google's Oauth workflow.
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
// Process within the session cookie and then used the cached ClientHandlerFunc,
// responseWriter and request to continue the original call of the calling client.
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

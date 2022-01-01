package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"rest/handlers"

	"golang.org/x/oauth2"
	"gopkg.in/boj/redistore.v1"
)

type Middleware struct {
	logger *log.Logger
	store  *redistore.RediStore
	gc     *handlers.GoogleClient
}

func NewMiddleWare(logger *log.Logger, store *redistore.RediStore, gc *handlers.GoogleClient) *Middleware {
	return &Middleware{
		logger: logger,
		store:  store,
		gc:     gc,
	}
}

func (mw *Middleware) Authorized(handler func(http.ResponseWriter, *http.Request, *http.Client)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session, err := mw.store.Get(r, "session-key")
		if err != nil {
			mw.logger.Fatalf("Error retrieveing session cookie: %v", err)
		}

		accessToken := session.Values["access-token"]

		if accessToken == nil {

			mw.gc.Authenticate(rw, r)
			return
		}

		token := new(oauth2.Token)
		token.AccessToken = accessToken.(string)
		token.RefreshToken = session.Values["refresh-token"].(string)
		token.TokenType = session.Values["token-type"].(string)
		token.Expiry, err = time.Parse(time.RFC3339Nano, session.Values["expiry"].(string))
		if err != nil {
			mw.logger.Fatalf("Error parsing time: %v", err)
			http.Redirect(rw, r, "/", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		config := mw.gc.GetConfiguration()
		client := config.Client(ctx, token)

		mw.logger.Println("Serving handler func")
		handler(rw, r, client)

		return

	}
}

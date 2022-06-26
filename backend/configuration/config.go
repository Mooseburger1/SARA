package configuration

import "os"

type GOAuthConfig struct {
	RedirectUrl  string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func NewGOAuthConfig() *GOAuthConfig {
	return &GOAuthConfig{
		RedirectUrl:  "http://localhost:9090/oauth-callback",
		ClientID:     os.Getenv("GOOGLE_API_ID"),
		ClientSecret: os.Getenv("GOOGLE_API_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
	}
}

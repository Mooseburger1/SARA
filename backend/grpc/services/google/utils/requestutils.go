package utils

import (
	"backend/grpc/proto/api/client"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// createClient is a package private function utilized
// to create an http client that has Google API
// oauth2 credentials bounded to it. It is utilized
// to make oauth2 verified REST requests to the Google
// Photos API server
func CreateClient(info *client.ClientInfo) (*http.Client, error) {
	token := new(oauth2.Token)
	token.AccessToken = info.GetTokenInfo().GetAccessToken()
	token.RefreshToken = info.GetTokenInfo().GetRefreshToken()
	token.TokenType = info.GetTokenInfo().GetTokenType()
	token.Expiry = info.GetTokenInfo().GetExpiry().AsTime()

	ctx := context.Background()
	client := configBuilder(info).Client(ctx, token)

	return client, nil
}

// configBuilder configures the server with the
// application registered credentials on Google's
// API developers dashboard.
func configBuilder(info *client.ClientInfo) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     info.GetAppCredentials().GetClientId(),
		ClientSecret: info.GetAppCredentials().GetClientSecret(),
		RedirectURL:  info.GetUrls().GetRedirectUrl(),
		Scopes:       info.GetAppScopes().GetScopes(),
		Endpoint:     google.Endpoint,
	}
}

func CreateClientCreationError(err error) *status.Status {
	st := status.New(codes.InvalidArgument, "Client creation error")
	desc := fmt.Sprintf("Error creating client for making REST calls to Google Photos RESTServer: %s", err)
	v := &errdetails.ErrorInfo{Reason: desc}
	st, err = st.WithDetails(v)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st
}

func CreateErrorResponseError(statusCode int, response []byte) *status.Status {
	var rpcErrCode codes.Code

	var desc errResponse

	json.Unmarshal(response, &desc)
	switch statusCode {
	case 400:
		rpcErrCode = codes.InvalidArgument
	default:
		rpcErrCode = codes.InvalidArgument
	}

	st := status.New(rpcErrCode, desc.Error.Message)

	return st
}

type errResponse struct {
	Error errDetails `json:"error"`
}

type errDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

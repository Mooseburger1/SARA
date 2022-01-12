package services

import (
	"backend/grpc/proto/api/POGO"
	"backend/grpc/proto/api/client"
	"backend/grpc/proto/api/photos"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
	PHOTOS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/mediaItems:search"
	GET             = "GET"
	POST            = "POST"
)

// listAlbums is a package private function utilized to make an
// http request to the google photos API server. The response
// is unmarshalled and converted into an AlbumsInfo protobuf
func listAlbums(info *client.ClientInfo, logger *log.Logger) *photos.AlbumsInfo {
	client, err := createClient(info)
	if err != nil {
		logger.Printf("Error creating client: %v", err)
	}

	req, err := http.NewRequest(GET, ALBUMS_ENDPOINT, nil)
	if err != nil {
		logger.Printf("Error creating new request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic("Didn't get code 200")
	}

	defer resp.Body.Close()

	result := albumListDecoder(resp.Body)

	return albumsPogo2Proto(&result)

}

// albumListDecoder takes in the response body from the
// Google Photos API server for listing albums. It unmarshals
// the JSON response into an AlbumsInfoPOGO struct. It utilized
// solely by the listAlbums function
func albumListDecoder(body io.ReadCloser) POGO.AlbumsInfoPOGO {
	var result POGO.AlbumsInfoPOGO
	json.NewDecoder(body).Decode(&result)
	return result
}

// albumsPogo2Proto converts an AlbumsInfoPOGO (plain old golang
// object) into a AlbumsInfo protobuf object
func albumsPogo2Proto(result *POGO.AlbumsInfoPOGO) *photos.AlbumsInfo {
	var slices []*photos.AlbumInfo
	for _, info := range result.AlbumsInfo {

		var count int32
		var countint int
		var err error
		if info.MediaItemsCount != "" {
			countint, err = strconv.Atoi(info.MediaItemsCount)
			if err != nil {
				panic(err)
			}
			count = int32(countint)
		} else {
			count = 0
		}
		slices = append(slices,
			&photos.AlbumInfo{
				Id:                    info.Id,
				Title:                 info.Title,
				ProductUrl:            info.ProductUrl,
				MediaItemsCount:       int32(count),
				CoverPhotoBaseUrl:     info.CoverPhotoBaseUrl,
				CoverPhotoMediaItemId: info.CoverPhotoMediaItemId})
	}

	return &photos.AlbumsInfo{AlbumsInfo: slices}
}

// createClient is a package private function utilized
// to create an http client that has Google API
// oauth2 credentials bounded to it. It is utilized
// to make oauth2 verified REST requests to the Google
// Photos API server
func createClient(info *client.ClientInfo) (*http.Client, error) {
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

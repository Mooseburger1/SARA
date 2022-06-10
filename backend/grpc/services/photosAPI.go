package services

import (
	"backend/grpc/proto/api/POGO"
	"backend/grpc/proto/api/client"
	"backend/grpc/proto/api/photos"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
	PHOTOS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/mediaItems:search"
	GET             = "GET"
	POST            = "POST"
)

type PhotosAPICallError struct {
	Response *http.Response
}

func (pe *PhotosAPICallError) Error() string {
	return "Did Not Get HTTP Status 200 OK"
}

// listPhotosFromAlbum is a package private funciton utilized to make
// an http request to the google photos API server, specifically the endpoint
// which returns all photos for a specified album. The response is unmarshaled
// and converted into an PhotosInfo protobuf
func listPhotosFromAlbum(rpc *photos.FromAlbumRequest, logger *log.Logger) (*photos.PhotosInfo, error) {
	client, err := createClient(rpc.GetClientInfo())
	if err != nil {
		logger.Printf("Error creating client: %v", err)
		st := createClientCreationError(err)
		return nil, st.Err()
	}

	requestBody := []byte(fmt.Sprintf(`{"albumId":"%v", "pageSize":"%v", "pageToken":"%v"}`, rpc.GetAlbumId(), rpc.GetPageSize(), rpc.GetPageToken()))

	req, err := http.NewRequest(POST, PHOTOS_ENDPOINT, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Printf("Call to List Photos From Album returned status code %v, not %v", resp.StatusCode, http.StatusOK)
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}
		st := createErrorResponseError(resp.StatusCode, bodyBytes)
		return &photos.PhotosInfo{}, st.Err()
	}
	result := photosFromAlbumDecoder(resp.Body)

	return photosPogo2Proto(&result), nil

}

// listAlbums is a package private function utilized to make an
// http request to the google photos API server. The response
// is unmarshalled and converted into an AlbumsInfo protobuf
func listAlbums(rpc *photos.AlbumListRequest, logger *log.Logger) (*photos.AlbumsInfo, error) {
	client, err := createClient(rpc.GetClientInfo())
	if err != nil {
		logger.Printf("Error creating client: %v", err)
		st := createClientCreationError(err)
		return nil, st.Err()
	}

	req, err := http.NewRequest(GET, ALBUMS_ENDPOINT, nil)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	query.Add("pageToken", rpc.GetPageToken())
	query.Add("pageSize", strconv.Itoa(int(rpc.GetPageSize())))
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		logger.Printf("Call to List Albums returned status code %v, not %v", resp.StatusCode, http.StatusOK)
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}
		st := createErrorResponseError(resp.StatusCode, bodyBytes)
		return nil, st.Err()
	}

	result := albumListDecoder(resp.Body)

	return albumsPogo2Proto(&result), nil

}

func createClientCreationError(err error) *status.Status {
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

func createErrorResponseError(statusCode int, response []byte) *status.Status {
	desc := fmt.Sprintf("Google Photos REST server responded with a non 200 error code: %d: %s", statusCode, response)
	st := status.New(codes.InvalidArgument, desc)
	v := &errdetails.ErrorInfo{Reason: desc}
	st, err := st.WithDetails(v)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st
}

func photosFromAlbumDecoder(body io.ReadCloser) POGO.PhotosInfoPOGO {
	var result POGO.PhotosInfoPOGO
	json.NewDecoder(body).Decode(&result)
	return result
}

// albumListDecoder takes in the response body from the
// Google Photos API server for listing albums. It unmarshals
// the JSON response into an AlbumsInfoPOGO struct. It is utilized
// solely by the listAlbums function
func albumListDecoder(body io.ReadCloser) POGO.AlbumsInfoPOGO {
	var result POGO.AlbumsInfoPOGO
	json.NewDecoder(body).Decode(&result)
	return result
}

func photosPogo2Proto(result *POGO.PhotosInfoPOGO) *photos.PhotosInfo {
	var slices []*photos.PhotoInfo
	for _, info := range result.MediaItems {
		slices = append(slices,
			&photos.PhotoInfo{
				Id:         info.Id,
				ProductUrl: info.ProductUrl,
				BaseUrl:    info.BaseUrl,
				MimeType:   info.MimeType,
				Filename:   info.Filename,
			})
	}

	return &photos.PhotosInfo{PhotosInfo: slices, NextPageToken: result.NextPageToken}
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

	return &photos.AlbumsInfo{AlbumsInfo: slices, NextPageToken: result.NextPageToken}
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

package photos

import (
	"backend/grpc/proto/api/POGO"
	"backend/grpc/proto/api/photos"
	"backend/grpc/services/google/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
	PHOTOS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/mediaItems:search"
	GET             = "GET"
	POST            = "POST"
)

// listPhotosFromAlbum is a package private funciton utilized to make
// an http request to the google photos API server, specifically the endpoint
// which returns all photos for a specified album. The response is unmarshaled
// and converted into an PhotosInfo protobuf
func listPhotosFromAlbum(rpc *photos.FromAlbumRequest, logger *log.Logger) (*photos.PhotosInfo, error) {
	client, err := utils.CreateClient(rpc.GetClientInfo())
	if err != nil {
		logger.Printf("Error creating client: %v", err)
		st := utils.CreateClientCreationError(err)
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
		st := utils.CreateErrorResponseError(resp.StatusCode, bodyBytes)
		return &photos.PhotosInfo{}, st.Err()
	}
	result := photosFromAlbumDecoder(resp.Body)

	return photosPogo2Proto(&result), nil

}

// listAlbums is a package private function utilized to make an
// http request to the google photos API server. The response
// is unmarshalled and converted into an AlbumsInfo protobuf
func listAlbums(rpc *photos.AlbumListRequest, logger *log.Logger) (*photos.AlbumsInfo, error) {
	client, err := utils.CreateClient(rpc.GetClientInfo())
	if err != nil {
		logger.Printf("Error creating client: %v", err)
		st := utils.CreateClientCreationError(err)
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
		st := utils.CreateErrorResponseError(resp.StatusCode, bodyBytes)
		return nil, st.Err()
	}

	result := albumListDecoder(resp.Body)

	return albumsPogo2Proto(&result), nil

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

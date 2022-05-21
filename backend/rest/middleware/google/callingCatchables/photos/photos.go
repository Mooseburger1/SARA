package callingCatchables

import (
	clientProto "backend/grpc/proto/api/client"
	photosProto "backend/grpc/proto/api/photos"
	auth "backend/rest/middleware/google/auth/OAuth"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// photosRpcCaller is the client responsible for making calls
// to the gRPC server for the Google Photos endpoints. Successful
// calls will be propogated to the injected handlers. Failed RPC
// calls will be caught and handled gracefully.
type photosRpcCaller struct {
	logger       *log.Logger
	photosClient *photosProto.GooglePhotoServiceClient
}

// AlbumsHandlerFunc is a http.HandlerFunc extended to handle the successful request
// to the gRPC server for Google Photos and specifically for the AlmbusInfo endpoint.
type AlbumsHandlerFunc func(http.ResponseWriter, *http.Request, *photosProto.AlbumsInfo)

// PhotosInfoHandlerFunc is a http.HandlerFunc extended to handle the successful request
// to the gRPC server for Google Photos and specifically for the PhotosInfo endpoint.
type PhotosInfoHandlerFunc func(http.ResponseWriter, *http.Request, *photosProto.PhotosInfo)

// NewPhotosRpcCaller is a builder for a photosRpcCaller client. Will create a new instance
// with each invocation. Does not follow the singleton pattern.
func NewPhotosRpcCaller(logger *log.Logger, pc *photosProto.GooglePhotoServiceClient) *photosRpcCaller {
	return &photosRpcCaller{
		logger:       logger,
		photosClient: pc,
	}
}

func (rpc *photosRpcCaller) CatchableListAlbums(handler AlbumsHandlerFunc) auth.ClientHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		listRequest := makeAlbumListRequest(r, clientInfo)
		pc := *rpc.photosClient
		albums, err := pc.ListAlbums(context.Background(), listRequest)
		if err != nil {
			panic(err)
		}
		if albums.FailedRequest != nil {
			rw.Write([]byte(albums.FailedRequest))
			return
		}
		handler(rw, r, albums)

	}
}

func (rpc *photosRpcCaller) CatchablePhotosFromAlbum(handler PhotosInfoHandlerFunc) auth.ClientHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		photoRequest := makePhotosFromAlbumRequest(r, clientInfo)
		pc := *rpc.photosClient
		photos, err := pc.ListPhotosFromAlbum(context.Background(), photoRequest)
		if err != nil {
			panic(err)
		}

		if photos.FailedRequest != nil {
			rw.Write([]byte(photos.FailedRequest))
			return
		}
		handler(rw, r, photos)
	}
}

// makeAlbumListRequest is a package private helper function
// utilized to extract variables from the API URL and generate
// an AlbumListRequest proto. More specifically, it is a parser
// for the REST endpoint of list-albums and constructs the necessary
// RPC proto.
func makeAlbumListRequest(r *http.Request, ci *clientProto.ClientInfo) *photosProto.AlbumListRequest {

	pageToken := r.URL.Query().Get("pageToken")
	pageSize := r.URL.Query().Get("pageSize")

	var req photosProto.AlbumListRequest
	req.ClientInfo = ci

	// Parse the pageSize Url variable
	if pageSize != "" {
		i, err := str2Int32(pageSize)
		if err != nil {
			panic(err)
		}
		req.PageSize = i
	}

	// Parse the pageToken URL variable
	if pageToken != "" {
		req.PageToken = pageToken
	}

	return &req
}

// str2Int32 is a package private helper function
// for type conversion
func str2Int32(val string) (int32, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

func makePhotosFromAlbumRequest(r *http.Request, ci *clientProto.ClientInfo) *photosProto.FromAlbumRequest {
	vars := mux.Vars(r)
	albumId := vars["albumId"]
	pageToken := r.URL.Query().Get("pageToken")
	pageSize := r.URL.Query().Get("pageSize")

	var req photosProto.FromAlbumRequest
	req.ClientInfo = ci
	req.AlbumId = albumId

	// Parse the pageSize Url variable
	if pageSize != "" {
		i, err := str2Int32(pageSize)
		if err != nil {
			panic(err)
		}
		req.PageSize = i
	}

	// Parse the pageToken URL variable
	if pageToken != "" {
		req.PageToken = pageToken
	}
	return &req
}

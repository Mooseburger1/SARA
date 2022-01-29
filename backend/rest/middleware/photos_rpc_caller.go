package middleware

import (
	clientProto "backend/grpc/proto/api/client"
	photosProto "backend/grpc/proto/api/photos"
	"context"
	"log"
	"net/http"
	"strconv"
)

type rpcCaller func(http.ResponseWriter, *http.Request, *clientProto.ClientInfo)

type photosRpcCaller struct {
	logger       *log.Logger
	photosClient *photosProto.GooglePhotoServiceClient
}

func NewPhotosRpcCaller(logger *log.Logger, pc *photosProto.GooglePhotoServiceClient) *photosRpcCaller {
	return &photosRpcCaller{
		logger:       logger,
		photosClient: pc,
	}
}

func (rpc *photosRpcCaller) ListAlbumsCallWithError(handler func(http.ResponseWriter, *http.Request, *photosProto.AlbumsInfo)) rpcCaller {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		listRequest := makeAlbumListRequest(r, clientInfo)
		pc := *rpc.photosClient
		albums, err := pc.ListAlbums(context.Background(), listRequest)
		rpc.logger.Printf("logging from rpc caller: %v", albums)
		if err != nil {
			panic(err)
		}
		if albums.FailedRequest != nil {
			rw.Write([]byte(albums.FailedRequest))
			return
		}
		rpc.logger.Println("Called into the proper middleware")
		handler(rw, r, albums)

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

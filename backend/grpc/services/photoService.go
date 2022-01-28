package services

import (
	"backend/grpc/proto/api/photos"
	"context"
	"log"
)

// GphotoServer is the implementation of the
// Google Photo RPC server. It implements the
// * ListAlbums service
// * ListPhotosFromAlbum service
type GphotoServer struct {
	logger *log.Logger
}

// Constructor for instantiating a GphotoServer
func NewGphotoServer(logger *log.Logger) *GphotoServer {
	return &GphotoServer{logger: logger}
}

// ListAlbums is a RPC service endpoint. It receives an AlbumListRequest
// proto and returns an AlbumsInfo proto. Internally it makes an Oauth2
// authorized REST request to the Google Photos API server for listing albums.
func (g *GphotoServer) ListAlbums(ctx context.Context, rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return listAlbums(rpc.GetClientInfo(), g.logger), nil

}

// ListPhotosFromAlbum is a RPC service endpoint. It receives a
// FromAlbumRequest proto and returns a PhotosInfo proto. Internally
// it makes an Oauth2 authorized rest request to the Google Photos API
// server for listing photos from a specific album
func (g *GphotoServer) ListPhotosFromAlbum(ctx context.Context, rpc *photos.FromAlbumRequest) (*photos.PhotosInfo, error) {
	return listPhotosFromAlbum(rpc, g.logger), nil
}

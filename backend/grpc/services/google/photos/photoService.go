package photos

import (
	"backend/grpc/proto/api/photos"
	"context"
	"log"
)

// GphotoStub is the implementation of the
// Google Photo RPC server. It implements the
// * ListAlbums service
// * ListPhotosFromAlbum service
type GphotoStub struct {
	logger *log.Logger
}

// Constructor for instantiating a GphotoStub
func NewGphotoStub(logger *log.Logger) *GphotoStub {
	return &GphotoStub{logger: logger}
}

// ListAlbums is a RPC service endpoint. It receives an AlbumListRequest
// proto and returns an AlbumsInfo proto. Internally it makes an Oauth2
// authorized REST request to the Google Photos API server for listing albums.
func (g *GphotoStub) ListAlbums(ctx context.Context,
	rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return listAlbums(rpc, g.logger)

}

// ListPhotosFromAlbum is a RPC service endpoint. It receives a
// FromAlbumRequest proto and returns a PhotosInfo proto. Internally
// it makes an Oauth2 authorized rest request to the Google Photos API
// server for listing photos from a specific album
func (g *GphotoStub) ListPhotosFromAlbum(ctx context.Context,
	rpc *photos.FromAlbumRequest) (*photos.PhotosInfo, error) {
	return listPhotosFromAlbum(rpc, g.logger)
}

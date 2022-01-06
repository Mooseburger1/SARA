package rpc_servers

import (
	"backend/grpc/proto/api/photos"
	"context"
	"fmt"
	"log"
)

type GphotoServer struct {
	logger *log.Logger
}

func NewGphotoServer(logger *log.Logger) *GphotoServer {
	return &GphotoServer{logger: logger}
}

func (g *GphotoServer) ListAlbums(ctx context.Context, req *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	client_info := req.GetClientInfo()
	page_size := req.GetPageSize()
	page_token := req.GetPageToken()
	//TODO
	fmt.Print(client_info, page_size, page_token)
	return &photos.AlbumsInfo{AlbumsInfo: nil}, nil
}

func (g *GphotoServer) ListPhotosFromAlbum(ctx context.Context, album_info *photos.FromAlbumInfo) (*photos.PhotosInfo, error) {
	album_id := album_info.GetAlbumId()
	page_token := album_info.GetPageToken()
	page_size := album_info.GetPageSize()
	//TODO
	fmt.Print(album_id, page_token, page_size)
	return &photos.PhotosInfo{}, nil
}

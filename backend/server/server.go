package server

import (
	"context"

	"github.com/Mooseburger1/SARA/backend/s3"
)

type Server struct {
}

func (*Server) PictureLinks(ctx context.Context, request *PictureLinksRequest) (*PictureLinksResponse, error) {
	s := s3.S3{}
	s.CreateSession()

	bucketObjectChannel := make(chan *string)
	objectURLChannel := make(chan []string)

	go s.ListPictures(bucketObjectChannel)
	go s.GenerateObjectURL(bucketObjectChannel, objectURLChannel)

	var results []string

	results = <-objectURLChannel

	response := &PictureLinksResponse{
		Urls: results,
	}

	return response, nil

}

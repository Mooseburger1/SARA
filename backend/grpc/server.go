package main

import (
	protoPhotos "backend/grpc/proto/api/photos"
	rpc "backend/grpc/services"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := log.New(os.Stdout, "rpc-server", log.LstdFlags)

	grpcServer := grpc.NewServer()
	photoServer := rpc.NewGphotoServer(logger)

	protoPhotos.RegisterGooglePhotoServiceServer(grpcServer, photoServer)

	reflection.Register(grpcServer)
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	grpcServer.Serve(l)
}

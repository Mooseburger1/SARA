package main

import (
	protoPhotos "backend/grpc/proto/api/photos"
	photos "backend/grpc/services/google/photos"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type photosGRPCServer struct {
	Server *grpc.Server
	Logger *log.Logger
}

func NewPhotosGRPCServer() *photosGRPCServer {
	ps := photosGRPCServer{}
	ps.initServer()
	return &ps
}

func (ps *photosGRPCServer) initServer() {
	logger := log.New(os.Stdout, "photos-rpc-server", log.LstdFlags)
	ps.Logger = logger
	grpcServer := grpc.NewServer()
	photoServer := photos.NewGphotoServer(logger)

	protoPhotos.RegisterGooglePhotoServiceServer(grpcServer, photoServer)
	ps.Server = grpcServer
}

func (ps *photosGRPCServer) StartServer() {
	reflection.Register(ps.Server)
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		ps.Logger.Fatal(err)
		os.Exit(1)
	}
	ps.Logger.Printf("Photos grpc listening on 9091")
	ps.Server.Serve(l)
}

func (ps *photosGRPCServer) ShutdownServer() {
	ps.Server.GracefulStop()
}
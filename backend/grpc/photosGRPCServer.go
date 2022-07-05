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

// NewPhotosGRPCServer creates a new instance of a photosGRPCServer.
// It intializes it first and returns to caller, ready to be started.
func NewPhotosGRPCServer() *photosGRPCServer {
	ps := photosGRPCServer{}
	ps.initServer()
	return &ps
}

func (ps *photosGRPCServer) initServer() {
	logger := log.New(os.Stdout, "photos-rpc-server", log.LstdFlags)
	ps.Logger = logger
	grpcServer := grpc.NewServer()
	photoServer := photos.NewGphotoStub(logger)

	protoPhotos.RegisterGooglePhotoServiceServer(grpcServer, photoServer)
	ps.Server = grpcServer
}

// StartServer will start the intialized photosGRPCServer. It listens
// on port :9091
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

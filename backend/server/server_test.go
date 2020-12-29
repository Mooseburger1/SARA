package server_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/Mooseburger1/SARA/backend/server"
	"google.golang.org/grpc"
)

func test_server() {
	address := "0.0.0.0:2001"
	lis, err := net.Listen("tcp", address)

	if err != nil {
		panic(fmt.Sprintf("Failed to create test gRPC server"))
	}

	s := grpc.NewServer()
	server.RegisterPictureLinksServiceServer(s, &server.Server{})

	s.Serve(lis)

}

func TestPictureLinks(t *testing.T) {

	go test_server()

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:2001", opts)
	if err != nil {
		t.Errorf("Failed to dial gRPC server: %v", err)
	}

	defer cc.Close()

	client := server.NewPictureLinksServiceClient(cc)
	request := &server.PictureLinksRequest{}

	resp, err := client.PictureLinks(context.Background(), request)

	if err != nil {
		t.Errorf("Failed to make PictureLinks request: %v", err)
	}

	t.Log(resp)

}

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Mooseburger1/SARA/backend/server"
	"google.golang.org/grpc"
)

func main() {
	address := "0.0.0.0:2001"

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	fmt.Printf("Sever is listening on %v ...", address)

	s := grpc.NewServer()

	server.RegisterPictureLinksServiceServer(s, &server.Server{})

	go Client()

	s.Serve(lis)

}

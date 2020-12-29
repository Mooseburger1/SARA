package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Mooseburger1/SARA/backend/server"
	"google.golang.org/grpc"
)

func Client() {
	fmt.Println("Hello Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:2001", opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := server.NewPictureLinksServiceClient(cc)
	request := &server.PictureLinksRequest{}

	resp, _ := client.PictureLinks(context.Background(), request)

	fmt.Println(resp.Urls)
}

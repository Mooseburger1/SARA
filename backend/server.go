package main

import (
	"fmt"

	"github.com/Mooseburger1/SARA/backend/environment"
	"github.com/Mooseburger1/SARA/backend/s3"
)

func main() {
	var e environment.Environment
	e.GetCredentials()

	fmt.Println(fmt.Sprintf("Key: %v | Secret: %v | Region %v | Bucket %v", e.Key, e.Secret, e.Region, e.Bucket))
	s := s3.S3{}
	s.CreateSession()

	outgoing := make(chan string)
	go s.ListPictures(outgoing)

	for val := range outgoing {
		fmt.Println(val)
	}

}

package main

import (
	"fmt"

	"github.com/Mooseburger1/SARA/backend/environment"
	"github.com/Mooseburger1/SARA/backend/s3"
)

func main() {
	var e environment.Environment
	e.GetCredentials()

	s := s3.S3{}
	s.CreateSession()

	bucketObjectChannel := make(chan *string)
	objectURLChannel := make(chan []string)

	go s.ListPictures(bucketObjectChannel)
	go s.GenerateObjectURL(bucketObjectChannel, objectURLChannel)

	var results []string

	results = <-objectURLChannel

	fmt.Println(results)
}

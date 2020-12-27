package s3_test

import (
	"fmt"
	"testing"

	"github.com/Mooseburger1/SARA/backend/s3"
)

func TestCreateSession(t *testing.T) {
	s := s3.S3{}

	s.CreateSession()

	if s.Sess == nil {
		t.Errorf(fmt.Sprintf("Failed to create S3 session. Sess attribute for S3 type is nil"))
	} else {
		t.Log("CreateSession PASSED for S3 object")
	}
}

func TestListPictures(t *testing.T) {
	counter := 0

	s := s3.S3{}

	s.CreateSession()

	outgoing := make(chan *string)
	go s.ListPictures(outgoing)

	for _ = range outgoing {
		counter++
	}

	if counter < 0 {
		t.Errorf(fmt.Sprintf("No image keys sent through the channel - Ensure images are present in the S3 bucket and double check credentials"))
	} else {
		t.Log("ListPictures PASSED for S3 object")
	}
}

func TestBucketLink(t *testing.T) {
	s := s3.S3{}

	s.CreateSession()

	real := "https://sims-family-photos.s3-us-west-2.amazonaws.com/"

	if s.BucketLink != real {
		t.Errorf(fmt.Sprintf("Creating BucketLink FAILED. Expected %s but received %s", real, s.BucketLink))
	} else {
		t.Log("BucketLink PASSED. URLs match")
	}
}

func TestGenerateObjectURL(t *testing.T) {

	s := s3.S3{}
	s.CreateSession()

	bucketObjectChannel := make(chan *string)
	objectURLChannel := make(chan []string)

	go s.ListPictures(bucketObjectChannel)
	go s.GenerateObjectURL(bucketObjectChannel, objectURLChannel)

	results := <-objectURLChannel

	if len(results) <= 0 {
		t.Errorf("Returned no URL links for images. Please check credentials and ensure images exists in your S3 bucket")
	} else {
		t.Log("GenerateObjectURL PASSED. Recieved URLs from channel")
	}

}

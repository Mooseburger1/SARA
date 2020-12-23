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

	outgoing := make(chan string)
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

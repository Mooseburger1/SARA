package s3

import (
	"fmt"

	"github.com/Mooseburger1/SARA/backend/environment"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 is the main interface
type S3 struct {
	ListOfPics *s3.ListBucketsOutput
	key        string
	secret     string
	region     string
	bucket     string
	Sess       *s3.S3
}

// CreateSession creates session
func (s *S3) CreateSession() {
	var e environment.Environment
	e.GetCredentials()

	key := e.Key
	secret := e.Secret
	region := e.Region
	bucket := e.Bucket

	if len(key) == 0 {
		panic(fmt.Sprintf("No valid AWS key found while trying to create session"))
	} else if len(secret) == 0 {
		panic(fmt.Sprintf("No valid AWS secret found while trying to create session"))
	} else if len(region) == 0 {
		panic(fmt.Sprintf("No valid AWS region found while trying to create session"))
	} else if len(bucket) == 0 {
		panic(fmt.Sprintf("No valid AWS region found while trying to create session"))
	} else {
		s.key = key
		s.secret = secret
		s.region = region
		s.bucket = bucket
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(s.region),
			Credentials: credentials.NewStaticCredentials(s.key, s.secret, ""),
		},
	)

	if err != nil {
		panic(fmt.Sprintf("Error creating NewSession to AWS: %v", err))
	}

	svc := s3.New(sess)

	s.Sess = svc

}

//ListPictures show all objects in the SARA bucket
func (s *S3) ListPictures(outgoing chan<- string) {
	defer close(outgoing)
	sess := s.Sess

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int64(25),
	}

	result, err := sess.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	for _, item := range result.Contents {
		size := item.Size
		key := item.Key

		if *size > 0 {
			outgoing <- *key
		}

	}

}

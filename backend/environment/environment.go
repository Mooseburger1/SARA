package environment

import (
	"os"
)

//Environment is used for getting and persisting set environment variables for AWS access
type Environment struct {
	Key    string
	Secret string
	Region string
	Bucket string
}

// GetCredentials gets AWS variables from environment variables
func (e *Environment) GetCredentials() *Environment {
	e.Key = os.Getenv("SARA_KEY")
	e.Secret = os.Getenv("SARA_SECRET")
	e.Region = os.Getenv("SARA_REGION")
	e.Bucket = os.Getenv("SARA_BUCKET")

	return e
}

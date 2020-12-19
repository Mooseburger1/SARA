package environment

import (
	"os"
)

//Environment is used for getting environment variables for AWS
type Environment struct {
	Key    string
	Secret string
}

// GetCredentials gets AWS variables from environment
func (e *Environment) GetCredentials() *Environment {
	e.Key = os.Getenv("SARA_KEY")
	e.Secret = os.Getenv("SARA_SECRET")
	return e
}

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
	e.Key = os.Getenv("KEY")
	e.Secret = os.Getenv("SECRET")
	return e
}

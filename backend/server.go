package main

import (
	"fmt"

	"github.com/Mooseburger1/SARA/backend/environment"
)

func main() {
	var e environment.Environment
	e.GetCredentials()

	fmt.Printf("KEY: %v", e.Key)
	fmt.Printf("SECRET: %v", e.Secret)
}

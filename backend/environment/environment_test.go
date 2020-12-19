package environment_test

import (
	"os"
	"testing"

	"github.com/Mooseburger1/SARA/backend/environment"
)

func TestGetCredentialsKey(t *testing.T) {
	var env environment.Environment
	env.GetCredentials()

	key := env.Key

	if len(key) == 0 {
		t.Errorf("GetCredentials FAILED for KEY variable, expected %v but got %v", os.Getenv("KEY"), key)
	} else if key != os.Getenv("KEY") {
		t.Errorf("GetCredentials FAILED for KEY variable, expected %v but got %v", os.Getenv("KEY"), key)
	} else {
		t.Log("GetCredentials PASSED for KEY variable")
	}
}

func TestGetCredentialsSecret(t *testing.T) {
	env := environment.Environment{}
	env.GetCredentials()

	secret := env.Secret

	if len(secret) == 0 {
		t.Errorf("GetCredentials FAILED for SECRET variable, expected %v but got %v", os.Getenv("SECRET"), secret)
	} else if secret != os.Getenv("SECRET") {
		t.Errorf("GetCredentials FAILED for SECRET variable, expected %v but got %v", os.Getenv("SECRET"), secret)
	} else {
		t.Log("GetCredentials PASSED for SECRET variable")
	}
}

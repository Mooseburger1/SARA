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
		t.Errorf("GetCredentials FAILED for KEY variable, expected %v but got %v", os.Getenv("SARA_KEY"), key)
	} else if key != os.Getenv("SARA_KEY") {
		t.Errorf("GetCredentials FAILED for KEY variable, expected %v but got %v", os.Getenv("SARA_KEY"), key)
	} else {
		t.Log("GetCredentials PASSED for KEY variable")
	}
}

func TestGetCredentialsSecret(t *testing.T) {
	env := environment.Environment{}
	env.GetCredentials()

	secret := env.Secret

	if len(secret) == 0 {
		t.Errorf("GetCredentials FAILED for SECRET variable, expected %v but got %v", os.Getenv("SARA_SECRET"), secret)
	} else if secret != os.Getenv("SARA_SECRET") {
		t.Errorf("GetCredentials FAILED for SECRET variable, expected %v but got %v", os.Getenv("SARA_SECRET"), secret)
	} else {
		t.Log("GetCredentials PASSED for SECRET variable")
	}
}

func TestGetCredentialsRegion(t *testing.T) {
	env := environment.Environment{}
	env.GetCredentials()

	region := env.Region

	if len(region) == 0 {
		t.Errorf("GetCredentials FAILED for REGION variable, expected %v but got %v", os.Getenv("SARA_REGION"), region)
	} else if region != os.Getenv("SARA_REGION") {
		t.Errorf("GetCredentials FAILED for REGION variable, expected %v but got %v", os.Getenv("SARA_REGION"), region)
	} else {
		t.Log("GetCredentials PASSED for REGION variable")
	}
}

func TestGetCredentialsBucket(t *testing.T) {
	env := environment.Environment{}
	env.GetCredentials()

	bucket := env.Bucket

	if len(bucket) == 0 {
		t.Errorf("GetCredentials FAILED for BUCKET variable, expected %v but got %v", os.Getenv("SARA_BUCKET"), bucket)
	} else if bucket != os.Getenv("SARA_BUCKET") {
		t.Errorf("GetCredentials FAILED for BUCKET variable, expected %v but got %v", os.Getenv("SARA_BUCKET"), bucket)
	} else {
		t.Log("GetCredentials PASSED for BUCKET variable")
	}
}

package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/boj/redistore.v1"
)

func TestNewGoogleClient(t *testing.T) {

	mockLogger := log.New(os.Stdout, "test", log.LstdFlags)
	mockconf := ConfigBuilder()
	mockClient := &http.Client{}
	mockStore := &redistore.RediStore{}

	actualGC := googleClient{
		logger: mockLogger,
		conf:   mockconf,
		Client: mockClient,
		store:  mockStore,
	}

	gc := NewGoogleClientBuilder().
		SetLogger(mockLogger).
		SetConfig(mockconf).
		SetClient(mockClient).
		SetStore(mockStore).
		Build()

	if actualGC != *gc {
		t.Fatalf("Expected %v but got %v", actualGC, gc)
	}

}

func TestAuthenticate(t *testing.T) {
	gc := NewGoogleClientBuilder().
		SetConfig(ConfigBuilder()).
		Build()

	req, err := http.NewRequest("GET", "/authenticate", nil)

	if err != nil {
		t.Fatalf("Error creating request for authenticate: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gc.Authenticate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

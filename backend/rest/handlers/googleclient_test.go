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
	mockStore := &redistore.RediStore{}

	actualGC := GoogleClient{
		logger: mockLogger,
		conf:   mockconf,
		store:  mockStore,
	}

	gc := NewGoogleClientBuilder().
		SetLogger(mockLogger).
		SetConfig(mockconf).
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

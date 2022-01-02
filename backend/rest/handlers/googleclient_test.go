package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewGoogleClient(t *testing.T) {

	mockLogger := log.New(os.Stdout, "test", log.LstdFlags)

	actualGC := GoogleClient{
		logger: mockLogger,
	}

	gc := NewGoogleClient(mockLogger)

	if actualGC != *gc {
		t.Fatalf("Expected %v but got %v", actualGC, gc)
	}

}

func TestListAlbums(t *testing.T) {
	mockLogger := log.New(os.Stdout, "test", log.LstdFlags)
	mockClient := &http.Client{}
	gc := NewGoogleClient(mockLogger)

	req, err := http.NewRequest("GET", "/list-albums", nil)

	if err != nil {
		t.Fatalf("Error creating request for authenticate: %v", err)
	}

	rr := httptest.NewRecorder()
	hfunc := convertTrueHandlerFunc(mockClient, gc.ListAlbums)
	handler := http.HandlerFunc(hfunc)

	handler.ServeHTTP(rr, req)
	var result errorJson

	err = json.NewDecoder(rr.Body).Decode(&result)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if result.Error.Code != 401 {
		t.Errorf("handler returned wrong http response code: got %v want %v", result.Error.Code, 401)
	}

	if result.Error.Status != "UNAUTHENTICATED" {
		t.Errorf("handler returned wrong status: got %v want %v", result.Error.Status, "UNAUTHENTICATED")
	}
}

func TestListPhotos(t *testing.T) {
	mockLogger := log.New(os.Stdout, "test2", log.LstdFlags)
	mockClient := &http.Client{}
	gc := NewGoogleClient(mockLogger)

	req, err := http.NewRequest("POST", "/list-photos-from-album", nil)

	if err != nil {
		t.Fatalf("Error creating request for authenticate: %v", err)
	}

	rr := httptest.NewRecorder()
	hfunc := convertTrueHandlerFunc(mockClient, gc.ListPicturesFromAlbum)
	handler := http.HandlerFunc(hfunc)

	handler.ServeHTTP(rr, req)
	var result errorJson
	err = json.NewDecoder(rr.Body).Decode(&result)
	if err != nil {
		t.Fatalf("fatal error %v", err)
		return
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if result.Error.Code != 401 {
		t.Errorf("handler returned wrong http response code: got %v want %v", result.Error.Code, 401)
	}

	if result.Error.Status != "UNAUTHENTICATED" {
		t.Errorf("handler returned wrong status: got %v want %v", result.Error.Status, "UNAUTHENTICATED")
	}
}

func convertTrueHandlerFunc(client *http.Client, f func(http.ResponseWriter, *http.Request, *http.Client)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, client)
	}
}

type errorJson struct {
	Error errorData `json:"error"`
}

type errorData struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

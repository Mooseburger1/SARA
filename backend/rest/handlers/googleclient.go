package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"log"

	clientProto "backend/grpc/proto/api/client"

	"github.com/gorilla/mux"
)

const (
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
	PHOTOS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/mediaItems:search"
	GET             = "GET"
	POST            = "POST"
)

type GoogleClient struct {
	logger *log.Logger
}

func NewGoogleClient(logger *log.Logger) *GoogleClient {
	return &GoogleClient{logger: logger}
}

// ListAlbums utilizes photoslibrary googleapis to list all albums in the
// Google photos account.
func (gc GoogleClient) ListAlbums(rw http.ResponseWriter, r *http.Request, client *clientProto.ClientInfo) {

	_ = makeAlbumListRequestBody(r, gc.logger)
	return

}

func (gc GoogleClient) ListPicturesFromAlbum(rw http.ResponseWriter, r *http.Request, client *clientProto.ClientInfo) {

	// var result datamodels.MediaItems

	// req, err := http.NewRequest("POST", PHOTOS_ENDPOINT, strings.NewReader(makeAlbumIdRequestBody(r)))
	// if err != nil {
	// 	gc.logger.Fatalf("Failed to create new request: %v", err)
	// }

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "application/json")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	gc.logger.Printf("Get: %v\n", err.Error())
	// 	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	response, _ := ioutil.ReadAll(resp.Body)
	// 	rw.Write([]byte(string(response)))
	// 	return
	// }

	// err = json.NewDecoder(resp.Body).Decode(&result)
	// if err != nil {
	// 	gc.logger.Printf("ReadAll: %v\n", err.Error())
	// 	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	// gc.logger.Printf("%v", result)

	return

}

func makeAlbumListRequestBody(r *http.Request, logger *log.Logger) string {
	vars := mux.Vars(r)
	pageSize := vars["pageSize"]
	pageToken := vars["pageToken"]

	logger.Printf("pageSize: %v", pageSize)
	logger.Printf("pageToken: %v", pageToken)
	return "hello"
}

func makeAlbumIdRequestBody(r *http.Request) string {
	vars := mux.Vars(r)
	return fmt.Sprintf(`{"albumId": "%v"}`, vars["albumId"])
}

// OhNo is the default Redirect handler for when a user has done something stupid
func (gc GoogleClient) OhNo(rw http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.New("Oh-No").Parse(`
	<h1>OH NO</h1>`))

	_ = templ.Execute(rw, nil)

	return
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"log"

	photosProto "backend/grpc/proto/api/photos"

	"github.com/gorilla/mux"
)

type GoogleClient struct {
	logger       *log.Logger
	photosClient *photosProto.GooglePhotoServiceClient
}

// NewGoogleClient creates a GoogleClient instance. The instance
// exposes methods to make RPC calls to the photos RPC server that
// interacts with the Google photos API
func NewGoogleClient(logger *log.Logger) *GoogleClient {
	return &GoogleClient{
		logger: logger,
	}
}

// ListAlbums makes RPC call to the photos RPC server. More specifically
// it invokes the ListAlbums endpoint of ther RPC server.
func (gc *GoogleClient) ListAlbums(rw http.ResponseWriter, r *http.Request, ai *photosProto.AlbumsInfo) {

	JSON, err := json.Marshal(ai)
	if err != nil {
		gc.logger.Printf("Unable to marshal: %v", err)
	}
	rw.Write(JSON)
}

func (gc *GoogleClient) ListPhotosFromAlbum(rw http.ResponseWriter, r *http.Request, pi *photosProto.PhotosInfo) {

	JSON, err := json.Marshal(pi)
	if err != nil {
		gc.logger.Printf("Unable to marshal: %v", err)
	}
	rw.Write(JSON)
}

// str2Int32 is a package private helper function
// for type conversion
func str2Int32(val string) (int32, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

func makeAlbumIdRequestBody(r *http.Request) string {
	vars := mux.Vars(r)
	return fmt.Sprintf(`{"albumId": "%v"}`, vars["albumId"])
}

// OhNo is the default Redirect handler for when a user has done something stupid
func (gc *GoogleClient) OhNo(rw http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.New("Oh-No").Parse(`
	<h1>OH NO</h1>`))

	_ = templ.Execute(rw, nil)

	return
}

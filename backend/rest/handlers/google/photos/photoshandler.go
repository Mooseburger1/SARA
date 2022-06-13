package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"log"

	photosProto "backend/grpc/proto/api/photos"
)

type PhotoHandler struct {
	logger *log.Logger
}

// NewPhotoHandler creates a PhotoHandler instance. The instance
// exposes methods to make RPC calls to the photos RPC server that
// interacts with the Google photos API
func NewPhotoHandler(logger *log.Logger) *PhotoHandler {
	return &PhotoHandler{
		logger: logger,
	}
}

// ListAlbums makes RPC call to the photos RPC server. More specifically
// it invokes the ListAlbums endpoint of ther RPC server.
func (gc *PhotoHandler) ListAlbums(rw http.ResponseWriter, r *http.Request, ai *photosProto.AlbumsInfo) {

	JSON, err := json.Marshal(ai)
	if err != nil {
		gc.logger.Printf("Unable to marshal: %v", err)
	}
	rw.Write(JSON)
}

func (gc *PhotoHandler) ListPhotosFromAlbum(rw http.ResponseWriter, r *http.Request, pi *photosProto.PhotosInfo) {

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

// OhNo is the default Redirect handler for when a user has done something stupid
func (gc *PhotoHandler) OhNo(rw http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.New("Oh-No").Parse(`
	<h1>OH NO</h1>`))

	_ = templ.Execute(rw, nil)

	return
}

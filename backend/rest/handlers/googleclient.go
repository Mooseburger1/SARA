package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"log"

	clientProto "backend/grpc/proto/api/client"
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
func NewGoogleClient(logger *log.Logger, pc *photosProto.GooglePhotoServiceClient) *GoogleClient {
	return &GoogleClient{
		logger:       logger,
		photosClient: pc,
	}
}

// ListAlbums makes RPC call to the photos RPC server. More specifically
// it invokes the ListAlbums endpoint of ther RPC server.
func (gc *GoogleClient) ListAlbums(rw http.ResponseWriter, r *http.Request, client *clientProto.ClientInfo) {
	listRequest := makeAlbumListRequest(r, client)
	pc := *gc.photosClient
	albums, err := pc.ListAlbums(context.Background(), listRequest)
	if err != nil {
		panic(err)
	}

	gc.logger.Printf("%v", albums)

	return

}

func (gc *GoogleClient) ListPicturesFromAlbum(rw http.ResponseWriter, r *http.Request, client *clientProto.ClientInfo) {
	photoRequest := makePhotosFromAlbumRequest(r, client)
	pc := *gc.photosClient
	photos, err := pc.ListPhotosFromAlbum(context.Background(), photoRequest)
	if err != nil {
		panic(err)
	}

	gc.logger.Printf("%v", photos)


	return

}


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



func makePhotosFromAlbumRequest(r *http.Request, ci *clientProto.ClientInfo) *photosProto.FromAlbumRequest {
	vars := mux.Vars(r)
	albumId := vars["albumId"]
	pageToken := vars["pageToken"]
	pageSize := vars["pageSize"]

	var req photosProto.FromAlbumRequest
	req.ClientInfo = ci
	req.AlbumId = albumId

	// Parse the pageSize Url variable
	if pageSize != "" {
		i, err := str2Int32(pageSize)
		if err != nil {
			panic(err)
		}
		req.PageSize = i
	}

	// Parse the pageToken URL variable
	if pageToken != "" {
		req.PageToken = pageToken
	}

	return &req
}




// makeAlbumListRequest is a package private helper function
// utilized to extract variables from the API URL and generate
// an AlbumListRequest proto. More specifically, it is a parser
// for the REST endpoint of list-albums and constructs the necessary
// RPC proto.
func makeAlbumListRequest(r *http.Request, ci *clientProto.ClientInfo) *photosProto.AlbumListRequest {
	vars := mux.Vars(r)
	pageSize := vars["pageSize"]
	pageToken := vars["pageToken"]

	var req photosProto.AlbumListRequest
	req.ClientInfo = ci

	// Parse the pageSize Url variable
	if pageSize != "" {
		i, err := str2Int32(pageSize)
		if err != nil {
			panic(err)
		}
		req.PageSize = i
	}

	// Parse the pageToken URL variable
	if pageToken != "" {
		req.PageToken = pageToken
	}

	return &req
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

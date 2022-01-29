package main

import (
	protos "backend/grpc/proto/api/photos"
	"backend/rest/handlers"
	"backend/rest/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/boj/redistore.v1"
)

func main() {
	// Main logger
	logger := log.New(os.Stdout, "rest-server", log.LstdFlags)

	// Initialize redis store
	store, err := redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))

	if err != nil {
		logger.Fatalf("Error processing redistore %v", err)
		return
	}

	/////// Initialize GRPC connections
	photoConn, err := grpc.Dial("grpc_backend:9091", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer photoConn.Close()
	gpsc := protos.NewGooglePhotoServiceClient(photoConn)

	/////// Initialize middleware and handlers here ///////
	gClient := handlers.NewGoogleClient(logger, &gpsc)
	mWare := middleware.NewMiddleWare(logger, store)

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()

	// GET SUBROUTER
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	//getRouter.HandleFunc("/", )
	getRouter.HandleFunc("/authenticate", mWare.Authenticate)
	getRouter.HandleFunc("/oauth-callback", mWare.RedirectCallback)

	//route for listing albums - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/albumsList", mWare.Authorized(gClient.ListAlbums))
	getRouter.HandleFunc("/oh-no", gClient.OhNo)

	// PUT SUBROUTER
	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/", mWare.Authenticate)

	// POST SUBROUTER
	//postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	//route for listing photos in an album - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/album/{albumId:[-_0-9A-Za-z]+}", mWare.Authorized(gClient.ListPicturesFromAlbum))

	// Configure the server {TODO: move these to an external configurable file/location}
	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Asynchronously expose the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	//Block while not receiving forecfull shutdown
	sig := <-sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	//Define context to provide server on how to shutdown all running processes
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

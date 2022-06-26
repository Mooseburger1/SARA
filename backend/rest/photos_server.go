package main

import (
	config "backend/configuration"
	protos "backend/grpc/proto/api/photos"
	photoshandlers "backend/rest/handlers/google/photos"
	gAuth "backend/rest/middleware/google/auth/OAuth"
	callingCatchables "backend/rest/middleware/google/callingCatchables/photos"

	gohandlers "github.com/gorilla/handlers"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/boj/redistore.v1"
)

type PhotoServer struct {
	server *http.Server
	logger *log.Logger
}

func GetPhotoServer() *PhotoServer {
	ps := PhotoServer{}
	ps.initPhotosServer()
	return &ps
}

func (ps *PhotoServer) initPhotosServer() {
	// Main logger
	ps.logger = log.New(os.Stdout, "rest-server-photos", log.LstdFlags)

	// Internal Configuration
	config := config.NewGOAuthConfig()

	// Initialize redis store
	store, err := redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))

	if err != nil {
		ps.logger.Fatalf("Error processing redistore %v", err)
		return
	}

	/////// Initialize GRPC connections
	photoConn, err := grpc.Dial("grpc_backend:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer photoConn.Close()
	gpsc := protos.NewGooglePhotoServiceClient(photoConn)

	/////// Initialize middleware and handlers here ///////
	pHandler := photoshandlers.NewPhotoHandler(ps.logger)
	gAuthMware := gAuth.NewAuthMiddleware(ps.logger, store, config)
	gPhotos := callingCatchables.NewPhotosRpcCaller(ps.logger, &gpsc)

	// CORS Handler
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:4200"}))

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()

	// GET SUBROUTER
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	//getRouter.HandleFunc("/", )
	getRouter.HandleFunc("/authenticate", gAuthMware.Authenticate)
	getRouter.HandleFunc("/oauth-callback", gAuthMware.RedirectCallback)

	//route for listing albums - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/albumsList", gAuthMware.IsAuthorized(gPhotos.CatchableListAlbums(pHandler.ListAlbums)))

	//route for listing photos in an album - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/album/{albumId:[-_0-9A-Za-z]+}", gAuthMware.IsAuthorized(gPhotos.CatchablePhotosFromAlbum(pHandler.ListPhotosFromAlbum)))

	getRouter.HandleFunc("/oh-no", pHandler.OhNo)

	// PUT SUBROUTER
	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/", gAuthMware.Authenticate)

	// POST SUBROUTER
	//postRouter := serveMux.Methods(http.MethodPost).Subrouter()

	// Configure the server {TODO: move these to an external configurable file/location}
	server := &http.Server{
		Addr:         ":9090",
		Handler:      corsHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	ps.server = server

	// // Asynchronously expose the server
	// go func() {
	// 	err := server.ListenAndServe()
	// 	if err != nil {
	// 		ps.logger.Fatal(err)
	// 	}
	// }()

	// sigChan := make(chan os.Signal)
	// signal.Notify(sigChan, os.Interrupt)
	// signal.Notify(sigChan, os.Kill)

	// //Block while not receiving forecfull shutdown
	// sig := <-sigChan

	// ps.logger.Println("Received terminate, graceful shutdown", sig)

	// //Define context to provide server on how to shutdown all running processes
	// tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// server.Shutdown(tc)
}

func (ps *PhotoServer) StartServer() {
	err := ps.server.ListenAndServe()
	if err != nil {
		ps.logger.Fatal(err)
	}
}

func (ps *PhotoServer) ShutdownServer(tc context.Context) {
	err := ps.server.Shutdown(tc)
	if err != nil {
		ps.logger.Fatal(err)
	}
}

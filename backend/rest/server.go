package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest/handlers"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
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

	/////// Initialize handlers here ///////
	gClient := handlers.NewGoogleClientBuilder().
		SetLogger(logger).
		SetStore(store).
		SetConfig(handlers.ConfigBuilder()).
		Build()

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()

	//Create filtered Routers to handle specific verbs
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	//getRouter.HandleFunc("/", )
	getRouter.HandleFunc("/authenticate", gClient.Authenticate)
	getRouter.HandleFunc("/oauth-callback", gClient.RedirectCallback)
	getRouter.HandleFunc("/list-albums", gClient.ListAlbums)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/", gClient.Authenticate)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", gClient.Authenticate)

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

package main

import (
	utils "backend/utils"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	photoServer := GetPhotoServer()

	calendarServer := GetCalendarServer()

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()

	// GET SUBROUTER
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()

	photoServer.RegisterGetRoutes(getRouter)

	calendarServer.RegisterGetRoutes(getRouter)

	// CORS Handler
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:4200"}))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      corsHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	sigChan := utils.GetOsKillerListener()

	sig := <-*sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(tc)

}

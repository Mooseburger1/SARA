package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type CalendarServer struct {
	server *http.Server
	logger *log.Logger
}

func GetCalendarServer() *CalendarServer {
	cs := CalendarServer{}
	cs.initCalendarServer()
	return &cs
}

func (cs *CalendarServer) initCalendarServer() {

	cs.logger = log.New(os.Stdout, "rest-server-calendar", log.LstdFlags)

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:4200"}))

	serveMux := mux.NewRouter()

	//getRrouter := serveMux.Methods(http.MethodGet).Subrouter()

	server :=
		&http.Server{
			Addr:         ":9091",
			Handler:      corsHandler(serveMux),
			IdleTimeout:  120 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		}

	cs.server = server
}

func (cs *CalendarServer) StartServer() {
	err := cs.server.ListenAndServe()
	if err != nil {
		cs.logger.Fatal(err)
	}
}

func (cs *CalendarServer) ShutdownServer(tc context.Context) {
	err := cs.server.Shutdown(tc)
	if err != nil {
		cs.logger.Fatal(err)
	}
}

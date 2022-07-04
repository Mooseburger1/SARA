package main

import (
	config "backend/configuration"
	protos "backend/grpc/proto/api/photos"
	photoshandlers "backend/rest/handlers/google/photos"
	gAuth "backend/rest/middleware/google/auth/OAuth"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/boj/redistore.v1"
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

	config := config.NewGOAuthConfig()

	store, err := redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))

	if err != nil {
		cs.logger.Fatalf("Error processing redistore %v", err)
		return
	}

	/////// Initialize GRPC connections
	calendarConn, err := grpc.Dial("grpc_backend:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer calendarConn.Close()
	gcsc := protos.NewGooglePhotoServiceClient(calendarConn)

	/////// Initialize middleware and handlers here ///////
	cHandler := photoshandlers.NewPhotoHandler(cs.logger)
	gAuthMware := gAuth.NewAuthMiddleware(cs.logger, store, config)
	gCalendar := callingCatchables.NewCalendarRpcCaller(cs.logger, &gcsc)

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:4200"}))

	serveMux := mux.NewRouter()

	// GET SUBROUTER
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()

	// route for listing calendars
	getRouter.HandleFunc("/calendar/listCalendars", test)

	server :=
		&http.Server{
			Addr:         ":9092",
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

func test(rw http.ResponseWriter, r *http.Request) {

	rw.Write([]byte("Hello, Testing 1..2..3.."))
}

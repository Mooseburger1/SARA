package main

import (
	protos "backend/grpc/proto/api/calendar"
	calendarhandlers "backend/rest/handlers/google/calendar"
	gAuth "backend/rest/middleware/google/auth/OAuth"
	callingCatchables "backend/rest/middleware/google/callingCatchables/calendar"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalendarServer struct {
	logger     *log.Logger
	gcsc       protos.GoogleCalendarServiceClient
	cHandler   *calendarhandlers.CalendarHandler
	gAuthMware *gAuth.AuthMiddleware
	gCalendar  *callingCatchables.CalendarRpcCaller
}

func GetCalendarServer() *CalendarServer {
	cs := CalendarServer{}
	cs.initService()
	return &cs
}

func (cs *CalendarServer) initService() {

	cs.logger = log.New(os.Stdout, "rest-server-calendar", log.LstdFlags)

	/////// Initialize GRPC connections
	calendarConn, err := grpc.Dial("grpc_backend:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// defer calendarConn.Close()
	cs.gcsc = protos.NewGoogleCalendarServiceClient(calendarConn)

	/////// Initialize middleware and handlers here ///////
	cs.cHandler = calendarhandlers.NewCalendarHandler(cs.logger)
	cs.gAuthMware = gAuth.GetAuthMiddleware()
	cs.gCalendar = callingCatchables.NewCalendarRpcCaller(cs.logger, &cs.gcsc)

}

func (cs *CalendarServer) RegisterGetRoutes(getRouter *mux.Router) {

	// route for tesing if this server is up and running
	getRouter.HandleFunc("/test", test)

	// route for listing calendars
	getRouter.HandleFunc("/calendar/listCalendars", cs.gAuthMware.IsAuthorized(cs.gCalendar.CatchableListCalendars(cs.cHandler.ListCalendars)))

}

func test(rw http.ResponseWriter, r *http.Request) {

	rw.Write([]byte("Hello, Testing 1..2..3.."))
}

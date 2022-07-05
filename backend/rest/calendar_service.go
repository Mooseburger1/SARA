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

type CalendarService struct {
	logger     *log.Logger
	gcsc       protos.GoogleCalendarServiceClient
	cHandler   *calendarhandlers.CalendarHandler
	gAuthMware *gAuth.AuthMiddleware
	gCalendar  *callingCatchables.CalendarRpcCaller
}

func GetCalendarService() *CalendarService {
	cs := CalendarService{}
	return &cs
}

//InitServiceAndReturnCloseConnectionFunc intializes the service. As a part of the
// initialization, it establishes a connection with the proper gRPC server. This
// connection needs to have a way to be closed at application termination. In order
// to defer those close, this function returns a function which encapsulates the closing
// of this conneciton. You should defer the returned function in your main function in order
// to keep the connection alive until the program is done.
func (cs *CalendarService) InitServiceAndReturnCloseConnectionFunc() func() {

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

	return func() { calendarConn.Close() }
}

func (cs *CalendarService) RegisterGetRoutes(getRouter *mux.Router) {

	// route for tesing if this server is up and running
	getRouter.HandleFunc("/test", test)

	// route for listing calendars - optional params {pageToken | maxResults | showDeleted | showHidden | syncToken}
	getRouter.HandleFunc("/calendar/listCalendars", cs.gAuthMware.IsAuthorized(cs.gCalendar.CatchableListCalendars(cs.cHandler.ListCalendars)))

}

func test(rw http.ResponseWriter, r *http.Request) {

	rw.Write([]byte("Hello, Testing 1..2..3.."))
}

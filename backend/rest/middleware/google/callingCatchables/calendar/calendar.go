package callingCatchables

import (
	calendarProto "backend/grpc/proto/api/calendar"
	clientProto "backend/grpc/proto/api/client"
	auth "backend/rest/middleware/google/auth/OAuth"
	utils "backend/rest/middleware/google/callingCatchables/common"
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc/status"
)

// CalendarRpcCaller is the client responsible for making calls
// to the gRPC server for the Google Calendar endpoints. Successful
// calls will be propogated to the injected handlers. Failed RPC
// calls will be caught and handled gracefully.
type CalendarRpcCaller struct {
	logger         *log.Logger
	calendarClient *calendarProto.GoogleCalendarServiceClient
}

// CalendarListHandlerFunc is a http.HandlerFunc extended to handle the successful request
// to the gRPC server for Google Calendar and specifically for the ListCalendars endpoint.
type CalendarListHandlerFunc func(http.ResponseWriter, *http.Request, *calendarProto.CalendarListResponse)

// NewCalendarRpcCaller is a builder for a CalendarRpcCaller client. It will create a new instance
// with each invocation. Does not follow the singleton pattern.
func NewCalendarRpcCaller(logger *log.Logger, cc *calendarProto.GoogleCalendarServiceClient) *CalendarRpcCaller {
	return &CalendarRpcCaller{
		logger:         logger,
		calendarClient: cc,
	}
}

// CatchableListCalendars makes a request to the RPC server for the ListCalendars endpoint. A successful
// request is propagated forward to the supplied CalendarListHandlerFunc. All errors will be caught and
// the error will be returned to the client caller
func (rpc *CalendarRpcCaller) CatchableListCalendarsList(handler CalendarListHandlerFunc) auth.ClientHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		listRequest := makeListCalendarsListRequest(r, clientInfo)
		cc := *rpc.calendarClient
		list, err := cc.ListCalendarList(context.Background(), listRequest)

		if err != nil {
			st := status.Convert(err)
			utils.Route404Error(st, rw)
			return
		}
		handler(rw, r, list)
	}
}

// CatchableListCalendars makes a request to the RPC server for the ListCalendars endpoint. A successful
// request is propagated forward to the supplied CalendarListHandlerFunc. All errors will be caught and
// the error will be returned to the client caller
func (rpc *CalendarRpcCaller) CatchableGetCalendarsList(handler CalendarListHandlerFunc) auth.ClientHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		getRequest := makeListCalendarsGetRequest(r, clientInfo)
		cc := *rpc.calendarClient
		calendar, err := cc.ListCalendarGet(context.Background(), getRequest)

		if err != nil {
			st := status.Convert(err)
			utils.Route404Error(st, rw)
			return
		}
		handler(rw, r, calendar)
	}
}

func makeListCalendarsGetRequest(r *http.Request, ci *clientProto.ClientInfo) *calendarProto.ListCalendarGetRequest {
	calendarId := r.URL.Query().Get("calendarId")
	var req calendarProto.ListCalendarGetRequest
	req.ClientInfo = ci

	if calendarId != "" {
		req.CalendarId = calendarId
	}

	return &req
}

// makeListCalendarRequest is a package private helper function
// utilized to extract query variables from the API URL and generate
// an CalendarListRequst proto. More specifically, it is a parser
// for the REST endpoint of list-calendars and constructs the necessary
// RPC proto.
func makeListCalendarsListRequest(r *http.Request, ci *clientProto.ClientInfo) *calendarProto.CalendarListRequest {
	pageToken := r.URL.Query().Get("pageToken")
	maxResults := r.URL.Query().Get("maxResults")
	showDeleted := r.URL.Query().Get("showDeleted")
	showHidden := r.URL.Query().Get("showHidden")
	syncToken := r.URL.Query().Get("syncToken")

	var req calendarProto.CalendarListRequest
	req.ClientInfo = ci

	if pageToken != "" {
		req.PageToken = pageToken
	}

	if maxResults != "" {
		i, err := utils.Str2Int32(maxResults)
		if err != nil {
			panic(err)
		}
		req.MaxResults = i
	}

	if showDeleted != "" {
		b, err := utils.Str2Bool(showDeleted)
		if err != nil {
			panic(err)
		}
		req.ShowDeleted = b
	}

	if showHidden != "" {
		b, err := utils.Str2Bool(showHidden)
		if err != nil {
			panic(err)
		}
		req.ShowHidden = b
	}

	if syncToken != "" {
		req.SyncToken = syncToken
	}

	return &req
}

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

type calendarRpcCaller struct {
	logger         *log.Logger
	calendarClient *calendarProto.GoogleCalendarServiceClient
}

type CalendarListHandlerFunc func(http.ResponseWriter, *http.Request, *calendarProto.CalendarListResponse)

func NewCalendarRpcCaller(logger *log.Logger, cc *calendarProto.GoogleCalendarServiceClient) *calendarRpcCaller {
	return &calendarRpcCaller{
		logger:         logger,
		calendarClient: cc,
	}
}

func (rpc *calendarRpcCaller) CatchableListCalendars(handler CalendarListHandlerFunc) auth.ClientHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, clientInfo *clientProto.ClientInfo) {
		listRequest := makeListCalendarRequest(r, clientInfo)
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

func makeListCalendarRequest(r *http.Request, ci *clientProto.ClientInfo) *calendarProto.CalendarListRequest {
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

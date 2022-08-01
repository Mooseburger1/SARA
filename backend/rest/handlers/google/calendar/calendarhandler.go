package handlers

import (
	calendarProto "backend/grpc/proto/api/calendar"
	"encoding/json"
	"log"
	"net/http"
)

type CalendarHandler struct {
	logger *log.Logger
}

// NewCalendarHandler creates a CalendarHandler instance. The instance is responsible
// for unmarshaling all responses from the gRPC calendar server and writing the response
// back to the client caller.
func NewCalendarHandler(logger *log.Logger) *CalendarHandler {
	return &CalendarHandler{
		logger: logger,
	}
}

// ListCalendars marshals the response from the gRPC server for the ListCalendars Endpoint and writes response back
// to client caller
func (ch *CalendarHandler) ListCalendarsList(rw http.ResponseWriter, r *http.Request, cl *calendarProto.CalendarListResponse) {
	JSON, err := json.Marshal(cl)
	if err != nil {
		ch.logger.Printf("Unable to marshal: %v", err)
	}

	rw.Write(JSON)
}

// GetCalendarsList marshals the response from the gRPC server for the ListCalendars Endpoint and writes response back
// to client caller
func (ch *CalendarHandler) GetCalendarsList(rw http.ResponseWriter, r *http.Request, cl *calendarProto.ListCalendarGetResponse) {
	JSON, err := json.Marshal(cl)
	if err != nil {
		ch.logger.Printf("Unable to marshal: %v", err)
	}

	rw.Write(JSON)
}

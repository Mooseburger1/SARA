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

func NewCalendarHandler(logger *log.Logger) *CalendarHandler {
	return &CalendarHandler{
		logger: logger,
	}
}

func (ch *CalendarHandler) ListCalendars(rw http.ResponseWriter, r *http.Request, cl *calendarProto.CalendarListResponse) {
	JSON, err := json.Marshal(cl)
	if err != nil {
		ch.logger.Printf("Unable to marshal: %v", err)
	}

	rw.Write(JSON)
}

package calendar

import (
	"backend/grpc/proto/api/calendar"
	"context"
	"log"
)

// GcalendarStub is the implementation of the
// Google Calendar RPC server. It implements the
// * ListCalendarList Service
type GcalendarStub struct {
	logger *log.Logger
}

// Constructor for instantiating a GcalendarStub
func NewGcalendarStub(logger *log.Logger) *GcalendarStub {
	return &GcalendarStub{logger: logger}
}

// ListCalendarList is a RPC service endpoint. It receives an CalendarListRequest
// proto and returns an CalendarListResponse proto. Internally it makes an Oauth2
// authorized REST request to the Google Calendar API server for listing calendars.
func (g *GcalendarStub) ListCalendarList(ctx context.Context, rpc *calendar.CalendarListRequest) (*calendar.CalendarListResponse, error) {
	return listCalendarList(rpc, g.logger)
}

package calendar

import (
	"backend/grpc/proto/api/calendar"
	"context"
	"log"
)

type GcalendarServer struct {
	logger *log.Logger
}

func NewGcalendarServer(logger *log.Logger) *GcalendarServer {
	return &GcalendarServer{logger: logger}
}

func (g *GcalendarServer) ListCalendarList(ctx context.Context, rpc *calendar.CalendarListRequest) (*calendar.CalendarListResponse, error) {
	return listCalendarList(rpc, g.logger)
}
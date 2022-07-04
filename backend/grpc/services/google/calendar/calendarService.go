package calendar

import (
	"backend/grpc/proto/api/calendar"
	"context"
	"log"
)

type GcalendarStub struct {
	logger *log.Logger
}

func NewGcalendarStub(logger *log.Logger) *GcalendarStub {
	return &GcalendarStub{logger: logger}
}

func (g *GcalendarStub) ListCalendarList(ctx context.Context, rpc *calendar.CalendarListRequest) (*calendar.CalendarListResponse, error) {
	return listCalendarList(rpc, g.logger)
}
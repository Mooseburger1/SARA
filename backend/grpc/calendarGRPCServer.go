package main

import (
	protoCalendar "backend/grpc/proto/api/calendar"
	calendar "backend/grpc/services/google/calendar"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type calendarGRPCServer struct {
	Server *grpc.Server
	Logger *log.Logger
}

func NewCalendarGRPCServer() *calendarGRPCServer {
	cs := calendarGRPCServer{}
	cs.initServer()
	return &cs
}

func (cs *calendarGRPCServer) initServer() {
	logger := log.New(os.Stdout, "calendar-rpc-server", log.LstdFlags)
	cs.Logger = logger

	grpcServer := grpc.NewServer()

	calendarStub := calendar.NewGcalendarStub(logger)

	protoCalendar.RegisterGoogleCalendarServiceServer(grpcServer, calendarStub)
	cs.Server = grpcServer
}

func (cs *calendarGRPCServer) StartServer() {
	l, err := net.Listen("tcp", ":9093")
	if err != nil {
		cs.Logger.Fatal(err)
		os.Exit(1)
	}
	cs.Server.Serve(l)
}

func (cs *calendarGRPCServer) ShutdownServer() {
	cs.Server.GracefulStop()
}

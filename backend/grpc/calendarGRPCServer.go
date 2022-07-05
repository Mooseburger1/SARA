package main

import (
	protoCalendar "backend/grpc/proto/api/calendar"
	calendar "backend/grpc/services/google/calendar"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type calendarGRPCServer struct {
	Server *grpc.Server
	Logger *log.Logger
}

// NewCalendarGRPCServer creates a new instance of a calendarGRPCServer.
// It intializes it first and returns to caller, ready to be started.
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

// StartServer will start the intialized calendarGRPCServer. It listens
// on port :9093
func (cs *calendarGRPCServer) StartServer() {
	reflection.Register(cs.Server)
	l, err := net.Listen("tcp", ":9093")
	if err != nil {
		cs.Logger.Fatal(err)
		os.Exit(1)
	}
	cs.Logger.Print("Calendar grpc listening on 9093")
	cs.Server.Serve(l)
}

func (cs *calendarGRPCServer) ShutdownServer() {
	cs.Server.GracefulStop()
}

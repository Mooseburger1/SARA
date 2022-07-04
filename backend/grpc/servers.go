package main

import (
	"backend/utils"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "grpc-server-manager", log.LstdFlags)

	photoGrpcServer := NewPhotosGRPCServer()

	calendarGrpcServer := NewCalendarGRPCServer()

	runServers(logger, photoGrpcServer, calendarGrpcServer)

	sigChan := utils.GetOsKillerListener()

	sig := <-*sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	shutdownServers(logger, photoGrpcServer, calendarGrpcServer)
}

func runServers(logger *log.Logger, servers ...SaraInterface) {
	for i, server := range servers {
		logger.Printf("Starting Server %d", i)
		go server.StartServer()
	}
}

func shutdownServers(logger *log.Logger, servers ...SaraInterface) {
	for i, server := range servers {
		logger.Printf("Shutting down server %d", i)
		server.ShutdownServer()
	}
}

type SaraInterface interface {
	StartServer()
	ShutdownServer()
}

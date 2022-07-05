package main

import (
	"backend/utils"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "grpc-server-manager", log.LstdFlags)

	//// Initialize all gRPC servers here
	photoGrpcServer := NewPhotosGRPCServer()
	calendarGrpcServer := NewCalendarGRPCServer()

	// Start em up
	runServers(logger, photoGrpcServer, calendarGrpcServer)

	// Listen for terminate signal
	sigChan := utils.GetOsKillerListener()
	sig := <-*sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	shutdownServers(logger, photoGrpcServer, calendarGrpcServer)
}

//runServers takes in a variable amount of SaraGRPCServers and starts them
// asynchronously
func runServers(logger *log.Logger, servers ...SaraGRPCServer) {
	for i, server := range servers {
		logger.Printf("Starting Server %d", i)
		go server.StartServer()
	}
}

func shutdownServers(logger *log.Logger, servers ...SaraGRPCServer) {
	for i, server := range servers {
		logger.Printf("Shutting down server %d", i)
		server.ShutdownServer()
	}
}

// SaraGRPCServer is an interface that any service supported by
// Sara must implement. If you create a new gRPC service, it must
// implement tis interface in order to be started / supported
type SaraGRPCServer interface {
	StartServer()
	ShutdownServer()
}

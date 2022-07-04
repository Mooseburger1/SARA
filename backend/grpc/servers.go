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

	runServers(photoGrpcServer, calendarGrpcServer)

	sigChan := utils.GetOsKillerListener()

	sig := <-*sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	shutdownServers(photoGrpcServer, calendarGrpcServer)
}

func runServers(servers ...SaraInterface) {
	for _, server := range servers {
		go server.StartServer()
	}
}

func shutdownServers(servers ...SaraInterface) {
	for _, server := range servers {
		server.ShutdownServer()
	}
}

type SaraInterface interface {
	StartServer()
	ShutdownServer()
}

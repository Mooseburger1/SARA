package main

import (
	utils "backend/utils"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	photoServer := GetPhotoServer()

	calendarServer := GetCalendarServer()

	runServers(photoServer,
		calendarServer)

	sigChan := utils.GetOsKillerListener()

	sig := <-*sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	shutdownServers(photoServer, calendarServer)

}

func runServers(servers ...SaraInterface) {
	for _, server := range servers {
		go server.StartServer()
	}
}

func shutdownServers(servers ...SaraInterface) {
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	for _, server := range servers {
		server.ShutdownServer(tc)
	}
}

type SaraInterface interface {
	StartServer()
	ShutdownServer(context.Context)
}

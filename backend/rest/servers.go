package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	photoServer := GetPhotoServer()

	calendarServer := GetCalendarServer()

	runServers(photoServer,
		calendarServer)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Println("Received terminate, graceful shutdown", sig)

	shutdownServers(photoServer)
	//shutdownServers(photoServer, calendarServer)

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

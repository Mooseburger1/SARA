package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest/handlers"
	"time"

	"golang.org/x/net/context"
)

func main() {

	l := log.New(os.Stdout, "rest-server", log.LstdFlags)

	hh := handlers.NewHelloStruct(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

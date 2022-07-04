package utils

import (
	"os"
	"os/signal"
)

type SingletonBackend struct {
	backend Backend
}

type Backend struct {
	signal chan os.Signal
}

var backend = SingletonBackend{}

func GetOsKillerListener() *chan os.Signal {
	if backend.backend.signal != nil {
		return &backend.backend.signal
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	backend.backend.signal = sigChan

	return &backend.backend.signal
}

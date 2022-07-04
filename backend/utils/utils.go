package utils

import (
	"os"
	"os/signal"

	"gopkg.in/boj/redistore.v1"
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

var store *redistore.RediStore
var err error

func GetDefaultRedisInstance() *redistore.RediStore {
	if store == nil {
		store, err = redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))
		if err != nil {
			panic(err)
		}
	}
	return store
}

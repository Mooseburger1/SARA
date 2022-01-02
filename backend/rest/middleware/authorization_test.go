package middleware

import (
	"log"
	"os"
	"testing"

	"gopkg.in/boj/redistore.v1"
)

func TestNewMiddleWare(t *testing.T) {
	mockLogger := log.New(os.Stdout, "rest-server", log.LstdFlags)
	mockStore := &redistore.RediStore{}

	actualMW := Middleware{logger: mockLogger, store: mockStore}

	mw := NewMiddleWare(mockLogger, mockStore)

	if actualMW.logger != mw.logger {
		t.Fatalf("Manually created Middleware != Constructor: Expected: %v , Got: %v", actualMW.logger, mw.logger)
	}

}

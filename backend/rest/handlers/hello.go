package handlers

import (
	"log"
	"net/http"
)

type HelloStruct struct {
	logger *log.Logger
}

func NewHelloStruct(l *log.Logger) *HelloStruct {
	return &HelloStruct{l}
}

func (h *HelloStruct) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello"))
}

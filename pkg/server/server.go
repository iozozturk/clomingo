package server

import (
	"net/http"
	"time"
)

func New(mux *http.ServeMux, host string) *http.Server {
	server := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server
}

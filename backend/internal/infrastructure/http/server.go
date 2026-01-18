package http

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		}}
}

func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

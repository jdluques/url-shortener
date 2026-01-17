package http

import "net/http"

type Server struct {
	handler http.Handler
}

func NewServer(handler http.Handler) *Server {
	return &Server{handler: handler}
}

func (server *Server) Start(addr string) error {
	return http.ListenAndServe(addr, server.handler)
}

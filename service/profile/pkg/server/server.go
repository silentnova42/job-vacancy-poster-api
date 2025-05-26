package server

import (
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(addr string, router http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		IdleTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

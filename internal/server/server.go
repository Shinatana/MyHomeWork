package server

import (
	"context"
	"net/http"
	"time"
)

const (
	DefaultServerShutdownTime = 5 * time.Second
)

type Server struct {
	server *http.Server
}

func NewHttpServer(addr string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.server.Addr
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultServerShutdownTime)
	defer cancel()

	return s.server.Shutdown(ctx)
}

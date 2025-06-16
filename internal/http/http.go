package http

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	server *http.Server
}

func NewHttpServer(handler http.Handler) *Server {
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

func (s *Server) Close(t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return s.server.Shutdown(ctx)
}

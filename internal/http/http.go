package http

import (
	"context"
	"net/http"
	"strconv"
)

type HTTP struct {
	server *http.Server
}

func NewHTTP(port int) *HTTP {
	return &HTTP{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
	}
}

func (h *HTTP) Start() error {
	return h.server.ListenAndServe()
}

func (h *HTTP) AddHandlerFunc(path string, fn http.HandlerFunc) {
	http.HandleFunc(path, fn)
}

func (h *HTTP) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

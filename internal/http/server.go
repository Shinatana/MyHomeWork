package http

import (
	"context"
	"errors"
	"log"
	"net/http"
)

const (
	DefaultServerShutdownTime 5*time.Second
)

func startServer(db *sqlDb, addr string, ctx) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/", handles.GetUserHandler(db))
	mux.HandleFunc("/users/", handles.PostUserHandler(db))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	log.Printf("Server started on :%v", addr)
	ctx, cancel := context.WithTimeout(ctx, DefaultServerShutdownTime)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
		return err
	}
	log.Printf("Server shutdown gracefully")
	return nil
}

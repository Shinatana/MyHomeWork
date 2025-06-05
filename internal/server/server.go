package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

func StartServer(addr string, handler http.Handler, shutdownTimeout time.Duration) func() error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Server started")

	return func() error {
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown server gracefully: %v", err)
			return err
		}
		log.Println("Server stopped gracefully")
		return nil
	}
}

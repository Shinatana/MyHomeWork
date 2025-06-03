package server

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"MyHomeWork/internal/handlers"
)

func RunServer(db *pgx.Conn, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", handlers.GetUserHandler(db))
	mux.HandleFunc("POST /users/{id}", handlers.PostUserHandler(db))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}

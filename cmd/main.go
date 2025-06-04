package main

import (
	"MyHomeWork/internal/handlers"
	"MyHomeWork/internal/migrate"
	"MyHomeWork/internal/server"
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultTimeoutShutdownServer = 5 * time.Second
	defaultTimeoutShutdownSQL    = 3 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutShutdownSQL)
	defer cancel()

	db, err := pgx.Connect(ctx, "postgres://Admin:123@localhost:5432/ID")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf(" Failed to close database: %v", err)
		}
	}()

	dbURL := "postgres://Admin:123@localhost:5432/ID?sslmode=disable"
	migrationsPath := "/Users/konoko/Documents/Go/MyHomeWork/internal/migrate/sqlMigrate"
	if err := migrate.RunMigrations(dbURL, migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/users/", handlers.GetUserHandler(db))
	r.HandleFunc("/users/", handlers.PostUserHandler(db))

	srv, err := server.StartServer(":8080", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if err := server.ShutdownServer(srv, defaultTimeoutShutdownServer); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}

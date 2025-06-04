package main

import (
	"MyHomeWork/internal/handlers"
	"MyHomeWork/internal/migrate"
	"MyHomeWork/internal/server"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, "postgres://Admin:123@localhost:5432/ID")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer func() {
		if err := db.Close(ctx); err != nil {
			log.Printf(" Failed to close database: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/users/", handlers.GetUserHandler(db))
	mux.HandleFunc("/users/", handlers.PostUserHandler(db))

	srv, err := server.StartServer(":8080", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if err := server.ShutdownServer(srv, 5*time.Second); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	dbURL := "postgres://Admin:123@localhost:5432/ID?sslmode=disable"
	migrationsPath := "/Users/konoko/Documents/Go/Homework3/study/migrations"
	if err := migrate.RunMigrations(dbURL, migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

package main

import (
	"MyHomeWork/internal/db"
	"MyHomeWork/internal/db/migrate"
	"MyHomeWork/internal/handler"
	"MyHomeWork/internal/server"
	"context"
	"log"
	"net/http"
)

func main() {
	dbURL := "postgres://Admin:123@localhost:5432/ID"

	err := migrate.RunMigrations(dbURL)
	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	ctx := context.Background()

	pool, err := db.ConnectDB(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	srv := server.NewHttpServer(":8080", mux)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/", handler.GetHandler(pool))
	mux.HandleFunc("/", handler.PostHandler(pool))

}

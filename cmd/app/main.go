package main

import (
	"MyHomeWork/internal/db/migrate"
	"MyHomeWork/internal/server"
	"flag"
	"log"
	"net/http"
)

func main() {
	dbURL := "postgres://Admin:123@localhost:5432/ID"
	direction := flag.String("migrate", "up", "Migration direction: up or down")
	err := migrate.RunMigrations(direction, dbURL)
	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	mux := http.NewServeMux()

	srv := server.NewHttpServer(":8080", mux)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/", GetHandler)
	mux.HandleFunc("/", PostHandler)
}

package main

import (
	"MyHomeWork/internal/migrate"
	"MyHomeWork/internal/server"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	db, err := pgx.Connect(context.Background(), "postgres://Admin:123@localhost:5432/ID")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	dbURL := "postgres://Admin:123@localhost:5432/ID?sslmode=disable"
	migrationsPath := "/Users/konoko/Documents/Go/Homework3/study/migrations"

	if err := migrate.RunMigrations(dbURL, migrationsPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := server.RunServer(db, ":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

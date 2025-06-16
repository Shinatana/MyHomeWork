package main

import (
	"MyHomeWork/internal/db"
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	dbURL := "postgres://Admin:123@localhost:5432/ID"
	ctx := context.Background()
	conn, ctxBd := db.Start(dbURL, ctx)
	defer db.Stop(conn, ctxBd)

	//проверяем миграции
	if err := migrations.RunMigrations("file://migrations", dbURL); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	//запускаем сервер
	go func() {
		if err := server.RunServer(conn); err != nil {
			log.Fatalf("Server run error: %v", err)
		}
	}()
	// ожидаем сигнала к завершению

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("App shutting down")
}

//	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutShutdownSQL)
//	defer cancel()
//
//	db, err := pgx.Connect(ctx, "postgres://Admin:123@localhost:5432/ID")
//	if err != nil {
//		log.Fatalf("Unable to connect to database: %v", err)
//	}
//
//	defer func() {
//		if err := db.Close(ctx); err != nil {
//			log.Printf(" Failed to close database: %v", err)
//		}
//	}()
//
//	dbURL := "postgres://Admin:123@localhost:5432/ID?sslmode=disable"
//	migrationsPath := "/Users/konoko/Documents/Go/MyHomeWork/internal/migrate/sqlMigrate"
//	if err := migrate.RunMigrations(dbURL, migrationsPath); err != nil {
//		log.Fatalf("Migration failed: %v", err)
//	}
//
//	r := mux.NewRouter()
//	r.HandleFunc("/users/", handlers.GetUserHandler(db))
//	r.HandleFunc("/users/", handlers.PostUserHandler(db))
//
//	srv, err := server.StartServer(":8080", r)
//	if err != nil {
//		log.Fatalf("Failed to start server: %v", err)
//	}
//
//	stop := make(chan os.Signal, 1)
//	signal.Notify(stop, os.Interrupt)
//	<-stop
//
//	if err := server.ShutdownServer(srv, defaultTimeoutShutdownServer); err != nil {
//		log.Fatalf("Failed to shutdown server: %v", err)
//	}
//}
//
//func GetUserHandler

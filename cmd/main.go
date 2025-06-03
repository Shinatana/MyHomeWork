package main

import (
	"MyHomeWork/migrate"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type errResponse struct {
	Error string `json:"error"`
}

type userIdResponse struct {
	Id string `json:"id"`
}

func getUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := context.Background()

		var userID bool
		// найди правильную ошибку
		err := db.QueryRow(ctx, "SELECT id FROM users WHERE id=$1", id).Scan(&userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusConflict)
				if err := json.NewEncoder(w).Encode(errResponse{"ID already exists"}); err != nil {
					log.Printf("Error encoding JSON response: %v", err)
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(errResponse{"Internal server error"}); err != nil {
				log.Printf("Error encoding JSON response: %v", err)

			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// верни нормальные данные
		if err := json.NewEncoder(w).Encode(userIdResponse{Id: id}); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}
	}
}

func postUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := context.Background()
		_, err := db.Exec(ctx, "INSERT INTO users (id) VALUES ($1)", id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				// Дубликат ключа — возвращаем 409 Conflict
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				if err := json.NewEncoder(w).Encode(errResponse{"ID already exists"}); err != nil {
					log.Printf("Error encoding JSON response: %v", err)
				}
				return
			}

			// Другие ошибки
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(errResponse{"Error adding ID"}); err != nil {
				log.Printf("Error encoding JSON response: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userIdResponse{Id: id}); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}
	}
}

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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", getUserHandler(db))
	mux.HandleFunc("POST /users/{id}", postUserHandler(db))

	srv := &http.Server{
		Addr:    ":8080",
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
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server stopped gracefully")
}

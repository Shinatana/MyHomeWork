package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

type userIdResponse struct {
	Id string `json:"id"`
}

func GetUserHandler(db *pgx.Conn) http.HandlerFunc {
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

func PostUserHandler(db *pgx.Conn) http.HandlerFunc {
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

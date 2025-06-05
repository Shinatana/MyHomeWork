package handlers

import (
	"bytes"
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

func writeJSON(w http.ResponseWriter, status int, data interface{}) {

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(data); err != nil {
		http.Error(w, "Error on server", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(buff.Bytes()); err != nil {
		log.Printf("Write error: %v", err)
	}
}

func GetUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := r.Context()

		var userID string
		err := db.QueryRow(ctx, "SELECT id FROM users WHERE id=$1", id).Scan(&userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writeJSON(w, http.StatusConflict, errResponse{"ID not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, errResponse{"Internal server error"})
			return
		}

		writeJSON(w, http.StatusOK, userIdResponse{Id: id})
	}
}

func PostUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		ctx := r.Context()

		_, err := db.Exec(ctx, "INSERT INTO users (id) VALUES ($1)", id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				writeJSON(w, http.StatusConflict, errResponse{"ID already exists"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, errResponse{"Error adding ID"})
			return
		}
		writeJSON(w, http.StatusCreated, userIdResponse{Id: id})
	}
}

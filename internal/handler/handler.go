package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"strconv"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "id is required",
			})
			log.Println("id is required")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "id is invalid",
			})
			log.Println("id is invalid")
			return
		}
		var exists bool
		err = pool.QueryRow(r.Context(), "SELECT EXISTS (SELECT 1 FROM your_table WHERE id=$1)", id).Scan(&exists)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "internal server error",
			})
			log.Println("id is invalid")
			return
		}
		if !exists {
			writeJSON(w, http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: "id not found",
			})
			log.Println("id not found")
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{
			Status:  "success",
			Message: "id exists",
		})
	}
}

func PostHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "id is required",
			})
			log.Println("id is required")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "id is invalid",
			})
			log.Println("id is invalid")
			return
		}
		var exists bool
		err = pool.QueryRow(r.Context(), "SELECT EXISTS (SELECT 1 FROM your_table WHERE id=$1)", id).Scan(&exists)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "internal server error",
			})
			log.Printf("db query error: %v", err)
			return
		}

		if exists {
			writeJSON(w, http.StatusConflict, APIResponse{
				Status:  "error",
				Message: "id is exists",
			})
			log.Println("id is already exists")
			return
		}
		_, err = pool.Exec(r.Context(), "INSERT INTO your_table(id) VALUES($1)", id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "internal server error",
			})
			log.Println("failed to insert id:", err)
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{
			Status:  "success",
			Message: "id created",
		})
	}
}

package user

import (
	"MyHomework/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"MyHomework/internal/repo"
)

const (
	defaultReqIDHeader = "X-Request-ID"
)

type User struct {
	db  repo.UsersDB
	log *slog.Logger
}

func NewUser(db repo.UsersDB, log *slog.Logger) *User {
	return &User{db: db, log: log}
}

func (u *User) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		userID := r.PathValue("id")
		w.Header().Set(defaultReqIDHeader, reqID)
		lg := u.log.With(
			"request_id", reqID,
			"user_id", userID,
		)

		if userID == "" {
			lg.Error("user id is empty")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		convID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			lg.Error("user id is not a number")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if convID < 0 {
			lg.Error("user id is negative")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		lg = lg.With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
		lg.Debug("request received")

		result, err := u.db.GetUserByID(r.Context(), convID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lg.Error("user not found")
				w.WriteHeader(http.StatusNotFound)
				return
			}

			lg.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		lg.Debug("user found", "user", result)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(result); err != nil {
			lg.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (u *User) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		userID := r.PathValue("id")
		w.Header().Set(defaultReqIDHeader, reqID)
		lg := u.log.With(
			"request_id", reqID,
			"user_id", userID,
		)

		var createUser models.PostUserRequest
		if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
			lg.Error("Invalid input")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if createUser.Name == "" {
			lg.Error("An empty string was passed instead of a name")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if createUser.Age < 0 {
			lg.Error("Age is negative value")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		id, err := u.db.UpsertUser(r.Context(), createUser.Name, createUser.Age)
		if err != nil {
			lg.Error("failed to write to database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if id == 0 {
			lg.Error("failed to write to database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := models.CreateUserResponse{
			ID:      id,
			Message: "User successfully created",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			lg.Error("failed to encode response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

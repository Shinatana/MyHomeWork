package user

import (
	"MyHomework/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
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

		if err = json.NewEncoder(w).Encode(result); err != nil {
			lg.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
			lg.Error("Invalid input", "error", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(&createUser); err != nil {
			lg.Error("Validation failed", "error", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		id, err := u.db.UpsertUser(r.Context(), createUser.Name, createUser.Age)
		if err != nil {
			lg.Error("failed to write to database", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := models.CreateUserResponse{
			ID: id,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			lg.Error("failed to encode response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

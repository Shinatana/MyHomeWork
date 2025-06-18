package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Shinatana/MyHomeWork/internal/models"
)

const (
	getUserByIDQuery = "SELECT id, name, age FROM users WHERE id = $1"
	upsertUserQuery  = "INSERT INTO users (name, age) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET age = EXCLUDED.age"
)

type UsersDB interface {
	GetUserByID(ctx context.Context, id int64) (*models.GetUserResponse, error)
	UpsertUser(ctx context.Context, name string, age int) error
	Close() error
}

type DB struct {
	db *sql.DB
}

func NewDB(dsn string) (UsersDB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (*models.GetUserResponse, error) {
	rows, err := db.db.QueryContext(ctx, getUserByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query db: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var user models.GetUserResponse
	err = rows.Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &user, nil
}

func (db *DB) UpsertUser(ctx context.Context, name string, age int) error {
	_, err := db.db.ExecContext(ctx, upsertUserQuery, name, age)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const defaultTimeout = 5 * time.Second

func ConnectDB(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	ctxBD, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctxBD, dbURL)
	if err != nil {
		return nil, err
	}

	// Пробуем сделать ping (optional)
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("Connected to DB")
	return pool, nil
}

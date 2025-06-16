package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

const defaultTimeout = 5 * time.Second

func Start(dbURL string, ctx context.Context) (*pgx.Conn, context.Context) {
	// Создаём контекст с таймаутом для подключения к базе
	ctxBD, cancel := context.WithTimeout(ctx, defaultTimeout)

	// Важно: отменять таймаут когда он не нужен
	// Но если возвращаем ctx, отмену должен делать вызывающий
	// чтобы таймаут был актуален при использовании

	// Подключаемся к базе
	conn, err := pgx.Connect(ctxBD, dbURL)
	if err != nil {
		log.Printf("Connect to %s failed: %v", dbURL, err)
		cancel()
		return nil, nil
	}

	return conn, ctxBD
}

func Stop(ctx context.Context, conn *pgx.Conn, cancel context.CancelFunc) {
	// Закрываем соединение с базой
	if conn != nil {
		_ = conn.Close(ctx)
	}

	// Отменяем контекст (если был таймаут/отмена)
	if cancel != nil {
		cancel()
	}
}

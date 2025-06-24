package main

import (
	"MyHomework/internal/conf"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalHttp "MyHomework/internal/http"
	userHandlers "MyHomework/internal/http/handlers/user"
	"MyHomework/internal/log"
	"MyHomework/internal/repo"
)

const (
	defaultHttpShutdownTimeout = 5 * time.Second
)

func main() {
	// init cfg
	cfg, err := conf.NewCfg("pkg/conf/config.yaml")
	if err != nil {
		log.NewLog("", "").Error(err.Error())
		os.Exit(1)
	}
	// create logger
	lg := log.NewLog(cfg.LogFormat, cfg.LogLevel)
	lg.Debug(
		"cfg initialized",
		"log_format", cfg.LogFormat,
		"log_level", cfg.LogLevel,
		"db_dsn", "********",
		"http_port", cfg.HttpPort,
	)

	// connect to db
	db, err := repo.NewDB(cfg.DSN)
	if err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = db.Close()
		lg.Info("db closed")
	}()
	lg.Info("db initialized")

	// init http server
	server := internalHttp.NewHTTP(cfg.HttpPort)

	user := userHandlers.NewUser(db, lg)
	server.AddHandlerFunc("GET /user/{id}", user.Get())
	server.AddHandlerFunc("POST /user", user.Post())

	// start http server
	go func() {
		lg.Info("http server started")
		err := server.Start()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				lg.Error(err.Error())
			}
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultHttpShutdownTimeout)
		defer cancel()
		_ = server.Stop(ctx)
		lg.Info("http server stopped")
	}()

	// graceful stop and quit
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(quit)
	<-quit
}

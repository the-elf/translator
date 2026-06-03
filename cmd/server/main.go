package main

import (
	"context"
	"log/slog"
	"os"
	"time"
	"translator/internal/handler"
	"translator/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

const databaseURLKey = "DATABASE_URL"

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug})))

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env not found, using environment")
	}

	databaseURL, err := util.RequireEnv(databaseURLKey)
	if err != nil {
		fatal(err.Error())
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fatal("failed to connect to database", "err", err)
	}
	defer pool.Close()

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = pool.Ping(timeout)
	if err != nil {
		fatal("failed to connect to database", "err", err)
	}

	th, err := handler.NewTelegramHandler(pool)
	if err != nil {
		fatal("failed to create telegram handler", "err", err)
	}

	th.Start()
}

func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

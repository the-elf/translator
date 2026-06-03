package main

import (
	"context"
	"log"
	"os"
	"time"
	"translator/internal/handler"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const databaseURLKey = "DATABASE_URL"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using environment")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv(databaseURLKey))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = pool.Ping(timeout)
	if err != nil {
		panic(err)
	}

	th, err := handler.NewTelegramHandler(pool)
	if err != nil {
		panic(err)
	}

	th.Start()
}

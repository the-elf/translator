package main

import (
	"log"
	"translator/internal/handler"
	"translator/internal/repository"
	"translator/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using environment")
	}

	cr := repository.NewClientRepositoryMemory()
	ts := service.NewTelegramService(cr)
	th, err := handler.NewTelegramHandler(ts)
	if err != nil {
		panic(err)
	}

	th.Start()
}

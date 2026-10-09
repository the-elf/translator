package config

import (
	"log/slog"
	"time"
)

type Config struct {
	TelegramToken   string
	DatabaseURL     string
	HTTPAddr        string
	HandlerTimeout  time.Duration
	ShutdownTimeout time.Duration
	OpenAI          OpenAI
	Log             Log
}

type OpenAI struct {
	APIKey         string
	Model          string
	AttemptTimeout time.Duration
	MaxAttempts    int
	MaxConcurrent  int
}

type Log struct {
	Level  slog.Level
	Format string // "text" | "json"
}

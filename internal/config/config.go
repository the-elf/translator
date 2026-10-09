package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	TelegramToken   string        `env:"TELEGRAM_BOT_TOKEN,required"`
	DatabaseURL     string        `env:"DATABASE_URL,required"`
	HTTPAddr        string        `env:"HTTP_ADDR,default=:8080"`
	HandlerTimeout  time.Duration `env:"HANDLER_TIMEOUT,default=60s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=70s"`
	OpenAI          OpenAI
	Log             Log
}

type OpenAI struct {
	APIKey         string        `env:"OPENAI_API_KEY,required"`
	Model          string        `env:"GPT_MODEL_VERSION,required"`
	AttemptTimeout time.Duration `env:"OPENAI_ATTEMPT_TIMEOUT,default=20s"`

	MaxAttempts   int
	MaxConcurrent int
}

type Log struct {
	Level  slog.Level `env:"LOG_LEVEL,default=info"`
	Format string     `env:"LOG_FORMAT,default=text"` // "text" | "json"
}

type nonEmptyLookuper struct {
	inner envconfig.Lookuper
}

func (l *nonEmptyLookuper) Lookup(key string) (string, bool) {
	val, found := l.inner.Lookup(key)
	if !found || strings.TrimSpace(val) == "" {
		return "", false
	}
	return val, true
}

func Load(ctx context.Context) (Config, error) {
	var cfg Config
	err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   &cfg,
		Lookuper: &nonEmptyLookuper{inner: envconfig.OsLookuper()},
	})
	if err != nil {
		return Config{}, fmt.Errorf("process env: %w", err)
	}
	if cfg.ShutdownTimeout <= cfg.HandlerTimeout {
		return Config{}, errors.New("shutdown timeout must be greater than handler timeout")
	}
	cfg.OpenAI.MaxAttempts = 3
	cfg.OpenAI.MaxConcurrent = 10
	return cfg, nil
}

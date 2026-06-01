package handler

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"translator/internal/service"

	"gopkg.in/telebot.v3"
)

const telegramBotTokenKey = "TELEGRAM_BOT_TOKEN"

type TelegramHandler struct {
	bot *telebot.Bot
	svc *service.TelegramService
}

func NewTelegramHandler(svc *service.TelegramService) (*TelegramHandler, error) {
	bot, err := newTelegramBot()
	if err != nil {
		return nil, err
	}

	return &TelegramHandler{bot: bot, svc: svc}, nil
}

func (t *TelegramHandler) Start() {
	log.Println("Bot started")
	t.registerHandlers()
	go t.gracefulShutdown()
	t.bot.Start()
}

func newTelegramBot() (*telebot.Bot, error) {
	token := os.Getenv(telegramBotTokenKey)
	if token == "" {
		return nil, fmt.Errorf("failed to create a new bot: %s not set", telegramBotTokenKey)
	}
	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	return telebot.NewBot(settings)
}

func (t *TelegramHandler) registerHandlers() {
	t.bot.Handle("/start", t.svc.Register)
	t.bot.Handle(telebot.OnText, t.svc.Translate)
	t.bot.Handle("/setLang", t.svc.SetLangCommandButtons)
	t.bot.Handle(&telebot.InlineButton{Unique: "set_lang"}, t.svc.SetLangCommand)
}

func (t *TelegramHandler) gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Printf("\rGot %s signal. Shutting down the bot\n", sig)
	t.bot.Stop()
}

package handler

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"translator/internal/model"
	"translator/internal/service"
	"translator/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v3"
)

const telegramBotTokenKey = "TELEGRAM_BOT_TOKEN"
const editDelay = 300 * time.Millisecond

type TelegramHandler struct {
	bot            *telebot.Bot
	userSvc        *service.UserService
	langSvc        *service.LanguageService
	mtSvc          *service.MessageTemplateService
	translationSvc *service.TranslationService
	btnSvc         *service.ButtonService
}

func NewTelegramHandler(db *pgxpool.Pool) (*TelegramHandler, error) {
	bot, err := newTelegramBot()
	if err != nil {
		return nil, err
	}

	translationService, err := service.NewTranslationService()
	if err != nil {
		return nil, err
	}

	return &TelegramHandler{
		bot:            bot,
		userSvc:        service.NewUserService(db),
		langSvc:        service.NewLanguageService(db),
		mtSvc:          service.NewMessageTemplate(db),
		translationSvc: translationService,
		btnSvc:         service.NewButtonService(db),
	}, nil
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
		Token:   token,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, context telebot.Context) {},
	}
	return telebot.NewBot(settings)
}

func (t *TelegramHandler) gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Printf("Got %s signal. Shutting down the bot\n", sig)
	t.bot.Stop()
}

func (t *TelegramHandler) registerHandlers() {
	t.bot.Handle("/start", t.handleRegister)
	t.bot.Handle(telebot.OnText, t.handleTranslate)
	t.bot.Handle("/setLang", t.handleSetLangCommandButtons)
	t.bot.Handle(&telebot.InlineButton{Unique: "set_lang"}, t.handleSetLangCommand)
}

func (t *TelegramHandler) handleRegister(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err == nil {
		return t.sendMessageTemplateText(ctx, model.AlreadyRegisteredMsg, user.Language)
	}

	if errors.Is(err, model.ErrUserNotFound) {
		defaultLang, err := t.langSvc.GetDefaultLang()
		if err != nil {
			return sendUnexpectedError(ctx, err)
		}

		user := model.User{ChatID: chatID, Language: defaultLang}
		err = t.userSvc.SaveUser(user)
		if err != nil {
			return sendUnexpectedError(ctx, err)
		}

		log.Printf("user registered: chatID=%d, lang=%s", chatID, defaultLang.Code)
		return t.sendMessageTemplateText(ctx, model.SuccessfullyRegisteredMsg, user.Language)
	}

	return sendUnexpectedError(ctx, err)
}

func (t *TelegramHandler) handleTranslate(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return t.sendNotRegisteredMsg(ctx)
		}
		return sendUnexpectedError(ctx, err)
	}

	text := ctx.Message().Text
	placeHolderTemplate, err := t.mtSvc.GetByCodeAndLang(model.TranslationPlaceholder, user.Language)
	if err != nil {
		logMsg := fmt.Sprintf("failed to get %s message template", model.TranslationPlaceholder)
		return sendUnexpectedErrorWithLogMessage(ctx, err, logMsg)
	}

	sent, err := ctx.Bot().Send(ctx.Recipient(), placeHolderTemplate.Text)
	if err != nil {
		logMsg := "failed to send translation placeholder to user"
		return sendUnexpectedErrorWithLogMessage(ctx, err, logMsg)
	}

	promptTemplate, err := t.mtSvc.GetByCodeAndLang(model.TranslationPromptText, user.Language)
	if err != nil {
		logMsg := fmt.Sprintf("failed to get %s message template", model.TranslationPromptText)
		return sendUnexpectedErrorWithLogMessage(ctx, err, logMsg)
	}

	translation, err := t.translationSvc.Translate(text, promptTemplate)
	if err != nil {
		log.Printf("translation failed: chatID=%d, err=%v", chatID, err)
		return t.sendMessageTemplateText(ctx, model.TranslationErrorMsg, user.Language)
	}

	//rate-limit outgoing edits
	time.Sleep(editDelay)
	_, err = ctx.Bot().Edit(sent, translation)
	if err != nil {
		log.Printf("failed to edit message: chatID=%d, err=%v", chatID, err)
	}

	return err
}

func (t *TelegramHandler) handleSetLangCommandButtons(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return t.sendNotRegisteredMsg(ctx)
		}
		return sendUnexpectedError(ctx, err)
	}

	setLangButtons, err := t.btnSvc.GetGroupByNameAndLang(model.SetLangButtonGroup, user.Language)
	if err != nil {
		return sendUnexpectedError(ctx, err)
	}

	menu := &telebot.ReplyMarkup{}
	row := telebot.Row{}
	for _, button := range setLangButtons {
		row = append(row, menu.Data(button.Title, button.GroupName, button.Data))
	}

	menu.Inline(row)

	template, tErr := t.mtSvc.GetByCodeAndLang(model.PreferredLanguageMsg, user.Language)
	if tErr != nil {
		return sendUnexpectedError(ctx, tErr)
	}

	return sendMessage(ctx, template.Text, menu)
}

func (t *TelegramHandler) handleSetLangCommand(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return t.sendNotRegisteredMsg(ctx)
		}
		return sendUnexpectedError(ctx, err)
	}

	lang, err := t.langSvc.GetLangByCode(model.LanguageCode(ctx.Data()))
	if err != nil {
		if errors.Is(err, model.ErrUnsupportedLanguage) {
			return t.sendMessageTemplateText(ctx, model.UnsupportedLanguageMsg, user.Language)
		}

		return sendUnexpectedError(ctx, err)
	}

	err = t.userSvc.SetLang(user, lang)
	if err != nil {
		return sendUnexpectedError(ctx, err)
	}

	log.Printf("language changed: chatID=%d, lang=%s", chatID, lang.Code)
	return t.sendMessageTemplateText(ctx, model.LanguageSetMsg, lang, lang.Code)
}

func sendUnexpectedErrorWithLogMessage(ctx telebot.Context, err error, logMsg string) error {
	logMsg = util.Ternary(
		logMsg == "",
		fmt.Sprintf("unexpected error: chatID=%d, err=%v", ctx.Sender().ID, err),
		fmt.Sprintf("%s: chatID=%d, err=%v", logMsg, ctx.Sender().ID, err),
	)
	log.Print(logMsg)

	_ = sendMessage(ctx, "⚠️ Произошла неожиданная ошибка. Повторите попытку позже.")
	return err
}

func sendMessage(ctx telebot.Context, text string, args ...any) error {
	sendErr := ctx.Send(text, args...)
	if sendErr != nil {
		log.Printf("failed to send a message to user: %v", sendErr)
	}

	return sendErr
}

func sendUnexpectedError(ctx telebot.Context, err error) error {
	return sendUnexpectedErrorWithLogMessage(ctx, err, "")
}

func (t *TelegramHandler) sendNotRegisteredMsg(ctx telebot.Context) error {
	log.Printf("user with chatID=%d not registered", ctx.Sender().ID)

	lang, err := t.langSvc.GetDefaultLang()
	if err != nil {
		return sendUnexpectedError(ctx, err)
	}

	return t.sendMessageTemplateText(ctx, model.NotRegisteredMsg, lang)
}

func (t *TelegramHandler) sendMessageTemplateText(ctx telebot.Context, code model.MessageTemplateCode, lang model.Language, args ...any) error {
	template, err := t.mtSvc.GetByCodeAndLang(code, lang)
	if err != nil {
		return sendUnexpectedError(ctx, err)
	}

	text := util.Ternary(
		len(args) > 0,
		fmt.Sprintf(template.Text, args...),
		template.Text,
	)

	return sendMessage(ctx, text)
}

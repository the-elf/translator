package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"translator/internal/model"
	"translator/internal/service"
	"translator/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"gopkg.in/telebot.v3"
)

const telegramBotTokenKey = "TELEGRAM_BOT_TOKEN"
const editDelay = 300 * time.Millisecond
const defaultUnexpectedErrMsg = "⚠️ Произошла неожиданная ошибка. Повторите попытку позже."

const startCommand = "/start"
const setLangCommand = "/setLang"

var setLangInlineButton = &telebot.InlineButton{Unique: string(model.SetLangButtonGroup)}
var geoTranslitButton = &telebot.InlineButton{Unique: string(model.GeoTranslitButtonGroup)}

type TelegramHandler struct {
	bot            *telebot.Bot
	translitCache  *cache.Cache
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
		translitCache:  cache.New(5*time.Minute, 10*time.Minute),
		userSvc:        service.NewUserService(db),
		langSvc:        service.NewLanguageService(db),
		mtSvc:          service.NewMessageTemplateService(db),
		translationSvc: translationService,
		btnSvc:         service.NewButtonService(db),
	}, nil
}

func (t *TelegramHandler) Start() {
	slog.Info("Bot started")
	t.registerHandlers()
	go t.gracefulShutdown()
	t.bot.Start()
}

func newTelegramBot() (*telebot.Bot, error) {
	token, err := util.RequireEnv(telegramBotTokenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %v", err)
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
	slog.Info("shutting down", "signal", sig)
	t.bot.Stop()
}

func (t *TelegramHandler) registerHandlers() {
	t.bot.Handle(startCommand, t.handleRegister)
	t.bot.Handle(telebot.OnText, t.handleTranslate)
	t.bot.Handle(setLangCommand, t.handleSetLangCommand)
	t.bot.Handle(setLangInlineButton, t.handleSetLangButton)
	t.bot.Handle(geoTranslitButton, t.handleGeoTranslitButton)
}

func (t *TelegramHandler) handleRegister(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err == nil {
		return t.sendMessageTemplateText(ctx, model.AlreadyRegisteredMsg, user.Language)
	}

	if errors.Is(err, model.ErrUserNotFound) {
		defaultLang, lErr := t.langSvc.GetDefaultLang()
		if lErr != nil {
			return t.sendDefaultUnexpectedError(ctx, lErr)
		}

		user = model.User{ChatID: chatID, Language: defaultLang}
		err = t.userSvc.SaveUser(user)
		if err != nil {
			return t.sendDefaultUnexpectedError(ctx, err)
		}

		slog.Info("user registered", "chatID", chatID, "lang", defaultLang.Code)
		return t.sendMessageTemplateText(ctx, model.SuccessfullyRegisteredMsg, user.Language)
	}

	return t.sendDefaultUnexpectedError(ctx, err)
}

func (t *TelegramHandler) handleTranslate(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return t.sendNotRegisteredMsg(ctx)
		}
		return t.sendDefaultUnexpectedError(ctx, err)
	}

	text := ctx.Message().Text
	placeHolderTemplate, err := t.mtSvc.GetByCodeAndLang(model.TranslationPlaceholder, user.Language)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, model.TranslationPlaceholder, user.Language)
	}

	sent, err := ctx.Bot().Send(ctx.Recipient(), placeHolderTemplate.Text)
	if err != nil {
		return t.sendError(
			ctx,
			err,
			"failed to send translation placeholder to user",
			user.Language,
		)
	}

	promptTemplate, err := t.mtSvc.GetByCodeAndLang(model.TranslationPromptText, user.Language)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, model.TranslationPromptText, user.Language)
	}

	translation, err := t.translationSvc.Translate(text, promptTemplate)
	if err != nil {
		if errors.Is(err, service.ErrNotGeorgianText) {
			return t.processNonGeorgianText(ctx, sent, user.Language)
		}
		slog.Error("translation failed", "chatID", chatID, "err", err)
		return t.sendMessageTemplateText(ctx, model.TranslationErrorMsg, user.Language)
	}

	return editMessage(ctx, sent, translation)
}

func (t *TelegramHandler) handleSetLangCommand(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		return t.processGetUserErr(ctx, err)
	}

	menu, err := t.createMenu(model.SetLangButtonGroup, user.Language)
	if err != nil {
		return t.sendMenuError(ctx, err, user.Language, model.SetLangButtonGroup)
	}

	template, err := t.mtSvc.GetByCodeAndLang(model.PreferredLanguageMsg, user.Language)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, model.PreferredLanguageMsg, user.Language)
	}

	return sendMessage(ctx, template.Text, menu)
}

func (t *TelegramHandler) handleSetLangButton(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		return t.processGetUserErr(ctx, err)
	}

	lang, err := t.langSvc.GetLangByCode(model.LanguageCode(ctx.Data()))
	if err != nil {
		if errors.Is(err, model.ErrUnsupportedLanguage) {
			return t.sendMessageTemplateText(ctx, model.UnsupportedLanguageMsg, user.Language)
		}

		return t.sendError(ctx, err, "failed to get language", user.Language)
	}

	err = t.userSvc.SetLang(user, lang)
	if err != nil {
		return t.sendError(ctx, err, "failed to set language", user.Language)
	}

	slog.Info("language changed", "chatID", chatID, "lang", lang.Code)
	return t.editMessageTemplateText(ctx, ctx.Message(), model.LanguageSetMsg, lang, lang.Code)
}

func (t *TelegramHandler) handleGeoTranslitButton(ctx telebot.Context) error {
	chatID := ctx.Sender().ID
	user, err := t.userSvc.GetUser(chatID)
	if err != nil {
		return t.processGetUserErr(ctx, err)
	}

	switch model.ButtonData(ctx.Data()) {
	case model.GeoTranslitButtonYes:
		text, ok := t.getFromTranslitCache(chatID, ctx.Message().ID)
		if !ok {
			// todo добавить сообщение для пользователя в чате
			return t.sendError(
				ctx,
				errors.New("cache miss"),
				"translit cache lookup failed",
				user.Language,
			)
		}

		promptTemplate, err := t.mtSvc.GetByCodeAndLang(model.TranslationPromptText, user.Language)
		if err != nil {
			return t.sendMessageTemplateError(ctx, err, model.TranslationPromptText, user.Language)
		}

		translation, err := t.translationSvc.Translate(util.ToGeorgian(text), promptTemplate)
		if err != nil {
			slog.Error("translation failed", "chatID", chatID, "err", err)
			if errors.Is(err, service.ErrNotGeorgianText) {
				return t.editMessageTemplateText(ctx, ctx.Message(), model.NotGeorgianTextErrorMsg, user.Language)
			}
			return t.sendMessageTemplateText(ctx, model.TranslationErrorMsg, user.Language)
		}

		return editMessage(ctx, ctx.Message(), translation)
	case model.GeoTranslitButtonNo:
		return t.editMessageTemplateText(ctx, ctx.Message(), model.NotGeorgianTextErrorMsg, user.Language)
	}

	return nil
}

func (t *TelegramHandler) processGetUserErr(ctx telebot.Context, err error) error {
	if errors.Is(err, model.ErrUserNotFound) {
		return t.sendNotRegisteredMsg(ctx)
	}
	return t.sendDefaultUnexpectedError(ctx, err)
}

func (t *TelegramHandler) processNonGeorgianText(ctx telebot.Context, sent *telebot.Message, lang model.Language) error {
	menu, err := t.createMenu(model.GeoTranslitButtonGroup, lang)
	if err != nil {
		return t.sendMenuError(ctx, err, lang, model.GeoTranslitButtonGroup)
	}

	template, err := t.mtSvc.GetByCodeAndLang(model.GeorgianTranslitMsg, lang)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, model.GeorgianTranslitMsg, lang)
	}

	t.addToTranslitCache(ctx.Sender().ID, sent.ID, ctx.Message().Text)
	return editMessage(ctx, sent, template.Text, menu)
}

func (t *TelegramHandler) createMenu(buttonGroup model.ButtonGroupName, lang model.Language) (*telebot.ReplyMarkup, error) {
	buttons, err := t.btnSvc.GetGroupByNameAndLang(buttonGroup, lang)
	if err != nil {
		return nil, err
	}

	menu := &telebot.ReplyMarkup{}
	row := telebot.Row{}
	for _, button := range buttons {
		row = append(row, menu.Data(button.Title, string(button.GroupName), string(button.Data)))
	}

	menu.Inline(row)

	return menu, nil
}

func (t *TelegramHandler) sendNotRegisteredMsg(ctx telebot.Context) error {
	slog.Warn("user is not registered", "chatID", ctx.Sender().ID)

	lang, err := t.langSvc.GetDefaultLang()
	if err != nil {
		return t.sendDefaultUnexpectedError(ctx, err)
	}

	return t.sendMessageTemplateText(ctx, model.NotRegisteredMsg, lang)
}

func (t *TelegramHandler) editMessageTemplateText(ctx telebot.Context, message *telebot.Message, code model.MessageTemplateCode, lang model.Language, args ...any) error {
	template, err := t.mtSvc.GetByCodeAndLang(code, lang)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, code, lang)
	}

	var text string
	if len(args) > 0 {
		text = fmt.Sprintf(template.Text, args...)
	} else {
		text = template.Text
	}

	return editMessage(ctx, message, text)
}

func (t *TelegramHandler) sendMessageTemplateText(ctx telebot.Context, code model.MessageTemplateCode, lang model.Language, args ...any) error {
	template, err := t.mtSvc.GetByCodeAndLang(code, lang)
	if err != nil {
		return t.sendMessageTemplateError(ctx, err, code, lang)
	}

	var text string
	if len(args) > 0 {
		text = fmt.Sprintf(template.Text, args...)
	} else {
		text = template.Text
	}

	return sendMessage(ctx, text)
}

func (t *TelegramHandler) addToTranslitCache(chatID int64, messageID int, text string) {
	key := translitCacheKey(chatID, messageID)
	t.translitCache.SetDefault(key, text)
}

func (t *TelegramHandler) getFromTranslitCache(chatID int64, messageID int) (string, bool) {
	key := translitCacheKey(chatID, messageID)
	text, ok := t.translitCache.Get(key)
	if !ok {
		return "", false
	}

	return text.(string), ok
}

func (t *TelegramHandler) sendError(ctx telebot.Context, err error, logMsg string, lang model.Language, args ...any) error {
	if logMsg == "" {
		logMsg = "unexpected error"
	}
	args = append(args, "lang", lang.Code, "chatID", ctx.Sender().ID, "err", err)
	slog.Error(logMsg, args...)

	var errMsg string
	if lang == model.EmptyLang() {
		errMsg = defaultUnexpectedErrMsg
	} else {
		template, tErr := t.mtSvc.GetByCodeAndLang(model.UnexpectedErrMsg, lang)
		if tErr != nil {
			_ = sendMessage(ctx, defaultUnexpectedErrMsg)
			return fmt.Errorf("errors: [%w, %v]", err, tErr)
		}
		errMsg = template.Text
	}

	_ = sendMessage(ctx, errMsg)
	return err
}

func (t *TelegramHandler) sendMenuError(ctx telebot.Context, err error, lang model.Language, menuName model.ButtonGroupName) error {
	return t.sendError(
		ctx,
		err,
		"failed to create menu",
		lang,
		"menuName",
		menuName,
	)
}

func (t *TelegramHandler) sendDefaultUnexpectedError(ctx telebot.Context, err error) error {
	return t.sendError(ctx, err, "", model.EmptyLang())
}

func (t *TelegramHandler) sendMessageTemplateError(ctx telebot.Context, err error, code model.MessageTemplateCode, lang model.Language) error {
	return t.sendError(
		ctx,
		err,
		"failed to get message template",
		lang,
		"code", code,
	)
}

func translitCacheKey(chatID int64, messageID int) string {
	return fmt.Sprintf("%d_%d", chatID, messageID)
}

func sendMessage(ctx telebot.Context, text string, args ...any) error {
	sendErr := ctx.Send(text, args...)
	if sendErr != nil {
		slog.Error("failed to send a message to user", "err", sendErr)
	}

	return sendErr
}

func editMessage(ctx telebot.Context, message *telebot.Message, text string, args ...any) error {
	// rate-limit outgoing edits
	time.Sleep(editDelay)
	_, err := ctx.Bot().Edit(message, text, args...)
	if err != nil {
		slog.Error("failed to edit message", "chatID", ctx.Sender().ID, "err", err)
	}

	return err
}

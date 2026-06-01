package service

import (
	"errors"
	"fmt"
	"log"
	"time"
	"translator/internal/model"
	"translator/internal/repository"

	"gopkg.in/telebot.v3"
)

const editDelay = 300 * time.Millisecond
const msgNotRegistered = "You are not registered yet. Please use /start command"

type TelegramService struct {
	clientRepository repository.ClientRepository
	openAiService    *OpenAiService
}

func NewTelegramService(clientRepository repository.ClientRepository) *TelegramService {
	return &TelegramService{
		clientRepository: clientRepository,
		openAiService:    NewOpenAiService(),
	}
}

func (t *TelegramService) Register(ctx telebot.Context) error {
	id := ctx.Sender().ID
	if !t.clientRepository.SaveIfAbsent(id) {
		return ctx.Send("You are already registered!")
	}

	return ctx.Send("You've successfully registered!")
}

func (t *TelegramService) Translate(ctx telebot.Context) error {
	lang, err := t.getLang(ctx)
	if errors.Is(err, repository.ErrClientNotFound) {
		return ctx.Send(msgNotRegistered)
	}
	if err != nil {
		return t.handleUnexpectedError(ctx, err)
	}

	text := ctx.Message().Text
	placeholder := translationPlaceholders[lang]
	if placeholder == "" {
		return fmt.Errorf("no translation placeholder found for '%s' language", lang)
	}

	sent, err := ctx.Bot().Send(ctx.Recipient(), placeholder)
	if err != nil {
		log.Printf("Failed to send translation placeholder: %v", err)
		return err
	}

	completion, err := t.openAiService.translationMessage(text, lang)
	if err != nil {
		log.Printf("Failed to send text: %v", err)
		return err
	}
	if len(completion.Choices) == 0 {
		return fmt.Errorf("openai returned no choices")
	}

	//rate-limit outgoing edits
	time.Sleep(editDelay)
	_, err = ctx.Bot().Edit(sent, completion.Choices[0].Message.Content)

	return err
}

func (t *TelegramService) SetLangCommandButtons(ctx telebot.Context) error {
	lang, err := t.getLang(ctx)
	if errors.Is(err, repository.ErrClientNotFound) {
		return ctx.Send(msgNotRegistered)
	}
	if err != nil {
		return t.handleUnexpectedError(ctx, err)
	}

	menu := &telebot.ReplyMarkup{}
	row := telebot.Row{}
	for _, langTitle := range getLangTitles(lang) {
		row = append(row, menu.Data(langTitle.Title, "set_lang", string(langTitle.Lang)))
	}

	menu.Inline(row)

	return ctx.Send("Choose preferred language", menu)
}

func (t *TelegramService) SetLangCommand(ctx telebot.Context) error {
	lang, err := t.setLang(ctx)
	if errors.Is(err, repository.ErrClientNotFound) {
		return ctx.Send(msgNotRegistered)
	}
	if errors.Is(err, model.ErrUnsupportedLanguage) {
		return ctx.Send(fmt.Sprintf("Unsupported language: '%s'", ctx.Data()))
	}
	if err != nil {
		return t.handleUnexpectedError(ctx, err)
	}

	return ctx.Edit(fmt.Sprintf("Language is set: %s", lang))
}

func (t *TelegramService) handleUnexpectedError(ctx telebot.Context, err error) error {
	log.Printf("unexpected error: %v", err)
	return ctx.Send("Unexpected error. Try again later.")
}

func (t *TelegramService) getLang(ctx telebot.Context) (model.Language, error) {
	return t.clientRepository.GetLang(ctx.Sender().ID)
}

func (t *TelegramService) setLang(ctx telebot.Context) (model.Language, error) {
	lang, err := model.ParseLang(ctx.Data())
	if err != nil {
		return "", err
	}

	return lang, t.clientRepository.SetLang(ctx.Sender().ID, lang)
}

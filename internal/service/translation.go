package service

import (
	"errors"
	"fmt"
	"strings"
	"translator/internal/model"
	"translator/internal/util"
)

type TranslationService struct {
	openAiService *OpenAiService
}

func NewTranslationService() (*TranslationService, error) {
	openAiService, err := NewOpenAiService()
	if err != nil {
		return nil, fmt.Errorf("failed to create translation service: %w", err)
	}

	return &TranslationService{openAiService: openAiService}, nil
}

var ErrNotGeorgianText = errors.New("not a georgian text")

func (t *TranslationService) Translate(text string, template model.MessageTemplate) (string, error) {
	completion, err := t.openAiService.translationMessage(text, template)
	if err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("openai returned no choices")
	}

	translation := completion.Choices[0].Message.Content
	return t.resolveTranslation(text, template, translation)
}

func (t *TranslationService) resolveTranslation(text string, template model.MessageTemplate, translation string) (string, error) {
	if translation == "NOT_GEORGIAN_TEXT" {
		return "", ErrNotGeorgianText
	}
	if translation == "GEORGIAN_TRANSLITERATION" {
		return t.Translate(util.ToGeorgian(text), template)
	}

	translation = strings.ReplaceAll(translation, "—", "-")
	return translation, nil
}

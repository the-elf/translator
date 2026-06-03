package service

import (
	"fmt"
	"translator/internal/model"
)

type TranslationService struct {
	openAiService *OpenAiService
}

func NewTranslationService() (*TranslationService, error) {
	openAiService, err := NewOpenAiService()
	if err != nil {
		return nil, fmt.Errorf("failed to create translation service: %v", err)
	}

	return &TranslationService{openAiService: openAiService}, nil
}

var ErrOpenAiEmptyResponse = fmt.Errorf("openai returned no choices")

func (t *TranslationService) Translate(text string, template model.MessageTemplate) (string, error) {
	completion, err := t.openAiService.translationMessage(text, template)
	if err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", ErrOpenAiEmptyResponse
	}

	return completion.Choices[0].Message.Content, nil
}

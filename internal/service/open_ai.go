package service

import (
	"context"
	"fmt"
	"os"
	"time"
	"translator/internal/model"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

type OpenAiService struct {
	client openai.Client
}

const (
	openAiApiKeyKey = "OPENAI_API_KEY"
	gptModelVersion = "gpt-5.4-nano"
)

func NewOpenAiService() *OpenAiService {
	_, ok := os.LookupEnv(openAiApiKeyKey)
	if !ok {
		panic(fmt.Sprintf("Missing %s", openAiApiKeyKey))
	}
	return &OpenAiService{client: openai.NewClient()}
}

var translationPrompts = map[model.Language]string{
	model.EnglishLanguage: "Translate the text into English without adding anything extra. Return only the translation of the provided text and nothing else.",
	model.RussianLanguage: "Переведи текст на русский язык, без отсебятины. В ответ верни только перевод предоставленного текста и ничего больше.",
}

func (o *OpenAiService) translationMessage(text string, lang model.Language) (*openai.ChatCompletion, error) {
	prompt := translationPrompts[lang]
	if prompt == "" {
		return nil, fmt.Errorf("no translation prompt found for '%s' language", lang)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	completion, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       gptModelVersion,
		Temperature: param.NewOpt(0.0),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(text),
		},
	})
	// todo обработать случай context deadline exceeded, чтобы не висело соообщение "Перевожу..."
	if err != nil {
		return nil, err
	}

	return completion, nil
}

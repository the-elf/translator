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

func NewOpenAiService() (*OpenAiService, error) {
	if _, ok := os.LookupEnv(openAiApiKeyKey); !ok {
		return nil, fmt.Errorf("missing %s", openAiApiKeyKey)
	}
	return &OpenAiService{client: openai.NewClient()}, nil
}

func (o *OpenAiService) translationMessage(text string, template model.MessageTemplate) (*openai.ChatCompletion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       gptModelVersion,
		Temperature: param.NewOpt(0.0),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(template.Text),
			openai.UserMessage(text),
		},
	})
}

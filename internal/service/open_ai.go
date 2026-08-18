package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"translator/internal/model"
	"translator/internal/util"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

type OpenAiService struct {
	client             openai.Client
	gptModelVersion    string
	translationTimeout time.Duration
}

const (
	openAiApiKeyKey       = "OPENAI_API_KEY"
	gptModelVersionKey    = "GPT_MODEL_VERSION"
	translationTimeoutKey = "TRANSLATION_TIMEOUT_SECONDS"
)

func NewOpenAiService() (*OpenAiService, error) {
	envs, err := util.RequireEnvs(openAiApiKeyKey, gptModelVersionKey, translationTimeoutKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create openai service: %w", err)
	}

	timeoutValue, err := strconv.ParseInt(envs[translationTimeoutKey], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse translation timeout: %w", err)
	}
	if timeoutValue <= 0 {
		return nil, fmt.Errorf("translation timeout must be greater than 0")
	}

	translationTimeout := time.Duration(timeoutValue) * time.Second
	return &OpenAiService{
		client:             openai.NewClient(),
		gptModelVersion:    envs[gptModelVersionKey],
		translationTimeout: translationTimeout,
	}, nil
}

func (o *OpenAiService) translationMessage(text string, template model.MessageTemplate) (*openai.ChatCompletion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.translationTimeout)
	defer cancel()

	return o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       o.gptModelVersion,
		Temperature: param.NewOpt(0.0),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(template.Text),
			openai.UserMessage(text),
		},
	})
}

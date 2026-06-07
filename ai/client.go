package ai

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

var ErrEmptyAPIKey = errors.New("openai api key is empty")

type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(apiKey string, model string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIClient{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
		model: model,
	}, nil
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: c.model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
	})
	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}

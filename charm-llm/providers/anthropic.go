package providers

import (
	"context"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	BaseProvider
	Client *anthropic.Client
}

func NewAnthropicProvider(model string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)

	return &AnthropicProvider{
		BaseProvider: BaseProvider{
			Name:  "anthropic",
			Model: model,
		},
		Client: &client,
	}
}

func (a *AnthropicProvider) GetResponse(ctx context.Context, prompt string) (string, error) {
	message, err := a.Client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_20250514,
		MaxTokens: 1000,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	if len(message.Content) == 0 {
		return "", nil
	}

	return message.Content[0].Text, nil
}

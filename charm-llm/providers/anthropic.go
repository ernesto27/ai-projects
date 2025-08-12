package providers

import (
	"context"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	BaseProvider
	Client *anthropic.Client
}

func parseModelName(shortName string) anthropic.Model {
	if shortName == "" {
		return anthropic.ModelClaudeSonnet4_0
	}

	shortName = strings.ToLower(strings.TrimSpace(shortName))

	switch shortName {
	case "claude-3-7":
		return anthropic.ModelClaude3_7SonnetLatest
	case "claude-3-5":
		return anthropic.ModelClaude3_5SonnetLatest
	case "claude-3-5-haiku":
		return anthropic.ModelClaude3_5HaikuLatest
	case "claude-4":
		return anthropic.ModelClaudeSonnet4_0
	case "claude-4-opus":
		return anthropic.ModelClaudeOpus4_0
	default:
		return anthropic.Model(shortName)
	}
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
	model := parseModelName(a.Model)

	message, err := a.Client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     model,
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

func (a *AnthropicProvider) GetResolvedModel() string {
	if a.Model == "" {
		return string(anthropic.ModelClaudeSonnet4_0)
	}
	return a.Model
}

package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type OpenAIProvider struct {
	BaseProvider
	Client openai.Client
}

func parseOpenAIModelName(shortName string) openai.ChatModel {
	if shortName == "" {
		return openai.ChatModelGPT4o
	}

	shortName = strings.ToLower(strings.TrimSpace(shortName))

	switch shortName {
	case "gpt-5":
		return openai.ChatModelGPT5
	case "gpt-4o":
		return openai.ChatModelGPT4o
	case "gpt-4o-2024-08-06":
		return openai.ChatModelGPT4o2024_08_06
	case "gpt-4o-mini":
		return openai.ChatModelGPT4oMini
	case "gpt-4":
		return openai.ChatModelGPT4
	case "gpt-4-turbo":
		return openai.ChatModelGPT4Turbo
	case "gpt-3":
		return openai.ChatModelGPT3_5Turbo
	default:
		return openai.ChatModel(shortName)
	}
}

func usesMaxCompletionTokens(model openai.ChatModel) bool {
	// Newer models like GPT-5 use MaxCompletionTokens instead of MaxTokens
	switch model {
	case openai.ChatModelGPT5:
		return true
	default:
		return false
	}
}

func createChatParams(model openai.ChatModel, prompt string) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: model,
	}

	if usesMaxCompletionTokens(model) {
		params.MaxCompletionTokens = openai.Int(MAX_TOKENS)
	} else {
		params.MaxTokens = openai.Int(MAX_TOKENS)
	}

	return params
}

func NewOpenAIProvider(model string, apiKey string) *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIProvider{
		BaseProvider: BaseProvider{
			Name:  "openai",
			Model: model,
		},
		Client: client,
	}
}

func (o *OpenAIProvider) GetResponse(ctx context.Context, prompt string) (string, error) {
	model := parseOpenAIModelName(o.Model)

	chatCompletion, err := o.Client.Chat.Completions.New(ctx, createChatParams(model, prompt))
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", nil
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (o *OpenAIProvider) GetStreamResponse(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	textChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(textChan)
		defer close(errChan)

		model := parseOpenAIModelName(o.Model)

		stream := o.Client.Chat.Completions.NewStreaming(ctx, createChatParams(model, prompt))

		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				select {
				case textChan <- chunk.Choices[0].Delta.Content:
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				}
			}
		}

		if stream.Err() != nil {
			errChan <- fmt.Errorf("streaming error: %w", stream.Err())
		}
	}()

	return textChan, errChan
}

func (o *OpenAIProvider) GetResolvedModel() string {
	if o.Model == "" {
		return string(openai.ChatModelGPT4o)
	}
	return o.Model
}

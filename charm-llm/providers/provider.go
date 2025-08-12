package providers

import "context"

type LLMProvider interface {
	GetResponse(ctx context.Context, prompt string) (string, error)
	GetStreamResponse(ctx context.Context, prompt string) (<-chan string, <-chan error)
	GetName() string
	SetModel(model string)
	GetModel() string
	GetResolvedModel() string
}

type BaseProvider struct {
	Name  string
	Model string
}

func (b *BaseProvider) GetName() string {
	return b.Name
}

func (b *BaseProvider) SetModel(model string) {
	b.Model = model
}

func (b *BaseProvider) GetModel() string {
	return b.Model
}
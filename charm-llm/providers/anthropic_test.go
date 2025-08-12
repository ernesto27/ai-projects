package providers

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestParseModelName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected anthropic.Model
	}{
		{
			name:     "empty string defaults to claude-4",
			input:    "",
			expected: anthropic.ModelClaudeSonnet4_0,
		},
		{
			name:     "claude-3-7 maps to latest",
			input:    "claude-3-7",
			expected: anthropic.ModelClaude3_7SonnetLatest,
		},
		{
			name:     "claude-3-5 maps to latest",
			input:    "claude-3-5",
			expected: anthropic.ModelClaude3_5SonnetLatest,
		},
		{
			name:     "claude-3-5-haiku maps to latest",
			input:    "claude-3-5-haiku",
			expected: anthropic.ModelClaude3_5HaikuLatest,
		},
		{
			name:     "claude-4 maps to sonnet 4.0",
			input:    "claude-4",
			expected: anthropic.ModelClaudeSonnet4_0,
		},
		{
			name:     "claude-4-opus maps to opus 4.0",
			input:    "claude-4-opus",
			expected: anthropic.ModelClaudeOpus4_0,
		},
		{
			name:     "case insensitive - CLAUDE-4",
			input:    "CLAUDE-4",
			expected: anthropic.ModelClaudeSonnet4_0,
		},
		{
			name:     "whitespace handling",
			input:    "  claude-3-7  ",
			expected: anthropic.ModelClaude3_7SonnetLatest,
		},
		{
			name:     "full model name passthrough",
			input:    "claude-sonnet-4-20250514",
			expected: anthropic.Model("claude-sonnet-4-20250514"),
		},
		{
			name:     "unknown model passthrough",
			input:    "custom-model",
			expected: anthropic.Model("custom-model"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseModelName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnthropicProvider_GetResolvedModel(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		expected string
	}{
		{
			name:     "empty model returns default",
			model:    "",
			expected: string(anthropic.ModelClaudeSonnet4_0),
		},
		{
			name:     "short model name returned as-is",
			model:    "claude-3-7",
			expected: "claude-3-7",
		},
		{
			name:     "full model name returned as-is",
			model:    "claude-sonnet-4-20250514",
			expected: "claude-sonnet-4-20250514",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &AnthropicProvider{
				BaseProvider: BaseProvider{
					Model: tt.model,
				},
			}
			result := provider.GetResolvedModel()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAnthropicProvider(t *testing.T) {
	model := "claude-4"
	provider := NewAnthropicProvider(model)

	assert.NotNil(t, provider)
	assert.Equal(t, "anthropic", provider.GetName())
	assert.Equal(t, model, provider.GetModel())
	assert.NotNil(t, provider.Client)
}
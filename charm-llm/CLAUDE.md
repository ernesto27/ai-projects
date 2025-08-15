# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool called `charm-llm` that provides a beautiful interface for interacting with various LLM providers (OpenAI, Anthropic). Built in Go using the Cobra CLI framework and Charm UI libraries for styling.

## Common Commands

### Build and Run
```bash
# Run directly
go run . -p anthropic -m claude-4 "Your prompt here"

# Build binary
go build -o charm-llm

# Run tests
go test ./...

# Run specific provider tests
go test ./providers/
```

### Usage Examples
```bash
# Basic usage with Anthropic
go run . -p anthropic -m claude-4 "Create a curl clone using python"

# With streaming
go run . -p anthropic -m claude-4 -s "Explain this code"

# OpenAI usage
go run . -p openai -m gpt-4o "Write a function to sort an array"

# Compare responses from all available providers
go run . -c "Create a curl clone using python"
go run . -c -m "anthropic:claude-4,openai:gpt-4o" "Compare with specific models"
go run . -c -m "anthropic:claude-3-5" "Specify model for some providers"

# Configuration management
go run . config set-anthropic-key YOUR_API_KEY
go run . config set-openai-key YOUR_API_KEY
go run . config show
```

## Architecture

### Core Components

- **main.go**: Entry point that delegates to cmd package
- **cmd/**: Cobra CLI commands and main application logic
  - `root.go`: Main command with provider/model selection, request handling, and parallel comparison
  - `config.go`: API key management subcommands
- **providers/**: LLM provider implementations
  - `provider.go`: Interface definition and base provider
  - `anthropic.go`: Anthropic Claude provider with model shortcuts
  - `openai.go`: OpenAI GPT provider with model shortcuts
- **config/**: Configuration management for API keys stored in `~/.config/charm-llm/config.json`
- **tui/**: Terminal UI styling using Charm libraries

### Provider System

The app uses a provider interface pattern where each LLM service implements:
- `GetResponse()`: Non-streaming response
- `GetStreamResponse()`: Streaming response with channels
- Model name resolution with shortcuts (e.g., "claude-4" → "claude-sonnet-4.0")

### Configuration

API keys are stored securely in `~/.config/charm-llm/config.json` with 0600 permissions. The app supports both file-based config and environment variables (ANTHROPIC_API_KEY, OPENAI_API_KEY).

### Compare Mode

The `-c` flag enables parallel comparison of responses from all configured providers:
- Executes requests to all available providers simultaneously using goroutines
- Displays responses in a beautifully formatted comparison view
- Shows individual response times and total execution time
- Handles partial failures gracefully (shows successful responses even if some providers fail)
- Supports provider-specific model selection with `provider:model` syntax
- Only non-streaming mode (streaming disabled for clean comparison output)
- Can save combined responses to markdown file with `-f` flag

### Provider-Specific Models
In compare mode, use the `provider:model` syntax to specify different models for each provider:
- `-m "anthropic:claude-4,openai:gpt-4o"` - specific models for both providers
- `-m "anthropic:claude-3-5"` - specific model for anthropic, default for openai
- No `-m` flag - uses default models for all providers

### UI/UX

Uses Charmbracelet libraries for rich terminal output:
- Glamour for markdown rendering of responses
- Lipgloss for styled terminal output
- Custom styling in tui/styles.go
- Color-coded provider sections in compare mode

## Model Shortcuts

### Anthropic Models
- `claude-3-7` → Claude 3.7 Sonnet Latest
- `claude-3-5` → Claude 3.5 Sonnet Latest  
- `claude-3-5-haiku` → Claude 3.5 Haiku Latest
- `claude-4` → Claude Sonnet 4.0 (default)
- `claude-4-opus` → Claude Opus 4.0

### OpenAI Models
- `gpt-5` → GPT-5
- `gpt-4o` → GPT-4o (default)
- `gpt-4o-mini` → GPT-4o Mini
- `gpt-4` → GPT-4
- `gpt-4-turbo` → GPT-4 Turbo
- `gpt-3` → GPT-3.5 Turbo
- do not test cli app using go build, i test myself
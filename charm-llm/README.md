# ✨ Charm LLM

A CLI tool for interacting with various LLM providers with an elegant terminal experience.

## 🚀 Features

- 🤖 **Multiple LLM Providers**: Support for Anthropic Claude and OpenAI GPT models
- 📡 **Streaming Responses**: Real-time streaming for faster interactions
- 🎨 **Rich UI**: Terminal interface with markdown rendering
- ⚡ **Model Shortcuts**: Easy-to-remember aliases for popular models
- 🛠️ **Simple Commands**: Intuitive CLI with helpful commands

## 📦 Installation

### Prerequisites
- Go 1.24.0 or later

### Install from Source

```bash
# Clone the repository
git clone <repository-url>
cd charm-llm

# Build and install
go build -o charm-llm
sudo mv charm-llm /usr/local/bin/

# Or install directly
go install
```

## ⚙️ Setup

### Configure API Keys

Before using charm-llm, you need to configure your API keys:

```bash
# Set Anthropic API key
charm-llm config set-anthropic-key YOUR_ANTHROPIC_API_KEY

# Set OpenAI API key  
charm-llm config set-openai-key YOUR_OPENAI_API_KEY

# View current configuration
charm-llm config show
```

## 🎯 Usage

### Basic Usage

```bash
# Chat with Claude
charm-llm -p anthropic -m claude-4 "Explain quantum computing"

# Chat with GPT
charm-llm -p openai -m gpt-4o "Write a Python function to sort a list"

# Enable streaming for real-time responses
charm-llm -p anthropic -m claude-4 -s "Create a REST API in Go"
```

### Model Shortcuts

#### Anthropic Models
- `claude-4` → Claude Sonnet 4.0 (default)
- `claude-4-opus` → Claude Opus 4.0
- `claude-3-7` → Claude 3.7 Sonnet Latest
- `claude-3-5` → Claude 3.5 Sonnet Latest
- `claude-3-5-haiku` → Claude 3.5 Haiku Latest

#### OpenAI Models
- `gpt-4o` → GPT-4o (default)
- `gpt-5` → GPT-5
- `gpt-4o-mini` → GPT-4o Mini
- `gpt-4` → GPT-4
- `gpt-4-turbo` → GPT-4 Turbo
- `gpt-3` → GPT-3.5 Turbo

### Examples

```bash
# Quick code explanation
charm-llm -p anthropic -m claude-4 "Explain this Go code: func main() { fmt.Println(\"Hello\") }"

# Get help with debugging
charm-llm -p openai -m gpt-4o -s "How to fix a segmentation fault in C?"

# Creative writing
charm-llm -p anthropic -m claude-4-opus "Write a short story about a robot learning to paint"
```

## 📋 Command Reference

```bash
# Main command
charm-llm [flags] "your prompt"

# Flags
-p, --provider string   LLM provider (openai, anthropic) [required]
-m, --model string     Model name (e.g., claude-4, gpt-4o)
-s, --stream          Enable streaming response

# Configuration commands
charm-llm config set-anthropic-key [key]   Set Anthropic API key
charm-llm config set-openai-key [key]      Set OpenAI API key  
charm-llm config show                       Show current configuration
```

## 🏗️ Development

```bash
# Run from source
go run . -p anthropic -m claude-4 "Your prompt"

# Run tests
go test ./...

# Build binary
go build -o charm-llm
```

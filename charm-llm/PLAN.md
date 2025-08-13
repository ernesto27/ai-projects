CLI TOOL in which you can connect to some LLM provider (openAI, anthropic, etc.) and run a prompt against.

Uses charm golang libs to UI 
https://github.com/charmbracelet 

example of use 

charmllm -p openai --m gpt-4 "Create a curl clone using python"

  ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY go run . -p anthropic -m claude-3-5-sonnet-20241022 "create a curl clone using python"

- Response in terminal stdo 
- Save response to file 

Pipe file 
echo main.go | charmllm -p openai -m gpt-4 "explain this code"

Run multiple providers in parallel
charmllm -p openai,anthropic -m gpt-4,claude-5 "Create a curl clone using python"


# Plan for CLI Tool: charmllm
## Features
- [x] Create providers interface to support multiple LLM providers
- [x] Implement command-line arguments using cobra lib 
- [x] command to get response in stdout
- [x] get model from parameters,  set default model if not provided
- [x] option stream respone
- [x] add open AI provider
- [] Save response to file
- [] Run multiple providers in parallel
- [] Gemini provider
- [] copy response to clipboard
- [] show loading
- [] Pipe file input to the command
- [] Handle multiple providers in parallel
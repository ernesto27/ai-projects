package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"charm-llm/config"
	"charm-llm/providers"
	"charm-llm/tui"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var (
	provider   string
	model      string
	stream     bool
	saveToFile bool
)

var rootCmd = &cobra.Command{
	Use:   "charm-llm [prompt]",
	Short: "A beautiful CLI tool for LLM interactions",
	Long:  "A CLI tool that provides a beautiful interface for interacting with various LLM providers.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		handleRequest(provider, model, prompt, stream, saveToFile)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider (openai, anthropic)")
	rootCmd.Flags().StringVarP(&model, "model", "m", "", "Model name (e.g., claude-4, gpt-4o, gpt-4o-mini)")
	rootCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Enable streaming response")
	rootCmd.Flags().BoolVarP(&saveToFile, "save-to-file", "f", false, "Save response to a random txt file")
	rootCmd.MarkFlagRequired("provider")
}

func generateRandomFilename() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("response_%s.txt", hex.EncodeToString(bytes))
}

func saveResponseToFile(content string) error {
	filename := generateRandomFilename()
	
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	
	filePath := filepath.Join(wd, filename)
	
	// Write content to file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	
	fmt.Printf("💾 Response saved to: %s\n", filename)
	return nil
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createProvider(providerName, model string) (providers.LLMProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %v", err)
	}

	switch strings.ToLower(providerName) {
	case "anthropic":
		apiKey := cfg.GetAnthropicKey()
		if apiKey == "" {
			return nil, fmt.Errorf("Anthropic API key not found. Set it with: charm-llm config set-anthropic-key YOUR_KEY")
		}
		return providers.NewAnthropicProvider(model, apiKey), nil
	case "openai":
		apiKey := cfg.GetOpenAIKey()
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API key not found. Set it with: charm-llm config set-openai-key YOUR_KEY")
		}
		return providers.NewOpenAIProvider(model, apiKey), nil
	default:
		return nil, fmt.Errorf("Provider '%s' is not supported. Available providers: anthropic, openai", providerName)
	}
}

func handleRequest(provider, model, prompt string, stream bool, saveToFile bool) {
	clearScreen()

	// Create provider first to get resolved model name
	llmProvider, err := createProvider(provider, model)
	if err != nil {
		errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error: %s", err.Error()))
		fmt.Println(errorMsg)
		return
	}

	// Initial display with resolved model name
	title := tui.HeaderStyle.Render("✨ Charm LLM")
	info := tui.InfoStyle.Render(fmt.Sprintf("Provider: %s • Model: %s", provider, llmProvider.GetResolvedModel()))
	styledPrompt := tui.PromptStyle.Render(fmt.Sprintf("💬 %s", prompt))

	// Show initial interface
	fmt.Println(title)
	fmt.Println()
	fmt.Println(info)
	fmt.Println(styledPrompt)
	fmt.Println()

	if stream {
		handleStreamingResponse(llmProvider, prompt, saveToFile)
	} else {
		// Handle regular response with spinner
		handleNonStreamingResponse(llmProvider, prompt, saveToFile)
	}
}

func handleStreamingResponse(llmProvider providers.LLMProvider, prompt string, saveToFile bool) {
	ctx := context.Background()

	fmt.Println("🤖 Response:")
	fmt.Println()

	textChan, errChan := llmProvider.GetStreamResponse(ctx, prompt)
	var fullResponse strings.Builder

	for {
		select {
		case chunk, ok := <-textChan:
			if !ok {
				// Streaming finished, render final markdown
				finalText := fullResponse.String()
				if finalText != "" {
					r, err := glamour.NewTermRenderer(
						glamour.WithAutoStyle(),
						glamour.WithWordWrap(tui.GetTerminalWidth()-6),
					)
					if err == nil {
						if out, err := r.Render(finalText); err == nil {
							// Clear the line and move cursor to beginning of response
							fmt.Print("\r\033[K")
							// Move cursor up to beginning of response
							lines := strings.Count(fullResponse.String(), "\n") + 1
							if lines > 1 {
								fmt.Printf("\033[%dA", lines)
							}
							fmt.Print(out)
						}
					}
				}
				
				// Save to file if requested
				if saveToFile {
					if err := saveResponseToFile(fullResponse.String()); err != nil {
						errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error saving to file: %s", err.Error()))
						fmt.Println(errorMsg)
					}
				}
				
				fmt.Println()
				return
			}
			fullResponse.WriteString(chunk)
			fmt.Print(chunk)
		case err := <-errChan:
			if err != nil {
				errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error getting streaming response: %s", err.Error()))
				fmt.Println(errorMsg)
				return
			}
		case <-ctx.Done():
			errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Request cancelled: %s", ctx.Err().Error()))
			fmt.Println(errorMsg)
			return
		}
	}
}

func handleNonStreamingResponse(llmProvider providers.LLMProvider, prompt string, saveToFile bool) {
	var response string

	err := tui.ShowSpinnerWhile("Thinking...", func(ctx context.Context) error {
		resp, err := llmProvider.GetResponse(ctx, prompt)
		response = resp
		return err
	})

	if err != nil {
		errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error getting response: %s", err.Error()))
		fmt.Println(errorMsg)
		return
	}

	// Render markdown response
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(tui.GetTerminalWidth()-6),
	)
	if err != nil {
		// Fallback to plain styling if glamour fails
		styledResponse := tui.ResponseStyle.Render(fmt.Sprintf("🤖 %s", response))
		fmt.Println(styledResponse)
		
		// Save to file if requested
		if saveToFile {
			if err := saveResponseToFile(response); err != nil {
				errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error saving to file: %s", err.Error()))
				fmt.Println(errorMsg)
			}
		}
		return
	}

	out, err := r.Render(response)
	if err != nil {
		// Fallback to plain styling if rendering fails
		styledResponse := tui.ResponseStyle.Render(fmt.Sprintf("🤖 %s", response))
		fmt.Println(styledResponse)
		
		// Save to file if requested
		if saveToFile {
			if err := saveResponseToFile(response); err != nil {
				errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error saving to file: %s", err.Error()))
				fmt.Println(errorMsg)
			}
		}
		return
	}

	// Display the beautifully rendered markdown
	fmt.Println("🤖 Response:")
	fmt.Println()
	fmt.Print(out)
	
	// Save to file if requested
	if saveToFile {
		if err := saveResponseToFile(response); err != nil {
			errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error saving to file: %s", err.Error()))
			fmt.Println(errorMsg)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

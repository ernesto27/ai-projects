package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"charm-llm/providers"
	"charm-llm/tui"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var (
	provider string
	model    string
	stream   bool
)

var rootCmd = &cobra.Command{
	Use:   "charm-llm [prompt]",
	Short: "A beautiful CLI tool for LLM interactions",
	Long:  "A CLI tool that provides a beautiful interface for interacting with various LLM providers.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		handleRequest(provider, model, prompt, stream)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider (e.g., openai, anthropic)")
	rootCmd.Flags().StringVarP(&model, "model", "m", "", "Model name (e.g., claude-3-7, claude-4)")
	rootCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Enable streaming response")
	rootCmd.MarkFlagRequired("provider")
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
	switch strings.ToLower(providerName) {
	case "anthropic":
		return providers.NewAnthropicProvider(model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

func handleRequest(provider, model, prompt string, stream bool) {
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

	// Show loading text
	fmt.Println("🤔 Thinking...")

	ctx := context.Background()
	
	if stream {
		// Handle streaming response
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
	} else {
		// Handle regular response
		response, err := llmProvider.GetResponse(ctx, prompt)

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
			return
		}

		out, err := r.Render(response)
		if err != nil {
			// Fallback to plain styling if rendering fails
			styledResponse := tui.ResponseStyle.Render(fmt.Sprintf("🤖 %s", response))
			fmt.Println(styledResponse)
			return
		}

		// Display the beautifully rendered markdown
		fmt.Println("🤖 Response:")
		fmt.Println()
		fmt.Print(out)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

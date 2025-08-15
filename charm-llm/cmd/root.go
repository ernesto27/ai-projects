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
	"sync"
	"time"

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
	compare    bool
	cfg        *config.Config
)

type ProviderConfig struct {
	name   string
	model  string
	apiKey string
}

type ProviderResponse struct {
	Provider string
	Model    string
	Response string
	Duration time.Duration
	Error    error
}

type ComparisonResult struct {
	Responses []ProviderResponse
	TotalTime time.Duration
}

var rootCmd = &cobra.Command{
	Use:   "charm-llm [prompt]",
	Short: "A beautiful CLI tool for LLM interactions",
	Long: `A CLI tool that provides a beautiful interface for interacting with various LLM providers.

Examples:
  # Single provider usage
  charm-llm -p anthropic -m claude-4 "Explain quantum computing"
  charm-llm -p openai -m gpt-4o "Write a sorting algorithm"

  # Compare responses from all available providers
  charm-llm -c "Create a curl clone using python"
  charm-llm -c -m "anthropic:claude-4,openai:gpt-4o" "Compare with specific models"
  charm-llm -c -m "anthropic:claude-3-5" "Specify model for some providers"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		if compare {
			handleCompareRequest(prompt, model, saveToFile)
		} else {
			handleRequest(provider, model, prompt, stream, saveToFile)
		}
	},
}

func validateModelFlag(compare bool, model string) error {
	if model == "" {
		return nil // Always valid (use defaults)
	}
	
	hasProviderSyntax := strings.Contains(model, ":")
	
	if compare && !hasProviderSyntax {
		return fmt.Errorf("compare mode requires provider:model format (e.g., 'anthropic:claude-4,openai:gpt-4o')")
	}
	
	if !compare && hasProviderSyntax {
		return fmt.Errorf("single provider mode uses simple model format (e.g., 'claude-4')")
	}
	
	return nil
}

func parseProviderModels(modelFlag string) (map[string]string, error) {
	if modelFlag == "" {
		return make(map[string]string), nil
	}
	
	modelMap := make(map[string]string)
	
	pairs := strings.Split(modelFlag, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format: '%s' (expected provider:model)", pair)
		}
		
		provider := strings.ToLower(strings.TrimSpace(parts[0]))
		model := strings.TrimSpace(parts[1])
		
		// Validate provider name
		if provider != "anthropic" && provider != "openai" {
			return nil, fmt.Errorf("unsupported provider: '%s' (available: anthropic, openai)", provider)
		}
		
		modelMap[provider] = model
	}
	
	return modelMap, nil
}

func getModelForProvider(modelMap map[string]string, provider string) string {
	if model, exists := modelMap[provider]; exists {
		return model
	}
	return "" // Use default model
}

func init() {
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider (openai, anthropic)")
	rootCmd.Flags().StringVarP(&model, "model", "m", "", "Model name: simple format for single provider (e.g., claude-4) or provider:model format for compare mode (e.g., anthropic:claude-4,openai:gpt-4o)")
	rootCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Enable streaming response")
	rootCmd.Flags().BoolVarP(&saveToFile, "save-to-file", "f", false, "Save response to a random txt file")
	rootCmd.Flags().BoolVarP(&compare, "compare", "c", false, "Compare responses from all available providers")

	// Make provider required only when not in compare mode
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Load global config
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %v", err)
		}

		if !compare && provider == "" {
			return fmt.Errorf("provider is required when not using compare mode. Use -p flag or -c for compare mode")
		}
		
		// Validate model flag syntax
		if err := validateModelFlag(compare, model); err != nil {
			return err
		}
		
		return nil
	}
}

func generateRandomFilename() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("response_%s.md", hex.EncodeToString(bytes))
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

func handleCompareRequest(prompt string, modelFlag string, saveToFile bool) {
	clearScreen()

	// Parse provider-specific models
	modelMap, err := parseProviderModels(modelFlag)
	if err != nil {
		errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("❌ Error: %s", err.Error()))
		fmt.Println(errorMsg)
		return
	}

	// Determine which providers to use with specified or default models
	providerConfigs := []ProviderConfig{
		{"anthropic", getModelForProvider(modelMap, "anthropic"), cfg.GetAnthropicKey()},
		{"openai", getModelForProvider(modelMap, "openai"), cfg.GetOpenAIKey()},
	}

	// Filter to only providers with API keys
	var availableProviders []ProviderConfig

	for _, pc := range providerConfigs {
		if pc.apiKey != "" {
			availableProviders = append(availableProviders, pc)
		}
	}

	if len(availableProviders) == 0 {
		errorMsg := tui.ErrorStyle.Render("❌ Error: No API keys configured. Set keys with: charm-llm config set-anthropic-key or set-openai-key")
		fmt.Println(errorMsg)
		return
	}

	// Display header
	title := tui.HeaderStyle.Render("✨ Charm LLM - Compare Mode")
	info := tui.InfoStyle.Render(fmt.Sprintf("Comparing %d providers", len(availableProviders)))
	styledPrompt := tui.PromptStyle.Render(fmt.Sprintf("💬 %s", prompt))

	fmt.Println(title)
	fmt.Println()
	fmt.Println(info)
	fmt.Println(styledPrompt)
	fmt.Println()

	// Execute requests in parallel
	result := executeParallelRequests(availableProviders, prompt)

	// Display results
	displayComparisonResults(result, saveToFile)
}

func executeParallelRequests(providerConfigs []ProviderConfig, prompt string) ComparisonResult {
	var wg sync.WaitGroup
	responses := make(chan ProviderResponse, len(providerConfigs))

	startTime := time.Now()

	for _, pc := range providerConfigs {
		wg.Add(1)
		go func(providerName, model, apiKey string) {
			defer wg.Done()

			start := time.Now()

			// Create provider with specified or default model
			llmProvider, err := createProvider(providerName, model)
			if err != nil {
				responses <- ProviderResponse{
					Provider: providerName,
					Model:    model,
					Error:    err,
					Duration: time.Since(start),
				}
				return
			}

			// Get response (non-streaming only in compare mode)
			ctx := context.Background()
			response, err := llmProvider.GetResponse(ctx, prompt)

			responses <- ProviderResponse{
				Provider: providerName,
				Model:    llmProvider.GetResolvedModel(),
				Response: response,
				Error:    err,
				Duration: time.Since(start),
			}
		}(pc.name, pc.model, pc.apiKey)
	}

	// Wait for all requests to complete and close channel
	wg.Wait()
	close(responses)

	// Collect all responses
	var allResponses []ProviderResponse
	for response := range responses {
		allResponses = append(allResponses, response)
	}

	return ComparisonResult{
		Responses: allResponses,
		TotalTime: time.Since(startTime),
	}
}

func displayComparisonResults(result ComparisonResult, saveToFile bool) {
	fmt.Printf("⏱️  Total execution time: %s\n\n", result.TotalTime.Round(time.Millisecond))

	var allResponses strings.Builder

	for i, response := range result.Responses {
		// Provider header
		providerTitle := fmt.Sprintf("🤖 %s (%s)", strings.Title(response.Provider), response.Model)
		if response.Error != nil {
			providerTitle += " - ❌ ERROR"
		}

		fmt.Println(tui.HeaderStyle.Render(providerTitle))
		fmt.Printf("⏱️  Response time: %s\n", response.Duration.Round(time.Millisecond))
		fmt.Println()

		if response.Error != nil {
			errorMsg := tui.ErrorStyle.Render(fmt.Sprintf("Error: %s", response.Error.Error()))
			fmt.Println(errorMsg)
		} else {
			// Render markdown response
			r, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(tui.GetTerminalWidth()-6),
			)
			if err != nil {
				// Fallback to plain text
				fmt.Println(response.Response)
			} else {
				out, err := r.Render(response.Response)
				if err != nil {
					fmt.Println(response.Response)
				} else {
					fmt.Print(out)
				}
			}

			// Add to combined responses for file saving with markdown formatting
			allResponses.WriteString(fmt.Sprintf("# %s (%s)\n\n", strings.ToUpper(response.Provider), response.Model))
			allResponses.WriteString(fmt.Sprintf("**Response Time:** %s\n\n", response.Duration.Round(time.Millisecond)))
			allResponses.WriteString(response.Response)
			allResponses.WriteString("\n\n---\n\n")
		}

		// Add separator between responses (except for the last one)
		if i < len(result.Responses)-1 {
			fmt.Println("\n" + strings.Repeat("─", 80) + "\n")
		}
	}

	// Save to file if requested
	if saveToFile && allResponses.Len() > 0 {
		if err := saveResponseToFile(allResponses.String()); err != nil {
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

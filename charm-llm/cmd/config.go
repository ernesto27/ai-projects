package cmd

import (
	"fmt"

	"charm-llm/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  "Set and view API keys and other configuration options",
}

var setAnthropicKeyCmd = &cobra.Command{
	Use:   "set-anthropic-key [key]",
	Short: "Set Anthropic API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		err = cfg.SetAnthropicKey(args[0])
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("✅ Anthropic API key saved successfully")
	},
}

var setOpenAIKeyCmd = &cobra.Command{
	Use:   "set-openai-key [key]",
	Short: "Set OpenAI API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		err = cfg.SetOpenAIKey(args[0])
		if err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("✅ OpenAI API key saved successfully")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Print(cfg.String())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setAnthropicKeyCmd)
	configCmd.AddCommand(setOpenAIKeyCmd)
	configCmd.AddCommand(showConfigCmd)
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AnthropicAPIKey string `json:"anthropic_api_key,omitempty"`
	OpenAIAPIKey    string `json:"openai_api_key,omitempty"`
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(homeDir, ".config", "charm-llm")
	return configDir, nil
}

func getConfigFile() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "config.json"), nil
}

func Load() (*Config, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return &Config{}, err
	}
	
	// If config file doesn't exist, return empty config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		return &Config{}, err
	}
	
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, err
	}
	
	return &config, nil
}

func (c *Config) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	
	// Create config directory if it doesn't exist
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(configFile, data, 0600)
}

func (c *Config) SetAnthropicKey(key string) error {
	c.AnthropicAPIKey = key
	return c.Save()
}

func (c *Config) SetOpenAIKey(key string) error {
	c.OpenAIAPIKey = key
	return c.Save()
}

func (c *Config) GetAnthropicKey() string {
	return c.AnthropicAPIKey
}

func (c *Config) GetOpenAIKey() string {
	return c.OpenAIAPIKey
}

func (c *Config) String() string {
	result := "Configuration:\n"
	
	if c.GetAnthropicKey() != "" {
		result += "  Anthropic API Key: ✓ (configured)\n"
	} else {
		result += "  Anthropic API Key: ✗ (not set)\n"
	}
	
	if c.GetOpenAIKey() != "" {
		result += "  OpenAI API Key: ✓ (configured)\n"
	} else {
		result += "  OpenAI API Key: ✗ (not set)\n"
	}
	
	configFile, _ := getConfigFile()
	result += fmt.Sprintf("  Config file: %s\n", configFile)
	
	return result
}
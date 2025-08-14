package audio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFile represents the structure of the audio configuration file
type ConfigFile struct {
	Audio AudioConfig `json:"audio"`
}

// SaveConfig saves the audio configuration to a file
func SaveConfig(config AudioConfig, filename string) error {
	configFile := ConfigFile{
		Audio: config,
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// LoadConfig loads the audio configuration from a file
func LoadConfig(filename string) (AudioConfig, error) {
	// Return default config if file doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return AudioConfig{}, fmt.Errorf("failed to read config file: %v", err)
	}

	// Unmarshal from JSON
	var configFile ConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return AudioConfig{}, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Validate configuration
	if err := ValidateConfig(configFile.Audio); err != nil {
		return AudioConfig{}, fmt.Errorf("invalid config: %v", err)
	}

	return configFile.Audio, nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	// Try to use system config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "gameboy-emulator", "audio.json")
	}

	// Fallback to current directory
	return "audio.json"
}

// AudioPresets contains predefined audio configurations
var AudioPresets = map[string]AudioConfig{
	"default": {
		SampleRate: 44100,
		BufferSize: 1024,
		Volume:     1.0,
		Enabled:    true,
	},
	"low_latency": {
		SampleRate: 44100,
		BufferSize: 256,
		Volume:     1.0,
		Enabled:    true,
	},
	"high_quality": {
		SampleRate: 48000,
		BufferSize: 2048,
		Volume:     1.0,
		Enabled:    true,
	},
	"retro": {
		SampleRate: 22050,
		BufferSize: 512,
		Volume:     0.8,
		Enabled:    true,
	},
}

// GetPreset returns a predefined audio configuration
func GetPreset(name string) (AudioConfig, bool) {
	config, exists := AudioPresets[name]
	return config, exists
}

// ListPresets returns the names of all available presets
func ListPresets() []string {
	presets := make([]string, 0, len(AudioPresets))
	for name := range AudioPresets {
		presets = append(presets, name)
	}
	return presets
}
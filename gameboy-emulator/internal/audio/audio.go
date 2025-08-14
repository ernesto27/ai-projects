// Package audio implements the Game Boy emulator audio output system
// for rendering APU audio samples to external audio libraries.
//
// The audio system handles:
// - Real-time audio output via SDL2 or other audio libraries
// - Sample rate conversion and buffering
// - Volume control and audio configuration
// - Threading and synchronization for audio callbacks
package audio

import (
	"sync"
)

// Audio constants
const (
	// Standard audio parameters
	DefaultSampleRate = 44100 // Standard 44.1kHz sample rate
	DefaultBufferSize = 1024  // Audio buffer size in samples
	MinSampleRate     = 8000  // Minimum supported sample rate
	MaxSampleRate     = 96000 // Maximum supported sample rate
	MinBufferSize     = 256   // Minimum buffer size
	MaxBufferSize     = 4096  // Maximum buffer size
	
	// Audio format
	SampleFormat = "int16" // 16-bit signed integer samples
	Channels     = 2       // Stereo output (left/right)
)

// AudioConfig holds audio system configuration
type AudioConfig struct {
	SampleRate int     // Sample rate in Hz (e.g., 44100)
	BufferSize int     // Buffer size in samples
	Volume     float32 // Master volume (0.0-1.0)
	Enabled    bool    // Enable/disable audio output
}

// AudioSample represents a stereo audio sample
type AudioSample struct {
	Left  int16 // Left channel sample (-32768 to 32767)
	Right int16 // Right channel sample (-32768 to 32767)
}

// AudioOutputInterface defines the contract for audio output implementations
// This allows the emulator to work with different audio libraries (SDL2, PortAudio, etc.)
type AudioOutputInterface interface {
	// Initialize sets up the audio system with the given configuration
	Initialize(config AudioConfig) error
	
	// Start begins audio playback
	Start() error
	
	// Stop pauses audio playback
	Stop() error
	
	// PushSamples sends audio samples to the output buffer
	// samples should be interleaved stereo (left, right, left, right, ...)
	PushSamples(samples []int16) error
	
	// SetVolume adjusts the master volume (0.0-1.0)
	SetVolume(volume float32) error
	
	// GetConfig returns the current audio configuration
	GetConfig() AudioConfig
	
	// IsPlaying returns true if audio is currently playing
	IsPlaying() bool
	
	// GetBufferLevel returns the current buffer fill level (0.0-1.0)
	GetBufferLevel() float32
	
	// Cleanup releases audio system resources
	Cleanup() error
}

// AudioOutput manages the emulator audio output system
type AudioOutput struct {
	config AudioConfig
	impl   AudioOutputInterface // Actual audio library implementation
	mutex  sync.RWMutex         // Thread safety for configuration access
}

// NewAudioOutput creates a new audio output manager with the specified implementation
func NewAudioOutput(impl AudioOutputInterface) *AudioOutput {
	return &AudioOutput{
		impl: impl,
		config: AudioConfig{
			SampleRate: DefaultSampleRate,
			BufferSize: DefaultBufferSize,
			Volume:     1.0,
			Enabled:    true,
		},
	}
}

// Initialize sets up the audio output with the given configuration
func (a *AudioOutput) Initialize(config AudioConfig) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return err
	}
	
	a.config = config
	return a.impl.Initialize(config)
}

// Start begins audio playback
func (a *AudioOutput) Start() error {
	return a.impl.Start()
}

// Stop pauses audio playback
func (a *AudioOutput) Stop() error {
	return a.impl.Stop()
}

// PushSamples sends audio samples to the output
func (a *AudioOutput) PushSamples(samples []int16) error {
	a.mutex.RLock()
	enabled := a.config.Enabled
	a.mutex.RUnlock()
	
	if !enabled {
		return nil // Audio disabled, discard samples
	}
	
	return a.impl.PushSamples(samples)
}

// SetVolume adjusts the master volume
func (a *AudioOutput) SetVolume(volume float32) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	a.config.Volume = volume
	return a.impl.SetVolume(volume)
}

// GetConfig returns the current audio configuration (thread-safe)
func (a *AudioOutput) GetConfig() AudioConfig {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.config
}

// IsPlaying returns true if audio is currently playing
func (a *AudioOutput) IsPlaying() bool {
	return a.impl.IsPlaying()
}

// GetBufferLevel returns the current buffer fill level
func (a *AudioOutput) GetBufferLevel() float32 {
	return a.impl.GetBufferLevel()
}

// Enable enables or disables audio output
func (a *AudioOutput) Enable(enabled bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.config.Enabled = enabled
}

// IsEnabled returns true if audio output is enabled
func (a *AudioOutput) IsEnabled() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.config.Enabled
}

// Cleanup releases audio system resources
func (a *AudioOutput) Cleanup() error {
	return a.impl.Cleanup()
}

// DefaultConfig returns a default audio configuration
func DefaultConfig() AudioConfig {
	return AudioConfig{
		SampleRate: DefaultSampleRate,
		BufferSize: DefaultBufferSize,
		Volume:     1.0,
		Enabled:    true,
	}
}

// ValidateConfig checks if the audio configuration is valid
func ValidateConfig(config AudioConfig) error {
	if config.SampleRate < MinSampleRate || config.SampleRate > MaxSampleRate {
		return ErrInvalidSampleRate
	}
	
	if config.BufferSize < MinBufferSize || config.BufferSize > MaxBufferSize {
		return ErrInvalidBufferSize
	}
	
	if config.Volume < 0.0 || config.Volume > 1.0 {
		return ErrInvalidVolume
	}
	
	return nil
}

// ConvertSamplesToStereo converts mono samples to stereo by duplicating each sample
func ConvertSamplesToStereo(monoSamples []int16) []int16 {
	stereoSamples := make([]int16, len(monoSamples)*2)
	for i, sample := range monoSamples {
		stereoSamples[i*2] = sample     // Left channel
		stereoSamples[i*2+1] = sample   // Right channel
	}
	return stereoSamples
}

// MixStereoSamples mixes left and right channel samples into interleaved stereo format
func MixStereoSamples(left, right []int16) []int16 {
	if len(left) != len(right) {
		// Handle mismatched lengths by using the shorter length
		minLen := len(left)
		if len(right) < minLen {
			minLen = len(right)
		}
		left = left[:minLen]
		right = right[:minLen]
	}
	
	stereoSamples := make([]int16, len(left)*2)
	for i := 0; i < len(left); i++ {
		stereoSamples[i*2] = left[i]    // Left channel
		stereoSamples[i*2+1] = right[i] // Right channel
	}
	return stereoSamples
}

// ApplyVolume applies volume scaling to audio samples
func ApplyVolume(samples []int16, volume float32) {
	if volume == 1.0 {
		return // No scaling needed
	}
	
	for i := range samples {
		// Scale sample and clamp to prevent overflow
		scaled := float32(samples[i]) * volume
		if scaled > 32767 {
			samples[i] = 32767
		} else if scaled < -32768 {
			samples[i] = -32768
		} else {
			samples[i] = int16(scaled)
		}
	}
}
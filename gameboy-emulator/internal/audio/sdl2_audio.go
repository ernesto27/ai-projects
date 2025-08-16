package audio

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// SDL2AudioOutput implements AudioOutputInterface using SDL2 audio
type SDL2AudioOutput struct {
	deviceID    sdl.AudioDeviceID
	spec        *sdl.AudioSpec
	config      AudioConfig
	playing     bool
	initialized bool
}

// NewSDL2AudioOutput creates a new SDL2 audio output implementation
func NewSDL2AudioOutput() *SDL2AudioOutput {
	return &SDL2AudioOutput{}
}

// Initialize sets up SDL2 audio with the given configuration
func (s *SDL2AudioOutput) Initialize(config AudioConfig) error {
	if s.initialized {
		return ErrAudioAlreadyStarted
	}

	// Initialize SDL2 audio subsystem
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		return fmt.Errorf("%w: %v", ErrAudioInitFailed, err)
	}

	// Create audio specification
	spec := &sdl.AudioSpec{
		Freq:     int32(config.SampleRate),
		Format:   sdl.AUDIO_S16LSB, // 16-bit signed little-endian
		Channels: Channels,         // Stereo
		Samples:  uint16(config.BufferSize),
		Callback: nil, // Use queue-based audio instead of callback
	}

	// Open audio device
	deviceID, err := sdl.OpenAudioDevice("", false, spec, nil, 0)
	if err != nil {
		sdl.Quit()
		return fmt.Errorf("%w: %v", ErrAudioDeviceNotFound, err)
	}

	s.deviceID = deviceID
	s.spec = spec
	s.config = config
	s.initialized = true

	return nil
}

// Start begins audio playback
func (s *SDL2AudioOutput) Start() error {
	if !s.initialized {
		return ErrAudioNotInitialized
	}

	if s.playing {
		return ErrAudioAlreadyStarted
	}

	// Start audio playback (unpause the device)
	sdl.PauseAudioDevice(s.deviceID, false)
	s.playing = true

	return nil
}

// Stop pauses audio playback
func (s *SDL2AudioOutput) Stop() error {
	if !s.initialized {
		return ErrAudioNotInitialized
	}

	if !s.playing {
		return ErrAudioNotStarted
	}

	// Pause audio playback
	sdl.PauseAudioDevice(s.deviceID, true)
	s.playing = false

	return nil
}

// PushSamples sends audio samples to the output buffer
func (s *SDL2AudioOutput) PushSamples(samples []int16) error {
	if !s.initialized {
		return ErrAudioNotInitialized
	}

	// Apply volume scaling
	volumeScaledSamples := make([]int16, len(samples))
	copy(volumeScaledSamples, samples)
	ApplyVolume(volumeScaledSamples, s.config.Volume)

	// Convert to bytes for SDL2
	sampleBytes := (*[1 << 30]byte)(unsafe.Pointer(&volumeScaledSamples[0]))[:len(volumeScaledSamples)*2:len(volumeScaledSamples)*2]
	
	// Queue audio data to SDL2
	if err := sdl.QueueAudio(s.deviceID, sampleBytes); err != nil {
		return fmt.Errorf("failed to queue audio: %v", err)
	}

	return nil
}

// SetVolume adjusts the master volume
func (s *SDL2AudioOutput) SetVolume(volume float32) error {
	if !s.initialized {
		return ErrAudioNotInitialized
	}

	if volume < 0.0 || volume > 1.0 {
		return ErrInvalidVolume
	}

	s.config.Volume = volume
	return nil
}

// GetConfig returns the current audio configuration
func (s *SDL2AudioOutput) GetConfig() AudioConfig {
	return s.config
}

// IsPlaying returns true if audio is currently playing
func (s *SDL2AudioOutput) IsPlaying() bool {
	return s.playing
}

// GetBufferLevel returns the current buffer fill level (0.0-1.0)
func (s *SDL2AudioOutput) GetBufferLevel() float32 {
	if !s.initialized {
		return 0.0
	}
	
	// Get queued audio size from SDL2
	queuedSize := sdl.GetQueuedAudioSize(s.deviceID)
	maxBufferSize := uint32(s.config.BufferSize * Channels * 2) // 2 bytes per sample
	
	if maxBufferSize == 0 {
		return 0.0
	}
	
	level := float32(queuedSize) / float32(maxBufferSize)
	if level > 1.0 {
		level = 1.0
	}
	
	return level
}

// Cleanup releases SDL2 audio resources
func (s *SDL2AudioOutput) Cleanup() error {
	if !s.initialized {
		return nil
	}

	// Stop playback if running
	if s.playing {
		s.Stop()
	}

	// Close audio device
	sdl.CloseAudioDevice(s.deviceID)

	// Quit SDL2 audio subsystem
	sdl.Quit()

	s.initialized = false
	return nil
}


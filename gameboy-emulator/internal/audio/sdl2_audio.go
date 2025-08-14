package audio

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// SDL2AudioOutput implements AudioOutputInterface using SDL2 audio
type SDL2AudioOutput struct {
	deviceID   sdl.AudioDeviceID
	spec       *sdl.AudioSpec
	config     AudioConfig
	playing    bool
	buffer     []int16
	bufferMux  sync.Mutex
	sampleChan chan []int16
	initialized bool
}

// NewSDL2AudioOutput creates a new SDL2 audio output implementation
func NewSDL2AudioOutput() *SDL2AudioOutput {
	return &SDL2AudioOutput{
		sampleChan: make(chan []int16, 10), // Buffered channel for samples
	}
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
		Callback: sdl.AudioCallback(s.audioCallback),
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
	s.buffer = make([]int16, config.BufferSize*Channels)
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

	// Send samples to the audio callback via channel
	select {
	case s.sampleChan <- volumeScaledSamples:
		return nil
	default:
		// Channel is full, drop samples to prevent blocking
		return ErrBufferOverflow
	}
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
	// Return the channel buffer level as a percentage
	return float32(len(s.sampleChan)) / 10.0 // Channel capacity is 10
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

	// Close channel
	close(s.sampleChan)

	s.initialized = false
	return nil
}

// audioCallback is called by SDL2 when it needs more audio data
// This function runs in a separate thread created by SDL2
func (s *SDL2AudioOutput) audioCallback(userdata unsafe.Pointer, stream *uint8, length int32) {
	// Calculate number of samples needed
	samplesNeeded := int(length) / 2 // 2 bytes per int16 sample

	// Lock buffer access
	s.bufferMux.Lock()
	defer s.bufferMux.Unlock()

	// Clear the buffer first
	for i := 0; i < samplesNeeded && i < len(s.buffer); i++ {
		s.buffer[i] = 0
	}

	// Try to get samples from the channel
	select {
	case samples := <-s.sampleChan:
		// Copy samples to buffer, handling different lengths
		copyLen := samplesNeeded
		if len(samples) < copyLen {
			copyLen = len(samples)
		}
		if copyLen > len(s.buffer) {
			copyLen = len(s.buffer)
		}

		copy(s.buffer[:copyLen], samples[:copyLen])

	default:
		// No samples available, buffer remains silent (zeros)
	}

	// Copy buffer to SDL2 audio stream
	// Convert from []int16 to []uint8 for SDL2
	bufferBytes := (*[2]uint8)(unsafe.Pointer(&s.buffer[0]))
	streamSlice := (*[1 << 30]uint8)(unsafe.Pointer(stream))[:length:length]

	copyBytes := int(length)
	if len(s.buffer)*2 < copyBytes {
		copyBytes = len(s.buffer) * 2
	}

	for i := 0; i < copyBytes && i < len(s.buffer)*2; i++ {
		if i < len(streamSlice) {
			streamSlice[i] = bufferBytes[i]
		}
	}
}
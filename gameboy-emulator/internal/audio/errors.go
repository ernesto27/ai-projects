package audio

import "errors"

// Audio system error definitions
var (
	// Configuration errors
	ErrInvalidSampleRate = errors.New("invalid sample rate: must be between 8000 and 96000 Hz")
	ErrInvalidBufferSize = errors.New("invalid buffer size: must be between 256 and 4096 samples")
	ErrInvalidVolume     = errors.New("invalid volume: must be between 0.0 and 1.0")
	
	// Runtime errors
	ErrAudioNotInitialized = errors.New("audio system not initialized")
	ErrAudioAlreadyStarted = errors.New("audio already started")
	ErrAudioNotStarted     = errors.New("audio not started")
	ErrBufferOverflow      = errors.New("audio buffer overflow")
	ErrBufferUnderflow     = errors.New("audio buffer underflow")
	
	// System errors
	ErrAudioDeviceNotFound = errors.New("audio device not found")
	ErrAudioInitFailed     = errors.New("failed to initialize audio system")
	ErrAudioStartFailed    = errors.New("failed to start audio playback")
	ErrAudioStopFailed     = errors.New("failed to stop audio playback")
)
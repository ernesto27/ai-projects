# Game Boy Emulator Audio System Guide

## Overview

The Game Boy emulator features a complete audio output system that delivers authentic Game Boy sound through modern audio APIs. The system consists of two main components:

1. **APU (Audio Processing Unit)** - Generates authentic Game Boy audio samples
2. **Audio Output System** - Delivers samples to speakers via SDL2

## Architecture

```
Game Boy ROM → APU → Audio Output → SDL2 → Speakers
     ↓           ↓         ↓         ↓
   Sound      4-Channel  Sample   Real-time
  Commands    Mixing    Conversion  Playback
```

### Component Breakdown

#### APU (Audio Processing Unit)
- **Location**: `internal/apu/`
- **Purpose**: Authentic Game Boy sound generation
- **Output**: Float32 stereo samples at emulation sample rate

#### Audio Output System  
- **Location**: `internal/audio/`
- **Purpose**: Modern audio delivery and configuration
- **Backend**: SDL2 with queue-based audio delivery

## Audio Pipeline

### 1. APU Sample Generation
```go
// APU generates samples every CPU cycle
apu.Update(cycles)
samples := apu.GetSamples() // []float32 stereo samples
```

### 2. Sample Conversion
```go
// Convert float32 to int16 for SDL2
int16Samples := make([]int16, len(samples))
for i, sample := range samples {
    // Clamp and convert
    clamped := clamp(sample, -1.0, 1.0)
    int16Samples[i] = int16(clamped * 32767)
}
```

### 3. Audio Output
```go
// Queue samples to SDL2 (non-blocking)
err := audioOutput.PushSamples(int16Samples)
```

## Configuration

### Audio Configuration Options

```go
type AudioConfig struct {
    SampleRate int     // Audio sample rate (8000-96000 Hz)
    BufferSize int     // Buffer size in samples (256-4096)
    Volume     float32 // Master volume (0.0-1.0)
    Enabled    bool    // Enable/disable audio output
}
```

### Default Configuration
```go
config := audio.DefaultConfig()
// SampleRate: 44100 Hz
// BufferSize: 1024 samples (~23ms latency)
// Volume: 1.0 (100%)
// Enabled: true
```

### Configuration Presets
```go
// Available presets
presets := audio.ListPresets()
// "default", "low_latency", "high_quality", "retro"

config, exists := audio.GetPreset("low_latency")
// SampleRate: 44100 Hz, BufferSize: 256 samples (~6ms latency)
```

## Usage Examples

### Basic Setup
```go
// Create audio system
audioImpl := audio.NewSDL2AudioOutput()
audioOutput := audio.NewAudioOutput(audioImpl)

// Initialize with configuration
config := audio.DefaultConfig()
err := audioOutput.Initialize(config)
if err != nil {
    log.Fatal("Audio initialization failed:", err)
}

// Start playback
err = audioOutput.Start()
if err != nil {
    log.Fatal("Audio start failed:", err)
}
```

### Runtime Configuration
```go
// Adjust volume
audioOutput.SetVolume(0.7) // 70% volume

// Get current configuration
config := audioOutput.GetConfig()
fmt.Printf("Sample rate: %d Hz\n", config.SampleRate)

// Check playback status
if audioOutput.IsPlaying() {
    fmt.Println("Audio is playing")
}
```

### Buffer Management
```go
// Monitor buffer level
level := audioOutput.GetBufferLevel() // 0.0-1.0
if level > 0.8 {
    fmt.Println("Audio buffer getting full")
}

// Handle buffer overflow gracefully
err := audioOutput.PushSamples(samples)
if err == audio.ErrBufferOverflow {
    // Samples dropped, but emulation continues
    log.Println("Audio buffer overflow - dropping samples")
}
```

## Performance Optimization

### Latency Control
- **Low Latency**: Buffer size 256-512 samples (6-12ms)
- **Balanced**: Buffer size 1024 samples (23ms) 
- **High Compatibility**: Buffer size 2048+ samples (46ms+)

### CPU Usage
- Audio processing runs in separate thread
- Non-blocking sample queuing prevents emulation slowdown
- Automatic sample dropping on buffer overflow

### Memory Management
```go
// Efficient sample handling
samples := apu.GetSamples()
if len(samples) > 0 {
    // Process samples immediately
    audioOutput.PushSamples(convertSamples(samples))
    // APU automatically clears buffer
}
```

## Error Handling

### Common Errors
```go
// Audio device not found
if err == audio.ErrAudioDeviceNotFound {
    log.Println("No audio device available")
}

// Buffer overflow (non-critical)
if err == audio.ErrBufferOverflow {
    // Continue emulation, samples dropped
}

// Invalid configuration
if err == audio.ErrInvalidSampleRate {
    log.Println("Unsupported sample rate")
}
```

### Recovery Strategies
```go
// Graceful audio disable on failure
if err := audioOutput.Initialize(config); err != nil {
    log.Printf("Audio init failed: %v", err)
    config.Enabled = false // Disable audio, continue emulation
}
```

## Testing

### Mock Implementation
```go
// For testing without audio hardware
type MockAudioImpl struct{}

func (m *MockAudioImpl) Initialize(config AudioConfig) error { return nil }
func (m *MockAudioImpl) PushSamples(samples []int16) error { return nil }
// ... other interface methods

// Use in tests
mockAudio := &MockAudioImpl{}
audioOutput := audio.NewAudioOutput(mockAudio)
```

### Audio Quality Testing
```go
// Test with known waveforms
testSamples := generateTestTone(440.0, 44100, 1.0) // 440Hz tone
err := audioOutput.PushSamples(testSamples)

// Verify sample conversion accuracy
original := []float32{0.5, -0.5, 1.0, -1.0}
converted := convertToInt16(original)
expected := []int16{16383, -16384, 32767, -32768}
```

## Troubleshooting

### No Audio Output
1. Check SDL2 installation: `pkg-config --exists sdl2`
2. Verify audio device availability
3. Check volume levels (system and emulator)
4. Ensure audio is enabled in configuration

### Audio Crackling/Distortion
1. Increase buffer size to reduce underruns
2. Check system audio settings
3. Monitor CPU usage - may need optimization
4. Verify sample rate compatibility with audio device

### High Latency
1. Decrease buffer size (minimum 256 samples)
2. Use "low_latency" preset
3. Check system audio buffer settings
4. Consider different audio drivers

### Memory Issues
1. Monitor buffer overflow errors
2. Check for memory leaks in long-running sessions
3. Verify proper cleanup on emulator shutdown

## Advanced Features

### Custom Audio Backend
```go
// Implement AudioOutputInterface for custom backends
type CustomAudioOutput struct {
    // Custom implementation
}

func (c *CustomAudioOutput) Initialize(config AudioConfig) error {
    // Custom initialization
}
// ... implement other interface methods
```

### Audio Recording
```go
// Extend for recording capabilities
type RecordingAudioOutput struct {
    underlying AudioOutputInterface
    recorder   *AudioRecorder
}

func (r *RecordingAudioOutput) PushSamples(samples []int16) error {
    r.recorder.Record(samples) // Save to file
    return r.underlying.PushSamples(samples) // Continue playback
}
```

### Real-time Audio Analysis
```go
// Monitor audio levels
func (a *AudioOutput) PushSamples(samples []int16) error {
    // Calculate RMS level
    level := calculateRMSLevel(samples)
    if level > threshold {
        triggerVolumeWarning()
    }
    return a.impl.PushSamples(samples)
}
```

## API Reference

### AudioOutputInterface
```go
type AudioOutputInterface interface {
    Initialize(config AudioConfig) error
    Start() error
    Stop() error
    PushSamples(samples []int16) error
    SetVolume(volume float32) error
    GetConfig() AudioConfig
    IsPlaying() bool
    GetBufferLevel() float32
    Cleanup() error
}
```

### Key Functions
- `audio.NewSDL2AudioOutput()` - Create SDL2 backend
- `audio.NewAudioOutput(impl)` - Create audio output wrapper
- `audio.DefaultConfig()` - Get default configuration
- `audio.ValidateConfig(config)` - Validate configuration
- `audio.GetPreset(name)` - Get preset configuration

### Utility Functions
- `audio.ConvertSamplesToStereo(mono)` - Convert mono to stereo
- `audio.MixStereoSamples(left, right)` - Mix separate channels
- `audio.ApplyVolume(samples, volume)` - Apply volume scaling

## Integration with Emulator

The audio system integrates seamlessly with the main emulator loop:

```go
// In emulator.Step()
apu.Update(cycles)
if samples := apu.GetSamples(); samples != nil {
    audioOutput.PushSamples(convertSamples(samples))
}
```

This ensures audio stays synchronized with emulation timing while maintaining 60 FPS performance.
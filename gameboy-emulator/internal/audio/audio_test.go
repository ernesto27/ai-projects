package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAudioOutput implements AudioOutputInterface for testing
type MockAudioOutput struct {
	initialized bool
	started     bool
	samples     [][]int16
	config      AudioConfig
	volume      float32
	bufferLevel float32
	initError   error
	startError  error
	stopError   error
	pushError   error
}

func NewMockAudioOutput() *MockAudioOutput {
	return &MockAudioOutput{
		samples: make([][]int16, 0),
		volume:  1.0,
	}
}

func (m *MockAudioOutput) Initialize(config AudioConfig) error {
	if m.initError != nil {
		return m.initError
	}
	m.initialized = true
	m.config = config
	return nil
}

func (m *MockAudioOutput) Start() error {
	if m.startError != nil {
		return m.startError
	}
	if !m.initialized {
		return ErrAudioNotInitialized
	}
	m.started = true
	return nil
}

func (m *MockAudioOutput) Stop() error {
	if m.stopError != nil {
		return m.stopError
	}
	m.started = false
	return nil
}

func (m *MockAudioOutput) PushSamples(samples []int16) error {
	if m.pushError != nil {
		return m.pushError
	}
	if !m.initialized {
		return ErrAudioNotInitialized
	}
	
	// Copy samples to avoid slice modifications
	samplesCopy := make([]int16, len(samples))
	copy(samplesCopy, samples)
	m.samples = append(m.samples, samplesCopy)
	return nil
}

func (m *MockAudioOutput) SetVolume(volume float32) error {
	if volume < 0.0 || volume > 1.0 {
		return ErrInvalidVolume
	}
	m.volume = volume
	return nil
}

func (m *MockAudioOutput) GetConfig() AudioConfig {
	return m.config
}

func (m *MockAudioOutput) IsPlaying() bool {
	return m.started
}

func (m *MockAudioOutput) GetBufferLevel() float32 {
	return m.bufferLevel
}

func (m *MockAudioOutput) Cleanup() error {
	m.initialized = false
	m.started = false
	return nil
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, 44100, config.SampleRate)
	assert.Equal(t, 1024, config.BufferSize)
	assert.Equal(t, float32(1.0), config.Volume)
	assert.True(t, config.Enabled)
}

func TestAudioOutputInitialization(t *testing.T) {
	config := DefaultConfig()
	mockImpl := NewMockAudioOutput()
	
	audioOutput := NewAudioOutput(mockImpl)

	err := audioOutput.Initialize(config)
	require.NoError(t, err)
	assert.True(t, mockImpl.initialized)
	assert.Equal(t, config.SampleRate, mockImpl.config.SampleRate)
	assert.Equal(t, config.BufferSize, mockImpl.config.BufferSize)

	// Test start
	err = audioOutput.Start()
	assert.NoError(t, err)
	assert.True(t, mockImpl.started)

	// Cleanup
	err = audioOutput.Stop()
	assert.NoError(t, err)
	err = audioOutput.Cleanup()
	assert.NoError(t, err)
}

func TestAudioOutputDisabled(t *testing.T) {
	config := DefaultConfig()
	config.Enabled = false
	
	mockImpl := NewMockAudioOutput()
	audioOutput := NewAudioOutput(mockImpl)

	err := audioOutput.Initialize(config)
	assert.NoError(t, err)

	// Audio should initialize but not push samples when disabled
	testSamples := []int16{100, 200, 300, 400}
	err = audioOutput.PushSamples(testSamples)
	assert.NoError(t, err) // Should not error, but samples are discarded

	assert.Len(t, mockImpl.samples, 0) // No samples should be pushed when disabled
}

func TestAudioOutputPushSamples(t *testing.T) {
	config := DefaultConfig()
	mockImpl := NewMockAudioOutput()
	
	audioOutput := NewAudioOutput(mockImpl)

	err := audioOutput.Initialize(config)
	require.NoError(t, err)

	// Test pushing samples
	testSamples := []int16{1000, 2000, 3000, 4000}
	err = audioOutput.PushSamples(testSamples)
	assert.NoError(t, err)

	assert.Len(t, mockImpl.samples, 1)
	assert.Equal(t, testSamples, mockImpl.samples[0])

	// Cleanup
	err = audioOutput.Cleanup()
	assert.NoError(t, err)
}

func TestAudioOutputVolumeControl(t *testing.T) {
	config := DefaultConfig()
	mockImpl := NewMockAudioOutput()
	
	audioOutput := NewAudioOutput(mockImpl)

	err := audioOutput.Initialize(config)
	require.NoError(t, err)

	// Test volume setting
	err = audioOutput.SetVolume(0.8)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.8), mockImpl.volume)

	// Test volume clamping
	err = audioOutput.SetVolume(-0.5)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.0), mockImpl.volume)

	err = audioOutput.SetVolume(1.5)
	assert.NoError(t, err)
	assert.Equal(t, float32(1.0), mockImpl.volume)
}

func TestAudioOutputConfig(t *testing.T) {
	config := DefaultConfig()
	config.SampleRate = 48000
	config.BufferSize = 2048
	config.Volume = 0.7
	
	mockImpl := NewMockAudioOutput()
	audioOutput := NewAudioOutput(mockImpl)

	err := audioOutput.Initialize(config)
	require.NoError(t, err)

	retrievedConfig := audioOutput.GetConfig()
	assert.Equal(t, config.SampleRate, retrievedConfig.SampleRate)
	assert.Equal(t, config.BufferSize, retrievedConfig.BufferSize)
	assert.Equal(t, config.Volume, retrievedConfig.Volume)
	assert.Equal(t, config.Enabled, retrievedConfig.Enabled)
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := DefaultConfig()
	err := ValidateConfig(validConfig)
	assert.NoError(t, err)

	// Test invalid sample rate
	invalidConfig := validConfig
	invalidConfig.SampleRate = 100000
	err = ValidateConfig(invalidConfig)
	assert.Equal(t, ErrInvalidSampleRate, err)

	// Test invalid buffer size
	invalidConfig = validConfig
	invalidConfig.BufferSize = 100
	err = ValidateConfig(invalidConfig)
	assert.Equal(t, ErrInvalidBufferSize, err)

	// Test invalid volume
	invalidConfig = validConfig
	invalidConfig.Volume = 2.0
	err = ValidateConfig(invalidConfig)
	assert.Equal(t, ErrInvalidVolume, err)
}

func TestAudioUtilityFunctions(t *testing.T) {
	// Test ConvertSamplesToStereo
	monoSamples := []int16{100, 200, 300}
	stereoSamples := ConvertSamplesToStereo(monoSamples)
	expected := []int16{100, 100, 200, 200, 300, 300}
	assert.Equal(t, expected, stereoSamples)

	// Test MixStereoSamples
	left := []int16{100, 200}
	right := []int16{150, 250}
	mixed := MixStereoSamples(left, right)
	expectedMixed := []int16{100, 150, 200, 250}
	assert.Equal(t, expectedMixed, mixed)

	// Test ApplyVolume
	samples := []int16{1000, 2000, 3000}
	ApplyVolume(samples, 0.5)
	expectedVolume := []int16{500, 1000, 1500}
	assert.Equal(t, expectedVolume, samples)
}
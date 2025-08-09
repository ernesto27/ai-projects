package apu

import (
	"fmt"
)

// APU represents the Game Boy Audio Processing Unit
type APU struct {
	// Sound channels
	channel1 *Channel1 // Square wave with sweep
	channel2 *Channel2 // Square wave
	channel3 *Channel3 // Wave pattern
	channel4 *Channel4 // Noise generator

	// Audio control
	mixer   *Mixer
	enabled bool

	// Master registers
	nr50 uint8 // Master volume & VIN panning (0xFF24)
	nr51 uint8 // Sound panning (0xFF25)
	nr52 uint8 // Sound on/off (0xFF26)

	// Timing
	frameSequencer uint8  // 8-step frame sequencer (512 Hz)
	frameCounter   uint16 // Counts to 8192 CPU cycles per step
	cycles         uint64 // Total cycles processed

	// Audio output
	sampleRate   float64   // Target sample rate (e.g., 44100 Hz)
	sampleBuffer []float32 // Audio sample buffer
	sampleIndex  int       // Current position in sample buffer
}

// AudioInterface defines the interface for audio output
type AudioInterface interface {
	Initialize(sampleRate int, bufferSize int) error
	QueueAudio(samples []float32) error
	GetQueuedBytes() int
	Close() error
}

// NewAPU creates a new APU instance
func NewAPU() *APU {
	apu := &APU{
		channel1:     NewChannel1(),
		channel2:     NewChannel2(),
		channel3:     NewChannel3(),
		channel4:     NewChannel4(),
		mixer:        NewMixer(),
		sampleRate:   44100.0, // Standard sample rate
		sampleBuffer: make([]float32, 1024), // 1KB audio buffer
	}

	apu.Reset()
	return apu
}

// Reset initializes the APU to its power-on state
func (apu *APU) Reset() {
	// Reset all channels
	apu.channel1.Reset()
	apu.channel2.Reset()
	apu.channel3.Reset()
	apu.channel4.Reset()

	// Reset master registers
	apu.nr50 = 0x77 // Max volume both channels
	apu.nr51 = 0xF3 // All channels enabled on both sides
	apu.nr52 = 0xF1 // APU enabled, all channels enabled

	// Reset timing
	apu.frameSequencer = 0
	apu.frameCounter = 0
	apu.cycles = 0

	// APU starts enabled
	apu.enabled = true

	// Reset mixer
	apu.mixer.Reset()
}

// Update processes APU for the given number of CPU cycles
func (apu *APU) Update(cycles uint8) {
	if !apu.enabled {
		return
	}

	apu.cycles += uint64(cycles)
	apu.frameCounter += uint16(cycles)

	// Frame sequencer runs at 512 Hz (8192 CPU cycles per step)
	if apu.frameCounter >= 8192 {
		apu.frameCounter -= 8192
		apu.stepFrameSequencer()
	}

	// Update all channels
	apu.channel1.Update(cycles)
	apu.channel2.Update(cycles)
	apu.channel3.Update(cycles)
	apu.channel4.Update(cycles)

	// Generate audio samples
	apu.generateSamples(cycles)
}

// stepFrameSequencer advances the frame sequencer one step
func (apu *APU) stepFrameSequencer() {
	// Frame sequencer pattern (8 steps, 512 Hz):
	// Step 0: Length
	// Step 1: Nothing  
	// Step 2: Length + Sweep
	// Step 3: Nothing
	// Step 4: Length
	// Step 5: Nothing
	// Step 6: Length + Sweep  
	// Step 7: Envelope

	switch apu.frameSequencer {
	case 0, 2, 4, 6: // Length counter steps
		apu.channel1.StepLength()
		apu.channel2.StepLength()
		apu.channel3.StepLength()
		apu.channel4.StepLength()

		if apu.frameSequencer == 2 || apu.frameSequencer == 6 { // Sweep steps
			apu.channel1.StepSweep()
		}

	case 7: // Envelope step
		apu.channel1.StepEnvelope()
		apu.channel2.StepEnvelope()
		apu.channel4.StepEnvelope()
	}

	apu.frameSequencer = (apu.frameSequencer + 1) % 8
}

// generateSamples creates audio samples for the given CPU cycles
func (apu *APU) generateSamples(cycles uint8) {
	// Game Boy CPU runs at ~4.194304 MHz
	// We want to generate samples at our target sample rate
	// Calculate how many samples to generate for these cycles
	
	cpuFreq := 4194304.0 // Game Boy CPU frequency
	samplesNeeded := float64(cycles) * apu.sampleRate / cpuFreq
	
	// Generate the required number of samples
	for i := 0.0; i < samplesNeeded; i++ {
		leftSample, rightSample := apu.mixer.Mix(
			apu.channel1.GetSample(),
			apu.channel2.GetSample(),
			apu.channel3.GetSample(),
			apu.channel4.GetSample(),
			apu.nr50,
			apu.nr51,
		)

		// Store samples (interleaved stereo)
		if apu.sampleIndex < len(apu.sampleBuffer)-1 {
			apu.sampleBuffer[apu.sampleIndex] = leftSample
			apu.sampleBuffer[apu.sampleIndex+1] = rightSample
			apu.sampleIndex += 2
		}
	}
}

// GetSamples returns the current audio samples and resets the buffer
func (apu *APU) GetSamples() []float32 {
	if apu.sampleIndex == 0 {
		return nil
	}

	// Copy samples and reset buffer
	samples := make([]float32, apu.sampleIndex)
	copy(samples, apu.sampleBuffer[:apu.sampleIndex])
	apu.sampleIndex = 0

	return samples
}

// ReadByte reads from an APU register
func (apu *APU) ReadByte(address uint16) uint8 {
	switch {
	case address >= 0xFF10 && address <= 0xFF14: // Channel 1
		return apu.channel1.ReadRegister(uint8(address - 0xFF10))
	case address >= 0xFF16 && address <= 0xFF19: // Channel 2
		return apu.channel2.ReadRegister(uint8(address - 0xFF16))
	case address >= 0xFF1A && address <= 0xFF1E: // Channel 3
		return apu.channel3.ReadRegister(uint8(address - 0xFF1A))
	case address >= 0xFF20 && address <= 0xFF23: // Channel 4
		return apu.channel4.ReadRegister(uint8(address - 0xFF20))
	case address == 0xFF24: // NR50 - Master volume & VIN panning
		return apu.nr50
	case address == 0xFF25: // NR51 - Sound panning  
		return apu.nr51
	case address == 0xFF26: // NR52 - Sound on/off
		return apu.nr52
	case address >= 0xFF30 && address <= 0xFF3F: // Wave RAM
		return apu.channel3.ReadWaveRAM(uint8(address - 0xFF30))
	default:
		return 0xFF // Unmapped APU register
	}
}

// WriteByte writes to an APU register
func (apu *APU) WriteByte(address uint16, value uint8) {
	// If APU is disabled, only NR52 writes are allowed
	if !apu.enabled && address != 0xFF26 {
		return
	}

	switch {
	case address >= 0xFF10 && address <= 0xFF14: // Channel 1
		apu.channel1.WriteRegister(uint8(address-0xFF10), value)
	case address >= 0xFF16 && address <= 0xFF19: // Channel 2
		apu.channel2.WriteRegister(uint8(address-0xFF16), value)
	case address >= 0xFF1A && address <= 0xFF1E: // Channel 3
		apu.channel3.WriteRegister(uint8(address-0xFF1A), value)
	case address >= 0xFF20 && address <= 0xFF23: // Channel 4
		apu.channel4.WriteRegister(uint8(address-0xFF20), value)
	case address == 0xFF24: // NR50 - Master volume & VIN panning
		apu.nr50 = value
	case address == 0xFF25: // NR51 - Sound panning
		apu.nr51 = value
	case address == 0xFF26: // NR52 - Sound on/off
		apu.writeNR52(value)
	case address >= 0xFF30 && address <= 0xFF3F: // Wave RAM
		apu.channel3.WriteWaveRAM(uint8(address-0xFF30), value)
	}
}

// writeNR52 handles writes to the master sound control register
func (apu *APU) writeNR52(value uint8) {
	wasEnabled := apu.enabled
	apu.enabled = (value & 0x80) != 0

	if wasEnabled && !apu.enabled {
		// APU was turned off - clear all registers except wave RAM
		apu.clearRegisters()
	}

	// Update NR52 with current channel status
	apu.updateNR52()
}

// clearRegisters clears all APU registers when APU is disabled
func (apu *APU) clearRegisters() {
	// Clear all APU registers except wave RAM (0xFF30-0xFF3F)
	for addr := uint16(0xFF10); addr <= 0xFF25; addr++ {
		if addr != 0xFF26 { // Don't clear NR52
			apu.WriteByte(addr, 0)
		}
	}

	// Clear channel states
	apu.channel1.Reset()
	apu.channel2.Reset()
	apu.channel3.Reset()
	apu.channel4.Reset()
}

// updateNR52 updates the NR52 register with current channel status
func (apu *APU) updateNR52() {
	apu.nr52 = 0
	if apu.enabled {
		apu.nr52 |= 0x80 // APU enabled bit
	}

	// Set channel enable bits based on channel status
	if apu.channel1.IsEnabled() {
		apu.nr52 |= 0x01
	}
	if apu.channel2.IsEnabled() {
		apu.nr52 |= 0x02
	}
	if apu.channel3.IsEnabled() {
		apu.nr52 |= 0x04
	}
	if apu.channel4.IsEnabled() {
		apu.nr52 |= 0x08
	}
}

// IsEnabled returns whether the APU is enabled
func (apu *APU) IsEnabled() bool {
	return apu.enabled
}

// GetChannelStatus returns the status of all channels
func (apu *APU) GetChannelStatus() (bool, bool, bool, bool) {
	return apu.channel1.IsEnabled(),
		apu.channel2.IsEnabled(),
		apu.channel3.IsEnabled(),
		apu.channel4.IsEnabled()
}

// SetSampleRate sets the target audio sample rate
func (apu *APU) SetSampleRate(rate float64) {
	apu.sampleRate = rate
}

// String returns a string representation of the APU state
func (apu *APU) String() string {
	return fmt.Sprintf("APU{enabled=%t, nr50=0x%02X, nr51=0x%02X, nr52=0x%02X, frame=%d}",
		apu.enabled, apu.nr50, apu.nr51, apu.nr52, apu.frameSequencer)
}
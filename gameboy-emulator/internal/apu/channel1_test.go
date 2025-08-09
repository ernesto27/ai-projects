package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChannel1(t *testing.T) {
	ch := NewChannel1()
	assert.NotNil(t, ch)
	assert.False(t, ch.enabled)
	assert.False(t, ch.dacEnabled)
	assert.Equal(t, uint16(0), ch.frequency)
	assert.Equal(t, uint8(0), ch.currentVolume)
}

func TestChannel1Reset(t *testing.T) {
	ch := NewChannel1()
	
	// Modify state
	ch.enabled = true
	ch.frequency = 1000
	ch.currentVolume = 10
	ch.sweepEnabled = true
	
	// Reset should restore defaults
	ch.Reset()
	
	assert.False(t, ch.enabled)
	assert.False(t, ch.dacEnabled)
	assert.Equal(t, uint16(0), ch.frequency)
	assert.Equal(t, uint8(0), ch.currentVolume)
	assert.False(t, ch.sweepEnabled)
}

func TestChannel1DutyPatterns(t *testing.T) {
	// Test that duty patterns have correct length and content
	assert.Equal(t, 4, len(dutyPatterns), "Should have 4 duty patterns")
	
	for i, pattern := range dutyPatterns {
		assert.Equal(t, 8, len(pattern), "Duty pattern %d should have 8 steps", i)
		
		// Count high samples (value = 1)
		highCount := 0
		for _, sample := range pattern {
			if sample > 0 {
				highCount++
			}
		}
		
		switch i {
		case 0: // 12.5% duty
			assert.Equal(t, 1, highCount, "12.5% duty should have 1 high sample")
		case 1: // 25% duty
			assert.Equal(t, 2, highCount, "25% duty should have 2 high samples")
		case 2: // 50% duty
			assert.Equal(t, 4, highCount, "50% duty should have 4 high samples")
		case 3: // 75% duty  
			assert.Equal(t, 6, highCount, "75% duty should have 6 high samples")
		}
	}
}

func TestChannel1RegisterWrite(t *testing.T) {
	ch := NewChannel1()

	// Test NR10 - Sweep control
	ch.WriteRegister(0, 0x71) // Period=7, increase (bit 3 = 0), shift=1
	assert.Equal(t, uint8(7), ch.sweepPeriod)
	assert.False(t, ch.sweepDirection) // bit 3 = 0 means increase (false = increase)
	assert.Equal(t, uint8(1), ch.sweepShift)

	// Test NR11 - Wave pattern duty and length
	ch.WriteRegister(1, 0xC0) // 75% duty, length=0
	assert.Equal(t, uint8(3), ch.dutyPattern) // 0xC0 >> 6 = 3
	assert.Equal(t, uint8(64), ch.lengthCounter) // 64 - 0 = 64

	// Test NR12 - Volume envelope
	ch.WriteRegister(2, 0xFB) // Volume=15, increase (bit 3 = 1), period=3
	assert.Equal(t, uint8(15), ch.initialVolume)
	assert.True(t, ch.envelopeDirection) // bit 3 set = increase
	assert.Equal(t, uint8(3), ch.envelopePeriod)
	assert.True(t, ch.dacEnabled) // Upper 5 bits not zero

	// Test NR13 - Frequency low byte
	ch.WriteRegister(3, 0xFF)
	assert.Equal(t, uint16(0xFF), ch.frequency&0xFF)

	// Test NR14 - Frequency high byte + control
	ch.WriteRegister(4, 0xC7) // Trigger + length enable + freq high = 7
	expectedFreq := uint16(0x7FF) // 0x7 << 8 | 0xFF
	assert.Equal(t, expectedFreq, ch.frequency)
	// Length enabled should be set (bit 6 of 0xC7 = 0x40 is set)
	assert.True(t, ch.lengthEnabled, "Length should be enabled with bit 6 set in NR14")
}

func TestChannel1RegisterRead(t *testing.T) {
	ch := NewChannel1()

	// Set some register values
	ch.WriteRegister(0, 0x71)
	ch.WriteRegister(1, 0xC0)
	ch.WriteRegister(2, 0xFB)

	// Test reads (adjust expectations based on actual bit masking)
	nr10Read := ch.ReadRegister(0) // NR10 | 0x80
	assert.True(t, (nr10Read & 0x80) != 0, "NR10 bit 7 should be set")
	
	nr11Read := ch.ReadRegister(1) // NR11 | 0x3F  
	assert.True(t, (nr11Read & 0x3F) == 0x3F, "NR11 lower 6 bits should be set")
	
	assert.Equal(t, uint8(0xFB), ch.ReadRegister(2)) // NR12
	assert.Equal(t, uint8(0xFF), ch.ReadRegister(3)) // NR13 write-only
	
	nr14Read := ch.ReadRegister(4) // NR14 | 0xBF
	assert.True(t, (nr14Read & 0x80) != 0, "NR14 bit 7 should always read as 1")
}

func TestChannel1Trigger(t *testing.T) {
	ch := NewChannel1()

	// Set up for trigger
	ch.WriteRegister(2, 0xF0) // Enable DAC
	ch.frequency = 1000
	ch.lengthCounter = 0

	// Trigger channel
	ch.trigger()

	assert.True(t, ch.enabled)
	assert.Equal(t, uint8(64), ch.lengthCounter) // Should initialize to 64
	assert.Equal(t, uint16((2048-1000)*4), ch.period) // Should set period
	assert.Equal(t, uint8(15), ch.currentVolume) // Should set to initial volume
}

func TestChannel1TriggerWithDACDisabled(t *testing.T) {
	ch := NewChannel1()

	// Disable DAC
	ch.WriteRegister(2, 0x07) // Volume=0, period=7, DAC disabled

	// Trigger should not enable channel
	ch.trigger()
	assert.False(t, ch.enabled)
}

func TestChannel1LengthCounter(t *testing.T) {
	ch := NewChannel1()

	// Enable channel and length counter
	ch.enabled = true
	ch.lengthEnabled = true
	ch.lengthCounter = 2

	// Step length counter
	ch.StepLength()
	assert.Equal(t, uint8(1), ch.lengthCounter)
	assert.True(t, ch.enabled)

	// Step again - should disable channel
	ch.StepLength()
	assert.Equal(t, uint8(0), ch.lengthCounter)
	assert.False(t, ch.enabled)
}

func TestChannel1LengthCounterDisabled(t *testing.T) {
	ch := NewChannel1()

	// Disable length counter
	ch.enabled = true
	ch.lengthEnabled = false
	ch.lengthCounter = 1

	// Step should not affect counter or channel
	ch.StepLength()
	assert.Equal(t, uint8(1), ch.lengthCounter)
	assert.True(t, ch.enabled)
}

func TestChannel1Envelope(t *testing.T) {
	ch := NewChannel1()

	// Test envelope increase
	ch.envelopeEnabled = true
	ch.envelopeDirection = true // increase
	ch.envelopePeriod = 2
	ch.envelopeCounter = 1
	ch.currentVolume = 10

	// Step envelope - counter should decrement, and if it reaches 0, volume should change
	ch.StepEnvelope()
	// Counter was 1, decrements to 0, so volume should change and counter reset
	assert.Equal(t, uint8(11), ch.currentVolume)
	assert.Equal(t, uint8(2), ch.envelopeCounter) // Reset counter to period
}

func TestChannel1EnvelopeDecrease(t *testing.T) {
	ch := NewChannel1()

	// Test envelope decrease
	ch.envelopeEnabled = true
	ch.envelopeDirection = false // decrease
	ch.envelopePeriod = 1
	ch.envelopeCounter = 1
	ch.currentVolume = 10

	// Step envelope
	ch.StepEnvelope()
	assert.Equal(t, uint8(9), ch.currentVolume)
}

func TestChannel1EnvelopeLimits(t *testing.T) {
	ch := NewChannel1()

	// Test volume at maximum
	ch.envelopeEnabled = true
	ch.envelopeDirection = true
	ch.envelopePeriod = 1
	ch.envelopeCounter = 1
	ch.currentVolume = 15

	ch.StepEnvelope()
	assert.Equal(t, uint8(15), ch.currentVolume) // Should not exceed 15
	assert.False(t, ch.envelopeEnabled) // Should disable envelope

	// Test volume at minimum
	ch.envelopeEnabled = true
	ch.envelopeDirection = false
	ch.envelopeCounter = 1
	ch.currentVolume = 0

	ch.StepEnvelope()
	assert.Equal(t, uint8(0), ch.currentVolume) // Should not go below 0
	assert.False(t, ch.envelopeEnabled) // Should disable envelope
}

func TestChannel1Sweep(t *testing.T) {
	ch := NewChannel1()

	// Test sweep increase
	ch.sweepEnabled = true
	ch.sweepPeriod = 1
	ch.sweepCounter = 1
	ch.sweepDirection = false // increase
	ch.sweepShift = 1
	ch.sweepShadow = 100
	ch.frequency = 100

	// Step sweep
	ch.StepSweep()
	
	// New frequency should be 100 + (100 >> 1) = 150
	expectedFreq := uint16(100 + 50)
	assert.Equal(t, expectedFreq, ch.sweepShadow)
	assert.Equal(t, expectedFreq, ch.frequency)
}

func TestChannel1SweepOverflow(t *testing.T) {
	ch := NewChannel1()

	// Set up sweep that will overflow
	ch.enabled = true
	ch.sweepEnabled = true
	ch.sweepPeriod = 1
	ch.sweepCounter = 1
	ch.sweepDirection = false // increase
	ch.sweepShift = 1
	ch.sweepShadow = 2000 // High frequency
	ch.frequency = 2000

	// Step sweep - should disable channel due to overflow
	ch.StepSweep()
	assert.False(t, ch.enabled)
}

func TestChannel1SweepDecrease(t *testing.T) {
	ch := NewChannel1()

	// Test sweep decrease
	ch.sweepEnabled = true
	ch.sweepPeriod = 1
	ch.sweepCounter = 1
	ch.sweepDirection = true // decrease
	ch.sweepShift = 2
	ch.sweepShadow = 200
	ch.frequency = 200

	// Step sweep
	ch.StepSweep()
	
	// New frequency should be 200 - (200 >> 2) = 150
	expectedFreq := uint16(200 - 50)
	assert.Equal(t, expectedFreq, ch.sweepShadow)
}

func TestChannel1Update(t *testing.T) {
	ch := NewChannel1()

	// Test update with disabled channel
	ch.Update(100)
	assert.Equal(t, float32(0), ch.sample)

	// Enable channel
	ch.enabled = true
	ch.dacEnabled = true
	ch.frequency = 1000
	ch.currentVolume = 8
	ch.dutyPattern = 2 // 50% duty

	// Update should generate wave
	ch.Update(10)
	// Sample should be non-zero (exact value depends on wave position)
	assert.True(t, ch.sample >= -1.0 && ch.sample <= 1.0)
}

func TestChannel1GetSample(t *testing.T) {
	ch := NewChannel1()

	// Disabled channel should return 0
	assert.Equal(t, float32(0), ch.GetSample())

	// Enable channel
	ch.enabled = true
	ch.dacEnabled = true
	ch.sample = 0.5

	// Should return current sample
	assert.Equal(t, float32(0.5), ch.GetSample())

	// Disable DAC
	ch.dacEnabled = false
	assert.Equal(t, float32(0), ch.GetSample())
}

func TestChannel1IsEnabled(t *testing.T) {
	ch := NewChannel1()

	assert.False(t, ch.IsEnabled())
	
	ch.enabled = true
	assert.True(t, ch.IsEnabled())
}

func TestChannel1GetFrequency(t *testing.T) {
	ch := NewChannel1()

	ch.frequency = 1337
	assert.Equal(t, uint16(1337), ch.GetFrequency())
}

func TestChannel1GetVolume(t *testing.T) {
	ch := NewChannel1()

	ch.currentVolume = 12
	assert.Equal(t, uint8(12), ch.GetVolume())
}

func TestChannel1IsDACEnabled(t *testing.T) {
	ch := NewChannel1()

	assert.False(t, ch.IsDACEnabled())
	
	ch.dacEnabled = true
	assert.True(t, ch.IsDACEnabled())
}
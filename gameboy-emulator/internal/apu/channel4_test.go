package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChannel4(t *testing.T) {
	ch := NewChannel4()
	assert.NotNil(t, ch)
	assert.False(t, ch.enabled)
	assert.False(t, ch.dacEnabled)
	assert.Equal(t, uint16(0x7FFF), ch.lfsr) // Should be initialized to all 1s except bit 15
}

func TestChannel4Reset(t *testing.T) {
	ch := NewChannel4()
	
	// Modify state
	ch.enabled = true
	ch.lfsr = 0x1234
	ch.currentVolume = 10
	ch.widthMode = true
	
	// Reset should restore defaults
	ch.Reset()
	
	assert.False(t, ch.enabled)
	assert.False(t, ch.dacEnabled)
	assert.Equal(t, uint16(0x7FFF), ch.lfsr)
	assert.Equal(t, uint8(0), ch.currentVolume)
	assert.False(t, ch.widthMode)
}

func TestNoiseDividersConstant(t *testing.T) {
	// Test that noise dividers are correctly defined
	expected := [8]uint16{8, 16, 32, 48, 64, 80, 96, 112}
	assert.Equal(t, expected, noiseDividers)
}

func TestChannel4RegisterWrite(t *testing.T) {
	ch := NewChannel4()

	// Test NR41 - Length
	ch.WriteRegister(0, 0x20) // Length = 32
	assert.Equal(t, uint8(32), ch.lengthCounter) // 64 - 32 = 32

	// Test NR42 - Volume envelope
	ch.WriteRegister(1, 0xFB) // Volume=15, increase (bit 3 = 1), period=3
	assert.Equal(t, uint8(15), ch.initialVolume)
	assert.True(t, ch.envelopeDirection) // bit 3 set = increase
	assert.Equal(t, uint8(3), ch.envelopePeriod)
	assert.True(t, ch.dacEnabled) // Upper 5 bits not zero

	// Test NR43 - Polynomial counter
	ch.WriteRegister(2, 0xAB) // Shift=10, width=1, divider=3
	assert.Equal(t, uint8(10), ch.clockShift) // 0xA0 >> 4 = 10
	assert.True(t, ch.widthMode) // Bit 3 set
	assert.Equal(t, uint8(3), ch.clockDivider) // 0x0B & 0x07 = 3

	// Test NR44 - Control
	ch.WriteRegister(3, 0xC0) // Trigger + length enable
	assert.True(t, ch.lengthEnabled)
}

func TestChannel4RegisterRead(t *testing.T) {
	ch := NewChannel4()

	// Set some register values
	ch.WriteRegister(1, 0xFB)
	ch.WriteRegister(2, 0xAB)

	// Test reads
	assert.Equal(t, uint8(0xFF), ch.ReadRegister(0)) // NR41 write-only
	assert.Equal(t, uint8(0xFB), ch.ReadRegister(1)) // NR42
	assert.Equal(t, uint8(0xAB), ch.ReadRegister(2)) // NR43
	
	nr44Read := ch.ReadRegister(3) // NR44 | 0xBF
	assert.True(t, (nr44Read & 0x80) != 0, "NR44 bit 7 should always read as 1")
}

func TestChannel4Trigger(t *testing.T) {
	ch := NewChannel4()

	// Set up for trigger
	ch.WriteRegister(1, 0xF0) // Enable DAC
	ch.lengthCounter = 0
	ch.lfsr = 0x1234 // Non-default value

	// Trigger channel
	ch.trigger()

	assert.True(t, ch.enabled)
	assert.Equal(t, uint8(64), ch.lengthCounter) // Should initialize to 64
	assert.Equal(t, uint16(0x7FFF), ch.lfsr) // Should reset LFSR
	assert.Equal(t, uint8(15), ch.currentVolume) // Should set to initial volume
}

func TestChannel4TriggerWithDACDisabled(t *testing.T) {
	ch := NewChannel4()

	// Disable DAC
	ch.WriteRegister(1, 0x07) // Volume=0, period=7, DAC disabled

	// Trigger should not enable channel
	ch.trigger()
	assert.False(t, ch.enabled)
}

func TestChannel4LengthCounter(t *testing.T) {
	ch := NewChannel4()

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

func TestChannel4Envelope(t *testing.T) {
	ch := NewChannel4()

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

func TestChannel4LFSR15Bit(t *testing.T) {
	ch := NewChannel4()

	// Set up 15-bit LFSR
	ch.widthMode = false
	ch.lfsr = 0x7FFE // All 1s except bits 0 and 15

	// Step LFSR
	ch.stepLFSR()

	// Check that LFSR shifted and XOR was applied
	assert.NotEqual(t, uint16(0x7FFE), ch.lfsr, "LFSR should have changed")
	
	// Bit 15 should be clear (original bit 15 was 0)
	assert.Equal(t, uint16(0), ch.lfsr&0x8000, "Bit 15 should be 0")
}

func TestChannel4LFSR7Bit(t *testing.T) {
	ch := NewChannel4()

	// Set up 7-bit LFSR
	ch.widthMode = true
	ch.lfsr = 0x7FFE // All 1s except bits 0 and 15

	// Step LFSR
	ch.stepLFSR()

	// In 7-bit mode, bit 6 should also be affected
	// The exact pattern depends on the XOR result
	assert.NotEqual(t, uint16(0x7FFE), ch.lfsr, "LFSR should have changed")
}

func TestChannel4LFSROutput(t *testing.T) {
	ch := NewChannel4()
	
	// Enable channel
	ch.enabled = true
	ch.dacEnabled = true
	ch.currentVolume = 15

	// Test with LFSR bit 0 = 0 (should output positive)
	ch.lfsr = 0x7FFE // Bit 0 = 0
	ch.updateNoise(0) // Don't advance, just generate sample
	assert.True(t, ch.sample >= 0, "Should output positive when LFSR bit 0 = 0")

	// Test with LFSR bit 0 = 1 (should output negative)
	ch.lfsr = 0x7FFF // Bit 0 = 1
	ch.updateNoise(0) // Don't advance, just generate sample
	assert.True(t, ch.sample <= 0, "Should output negative when LFSR bit 0 = 1")
}

func TestChannel4Update(t *testing.T) {
	ch := NewChannel4()

	// Test update with disabled channel
	ch.Update(100)
	assert.Equal(t, float32(0), ch.sample)

	// Enable channel
	ch.enabled = true
	ch.dacEnabled = true
	ch.currentVolume = 8
	ch.clockDivider = 0 // Divider = 8
	ch.clockShift = 0   // No shift

	// Update should generate noise
	ch.Update(10)
	// Sample should be non-zero and within range
	assert.True(t, ch.sample >= -1.0 && ch.sample <= 1.0)
}

func TestChannel4PeriodCalculation(t *testing.T) {
	ch := NewChannel4()

	testCases := []struct {
		divider     uint8
		shift       uint8
		expectedMin uint16 // Period should be at least this
	}{
		{0, 0, 8},   // divider=8, shift=0 -> period=8
		{0, 1, 16},  // divider=8, shift=1 -> period=16
		{1, 0, 16},  // divider=16, shift=0 -> period=16
		{1, 1, 32},  // divider=16, shift=1 -> period=32
		{7, 4, 1792}, // divider=112, shift=4 -> period=1792
	}

	for _, tc := range testCases {
		ch.clockDivider = tc.divider
		ch.clockShift = tc.shift
		
		// Calculate period
		divider := noiseDividers[tc.divider]
		expectedPeriod := divider << tc.shift
		
		assert.Equal(t, tc.expectedMin, expectedPeriod,
			"Period calculation for divider=%d, shift=%d", tc.divider, tc.shift)
	}
}

func TestChannel4GetSample(t *testing.T) {
	ch := NewChannel4()

	// Disabled channel should return 0
	assert.Equal(t, float32(0), ch.GetSample())

	// Enable channel
	ch.enabled = true
	ch.dacEnabled = true
	ch.sample = 0.75

	// Should return current sample
	assert.Equal(t, float32(0.75), ch.GetSample())

	// Disable DAC
	ch.dacEnabled = false
	assert.Equal(t, float32(0), ch.GetSample())
}

func TestChannel4IsEnabled(t *testing.T) {
	ch := NewChannel4()

	assert.False(t, ch.IsEnabled())
	
	ch.enabled = true
	assert.True(t, ch.IsEnabled())
}

func TestChannel4IsDACEnabled(t *testing.T) {
	ch := NewChannel4()

	assert.False(t, ch.IsDACEnabled())
	
	ch.dacEnabled = true
	assert.True(t, ch.IsDACEnabled())
}

func TestChannel4GetLFSR(t *testing.T) {
	ch := NewChannel4()

	ch.lfsr = 0x1234
	assert.Equal(t, uint16(0x1234), ch.GetLFSR())
}

func TestChannel4GetVolume(t *testing.T) {
	ch := NewChannel4()

	ch.currentVolume = 9
	assert.Equal(t, uint8(9), ch.GetVolume())
}

func TestChannel4IsWidthMode(t *testing.T) {
	ch := NewChannel4()

	assert.False(t, ch.IsWidthMode())
	
	ch.widthMode = true
	assert.True(t, ch.IsWidthMode())
}

func TestChannel4DACEnableDisable(t *testing.T) {
	ch := NewChannel4()

	// Initially DAC should be disabled
	assert.False(t, ch.dacEnabled)

	// Write volume envelope with non-zero upper bits
	ch.WriteRegister(1, 0xF0) // Volume=15, no envelope
	assert.True(t, ch.dacEnabled)

	// Write volume envelope with zero upper bits
	ch.WriteRegister(1, 0x07) // Volume=0, envelope period=7
	assert.False(t, ch.dacEnabled)
	assert.False(t, ch.enabled) // Should also disable channel
}

func TestChannel4NoiseFrequencies(t *testing.T) {
	// Test different noise configurations produce different periods
	periods := make(map[uint16]bool)
	
	for divider := uint8(0); divider < 8; divider++ {
		for shift := uint8(0); shift < 8; shift++ {
			expectedPeriod := noiseDividers[divider] << shift
			periods[expectedPeriod] = true
		}
	}
	
	// Should have many unique periods
	assert.True(t, len(periods) > 20, "Should generate many different periods")
}
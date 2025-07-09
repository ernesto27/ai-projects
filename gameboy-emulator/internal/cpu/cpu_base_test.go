package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCPU tests CPU initialization
func TestNewCPU(t *testing.T) {
	cpu := NewCPU()

	// Test initial register values (Game Boy boot state)
	assert.Equal(t, uint8(0x01), cpu.A, "A register should be 0x01")
	assert.Equal(t, uint8(0xB0), cpu.F, "F register should be 0xB0")
	assert.Equal(t, uint8(0x00), cpu.B, "B register should be 0x00")
	assert.Equal(t, uint8(0x13), cpu.C, "C register should be 0x13")
	assert.Equal(t, uint8(0x00), cpu.D, "D register should be 0x00")
	assert.Equal(t, uint8(0xD8), cpu.E, "E register should be 0xD8")
	assert.Equal(t, uint8(0x01), cpu.H, "H register should be 0x01")
	assert.Equal(t, uint8(0x4D), cpu.L, "L register should be 0x4D")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be 0xFFFE")
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should be 0x0100")
	assert.False(t, cpu.Halted, "CPU should not be halted")
	assert.False(t, cpu.Stopped, "CPU should not be stopped")
}

// TestCPUReset tests CPU reset functionality
func TestCPUReset(t *testing.T) {
	cpu := NewCPU()

	// Modify CPU state
	cpu.A = 0xFF
	cpu.PC = 0x1234
	cpu.SP = 0x5678
	cpu.Halted = true
	cpu.Stopped = true

	// Reset CPU
	cpu.Reset()

	// Verify reset to initial state
	assert.Equal(t, uint8(0x01), cpu.A, "A should be reset to 0x01")
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should be reset to 0x0100")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be reset to 0xFFFE")
	assert.False(t, cpu.Halted, "CPU should not be halted after reset")
	assert.False(t, cpu.Stopped, "CPU should not be stopped after reset")
}

// TestFlagConstants tests that flag constants have correct values
func TestFlagConstants(t *testing.T) {
	assert.Equal(t, 0x80, int(FlagZ), "Zero flag constant should be 0x80")
	assert.Equal(t, 0x40, int(FlagN), "Subtract flag constant should be 0x40")
	assert.Equal(t, 0x20, int(FlagH), "Half-carry flag constant should be 0x20")
	assert.Equal(t, 0x10, int(FlagC), "Carry flag constant should be 0x10")
}

// TestMultipleFlags tests setting multiple flags at once
func TestMultipleFlags(t *testing.T) {
	cpu := NewCPU()

	// Start with a clean flag register for this test
	cpu.F = 0x00

	// Set multiple flags
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagC, true)

	// Check both are set
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")

	// Check other flags are not affected
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")
}

// TestRegisterPairBoundaries tests edge cases with register pairs
func TestRegisterPairBoundaries(t *testing.T) {
	cpu := NewCPU()

	// Test maximum values
	cpu.SetAF(0xFFFF)
	assert.Equal(t, uint8(0xFF), cpu.A, "A should be set to 0xFF")
	assert.Equal(t, uint8(0xFF), cpu.F, "F should be set to 0xFF")

	// Test minimum values
	cpu.SetBC(0x0000)
	assert.Equal(t, uint8(0x00), cpu.B, "B should be set to 0x00")
	assert.Equal(t, uint8(0x00), cpu.C, "C should be set to 0x00")

	// Test mixed values
	cpu.SetDE(0x00FF)
	assert.Equal(t, uint8(0x00), cpu.D, "D should be set to 0x00")
	assert.Equal(t, uint8(0xFF), cpu.E, "E should be set to 0xFF")

	cpu.SetHL(0xFF00)
	assert.Equal(t, uint8(0xFF), cpu.H, "H should be set to 0xFF")
	assert.Equal(t, uint8(0x00), cpu.L, "L should be set to 0x00")
}

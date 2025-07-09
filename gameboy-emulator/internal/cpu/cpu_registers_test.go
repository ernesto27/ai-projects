package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAFRegisterPair tests AF register pair operations
func TestAFRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting AF pair
	cpu.A = 0x12
	cpu.F = 0x34
	assert.Equal(t, uint16(0x1234), cpu.GetAF(), "AF pair should combine A and F registers")

	// Test setting AF pair
	cpu.SetAF(0x5678)
	assert.Equal(t, uint8(0x56), cpu.A, "A should be set to high byte")
	assert.Equal(t, uint8(0x78), cpu.F, "F should be set to low byte")
}

// TestBCRegisterPair tests BC register pair operations
func TestBCRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting BC pair
	cpu.B = 0xAB
	cpu.C = 0xCD
	assert.Equal(t, uint16(0xABCD), cpu.GetBC(), "BC pair should combine B and C registers")

	// Test setting BC pair
	cpu.SetBC(0x1234)
	assert.Equal(t, uint8(0x12), cpu.B, "B should be set to high byte")
	assert.Equal(t, uint8(0x34), cpu.C, "C should be set to low byte")
}

// TestHLRegisterPair tests HL register pair operations
func TestHLRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting HL pair
	cpu.H = 0x42
	cpu.L = 0x24
	assert.Equal(t, uint16(0x4224), cpu.GetHL(), "HL pair should combine H and L registers")

	// Test setting HL pair
	cpu.SetHL(0xBEEF)
	assert.Equal(t, uint8(0xBE), cpu.H, "H should be set to high byte")
	assert.Equal(t, uint8(0xEF), cpu.L, "L should be set to low byte")
}

// TestDERegisterPair tests DE register pair operations
func TestDERegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting DE pair
	cpu.D = 0xEF
	cpu.E = 0x01
	assert.Equal(t, uint16(0xEF01), cpu.GetDE(), "DE pair should combine D and E registers")

	// Test setting DE pair
	cpu.SetDE(0x9876)
	assert.Equal(t, uint8(0x98), cpu.D, "D should be set to high byte")
	assert.Equal(t, uint8(0x76), cpu.E, "E should be set to low byte")
}

// TestFlagOperations tests flag register operations
func TestFlagOperations(t *testing.T) {
	cpu := NewCPU()

	// Test setting flags
	cpu.SetFlag(FlagZ, true)
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")

	cpu.SetFlag(FlagN, true)
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")

	cpu.SetFlag(FlagH, true)
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set")

	cpu.SetFlag(FlagC, true)
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")

	// Test clearing flags
	cpu.SetFlag(FlagZ, false)
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be cleared")

	cpu.SetFlag(FlagN, false)
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be cleared")

	cpu.SetFlag(FlagH, false)
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be cleared")

	cpu.SetFlag(FlagC, false)
	assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be cleared")
}

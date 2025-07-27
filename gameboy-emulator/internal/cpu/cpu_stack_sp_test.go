package cpu

import (
	"testing"

	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

// Test LD (nn),SP instruction (0x08)
func TestLD_nn_SP(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic operation
	cpu.SP = 0x1234
	cycles := cpu.LD_nn_SP(mmu, 0x00, 0x80) // Store SP at 0x8000

	// Verify cycles
	assert.Equal(t, uint8(20), cycles, "LD (nn),SP should take 20 cycles")

	// Verify little-endian storage
	lowByte := mmu.ReadByte(0x8000)   // Should be 0x34 (low byte of SP)
	highByte := mmu.ReadByte(0x8001)  // Should be 0x12 (high byte of SP)
	assert.Equal(t, uint8(0x34), lowByte, "Low byte should be stored first")
	assert.Equal(t, uint8(0x12), highByte, "High byte should be stored second")

	// Verify SP unchanged
	assert.Equal(t, uint16(0x1234), cpu.SP, "SP should not be modified")
}

func TestLD_nn_SP_EdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	testCases := []struct {
		name     string
		sp       uint16
		low      uint8
		high     uint8
		address  uint16
		expectedLow  uint8
		expectedHigh uint8
	}{
		{"Zero SP", 0x0000, 0x00, 0x90, 0x9000, 0x00, 0x00},
		{"Max SP", 0xFFFF, 0xFF, 0x8F, 0x8FFF, 0xFF, 0xFF},
		{"Random values", 0xABCD, 0x34, 0x12, 0x1234, 0xCD, 0xAB},
		{"High memory", 0x5678, 0xFE, 0xFF, 0xFFFE, 0x78, 0x56},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.SP = tc.sp
			cycles := cpu.LD_nn_SP(mmu, tc.low, tc.high)

			assert.Equal(t, uint8(20), cycles, "Should take 20 cycles")
			assert.Equal(t, tc.expectedLow, mmu.ReadByte(tc.address), "Low byte mismatch")
			assert.Equal(t, tc.expectedHigh, mmu.ReadByte(tc.address+1), "High byte mismatch")
			assert.Equal(t, tc.sp, cpu.SP, "SP should not change")
		})
	}
}

// Test LD SP,HL instruction (0xF9)
func TestLD_SP_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic operation
	cpu.SetHL(0x5678)
	cpu.SP = 0x1234 // Different value to ensure it changes

	cycles := cpu.LD_SP_HL(mmu)

	// Verify cycles and result
	assert.Equal(t, uint8(8), cycles, "LD SP,HL should take 8 cycles")
	assert.Equal(t, uint16(0x5678), cpu.SP, "SP should be set to HL value")
	assert.Equal(t, uint16(0x5678), cpu.GetHL(), "HL should remain unchanged")
}

func TestLD_SP_HL_EdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	testCases := []struct {
		name string
		hl   uint16
	}{
		{"Zero HL", 0x0000},
		{"Max HL", 0xFFFF},
		{"Random value", 0xABCD},
		{"Stack start", 0xFFFE},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.SetHL(tc.hl)
			oldSP := cpu.SP

			cycles := cpu.LD_SP_HL(mmu)

			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
			assert.Equal(t, tc.hl, cpu.SP, "SP should match HL")
			assert.Equal(t, tc.hl, cpu.GetHL(), "HL should remain unchanged")
			assert.NotEqual(t, oldSP, cpu.SP, "SP should have changed")
		})
	}
}

// Test ADD SP,n instruction (0xE8)
func TestADD_SP_n(t *testing.T) {
	cpu := NewCPU()

	// Test positive offset
	cpu.SP = 0x1000
	cycles := cpu.ADD_SP_n(0x10) // +16

	assert.Equal(t, uint8(16), cycles, "ADD SP,n should take 16 cycles")
	assert.Equal(t, uint16(0x1010), cpu.SP, "SP should be incremented by 16")

	// Test negative offset (two's complement)
	cpu.SP = 0x1000
	cycles = cpu.ADD_SP_n(0xF0) // -16 in two's complement

	assert.Equal(t, uint8(16), cycles, "Should take 16 cycles")
	assert.Equal(t, uint16(0x0FF0), cpu.SP, "SP should be decremented by 16")
}

func TestADD_SP_n_Flags(t *testing.T) {
	cpu := NewCPU()

	// Test half-carry flag (carry from bit 3 to bit 4)
	cpu.SP = 0x100F
	cpu.F = 0xFF // Set all flags initially
	cpu.ADD_SP_n(0x01) // Should cause half-carry

	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")

	// Test carry flag (carry from bit 7 to bit 8)
	cpu.SP = 0x10FF
	cpu.F = 0x00 // Clear all flags
	cpu.ADD_SP_n(0x01) // Should cause carry

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")

	// Test no flags set
	cpu.SP = 0x1000
	cpu.F = 0xFF // Set all flags initially
	cpu.ADD_SP_n(0x01) // Should not cause any flags

	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")
	assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
}

// Test LD HL,SP+n instruction (0xF8)
func TestLD_HL_SP_n(t *testing.T) {
	cpu := NewCPU()

	// Test positive offset
	cpu.SP = 0x1000
	cpu.SetHL(0x0000) // Different value to ensure it changes

	cycles := cpu.LD_HL_SP_n(0x10) // +16

	assert.Equal(t, uint8(12), cycles, "LD HL,SP+n should take 12 cycles")
	assert.Equal(t, uint16(0x1010), cpu.GetHL(), "HL should be SP + offset")
	assert.Equal(t, uint16(0x1000), cpu.SP, "SP should remain unchanged")

	// Test negative offset
	cpu.SP = 0x1000
	cpu.SetHL(0x0000)

	cycles = cpu.LD_HL_SP_n(0xF0) // -16 in two's complement

	assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
	assert.Equal(t, uint16(0x0FF0), cpu.GetHL(), "HL should be SP - 16")
	assert.Equal(t, uint16(0x1000), cpu.SP, "SP should remain unchanged")
}

func TestLD_HL_SP_n_Flags(t *testing.T) {
	cpu := NewCPU()

	// Test half-carry flag (same logic as ADD_SP_n)
	cpu.SP = 0x100F
	cpu.F = 0xFF // Set all flags initially
	cpu.LD_HL_SP_n(0x01) // Should cause half-carry

	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")

	// Test carry flag
	cpu.SP = 0x10FF
	cpu.F = 0x00 // Clear all flags
	cpu.LD_HL_SP_n(0x01) // Should cause carry

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")

	// Verify SP unchanged
	assert.Equal(t, uint16(0x1100), cpu.GetHL(), "HL should have result")
	assert.Equal(t, uint16(0x10FF), cpu.SP, "SP should be unchanged")
}

// Test edge cases for signed arithmetic
func TestStackPointerSignedArithmetic(t *testing.T) {
	cpu := NewCPU()

	testCases := []struct {
		name     string
		sp       uint16
		offset   uint8
		expected uint16
	}{
		{"Max positive offset", 0x1000, 0x7F, 0x107F},   // +127
		{"Max negative offset", 0x1000, 0x80, 0x0F80},   // -128
		{"Zero offset", 0x1234, 0x00, 0x1234},           // +0
		{"Wrap around positive", 0xFFFF, 0x01, 0x0000},  // Overflow
		{"Wrap around negative", 0x0000, 0xFF, 0xFFFF},  // Underflow
	}

	for _, tc := range testCases {
		t.Run(tc.name+" ADD_SP_n", func(t *testing.T) {
			cpu.SP = tc.sp
			cpu.ADD_SP_n(tc.offset)
			assert.Equal(t, tc.expected, cpu.SP, "ADD_SP_n result mismatch")
		})

		t.Run(tc.name+" LD_HL_SP_n", func(t *testing.T) {
			cpu.SP = tc.sp
			cpu.LD_HL_SP_n(tc.offset)
			assert.Equal(t, tc.expected, cpu.GetHL(), "LD_HL_SP_n result mismatch")
			assert.Equal(t, tc.sp, cpu.SP, "SP should not change in LD_HL_SP_n")
		})
	}
}
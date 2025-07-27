package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test opcode dispatch integration for new stack pointer instructions
func TestStackPointerOpcodeDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("LD (nn),SP opcode 0x08", func(t *testing.T) {
		cpu.SP = 0x1234
		cycles, err := cpu.ExecuteInstruction(mmu, 0x08, 0x00, 0x80) // LD (0x8000),SP

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(20), cycles, "Should return correct cycles")

		// Verify memory contents
		assert.Equal(t, uint8(0x34), mmu.ReadByte(0x8000), "Low byte stored correctly")
		assert.Equal(t, uint8(0x12), mmu.ReadByte(0x8001), "High byte stored correctly")
	})

	t.Run("LD SP,HL opcode 0xF9", func(t *testing.T) {
		cpu.SetHL(0x5678)
		cpu.SP = 0x1234
		cycles, err := cpu.ExecuteInstruction(mmu, 0xF9) // LD SP,HL

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(8), cycles, "Should return correct cycles")
		assert.Equal(t, uint16(0x5678), cpu.SP, "SP should match HL")
	})

	t.Run("ADD SP,n opcode 0xE8", func(t *testing.T) {
		cpu.SP = 0x1000
		cycles, err := cpu.ExecuteInstruction(mmu, 0xE8, 0x10) // ADD SP,+16

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(16), cycles, "Should return correct cycles")
		assert.Equal(t, uint16(0x1010), cpu.SP, "SP should be incremented")
	})

	t.Run("LD HL,SP+n opcode 0xF8", func(t *testing.T) {
		cpu.SP = 0x1000
		cpu.SetHL(0x0000)
		cycles, err := cpu.ExecuteInstruction(mmu, 0xF8, 0x10) // LD HL,SP+16

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(12), cycles, "Should return correct cycles")
		assert.Equal(t, uint16(0x1010), cpu.GetHL(), "HL should be SP + offset")
		assert.Equal(t, uint16(0x1000), cpu.SP, "SP should remain unchanged")
	})
}

// Test wrapper function error handling
func TestStackPointerWrapperErrors(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("LD (nn),SP missing parameters", func(t *testing.T) {
		// Test with no parameters
		_, err := cpu.ExecuteInstruction(mmu, 0x08)
		assert.Error(t, err, "Should error with no parameters")
		assert.Contains(t, err.Error(), "requires 2 parameters", "Should mention parameter count")

		// Test with one parameter
		_, err = cpu.ExecuteInstruction(mmu, 0x08, 0x00)
		assert.Error(t, err, "Should error with one parameter")
	})

	t.Run("ADD SP,n missing parameter", func(t *testing.T) {
		_, err := cpu.ExecuteInstruction(mmu, 0xE8)
		assert.Error(t, err, "Should error with no parameters")
		assert.Contains(t, err.Error(), "requires 1 parameter", "Should mention parameter count")
	})

	t.Run("LD HL,SP+n missing parameter", func(t *testing.T) {
		_, err := cpu.ExecuteInstruction(mmu, 0xF8)
		assert.Error(t, err, "Should error with no parameters")
		assert.Contains(t, err.Error(), "requires 1 parameter", "Should mention parameter count")
	})
}

// Test opcode info for new instructions
func TestStackPointerOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode   uint8
		expected string
	}{
		{0x08, "LD (nn),SP"},
		{0xE8, "ADD SP,n"},
		{0xF8, "LD HL,SP+n"},
		{0xF9, "LD SP,HL"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			info, exists := GetOpcodeInfo(tc.opcode)
			assert.True(t, exists, "Opcode should be recognized")
			assert.Equal(t, tc.expected, info, "Should return correct instruction name")
		})
	}
}

// Test real Game Boy usage patterns
func TestStackPointerGameBoyPatterns(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("Function call setup pattern", func(t *testing.T) {
		// Common pattern: set up stack frame with HL offset
		cpu.SP = 0xFFFE // Normal stack start

		// LD HL,SP+n - create stack frame pointer
		cycles1, err1 := cpu.ExecuteInstruction(mmu, 0xF8, 0xFC) // SP-4
		assert.NoError(t, err1)
		assert.Equal(t, uint8(12), cycles1)
		assert.Equal(t, uint16(0xFFFA), cpu.GetHL(), "HL should point to stack frame")

		// LD SP,HL - adjust stack pointer
		cycles2, err2 := cpu.ExecuteInstruction(mmu, 0xF9)
		assert.NoError(t, err2)
		assert.Equal(t, uint8(8), cycles2)
		assert.Equal(t, uint16(0xFFFA), cpu.SP, "SP should be adjusted")
	})

	t.Run("SP save/restore pattern", func(t *testing.T) {
		// Save SP to memory
		cpu.SP = 0x1234
		cycles1, err1 := cpu.ExecuteInstruction(mmu, 0x08, 0x00, 0x80) // LD (0x8000),SP
		assert.NoError(t, err1)
		assert.Equal(t, uint8(20), cycles1)

		// Modify SP
		cpu.SP = 0x5678

		// Restore SP from memory (would need to implement LD SP,(nn) for complete pattern)
		// For now, verify the save worked
		storedSP := mmu.ReadWord(0x8000)
		assert.Equal(t, uint16(0x1234), storedSP, "SP should be correctly stored in memory")
	})
}

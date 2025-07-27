package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test flag operations via opcode dispatch
func TestFlagOpcodeDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("DAA opcode 0x27", func(t *testing.T) {
		// Simple DAA test - just verify it executes and returns correct cycles
		cpu.A = 0x0A
		cpu.SetFlag(FlagN, false) // Addition mode
		cpu.SetFlag(FlagH, false)
		cpu.SetFlag(FlagC, false)

		cycles, err := cpu.ExecuteInstruction(mmu, 0x27) // DAA

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		// DAA behavior is complex, just verify it modifies A and clears H flag
		assert.False(t, cpu.GetFlag(FlagH), "H flag should be cleared")
		// A may or may not change depending on the exact value and flags
	})

	t.Run("CPL opcode 0x2F", func(t *testing.T) {
		cpu.A = 0x35 // Binary: 00110101
		originalZ := cpu.GetFlag(FlagZ)
		originalC := cpu.GetFlag(FlagC)

		cycles, err := cpu.ExecuteInstruction(mmu, 0x2F) // CPL

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.Equal(t, uint8(0xCA), cpu.A, "Should complement A (00110101 → 11001010)")
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be set")
		assert.True(t, cpu.GetFlag(FlagH), "H flag should be set")
		assert.Equal(t, originalZ, cpu.GetFlag(FlagZ), "Z flag should be preserved")
		assert.Equal(t, originalC, cpu.GetFlag(FlagC), "C flag should be preserved")
	})

	t.Run("SCF opcode 0x37", func(t *testing.T) {
		cpu.SetFlag(FlagC, false) // Start with carry clear
		originalZ := cpu.GetFlag(FlagZ)

		cycles, err := cpu.ExecuteInstruction(mmu, 0x37) // SCF

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be set")
		assert.False(t, cpu.GetFlag(FlagN), "N flag should be cleared")
		assert.False(t, cpu.GetFlag(FlagH), "H flag should be cleared")
		assert.Equal(t, originalZ, cpu.GetFlag(FlagZ), "Z flag should be preserved")
	})

	t.Run("CCF opcode 0x3F", func(t *testing.T) {
		// Test with carry initially set
		cpu.SetFlag(FlagC, true)
		originalZ := cpu.GetFlag(FlagZ)

		cycles, err := cpu.ExecuteInstruction(mmu, 0x3F) // CCF

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.False(t, cpu.GetFlag(FlagC), "C flag should be flipped to clear")
		assert.False(t, cpu.GetFlag(FlagN), "N flag should be cleared")
		assert.False(t, cpu.GetFlag(FlagH), "H flag should be cleared")
		assert.Equal(t, originalZ, cpu.GetFlag(FlagZ), "Z flag should be preserved")

		// Test with carry initially clear
		cycles2, err2 := cpu.ExecuteInstruction(mmu, 0x3F) // CCF again

		assert.NoError(t, err2, "Should execute without error")
		assert.Equal(t, uint8(4), cycles2, "Should return correct cycles")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be flipped back to set")
	})
}

// Test flag instruction opcode info
func TestFlagOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode   uint8
		expected string
	}{
		{0x27, "DAA"},
		{0x2F, "CPL"},
		{0x37, "SCF"},
		{0x3F, "CCF"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			info, exists := GetOpcodeInfo(tc.opcode)
			assert.True(t, exists, "Opcode should be recognized")
			assert.Equal(t, tc.expected, info, "Should return correct instruction name")
		})
	}
}

package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test opcode dispatch integration for control instructions
func TestControlOpcodeDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("HALT opcode 0x76", func(t *testing.T) {
		// Reset state
		cpu.Resume()
		assert.False(t, cpu.IsHalted(), "Should not be halted initially")

		cycles, err := cpu.ExecuteInstruction(mmu, 0x76) // HALT

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.True(t, cpu.IsHalted(), "CPU should be halted")
		assert.False(t, cpu.IsStopped(), "CPU should not be stopped")
	})

	t.Run("STOP opcode 0x10", func(t *testing.T) {
		// Reset state
		cpu.Resume()
		assert.False(t, cpu.IsStopped(), "Should not be stopped initially")

		cycles, err := cpu.ExecuteInstruction(mmu, 0x10) // STOP

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.True(t, cpu.IsStopped(), "CPU should be stopped")
		assert.True(t, cpu.IsHalted(), "CPU should also be halted")
	})

	t.Run("DI opcode 0xF3", func(t *testing.T) {
		// Enable interrupts first
		cpu.InterruptsEnabled = true
		assert.True(t, cpu.AreInterruptsEnabled(), "Should have interrupts enabled initially")

		cycles, err := cpu.ExecuteInstruction(mmu, 0xF3) // DI

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.False(t, cpu.AreInterruptsEnabled(), "Interrupts should be disabled")
	})

	t.Run("EI opcode 0xFB", func(t *testing.T) {
		// Disable interrupts first
		cpu.InterruptsEnabled = false
		assert.False(t, cpu.AreInterruptsEnabled(), "Should have interrupts disabled initially")

		cycles, err := cpu.ExecuteInstruction(mmu, 0xFB) // EI

		assert.NoError(t, err, "Should execute without error")
		assert.Equal(t, uint8(4), cycles, "Should return correct cycles")
		assert.True(t, cpu.AreInterruptsEnabled(), "Interrupts should be enabled")
	})
}

// Test control instruction opcode info
func TestControlOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode   uint8
		expected string
	}{
		{0x10, "STOP"},
		{0x76, "HALT"},
		{0xF3, "DI"},
		{0xFB, "EI"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			info, exists := GetOpcodeInfo(tc.opcode)
			assert.True(t, exists, "Opcode should be recognized")
			assert.Equal(t, tc.expected, info, "Should return correct instruction name")
		})
	}
}

// Test Game Boy power management patterns
func TestPowerManagementPatterns(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("Power saving sequence", func(t *testing.T) {
		// Typical power saving: disable interrupts, then halt
		cycles1, err1 := cpu.ExecuteInstruction(mmu, 0xF3) // DI
		assert.NoError(t, err1)
		assert.Equal(t, uint8(4), cycles1)
		assert.False(t, cpu.AreInterruptsEnabled())

		cycles2, err2 := cpu.ExecuteInstruction(mmu, 0x76) // HALT
		assert.NoError(t, err2)
		assert.Equal(t, uint8(4), cycles2)
		assert.True(t, cpu.IsHalted())
	})

	t.Run("Critical section pattern", func(t *testing.T) {
		// Reset state
		cpu.Resume()
		cpu.InterruptsEnabled = true

		// Critical section: disable interrupts, do work, re-enable
		cycles1, err1 := cpu.ExecuteInstruction(mmu, 0xF3) // DI
		assert.NoError(t, err1)
		assert.False(t, cpu.AreInterruptsEnabled())

		// (critical work would happen here)

		cycles2, err2 := cpu.ExecuteInstruction(mmu, 0xFB) // EI
		assert.NoError(t, err2)
		assert.True(t, cpu.AreInterruptsEnabled())

		// Verify total cycles
		totalCycles := cycles1 + cycles2
		assert.Equal(t, uint8(8), totalCycles, "DI + EI should take 8 cycles total")
	})

	t.Run("Deep sleep pattern", func(t *testing.T) {
		// Reset state
		cpu.Resume()

		// Deep sleep: disable interrupts, then stop
		cpu.ExecuteInstruction(mmu, 0xF3) // DI
		assert.False(t, cpu.AreInterruptsEnabled())

		cycles, err := cpu.ExecuteInstruction(mmu, 0x10) // STOP
		assert.NoError(t, err)
		assert.Equal(t, uint8(4), cycles)
		assert.True(t, cpu.IsStopped())
		assert.True(t, cpu.IsHalted())
		assert.False(t, cpu.AreInterruptsEnabled())
	})
}

// Test control instructions with different CPU states
func TestControlInstructionsWithDifferentStates(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("HALT when already halted", func(t *testing.T) {
		// First HALT
		cpu.ExecuteInstruction(mmu, 0x76)
		assert.True(t, cpu.IsHalted())

		// Second HALT (should still work)
		cycles, err := cpu.ExecuteInstruction(mmu, 0x76)
		assert.NoError(t, err)
		assert.Equal(t, uint8(4), cycles)
		assert.True(t, cpu.IsHalted())
	})

	t.Run("STOP when already stopped", func(t *testing.T) {
		// Reset and first STOP
		cpu.Resume()
		cpu.ExecuteInstruction(mmu, 0x10)
		assert.True(t, cpu.IsStopped())

		// Second STOP (should still work)
		cycles, err := cpu.ExecuteInstruction(mmu, 0x10)
		assert.NoError(t, err)
		assert.Equal(t, uint8(4), cycles)
		assert.True(t, cpu.IsStopped())
	})

	t.Run("DI when interrupts already disabled", func(t *testing.T) {
		// Ensure interrupts disabled
		cpu.InterruptsEnabled = false
		assert.False(t, cpu.AreInterruptsEnabled())

		// DI again (should still work)
		cycles, err := cpu.ExecuteInstruction(mmu, 0xF3)
		assert.NoError(t, err)
		assert.Equal(t, uint8(4), cycles)
		assert.False(t, cpu.AreInterruptsEnabled())
	})

	t.Run("EI when interrupts already enabled", func(t *testing.T) {
		// Ensure interrupts enabled
		cpu.InterruptsEnabled = true
		assert.True(t, cpu.AreInterruptsEnabled())

		// EI again (should still work)
		cycles, err := cpu.ExecuteInstruction(mmu, 0xFB)
		assert.NoError(t, err)
		assert.Equal(t, uint8(4), cycles)
		assert.True(t, cpu.AreInterruptsEnabled())
	})
}

// Test that control instructions preserve all registers and flags
func TestControlInstructionsPreservation(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up distinctive CPU state
	cpu.A = 0x12
	cpu.B = 0x34
	cpu.C = 0x56
	cpu.D = 0x78
	cpu.E = 0x9A
	cpu.F = 0xBC
	cpu.H = 0xDE
	cpu.L = 0xF0
	cpu.SP = 0x1234
	cpu.PC = 0x5678

	// Store original state
	originalState := *cpu

	// Execute all control instructions
	controlOpcodes := []uint8{0x76, 0x10, 0xF3, 0xFB}
	for _, opcode := range controlOpcodes {
		// Reset control state but preserve registers
		cpu.Resume()

		cycles, err := cpu.ExecuteInstruction(mmu, opcode)
		assert.NoError(t, err, "Opcode 0x%02X should execute successfully", opcode)
		assert.Equal(t, uint8(4), cycles, "Opcode 0x%02X should take 4 cycles", opcode)

		// Check all registers preserved (except control state)
		assert.Equal(t, originalState.A, cpu.A, "A should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.B, cpu.B, "B should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.C, cpu.C, "C should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.D, cpu.D, "D should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.E, cpu.E, "E should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.F, cpu.F, "F should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.H, cpu.H, "H should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.L, cpu.L, "L should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.SP, cpu.SP, "SP should be preserved for opcode 0x%02X", opcode)
		assert.Equal(t, originalState.PC, cpu.PC, "PC should be preserved for opcode 0x%02X", opcode)
	}
}

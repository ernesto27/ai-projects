package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test HALT instruction (0x76)
func TestHALT(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Verify initial state
	assert.False(t, cpu.IsHalted(), "CPU should not be halted initially")

	// Execute HALT
	cycles := cpu.HALT(mmu)

	// Verify result
	assert.Equal(t, uint8(4), cycles, "HALT should take 4 cycles")
	assert.True(t, cpu.IsHalted(), "CPU should be halted after HALT")
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped by HALT")

	// Verify other state unchanged
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should not change")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should not change")
}

func TestHALTFlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags
	cpu.F = 0xFF
	originalFlags := cpu.F

	// Execute HALT
	cpu.HALT(mmu)

	// Verify flags unchanged
	assert.Equal(t, originalFlags, cpu.F, "HALT should not affect flags")
}

// Test STOP instruction (0x10)
func TestSTOP(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Verify initial state
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped initially")
	assert.False(t, cpu.IsHalted(), "CPU should not be halted initially")

	// Execute STOP
	cycles := cpu.STOP(mmu)

	// Verify result
	assert.Equal(t, uint8(4), cycles, "STOP should take 4 cycles")
	assert.True(t, cpu.IsStopped(), "CPU should be stopped after STOP")
	assert.True(t, cpu.IsHalted(), "CPU should also be halted after STOP")

	// Verify other state unchanged
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should not change")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should not change")
}

func TestSTOPFlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags
	cpu.F = 0xFF
	originalFlags := cpu.F

	// Execute STOP
	cpu.STOP(mmu)

	// Verify flags unchanged
	assert.Equal(t, originalFlags, cpu.F, "STOP should not affect flags")
}

// Test DI instruction (0xF3)
func TestDI(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Enable interrupts first
	cpu.InterruptsEnabled = true
	assert.True(t, cpu.AreInterruptsEnabled(), "Interrupts should be enabled initially")

	// Execute DI
	cycles := cpu.DI(mmu)

	// Verify result
	assert.Equal(t, uint8(4), cycles, "DI should take 4 cycles")
	assert.False(t, cpu.AreInterruptsEnabled(), "Interrupts should be disabled after DI")

	// Verify other state unchanged
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should not change")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should not change")
	assert.False(t, cpu.IsHalted(), "CPU should not be halted")
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped")
}

func TestDIFlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags
	cpu.F = 0xFF
	originalFlags := cpu.F

	// Execute DI
	cpu.DI(mmu)

	// Verify flags unchanged
	assert.Equal(t, originalFlags, cpu.F, "DI should not affect flags")
}

// Test EI instruction (0xFB)
func TestEI(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Verify initial state (interrupts disabled at boot)
	assert.False(t, cpu.AreInterruptsEnabled(), "Interrupts should be disabled initially")

	// Execute EI
	cycles := cpu.EI(mmu)

	// Verify result
	assert.Equal(t, uint8(4), cycles, "EI should take 4 cycles")
	assert.True(t, cpu.AreInterruptsEnabled(), "Interrupts should be enabled after EI")

	// Verify other state unchanged
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should not change")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should not change")
	assert.False(t, cpu.IsHalted(), "CPU should not be halted")
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped")
}

func TestEIFlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags
	cpu.F = 0xFF
	originalFlags := cpu.F

	// Execute EI
	cpu.EI(mmu)

	// Verify flags unchanged
	assert.Equal(t, originalFlags, cpu.F, "EI should not affect flags")
}

// Test DI/EI sequence
func TestDI_EI_Sequence(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test initial state
	assert.False(t, cpu.AreInterruptsEnabled(), "Should start with interrupts disabled")

	// Enable interrupts
	cpu.EI(mmu)
	assert.True(t, cpu.AreInterruptsEnabled(), "EI should enable interrupts")

	// Disable interrupts
	cpu.DI(mmu)
	assert.False(t, cpu.AreInterruptsEnabled(), "DI should disable interrupts")

	// Enable again
	cpu.EI(mmu)
	assert.True(t, cpu.AreInterruptsEnabled(), "EI should enable interrupts again")
}

// Test Resume function
func TestResume(t *testing.T) {
	cpu := NewCPU()

	// Put CPU in halt state
	cpu.Halted = true
	assert.True(t, cpu.IsHalted(), "CPU should be halted")

	// Resume
	cpu.Resume()
	assert.False(t, cpu.IsHalted(), "CPU should not be halted after resume")
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped after resume")

	// Put CPU in stop state
	cpu.Stopped = true
	cpu.Halted = true
	assert.True(t, cpu.IsStopped(), "CPU should be stopped")
	assert.True(t, cpu.IsHalted(), "CPU should be halted")

	// Resume
	cpu.Resume()
	assert.False(t, cpu.IsHalted(), "CPU should not be halted after resume")
	assert.False(t, cpu.IsStopped(), "CPU should not be stopped after resume")
}

// Test state query functions
func TestStateQueryFunctions(t *testing.T) {
	cpu := NewCPU()

	// Test initial state
	assert.False(t, cpu.IsHalted(), "IsHalted should return false initially")
	assert.False(t, cpu.IsStopped(), "IsStopped should return false initially")
	assert.False(t, cpu.AreInterruptsEnabled(), "AreInterruptsEnabled should return false initially")

	// Test after setting states
	cpu.Halted = true
	assert.True(t, cpu.IsHalted(), "IsHalted should return true when halted")

	cpu.Stopped = true
	assert.True(t, cpu.IsStopped(), "IsStopped should return true when stopped")

	cpu.InterruptsEnabled = true
	assert.True(t, cpu.AreInterruptsEnabled(), "AreInterruptsEnabled should return true when enabled")
}

// Test control instructions don't affect registers
func TestControlInstructionsRegisterPreservation(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set distinctive register values
	cpu.A = 0xAA
	cpu.B = 0xBB
	cpu.C = 0xCC
	cpu.D = 0xDD
	cpu.E = 0xEE
	cpu.H = 0x11
	cpu.L = 0x22
	cpu.SP = 0x1234
	cpu.PC = 0x5678

	// Store original values
	originalA := cpu.A
	originalB := cpu.B
	originalC := cpu.C
	originalD := cpu.D
	originalE := cpu.E
	originalH := cpu.H
	originalL := cpu.L
	originalSP := cpu.SP
	originalPC := cpu.PC

	// Execute all control instructions
	cpu.HALT(mmu)
	cpu.Resume() // Reset state
	cpu.STOP(mmu)
	cpu.Resume() // Reset state
	cpu.DI(mmu)
	cpu.EI(mmu)

	// Verify all registers preserved
	assert.Equal(t, originalA, cpu.A, "A register should be preserved")
	assert.Equal(t, originalB, cpu.B, "B register should be preserved")
	assert.Equal(t, originalC, cpu.C, "C register should be preserved")
	assert.Equal(t, originalD, cpu.D, "D register should be preserved")
	assert.Equal(t, originalE, cpu.E, "E register should be preserved")
	assert.Equal(t, originalH, cpu.H, "H register should be preserved")
	assert.Equal(t, originalL, cpu.L, "L register should be preserved")
	assert.Equal(t, originalSP, cpu.SP, "SP should be preserved")
	assert.Equal(t, originalPC, cpu.PC, "PC should be preserved")
}

// Test edge cases and combinations
func TestControlInstructionEdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	t.Run("Multiple HALT calls", func(t *testing.T) {
		cpu.HALT(mmu)
		assert.True(t, cpu.IsHalted(), "Should be halted after first HALT")

		cpu.HALT(mmu) // Second HALT
		assert.True(t, cpu.IsHalted(), "Should still be halted after second HALT")
	})

	t.Run("Multiple STOP calls", func(t *testing.T) {
		cpu.Resume() // Reset state
		cpu.STOP(mmu)
		assert.True(t, cpu.IsStopped(), "Should be stopped after first STOP")

		cpu.STOP(mmu) // Second STOP
		assert.True(t, cpu.IsStopped(), "Should still be stopped after second STOP")
	})

	t.Run("Multiple DI calls", func(t *testing.T) {
		cpu.InterruptsEnabled = true // Enable first
		cpu.DI(mmu)
		assert.False(t, cpu.AreInterruptsEnabled(), "Should be disabled after first DI")

		cpu.DI(mmu) // Second DI
		assert.False(t, cpu.AreInterruptsEnabled(), "Should still be disabled after second DI")
	})

	t.Run("Multiple EI calls", func(t *testing.T) {
		cpu.InterruptsEnabled = false // Disable first
		cpu.EI(mmu)
		assert.True(t, cpu.AreInterruptsEnabled(), "Should be enabled after first EI")

		cpu.EI(mmu) // Second EI
		assert.True(t, cpu.AreInterruptsEnabled(), "Should still be enabled after second EI")
	})
}

package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNOPInstruction tests the NOP instruction
func TestNOPInstruction(t *testing.T) {
	cpu := NewCPU()

	// Store initial state
	initialPC := cpu.PC
	initialA := cpu.A
	initialF := cpu.F
	initialSP := cpu.SP

	// Execute NOP instruction
	cycles := cpu.NOP()

	// NOP should take 4 cycles
	assert.Equal(t, uint8(4), cycles, "NOP should take 4 cycles")

	// NOP should not change any registers
	assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged after NOP")
	assert.Equal(t, initialA, cpu.A, "A should be unchanged after NOP")
	assert.Equal(t, initialF, cpu.F, "F should be unchanged after NOP")
	assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged after NOP")
}

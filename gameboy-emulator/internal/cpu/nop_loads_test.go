package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test register self-load (NOP-like) operations via opcode dispatch
func TestNOPLoadOpcodeDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	testCases := []struct {
		name     string
		opcode   uint8
		register *uint8
		expected string
	}{
		{"LD B,B", 0x40, &cpu.B, "LD B,B"},
		{"LD C,C", 0x49, &cpu.C, "LD C,C"},
		{"LD D,D", 0x52, &cpu.D, "LD D,D"},
		{"LD E,E", 0x5B, &cpu.E, "LD E,E"},
		{"LD H,H", 0x64, &cpu.H, "LD H,H"},
		{"LD L,L", 0x6D, &cpu.L, "LD L,L"},
		{"LD A,A", 0x7F, &cpu.A, "LD A,A"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set register to distinctive value
			*tc.register = 0x55
			originalValue := *tc.register

			// Store all other register values to verify they're unchanged
			originalA := cpu.A
			originalB := cpu.B
			originalC := cpu.C
			originalD := cpu.D
			originalE := cpu.E
			originalF := cpu.F
			originalH := cpu.H
			originalL := cpu.L
			originalSP := cpu.SP
			originalPC := cpu.PC

			cycles, err := cpu.ExecuteInstruction(mmu, tc.opcode)

			assert.NoError(t, err, "Should execute without error")
			assert.Equal(t, uint8(4), cycles, "Should return 4 cycles")
			assert.Equal(t, originalValue, *tc.register, "Register should remain unchanged")

			// Verify all other registers unchanged
			assert.Equal(t, originalA, cpu.A, "A register should be preserved")
			assert.Equal(t, originalB, cpu.B, "B register should be preserved")
			assert.Equal(t, originalC, cpu.C, "C register should be preserved")
			assert.Equal(t, originalD, cpu.D, "D register should be preserved")
			assert.Equal(t, originalE, cpu.E, "E register should be preserved")
			assert.Equal(t, originalF, cpu.F, "F register should be preserved")
			assert.Equal(t, originalH, cpu.H, "H register should be preserved")
			assert.Equal(t, originalL, cpu.L, "L register should be preserved")
			assert.Equal(t, originalSP, cpu.SP, "SP should be preserved")
			assert.Equal(t, originalPC, cpu.PC, "PC should be preserved")
		})
	}
}

// Test NOP load instruction opcode info
func TestNOPLoadOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode   uint8
		expected string
	}{
		{0x40, "LD B,B"},
		{0x49, "LD C,C"},
		{0x52, "LD D,D"},
		{0x5B, "LD E,E"},
		{0x64, "LD H,H"},
		{0x6D, "LD L,L"},
		{0x7F, "LD A,A"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			info, exists := GetOpcodeInfo(tc.opcode)
			assert.True(t, exists, "Opcode should be recognized")
			assert.Equal(t, tc.expected, info, "Should return correct instruction name")
		})
	}
}

// Test that NOP loads don't affect flags
func TestNOPLoadsFlagPreservation(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags to distinctive pattern
	cpu.F = 0xF0

	opcodes := []uint8{0x40, 0x49, 0x52, 0x5B, 0x64, 0x6D, 0x7F}

	for _, opcode := range opcodes {
		originalFlags := cpu.F
		cycles, err := cpu.ExecuteInstruction(mmu, opcode)

		assert.NoError(t, err, "Opcode 0x%02X should execute without error", opcode)
		assert.Equal(t, uint8(4), cycles, "Opcode 0x%02X should take 4 cycles", opcode)
		assert.Equal(t, originalFlags, cpu.F, "Opcode 0x%02X should not affect flags", opcode)
	}
}

// Test NOP loads with edge case values
func TestNOPLoadsEdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	edgeValues := []uint8{0x00, 0xFF, 0x80, 0x7F}

	for _, value := range edgeValues {
		t.Run("Value_0x"+fmt.Sprintf("%02X", value), func(t *testing.T) {
			// Test all registers with this edge value
			cpu.A = value
			cpu.B = value
			cpu.C = value
			cpu.D = value
			cpu.E = value
			cpu.H = value
			cpu.L = value

			opcodes := []uint8{0x40, 0x49, 0x52, 0x5B, 0x64, 0x6D, 0x7F}
			for _, opcode := range opcodes {
				cycles, err := cpu.ExecuteInstruction(mmu, opcode)
				assert.NoError(t, err, "Should handle edge value 0x%02X", value)
				assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			}

			// Verify all registers still have the edge value
			assert.Equal(t, value, cpu.A, "A should retain edge value")
			assert.Equal(t, value, cpu.B, "B should retain edge value")
			assert.Equal(t, value, cpu.C, "C should retain edge value")
			assert.Equal(t, value, cpu.D, "D should retain edge value")
			assert.Equal(t, value, cpu.E, "E should retain edge value")
			assert.Equal(t, value, cpu.H, "H should retain edge value")
			assert.Equal(t, value, cpu.L, "L should retain edge value")
		})
	}
}

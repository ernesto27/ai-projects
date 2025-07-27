package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for Immediate Value Wrapper Functions ===
// These test wrapper functions that need to extract parameters from params[]

// === Test Load Immediate Wrapper Functions ===

func TestWrapLD_A_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test basic functionality
	cycles, err := wrapLD_A_n(cpu, mmu, 0x99)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD A,n should return 8 cycles")
	assert.Equal(t, uint8(0x99), cpu.A, "A should contain the immediate value")

	// Test error handling - no parameters
	_, err = wrapLD_A_n(cpu, mmu)
	assert.Error(t, err, "Should return error when no parameters provided")
	assert.Contains(t, err.Error(), "requires 1 parameter, got 0")
}

func TestWrapLD_B_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test basic functionality
	cycles, err := wrapLD_B_n(cpu, mmu, 0x55)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x55), cpu.B, "B should contain the immediate value")

	// Test error handling
	_, err = wrapLD_B_n(cpu, mmu)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LD B,n requires 1 parameter")
}

func TestWrapLD_C_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cycles, err := wrapLD_C_n(cpu, mmu, 0xAA)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0xAA), cpu.C)

	// Test error handling
	_, err = wrapLD_C_n(cpu, mmu)
	assert.Error(t, err)
}

func TestWrapLD_D_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cycles, err := wrapLD_D_n(cpu, mmu, 0x33)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x33), cpu.D)

	// Test error handling
	_, err = wrapLD_D_n(cpu, mmu)
	assert.Error(t, err)
}

func TestWrapLD_E_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cycles, err := wrapLD_E_n(cpu, mmu, 0x77)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x77), cpu.E)

	// Test error handling
	_, err = wrapLD_E_n(cpu, mmu)
	assert.Error(t, err)
}

func TestWrapLD_H_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cycles, err := wrapLD_H_n(cpu, mmu, 0x11)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x11), cpu.H)

	// Test error handling
	_, err = wrapLD_H_n(cpu, mmu)
	assert.Error(t, err)
}

func TestWrapLD_L_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cycles, err := wrapLD_L_n(cpu, mmu, 0xEE)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0xEE), cpu.L)

	// Test error handling
	_, err = wrapLD_L_n(cpu, mmu)
	assert.Error(t, err)
}

// === Test Arithmetic Immediate Wrapper Functions ===

func TestWrapADD_A_n(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set A to a known value
	cpu.A = 0x10

	// Test basic functionality
	cycles, err := wrapADD_A_n(cpu, mmu, 0x20)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "ADD A,n should return 8 cycles")
	assert.Equal(t, uint8(0x30), cpu.A, "A should contain the sum")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear for ADD")

	// Test error handling
	_, err = wrapADD_A_n(cpu, mmu)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ADD A,n requires 1 parameter")
}

// === Test Parameter Handling ===

func TestImmediateParameterHandling(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test that extra parameters are ignored (should only use first one)
	cycles, err := wrapLD_A_n(cpu, mmu, 0x42, 0x99, 0x11) // Extra params ignored

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x42), cpu.A, "Should use first parameter only")

	// Test all possible immediate values (0x00 to 0xFF)
	for i := 0; i <= 255; i++ {
		cpu.A = 0x00 // Reset A
		cycles, err := wrapLD_A_n(cpu, mmu, uint8(i))

		assert.NoError(t, err)
		assert.Equal(t, uint8(8), cycles)
		assert.Equal(t, uint8(i), cpu.A, "Should handle all 8-bit values")
	}
}

// === Comparison Tests: Wrapper vs Original ===

func TestImmediateWrappersVsOriginals(t *testing.T) {
	tests := []struct {
		name     string
		value    uint8
		wrapper  func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		original func(*CPU, uint8) uint8
		checkReg func(*CPU) uint8
	}{
		{
			name:     "LD A,n comparison",
			value:    0x42,
			wrapper:  wrapLD_A_n,
			original: (*CPU).LD_A_n,
			checkReg: func(cpu *CPU) uint8 { return cpu.A },
		},
		{
			name:     "LD B,n comparison",
			value:    0x55,
			wrapper:  wrapLD_B_n,
			original: (*CPU).LD_B_n,
			checkReg: func(cpu *CPU) uint8 { return cpu.B },
		},
		{
			name:     "LD C,n comparison",
			value:    0x99,
			wrapper:  wrapLD_C_n,
			original: (*CPU).LD_C_n,
			checkReg: func(cpu *CPU) uint8 { return cpu.C },
		},
		{
			name:     "ADD A,n comparison",
			value:    0x30,
			wrapper:  wrapADD_A_n,
			original: (*CPU).ADD_A_n,
			checkReg: func(cpu *CPU) uint8 { return cpu.A },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create two identical CPUs
			cpu1 := NewCPU()
			cpu2 := NewCPU()
			mmu := createTestMMU()

			// For ADD operations, set initial A value
			if tt.name == "ADD A,n comparison" {
				cpu1.A = 0x10
				cpu2.A = 0x10
			}

			// Call original on cpu1, wrapper on cpu2
			originalCycles := tt.original(cpu1, tt.value)
			wrapperCycles, err := tt.wrapper(cpu2, mmu, tt.value)

			// They should behave identically
			assert.NoError(t, err)
			assert.Equal(t, originalCycles, wrapperCycles, "cycles should match")

			// Check the target register
			assert.Equal(t, tt.checkReg(cpu1), tt.checkReg(cpu2), "register values should match")

			// All other registers should be identical
			assert.Equal(t, cpu1.F, cpu2.F, "F (flags) register should match")
		})
	}
}

// === Edge Cases Tests ===

func TestImmediateEdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test with boundary values
	testCases := []struct {
		name     string
		value    uint8
		expected uint8
	}{
		{"Min value", 0x00, 0x00},
		{"Max value", 0xFF, 0xFF},
		{"Mid value", 0x80, 0x80},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.A = 0x00 // Reset A

			cycles, err := wrapLD_A_n(cpu, mmu, tc.value)

			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tc.expected, cpu.A)
		})
	}

	// Test ADD_A_n with overflow
	cpu.A = 0xFF
	cycles, err := wrapADD_A_n(cpu, mmu, 0x01)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x00), cpu.A, "Should wrap around on overflow")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")
}

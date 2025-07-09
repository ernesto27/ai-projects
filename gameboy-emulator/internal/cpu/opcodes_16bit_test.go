package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for 16-bit Immediate Value Wrapper Functions ===
// These test wrapper functions that need to extract 2 parameters from params[]

// === Test 16-bit Load Immediate Wrapper Functions ===

func TestWrapLD_BC_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic functionality
	cycles, err := wrapLD_BC_nn(cpu, mmu, 0x34, 0x12) // low=0x34, high=0x12

	assert.NoError(t, err)
	assert.Equal(t, uint8(12), cycles, "LD BC,nn should return 12 cycles")
	assert.Equal(t, uint8(0x34), cpu.C, "C should contain low byte")
	assert.Equal(t, uint8(0x12), cpu.B, "B should contain high byte")
	assert.Equal(t, uint16(0x1234), cpu.GetBC(), "BC should contain 0x1234")

	// Test error handling - no parameters
	_, err = wrapLD_BC_nn(cpu, mmu)
	assert.Error(t, err, "Should return error when no parameters provided")
	assert.Contains(t, err.Error(), "requires 2 parameters, got 0")

	// Test error handling - only 1 parameter
	_, err = wrapLD_BC_nn(cpu, mmu, 0x34)
	assert.Error(t, err, "Should return error when only 1 parameter provided")
	assert.Contains(t, err.Error(), "requires 2 parameters, got 1")
}

func TestWrapLD_DE_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic functionality
	cycles, err := wrapLD_DE_nn(cpu, mmu, 0xAD, 0xDE) // low=0xAD, high=0xDE

	assert.NoError(t, err)
	assert.Equal(t, uint8(12), cycles, "LD DE,nn should return 12 cycles")
	assert.Equal(t, uint8(0xAD), cpu.E, "E should contain low byte")
	assert.Equal(t, uint8(0xDE), cpu.D, "D should contain high byte")
	assert.Equal(t, uint16(0xDEAD), cpu.GetDE(), "DE should contain 0xDEAD")

	// Test error handling
	_, err = wrapLD_DE_nn(cpu, mmu)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LD DE,nn requires 2 parameters")

	_, err = wrapLD_DE_nn(cpu, mmu, 0xAD)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires 2 parameters, got 1")
}

func TestWrapLD_HL_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic functionality
	cycles, err := wrapLD_HL_nn(cpu, mmu, 0xEF, 0xBE) // low=0xEF, high=0xBE

	assert.NoError(t, err)
	assert.Equal(t, uint8(12), cycles, "LD HL,nn should return 12 cycles")
	assert.Equal(t, uint8(0xEF), cpu.L, "L should contain low byte")
	assert.Equal(t, uint8(0xBE), cpu.H, "H should contain high byte")
	assert.Equal(t, uint16(0xBEEF), cpu.GetHL(), "HL should contain 0xBEEF")

	// Test error handling
	_, err = wrapLD_HL_nn(cpu, mmu)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LD HL,nn requires 2 parameters")

	_, err = wrapLD_HL_nn(cpu, mmu, 0xEF)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires 2 parameters, got 1")
}

func TestWrapLD_SP_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic functionality
	cycles, err := wrapLD_SP_nn(cpu, mmu, 0xFE, 0xFF) // low=0xFE, high=0xFF

	assert.NoError(t, err)
	assert.Equal(t, uint8(12), cycles, "LD SP,nn should return 12 cycles")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should contain 0xFFFE")

	// Test error handling
	_, err = wrapLD_SP_nn(cpu, mmu)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LD SP,nn requires 2 parameters")

	_, err = wrapLD_SP_nn(cpu, mmu, 0xFE)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires 2 parameters, got 1")
}

// === Test 16-bit Parameter Handling Edge Cases ===

func TestSixteenBitParameterHandling(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	tests := []struct {
		name      string
		wrapper   func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		params    []uint8
		expected  uint16
		getResult func(*CPU) uint16
	}{
		{
			name:      "BC with 0x0000",
			wrapper:   wrapLD_BC_nn,
			params:    []uint8{0x00, 0x00},
			expected:  0x0000,
			getResult: (*CPU).GetBC,
		},
		{
			name:      "BC with 0xFFFF",
			wrapper:   wrapLD_BC_nn,
			params:    []uint8{0xFF, 0xFF},
			expected:  0xFFFF,
			getResult: (*CPU).GetBC,
		},
		{
			name:      "DE with 0x8000",
			wrapper:   wrapLD_DE_nn,
			params:    []uint8{0x00, 0x80},
			expected:  0x8000,
			getResult: (*CPU).GetDE,
		},
		{
			name:      "HL with 0x007F",
			wrapper:   wrapLD_HL_nn,
			params:    []uint8{0x7F, 0x00},
			expected:  0x007F,
			getResult: (*CPU).GetHL,
		},
		{
			name:      "SP with 0xFFFE (typical stack)",
			wrapper:   wrapLD_SP_nn,
			params:    []uint8{0xFE, 0xFF},
			expected:  0xFFFE,
			getResult: func(cpu *CPU) uint16 { return cpu.SP },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset CPU state
			cpu.Reset()

			// Execute wrapper
			cycles, err := tt.wrapper(cpu, mmu, tt.params...)

			// Verify results
			assert.NoError(t, err)
			assert.Equal(t, uint8(12), cycles, "All 16-bit loads should take 12 cycles")
			assert.Equal(t, tt.expected, tt.getResult(cpu), "Should load correct 16-bit value")
		})
	}
}

// === Comparison Test: 16-bit Wrappers vs Original Functions ===

func TestSixteenBitWrappersVsOriginals(t *testing.T) {
	tests := []struct {
		name        string
		wrapper     func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		original    func(*CPU, uint8, uint8) uint8
		lowByte     uint8
		highByte    uint8
		checkResult func(*CPU, *CPU) bool
	}{
		{
			name:     "LD BC,nn comparison",
			wrapper:  wrapLD_BC_nn,
			original: (*CPU).LD_BC_nn,
			lowByte:  0x34,
			highByte: 0x12,
			checkResult: func(cpu1, cpu2 *CPU) bool {
				return cpu1.B == cpu2.B && cpu1.C == cpu2.C
			},
		},
		{
			name:     "LD DE,nn comparison",
			wrapper:  wrapLD_DE_nn,
			original: (*CPU).LD_DE_nn,
			lowByte:  0xAD,
			highByte: 0xDE,
			checkResult: func(cpu1, cpu2 *CPU) bool {
				return cpu1.D == cpu2.D && cpu1.E == cpu2.E
			},
		},
		{
			name:     "LD HL,nn comparison",
			wrapper:  wrapLD_HL_nn,
			original: (*CPU).LD_HL_nn,
			lowByte:  0xEF,
			highByte: 0xBE,
			checkResult: func(cpu1, cpu2 *CPU) bool {
				return cpu1.H == cpu2.H && cpu1.L == cpu2.L
			},
		},
		{
			name:     "LD SP,nn comparison",
			wrapper:  wrapLD_SP_nn,
			original: (*CPU).LD_SP_nn,
			lowByte:  0xFE,
			highByte: 0xFF,
			checkResult: func(cpu1, cpu2 *CPU) bool {
				return cpu1.SP == cpu2.SP
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create two identical CPUs
			cpu1 := NewCPU()
			cpu2 := NewCPU()
			mmu := memory.NewMMU()

			// Call original on cpu1, wrapper on cpu2
			originalCycles := tt.original(cpu1, tt.lowByte, tt.highByte)
			wrapperCycles, err := tt.wrapper(cpu2, mmu, tt.lowByte, tt.highByte)

			// They should behave identically
			assert.NoError(t, err)
			assert.Equal(t, originalCycles, wrapperCycles, "cycles should match")
			assert.True(t, tt.checkResult(cpu1, cpu2), "register values should match")
		})
	}
}

// === Test 16-bit Endianness Understanding ===

func TestSixteenBitEndianness(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test that we understand Game Boy's little-endian format
	// In little-endian: the first byte is the LOW byte, second byte is HIGH byte
	// So for 0x1234: low=0x34, high=0x12

	t.Run("Endianness verification", func(t *testing.T) {
		// Load 0x1234 (4660 in decimal)
		cycles, err := wrapLD_BC_nn(cpu, mmu, 0x34, 0x12) // low=0x34, high=0x12

		assert.NoError(t, err)
		assert.Equal(t, uint8(12), cycles)
		assert.Equal(t, uint8(0x34), cpu.C, "C (low byte) should be 0x34")
		assert.Equal(t, uint8(0x12), cpu.B, "B (high byte) should be 0x12")
		assert.Equal(t, uint16(0x1234), cpu.GetBC(), "BC should form 0x1234")
	})

	t.Run("Multiple endianness examples", func(t *testing.T) {
		examples := []struct {
			name      string
			low, high uint8
			expected  uint16
		}{
			{"0x0000", 0x00, 0x00, 0x0000},
			{"0x00FF", 0xFF, 0x00, 0x00FF},
			{"0xFF00", 0x00, 0xFF, 0xFF00},
			{"0xFFFF", 0xFF, 0xFF, 0xFFFF},
			{"0x1234", 0x34, 0x12, 0x1234},
			{"0xABCD", 0xCD, 0xAB, 0xABCD},
		}

		for _, ex := range examples {
			t.Run(ex.name, func(t *testing.T) {
				cpu.Reset()
				_, err := wrapLD_HL_nn(cpu, mmu, ex.low, ex.high)
				assert.NoError(t, err)
				assert.Equal(t, ex.expected, cpu.GetHL(), "HL should contain expected value")
			})
		}
	})
}

// === Test Error Handling Edge Cases ===

func TestSixteenBitErrorHandling(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	wrappers := []struct {
		name    string
		wrapper func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
	}{
		{"LD BC,nn", wrapLD_BC_nn},
		{"LD DE,nn", wrapLD_DE_nn},
		{"LD HL,nn", wrapLD_HL_nn},
		{"LD SP,nn", wrapLD_SP_nn},
	}

	for _, w := range wrappers {
		t.Run(w.name+" error cases", func(t *testing.T) {
			// Test no parameters
			_, err := w.wrapper(cpu, mmu)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requires 2 parameters, got 0")

			// Test 1 parameter
			_, err = w.wrapper(cpu, mmu, 0x34)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "requires 2 parameters, got 1")

			// Test 3 parameters (should still work, extras ignored)
			_, err = w.wrapper(cpu, mmu, 0x34, 0x12, 0x99)
			assert.NoError(t, err, "Extra parameters should be ignored")
		})
	}
}

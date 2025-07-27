package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestXORInstructionDispatch tests XOR instructions through the opcode dispatch system
func TestXORInstructionDispatch(t *testing.T) {
	tests := []struct {
		name      string
		opcode    uint8
		setupA    uint8
		setupReg  func(*CPU)
		params    []uint8
		expectedA uint8
		expectedZ bool
	}{
		{
			name:   "XOR A,B instruction (0xA8)",
			opcode: 0xA8,
			setupA: 0b11110000,
			setupReg: func(cpu *CPU) {
				cpu.B = 0b00001111
			},
			params:    []uint8{},
			expectedA: 0b11111111,
			expectedZ: false,
		},
		{
			name:   "XOR A,C instruction (0xA9)",
			opcode: 0xA9,
			setupA: 0x42,
			setupReg: func(cpu *CPU) {
				cpu.C = 0x42
			},
			params:    []uint8{},
			expectedA: 0x00, // Same values = zero
			expectedZ: true,
		},
		{
			name:   "XOR A,A instruction (0xAF) - Zero test",
			opcode: 0xAF,
			setupA: 0x55,
			setupReg: func(cpu *CPU) {
				// A is already set, XOR A,A always results in 0
			},
			params:    []uint8{},
			expectedA: 0x00,
			expectedZ: true,
		},
		{
			name:   "XOR A,(HL) instruction (0xAE)",
			opcode: 0xAE,
			setupA: 0b10101010,
			setupReg: func(cpu *CPU) {
				cpu.SetHL(0x8000)
			},
			params:    []uint8{},
			expectedA: 0b01010101, // Inverted bits
			expectedZ: false,
		},
		{
			name:   "XOR A,n instruction (0xEE)",
			opcode: 0xEE,
			setupA: 0x00,
			setupReg: func(cpu *CPU) {
				// No register setup needed for immediate
			},
			params:    []uint8{0xFF},
			expectedA: 0xFF, // 0x00 XOR 0xFF = 0xFF
			expectedZ: false,
		},
		{
			name:   "XOR A,n instruction (0xEE) - Zero result",
			opcode: 0xEE,
			setupA: 0xAA,
			setupReg: func(cpu *CPU) {
				// No register setup needed for immediate
			},
			params:    []uint8{0xAA},
			expectedA: 0x00, // 0xAA XOR 0xAA = 0x00
			expectedZ: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()

			// Setup
			cpu.A = tt.setupA
			tt.setupReg(cpu)

			// For memory operations, setup memory
			if tt.opcode == 0xAE {
				mmu.WriteByte(0x8000, 0b11111111) // Memory value for XOR
			}

			// Execute instruction
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			// Verify results
			assert.NoError(t, err, "Instruction should execute without error")
			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be set correctly")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
			assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")

			// Verify cycle count
			expectedCycles := uint8(4)
			if tt.opcode == 0xAE || tt.opcode == 0xEE {
				expectedCycles = 8 // Memory and immediate operations take 8 cycles
			}
			assert.Equal(t, expectedCycles, cycles, "Should return correct cycle count")
		})
	}
}

// TestXORInstructionBitPatterns tests various bit patterns with XOR through dispatch
func TestXORInstructionBitPatterns(t *testing.T) {
	t.Run("XOR instruction bit pattern tests", func(t *testing.T) {
		bitTests := []struct {
			name     string
			opcode   uint8
			valueA   uint8
			valueB   uint8
			expected uint8
			setupReg func(*CPU, uint8)
		}{
			{
				name:     "XOR A,B: Toggle alternating bits",
				opcode:   0xA8,
				valueA:   0b10101010,
				valueB:   0b01010101,
				expected: 0b11111111,
				setupReg: func(cpu *CPU, val uint8) { cpu.B = val },
			},
			{
				name:     "XOR A,B: Clear with identical values",
				opcode:   0xA8,
				valueA:   0xFF,
				valueB:   0xFF,
				expected: 0x00,
				setupReg: func(cpu *CPU, val uint8) { cpu.B = val },
			},
			{
				name:     "XOR A,C: Encryption pattern",
				opcode:   0xA9,
				valueA:   0x42,
				valueB:   0xAA,
				expected: 0xE8,
				setupReg: func(cpu *CPU, val uint8) { cpu.C = val },
			},
			{
				name:     "XOR A,D: Toggle high bit only",
				opcode:   0xAA,
				valueA:   0x00,
				valueB:   0x80,
				expected: 0x80,
				setupReg: func(cpu *CPU, val uint8) { cpu.D = val },
			},
			{
				name:     "XOR A,E: Toggle low bit only",
				opcode:   0xAB,
				valueA:   0x00,
				valueB:   0x01,
				expected: 0x01,
				setupReg: func(cpu *CPU, val uint8) { cpu.E = val },
			},
			{
				name:     "XOR A,H: Nibble combination",
				opcode:   0xAC,
				valueA:   0x0F,
				valueB:   0xF0,
				expected: 0xFF,
				setupReg: func(cpu *CPU, val uint8) { cpu.H = val },
			},
			{
				name:     "XOR A,L: No change with zero",
				opcode:   0xAD,
				valueA:   0x55,
				valueB:   0x00,
				expected: 0x55,
				setupReg: func(cpu *CPU, val uint8) { cpu.L = val },
			},
		}

		for _, tt := range bitTests {
			t.Run(tt.name, func(t *testing.T) {
				cpu := NewCPU()
				mmu := createTestMMU()

				cpu.A = tt.valueA
				tt.setupReg(cpu, tt.valueB)

				cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode)

				assert.NoError(t, err, "Instruction should execute without error")
				assert.Equal(t, tt.expected, cpu.A, "XOR result should match expected pattern")
				assert.Equal(t, uint8(4), cycles, "Register XOR operations should take 4 cycles")
			})
		}
	})
}

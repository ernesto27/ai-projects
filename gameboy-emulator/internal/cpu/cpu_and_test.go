package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAND_A_Register tests all register-to-register AND operations
func TestAND_A_Register(t *testing.T) {
	tests := []struct {
		name        string
		setupA      uint8
		setupOther  uint8
		setOtherReg func(*CPU, uint8) // Function to set the other register
		instruction func(*CPU) uint8
		expectedA   uint8
		expectedZ   bool
		expectedN   bool
		expectedH   bool
		expectedC   bool
	}{
		// AND_A_A tests
		{
			name:        "AND_A_A - non-zero value",
			setupA:      0x42,
			setupOther:  0x42,                         // Not used for AND_A_A
			setOtherReg: func(cpu *CPU, val uint8) {}, // No-op
			instruction: (*CPU).AND_A_A,
			expectedA:   0x42, // 0x42 & 0x42 = 0x42
			expectedZ:   false,
			expectedN:   false, // Always reset for AND
			expectedH:   true,  // Always set for AND
			expectedC:   false, // Always reset for AND
		},
		{
			name:        "AND_A_A - zero value",
			setupA:      0x00,
			setupOther:  0x00,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).AND_A_A,
			expectedA:   0x00,
			expectedZ:   true, // Result is zero
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_A - maximum value",
			setupA:      0xFF,
			setupOther:  0xFF,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).AND_A_A,
			expectedA:   0xFF, // 0xFF & 0xFF = 0xFF
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_B tests
		{
			name:        "AND_A_B - normal masking",
			setupA:      0xF0, // 11110000
			setupOther:  0x0F, // 00001111
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).AND_A_B,
			expectedA:   0x00, // 11110000 & 00001111 = 00000000
			expectedZ:   true, // Result is zero
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_B - partial overlap",
			setupA:      0xAA, // 10101010
			setupOther:  0x55, // 01010101
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).AND_A_B,
			expectedA:   0x00, // 10101010 & 01010101 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_B - identical values",
			setupA:      0x33, // 00110011
			setupOther:  0x33, // 00110011
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).AND_A_B,
			expectedA:   0x33, // 00110011 & 00110011 = 00110011
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_B - mask lower nibble",
			setupA:      0x87, // 10000111
			setupOther:  0x0F, // 00001111
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).AND_A_B,
			expectedA:   0x07, // 10000111 & 00001111 = 00000111
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_C tests
		{
			name:        "AND_A_C - mask upper nibble",
			setupA:      0x87, // 10000111
			setupOther:  0xF0, // 11110000
			setOtherReg: func(cpu *CPU, val uint8) { cpu.C = val },
			instruction: (*CPU).AND_A_C,
			expectedA:   0x80, // 10000111 & 11110000 = 10000000
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_C - clear all bits",
			setupA:      0xFF,
			setupOther:  0x00,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.C = val },
			instruction: (*CPU).AND_A_C,
			expectedA:   0x00, // 11111111 & 00000000 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_D tests
		{
			name:        "AND_A_D - bit pattern",
			setupA:      0xCC, // 11001100
			setupOther:  0xAA, // 10101010
			setOtherReg: func(cpu *CPU, val uint8) { cpu.D = val },
			instruction: (*CPU).AND_A_D,
			expectedA:   0x88, // 11001100 & 10101010 = 10001000
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_E tests
		{
			name:        "AND_A_E - single bit isolation",
			setupA:      0x81, // 10000001
			setupOther:  0x01, // 00000001
			setOtherReg: func(cpu *CPU, val uint8) { cpu.E = val },
			instruction: (*CPU).AND_A_E,
			expectedA:   0x01, // 10000001 & 00000001 = 00000001
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_H tests
		{
			name:        "AND_A_H - alternating bits",
			setupA:      0x55, // 01010101
			setupOther:  0x33, // 00110011
			setOtherReg: func(cpu *CPU, val uint8) { cpu.H = val },
			instruction: (*CPU).AND_A_H,
			expectedA:   0x11, // 01010101 & 00110011 = 00010001
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},

		// AND_A_L tests
		{
			name:        "AND_A_L - complex pattern",
			setupA:      0x7E, // 01111110
			setupOther:  0x42, // 01000010
			setOtherReg: func(cpu *CPU, val uint8) { cpu.L = val },
			instruction: (*CPU).AND_A_L,
			expectedA:   0x42, // 01111110 & 01000010 = 01000010
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_L - zero result from different inputs",
			setupA:      0x0F, // 00001111
			setupOther:  0xF0, // 11110000
			setOtherReg: func(cpu *CPU, val uint8) { cpu.L = val },
			instruction: (*CPU).AND_A_L,
			expectedA:   0x00, // 00001111 & 11110000 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			tt.setOtherReg(cpu, tt.setupOther)

			cycles := tt.instruction(cpu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Instruction cycles")
		})
	}
}

// TestAND_A_HL tests the AND A,(HL) memory operation
func TestAND_A_HL(t *testing.T) {
	tests := []struct {
		name        string
		setupA      uint8
		setupH      uint8
		setupL      uint8
		memoryValue uint8
		expectedA   uint8
		expectedZ   bool
		expectedN   bool
		expectedH   bool
		expectedC   bool
	}{
		{
			name:        "AND_A_HL - normal masking",
			setupA:      0xF0, // 11110000
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x0F, // 00001111
			expectedA:   0x00, // 11110000 & 00001111 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_HL - preserve bits",
			setupA:      0xFF, // 11111111
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x42, // 01000010
			expectedA:   0x42, // 11111111 & 01000010 = 01000010
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_HL - isolate single bit",
			setupA:      0x87, // 10000111
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x80, // 10000000
			expectedA:   0x80, // 10000111 & 10000000 = 10000000
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_HL - zero from memory",
			setupA:      0x42,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x00,
			expectedA:   0x00, // 01000010 & 00000000 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_HL - alternating pattern",
			setupA:      0xAA, // 10101010
			setupH:      0x90,
			setupL:      0x00,
			memoryValue: 0x55, // 01010101
			expectedA:   0x00, // 10101010 & 01010101 = 00000000
			expectedZ:   true,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "AND_A_HL - nibble masking",
			setupA:      0x3C, // 00111100
			setupH:      0x90,
			setupL:      0x00,
			memoryValue: 0x0F, // 00001111
			expectedA:   0x0C, // 00111100 & 00001111 = 00001100
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()

			cpu.A = tt.setupA
			cpu.H = tt.setupH
			cpu.L = tt.setupL

			address := uint16(tt.setupH)<<8 | uint16(tt.setupL)
			mmu.WriteByte(address, tt.memoryValue)

			cycles := cpu.AND_A_HL(mmu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Instruction cycles")
		})
	}
}

// TestAND_A_n tests the AND A,n immediate value operation
func TestAND_A_n(t *testing.T) {
	tests := []struct {
		name         string
		setupA       uint8
		immediateVal uint8
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
	}{
		{
			name:         "AND_A_n - mask lower nibble",
			setupA:       0x87, // 10000111
			immediateVal: 0x0F, // 00001111
			expectedA:    0x07, // 10000111 & 00001111 = 00000111
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - mask upper nibble",
			setupA:       0x87, // 10000111
			immediateVal: 0xF0, // 11110000
			expectedA:    0x80, // 10000111 & 11110000 = 10000000
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - clear all bits",
			setupA:       0xFF,
			immediateVal: 0x00,
			expectedA:    0x00, // 11111111 & 00000000 = 00000000
			expectedZ:    true,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - preserve all bits",
			setupA:       0x42,
			immediateVal: 0xFF,
			expectedA:    0x42, // 01000010 & 11111111 = 01000010
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - isolate single bit",
			setupA:       0x81, // 10000001
			immediateVal: 0x01, // 00000001
			expectedA:    0x01, // 10000001 & 00000001 = 00000001
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - check bit 7",
			setupA:       0x81, // 10000001
			immediateVal: 0x80, // 10000000
			expectedA:    0x80, // 10000001 & 10000000 = 10000000
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - zero from zero A",
			setupA:       0x00,
			immediateVal: 0xFF,
			expectedA:    0x00, // 00000000 & 11111111 = 00000000
			expectedZ:    true,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - alternating bits",
			setupA:       0xAA, // 10101010
			immediateVal: 0x55, // 01010101
			expectedA:    0x00, // 10101010 & 01010101 = 00000000
			expectedZ:    true,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "AND_A_n - bit pattern matching",
			setupA:       0x3C, // 00111100
			immediateVal: 0x18, // 00011000
			expectedA:    0x18, // 00111100 & 00011000 = 00011000
			expectedZ:    false,
			expectedN:    false,
			expectedH:    true,
			expectedC:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := cpu.AND_A_n(tt.immediateVal)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Instruction cycles")
		})
	}
}

// TestAND_FlagBehavior tests specific flag behavior for AND operations
func TestAND_FlagBehavior(t *testing.T) {
	t.Run("AND always resets N and C flags", func(t *testing.T) {
		cpu := NewCPU()

		// Set all flags initially
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, false)
		cpu.SetFlag(FlagC, true)

		cpu.A = 0xFF
		cpu.B = 0xFF

		cpu.AND_A_B()

		// Verify flag behavior
		assert.False(t, cpu.GetFlag(FlagZ)) // Result is 0xFF (not zero)
		assert.False(t, cpu.GetFlag(FlagN)) // Always reset for AND
		assert.True(t, cpu.GetFlag(FlagH))  // Always set for AND
		assert.False(t, cpu.GetFlag(FlagC)) // Always reset for AND
	})

	t.Run("AND always sets H flag", func(t *testing.T) {
		cpu := NewCPU()

		// Clear all flags initially
		cpu.SetFlag(FlagZ, false)
		cpu.SetFlag(FlagN, false)
		cpu.SetFlag(FlagH, false)
		cpu.SetFlag(FlagC, false)

		cpu.A = 0x00
		cpu.B = 0x00

		cpu.AND_A_B()

		// Verify H flag is always set for AND operations
		assert.True(t, cpu.GetFlag(FlagZ))  // Result is zero
		assert.False(t, cpu.GetFlag(FlagN)) // Always reset for AND
		assert.True(t, cpu.GetFlag(FlagH))  // Always set for AND
		assert.False(t, cpu.GetFlag(FlagC)) // Always reset for AND
	})

	t.Run("AND zero flag accuracy", func(t *testing.T) {
		// Test cases that should produce zero vs non-zero results
		testCases := []struct {
			a, b    uint8
			expectZ bool
			name    string
		}{
			{0x00, 0xFF, true, "0x00 AND 0xFF should be zero"},
			{0xFF, 0x00, true, "0xFF AND 0x00 should be zero"},
			{0xF0, 0x0F, true, "0xF0 AND 0x0F should be zero (no bit overlap)"},
			{0xAA, 0x55, true, "0xAA AND 0x55 should be zero (alternating bits)"},
			{0xFF, 0xFF, false, "0xFF AND 0xFF should not be zero"},
			{0x01, 0x01, false, "0x01 AND 0x01 should not be zero"},
			{0x80, 0x80, false, "0x80 AND 0x80 should not be zero"},
			{0x0F, 0x07, false, "0x0F AND 0x07 should not be zero"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tc.a
				cpu.B = tc.b
				cpu.AND_A_B()
				assert.Equal(t, tc.expectZ, cpu.GetFlag(FlagZ), tc.name)
			})
		}
	})
}

// TestAND_EdgeCases tests edge cases and boundary conditions for AND operations
func TestAND_EdgeCases(t *testing.T) {
	t.Run("AND with maximum values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF
		cpu.B = 0xFF

		cpu.AND_A_B()

		assert.Equal(t, uint8(0xFF), cpu.A, "0xFF AND 0xFF should be 0xFF")
		assert.False(t, cpu.GetFlag(FlagZ), "Result should not be zero")
	})

	t.Run("AND with zero values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		cpu.B = 0x00

		cpu.AND_A_B()

		assert.Equal(t, uint8(0x00), cpu.A, "0x00 AND 0x00 should be 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Result should be zero")
	})

	t.Run("AND bit masking patterns", func(t *testing.T) {
		// Common bit masking scenarios in Game Boy programming
		testCases := []struct {
			a, mask, expected uint8
			description       string
		}{
			{0x87, 0x0F, 0x07, "Extract lower nibble"},
			{0x87, 0xF0, 0x80, "Extract upper nibble"},
			{0xFF, 0x01, 0x01, "Check bit 0"},
			{0xFF, 0x80, 0x80, "Check bit 7"},
			{0x3C, 0x18, 0x18, "Extract middle bits"},
			{0xC3, 0x81, 0x81, "Extract corner bits"},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tc.a
				result := cpu.AND_A_n(tc.mask)
				assert.Equal(t, tc.expected, cpu.A, tc.description)
				assert.Equal(t, uint8(8), result, "Should take 8 cycles")
			})
		}
	})
}

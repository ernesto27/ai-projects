package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSUB_A_Register tests all register-to-register SUB operations
func TestSUB_A_Register(t *testing.T) {
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
		// SUB_A_A tests
		{
			name:        "SUB_A_A - always zero",
			setupA:      0x42,
			setupOther:  0x42, // Not used for SUB_A_A
			setOtherReg: func(cpu *CPU, val uint8) {}, // No-op
			instruction: (*CPU).SUB_A_A,
			expectedA:   0x00,
			expectedZ:   true,  // Always zero
			expectedN:   true,  // Always set for subtraction
			expectedH:   false, // No borrow needed for 0-0
			expectedC:   false, // No underflow for 0-0
		},
		{
			name:        "SUB_A_A - zero input",
			setupA:      0x00,
			setupOther:  0x00,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).SUB_A_A,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_A - maximum value",
			setupA:      0xFF,
			setupOther:  0xFF,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).SUB_A_A,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},

		// SUB_A_B tests
		{
			name:        "SUB_A_B - normal case",
			setupA:      0x30,
			setupOther:  0x10,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).SUB_A_B,
			expectedA:   0x20,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_B - zero result",
			setupA:      0x42,
			setupOther:  0x42,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).SUB_A_B,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_B - half carry (borrow from bit 4)",
			setupA:      0x10,
			setupOther:  0x01,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).SUB_A_B,
			expectedA:   0x0F,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true, // Borrow from bit 4 to bit 3
			expectedC:   false,
		},
		{
			name:        "SUB_A_B - carry (underflow)",
			setupA:      0x10,
			setupOther:  0x20,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).SUB_A_B,
			expectedA:   0xF0, // 0x10 - 0x20 = -0x10 = 0xF0 in 8-bit
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   true, // A < B, so underflow
		},
		{
			name:        "SUB_A_B - both half carry and carry",
			setupA:      0x00,
			setupOther:  0x01,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.B = val },
			instruction: (*CPU).SUB_A_B,
			expectedA:   0xFF,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true, // Borrow needed
			expectedC:   true, // Underflow
		},

		// SUB_A_C tests
		{
			name:        "SUB_A_C - normal case",
			setupA:      0x80,
			setupOther:  0x40,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.C = val },
			instruction: (*CPU).SUB_A_C,
			expectedA:   0x40,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_C - zero result",
			setupA:      0x7F,
			setupOther:  0x7F,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.C = val },
			instruction: (*CPU).SUB_A_C,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_C - half carry",
			setupA:      0x20,
			setupOther:  0x01,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.C = val },
			instruction: (*CPU).SUB_A_C,
			expectedA:   0x1F,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   false,
		},

		// SUB_A_D tests
		{
			name:        "SUB_A_D - normal case",
			setupA:      0xFF,
			setupOther:  0x01,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.D = val },
			instruction: (*CPU).SUB_A_D,
			expectedA:   0xFE,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_D - underflow",
			setupA:      0x05,
			setupOther:  0x10,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.D = val },
			instruction: (*CPU).SUB_A_D,
			expectedA:   0xF5, // 0x05 - 0x10 = -0x0B = 0xF5
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false, // No borrow needed in low nibble: 0x5 - 0x0 = 0x5
			expectedC:   true, // Underflow: 0x05 < 0x10
		},

		// SUB_A_E tests
		{
			name:        "SUB_A_E - normal case",
			setupA:      0x50,
			setupOther:  0x30,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.E = val },
			instruction: (*CPU).SUB_A_E,
			expectedA:   0x20,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_E - zero from maximum values",
			setupA:      0xFF,
			setupOther:  0xFF,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.E = val },
			instruction: (*CPU).SUB_A_E,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},

		// SUB_A_H tests
		{
			name:        "SUB_A_H - normal case",
			setupA:      0x88,
			setupOther:  0x44,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.H = val },
			instruction: (*CPU).SUB_A_H,
			expectedA:   0x44,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_H - half carry edge case",
			setupA:      0x40,
			setupOther:  0x01,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.H = val },
			instruction: (*CPU).SUB_A_H,
			expectedA:   0x3F,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   false,
		},

		// SUB_A_L tests
		{
			name:        "SUB_A_L - normal case",
			setupA:      0x99,
			setupOther:  0x33,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.L = val },
			instruction: (*CPU).SUB_A_L,
			expectedA:   0x66,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_L - maximum underflow",
			setupA:      0x00,
			setupOther:  0xFF,
			setOtherReg: func(cpu *CPU, val uint8) { cpu.L = val },
			instruction: (*CPU).SUB_A_L,
			expectedA:   0x01, // 0x00 - 0xFF = -0xFF = 0x01 in 8-bit
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   true,
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

// TestSUB_A_HL tests the SUB A,(HL) memory operation
func TestSUB_A_HL(t *testing.T) {
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
			name:        "SUB_A_HL - normal case",
			setupA:      0x50,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x20,
			expectedA:   0x30,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_HL - zero result",
			setupA:      0x42,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x42,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "SUB_A_HL - half carry",
			setupA:      0x10,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x01,
			expectedA:   0x0F,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "SUB_A_HL - carry (underflow)",
			setupA:      0x10,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x20,
			expectedA:   0xF0,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   true,
		},
		{
			name:        "SUB_A_HL - both half carry and carry",
			setupA:      0x00,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x01,
			expectedA:   0xFF,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   true,
		},
		{
			name:        "SUB_A_HL - maximum values",
			setupA:      0xFF,
			setupH:      0x90,
			setupL:      0x00,
			memoryValue: 0xFF,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()

			cpu.A = tt.setupA
			cpu.H = tt.setupH
			cpu.L = tt.setupL

			address := uint16(tt.setupH)<<8 | uint16(tt.setupL)
			mmu.WriteByte(address, tt.memoryValue)

			cycles := cpu.SUB_A_HL(mmu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Instruction cycles")
		})
	}
}

// TestSUB_A_n tests the SUB A,n immediate value operation
func TestSUB_A_n(t *testing.T) {
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
			name:         "SUB_A_n - normal case",
			setupA:       0x50,
			immediateVal: 0x20,
			expectedA:    0x30,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "SUB_A_n - zero result",
			setupA:       0x42,
			immediateVal: 0x42,
			expectedA:    0x00,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "SUB_A_n - half carry",
			setupA:       0x10,
			immediateVal: 0x01,
			expectedA:    0x0F,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    false,
		},
		{
			name:         "SUB_A_n - carry (underflow)",
			setupA:       0x10,
			immediateVal: 0x20,
			expectedA:    0xF0,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    true,
		},
		{
			name:         "SUB_A_n - both half carry and carry",
			setupA:       0x00,
			immediateVal: 0x01,
			expectedA:    0xFF,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    true,
		},
		{
			name:         "SUB_A_n - maximum values",
			setupA:       0xFF,
			immediateVal: 0xFF,
			expectedA:    0x00,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "SUB_A_n - subtract zero",
			setupA:       0x42,
			immediateVal: 0x00,
			expectedA:    0x42,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "SUB_A_n - subtract from zero",
			setupA:       0x00,
			immediateVal: 0x80,
			expectedA:    0x80,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := cpu.SUB_A_n(tt.immediateVal)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Instruction cycles")
		})
	}
}

// TestSUB_EdgeCases tests edge cases and boundary conditions for SUB operations
func TestSUB_EdgeCases(t *testing.T) {
	t.Run("SUB flag preservation", func(t *testing.T) {
		cpu := NewCPU()
		
		// Set some initial flags that should be preserved or overwritten correctly
		cpu.SetFlag(FlagZ, false)
		cpu.SetFlag(FlagN, false)
		cpu.SetFlag(FlagH, false)
		cpu.SetFlag(FlagC, false)
		
		cpu.A = 0x10
		cpu.B = 0x08
		
		cpu.SUB_A_B()
		
		// Verify that all flags are properly set/cleared
		assert.Equal(t, uint8(0x08), cpu.A)
		assert.False(t, cpu.GetFlag(FlagZ)) // Result is not zero
		assert.True(t, cpu.GetFlag(FlagN))  // Always set for SUB
		assert.True(t, cpu.GetFlag(FlagH))  // Half-carry: 0x0 < 0x8
		assert.False(t, cpu.GetFlag(FlagC)) // No carry
	})
	
	t.Run("SUB half-carry calculation accuracy", func(t *testing.T) {
		// Test the exact boundary conditions for half-carry
		testCases := []struct {
			a, b     uint8
			expectH  bool
			name     string
		}{
			{0x10, 0x01, true, "0x10 - 0x01 should have half-carry"},
			{0x20, 0x01, true, "0x20 - 0x01 should have half-carry"},
			{0x18, 0x08, false, "0x18 - 0x08 should not have half-carry"},
			{0x1F, 0x0F, false, "0x1F - 0x0F should not have half-carry"},
			{0x0F, 0x0F, false, "0x0F - 0x0F should not have half-carry"},
			{0x08, 0x09, true, "0x08 - 0x09 should have half-carry (borrow)"},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tc.a
				cpu.B = tc.b
				cpu.SUB_A_B()
				assert.Equal(t, tc.expectH, cpu.GetFlag(FlagH), tc.name)
			})
		}
	})
	
	t.Run("SUB carry calculation accuracy", func(t *testing.T) {
		// Test carry flag for underflow conditions
		testCases := []struct {
			a, b      uint8
			expectC   bool
			name      string
		}{
			{0x50, 0x30, false, "0x50 - 0x30 should not have carry"},
			{0x30, 0x50, true, "0x30 - 0x50 should have carry (underflow)"},
			{0x00, 0x01, true, "0x00 - 0x01 should have carry"},
			{0xFF, 0xFF, false, "0xFF - 0xFF should not have carry"},
			{0x01, 0xFF, true, "0x01 - 0xFF should have carry"},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tc.a
				cpu.B = tc.b
				cpu.SUB_A_B()
				assert.Equal(t, tc.expectC, cpu.GetFlag(FlagC), tc.name)
			})
		}
	})
}

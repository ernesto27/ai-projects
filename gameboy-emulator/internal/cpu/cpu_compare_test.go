package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === CP (Compare) Instruction Tests ===
// These tests verify the compare operations which are crucial for conditional logic
// CP instructions perform A - operand but don't store the result (A unchanged)

func TestCP_A_A(t *testing.T) {
	tests := []struct {
		name           string
		setupA         uint8
		expectedZ      bool
		expectedN      bool
		expectedH      bool
		expectedC      bool
		expectedCycles uint8
	}{
		{
			name:           "Compare A with A - zero value",
			setupA:         0x00,
			expectedZ:      true,  // A == A always true
			expectedN:      true,  // Always set for compare
			expectedH:      false, // No half-carry when equal
			expectedC:      false, // No carry when equal
			expectedCycles: 4,
		},
		{
			name:           "Compare A with A - max value",
			setupA:         0xFF,
			expectedZ:      true,  // A == A always true
			expectedN:      true,  // Always set for compare
			expectedH:      false, // No half-carry when equal
			expectedC:      false, // No carry when equal
			expectedCycles: 4,
		},
		{
			name:           "Compare A with A - middle value",
			setupA:         0x80,
			expectedZ:      true,  // A == A always true
			expectedN:      true,  // Always set for compare
			expectedH:      false, // No half-carry when equal
			expectedC:      false, // No carry when equal
			expectedCycles: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := cpu.CP_A_A()

			// Verify A register is unchanged
			assert.Equal(t, tt.setupA, cpu.A, "A register should be unchanged")

			// Verify flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")

			// Verify cycles
			assert.Equal(t, tt.expectedCycles, cycles, "Instruction cycles")
		})
	}
}

func TestCP_A_Register(t *testing.T) {
	tests := []struct {
		name       string
		setupA     uint8
		setupOther uint8
		expectedZ  bool
		expectedN  bool
		expectedH  bool
		expectedC  bool
	}{
		// Equal values tests
		{
			name:       "Equal values - zero",
			setupA:     0x00,
			setupOther: 0x00,
			expectedZ:  true,
			expectedN:  true,
			expectedH:  false,
			expectedC:  false,
		},
		{
			name:       "Equal values - max",
			setupA:     0xFF,
			setupOther: 0xFF,
			expectedZ:  true,
			expectedN:  true,
			expectedH:  false,
			expectedC:  false,
		},

		// A > other tests
		{
			name:       "A greater than other",
			setupA:     0x10,
			setupOther: 0x05,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  true, // 0x0 < 0x5 in lower nibble, half-carry needed
			expectedC:  false,
		},
		{
			name:       "A greater - no half carry",
			setupA:     0x20,
			setupOther: 0x10,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  false, // 0x0 >= 0x0, no half-carry
			expectedC:  false,
		},

		// A < other tests
		{
			name:       "A less than other",
			setupA:     0x05,
			setupOther: 0x10,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  false, // 0x5 >= 0x0 in lower nibble, no half-carry needed
			expectedC:  true,  // A < other, so carry set
		},
		{
			name:       "A less - with half carry",
			setupA:     0x0F,
			setupOther: 0x10,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  false, // 0xF >= 0x0, no half-carry needed
			expectedC:  true,  // A < other, so carry set
		},

		// Half-carry edge cases
		{
			name:       "Half-carry boundary test 1",
			setupA:     0x10,
			setupOther: 0x01,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  true, // 0x0 < 0x1, half-carry needed
			expectedC:  false,
		},
		{
			name:       "Half-carry boundary test 2",
			setupA:     0x08,
			setupOther: 0x09,
			expectedZ:  false,
			expectedN:  true,
			expectedH:  true, // 0x8 < 0x9, half-carry needed
			expectedC:  true, // A < other
		},
	}

	for _, tt := range tests {
		t.Run("CP_A_B "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.B = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_B()

			// Verify A register is unchanged
			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")

			// Verify flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})

		t.Run("CP_A_C "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.C = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_C()

			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})

		t.Run("CP_A_D "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.D = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_D()

			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})

		t.Run("CP_A_E "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.E = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_E()

			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})

		t.Run("CP_A_H "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.H = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_H()

			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})

		t.Run("CP_A_L "+tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.L = tt.setupOther
			originalA := cpu.A

			cycles := cpu.CP_A_L()

			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(4), cycles, "Cycles should be 4")
		})
	}
}

func TestCP_A_HL(t *testing.T) {
	tests := []struct {
		name        string
		setupA      uint8
		memoryValue uint8
		memoryAddr  uint16
		expectedZ   bool
		expectedN   bool
		expectedH   bool
		expectedC   bool
	}{
		{
			name:        "Compare A with memory - equal values",
			setupA:      0x42,
			memoryValue: 0x42,
			memoryAddr:  0x8000,
			expectedZ:   true,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "Compare A with memory - A greater",
			setupA:      0x50,
			memoryValue: 0x30,
			memoryAddr:  0x8100,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "Compare A with memory - A less",
			setupA:      0x30,
			memoryValue: 0x50,
			memoryAddr:  0x8200,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   false, // 0x0 >= 0x0 in lower nibble, no half-carry needed
			expectedC:   true,
		},
		{
			name:        "Compare A with memory - half carry test",
			setupA:      0x10,
			memoryValue: 0x01,
			memoryAddr:  0x8300,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true, // 0x0 < 0x1 in lower nibble
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()

			cpu.A = tt.setupA
			cpu.SetHL(tt.memoryAddr)
			mmu.WriteByte(tt.memoryAddr, tt.memoryValue)
			originalA := cpu.A

			cycles := cpu.CP_A_HL(mmu)

			// Verify A register is unchanged
			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")

			// Verify memory is unchanged
			assert.Equal(t, tt.memoryValue, mmu.ReadByte(tt.memoryAddr), "Memory should be unchanged")

			// Verify flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Cycles should be 8")
		})
	}
}

func TestCP_A_n(t *testing.T) {
	tests := []struct {
		name         string
		setupA       uint8
		immediateVal uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
	}{
		{
			name:         "Compare A with immediate - equal",
			setupA:       0x7F,
			immediateVal: 0x7F,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "Compare A with immediate - A greater",
			setupA:       0x80,
			immediateVal: 0x40,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "Compare A with immediate - A less",
			setupA:       0x40,
			immediateVal: 0x80,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false, // 0x0 >= 0x0 in lower nibble, no half-carry needed
			expectedC:    true,
		},
		{
			name:         "Compare A with 0",
			setupA:       0x01,
			immediateVal: 0x00,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
		},
		{
			name:         "Compare 0 with immediate",
			setupA:       0x00,
			immediateVal: 0x01,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true, // 0x0 < 0x1
			expectedC:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			originalA := cpu.A

			cycles := cpu.CP_A_n(tt.immediateVal)

			// Verify A register is unchanged
			assert.Equal(t, originalA, cpu.A, "A register should be unchanged")

			// Verify flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag")
			assert.Equal(t, uint8(8), cycles, "Cycles should be 8")
		})
	}
}

// Test edge cases and boundary conditions
func TestCP_EdgeCases(t *testing.T) {
	t.Run("All registers preserve their values", func(t *testing.T) {
		cpu := NewCPU()
		mmu := memory.NewMMU()

		// Set up all registers with specific values
		cpu.A = 0x42
		cpu.B = 0x10
		cpu.C = 0x20
		cpu.D = 0x30
		cpu.E = 0x40
		cpu.H = 0x50
		cpu.L = 0x60
		cpu.SetHL(0x8000)
		mmu.WriteByte(0x8000, 0x70)

		// Store original values
		origA, origB, origC, origD, origE, origH, origL := cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L

		// Execute all CP instructions
		cpu.CP_A_A()
		cpu.CP_A_B()
		cpu.CP_A_C()
		cpu.CP_A_D()
		cpu.CP_A_E()
		cpu.CP_A_H()
		cpu.CP_A_L()
		cpu.CP_A_HL(mmu)
		cpu.CP_A_n(0x80)

		// Verify all registers are unchanged
		assert.Equal(t, origA, cpu.A, "A register should be unchanged")
		assert.Equal(t, origB, cpu.B, "B register should be unchanged")
		assert.Equal(t, origC, cpu.C, "C register should be unchanged")
		assert.Equal(t, origD, cpu.D, "D register should be unchanged")
		assert.Equal(t, origE, cpu.E, "E register should be unchanged")
		assert.Equal(t, origH, cpu.H, "H register should be unchanged")
		assert.Equal(t, origL, cpu.L, "L register should be unchanged")
		assert.Equal(t, uint8(0x70), mmu.ReadByte(0x8000), "Memory should be unchanged")
	})

	t.Run("Flag behavior consistency", func(t *testing.T) {
		cpu := NewCPU()

		// Test that N flag is always set for all CP operations
		cpu.A = 0x50
		cpu.B = 0x30

		cpu.CP_A_A()
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be set for CP_A_A")

		cpu.CP_A_B()
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be set for CP_A_B")

		cpu.CP_A_n(0x20)
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be set for CP_A_n")
	})
}

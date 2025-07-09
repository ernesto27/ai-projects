package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestADD_A_Register tests all register-to-register ADD operations
func TestADD_A_Register(t *testing.T) {
	tests := []struct {
		name        string
		setupA      uint8
		setupOther  uint8
		instruction func(*CPU) uint8
		expectedA   uint8
		expectedZ   bool
		expectedN   bool
		expectedH   bool
		expectedC   bool
	}{
		// ADD_A_A tests
		{
			name:        "ADD_A_A - normal case",
			setupA:      0x10,
			setupOther:  0x10, // Not used for ADD_A_A
			instruction: (*CPU).ADD_A_A,
			expectedA:   0x20,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_A - zero result",
			setupA:      0x00,
			setupOther:  0x00,
			instruction: (*CPU).ADD_A_A,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_A - half carry",
			setupA:      0x08,
			setupOther:  0x08,
			instruction: (*CPU).ADD_A_A,
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "ADD_A_A - carry",
			setupA:      0x80,
			setupOther:  0x80,
			instruction: (*CPU).ADD_A_A,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   true,
		},
		{
			name:        "ADD_A_A - both carry and half carry",
			setupA:      0x88,
			setupOther:  0x88,
			instruction: (*CPU).ADD_A_A,
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   true,
		},
		// ADD_A_B tests
		{
			name:        "ADD_A_B - normal case",
			setupA:      0x10,
			setupOther:  0x20,
			instruction: func(cpu *CPU) uint8 { cpu.B = 0x20; return cpu.ADD_A_B() },
			expectedA:   0x30,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_B - zero result",
			setupA:      0x00,
			setupOther:  0x00,
			instruction: func(cpu *CPU) uint8 { cpu.B = 0x00; return cpu.ADD_A_B() },
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_B - half carry",
			setupA:      0x08,
			setupOther:  0x08,
			instruction: func(cpu *CPU) uint8 { cpu.B = 0x08; return cpu.ADD_A_B() },
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "ADD_A_B - carry",
			setupA:      0x80,
			setupOther:  0x80,
			instruction: func(cpu *CPU) uint8 { cpu.B = 0x80; return cpu.ADD_A_B() },
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   true,
		},
		{
			name:        "ADD_A_B - both carry and half carry",
			setupA:      0x88,
			setupOther:  0x88,
			instruction: func(cpu *CPU) uint8 { cpu.B = 0x88; return cpu.ADD_A_B() },
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   true,
		},
		// ADD_A_C tests
		{
			name:        "ADD_A_C - normal case",
			setupA:      0x15,
			setupOther:  0x25,
			instruction: func(cpu *CPU) uint8 { cpu.C = 0x25; return cpu.ADD_A_C() },
			expectedA:   0x3A,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// ADD_A_D tests
		{
			name:        "ADD_A_D - normal case",
			setupA:      0x30,
			setupOther:  0x40,
			instruction: func(cpu *CPU) uint8 { cpu.D = 0x40; return cpu.ADD_A_D() },
			expectedA:   0x70,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// ADD_A_E tests
		{
			name:        "ADD_A_E - normal case",
			setupA:      0x05,
			setupOther:  0x0A,
			instruction: func(cpu *CPU) uint8 { cpu.E = 0x0A; return cpu.ADD_A_E() },
			expectedA:   0x0F,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// ADD_A_H tests
		{
			name:        "ADD_A_H - normal case",
			setupA:      0x12,
			setupOther:  0x34,
			instruction: func(cpu *CPU) uint8 { cpu.H = 0x34; return cpu.ADD_A_H() },
			expectedA:   0x46,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// ADD_A_L tests
		{
			name:        "ADD_A_L - normal case",
			setupA:      0x56,
			setupOther:  0x78,
			instruction: func(cpu *CPU) uint8 { cpu.L = 0x78; return cpu.ADD_A_L() },
			expectedA:   0xCE,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := tt.instruction(cpu)

			// Check result
			assert.Equal(t, tt.expectedA, cpu.A, "A register value")

			// Check flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Z flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "N flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "H flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "C flag")

			// Check cycles
			assert.Equal(t, uint8(4), cycles, "cycles")
		})
	}
}

// TestADD_A_HL tests memory-based ADD operation
func TestADD_A_HL(t *testing.T) {
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
			name:        "ADD_A_HL - normal case",
			setupA:      0x10,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x20,
			expectedA:   0x30,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_HL - zero result",
			setupA:      0x00,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x00,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "ADD_A_HL - half carry",
			setupA:      0x08,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x08,
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "ADD_A_HL - carry",
			setupA:      0x80,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x80,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   true,
		},
		{
			name:        "ADD_A_HL - both carry and half carry",
			setupA:      0x88,
			setupH:      0x80,
			setupL:      0x00,
			memoryValue: 0x88,
			expectedA:   0x10,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   true,
			expectedC:   true,
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

			cycles := cpu.ADD_A_HL(mmu)

			// Check result
			assert.Equal(t, tt.expectedA, cpu.A, "A register value")

			// Check flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Z flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "N flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "H flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "C flag")

			// Check cycles
			assert.Equal(t, uint8(8), cycles, "cycles")
		})
	}
}

// TestADD_A_n tests immediate value ADD operation
func TestADD_A_n(t *testing.T) {
	tests := []struct {
		name      string
		setupA    uint8
		immediate uint8
		expectedA uint8
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
	}{
		{
			name:      "ADD_A_n - normal case",
			setupA:    0x10,
			immediate: 0x20,
			expectedA: 0x30,
			expectedZ: false,
			expectedN: false,
			expectedH: false,
			expectedC: false,
		},
		{
			name:      "ADD_A_n - zero result",
			setupA:    0x00,
			immediate: 0x00,
			expectedA: 0x00,
			expectedZ: true,
			expectedN: false,
			expectedH: false,
			expectedC: false,
		},
		{
			name:      "ADD_A_n - half carry",
			setupA:    0x08,
			immediate: 0x08,
			expectedA: 0x10,
			expectedZ: false,
			expectedN: false,
			expectedH: true,
			expectedC: false,
		},
		{
			name:      "ADD_A_n - carry",
			setupA:    0x80,
			immediate: 0x80,
			expectedA: 0x00,
			expectedZ: true,
			expectedN: false,
			expectedH: false,
			expectedC: true,
		},
		{
			name:      "ADD_A_n - both carry and half carry",
			setupA:    0x88,
			immediate: 0x88,
			expectedA: 0x10,
			expectedZ: false,
			expectedN: false,
			expectedH: true,
			expectedC: true,
		},
		{
			name:      "ADD_A_n - maximum values",
			setupA:    0xFF,
			immediate: 0x01,
			expectedA: 0x00,
			expectedZ: true,
			expectedN: false,
			expectedH: true,
			expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := cpu.ADD_A_n(tt.immediate)

			// Check result
			assert.Equal(t, tt.expectedA, cpu.A, "A register value")

			// Check flags
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Z flag")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "N flag")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "H flag")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "C flag")

			// Check cycles
			assert.Equal(t, uint8(8), cycles, "cycles")
		})
	}
}

// TestADD_EdgeCases tests edge cases for ADD operations
func TestADD_EdgeCases(t *testing.T) {
	t.Run("ADD_A_A with 0xFF", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF

		cycles := cpu.ADD_A_A()

		assert.Equal(t, uint8(0xFE), cpu.A, "A register value")
		assert.True(t, cpu.GetFlag(FlagC), "carry flag should be set")
		assert.True(t, cpu.GetFlag(FlagH), "half-carry flag should be set")
		assert.False(t, cpu.GetFlag(FlagZ), "zero flag should be clear")
		assert.Equal(t, uint8(4), cycles, "cycles")
	})

	t.Run("ADD_A_B with maximum boundary", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x7F
		cpu.B = 0x80

		cycles := cpu.ADD_A_B()

		assert.Equal(t, uint8(0xFF), cpu.A, "A register value")
		assert.False(t, cpu.GetFlag(FlagC), "carry flag should be clear")
		assert.False(t, cpu.GetFlag(FlagH), "half-carry flag should be clear")
		assert.False(t, cpu.GetFlag(FlagZ), "zero flag should be clear")
		assert.Equal(t, uint8(4), cycles, "cycles")
	})
}

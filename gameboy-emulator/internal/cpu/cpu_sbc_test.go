package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSBC_A_A(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		carryFlag    bool
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
		expectedCycles uint8
	}{
		{
			name:         "SBC A,A with no carry - results in 0",
			initialA:     0x42,
			carryFlag:    false,
			expectedA:    0x00,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,A with carry - results in 0xFF (underflow)",
			initialA:     0x42,
			carryFlag:    true,
			expectedA:    0xFF,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    true,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,A with A=0 and no carry",
			initialA:     0x00,
			carryFlag:    false,
			expectedA:    0x00,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,A with A=0 and carry",
			initialA:     0x00,
			carryFlag:    true,
			expectedA:    0xFF,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    true,
			expectedCycles: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.initialA
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := cpu.SBC_A_A()

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should be %02X", tt.expectedA)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %v", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %v", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %v", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %v", tt.expectedC)
			assert.Equal(t, tt.expectedCycles, cycles, "Should take %d cycles", tt.expectedCycles)
		})
	}
}

func TestSBC_A_B(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		initialB     uint8
		carryFlag    bool
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
		expectedCycles uint8
	}{
		{
			name:         "SBC A,B basic subtraction with no carry",
			initialA:     0x50,
			initialB:     0x30,
			carryFlag:    false,
			expectedA:    0x20,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,B basic subtraction with carry",
			initialA:     0x50,
			initialB:     0x30,
			carryFlag:    true,
			expectedA:    0x1F,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    false,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,B underflow with no carry",
			initialA:     0x20,
			initialB:     0x30,
			carryFlag:    false,
			expectedA:    0xF0,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    false,
			expectedC:    true,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,B underflow with carry",
			initialA:     0x20,
			initialB:     0x30,
			carryFlag:    true,
			expectedA:    0xEF,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    true,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,B resulting in zero",
			initialA:     0x42,
			initialB:     0x42,
			carryFlag:    false,
			expectedA:    0x00,
			expectedZ:    true,
			expectedN:    true,
			expectedH:    false,
			expectedC:    false,
			expectedCycles: 4,
		},
		{
			name:         "SBC A,B half-carry test",
			initialA:     0x10,
			initialB:     0x01,
			carryFlag:    true,
			expectedA:    0x0E,
			expectedZ:    false,
			expectedN:    true,
			expectedH:    true,
			expectedC:    false,
			expectedCycles: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.initialA
			cpu.B = tt.initialB
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := cpu.SBC_A_B()

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should be %02X", tt.expectedA)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %v", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %v", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %v", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %v", tt.expectedC)
			assert.Equal(t, tt.expectedCycles, cycles, "Should take %d cycles", tt.expectedCycles)
		})
	}
}

func TestSBC_A_C(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		initialC     uint8
		carryFlag    bool
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
	}{
		{
			name:      "SBC A,C with carry produces correct result",
			initialA:  0x80,
			initialC:  0x40,
			carryFlag: true,
			expectedA: 0x3F,
			expectedZ: false,
			expectedN: true,
			expectedH: true,
			expectedC: false,
		},
		{
			name:      "SBC A,C maximum values with carry",
			initialA:  0xFF,
			initialC:  0xFF,
			carryFlag: true,
			expectedA: 0xFF,
			expectedZ: false,
			expectedN: true,
			expectedH: true,
			expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.initialA
			cpu.C = tt.initialC
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := cpu.SBC_A_C()

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should be %02X", tt.expectedA)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %v", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %v", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %v", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %v", tt.expectedC)
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
		})
	}
}

func TestSBC_A_D(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x60
	cpu.D = 0x30
	cpu.SetFlag(FlagC, true)

	cycles := cpu.SBC_A_D()

	assert.Equal(t, uint8(0x2F), cpu.A, "A should be 0x2F after SBC A,D")
	assert.Equal(t, false, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be true")
	assert.Equal(t, true, cpu.GetFlag(FlagH), "Half-carry flag should be true")
	assert.Equal(t, false, cpu.GetFlag(FlagC), "Carry flag should be false")
	assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
}

func TestSBC_A_E(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x10
	cpu.E = 0x20
	cpu.SetFlag(FlagC, false)

	cycles := cpu.SBC_A_E()

	assert.Equal(t, uint8(0xF0), cpu.A, "A should be 0xF0 after underflow")
	assert.Equal(t, false, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be true")
	assert.Equal(t, false, cpu.GetFlag(FlagH), "Half-carry flag should be false")
	assert.Equal(t, true, cpu.GetFlag(FlagC), "Carry flag should be true")
	assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
}

func TestSBC_A_H(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xFF
	cpu.H = 0x01
	cpu.SetFlag(FlagC, true)

	cycles := cpu.SBC_A_H()

	assert.Equal(t, uint8(0xFD), cpu.A, "A should be 0xFD")
	assert.Equal(t, false, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be true")
	assert.Equal(t, false, cpu.GetFlag(FlagH), "Half-carry flag should be false")
	assert.Equal(t, false, cpu.GetFlag(FlagC), "Carry flag should be false")
	assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
}

func TestSBC_A_L(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x80
	cpu.L = 0x80
	cpu.SetFlag(FlagC, true)

	cycles := cpu.SBC_A_L()

	assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF after underflow")
	assert.Equal(t, false, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be true")
	assert.Equal(t, true, cpu.GetFlag(FlagH), "Half-carry flag should be true")
	assert.Equal(t, true, cpu.GetFlag(FlagC), "Carry flag should be true")
	assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
}

func TestSBC_A_HL(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		memoryValue  uint8
		carryFlag    bool
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
	}{
		{
			name:        "SBC A,(HL) basic memory operation",
			initialA:    0x50,
			memoryValue: 0x25,
			carryFlag:   false,
			expectedA:   0x2B,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "SBC A,(HL) with carry",
			initialA:    0x50,
			memoryValue: 0x25,
			carryFlag:   true,
			expectedA:   0x2A,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   false,
		},
		{
			name:        "SBC A,(HL) underflow with memory",
			initialA:    0x10,
			memoryValue: 0x20,
			carryFlag:   true,
			expectedA:   0xEF,
			expectedZ:   false,
			expectedN:   true,
			expectedH:   true,
			expectedC:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()
			
			cpu.A = tt.initialA
			cpu.SetHL(0x8000)
			cpu.SetFlag(FlagC, tt.carryFlag)
			mmu.WriteByte(0x8000, tt.memoryValue)

			cycles := cpu.SBC_A_HL(mmu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should be %02X", tt.expectedA)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %v", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %v", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %v", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %v", tt.expectedC)
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		})
	}
}

func TestSBC_A_n(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		value        uint8
		carryFlag    bool
		expectedA    uint8
		expectedZ    bool
		expectedN    bool
		expectedH    bool
		expectedC    bool
	}{
		{
			name:      "SBC A,n basic immediate operation",
			initialA:  0x75,
			value:     0x25,
			carryFlag: false,
			expectedA: 0x50,
			expectedZ: false,
			expectedN: true,
			expectedH: false,
			expectedC: false,
		},
		{
			name:      "SBC A,n with carry",
			initialA:  0x75,
			value:     0x25,
			carryFlag: true,
			expectedA: 0x4F,
			expectedZ: false,
			expectedN: true,
			expectedH: true,
			expectedC: false,
		},
		{
			name:      "SBC A,n resulting in zero",
			initialA:  0x42,
			value:     0x42,
			carryFlag: false,
			expectedA: 0x00,
			expectedZ: true,
			expectedN: true,
			expectedH: false,
			expectedC: false,
		},
		{
			name:      "SBC A,n underflow scenario",
			initialA:  0x05,
			value:     0x10,
			carryFlag: true,
			expectedA: 0xF4,
			expectedZ: false,
			expectedN: true,
			expectedH: false,
			expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.initialA
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := cpu.SBC_A_n(tt.value)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should be %02X", tt.expectedA)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %v", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %v", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %v", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %v", tt.expectedC)
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		})
	}
}

// Test all SBC instructions preserve other registers
func TestSBC_RegisterPreservation(t *testing.T) {
	registers := []struct {
		name string
		test func(*CPU)
	}{
		{"SBC_A_B", func(cpu *CPU) { cpu.SBC_A_B() }},
		{"SBC_A_C", func(cpu *CPU) { cpu.SBC_A_C() }},
		{"SBC_A_D", func(cpu *CPU) { cpu.SBC_A_D() }},
		{"SBC_A_E", func(cpu *CPU) { cpu.SBC_A_E() }},
		{"SBC_A_H", func(cpu *CPU) { cpu.SBC_A_H() }},
		{"SBC_A_L", func(cpu *CPU) { cpu.SBC_A_L() }},
		{"SBC_A_A", func(cpu *CPU) { cpu.SBC_A_A() }},
	}

	for _, reg := range registers {
		t.Run(reg.name, func(t *testing.T) {
			cpu := NewCPU()
			originalB := uint8(0x11)
			originalC := uint8(0x22)
			originalD := uint8(0x33)
			originalE := uint8(0x44)
			originalH := uint8(0x55)
			originalL := uint8(0x66)
			originalSP := uint16(0x7788)
			originalPC := uint16(0x9900)

			cpu.B = originalB
			cpu.C = originalC
			cpu.D = originalD
			cpu.E = originalE
			cpu.H = originalH
			cpu.L = originalL
			cpu.SP = originalSP
			cpu.PC = originalPC

			reg.test(cpu)

			// All registers except A should be preserved
			if reg.name != "SBC_A_B" {
				assert.Equal(t, originalB, cpu.B, "Register B should be preserved")
			}
			if reg.name != "SBC_A_C" {
				assert.Equal(t, originalC, cpu.C, "Register C should be preserved")
			}
			if reg.name != "SBC_A_D" {
				assert.Equal(t, originalD, cpu.D, "Register D should be preserved")
			}
			if reg.name != "SBC_A_E" {
				assert.Equal(t, originalE, cpu.E, "Register E should be preserved")
			}
			if reg.name != "SBC_A_H" {
				assert.Equal(t, originalH, cpu.H, "Register H should be preserved")
			}
			if reg.name != "SBC_A_L" {
				assert.Equal(t, originalL, cpu.L, "Register L should be preserved")
			}
			assert.Equal(t, originalSP, cpu.SP, "Stack Pointer should be preserved")
			assert.Equal(t, originalPC, cpu.PC, "Program Counter should be preserved")
		})
	}
}

// Test SBC edge cases
func TestSBC_EdgeCases(t *testing.T) {
	t.Run("SBC maximum values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF
		cpu.B = 0xFF
		cpu.SetFlag(FlagC, false)

		cycles := cpu.SBC_A_B()

		assert.Equal(t, uint8(0x00), cpu.A, "0xFF - 0xFF should equal 0x00")
		assert.Equal(t, true, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be set")
		assert.Equal(t, false, cpu.GetFlag(FlagH), "Half-carry flag should be clear")
		assert.Equal(t, false, cpu.GetFlag(FlagC), "Carry flag should be clear")
		assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
	})

	t.Run("SBC minimum values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		cpu.B = 0x00
		cpu.SetFlag(FlagC, false)

		cycles := cpu.SBC_A_B()

		assert.Equal(t, uint8(0x00), cpu.A, "0x00 - 0x00 should equal 0x00")
		assert.Equal(t, true, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.Equal(t, true, cpu.GetFlag(FlagN), "Subtract flag should be set")
		assert.Equal(t, false, cpu.GetFlag(FlagH), "Half-carry flag should be clear")
		assert.Equal(t, false, cpu.GetFlag(FlagC), "Carry flag should be clear")
		assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
	})
}
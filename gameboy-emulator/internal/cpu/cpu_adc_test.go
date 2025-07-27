package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestADC_A_Register tests all ADC register-to-register operations
func TestADC_A_Register(t *testing.T) {
	tests := []struct {
		name      string
		setupA    uint8
		setupReg  uint8
		regName   string
		carryFlag bool
		expectedA uint8
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
		adcFunc   func(*CPU) uint8
	}{
		// ADC_A_B Tests
		{
			name:   "ADC_A_B - normal case without initial carry",
			setupA: 0x10, setupReg: 0x20, regName: "B", carryFlag: false,
			expectedA: 0x30, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x20; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - normal case with initial carry",
			setupA: 0x10, setupReg: 0x20, regName: "B", carryFlag: true,
			expectedA: 0x31, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x20; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - zero result without carry",
			setupA: 0x00, setupReg: 0x00, regName: "B", carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x00; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - zero result with carry",
			setupA: 0x00, setupReg: 0x00, regName: "B", carryFlag: true,
			expectedA: 0x01, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x00; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - half carry without initial carry",
			setupA: 0x0F, setupReg: 0x01, regName: "B", carryFlag: false,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x01; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - half carry with initial carry",
			setupA: 0x0E, setupReg: 0x01, regName: "B", carryFlag: true,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x01; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - carry without initial carry",
			setupA: 0xFF, setupReg: 0x01, regName: "B", carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x01; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - carry with initial carry",
			setupA: 0xFE, setupReg: 0x01, regName: "B", carryFlag: true,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0x01; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - both half carry and carry",
			setupA: 0xFF, setupReg: 0xFF, regName: "B", carryFlag: false,
			expectedA: 0xFE, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0xFF; return cpu.ADC_A_B() },
		},
		{
			name:   "ADC_A_B - both half carry and carry with initial carry",
			setupA: 0xFF, setupReg: 0xFF, regName: "B", carryFlag: true,
			expectedA: 0xFF, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
			adcFunc: func(cpu *CPU) uint8 { cpu.B = 0xFF; return cpu.ADC_A_B() },
		},

		// ADC_A_C Tests
		{
			name:   "ADC_A_C - normal case",
			setupA: 0x25, setupReg: 0x13, regName: "C", carryFlag: false,
			expectedA: 0x38, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.C = 0x13; return cpu.ADC_A_C() },
		},
		{
			name:   "ADC_A_C - with carry flag",
			setupA: 0x25, setupReg: 0x13, regName: "C", carryFlag: true,
			expectedA: 0x39, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.C = 0x13; return cpu.ADC_A_C() },
		},

		// ADC_A_D Tests
		{
			name:   "ADC_A_D - normal case",
			setupA: 0x42, setupReg: 0x33, regName: "D", carryFlag: false,
			expectedA: 0x75, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.D = 0x33; return cpu.ADC_A_D() },
		},

		// ADC_A_E Tests
		{
			name:   "ADC_A_E - normal case",
			setupA: 0x12, setupReg: 0x34, regName: "E", carryFlag: false,
			expectedA: 0x46, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.E = 0x34; return cpu.ADC_A_E() },
		},

		// ADC_A_H Tests
		{
			name:   "ADC_A_H - normal case",
			setupA: 0x55, setupReg: 0x22, regName: "H", carryFlag: false,
			expectedA: 0x77, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.H = 0x22; return cpu.ADC_A_H() },
		},

		// ADC_A_L Tests
		{
			name:   "ADC_A_L - normal case",
			setupA: 0x66, setupReg: 0x11, regName: "L", carryFlag: false,
			expectedA: 0x77, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { cpu.L = 0x11; return cpu.ADC_A_L() },
		},

		// ADC_A_A Tests
		{
			name:   "ADC_A_A - normal case (double A)",
			setupA: 0x40, setupReg: 0x40, regName: "A", carryFlag: false,
			expectedA: 0x80, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { return cpu.ADC_A_A() },
		},
		{
			name:   "ADC_A_A - with carry (double A + 1)",
			setupA: 0x40, setupReg: 0x40, regName: "A", carryFlag: true,
			expectedA: 0x81, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { return cpu.ADC_A_A() },
		},
		{
			name:   "ADC_A_A - overflow case",
			setupA: 0x80, setupReg: 0x80, regName: "A", carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: false, expectedC: true,
			adcFunc: func(cpu *CPU) uint8 { return cpu.ADC_A_A() },
		},
		{
			name:   "ADC_A_A - overflow with carry",
			setupA: 0x7F, setupReg: 0x7F, regName: "A", carryFlag: true,
			expectedA: 0xFF, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
			adcFunc: func(cpu *CPU) uint8 { return cpu.ADC_A_A() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := tt.adcFunc(cpu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %t", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %t", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %t", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %t", tt.expectedC)
			assert.Equal(t, uint8(4), cycles, "ADC register operations should take 4 cycles")
		})
	}
}

// TestADC_A_HL tests ADC with memory operations
func TestADC_A_HL(t *testing.T) {
	tests := []struct {
		name      string
		setupA    uint8
		memValue  uint8
		carryFlag bool
		expectedA uint8
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
	}{
		{
			name:   "ADC_A_HL - normal case without carry",
			setupA: 0x30, memValue: 0x40, carryFlag: false,
			expectedA: 0x70, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_HL - normal case with carry",
			setupA: 0x30, memValue: 0x40, carryFlag: true,
			expectedA: 0x71, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_HL - half carry without initial carry",
			setupA: 0x0F, memValue: 0x01, carryFlag: false,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
		},
		{
			name:   "ADC_A_HL - half carry with initial carry",
			setupA: 0x0E, memValue: 0x01, carryFlag: true,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
		},
		{
			name:   "ADC_A_HL - carry without initial carry",
			setupA: 0xFF, memValue: 0x01, carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_HL - carry with initial carry",
			setupA: 0xFE, memValue: 0x01, carryFlag: true,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_HL - zero result with memory zero",
			setupA: 0x00, memValue: 0x00, carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_HL - carry only from flag",
			setupA: 0x00, memValue: 0x00, carryFlag: true,
			expectedA: 0x01, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_HL - maximum values",
			setupA: 0xFF, memValue: 0xFF, carryFlag: false,
			expectedA: 0xFE, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_HL - maximum values with carry",
			setupA: 0xFF, memValue: 0xFF, carryFlag: true,
			expectedA: 0xFF, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()

			cpu.A = tt.setupA
			cpu.SetFlag(FlagC, tt.carryFlag)
			cpu.SetHL(0x8000) // Set HL to valid memory address
			mmu.WriteByte(0x8000, tt.memValue)

			cycles := cpu.ADC_A_HL(mmu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %t", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %t", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %t", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %t", tt.expectedC)
			assert.Equal(t, uint8(8), cycles, "ADC memory operations should take 8 cycles")
		})
	}
}

// TestADC_A_n tests ADC with immediate values
func TestADC_A_n(t *testing.T) {
	tests := []struct {
		name      string
		setupA    uint8
		immediate uint8
		carryFlag bool
		expectedA uint8
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
	}{
		{
			name:   "ADC_A_n - normal case without carry",
			setupA: 0x50, immediate: 0x25, carryFlag: false,
			expectedA: 0x75, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_n - normal case with carry",
			setupA: 0x50, immediate: 0x25, carryFlag: true,
			expectedA: 0x76, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_n - half carry boundary",
			setupA: 0x0F, immediate: 0x01, carryFlag: false,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
		},
		{
			name:   "ADC_A_n - half carry with initial carry",
			setupA: 0x0E, immediate: 0x01, carryFlag: true,
			expectedA: 0x10, expectedZ: false, expectedN: false, expectedH: true, expectedC: false,
		},
		{
			name:   "ADC_A_n - carry boundary",
			setupA: 0xFF, immediate: 0x01, carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_n - carry boundary with initial carry",
			setupA: 0xFE, immediate: 0x01, carryFlag: true,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_n - zero result",
			setupA: 0x00, immediate: 0x00, carryFlag: false,
			expectedA: 0x00, expectedZ: true, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_n - carry makes non-zero",
			setupA: 0x00, immediate: 0x00, carryFlag: true,
			expectedA: 0x01, expectedZ: false, expectedN: false, expectedH: false, expectedC: false,
		},
		{
			name:   "ADC_A_n - maximum values",
			setupA: 0xFF, immediate: 0xFF, carryFlag: false,
			expectedA: 0xFE, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
		},
		{
			name:   "ADC_A_n - maximum values with carry",
			setupA: 0xFF, immediate: 0xFF, carryFlag: true,
			expectedA: 0xFF, expectedZ: false, expectedN: false, expectedH: true, expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			cpu.SetFlag(FlagC, tt.carryFlag)

			cycles := cpu.ADC_A_n(tt.immediate)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be %t", tt.expectedZ)
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be %t", tt.expectedN)
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be %t", tt.expectedH)
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be %t", tt.expectedC)
			assert.Equal(t, uint8(8), cycles, "ADC immediate operations should take 8 cycles")
		})
	}
}

// TestADC_EdgeCases tests edge cases and boundary conditions
func TestADC_EdgeCases(t *testing.T) {
	t.Run("ADC with all flags initially set", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x7F
		cpu.B = 0x01

		// Set all flags initially
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)

		cycles := cpu.ADC_A_B()

		// Result: 0x7F + 0x01 + 0x01 (carry) = 0x81
		assert.Equal(t, uint8(0x81), cpu.A)
		assert.Equal(t, false, cpu.GetFlag(FlagZ), "Zero flag should be reset")
		assert.Equal(t, false, cpu.GetFlag(FlagN), "Subtract flag should be reset")
		assert.Equal(t, true, cpu.GetFlag(FlagH), "Half-carry flag should be set (0xF + 0x1 + 0x1 = 0x11)")
		assert.Equal(t, false, cpu.GetFlag(FlagC), "Carry flag should be reset")
		assert.Equal(t, uint8(4), cycles)
	})

	t.Run("ADC chain operation (multi-byte addition simulation)", func(t *testing.T) {
		cpu := NewCPU()

		// Simulate adding two 16-bit numbers: 0x12FF + 0x0101 = 0x1400
		// First byte: 0xFF + 0x01 = 0x00 (carry set)
		cpu.A = 0xFF
		cpu.B = 0x01
		cpu.SetFlag(FlagC, false)
		cpu.ADD_A_B() // This should set carry flag

		firstResult := cpu.A
		carryFromFirst := cpu.GetFlag(FlagC)

		// Second byte: 0x12 + 0x01 + carry = 0x14
		cpu.A = 0x12
		cpu.C = 0x01
		// Carry flag should still be set from previous operation
		cpu.ADC_A_C() // This should use the carry from first addition

		assert.Equal(t, uint8(0x00), firstResult, "First byte should be 0x00")
		assert.Equal(t, true, carryFromFirst, "First operation should set carry")
		assert.Equal(t, uint8(0x14), cpu.A, "Second byte should be 0x14 (0x12 + 0x01 + carry)")
		assert.Equal(t, false, cpu.GetFlag(FlagC), "Final carry should be clear")
	})

	t.Run("ADC preserves other register values", func(t *testing.T) {
		cpu := NewCPU()

		// Set up all registers with known values
		cpu.A = 0x10
		cpu.B = 0x20
		cpu.C = 0x30
		cpu.D = 0x40
		cpu.E = 0x50
		cpu.H = 0x60
		cpu.L = 0x70
		cpu.SetFlag(FlagC, true)

		cpu.ADC_A_B()

		// Only A should change, others should remain
		assert.Equal(t, uint8(0x31), cpu.A, "A should be 0x10 + 0x20 + 0x01")
		assert.Equal(t, uint8(0x20), cpu.B, "B should be unchanged")
		assert.Equal(t, uint8(0x30), cpu.C, "C should be unchanged")
		assert.Equal(t, uint8(0x40), cpu.D, "D should be unchanged")
		assert.Equal(t, uint8(0x50), cpu.E, "E should be unchanged")
		assert.Equal(t, uint8(0x60), cpu.H, "H should be unchanged")
		assert.Equal(t, uint8(0x70), cpu.L, "L should be unchanged")
	})
}

// TestADC_OpcodeDispatch tests that ADC instructions work through the opcode dispatch system
func TestADC_OpcodeDispatch(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint8
		setupCPU func(*CPU)
		params   []uint8
		expected uint8
		cycles   uint8
	}{
		{
			name:     "ADC A,B (0x88)",
			opcode:   0x88,
			setupCPU: func(cpu *CPU) { cpu.A = 0x10; cpu.B = 0x20; cpu.SetFlag(FlagC, false) },
			params:   []uint8{},
			expected: 0x30,
			cycles:   4,
		},
		{
			name:     "ADC A,C (0x89)",
			opcode:   0x89,
			setupCPU: func(cpu *CPU) { cpu.A = 0x15; cpu.C = 0x25; cpu.SetFlag(FlagC, true) },
			params:   []uint8{},
			expected: 0x3B, // 0x15 + 0x25 + 0x01 = 0x3B
			cycles:   4,
		},
		{
			name:     "ADC A,D (0x8A)",
			opcode:   0x8A,
			setupCPU: func(cpu *CPU) { cpu.A = 0x30; cpu.D = 0x40; cpu.SetFlag(FlagC, false) },
			params:   []uint8{},
			expected: 0x70,
			cycles:   4,
		},
		{
			name:     "ADC A,E (0x8B)",
			opcode:   0x8B,
			setupCPU: func(cpu *CPU) { cpu.A = 0x22; cpu.E = 0x33; cpu.SetFlag(FlagC, false) },
			params:   []uint8{},
			expected: 0x55,
			cycles:   4,
		},
		{
			name:     "ADC A,H (0x8C)",
			opcode:   0x8C,
			setupCPU: func(cpu *CPU) { cpu.A = 0x11; cpu.H = 0x22; cpu.SetFlag(FlagC, true) },
			params:   []uint8{},
			expected: 0x34, // 0x11 + 0x22 + 0x01 = 0x34
			cycles:   4,
		},
		{
			name:     "ADC A,L (0x8D)",
			opcode:   0x8D,
			setupCPU: func(cpu *CPU) { cpu.A = 0x44; cpu.L = 0x33; cpu.SetFlag(FlagC, false) },
			params:   []uint8{},
			expected: 0x77,
			cycles:   4,
		},
		{
			name:     "ADC A,(HL) (0x8E)",
			opcode:   0x8E,
			setupCPU: func(cpu *CPU) { cpu.A = 0x50; cpu.SetHL(0x8000); cpu.SetFlag(FlagC, false) },
			params:   []uint8{},
			expected: 0x80, // 0x50 + 0x30 = 0x80
			cycles:   8,
		},
		{
			name:     "ADC A,A (0x8F)",
			opcode:   0x8F,
			setupCPU: func(cpu *CPU) { cpu.A = 0x40; cpu.SetFlag(FlagC, true) },
			params:   []uint8{},
			expected: 0x81, // 0x40 + 0x40 + 0x01 = 0x81
			cycles:   4,
		},
		{
			name:     "ADC A,n (0xCE)",
			opcode:   0xCE,
			setupCPU: func(cpu *CPU) { cpu.A = 0x60; cpu.SetFlag(FlagC, false) },
			params:   []uint8{0x35},
			expected: 0x95, // 0x60 + 0x35 = 0x95
			cycles:   8,
		},
		{
			name:     "ADC A,n with carry (0xCE)",
			opcode:   0xCE,
			setupCPU: func(cpu *CPU) { cpu.A = 0x60; cpu.SetFlag(FlagC, true) },
			params:   []uint8{0x35},
			expected: 0x96, // 0x60 + 0x35 + 0x01 = 0x96
			cycles:   8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()

			// Special setup for memory test
			if tt.opcode == 0x8E {
				mmu.WriteByte(0x8000, 0x30) // Set memory value for ADC A,(HL)
			}

			tt.setupCPU(cpu)

			// Verify opcode is implemented
			assert.True(t, IsOpcodeImplemented(tt.opcode), "Opcode 0x%02X should be implemented", tt.opcode)

			// Execute instruction
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			assert.NoError(t, err, "Execution should not return error")
			assert.Equal(t, tt.expected, cpu.A, "Register A should have expected value")
			assert.Equal(t, tt.cycles, cycles, "Should return correct cycle count")
			assert.Equal(t, false, cpu.GetFlag(FlagN), "N flag should always be reset for ADC")
		})
	}
}

// TestADC_HalfCarryAccuracy tests the accuracy of half-carry flag calculation
func TestADC_HalfCarryAccuracy(t *testing.T) {
	tests := []struct {
		name    string
		a, b    uint8
		carry   bool
		expectH bool
		desc    string
	}{
		{name: "No half-carry: 0x00+0x00+0", a: 0x00, b: 0x00, carry: false, expectH: false, desc: "0+0+0=0"},
		{name: "No half-carry: 0x00+0x00+1", a: 0x00, b: 0x00, carry: true, expectH: false, desc: "0+0+1=1"},
		{name: "No half-carry: 0x07+0x07+0", a: 0x07, b: 0x07, carry: false, expectH: false, desc: "7+7+0=14"},
		{name: "No half-carry: 0x06+0x07+1", a: 0x06, b: 0x07, carry: true, expectH: false, desc: "6+7+1=14"},
		{name: "Half-carry: 0x08+0x08+0", a: 0x08, b: 0x08, carry: false, expectH: true, desc: "8+8+0=16"},
		{name: "Half-carry: 0x07+0x08+1", a: 0x07, b: 0x08, carry: true, expectH: true, desc: "7+8+1=16"},
		{name: "Half-carry: 0x0F+0x01+0", a: 0x0F, b: 0x01, carry: false, expectH: true, desc: "15+1+0=16"},
		{name: "Half-carry: 0x0E+0x01+1", a: 0x0E, b: 0x01, carry: true, expectH: true, desc: "14+1+1=16"},
		{name: "Half-carry: 0x0F+0x0F+0", a: 0x0F, b: 0x0F, carry: false, expectH: true, desc: "15+15+0=30"},
		{name: "Half-carry: 0x0E+0x0F+1", a: 0x0E, b: 0x0F, carry: true, expectH: true, desc: "14+15+1=30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a
			cpu.B = tt.b
			cpu.SetFlag(FlagC, tt.carry)

			cpu.ADC_A_B()

			assert.Equal(t, tt.expectH, cpu.GetFlag(FlagH),
				"Half-carry calculation wrong for %s: (0x%02X & 0x0F) + (0x%02X & 0x0F) + %d should %s half-carry",
				tt.desc, tt.a, tt.b, map[bool]int{false: 0, true: 1}[tt.carry],
				map[bool]string{false: "not set", true: "set"}[tt.expectH])
		})
	}
}

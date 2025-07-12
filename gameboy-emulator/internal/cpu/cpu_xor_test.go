package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestXOR_A_Register tests all register-to-register XOR operations
func TestXOR_A_Register(t *testing.T) {
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
		// XOR_A_A tests
		{
			name:        "XOR_A_A - always zero (fast clear)",
			setupA:      0x42,
			setupOther:  0x42,                         // Not used for XOR_A_A
			setOtherReg: func(cpu *CPU, val uint8) {}, // No-op
			instruction: (*CPU).XOR_A_A,
			expectedA:   0x00, // A ^ A = 0 (always)
			expectedZ:   true, // Result is always zero
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "XOR_A_A - zero input",
			setupA:      0x00,
			setupOther:  0x00,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).XOR_A_A,
			expectedA:   0x00,
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:        "XOR_A_A - maximum value",
			setupA:      0xFF,
			setupOther:  0xFF,
			setOtherReg: func(cpu *CPU, val uint8) {},
			instruction: (*CPU).XOR_A_A,
			expectedA:   0x00, // 0xFF ^ 0xFF = 0x00
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_B tests
		{
			name:       "XOR_A_B - bit toggling",
			setupA:     0b11110000,
			setupOther: 0b00001111,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.B = val
			},
			instruction: (*CPU).XOR_A_B,
			expectedA:   0b11111111, // Toggle all bits
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:       "XOR_A_B - identical values",
			setupA:     0x55,
			setupOther: 0x55,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.B = val
			},
			instruction: (*CPU).XOR_A_B,
			expectedA:   0x00, // 0x55 ^ 0x55 = 0x00
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:       "XOR_A_B - selective bit toggling",
			setupA:     0b10101010,
			setupOther: 0b01010101,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.B = val
			},
			instruction: (*CPU).XOR_A_B,
			expectedA:   0b11111111, // Toggle alternating bits
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		{
			name:       "XOR_A_B - encryption pattern",
			setupA:     0x42,
			setupOther: 0xAA,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.B = val
			},
			instruction: (*CPU).XOR_A_B,
			expectedA:   0xE8, // 0x42 ^ 0xAA = 0xE8
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_C tests
		{
			name:       "XOR_A_C - nibble swap simulation",
			setupA:     0b11110000,
			setupOther: 0b00001111,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.C = val
			},
			instruction: (*CPU).XOR_A_C,
			expectedA:   0b11111111,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_D tests
		{
			name:       "XOR_A_D - checksum operation",
			setupA:     0x12,
			setupOther: 0x34,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.D = val
			},
			instruction: (*CPU).XOR_A_D,
			expectedA:   0x26, // 0x12 ^ 0x34 = 0x26
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_E tests
		{
			name:       "XOR_A_E - data manipulation",
			setupA:     0b11001100,
			setupOther: 0b00110011,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.E = val
			},
			instruction: (*CPU).XOR_A_E,
			expectedA:   0b11111111,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_H tests
		{
			name:       "XOR_A_H - address manipulation",
			setupA:     0x80,
			setupOther: 0x7F,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.H = val
			},
			instruction: (*CPU).XOR_A_H,
			expectedA:   0xFF,
			expectedZ:   false,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
		// XOR_A_L tests
		{
			name:       "XOR_A_L - low-byte manipulation",
			setupA:     0x01,
			setupOther: 0x01,
			setOtherReg: func(cpu *CPU, val uint8) {
				cpu.L = val
			},
			instruction: (*CPU).XOR_A_L,
			expectedA:   0x00, // Same bits cancel out
			expectedZ:   true,
			expectedN:   false,
			expectedH:   false,
			expectedC:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA
			tt.setOtherReg(cpu, tt.setupOther)

			cycles := tt.instruction(cpu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be set correctly")
			assert.Equal(t, tt.expectedN, cpu.GetFlag(FlagN), "Subtract flag should be reset")
			assert.Equal(t, tt.expectedH, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should be reset")
		})
	}
}

// TestXOR_A_HL tests XOR operation with memory
func TestXOR_A_HL(t *testing.T) {
	tests := []struct {
		name        string
		setupA      uint8
		memoryValue uint8
		hlAddress   uint16
		expectedA   uint8
		expectedZ   bool
	}{
		{
			name:        "XOR_A_HL - memory encryption",
			setupA:      0x42,
			memoryValue: 0xAA,
			hlAddress:   0x8000,
			expectedA:   0xE8, // 0x42 ^ 0xAA = 0xE8
			expectedZ:   false,
		},
		{
			name:        "XOR_A_HL - decrypt same key",
			setupA:      0xE8,
			memoryValue: 0xAA,
			hlAddress:   0x8000,
			expectedA:   0x42, // 0xE8 ^ 0xAA = 0x42 (decryption)
			expectedZ:   false,
		},
		{
			name:        "XOR_A_HL - zero result from memory",
			setupA:      0x55,
			memoryValue: 0x55,
			hlAddress:   0x9000,
			expectedA:   0x00, // Same values = zero
			expectedZ:   true,
		},
		{
			name:        "XOR_A_HL - bit flip pattern",
			setupA:      0b11110000,
			memoryValue: 0b00001111,
			hlAddress:   0xC000,
			expectedA:   0b11111111, // All bits flipped
			expectedZ:   false,
		},
		{
			name:        "XOR_A_HL - toggle single bit",
			setupA:      0b10000000,
			memoryValue: 0b00000001,
			hlAddress:   0xFF80,
			expectedA:   0b10000001, // Toggle bit 0
			expectedZ:   false,
		},
		{
			name:        "XOR_A_HL - maximum values",
			setupA:      0xFF,
			memoryValue: 0xFF,
			hlAddress:   0xFFFF,
			expectedA:   0x00, // 0xFF ^ 0xFF = 0x00
			expectedZ:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()

			cpu.A = tt.setupA
			cpu.SetHL(tt.hlAddress)
			mmu.WriteByte(tt.hlAddress, tt.memoryValue)

			cycles := cpu.XOR_A_HL(mmu)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be set correctly")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
			assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
		})
	}
}

// TestXOR_A_n tests XOR operation with immediate values
func TestXOR_A_n(t *testing.T) {
	tests := []struct {
		name         string
		setupA       uint8
		immediateVal uint8
		expectedA    uint8
		expectedZ    bool
	}{
		{
			name:         "XOR_A_n - bit toggling with constant",
			setupA:       0b10101010,
			immediateVal: 0b01010101,
			expectedA:    0b11111111, // Toggle alternating bits
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - encrypt with key",
			setupA:       0x42,
			immediateVal: 0xAA,
			expectedA:    0xE8, // Simple encryption
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - invert all bits",
			setupA:       0b10101010,
			immediateVal: 0xFF,
			expectedA:    0b01010101, // Invert all bits
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - toggle high bit",
			setupA:       0x00,
			immediateVal: 0x80,
			expectedA:    0x80, // Set bit 7
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - toggle low bit",
			setupA:       0x00,
			immediateVal: 0x01,
			expectedA:    0x01, // Set bit 0
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - clear with same value",
			setupA:       0x42,
			immediateVal: 0x42,
			expectedA:    0x00, // Clear A register
			expectedZ:    true,
		},
		{
			name:         "XOR_A_n - no change with zero",
			setupA:       0x55,
			immediateVal: 0x00,
			expectedA:    0x55, // No change
			expectedZ:    false,
		},
		{
			name:         "XOR_A_n - zero result from zero A",
			setupA:       0x00,
			immediateVal: 0x00,
			expectedA:    0x00,
			expectedZ:    true,
		},
		{
			name:         "XOR_A_n - checksum calculation",
			setupA:       0x12,
			immediateVal: 0x34,
			expectedA:    0x26, // 0x12 ^ 0x34 = 0x26
			expectedZ:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.setupA

			cycles := cpu.XOR_A_n(tt.immediateVal)

			assert.Equal(t, tt.expectedA, cpu.A, "Register A should have expected value")
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should be set correctly")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
			assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
		})
	}
}

// TestXOR_FlagBehavior tests specific flag behavior for XOR operations
func TestXOR_FlagBehavior(t *testing.T) {
	t.Run("XOR always resets N, H, and C flags", func(t *testing.T) {
		cpu := NewCPU()
		// Set all flags first
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)

		// Any XOR operation should reset N, H, C
		cpu.A = 0x42
		cpu.B = 0x55
		cpu.XOR_A_B()

		assert.False(t, cpu.GetFlag(FlagN), "N flag should always be reset for XOR")
		assert.False(t, cpu.GetFlag(FlagH), "H flag should always be reset for XOR")
		assert.False(t, cpu.GetFlag(FlagC), "C flag should always be reset for XOR")
	})

	t.Run("XOR zero flag accuracy", func(t *testing.T) {
		zeroTests := []struct {
			name      string
			valueA    uint8
			valueB    uint8
			shouldBeZ bool
		}{
			{"0x00 XOR 0x00 should be zero", 0x00, 0x00, true},
			{"0xFF XOR 0xFF should be zero", 0xFF, 0xFF, true},
			{"0xAA XOR 0xAA should be zero (identical values)", 0xAA, 0xAA, true},
			{"0x55 XOR 0x55 should be zero (identical values)", 0x55, 0x55, true},
			{"0xFF XOR 0x00 should not be zero", 0xFF, 0x00, false},
			{"0x01 XOR 0x02 should not be zero", 0x01, 0x02, false},
			{"0x80 XOR 0x01 should not be zero", 0x80, 0x01, false},
			{"0x0F XOR 0xF0 should not be zero", 0x0F, 0xF0, false},
		}

		for _, tt := range zeroTests {
			t.Run(tt.name, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tt.valueA
				cpu.B = tt.valueB
				cpu.XOR_A_B()
				assert.Equal(t, tt.shouldBeZ, cpu.GetFlag(FlagZ), "Zero flag should be set correctly")
			})
		}
	})
}

// TestXOR_EdgeCases tests edge cases and patterns for XOR operations
func TestXOR_EdgeCases(t *testing.T) {
	t.Run("XOR with maximum values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF
		cpu.B = 0xFF
		cpu.XOR_A_B()
		assert.Equal(t, uint8(0x00), cpu.A, "0xFF XOR 0xFF should be 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
	})

	t.Run("XOR with zero values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		cpu.B = 0x00
		cpu.XOR_A_B()
		assert.Equal(t, uint8(0x00), cpu.A, "0x00 XOR 0x00 should be 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
	})

	t.Run("XOR bit manipulation patterns", func(t *testing.T) {
		bitTests := []struct {
			name     string
			valueA   uint8
			valueB   uint8
			expected uint8
		}{
			{"Toggle all bits", 0b10101010, 0b11111111, 0b01010101},
			{"Toggle no bits", 0b10101010, 0b00000000, 0b10101010},
			{"Toggle alternating bits", 0b11110000, 0b10101010, 0b01011010},
			{"Single bit toggle", 0b00000000, 0b10000000, 0b10000000},
			{"Clear specific bits", 0b11111111, 0b11111111, 0b00000000},
			{"Complement operation", 0b00001111, 0b11110000, 0b11111111},
		}

		for _, tt := range bitTests {
			t.Run(tt.name, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tt.valueA
				cpu.B = tt.valueB
				cpu.XOR_A_B()
				assert.Equal(t, tt.expected, cpu.A, "XOR result should match expected pattern")
			})
		}
	})
}

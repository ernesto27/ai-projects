package cpu

import (
	"testing"
	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

// === SRA Instruction Tests ===
// SRA shifts right arithmetic - preserves sign bit for signed numbers

func TestSRA_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		expected    uint8
		expectedZ   bool
		expectedC   bool
		description string
	}{
		// Test positive numbers (sign bit = 0)
		{"SRA_B_positive", (*CPU).SRA_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x7E, 0x3F, false, false, "0111 1110 -> 0011 1111"},
		{"SRA_C_positive", (*CPU).SRA_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x7E, 0x3F, false, false, "positive number"},
		{"SRA_D_positive", (*CPU).SRA_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x7E, 0x3F, false, false, "positive number"},
		{"SRA_E_positive", (*CPU).SRA_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x7E, 0x3F, false, false, "positive number"},
		{"SRA_H_positive", (*CPU).SRA_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x7E, 0x3F, false, false, "positive number"},
		{"SRA_L_positive", (*CPU).SRA_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x7E, 0x3F, false, false, "positive number"},
		{"SRA_A_positive", (*CPU).SRA_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x7E, 0x3F, false, false, "positive number"},

		// Test negative numbers (sign bit = 1) - sign bit should be preserved
		{"SRA_B_negative", (*CPU).SRA_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x84, 0xC2, false, false, "1000 0100 -> 1100 0010"},
		{"SRA_C_negative", (*CPU).SRA_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x84, 0xC2, false, false, "negative number"},
		{"SRA_D_negative", (*CPU).SRA_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0xFF, 0xFF, false, true, "all 1s stays all 1s"},
		{"SRA_E_negative", (*CPU).SRA_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x85, 0xC2, false, true, "1000 0101 -> 1100 0010"},

		// Test carry flag cases
		{"SRA_H_carry", (*CPU).SRA_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x85, 0xC2, false, true, "odd number sets carry"},
		{"SRA_L_no_carry", (*CPU).SRA_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x84, 0xC2, false, false, "even number no carry"},

		// Test zero result
		{"SRA_A_zero", (*CPU).SRA_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x01, 0x00, true, true, "0000 0001 -> 0000 0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "SRA should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected: %s", tt.description)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestSRA_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.SetHL(0x8000)
	
	// Test positive number in memory
	mmu.WriteByte(0x8000, 0x7E) // 01111110
	cycles := cpu.SRA_HL(mmu)
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x3F), result, "Positive: 01111110 -> 00111111")
	assert.Equal(t, uint8(16), cycles, "SRA (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be false (even number)")
	
	// Test negative number in memory
	mmu.WriteByte(0x8000, 0x85) // 10000101
	cycles = cpu.SRA_HL(mmu)
	result = mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0xC2), result, "Negative: 10000101 -> 11000010 (sign preserved)")
	assert.Equal(t, uint8(16), cycles, "SRA (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (odd number)")
}

// === SRL Instruction Tests ===
// SRL shifts right logical - always fills with zero (for unsigned numbers)

func TestSRL_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		expected    uint8
		expectedZ   bool
		expectedC   bool
		description string
	}{
		// Test logical shift (always zero fill)
		{"SRL_B_logical", (*CPU).SRL_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x84, 0x42, false, false, "1000 0100 -> 0100 0010"},
		{"SRL_C_logical", (*CPU).SRL_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x84, 0x42, false, false, "logical shift"},
		{"SRL_D_logical", (*CPU).SRL_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x84, 0x42, false, false, "logical shift"},
		{"SRL_E_logical", (*CPU).SRL_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x84, 0x42, false, false, "logical shift"},
		{"SRL_H_logical", (*CPU).SRL_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x84, 0x42, false, false, "logical shift"},
		{"SRL_L_logical", (*CPU).SRL_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x84, 0x42, false, false, "logical shift"},
		{"SRL_A_logical", (*CPU).SRL_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x84, 0x42, false, false, "logical shift"},

		// Test with negative numbers (should NOT preserve sign)
		{"SRL_B_negative", (*CPU).SRL_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0xFF, 0x7F, false, true, "1111 1111 -> 0111 1111"},
		{"SRL_C_negative", (*CPU).SRL_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x85, 0x42, false, true, "1000 0101 -> 0100 0010"},

		// Test carry flag cases
		{"SRL_D_carry", (*CPU).SRL_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x85, 0x42, false, true, "odd number sets carry"},
		{"SRL_E_no_carry", (*CPU).SRL_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x84, 0x42, false, false, "even number no carry"},

		// Test zero result
		{"SRL_H_zero", (*CPU).SRL_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x01, 0x00, true, true, "0000 0001 -> 0000 0000"},
		{"SRL_L_zero_input", (*CPU).SRL_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x00, 0x00, true, false, "0000 0000 -> 0000 0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "SRL should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected: %s", tt.description)
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestSRL_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.SetHL(0x8000)
	
	// Test logical shift in memory
	mmu.WriteByte(0x8000, 0xFF) // 11111111
	cycles := cpu.SRL_HL(mmu)
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x7F), result, "Logical: 11111111 -> 01111111 (zero fill)")
	assert.Equal(t, uint8(16), cycles, "SRL (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (odd number)")
	
	// Test zero result in memory
	mmu.WriteByte(0x8000, 0x01) // 00000001
	cycles = cpu.SRL_HL(mmu)
	result = mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x00), result, "Should become zero")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be true")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (bit 0 was 1)")
}

// === SRA vs SRL Comparison Tests ===

func TestSRA_vs_SRL_Difference(t *testing.T) {
	// Test with negative number to show the difference
	input := uint8(0x85) // 10000101 (negative if signed)
	
	// SRA should preserve sign bit
	cpu1 := NewCPU()
	cpu1.A = input
	cpu1.SRA_A()
	sraResult := cpu1.A
	assert.Equal(t, uint8(0xC2), sraResult, "SRA: 10000101 -> 11000010 (sign preserved)")
	
	// SRL should zero-fill
	cpu2 := NewCPU()
	cpu2.A = input
	cpu2.SRL_A()
	srlResult := cpu2.A
	assert.Equal(t, uint8(0x42), srlResult, "SRL: 10000101 -> 01000010 (zero fill)")
	
	// Both should set carry flag (bit 0 was 1)
	assert.True(t, cpu1.GetFlag(FlagC), "SRA should set carry")
	assert.True(t, cpu2.GetFlag(FlagC), "SRL should set carry")
}

// === CB Dispatch Integration Tests ===

func TestCB_Shift_Dispatch(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint8
		setup    func(*CPU, *memory.MMU)
		expected func(*CPU, *memory.MMU) bool
	}{
		{
			"CB_SRA_B",
			0x28,
			func(cpu *CPU, mmu *memory.MMU) { cpu.B = 0x85 },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.B == 0xC2 && cpu.GetFlag(FlagC) },
		},
		{
			"CB_SRA_C", 
			0x29,
			func(cpu *CPU, mmu *memory.MMU) { cpu.C = 0x7E },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.C == 0x3F && !cpu.GetFlag(FlagC) },
		},
		{
			"CB_SRA_HL",
			0x2E,
			func(cpu *CPU, mmu *memory.MMU) { cpu.SetHL(0x8000); mmu.WriteByte(0x8000, 0x85) },
			func(cpu *CPU, mmu *memory.MMU) bool { return mmu.ReadByte(0x8000) == 0xC2 && cpu.GetFlag(FlagC) },
		},
		{
			"CB_SRL_B",
			0x38,
			func(cpu *CPU, mmu *memory.MMU) { cpu.B = 0x85 },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.B == 0x42 && cpu.GetFlag(FlagC) },
		},
		{
			"CB_SRL_A",
			0x3F,
			func(cpu *CPU, mmu *memory.MMU) { cpu.A = 0xFF },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.A == 0x7F && cpu.GetFlag(FlagC) },
		},
		{
			"CB_SRL_HL",
			0x3E,
			func(cpu *CPU, mmu *memory.MMU) { cpu.SetHL(0x8000); mmu.WriteByte(0x8000, 0xFF) },
			func(cpu *CPU, mmu *memory.MMU) bool { return mmu.ReadByte(0x8000) == 0x7F && cpu.GetFlag(FlagC) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()
			
			tt.setup(cpu, mmu)
			
			cycles, err := cpu.ExecuteCBInstruction(mmu, tt.opcode)
			
			assert.NoError(t, err, "CB instruction should execute without error")
			assert.True(t, cycles == 8 || cycles == 16, "Cycles should be 8 or 16")
			assert.True(t, tt.expected(cpu, mmu), "Instruction should produce expected result")
		})
	}
}

// === Edge Case Tests ===

func TestShift_EdgeCases(t *testing.T) {
	t.Run("SRA_all_zeros", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		
		cycles := cpu.SRA_A()
		
		assert.Equal(t, uint8(0x00), cpu.A, "0x00 should stay 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("SRL_all_zeros", func(t *testing.T) {
		cpu := NewCPU()
		cpu.B = 0x00
		
		cycles := cpu.SRL_B()
		
		assert.Equal(t, uint8(0x00), cpu.B, "0x00 should stay 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("SRA_single_bit", func(t *testing.T) {
		cpu := NewCPU()
		cpu.C = 0x80 // 10000000 (most negative 8-bit signed number)
		
		cycles := cpu.SRA_C()
		
		assert.Equal(t, uint8(0xC0), cpu.C, "10000000 -> 11000000 (sign preserved)")
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear (bit 0 was 0)")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("SRL_single_bit", func(t *testing.T) {
		cpu := NewCPU()
		cpu.D = 0x80 // 10000000
		
		cycles := cpu.SRL_D()
		
		assert.Equal(t, uint8(0x40), cpu.D, "10000000 -> 01000000 (zero fill)")
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear (bit 0 was 0)")
		assert.Equal(t, uint8(8), cycles)
	})
}

// === Performance Tests ===

func BenchmarkSRA_Register(b *testing.B) {
	cpu := NewCPU()
	cpu.A = 0x85
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpu.A = 0x85
		cpu.SRA_A()
	}
}

func BenchmarkSRL_Register(b *testing.B) {
	cpu := NewCPU()
	cpu.B = 0x85
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpu.B = 0x85
		cpu.SRL_B()
	}
}

func BenchmarkSRA_Memory(b *testing.B) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	cpu.SetHL(0x8000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mmu.WriteByte(0x8000, 0x85)
		cpu.SRA_HL(mmu)
	}
}
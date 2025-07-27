package cpu

import (
	"testing"
	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

// === RLC Instruction Tests ===
// RLC rotates left circular - bit 7 wraps to bit 0

func TestRLC_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		expected    uint8
		expectedZ   bool
		expectedC   bool
	}{
		{"RLC_D", (*CPU).RLC_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x85, 0x0B, false, true},
		{"RLC_E", (*CPU).RLC_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x85, 0x0B, false, true},
		{"RLC_H", (*CPU).RLC_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x85, 0x0B, false, true},
		{"RLC_L", (*CPU).RLC_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x85, 0x0B, false, true},
		{"RLC_A", (*CPU).RLC_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x85, 0x0B, false, true},
		
		// Test zero result
		{"RLC_D_zero", (*CPU).RLC_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x00, 0x00, true, false},
		
		// Test no carry
		{"RLC_E_no_carry", (*CPU).RLC_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x42, 0x84, false, false},
		
		// Test all bits set
		{"RLC_H_all_bits", (*CPU).RLC_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0xFF, 0xFF, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "RLC should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestRLC_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()
	
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x85) // 10000101
	
	cycles := cpu.RLC_HL(mmu)
	
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x0B), result, "Memory value should be rotated: 10000101 -> 00001011")
	assert.Equal(t, uint8(16), cycles, "RLC (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (bit 7 was 1)")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be false")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be false")
}

// === RRC Instruction Tests ===
// RRC rotates right circular - bit 0 wraps to bit 7

func TestRRC_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		expected    uint8
		expectedZ   bool
		expectedC   bool
	}{
		{"RRC_D", (*CPU).RRC_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x85, 0xC2, false, true},
		{"RRC_E", (*CPU).RRC_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x85, 0xC2, false, true},
		{"RRC_H", (*CPU).RRC_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x85, 0xC2, false, true},
		{"RRC_L", (*CPU).RRC_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x85, 0xC2, false, true},
		{"RRC_A", (*CPU).RRC_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x85, 0xC2, false, true},
		
		// Test zero result
		{"RRC_D_zero", (*CPU).RRC_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x00, 0x00, true, false},
		
		// Test no carry (even number)
		{"RRC_E_no_carry", (*CPU).RRC_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x42, 0x21, false, false},
		
		// Test single bit
		{"RRC_H_single_bit", (*CPU).RRC_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x01, 0x80, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "RRC should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestRRC_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()
	
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x85) // 10000101
	
	cycles := cpu.RRC_HL(mmu)
	
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0xC2), result, "Memory value should be rotated: 10000101 -> 11000010")
	assert.Equal(t, uint8(16), cycles, "RRC (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (bit 0 was 1)")
}

// === RL Instruction Tests ===
// RL rotates left through carry - old carry becomes bit 0

func TestRL_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		initialCarry bool
		expected    uint8
		expectedZ   bool
		expectedC   bool
	}{
		// Test with carry=0
		{"RL_B_carry0", (*CPU).RL_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x85, false, 0x0A, false, true},
		{"RL_C_carry0", (*CPU).RL_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x85, false, 0x0A, false, true},
		{"RL_D_carry0", (*CPU).RL_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x85, false, 0x0A, false, true},
		{"RL_E_carry0", (*CPU).RL_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x85, false, 0x0A, false, true},
		{"RL_H_carry0", (*CPU).RL_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x85, false, 0x0A, false, true},
		{"RL_L_carry0", (*CPU).RL_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x85, false, 0x0A, false, true},
		{"RL_A_carry0", (*CPU).RL_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x85, false, 0x0A, false, true},
		
		// Test with carry=1
		{"RL_B_carry1", (*CPU).RL_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x85, true, 0x0B, false, true},
		{"RL_C_carry1", (*CPU).RL_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x42, true, 0x85, false, false},
		
		// Test zero result  
		{"RL_D_zero", (*CPU).RL_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x80, false, 0x00, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetFlag(FlagC, tt.initialCarry)
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "RL should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestRL_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()
	
	cpu.SetHL(0x8000)
	cpu.SetFlag(FlagC, true) // Initial carry = 1
	mmu.WriteByte(0x8000, 0x85) // 10000101
	
	cycles := cpu.RL_HL(mmu)
	
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x0B), result, "Memory should be rotated left with carry: 10000101 -> 00001011")
	assert.Equal(t, uint8(16), cycles, "RL (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (bit 7 was 1)")
}

// === RR Instruction Tests ===
// RR rotates right through carry - old carry becomes bit 7

func TestRR_Instructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction func(*CPU) uint8
		setValue    func(*CPU, uint8)
		getValue    func(*CPU) uint8
		input       uint8
		initialCarry bool
		expected    uint8
		expectedZ   bool
		expectedC   bool
	}{
		// Test with carry=0
		{"RR_B_carry0", (*CPU).RR_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x85, false, 0x42, false, true},
		{"RR_C_carry0", (*CPU).RR_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x85, false, 0x42, false, true},
		{"RR_D_carry0", (*CPU).RR_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x85, false, 0x42, false, true},
		{"RR_E_carry0", (*CPU).RR_E, func(cpu *CPU, v uint8) { cpu.E = v }, func(cpu *CPU) uint8 { return cpu.E }, 0x85, false, 0x42, false, true},
		{"RR_H_carry0", (*CPU).RR_H, func(cpu *CPU, v uint8) { cpu.H = v }, func(cpu *CPU) uint8 { return cpu.H }, 0x85, false, 0x42, false, true},
		{"RR_L_carry0", (*CPU).RR_L, func(cpu *CPU, v uint8) { cpu.L = v }, func(cpu *CPU) uint8 { return cpu.L }, 0x85, false, 0x42, false, true},
		{"RR_A_carry0", (*CPU).RR_A, func(cpu *CPU, v uint8) { cpu.A = v }, func(cpu *CPU) uint8 { return cpu.A }, 0x85, false, 0x42, false, true},
		
		// Test with carry=1
		{"RR_B_carry1", (*CPU).RR_B, func(cpu *CPU, v uint8) { cpu.B = v }, func(cpu *CPU) uint8 { return cpu.B }, 0x85, true, 0xC2, false, true},
		{"RR_C_carry1", (*CPU).RR_C, func(cpu *CPU, v uint8) { cpu.C = v }, func(cpu *CPU) uint8 { return cpu.C }, 0x42, true, 0xA1, false, false},
		
		// Test zero result
		{"RR_D_zero", (*CPU).RR_D, func(cpu *CPU, v uint8) { cpu.D = v }, func(cpu *CPU) uint8 { return cpu.D }, 0x01, false, 0x00, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetFlag(FlagC, tt.initialCarry)
			tt.setValue(cpu, tt.input)
			
			cycles := tt.instruction(cpu)
			
			assert.Equal(t, uint8(8), cycles, "RR should take 8 cycles")
			assert.Equal(t, tt.expected, tt.getValue(cpu), "Register value should match expected")
			assert.Equal(t, tt.expectedZ, cpu.GetFlag(FlagZ), "Zero flag should match expected")
			assert.Equal(t, tt.expectedC, cpu.GetFlag(FlagC), "Carry flag should match expected")
			assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should always be false")
			assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should always be false")
		})
	}
}

func TestRR_HL_Memory(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()
	
	cpu.SetHL(0x8000)
	cpu.SetFlag(FlagC, true) // Initial carry = 1
	mmu.WriteByte(0x8000, 0x85) // 10000101
	
	cycles := cpu.RR_HL(mmu)
	
	result := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0xC2), result, "Memory should be rotated right with carry: 10000101 -> 11000010")
	assert.Equal(t, uint8(16), cycles, "RR (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be true (bit 0 was 1)")
}

// === CB Dispatch Integration Tests ===

func TestCB_Rotation_Dispatch(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint8
		setup    func(*CPU, *memory.MMU)
		expected func(*CPU, *memory.MMU) bool
	}{
		{
			"CB_RLC_D",
			0x02,
			func(cpu *CPU, mmu *memory.MMU) { cpu.D = 0x85 },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.D == 0x0B && cpu.GetFlag(FlagC) },
		},
		{
			"CB_RRC_E", 
			0x0B,
			func(cpu *CPU, mmu *memory.MMU) { cpu.E = 0x85 },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.E == 0xC2 && cpu.GetFlag(FlagC) },
		},
		{
			"CB_RL_H",
			0x14,
			func(cpu *CPU, mmu *memory.MMU) { cpu.H = 0x85; cpu.SetFlag(FlagC, true) },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.H == 0x0B && cpu.GetFlag(FlagC) },
		},
		{
			"CB_RR_L",
			0x1D,
			func(cpu *CPU, mmu *memory.MMU) { cpu.L = 0x85; cpu.SetFlag(FlagC, true) },
			func(cpu *CPU, mmu *memory.MMU) bool { return cpu.L == 0xC2 && cpu.GetFlag(FlagC) },
		},
		{
			"CB_RLC_HL",
			0x06,
			func(cpu *CPU, mmu *memory.MMU) { cpu.SetHL(0x8000); mmu.WriteByte(0x8000, 0x85) },
			func(cpu *CPU, mmu *memory.MMU) bool { return mmu.ReadByte(0x8000) == 0x0B && cpu.GetFlag(FlagC) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			
			tt.setup(cpu, mmu)
			
			cycles, err := cpu.ExecuteCBInstruction(mmu, tt.opcode)
			
			assert.NoError(t, err, "CB instruction should execute without error")
			assert.True(t, cycles == 8 || cycles == 16, "Cycles should be 8 or 16")
			assert.True(t, tt.expected(cpu, mmu), "Instruction should produce expected result")
		})
	}
}

// === Edge Case Tests ===

func TestRotation_EdgeCases(t *testing.T) {
	t.Run("RLC_all_zeros", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		
		cycles := cpu.RLC_A()
		
		assert.Equal(t, uint8(0x00), cpu.A, "0x00 should stay 0x00")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("RLC_all_ones", func(t *testing.T) {
		cpu := NewCPU()
		cpu.B = 0xFF
		
		cycles := cpu.RLC_B()
		
		assert.Equal(t, uint8(0xFF), cpu.B, "0xFF should stay 0xFF")
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
		assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("RL_alternating_pattern", func(t *testing.T) {
		cpu := NewCPU()
		cpu.C = 0xAA // 10101010
		cpu.SetFlag(FlagC, false)
		
		cycles := cpu.RL_C()
		
		assert.Equal(t, uint8(0x54), cpu.C, "10101010 -> 01010100") // 0xAA -> 0x54
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
		assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set (bit 7 was 1)")
		assert.Equal(t, uint8(8), cycles)
	})
	
	t.Run("RR_single_bit_rotation", func(t *testing.T) {
		cpu := NewCPU()
		cpu.D = 0x01 // 00000001
		cpu.SetFlag(FlagC, false)
		
		cycles := cpu.RR_D()
		
		assert.Equal(t, uint8(0x00), cpu.D, "00000001 -> 00000000")
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
		assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set (bit 0 was 1)")
		assert.Equal(t, uint8(8), cycles)
	})
}

// === Performance Tests ===

func BenchmarkRLC_Register(b *testing.B) {
	cpu := NewCPU()
	cpu.A = 0x85
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpu.A = 0x85
		cpu.RLC_A()
	}
}

func BenchmarkRL_Register(b *testing.B) {
	cpu := NewCPU()
	cpu.B = 0x85
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpu.B = 0x85
		cpu.SetFlag(FlagC, true)
		cpu.RL_B()
	}
}

func BenchmarkRLC_Memory(b *testing.B) {
	cpu := NewCPU()
	mmu := createTestMMU()
	cpu.SetHL(0x8000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mmu.WriteByte(0x8000, 0x85)
		cpu.RLC_HL(mmu)
	}
}
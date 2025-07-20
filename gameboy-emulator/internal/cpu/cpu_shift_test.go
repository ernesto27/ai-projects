package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === SLA (Shift Left Arithmetic) Tests ===
// SLA shifts all bits left by 1, bit 7 goes to carry, bit 0 becomes 0
// Flags: Z = result==0, N = 0, H = 0, C = original bit 7

func TestSLA_B_BasicOperation(t *testing.T) {
	cpu := NewCPU()
	
	// Test normal case: 0b01010101 -> 0b10101010
	cpu.B = 0x55 // 0b01010101
	cycles := cpu.SLA_B()
	
	assert.Equal(t, uint8(0xAA), cpu.B, "SLA B: 0x55 should become 0xAA")
	assert.Equal(t, uint8(8), cycles, "SLA B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be clear (result not zero)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should always be clear for SLA")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should always be clear for SLA") 
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be clear (bit 7 was 0)")
}

func TestSLA_B_CarryFlag(t *testing.T) {
	cpu := NewCPU()
	
	// Test carry case: 0b10110101 -> 0b01101010, C=1
	cpu.B = 0xB5 // 0b10110101 (bit 7 = 1)
	cycles := cpu.SLA_B()
	
	assert.Equal(t, uint8(0x6A), cpu.B, "SLA B: 0xB5 should become 0x6A")
	assert.Equal(t, uint8(8), cycles, "SLA B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be clear")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be set (original bit 7 was 1)")
}

func TestSLA_B_ZeroResult(t *testing.T) {
	cpu := NewCPU()
	
	// Test zero result: 0b10000000 -> 0b00000000, C=1, Z=1
	cpu.B = 0x80 // 0b10000000
	cycles := cpu.SLA_B()
	
	assert.Equal(t, uint8(0x00), cpu.B, "SLA B: 0x80 should become 0x00")
	assert.Equal(t, uint8(8), cycles, "SLA B should take 8 cycles")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be set (result is zero)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be clear")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be set (original bit 7 was 1)")
}

func TestSLA_AllRegisters(t *testing.T) {
	testCases := []struct {
		name     string
		setupReg func(*CPU, uint8)
		getReg   func(*CPU) uint8
		slaFunc  func(*CPU) uint8
	}{
		{"SLA_B", func(cpu *CPU, val uint8) { cpu.B = val }, func(cpu *CPU) uint8 { return cpu.B }, (*CPU).SLA_B},
		{"SLA_C", func(cpu *CPU, val uint8) { cpu.C = val }, func(cpu *CPU) uint8 { return cpu.C }, (*CPU).SLA_C},
		{"SLA_D", func(cpu *CPU, val uint8) { cpu.D = val }, func(cpu *CPU) uint8 { return cpu.D }, (*CPU).SLA_D},
		{"SLA_E", func(cpu *CPU, val uint8) { cpu.E = val }, func(cpu *CPU) uint8 { return cpu.E }, (*CPU).SLA_E},
		{"SLA_H", func(cpu *CPU, val uint8) { cpu.H = val }, func(cpu *CPU) uint8 { return cpu.H }, (*CPU).SLA_H},
		{"SLA_L", func(cpu *CPU, val uint8) { cpu.L = val }, func(cpu *CPU) uint8 { return cpu.L }, (*CPU).SLA_L},
		{"SLA_A", func(cpu *CPU, val uint8) { cpu.A = val }, func(cpu *CPU) uint8 { return cpu.A }, (*CPU).SLA_A},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := NewCPU()
			
			// Test with 0x42 (0b01000010) -> 0x84 (0b10000100), no carry
			tc.setupReg(cpu, 0x42)
			cycles := tc.slaFunc(cpu)
			
			assert.Equal(t, uint8(0x84), tc.getReg(cpu), "%s: 0x42 should become 0x84", tc.name)
			assert.Equal(t, uint8(8), cycles, "%s should take 8 cycles", tc.name)
			assert.False(t, cpu.GetFlag(FlagZ), "%s: Z flag should be clear", tc.name)
			assert.False(t, cpu.GetFlag(FlagC), "%s: C flag should be clear (bit 7 was 0)", tc.name)
		})
	}
}

func TestSLA_HL_MemoryOperation(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set up HL to point to memory address 0x8000
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x33) // 0b00110011
	
	cycles := cpu.SLA_HL(mmu)
	result := mmu.ReadByte(0x8000)
	
	assert.Equal(t, uint8(0x66), result, "SLA (HL): 0x33 should become 0x66")
	assert.Equal(t, uint8(16), cycles, "SLA (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be clear")
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be clear (bit 7 was 0)")
}

func TestSLA_HL_MemoryWithCarry(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set up HL and test value with bit 7 set
	cpu.SetHL(0x9000)
	mmu.WriteByte(0x9000, 0xFF) // 0b11111111
	
	cycles := cpu.SLA_HL(mmu)
	result := mmu.ReadByte(0x9000)
	
	assert.Equal(t, uint8(0xFE), result, "SLA (HL): 0xFF should become 0xFE")
	assert.Equal(t, uint8(16), cycles, "SLA (HL) should take 16 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be clear")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be set (original bit 7 was 1)")
}

func TestSLA_EdgeCases(t *testing.T) {
	t.Run("Zero input", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		
		cycles := cpu.SLA_A()
		
		assert.Equal(t, uint8(0x00), cpu.A, "SLA A: 0x00 should remain 0x00")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be set (result is zero)")
		assert.False(t, cpu.GetFlag(FlagC), "C flag should be clear (bit 7 was 0)")
	})

	t.Run("Maximum value", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF // 0b11111111
		
		cycles := cpu.SLA_A()
		
		assert.Equal(t, uint8(0xFE), cpu.A, "SLA A: 0xFF should become 0xFE")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be clear")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be set")
	})

	t.Run("Flag preservation", func(t *testing.T) {
		cpu := NewCPU()
		// Set some initial flags
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		
		cpu.A = 0x01
		cpu.SLA_A()
		
		assert.False(t, cpu.GetFlag(FlagN), "N flag should always be cleared by SLA")
		assert.False(t, cpu.GetFlag(FlagH), "H flag should always be cleared by SLA")
	})
}

// === CB Dispatch Integration Tests ===

func TestSLA_CBDispatchIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	testCases := []struct {
		opcode   uint8
		name     string
		setupReg func(*CPU)
		checkReg func(*CPU) uint8
		expected uint8
	}{
		{0x20, "SLA B", func(cpu *CPU) { cpu.B = 0x11 }, func(cpu *CPU) uint8 { return cpu.B }, 0x22},
		{0x21, "SLA C", func(cpu *CPU) { cpu.C = 0x11 }, func(cpu *CPU) uint8 { return cpu.C }, 0x22},
		{0x22, "SLA D", func(cpu *CPU) { cpu.D = 0x11 }, func(cpu *CPU) uint8 { return cpu.D }, 0x22},
		{0x23, "SLA E", func(cpu *CPU) { cpu.E = 0x11 }, func(cpu *CPU) uint8 { return cpu.E }, 0x22},
		{0x24, "SLA H", func(cpu *CPU) { cpu.H = 0x11 }, func(cpu *CPU) uint8 { return cpu.H }, 0x22},
		{0x25, "SLA L", func(cpu *CPU) { cpu.L = 0x11 }, func(cpu *CPU) uint8 { return cpu.L }, 0x22},
		{0x27, "SLA A", func(cpu *CPU) { cpu.A = 0x11 }, func(cpu *CPU) uint8 { return cpu.A }, 0x22},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.Reset()
			tc.setupReg(cpu)
			
			cycles, err := cpu.ExecuteCBInstruction(mmu, tc.opcode)
			
			assert.NoError(t, err, "CB instruction should execute without error")
			assert.Equal(t, uint8(8), cycles, "Should return 8 cycles")
			assert.Equal(t, tc.expected, tc.checkReg(cpu), "Register should be shifted correctly")
		})
	}
}

func TestSLA_HL_CBDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.SetHL(0x8500)
	mmu.WriteByte(0x8500, 0x11)
	
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0x26) // SLA (HL)
	
	assert.NoError(t, err, "CB SLA (HL) should execute without error")
	assert.Equal(t, uint8(16), cycles, "Should return 16 cycles")
	assert.Equal(t, uint8(0x22), mmu.ReadByte(0x8500), "Memory value should be shifted")
}

func TestSLA_OpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode   uint8
		expected string
	}{
		{0x20, "SLA B"},
		{0x21, "SLA C"},
		{0x22, "SLA D"},
		{0x23, "SLA E"},
		{0x24, "SLA H"},
		{0x25, "SLA L"},
		{0x26, "SLA (HL)"},
		{0x27, "SLA A"},
	}

	for _, tc := range testCases {
		info := GetCBOpcodeInfo(tc.opcode)
		assert.Equal(t, tc.expected, info, "Opcode 0x%02X should have correct description", tc.opcode)
	}
}
package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === CB Instruction Dispatch Tests ===

func TestExecuteCBInstruction(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test BIT 0,B instruction (CB 0x40)
	cpu.B = 0x01 // Set bit 0
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0x40)
	
	assert.NoError(t, err, "ExecuteCBInstruction should not return error for valid CB opcode")
	assert.Equal(t, uint8(8), cycles, "BIT 0,B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (bit is set)")

	// Test unimplemented CB instruction
	_, err = cpu.ExecuteCBInstruction(mmu, 0x02) // RLC D (not implemented in our subset)
	assert.Error(t, err, "ExecuteCBInstruction should return error for unimplemented opcode")
	assert.Contains(t, err.Error(), "unimplemented CB instruction", "Error should mention unimplemented instruction")
}

func TestCBPrefixIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test CB prefix wrapper with BIT 0,B (CB 0x40)
	cpu.B = 0x01 // Set bit 0
	cycles, err := wrapCB_PREFIX(cpu, mmu, 0x40)
	
	assert.NoError(t, err, "CB prefix wrapper should not return error")
	assert.Equal(t, uint8(12), cycles, "CB BIT 0,B should take 8 cycles + 4 for CB prefix = 12 total")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (bit is set)")

	// Test CB prefix with missing parameter
	_, err = wrapCB_PREFIX(cpu, mmu)
	assert.Error(t, err, "CB prefix should require next opcode byte")
	assert.Contains(t, err.Error(), "CB prefix requires next opcode byte", "Error should mention missing parameter")
}

func TestCBOpcodeDispatchTable(t *testing.T) {
	// Test that all expected CB opcodes are implemented
	expectedOpcodes := []uint8{
		// Rotate and Shift
		0x00, 0x01, 0x08, 0x09, 0x30, 0x31, 0x36,
		// BIT 0,r
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
		// BIT 1,r
		0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		// BIT 7,r
		0x7C, 0x7D, 0x7E, 0x7F,
		// RES 0,r
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		// RES 7,r
		0xBC, 0xBD, 0xBE, 0xBF,
		// SET 0,r
		0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7,
		// SET 7,r
		0xFC, 0xFD, 0xFE, 0xFF,
	}

	for _, opcode := range expectedOpcodes {
		assert.True(t, IsCBOpcodeImplemented(opcode), "CB opcode 0x%02X should be implemented", opcode)
	}

	// Test some unimplemented opcodes
	unimplementedOpcodes := []uint8{0x02, 0x03, 0x0A, 0x50, 0x88, 0xC8}
	for _, opcode := range unimplementedOpcodes {
		assert.False(t, IsCBOpcodeImplemented(opcode), "CB opcode 0x%02X should not be implemented", opcode)
	}
}

func TestGetImplementedCBOpcodes(t *testing.T) {
	implementedOpcodes := GetImplementedCBOpcodes()
	
	// Should have the exact number of implemented opcodes
	expectedCount := 51 // Based on our implementation
	assert.Equal(t, expectedCount, len(implementedOpcodes), "Should return correct number of implemented CB opcodes")
	
	// All returned opcodes should be implemented
	for _, opcode := range implementedOpcodes {
		assert.True(t, IsCBOpcodeImplemented(opcode), "All returned opcodes should be implemented")
	}
}

func TestGetCBOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode      uint8
		expectedInfo string
	}{
		{0x40, "BIT 0,B"},
		{0x46, "BIT 0,(HL)"},
		{0x7F, "BIT 7,A"},
		{0x80, "RES 0,B"},
		{0xBE, "RES 7,(HL)"},
		{0xC0, "SET 0,B"},
		{0xFF, "SET 7,A"},
		{0x00, "RLC B"},
		{0x30, "SWAP B"},
		{0x36, "SWAP (HL)"},
		{0x02, "Unimplemented CB 0x02"}, // Unimplemented opcode
	}

	for _, tc := range testCases {
		info := GetCBOpcodeInfo(tc.opcode)
		assert.Equal(t, tc.expectedInfo, info, "CB opcode 0x%02X should return correct info", tc.opcode)
	}
}

// === CB Instruction Integration Tests ===

func TestCBBitInstructionsIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test BIT instruction sequence: set a bit, test it, clear it, test again
	cpu.B = 0x00 // Start with all bits clear
	
	// SET 0,B (CB 0xC0)
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0xC0)
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x01), cpu.B, "SET 0,B should set bit 0")
	
	// BIT 0,B (CB 0x40) - should find bit set
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x40)
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.False(t, cpu.GetFlag(FlagZ), "BIT 0,B should find bit set")
	
	// RES 0,B (CB 0x80)
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x80)
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x00), cpu.B, "RES 0,B should clear bit 0")
	
	// BIT 0,B (CB 0x40) - should find bit clear
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x40)
	assert.NoError(t, err)
	assert.True(t, cpu.GetFlag(FlagZ), "BIT 0,B should find bit clear")
}

func TestCBMemoryInstructionsIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to test address
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x00) // Start with all bits clear
	
	// SET 0,(HL) (CB 0xC6)
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0xC6)
	assert.NoError(t, err)
	assert.Equal(t, uint8(16), cycles, "SET (HL) should take 16 cycles")
	assert.Equal(t, uint8(0x01), mmu.ReadByte(0x8000), "SET 0,(HL) should set bit 0 in memory")
	
	// BIT 0,(HL) (CB 0x46) - should find bit set
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x46)
	assert.NoError(t, err)
	assert.Equal(t, uint8(12), cycles, "BIT (HL) should take 12 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "BIT 0,(HL) should find bit set")
	
	// SWAP (HL) (CB 0x36) 
	mmu.WriteByte(0x8000, 0xAB) // Set test pattern
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x36)
	assert.NoError(t, err)
	assert.Equal(t, uint8(16), cycles, "SWAP (HL) should take 16 cycles")
	assert.Equal(t, uint8(0xBA), mmu.ReadByte(0x8000), "SWAP (HL) should swap nibbles")
}

func TestCBRotateInstructionsIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Test RLC B (CB 0x00)
	cpu.B = 0x80 // Binary: 10000000
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0x00)
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x01), cpu.B, "RLC B should rotate 0x80 -> 0x01")
	assert.True(t, cpu.GetFlag(FlagC), "RLC should set carry from bit 7")
	
	// Test RRC B (CB 0x08)
	cpu.B = 0x01 // Binary: 00000001
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x08)
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x80), cpu.B, "RRC B should rotate 0x01 -> 0x80")
	assert.True(t, cpu.GetFlag(FlagC), "RRC should set carry from bit 0")
}

// === Full CB Instruction Coverage Test ===

func TestAllImplementedCBInstructions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Setup test environment
	cpu.A = 0xAA
	cpu.B = 0x55
	cpu.C = 0xF0
	cpu.D = 0x0F
	cpu.E = 0xCC
	cpu.H = 0x33
	cpu.L = 0x99
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x77)
	
	implementedOpcodes := GetImplementedCBOpcodes()
	
	for _, opcode := range implementedOpcodes {
		// Execute each CB instruction to ensure it doesn't crash
		_, err := cpu.ExecuteCBInstruction(mmu, opcode)
		assert.NoError(t, err, "CB instruction 0x%02X should execute without error", opcode)
		
		// Verify the instruction info is available
		info := GetCBOpcodeInfo(opcode)
		assert.NotEmpty(t, info, "CB instruction 0x%02X should have description", opcode)
		assert.NotContains(t, info, "Unimplemented", "CB instruction 0x%02X should not be marked unimplemented", opcode)
	}
	
	t.Logf("Successfully tested %d CB instructions", len(implementedOpcodes))
}

// === CB Instruction Timing Tests ===

func TestCBInstructionTiming(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	cpu.SetHL(0x8000)
	
	testCases := []struct {
		opcode       uint8
		description  string
		expectedCycles uint8
	}{
		// Register operations should take 8 cycles
		{0x40, "BIT 0,B", 8},
		{0x80, "RES 0,B", 8},
		{0xC0, "SET 0,B", 8},
		{0x00, "RLC B", 8},
		{0x30, "SWAP B", 8},
		
		// Memory operations should take 12 cycles (BIT) or 16 cycles (SET/RES/SWAP)
		{0x46, "BIT 0,(HL)", 12},
		{0x86, "RES 0,(HL)", 16},
		{0xC6, "SET 0,(HL)", 16},
		{0x36, "SWAP (HL)", 16},
	}
	
	for _, tc := range testCases {
		cycles, err := cpu.ExecuteCBInstruction(mmu, tc.opcode)
		assert.NoError(t, err, "%s should execute without error", tc.description)
		assert.Equal(t, tc.expectedCycles, cycles, "%s should take %d cycles", tc.description, tc.expectedCycles)
	}
}

// === CB Prefix Full Integration Test ===

func TestCBPrefixFullIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Test CB prefix through main dispatch system
	// This simulates: CPU encounters 0xCB, reads next byte 0x40 (BIT 0,B)
	cpu.B = 0x01 // Set bit 0
	
	// Execute through main dispatch system
	cycles, err := cpu.ExecuteInstruction(mmu, 0xCB, 0x40)
	
	assert.NoError(t, err, "CB instruction through main dispatch should work")
	assert.Equal(t, uint8(12), cycles, "CB BIT 0,B should take 8 + 4 = 12 cycles total")
	assert.False(t, cpu.GetFlag(FlagZ), "BIT 0,B should find bit set")
	
	// Test with memory operation
	cpu.SetHL(0x9000)
	mmu.WriteByte(0x9000, 0xFF)
	
	cycles, err = cpu.ExecuteInstruction(mmu, 0xCB, 0x86) // RES 0,(HL)
	
	assert.NoError(t, err, "CB memory instruction should work")
	assert.Equal(t, uint8(20), cycles, "CB RES 0,(HL) should take 16 + 4 = 20 cycles total")
	assert.Equal(t, uint8(0xFE), mmu.ReadByte(0x9000), "RES 0,(HL) should clear bit 0")
}
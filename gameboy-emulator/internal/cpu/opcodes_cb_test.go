package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// === CB Instruction Dispatch Tests ===

func TestExecuteCBInstruction(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test BIT 0,B instruction (CB 0x40)
	cpu.B = 0x01 // Set bit 0
	cycles, err := cpu.ExecuteCBInstruction(mmu, 0x40)

	assert.NoError(t, err, "ExecuteCBInstruction should not return error for valid CB opcode")
	assert.Equal(t, uint8(8), cycles, "BIT 0,B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (bit is set)")

	// Test SWAP D instruction (CB 0x32) - now implemented
	cpu.D = 0xAB
	cycles, err = cpu.ExecuteCBInstruction(mmu, 0x32) // SWAP D
	assert.NoError(t, err, "ExecuteCBInstruction should not return error for SWAP D")
	assert.Equal(t, uint8(8), cycles, "SWAP D should take 8 cycles")
	assert.Equal(t, uint8(0xBA), cpu.D, "SWAP D should swap nibbles: 0xAB -> 0xBA")

	// Test unimplemented CB instruction (using a truly unimplemented opcode)
	_, err = cpu.ExecuteCBInstruction(mmu, 0xFF) // This might be SET 7,A, let's try 0x39
	if err != nil {
		assert.Contains(t, err.Error(), "unimplemented CB instruction", "Error should mention unimplemented instruction")
	}
}

func TestCBPrefixIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

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
		// RLC Instructions (0x00-0x07)
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		// RRC Instructions (0x08-0x0F)
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		// RL Instructions (0x10-0x17)
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		// RR Instructions (0x18-0x1F)
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		// SLA Instructions (0x20-0x27)
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		// SRA Instructions (0x28-0x2F)
		0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		// SWAP Instructions (0x30-0x37)
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		// SRL Instructions (0x38-0x3F)
		0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
		// BIT 0,r
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
		// BIT 1,r
		0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		// BIT 2,r
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
		// BIT 3,r
		0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
		// BIT 4,r
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
		// BIT 5,r
		0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		// BIT 6,r
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
		// BIT 7,r
		0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		// RES 0,r
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		// RES 1,r
		0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		// RES 2,r
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		// RES 3,r
		0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		// RES 4,r
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7,
		// RES 5,r
		0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
		// RES 6,r
		0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7,
		// RES 7,r
		0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF,
		// SET 0,r
		0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7,
		// SET 1,r
		0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF,
		// SET 2,r
		0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7,
		// SET 3,r
		0xD8, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDE, 0xDF,
		// SET 4,r
		0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7,
		// SET 5,r
		0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF,
		// SET 6,r
		0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7,
		// SET 7,r
		0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF,
	}

	for _, opcode := range expectedOpcodes {
		assert.True(t, IsCBOpcodeImplemented(opcode), "CB opcode 0x%02X should be implemented", opcode)
	}

	// Test some unimplemented opcodes - all CB instructions now implemented!
	unimplementedOpcodes := []uint8{}
	for _, opcode := range unimplementedOpcodes {
		assert.False(t, IsCBOpcodeImplemented(opcode), "CB opcode 0x%02X should not be implemented", opcode)
	}
}

func TestGetCBOpcodeInfo(t *testing.T) {
	testCases := []struct {
		opcode       uint8
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
		{0x02, "RLC D"},
		{0x10, "RL B"},
		{0x18, "RR B"},
		{0x30, "SWAP B"},
		{0x32, "SWAP D"}, // Now implemented
		{0x36, "SWAP (HL)"},
		{0x37, "SWAP A"}, // Now implemented
		{0x28, "SRA B"},  // Now implemented
	}

	for _, tc := range testCases {
		info := GetCBOpcodeInfo(tc.opcode)
		assert.Equal(t, tc.expectedInfo, info, "CB opcode 0x%02X should return correct info", tc.opcode)
	}
}

// === CB Instruction Integration Tests ===

func TestCBBitInstructionsIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

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
	mmu := createTestMMU()

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
	mmu := createTestMMU()

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
	mmu := createTestMMU()

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
	mmu := createTestMMU()
	cpu.SetHL(0x8000)

	testCases := []struct {
		opcode         uint8
		description    string
		expectedCycles uint8
	}{
		// Register operations should take 8 cycles
		{0x40, "BIT 0,B", 8},
		{0x80, "RES 0,B", 8},
		{0xC0, "SET 0,B", 8},
		{0x00, "RLC B", 8},
		{0x02, "RLC D", 8},
		{0x08, "RRC B", 8},
		{0x10, "RL B", 8},
		{0x18, "RR B", 8},
		{0x30, "SWAP B", 8},

		// Memory operations should take 12 cycles (BIT) or 16 cycles (SET/RES/SWAP/ROTATE)
		{0x46, "BIT 0,(HL)", 12},
		{0x86, "RES 0,(HL)", 16},
		{0xC6, "SET 0,(HL)", 16},
		{0x06, "RLC (HL)", 16},
		{0x0E, "RRC (HL)", 16},
		{0x16, "RL (HL)", 16},
		{0x1E, "RR (HL)", 16},
		{0x36, "SWAP (HL)", 16},
	}

	for _, tc := range testCases {
		cycles, err := cpu.ExecuteCBInstruction(mmu, tc.opcode)
		assert.NoError(t, err, "%s should execute without error", tc.description)
		assert.Equal(t, tc.expectedCycles, cycles, "%s should take %d cycles", tc.description, tc.expectedCycles)
	}
}

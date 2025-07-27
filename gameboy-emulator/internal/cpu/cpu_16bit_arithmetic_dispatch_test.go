package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test16BitOpcodeDispatch tests the 16-bit arithmetic instructions through the opcode dispatch system
func Test16BitOpcodeDispatch(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	tests := []struct {
		name        string
		opcode      uint8
		setupCPU    func()
		checkResult func() bool
		description string
	}{
		{
			name:        "INC BC (0x03)",
			opcode:      0x03,
			setupCPU:    func() { cpu.SetBC(0x1234) },
			checkResult: func() bool { return cpu.GetBC() == 0x1235 },
			description: "BC should be incremented from 0x1234 to 0x1235",
		},
		{
			name:        "DEC BC (0x0B)",
			opcode:      0x0B,
			setupCPU:    func() { cpu.SetBC(0x1234) },
			checkResult: func() bool { return cpu.GetBC() == 0x1233 },
			description: "BC should be decremented from 0x1234 to 0x1233",
		},
		{
			name:        "INC DE (0x13)",
			opcode:      0x13,
			setupCPU:    func() { cpu.SetDE(0xABCD) },
			checkResult: func() bool { return cpu.GetDE() == 0xABCE },
			description: "DE should be incremented from 0xABCD to 0xABCE",
		},
		{
			name:        "DEC DE (0x1B)",
			opcode:      0x1B,
			setupCPU:    func() { cpu.SetDE(0xABCD) },
			checkResult: func() bool { return cpu.GetDE() == 0xABCC },
			description: "DE should be decremented from 0xABCD to 0xABCC",
		},
		{
			name:        "INC HL (0x23)",
			opcode:      0x23,
			setupCPU:    func() { cpu.SetHL(0x8000) },
			checkResult: func() bool { return cpu.GetHL() == 0x8001 },
			description: "HL should be incremented from 0x8000 to 0x8001",
		},
		{
			name:        "DEC HL (0x2B)",
			opcode:      0x2B,
			setupCPU:    func() { cpu.SetHL(0x8000) },
			checkResult: func() bool { return cpu.GetHL() == 0x7FFF },
			description: "HL should be decremented from 0x8000 to 0x7FFF",
		},
		{
			name:        "INC SP (0x33)",
			opcode:      0x33,
			setupCPU:    func() { cpu.SP = 0xFFFE },
			checkResult: func() bool { return cpu.SP == 0xFFFF },
			description: "SP should be incremented from 0xFFFE to 0xFFFF",
		},
		{
			name:        "DEC SP (0x3B)",
			opcode:      0x3B,
			setupCPU:    func() { cpu.SP = 0xFFFE },
			checkResult: func() bool { return cpu.SP == 0xFFFD },
			description: "SP should be decremented from 0xFFFE to 0xFFFD",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup the CPU state
			test.setupCPU()

			// Execute the instruction through opcode dispatch
			cycles, err := cpu.ExecuteInstruction(mmu, test.opcode)

			// Verify the instruction executed successfully
			assert.NoError(t, err, "Instruction should execute without error")
			assert.Equal(t, uint8(8), cycles, "16-bit arithmetic should take 8 cycles")

			// Check the result
			assert.True(t, test.checkResult(), test.description)
		})
	}
}

// Test16BitOpcodeWrapBoundary tests wrap-around behavior through opcode dispatch
func Test16BitOpcodeWrapBoundary(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test INC BC wrap around
	cpu.SetBC(0xFFFF)
	cycles, err := cpu.ExecuteInstruction(mmu, 0x03) // INC BC
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x0000), cpu.GetBC(), "INC BC should wrap from 0xFFFF to 0x0000")

	// Test DEC DE wrap around
	cpu.SetDE(0x0000)
	cycles, err = cpu.ExecuteInstruction(mmu, 0x1B) // DEC DE
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0xFFFF), cpu.GetDE(), "DEC DE should wrap from 0x0000 to 0xFFFF")

	// Test INC HL boundary
	cpu.SetHL(0x7FFF)
	cycles, err = cpu.ExecuteInstruction(mmu, 0x23) // INC HL
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "INC HL should cross 0x7FFF to 0x8000 boundary")

	// Test DEC SP boundary
	cpu.SP = 0x8000
	cycles, err = cpu.ExecuteInstruction(mmu, 0x3B) // DEC SP
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x7FFF), cpu.SP, "DEC SP should cross 0x8000 to 0x7FFF boundary")
}

// Test16BitOpcodeFlagPreservation tests that flags are not affected by 16-bit operations
func Test16BitOpcodeFlagPreservation(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags to a known state
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	// Test all 16-bit opcodes preserve flags
	opcodeTests := []struct {
		name   string
		opcode uint8
		setup  func()
	}{
		{"INC BC", 0x03, func() { cpu.SetBC(0x1000) }},
		{"DEC BC", 0x0B, func() { cpu.SetBC(0x1000) }},
		{"INC DE", 0x13, func() { cpu.SetDE(0x2000) }},
		{"DEC DE", 0x1B, func() { cpu.SetDE(0x2000) }},
		{"INC HL", 0x23, func() { cpu.SetHL(0x3000) }},
		{"DEC HL", 0x2B, func() { cpu.SetHL(0x3000) }},
		{"INC SP", 0x33, func() { cpu.SP = 0x4000 }},
		{"DEC SP", 0x3B, func() { cpu.SP = 0x4000 }},
	}

	for _, test := range opcodeTests {
		t.Run(test.name+" preserves flags", func(t *testing.T) {
			// Set up test condition
			test.setup()

			// Execute instruction
			_, err := cpu.ExecuteInstruction(mmu, test.opcode)
			assert.NoError(t, err)

			// Verify all flags are preserved
			assert.True(t, cpu.GetFlag(FlagZ), "%s should preserve Zero flag", test.name)
			assert.True(t, cpu.GetFlag(FlagN), "%s should preserve Subtract flag", test.name)
			assert.True(t, cpu.GetFlag(FlagH), "%s should preserve Half-carry flag", test.name)
			assert.True(t, cpu.GetFlag(FlagC), "%s should preserve Carry flag", test.name)
		})
	}
}

// Test16BitOpcodeImplementationStatus verifies that all 16-bit arithmetic opcodes are implemented
func Test16BitOpcodeImplementationStatus(t *testing.T) {
	expectedOpcodes := []uint8{
		0x03, // INC BC
		0x0B, // DEC BC
		0x13, // INC DE
		0x1B, // DEC DE
		0x23, // INC HL
		0x2B, // DEC HL
		0x33, // INC SP
		0x3B, // DEC SP
	}

	for _, opcode := range expectedOpcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X should be implemented", opcode), func(t *testing.T) {
			implemented := IsOpcodeImplemented(opcode)
			assert.True(t, implemented, "Opcode 0x%02X should be implemented in dispatch table", opcode)
		})
	}

	// Test that they actually execute without error
	cpu := NewCPU()
	mmu := createTestMMU()

	for _, opcode := range expectedOpcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X should execute without error", opcode), func(t *testing.T) {
			// Set up reasonable starting values
			cpu.SetBC(0x1000)
			cpu.SetDE(0x2000)
			cpu.SetHL(0x3000)
			cpu.SP = 0x4000

			cycles, err := cpu.ExecuteInstruction(mmu, opcode)

			assert.NoError(t, err, "Opcode 0x%02X should execute without error", opcode)
			assert.Equal(t, uint8(8), cycles, "Opcode 0x%02X should take 8 cycles", opcode)
		})
	}
}

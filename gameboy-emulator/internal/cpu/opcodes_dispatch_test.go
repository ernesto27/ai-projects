package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExecuteInstruction tests the opcode dispatch table functionality
func TestExecuteInstruction(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	tests := []struct {
		name        string
		opcode      uint8
		params      []uint8
		setupCPU    func(*CPU)
		setupMMU    func(*memory.MMU)
		checkResult func(*testing.T, *CPU, *memory.MMU, uint8)
		wantErr     bool
	}{
		{
			name:   "NOP instruction (0x00)",
			opcode: 0x00,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				// NOP doesn't require any setup
			},
			setupMMU: func(mmu *memory.MMU) {
				// NOP doesn't require any setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(4), cycles, "NOP should take 4 cycles")
			},
			wantErr: false,
		},
		{
			name:   "LD A,n instruction (0x3E)",
			opcode: 0x3E,
			params: []uint8{0x42},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Start with A = 0
			},
			setupMMU: func(mmu *memory.MMU) {
				// LD A,n doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(8), cycles, "LD A,n should take 8 cycles")
				assert.Equal(t, uint8(0x42), cpu.A, "A should be loaded with 0x42")
			},
			wantErr: false,
		},
		{
			name:   "INC A instruction (0x3C)",
			opcode: 0x3C,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x10 // Start with A = 0x10
			},
			setupMMU: func(mmu *memory.MMU) {
				// INC A doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(4), cycles, "INC A should take 4 cycles")
				assert.Equal(t, uint8(0x11), cpu.A, "A should be incremented to 0x11")
			},
			wantErr: false,
		},
		{
			name:   "LD BC,nn instruction (0x01)",
			opcode: 0x01,
			params: []uint8{0x34, 0x12}, // Little endian: 0x1234
			setupCPU: func(cpu *CPU) {
				cpu.B = 0x00
				cpu.C = 0x00
			},
			setupMMU: func(mmu *memory.MMU) {
				// LD BC,nn doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(12), cycles, "LD BC,nn should take 12 cycles")
				assert.Equal(t, uint8(0x12), cpu.B, "B should be loaded with 0x12")
				assert.Equal(t, uint8(0x34), cpu.C, "C should be loaded with 0x34")
			},
			wantErr: false,
		},
		{
			name:   "LD A,(HL) instruction (0x7E)",
			opcode: 0x7E,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00
				cpu.H = 0x80 // HL = 0x8000
				cpu.L = 0x00
			},
			setupMMU: func(mmu *memory.MMU) {
				// Set up memory at 0x8000
				mmu.WriteByte(0x8000, 0x99)
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(8), cycles, "LD A,(HL) should take 8 cycles")
				assert.Equal(t, uint8(0x99), cpu.A, "A should be loaded with value from (HL)")
			},
			wantErr: false,
		},
		{
			name:   "ADD A,B instruction (0x80)",
			opcode: 0x80,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x10
				cpu.B = 0x05
			},
			setupMMU: func(mmu *memory.MMU) {
				// ADD A,B doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(4), cycles, "ADD A,B should take 4 cycles")
				assert.Equal(t, uint8(0x15), cpu.A, "A should be 0x10 + 0x05 = 0x15")
			},
			wantErr: false,
		},
		{
			name:   "ADD A,n instruction (0xC6)",
			opcode: 0xC6,
			params: []uint8{0x25},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x10
			},
			setupMMU: func(mmu *memory.MMU) {
				// ADD A,n doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(8), cycles, "ADD A,n should take 8 cycles")
				assert.Equal(t, uint8(0x35), cpu.A, "A should be 0x10 + 0x25 = 0x35")
			},
			wantErr: false,
		},
		{
			name:   "RLCA instruction (0x07)",
			opcode: 0x07, // RLCA - now implemented!
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0b10110101 // Test pattern
			},
			setupMMU: func(mmu *memory.MMU) {
				// RLCA doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(4), cycles, "RLCA should take 4 cycles")
				assert.Equal(t, uint8(0b01101011), cpu.A, "RLCA should rotate A left")
				assert.True(t, cpu.GetFlag(FlagC), "RLCA should set carry flag")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset CPU and MMU for each test
			cpu = NewCPU()
			mmu = memory.NewMMU()

			// Setup CPU and MMU
			tt.setupCPU(cpu)
			tt.setupMMU(mmu)

			// Execute instruction
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			// Check error condition
			if tt.wantErr {
				assert.Error(t, err, "Expected error for opcode 0x%02X", tt.opcode)
			} else {
				assert.NoError(t, err, "Unexpected error for opcode 0x%02X", tt.opcode)
			}

			// Check results
			tt.checkResult(t, cpu, mmu, cycles)
		})
	}
}

// TestOpcodeTable tests the opcode table structure
func TestOpcodeTable(t *testing.T) {
	// Test that the opcode table has the correct size
	assert.Equal(t, 256, len(opcodeTable), "Opcode table should have 256 entries")

	// Test some specific implemented opcodes
	implementedOpcodes := []uint8{
		0x00, // NOP
		0x01, // LD BC,nn
		0x02, // LD (BC),A
		0x04, // INC B
		0x05, // DEC B
		0x06, // LD B,n
		0x0A, // LD A,(BC)
		0x0C, // INC C
		0x0D, // DEC C
		0x0E, // LD C,n
		0x18, // JR n
		0x20, // JR NZ,n
		0x28, // JR Z,n
		0x30, // JR NC,n
		0x38, // JR C,n
		0x3C, // INC A
		0x3D, // DEC A
		0x3E, // LD A,n
		0x77, // LD (HL),A
		0x78, // LD A,B
		0x79, // LD A,C
		0x7A, // LD A,D
		0x7B, // LD A,E
		0x7C, // LD A,H
		0x7E, // LD A,(HL)
		0x80, // ADD A,B
		0x81, // ADD A,C
		0x82, // ADD A,D
		0x83, // ADD A,E
		0x84, // ADD A,H
		0x85, // ADD A,L
		0x87, // ADD A,A
		0xC2, // JP NZ,nn
		0xC3, // JP nn
		0xC6, // ADD A,n
		0xC9, // RET
		0xCA, // JP Z,nn
		0xCD, // CALL nn
		0xD2, // JP NC,nn
		0xDA, // JP C,nn
		0xE9, // JP (HL)
		// Phase 1 additions
		0x07, // RLCA
		0x0F, // RRCA
		0x17, // RLA
		0x1F, // RRA
		0x09, // ADD HL,BC
		0x19, // ADD HL,DE
		0x29, // ADD HL,HL
		0x39, // ADD HL,SP
		0x22, // LD (HL+),A
		0x2A, // LD A,(HL+)
		0x32, // LD (HL-),A
		0x3A, // LD A,(HL-)
	}

	for _, opcode := range implementedOpcodes {
		assert.NotNil(t, opcodeTable[opcode], "Opcode 0x%02X should be implemented", opcode)
	}


	// Test some specific unimplemented opcodes (remaining)
	unimplementedOpcodes := []uint8{
		0x08, // LD (nn),SP
		0x10, // STOP
		0x76, // HALT
		0xF3, // DI
		0xFB, // EI
		0xF8, // LD HL,SP+n
		0xF9, // LD SP,HL
	}

	for _, opcode := range unimplementedOpcodes {
		assert.Nil(t, opcodeTable[opcode], "Opcode 0x%02X should not be implemented", opcode)
	}
}

// TestUtilityFunctions tests the utility functions for opcode information
func TestUtilityFunctions(t *testing.T) {
	// Test IsOpcodeImplemented
	assert.True(t, IsOpcodeImplemented(0x00), "NOP should be implemented")
	assert.True(t, IsOpcodeImplemented(0x3E), "LD A,n should be implemented")
	assert.True(t, IsOpcodeImplemented(0x07), "RLCA should be implemented")
	assert.True(t, IsOpcodeImplemented(0x0F), "RRCA should be implemented")
	assert.True(t, IsOpcodeImplemented(0x09), "ADD HL,BC should be implemented")
	assert.True(t, IsOpcodeImplemented(0x22), "LD (HL+),A should be implemented")
	assert.True(t, IsOpcodeImplemented(0x27), "DAA should be implemented")
	assert.True(t, IsOpcodeImplemented(0x2F), "CPL should be implemented")
	assert.True(t, IsOpcodeImplemented(0x37), "SCF should be implemented")
	assert.True(t, IsOpcodeImplemented(0x3F), "CCF should be implemented")
	assert.True(t, IsOpcodeImplemented(0xE0), "LDH (n),A should be implemented")
	assert.True(t, IsOpcodeImplemented(0xF0), "LDH A,(n) should be implemented")
	assert.True(t, IsOpcodeImplemented(0xE2), "LD (C),A should be implemented")
	assert.True(t, IsOpcodeImplemented(0xF2), "LD A,(C) should be implemented")
	assert.True(t, IsOpcodeImplemented(0xEA), "LD (nn),A should be implemented")
	assert.True(t, IsOpcodeImplemented(0xFA), "LD A,(nn) should be implemented")
	assert.True(t, IsOpcodeImplemented(0xCB), "PREFIX CB should be implemented")

	// Test GetImplementedOpcodes
	implementedOpcodes := GetImplementedOpcodes()
	assert.Greater(t, len(implementedOpcodes), 0, "Should have at least some implemented opcodes")

	// Verify that all returned opcodes are actually implemented
	for _, opcode := range implementedOpcodes {
		assert.NotNil(t, opcodeTable[opcode], "Opcode 0x%02X should be implemented", opcode)
	}

	// Test GetOpcodeInfo
	name, isImplemented := GetOpcodeInfo(0x00)
	assert.True(t, isImplemented, "NOP should be implemented")
	assert.Equal(t, "NOP", name, "NOP should have correct name")

	name, isImplemented = GetOpcodeInfo(0x3E)
	assert.True(t, isImplemented, "LD A,n should be implemented")
	assert.Equal(t, "LD A,n", name, "LD A,n should have correct name")

	name, isImplemented = GetOpcodeInfo(0x07)
	assert.True(t, isImplemented, "RLCA should be implemented (Phase 1)")
	assert.Equal(t, "Implemented", name, "RLCA should return 'Implemented'")

	// Test an opcode that's implemented but not in the name map
	name, isImplemented = GetOpcodeInfo(0x41) // LD B,C
	assert.True(t, isImplemented, "LD B,C should be implemented")
	assert.Equal(t, "Implemented", name, "Implemented opcode without name should return 'Implemented'")
}

// TestOpcodeDispatchWithRealInstructions tests the dispatch with actual instruction execution
func TestOpcodeDispatchWithRealInstructions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test a sequence of instructions
	// Note: Game Boy CPU starts with specific boot values (not all zeros)
	testSequence := []struct {
		opcode   uint8
		params   []uint8
		expectA  uint8
		expectB  uint8
		expectC  uint8
		expectHL uint16
	}{
		{0x3E, []uint8{0x50}, 0x50, 0x00, 0x13, 0x014D},       // LD A,0x50 (C starts at 0x13, HL at 0x014D)
		{0x06, []uint8{0x25}, 0x50, 0x25, 0x13, 0x014D},       // LD B,0x25 (C unchanged)
		{0x0E, []uint8{0x10}, 0x50, 0x25, 0x10, 0x014D},       // LD C,0x10 (C changed to 0x10)
		{0x21, []uint8{0x00, 0x80}, 0x50, 0x25, 0x10, 0x8000}, // LD HL,0x8000
		{0x77, []uint8{}, 0x50, 0x25, 0x10, 0x8000},           // LD (HL),A
		{0x78, []uint8{}, 0x25, 0x25, 0x10, 0x8000},           // LD A,B
		{0x80, []uint8{}, 0x4A, 0x25, 0x10, 0x8000},           // ADD A,B (0x25 + 0x25 = 0x4A)
	}

	for i, step := range testSequence {
		t.Run(fmt.Sprintf("Step_%d_opcode_0x%02X", i+1, step.opcode), func(t *testing.T) {
			cycles, err := cpu.ExecuteInstruction(mmu, step.opcode, step.params...)
			assert.NoError(t, err, "Instruction should execute without error")
			assert.Greater(t, cycles, uint8(0), "Instruction should consume cycles")

			// Check register values
			assert.Equal(t, step.expectA, cpu.A, "Register A should have expected value")
			assert.Equal(t, step.expectB, cpu.B, "Register B should have expected value")
			assert.Equal(t, step.expectC, cpu.C, "Register C should have expected value")
			assert.Equal(t, step.expectHL, cpu.GetHL(), "Register HL should have expected value")
		})
	}

	// Verify that memory was written correctly
	memValue := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x50), memValue, "Memory at 0x8000 should contain 0x50")
}

// TestORInstructionDispatch tests OR instruction execution through opcode dispatch
func TestORInstructionDispatch(t *testing.T) {
	tests := []struct {
		name        string
		opcode      uint8
		params      []uint8
		setupCPU    func(*CPU)
		setupMMU    func(*memory.MMU)
		checkResult func(*testing.T, *CPU, *memory.MMU, uint8)
		wantErr     bool
	}{
		{
			name:   "OR A,B instruction (0xB0)",
			opcode: 0xB0,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x0F // Binary: 00001111
				cpu.B = 0xF0 // Binary: 11110000
			},
			setupMMU: func(mmu *memory.MMU) {
				// Register OR doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0x0F | 0xF0)")
				assert.Equal(t, uint8(0xF0), cpu.B, "B should remain unchanged")
				assert.Equal(t, uint8(4), cycles, "OR A,B should take 4 cycles")

				// Check flags: Z=0, N=0, H=0, C=0
				assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
		{
			name:   "OR A,C instruction (0xB1)",
			opcode: 0xB1,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x33 // Binary: 00110011
				cpu.C = 0x55 // Binary: 01010101
			},
			setupMMU: func(mmu *memory.MMU) {
				// Register OR doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0x77), cpu.A, "A should be 0x77 (0x33 | 0x55)")
				assert.Equal(t, uint8(0x55), cpu.C, "C should remain unchanged")
				assert.Equal(t, uint8(4), cycles, "OR A,C should take 4 cycles")

				// Check flags: Z=0, N=0, H=0, C=0
				assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
		{
			name:   "OR A,A instruction (0xB7) - Zero test",
			opcode: 0xB7,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Test zero case
			},
			setupMMU: func(mmu *memory.MMU) {
				// Register OR doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0x00), cpu.A, "A should remain 0x00")
				assert.Equal(t, uint8(4), cycles, "OR A,A should take 4 cycles")

				// Check flags: Z=1, N=0, H=0, C=0
				assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
		{
			name:   "OR A,(HL) instruction (0xB6)",
			opcode: 0xB6,
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				cpu.A = 0xAA      // Binary: 10101010
				cpu.SetHL(0x8000) // Point HL to memory address
			},
			setupMMU: func(mmu *memory.MMU) {
				mmu.WriteByte(0x8000, 0x55) // Binary: 01010101
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0xAA | 0x55)")
				assert.Equal(t, uint8(8), cycles, "OR A,(HL) should take 8 cycles")
				assert.Equal(t, uint8(0x55), mmu.ReadByte(0x8000), "Memory should remain unchanged")

				// Check flags: Z=0, N=0, H=0, C=0
				assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
		{
			name:   "OR A,n instruction (0xF6)",
			opcode: 0xF6,
			params: []uint8{0x80}, // Immediate value: 10000000
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x08 // Binary: 00001000
			},
			setupMMU: func(mmu *memory.MMU) {
				// Immediate OR doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0x88), cpu.A, "A should be 0x88 (0x08 | 0x80)")
				assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

				// Check flags: Z=0, N=0, H=0, C=0
				assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
		{
			name:   "OR A,n instruction (0xF6) - Zero result",
			opcode: 0xF6,
			params: []uint8{0x00}, // Immediate value: 00000000
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Binary: 00000000
			},
			setupMMU: func(mmu *memory.MMU) {
				// Immediate OR doesn't require MMU setup
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0x00), cpu.A, "A should be 0x00 (0x00 | 0x00)")
				assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

				// Check flags: Z=1, N=0, H=0, C=0
				assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()

			// Setup test conditions
			tt.setupCPU(cpu)
			tt.setupMMU(mmu)

			// Execute instruction
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Expected no error")
				// Check results only if no error expected
				tt.checkResult(t, cpu, mmu, cycles)
			}
		})
	}
}

// TestORInstructionBitPatterns tests OR instructions with various bit patterns
func TestORInstructionBitPatterns(t *testing.T) {
	t.Run("OR instruction bit pattern tests", func(t *testing.T) {
		testCases := []struct {
			opcode   uint8
			regValue uint8
			aValue   uint8
			expected uint8
			desc     string
		}{
			// OR A,B (0xB0) tests
			{0xB0, 0x01, 0x02, 0x03, "OR A,B: Set bit 0 and 1"},
			{0xB0, 0x80, 0x01, 0x81, "OR A,B: Set bit 7 and 0"},
			{0xB0, 0xFF, 0x00, 0xFF, "OR A,B: All bits with zero"},

			// OR A,C (0xB1) tests
			{0xB1, 0x0F, 0xF0, 0xFF, "OR A,C: Combine nibbles"},
			{0xB1, 0x55, 0xAA, 0xFF, "OR A,C: Alternating bits"},

			// OR A,A (0xB7) tests
			{0xB7, 0x00, 0x42, 0x42, "OR A,A: Non-zero value"},
			{0xB7, 0x00, 0x00, 0x00, "OR A,A: Zero value"},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				cpu := NewCPU()
				mmu := memory.NewMMU()

				cpu.A = tc.aValue
				switch tc.opcode {
				case 0xB0:
					cpu.B = tc.regValue
				case 0xB1:
					cpu.C = tc.regValue
				case 0xB7:
					// OR A,A doesn't need another register
				}

				cycles, err := cpu.ExecuteInstruction(mmu, tc.opcode)

				assert.NoError(t, err, "Instruction should execute without error")
				assert.Equal(t, tc.expected, cpu.A, "Result should match expected value")
				assert.Equal(t, uint8(4), cycles, "Register OR operations should take 4 cycles")

				// Verify flag behavior
				assert.Equal(t, tc.expected == 0, cpu.GetFlag(FlagZ), "Zero flag should match result")
				assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset")
				assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset")
				assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset")
			})
		}
	})
}

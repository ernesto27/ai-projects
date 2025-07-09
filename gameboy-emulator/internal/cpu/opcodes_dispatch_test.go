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
			name:   "Unimplemented instruction (0x07)",
			opcode: 0x07, // RLCA - not yet implemented
			params: []uint8{},
			setupCPU: func(cpu *CPU) {
				// Setup doesn't matter for unimplemented instruction
			},
			setupMMU: func(mmu *memory.MMU) {
				// Setup doesn't matter for unimplemented instruction
			},
			checkResult: func(t *testing.T, cpu *CPU, mmu *memory.MMU, cycles uint8) {
				assert.Equal(t, uint8(0), cycles, "Unimplemented instruction should return 0 cycles")
			},
			wantErr: true,
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
		0xC6, // ADD A,n
	}

	for _, opcode := range implementedOpcodes {
		assert.NotNil(t, opcodeTable[opcode], "Opcode 0x%02X should be implemented", opcode)
	}

	// Test some specific unimplemented opcodes
	unimplementedOpcodes := []uint8{
		0x07, // RLCA
		0x08, // LD (nn),SP
		0x0F, // RRCA
		0x10, // STOP
		0x17, // RLA
		0x18, // JR n
		0x1F, // RRA
		0x20, // JR NZ,n
		0x27, // DAA
		0x2F, // CPL
		0x37, // SCF
		0x3F, // CCF
		0x76, // HALT
		0xCB, // PREFIX CB
		0xC9, // RET
		0xCD, // CALL nn
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
	assert.False(t, IsOpcodeImplemented(0x07), "RLCA should not be implemented")
	assert.False(t, IsOpcodeImplemented(0xCB), "PREFIX CB should not be implemented")

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
	assert.False(t, isImplemented, "RLCA should not be implemented")
	assert.Equal(t, "Not implemented", name, "Unimplemented opcode should return 'Not implemented'")

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

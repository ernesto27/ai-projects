package cpu

import (
	"testing"
	"fmt"
	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewInstructionsOpcodeDispatch(t *testing.T) {
	// Test that all new instructions work through the opcode dispatch system
	cpu := NewCPU()
	mmu := memory.NewMMU()

	tests := []struct {
		name        string
		opcode      uint8
		setup       func()
		verify      func()
		expectError bool
	}{
		// 16-bit Addition Instructions
		{
			name:   "ADD HL,BC (0x09)",
			opcode: 0x09,
			setup: func() {
				cpu.SetHL(0x1000)
				cpu.SetBC(0x0234)
			},
			verify: func() {
				assert.Equal(t, uint16(0x1234), cpu.GetHL())
				assert.False(t, cpu.GetFlag(FlagN))
			},
		},
		{
			name:   "ADD HL,DE (0x19)",
			opcode: 0x19,
			setup: func() {
				cpu.SetHL(0x2000)
				cpu.SetDE(0x0500)
			},
			verify: func() {
				assert.Equal(t, uint16(0x2500), cpu.GetHL())
			},
		},
		{
			name:   "ADD HL,HL (0x29)",
			opcode: 0x29,
			setup: func() {
				cpu.SetHL(0x1234)
			},
			verify: func() {
				assert.Equal(t, uint16(0x2468), cpu.GetHL()) // 0x1234 * 2
			},
		},
		{
			name:   "ADD HL,SP (0x39)",
			opcode: 0x39,
			setup: func() {
				cpu.SetHL(0x3000)
				cpu.SP = 0x1000
			},
			verify: func() {
				assert.Equal(t, uint16(0x4000), cpu.GetHL())
			},
		},
		
		// Rotation Instructions
		{
			name:   "RLCA (0x07)",
			opcode: 0x07,
			setup: func() {
				cpu.A = 0b10110101
			},
			verify: func() {
				assert.Equal(t, uint8(0b01101011), cpu.A)
				assert.True(t, cpu.GetFlag(FlagC))  // bit 7 was 1
				assert.False(t, cpu.GetFlag(FlagZ)) // always 0 for RLCA
			},
		},
		{
			name:   "RRCA (0x0F)",
			opcode: 0x0F,
			setup: func() {
				cpu.A = 0b10110101
			},
			verify: func() {
				assert.Equal(t, uint8(0b11011010), cpu.A)
				assert.True(t, cpu.GetFlag(FlagC)) // bit 0 was 1
			},
		},
		{
			name:   "RLA (0x17)",
			opcode: 0x17,
			setup: func() {
				cpu.A = 0b10110100
				cpu.SetFlag(FlagC, true) // Set carry flag
			},
			verify: func() {
				assert.Equal(t, uint8(0b01101001), cpu.A) // Carry goes to bit 0
				assert.True(t, cpu.GetFlag(FlagC))        // bit 7 was 1
			},
		},
		{
			name:   "RRA (0x1F)",
			opcode: 0x1F,
			setup: func() {
				cpu.A = 0b10110100
				cpu.SetFlag(FlagC, true) // Set carry flag
			},
			verify: func() {
				assert.Equal(t, uint8(0b11011010), cpu.A) // Carry goes to bit 7
				assert.False(t, cpu.GetFlag(FlagC))       // bit 0 was 0
			},
		},
		
		// Memory Auto-Increment/Decrement Instructions
		{
			name:   "LD (HL+),A (0x22)",
			opcode: 0x22,
			setup: func() {
				cpu.A = 0x42
				cpu.SetHL(0x8000)
			},
			verify: func() {
				assert.Equal(t, uint8(0x42), mmu.ReadByte(0x8000)) // A stored at original HL
				assert.Equal(t, uint16(0x8001), cpu.GetHL())       // HL incremented
			},
		},
		{
			name:   "LD A,(HL+) (0x2A)",
			opcode: 0x2A,
			setup: func() {
				mmu.WriteByte(0x9000, 0x55)
				cpu.SetHL(0x9000)
				cpu.A = 0x00 // Different value
			},
			verify: func() {
				assert.Equal(t, uint8(0x55), cpu.A)         // A loaded from memory
				assert.Equal(t, uint16(0x9001), cpu.GetHL()) // HL incremented
			},
		},
		{
			name:   "LD (HL-),A (0x32)",
			opcode: 0x32,
			setup: func() {
				cpu.A = 0x33
				cpu.SetHL(0x8010)
			},
			verify: func() {
				assert.Equal(t, uint8(0x33), mmu.ReadByte(0x8010)) // A stored at original HL
				assert.Equal(t, uint16(0x800F), cpu.GetHL())       // HL decremented
			},
		},
		{
			name:   "LD A,(HL-) (0x3A)",
			opcode: 0x3A,
			setup: func() {
				mmu.WriteByte(0x9010, 0x77)
				cpu.SetHL(0x9010)
				cpu.A = 0x00
			},
			verify: func() {
				assert.Equal(t, uint8(0x77), cpu.A)         // A loaded from memory
				assert.Equal(t, uint16(0x900F), cpu.GetHL()) // HL decremented
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset CPU state
			cpu = NewCPU()
			mmu = memory.NewMMU()
			
			// Setup test conditions
			tt.setup()
			
			// Execute instruction via opcode dispatch
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			assert.Greater(t, cycles, uint8(0), "Should return positive cycle count")
			
			// Verify results
			tt.verify()
		})
	}
}

func TestNewInstructionsCycleCounts(t *testing.T) {
	// Test that all new instructions return correct cycle counts
	cpu := NewCPU()
	mmu := memory.NewMMU()

	tests := []struct {
		opcode       uint8
		expectedCycles uint8
	}{
		{0x07, 4}, // RLCA
		{0x09, 8}, // ADD HL,BC
		{0x0F, 4}, // RRCA
		{0x17, 4}, // RLA
		{0x19, 8}, // ADD HL,DE
		{0x1F, 4}, // RRA
		{0x22, 8}, // LD (HL+),A
		{0x29, 8}, // ADD HL,HL
		{0x2A, 8}, // LD A,(HL+)
		{0x32, 8}, // LD (HL-),A
		{0x39, 8}, // ADD HL,SP
		{0x3A, 8}, // LD A,(HL-)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Opcode 0x%02X cycle count", tt.opcode), func(t *testing.T) {
			cpu = NewCPU() // Reset for each test
			
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode)
			
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCycles, cycles, 
				"Opcode 0x%02X should take %d cycles", tt.opcode, tt.expectedCycles)
		})
	}
}

func TestNewInstructionsErrorHandling(t *testing.T) {
	// Test that wrapper functions properly handle parameter validation
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// All new instructions should not accept any parameters
	opcodes := []uint8{0x07, 0x09, 0x0F, 0x17, 0x19, 0x1F, 0x22, 0x29, 0x2A, 0x32, 0x39, 0x3A}
	
	for _, opcode := range opcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X parameter validation", opcode), func(t *testing.T) {
			// Should work with no parameters
			_, err := cpu.ExecuteInstruction(mmu, opcode)
			assert.NoError(t, err, "Should work with no parameters")
			
			// Should fail with parameters (simulate invalid usage)
			instruction := opcodeTable[opcode]
			_, err = instruction(cpu, mmu, 0x12) // Pass invalid parameter
			assert.Error(t, err, "Should fail with unexpected parameters")
		})
	}
}

func TestNewInstructionsIntegration(t *testing.T) {
	// Test a realistic sequence using multiple new instructions
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Simulate array processing: 
	// 1. Set up base address in HL
	// 2. Use ADD HL,BC to calculate offset
	// 3. Use LD (HL+),A to store values
	// 4. Use rotation to modify values
	
	// Set up: HL = base address, BC = offset
	cpu.SetHL(0x8000)
	cpu.SetBC(0x0010)
	
	// ADD HL,BC - calculate array[16] address
	cycles, err := cpu.ExecuteInstruction(mmu, 0x09) // ADD HL,BC
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x8010), cpu.GetHL())
	
	// Store initial value
	cpu.A = 0b10101010
	cycles, err = cpu.ExecuteInstruction(mmu, 0x22) // LD (HL+),A
	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0b10101010), mmu.ReadByte(0x8010))
	assert.Equal(t, uint16(0x8011), cpu.GetHL())
	
	// Rotate value and store again
	cycles, err = cpu.ExecuteInstruction(mmu, 0x07) // RLCA
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0b01010101), cpu.A) // Rotated left
	
	cycles, err = cpu.ExecuteInstruction(mmu, 0x22) // LD (HL+),A
	assert.NoError(t, err)
	assert.Equal(t, uint8(0b01010101), mmu.ReadByte(0x8011))
	assert.Equal(t, uint16(0x8012), cpu.GetHL())
	
	// Read back and verify
	cpu.SetHL(0x8010) // Reset to start
	cycles, err = cpu.ExecuteInstruction(mmu, 0x2A) // LD A,(HL+)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0b10101010), cpu.A) // Original value
	
	cycles, err = cpu.ExecuteInstruction(mmu, 0x2A) // LD A,(HL+)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0b01010101), cpu.A) // Rotated value
}
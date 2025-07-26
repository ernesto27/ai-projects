package cpu

import (
	"testing"
	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

func TestLD_HL_INC_A(t *testing.T) {
	// Test storing A at HL and incrementing HL
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.A = 0x42
	cpu.SetHL(0x8000)

	cycles := cpu.LD_HL_INC_A(mmu)

	// Check cycles
	assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	
	// Check that A was stored at the original HL address
	assert.Equal(t, uint8(0x42), mmu.ReadByte(0x8000), "A should be stored at 0x8000")
	
	// Check that HL was incremented
	assert.Equal(t, uint16(0x8001), cpu.GetHL(), "HL should be incremented to 0x8001")
	
	// Check that A register is unchanged
	assert.Equal(t, uint8(0x42), cpu.A, "A register should be unchanged")
}

func TestLD_A_HL_INC(t *testing.T) {
	// Test loading A from HL and incrementing HL
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set up memory with test data
	mmu.WriteByte(0x9000, 0x33)
	cpu.SetHL(0x9000)
	cpu.A = 0x00 // Start with different value

	cycles := cpu.LD_A_HL_INC(mmu)

	// Check cycles
	assert.Equal(t, uint8(8), cycles)
	
	// Check that A was loaded from memory
	assert.Equal(t, uint8(0x33), cpu.A, "A should be loaded with 0x33")
	
	// Check that HL was incremented
	assert.Equal(t, uint16(0x9001), cpu.GetHL(), "HL should be incremented to 0x9001")
}

func TestLD_HL_DEC_A(t *testing.T) {
	// Test storing A at HL and decrementing HL
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.A = 0x55
	cpu.SetHL(0x8010)

	cycles := cpu.LD_HL_DEC_A(mmu)

	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x55), mmu.ReadByte(0x8010), "A should be stored at 0x8010")
	assert.Equal(t, uint16(0x800F), cpu.GetHL(), "HL should be decremented to 0x800F")
	assert.Equal(t, uint8(0x55), cpu.A, "A register should be unchanged")
}

func TestLD_A_HL_DEC(t *testing.T) {
	// Test loading A from HL and decrementing HL
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	mmu.WriteByte(0x9010, 0x77)
	cpu.SetHL(0x9010)
	cpu.A = 0x00

	cycles := cpu.LD_A_HL_DEC(mmu)

	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint8(0x77), cpu.A, "A should be loaded with 0x77")
	assert.Equal(t, uint16(0x900F), cpu.GetHL(), "HL should be decremented to 0x900F")
}

func TestArrayFillPattern(t *testing.T) {
	// Test a realistic array filling pattern
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Fill array [0x8000, 0x8001, 0x8002] with values [0xAA, 0xBB, 0xCC]
	cpu.SetHL(0x8000)
	
	// Fill first element
	cpu.A = 0xAA
	cpu.LD_HL_INC_A(mmu)
	assert.Equal(t, uint16(0x8001), cpu.GetHL())
	
	// Fill second element
	cpu.A = 0xBB
	cpu.LD_HL_INC_A(mmu)
	assert.Equal(t, uint16(0x8002), cpu.GetHL())
	
	// Fill third element
	cpu.A = 0xCC
	cpu.LD_HL_INC_A(mmu)
	assert.Equal(t, uint16(0x8003), cpu.GetHL())
	
	// Verify array contents
	assert.Equal(t, uint8(0xAA), mmu.ReadByte(0x8000))
	assert.Equal(t, uint8(0xBB), mmu.ReadByte(0x8001))
	assert.Equal(t, uint8(0xCC), mmu.ReadByte(0x8002))
}

func TestArrayReadPattern(t *testing.T) {
	// Test a realistic array reading pattern
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set up array data
	mmu.WriteByte(0x9000, 0x11)
	mmu.WriteByte(0x9001, 0x22)
	mmu.WriteByte(0x9002, 0x33)
	
	// Read array elements
	cpu.SetHL(0x9000)
	
	// Read first element
	cpu.LD_A_HL_INC(mmu)
	assert.Equal(t, uint8(0x11), cpu.A)
	assert.Equal(t, uint16(0x9001), cpu.GetHL())
	
	// Read second element
	cpu.LD_A_HL_INC(mmu)
	assert.Equal(t, uint8(0x22), cpu.A)
	assert.Equal(t, uint16(0x9002), cpu.GetHL())
	
	// Read third element
	cpu.LD_A_HL_INC(mmu)
	assert.Equal(t, uint8(0x33), cpu.A)
	assert.Equal(t, uint16(0x9003), cpu.GetHL())
}

func TestStringBuildingBackward(t *testing.T) {
	// Test building a string backward using HL-
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Build "HI!" backward starting at 0x8002
	cpu.SetHL(0x8002)
	
	// Store '!' at 0x8002, move to 0x8001
	cpu.A = '!'
	cpu.LD_HL_DEC_A(mmu)
	assert.Equal(t, uint16(0x8001), cpu.GetHL())
	
	// Store 'I' at 0x8001, move to 0x8000
	cpu.A = 'I'
	cpu.LD_HL_DEC_A(mmu)
	assert.Equal(t, uint16(0x8000), cpu.GetHL())
	
	// Store 'H' at 0x8000, move to 0x7FFF
	cpu.A = 'H'
	cpu.LD_HL_DEC_A(mmu)
	assert.Equal(t, uint16(0x7FFF), cpu.GetHL())
	
	// Verify string "HI!" is built correctly
	assert.Equal(t, uint8('H'), mmu.ReadByte(0x8000))
	assert.Equal(t, uint8('I'), mmu.ReadByte(0x8001))
	assert.Equal(t, uint8('!'), mmu.ReadByte(0x8002))
}

func TestWrapAroundBehavior(t *testing.T) {
	// Test HL wrap-around behavior
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	t.Run("Increment wrap-around", func(t *testing.T) {
		cpu.SetHL(0xFFFF)
		cpu.A = 0x99
		
		cpu.LD_HL_INC_A(mmu)
		
		assert.Equal(t, uint8(0x99), mmu.ReadByte(0xFFFF))
		assert.Equal(t, uint16(0x0000), cpu.GetHL()) // Should wrap to 0x0000
	})
	
	t.Run("Decrement wrap-around", func(t *testing.T) {
		cpu.SetHL(0x0000)
		cpu.A = 0x88
		
		cpu.LD_HL_DEC_A(mmu)
		
		assert.Equal(t, uint8(0x88), mmu.ReadByte(0x0000))
		assert.Equal(t, uint16(0xFFFF), cpu.GetHL()) // Should wrap to 0xFFFF
	})
}

func TestFlagsNotAffected(t *testing.T) {
	// Test that these instructions don't affect any flags
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set all flags to known state
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)
	
	cpu.A = 0x42
	cpu.SetHL(0x8000)
	
	// Execute all four instructions
	cpu.LD_HL_INC_A(mmu)
	cpu.LD_A_HL_INC(mmu)
	cpu.LD_HL_DEC_A(mmu)
	cpu.LD_A_HL_DEC(mmu)
	
	// All flags should remain unchanged
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
}
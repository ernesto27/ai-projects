package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === BIT Instruction Tests ===

func TestBIT_0_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test bit 0 set (bit value = 1, Z flag should be false)
	cpu.B = 0x01 // Binary: 00000001, bit 0 = 1
	cycles := cpu.BIT_0_B()
	
	assert.Equal(t, uint8(8), cycles, "BIT 0,B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
	
	// Test bit 0 clear (bit value = 0, Z flag should be true)
	cpu.B = 0xFE // Binary: 11111110, bit 0 = 0
	cycles = cpu.BIT_0_B()
	
	assert.Equal(t, uint8(8), cycles, "BIT 0,B should take 8 cycles")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit is clear")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
	
	// Verify register B is unchanged
	assert.Equal(t, uint8(0xFE), cpu.B, "Register B should be unchanged")
}

func TestBIT_1_C(t *testing.T) {
	cpu := NewCPU()
	
	// Test bit 1 set
	cpu.C = 0x02 // Binary: 00000010, bit 1 = 1
	cycles := cpu.BIT_1_C()
	
	assert.Equal(t, uint8(8), cycles, "BIT 1,C should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
	
	// Test bit 1 clear
	cpu.C = 0xFD // Binary: 11111101, bit 1 = 0
	cycles = cpu.BIT_1_C()
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit is clear")
	assert.Equal(t, uint8(0xFD), cpu.C, "Register C should be unchanged")
}

func TestBIT_7_A(t *testing.T) {
	cpu := NewCPU()
	
	// Test bit 7 set (most significant bit)
	cpu.A = 0x80 // Binary: 10000000, bit 7 = 1
	cycles := cpu.BIT_7_A()
	
	assert.Equal(t, uint8(8), cycles, "BIT 7,A should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
	
	// Test bit 7 clear
	cpu.A = 0x7F // Binary: 01111111, bit 7 = 0
	cycles = cpu.BIT_7_A()
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit is clear")
	assert.Equal(t, uint8(0x7F), cpu.A, "Register A should be unchanged")
}

func TestBIT_0_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x8000)
	
	// Test bit 0 set in memory
	mmu.WriteByte(0x8000, 0x01) // Binary: 00000001, bit 0 = 1
	cycles := cpu.BIT_0_HL(mmu)
	
	assert.Equal(t, uint8(12), cycles, "BIT 0,(HL) should take 12 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
	
	// Test bit 0 clear in memory
	mmu.WriteByte(0x8000, 0xFE) // Binary: 11111110, bit 0 = 0
	cycles = cpu.BIT_0_HL(mmu)
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit is clear")
	
	// Verify memory is unchanged
	assert.Equal(t, uint8(0xFE), mmu.ReadByte(0x8000), "Memory should be unchanged")
}

func TestBIT_7_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x9000)
	
	// Test bit 7 set in memory
	mmu.WriteByte(0x9000, 0x80) // Binary: 10000000, bit 7 = 1
	cycles := cpu.BIT_7_HL(mmu)
	
	assert.Equal(t, uint8(12), cycles, "BIT 7,(HL) should take 12 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit is set")
	
	// Test bit 7 clear in memory
	mmu.WriteByte(0x9000, 0x7F) // Binary: 01111111, bit 7 = 0
	cycles = cpu.BIT_7_HL(mmu)
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit is clear")
}

// === SET Instruction Tests ===

func TestSET_0_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test setting bit 0 when it's already clear
	cpu.B = 0xFE // Binary: 11111110, bit 0 = 0
	cycles := cpu.SET_0_B()
	
	assert.Equal(t, uint8(8), cycles, "SET 0,B should take 8 cycles")
	assert.Equal(t, uint8(0xFF), cpu.B, "Bit 0 should be set to 1")
	
	// Test setting bit 0 when it's already set (should remain set)
	cpu.B = 0x01 // Binary: 00000001, bit 0 = 1
	cycles = cpu.SET_0_B()
	
	assert.Equal(t, uint8(0x01), cpu.B, "Bit 0 should remain set")
	
	// Verify no flags are affected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)
	
	cpu.SET_0_B()
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
}

func TestSET_7_A(t *testing.T) {
	cpu := NewCPU()
	
	// Test setting bit 7 (most significant bit)
	cpu.A = 0x7F // Binary: 01111111, bit 7 = 0
	cycles := cpu.SET_7_A()
	
	assert.Equal(t, uint8(8), cycles, "SET 7,A should take 8 cycles")
	assert.Equal(t, uint8(0xFF), cpu.A, "Bit 7 should be set to 1")
	
	// Test setting bit 7 when it's already set
	cpu.A = 0x80 // Binary: 10000000, bit 7 = 1
	cycles = cpu.SET_7_A()
	
	assert.Equal(t, uint8(0x80), cpu.A, "Bit 7 should remain set")
}

func TestSET_0_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x8000)
	
	// Test setting bit 0 in memory
	mmu.WriteByte(0x8000, 0xFE) // Binary: 11111110, bit 0 = 0
	cycles := cpu.SET_0_HL(mmu)
	
	assert.Equal(t, uint8(16), cycles, "SET 0,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0x8000), "Bit 0 should be set in memory")
	
	// Test setting bit 0 when it's already set
	mmu.WriteByte(0x8000, 0x01) // Binary: 00000001, bit 0 = 1
	cycles = cpu.SET_0_HL(mmu)
	
	assert.Equal(t, uint8(0x01), mmu.ReadByte(0x8000), "Bit 0 should remain set")
}

func TestSET_7_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x9000)
	
	// Test setting bit 7 in memory
	mmu.WriteByte(0x9000, 0x7F) // Binary: 01111111, bit 7 = 0
	cycles := cpu.SET_7_HL(mmu)
	
	assert.Equal(t, uint8(16), cycles, "SET 7,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0x9000), "Bit 7 should be set in memory")
}

// === RES Instruction Tests ===

func TestRES_0_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test resetting bit 0 when it's set
	cpu.B = 0xFF // Binary: 11111111, bit 0 = 1
	cycles := cpu.RES_0_B()
	
	assert.Equal(t, uint8(8), cycles, "RES 0,B should take 8 cycles")
	assert.Equal(t, uint8(0xFE), cpu.B, "Bit 0 should be reset to 0")
	
	// Test resetting bit 0 when it's already clear (should remain clear)
	cpu.B = 0xFE // Binary: 11111110, bit 0 = 0
	cycles = cpu.RES_0_B()
	
	assert.Equal(t, uint8(0xFE), cpu.B, "Bit 0 should remain clear")
	
	// Verify no flags are affected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)
	
	cpu.RES_0_B()
	
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
}

func TestRES_7_A(t *testing.T) {
	cpu := NewCPU()
	
	// Test resetting bit 7 (most significant bit)
	cpu.A = 0xFF // Binary: 11111111, bit 7 = 1
	cycles := cpu.RES_7_A()
	
	assert.Equal(t, uint8(8), cycles, "RES 7,A should take 8 cycles")
	assert.Equal(t, uint8(0x7F), cpu.A, "Bit 7 should be reset to 0")
	
	// Test resetting bit 7 when it's already clear
	cpu.A = 0x7F // Binary: 01111111, bit 7 = 0
	cycles = cpu.RES_7_A()
	
	assert.Equal(t, uint8(0x7F), cpu.A, "Bit 7 should remain clear")
}

func TestRES_0_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x8000)
	
	// Test resetting bit 0 in memory
	mmu.WriteByte(0x8000, 0xFF) // Binary: 11111111, bit 0 = 1
	cycles := cpu.RES_0_HL(mmu)
	
	assert.Equal(t, uint8(16), cycles, "RES 0,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0xFE), mmu.ReadByte(0x8000), "Bit 0 should be reset in memory")
	
	// Test resetting bit 0 when it's already clear
	mmu.WriteByte(0x8000, 0xFE) // Binary: 11111110, bit 0 = 0
	cycles = cpu.RES_0_HL(mmu)
	
	assert.Equal(t, uint8(0xFE), mmu.ReadByte(0x8000), "Bit 0 should remain clear")
}

func TestRES_7_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x9000)
	
	// Test resetting bit 7 in memory
	mmu.WriteByte(0x9000, 0xFF) // Binary: 11111111, bit 7 = 1
	cycles := cpu.RES_7_HL(mmu)
	
	assert.Equal(t, uint8(16), cycles, "RES 7,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0x7F), mmu.ReadByte(0x9000), "Bit 7 should be reset in memory")
}

// === Rotate Instruction Tests ===

func TestRLC_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test rotate left with carry - bit 7 becomes carry and bit 0
	cpu.B = 0x85 // Binary: 10000101
	cycles := cpu.RLC_B()
	
	assert.Equal(t, uint8(8), cycles, "RLC B should take 8 cycles")
	assert.Equal(t, uint8(0x0B), cpu.B, "B should be rotated left: 10000101 -> 00001011")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (result != 0)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be true (old bit 7 was 1)")
	
	// Test rotate with zero result
	cpu.B = 0x00
	cycles = cpu.RLC_B()
	
	assert.Equal(t, uint8(0x00), cpu.B, "B should remain 0")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true (result == 0)")
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be false (old bit 7 was 0)")
}

func TestRLC_C(t *testing.T) {
	cpu := NewCPU()
	
	// Test rotate left with different pattern
	cpu.C = 0xAA // Binary: 10101010
	cycles := cpu.RLC_C()
	
	assert.Equal(t, uint8(8), cycles, "RLC C should take 8 cycles")
	assert.Equal(t, uint8(0x55), cpu.C, "C should be rotated left: 10101010 -> 01010101")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be true (old bit 7 was 1)")
}

func TestRRC_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test rotate right with carry - bit 0 becomes carry and bit 7
	cpu.B = 0xA1 // Binary: 10100001
	cycles := cpu.RRC_B()
	
	assert.Equal(t, uint8(8), cycles, "RRC B should take 8 cycles")
	assert.Equal(t, uint8(0xD0), cpu.B, "B should be rotated right: 10100001 -> 11010000")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (result != 0)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be true (old bit 0 was 1)")
	
	// Test rotate with zero result
	cpu.B = 0x00
	cycles = cpu.RRC_B()
	
	assert.Equal(t, uint8(0x00), cpu.B, "B should remain 0")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true (result == 0)")
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be false (old bit 0 was 0)")
}

func TestRRC_C(t *testing.T) {
	cpu := NewCPU()
	
	// Test rotate right with different pattern
	cpu.C = 0x55 // Binary: 01010101
	cycles := cpu.RRC_C()
	
	assert.Equal(t, uint8(8), cycles, "RRC C should take 8 cycles")
	assert.Equal(t, uint8(0xAA), cpu.C, "C should be rotated right: 01010101 -> 10101010")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be true (old bit 0 was 1)")
}

// === SWAP Instruction Tests ===

func TestSWAP_B(t *testing.T) {
	cpu := NewCPU()
	
	// Test swapping nibbles
	cpu.B = 0xAB // Binary: 10101011 (upper nibble: 1010, lower nibble: 0011)
	cycles := cpu.SWAP_B()
	
	assert.Equal(t, uint8(8), cycles, "SWAP B should take 8 cycles")
	assert.Equal(t, uint8(0xBA), cpu.B, "Nibbles should be swapped: 0xAB -> 0xBA")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (result != 0)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be false")
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be false")
	
	// Test swap with zero result
	cpu.B = 0x00
	cycles = cpu.SWAP_B()
	
	assert.Equal(t, uint8(0x00), cpu.B, "B should remain 0")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true (result == 0)")
	
	// Test swap with symmetric nibbles
	cpu.B = 0x33 // Both nibbles are the same
	cycles = cpu.SWAP_B()
	
	assert.Equal(t, uint8(0x33), cpu.B, "B should remain 0x33 (symmetric)")
}

func TestSWAP_C(t *testing.T) {
	cpu := NewCPU()
	
	// Test swapping nibbles with different pattern
	cpu.C = 0x12 // Binary: 00010010 (upper nibble: 0001, lower nibble: 0010)
	cycles := cpu.SWAP_C()
	
	assert.Equal(t, uint8(8), cycles, "SWAP C should take 8 cycles")
	assert.Equal(t, uint8(0x21), cpu.C, "Nibbles should be swapped: 0x12 -> 0x21")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (result != 0)")
}

func TestSWAP_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	// Set HL to a test address
	cpu.SetHL(0x8000)
	
	// Test swapping nibbles in memory
	mmu.WriteByte(0x8000, 0xCD) // Binary: 11001101 (upper nibble: 1100, lower nibble: 1101)
	cycles := cpu.SWAP_HL(mmu)
	
	assert.Equal(t, uint8(16), cycles, "SWAP (HL) should take 16 cycles")
	assert.Equal(t, uint8(0xDC), mmu.ReadByte(0x8000), "Nibbles should be swapped in memory: 0xCD -> 0xDC")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false (result != 0)")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.False(t, cpu.GetFlag(FlagH), "H flag should be false")
	assert.False(t, cpu.GetFlag(FlagC), "C flag should be false")
	
	// Test swap with zero in memory
	mmu.WriteByte(0x8000, 0x00)
	cycles = cpu.SWAP_HL(mmu)
	
	assert.Equal(t, uint8(0x00), mmu.ReadByte(0x8000), "Memory should remain 0")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true (result == 0)")
}

// === Edge Case Tests ===

func TestBitInstructionsPreserveOtherBits(t *testing.T) {
	cpu := NewCPU()
	
	// Test that SET only affects the target bit
	cpu.B = 0xAA // Binary: 10101010
	cpu.SET_0_B()
	assert.Equal(t, uint8(0xAB), cpu.B, "SET should only change bit 0: 10101010 -> 10101011")
	
	// Test that RES only affects the target bit
	cpu.B = 0x55 // Binary: 01010101
	cpu.RES_0_B()
	assert.Equal(t, uint8(0x54), cpu.B, "RES should only change bit 0: 01010101 -> 01010100")
}

func TestBitInstructionsWithAllPatterns(t *testing.T) {
	cpu := NewCPU()
	
	// Test BIT instruction with all bits set
	cpu.A = 0xFF
	cpu.BIT_7_A()
	assert.False(t, cpu.GetFlag(FlagZ), "BIT 7 should find the bit set")
	
	cpu.BIT_0_A()
	assert.False(t, cpu.GetFlag(FlagZ), "BIT 0 should find the bit set")
	
	// Test BIT instruction with no bits set
	cpu.A = 0x00
	cpu.BIT_7_A()
	assert.True(t, cpu.GetFlag(FlagZ), "BIT 7 should find the bit clear")
	
	cpu.BIT_0_A()
	assert.True(t, cpu.GetFlag(FlagZ), "BIT 0 should find the bit clear")
}

func TestRotateInstructionsEdgeCases(t *testing.T) {
	cpu := NewCPU()
	
	// Test RLC with 0x80 (only bit 7 set)
	cpu.B = 0x80
	cpu.RLC_B()
	assert.Equal(t, uint8(0x01), cpu.B, "RLC 0x80 should become 0x01")
	assert.True(t, cpu.GetFlag(FlagC), "Carry should be set from bit 7")
	
	// Test RRC with 0x01 (only bit 0 set)
	cpu.C = 0x01
	cpu.RRC_C()
	assert.Equal(t, uint8(0x80), cpu.C, "RRC 0x01 should become 0x80")
	assert.True(t, cpu.GetFlag(FlagC), "Carry should be set from bit 0")
}
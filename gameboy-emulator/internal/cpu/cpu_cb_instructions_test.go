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

// === Additional BIT Instruction Tests for Missing Coverage ===

func TestBIT_2_Instructions(t *testing.T) {
	cpu := NewCPU()

	// Test BIT 2 on register B (bit position 2 = value 0x04)
	cpu.B = 0x04 // Binary: 00000100, bit 2 = 1
	cycles := cpu.BIT_2_B()

	assert.Equal(t, uint8(8), cycles, "BIT 2,B should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit 2 is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")

	// Test BIT 2 when bit is clear
	cpu.B = 0xFB // Binary: 11111011, bit 2 = 0
	cycles = cpu.BIT_2_B()

	assert.Equal(t, uint8(8), cycles, "BIT 2,B should take 8 cycles")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit 2 is clear")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
}

func TestBIT_3_Instructions(t *testing.T) {
	cpu := NewCPU()

	// Test BIT 3 on register D (bit position 3 = value 0x08)
	cpu.D = 0x08 // Binary: 00001000, bit 3 = 1
	cycles := cpu.BIT_3_D()

	assert.Equal(t, uint8(8), cycles, "BIT 3,D should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit 3 is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")

	// Test BIT 3 when bit is clear
	cpu.D = 0xF7 // Binary: 11110111, bit 3 = 0
	cycles = cpu.BIT_3_D()

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be true when bit 3 is clear")
}

func TestBIT_5_Instructions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test BIT 5 on register E (bit position 5 = value 0x20)
	cpu.E = 0x20 // Binary: 00100000, bit 5 = 1
	cycles := cpu.BIT_5_E()

	assert.Equal(t, uint8(8), cycles, "BIT 5,E should take 8 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit 5 is set")

	// Test BIT 5 on memory location (HL)
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x20) // Set bit 5 in memory
	cycles = cpu.BIT_5_HL(mmu)

	assert.Equal(t, uint8(12), cycles, "BIT 5,(HL) should take 12 cycles")
	assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be false when bit 5 is set")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be true")
}

// === RES Instruction Tests for Missing Coverage ===

func TestRES_2_Instructions(t *testing.T) {
	cpu := NewCPU()

	// Test RES 2 on register C - clear bit 2
	cpu.C = 0xFF // All bits set
	cycles := cpu.RES_2_C()

	assert.Equal(t, uint8(8), cycles, "RES 2,C should take 8 cycles")
	assert.Equal(t, uint8(0xFB), cpu.C, "Bit 2 should be cleared: 0xFF -> 0xFB")

	// Test RES 2 when bit is already clear
	cpu.C = 0x00 // All bits clear
	cycles = cpu.RES_2_C()

	assert.Equal(t, uint8(0x00), cpu.C, "RES should not affect already clear bit")
}

func TestRES_4_Instructions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test RES 4 on register H - clear bit 4 (value 0x10)
	cpu.H = 0x1F // Binary: 00011111, has bit 4 set
	cycles := cpu.RES_4_H()

	assert.Equal(t, uint8(8), cycles, "RES 4,H should take 8 cycles")
	assert.Equal(t, uint8(0x0F), cpu.H, "Bit 4 should be cleared: 0x1F -> 0x0F")

	// Test RES 4 on memory location (HL)
	cpu.SetHL(0x9000)
	mmu.WriteByte(0x9000, 0xFF) // All bits set
	cycles = cpu.RES_4_HL(mmu)

	assert.Equal(t, uint8(16), cycles, "RES 4,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0xEF), mmu.ReadByte(0x9000), "Bit 4 should be cleared in memory")
}

// === SET Instruction Tests for Missing Coverage ===

func TestSET_3_Instructions(t *testing.T) {
	cpu := NewCPU()

	// Test SET 3 on register A - set bit 3 (value 0x08)
	cpu.A = 0x00 // All bits clear
	cycles := cpu.SET_3_A()

	assert.Equal(t, uint8(8), cycles, "SET 3,A should take 8 cycles")
	assert.Equal(t, uint8(0x08), cpu.A, "Bit 3 should be set: 0x00 -> 0x08")

	// Test SET 3 when bit is already set
	cpu.A = 0xFF // All bits set
	cycles = cpu.SET_3_A()

	assert.Equal(t, uint8(0xFF), cpu.A, "SET should not affect already set bit")
}

func TestSET_5_Instructions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test SET 5 on register L - set bit 5 (value 0x20)
	cpu.L = 0x00 // All bits clear
	cycles := cpu.SET_5_L()

	assert.Equal(t, uint8(8), cycles, "SET 5,L should take 8 cycles")
	assert.Equal(t, uint8(0x20), cpu.L, "Bit 5 should be set: 0x00 -> 0x20")

	// Test SET 5 on memory location (HL)
	cpu.SetHL(0xA000)
	mmu.WriteByte(0xA000, 0x00) // All bits clear
	cycles = cpu.SET_5_HL(mmu)

	assert.Equal(t, uint8(16), cycles, "SET 5,(HL) should take 16 cycles")
	assert.Equal(t, uint8(0x20), mmu.ReadByte(0xA000), "Bit 5 should be set in memory")
}

func TestSET_6_Instructions(t *testing.T) {
	cpu := NewCPU()

	// Test SET 6 on register B - set bit 6 (value 0x40)
	cpu.B = 0x3F // Binary: 00111111, all bits except 6 and 7
	cycles := cpu.SET_6_B()

	assert.Equal(t, uint8(8), cycles, "SET 6,B should take 8 cycles")
	assert.Equal(t, uint8(0x7F), cpu.B, "Bit 6 should be set: 0x3F -> 0x7F")
}

// === Comprehensive CB Instruction Integration Test ===

func TestCBInstructionBitManipulationWorkflow(t *testing.T) {
	cpu := NewCPU()

	// Test a complete bit manipulation workflow
	cpu.A = 0x00 // Start with all bits clear

	// Set some bits: 0, 3, 5, 7
	cpu.SET_0_A()
	cpu.SET_3_A()
	cpu.SET_5_A()
	cpu.SET_7_A()

	expectedValue := uint8(0x01 | 0x08 | 0x20 | 0x80) // bits 0,3,5,7 = 0xA9
	assert.Equal(t, expectedValue, cpu.A, "Multiple SET operations should set correct bits")

	// Test the bits are set
	cpu.BIT_0_A()
	assert.False(t, cpu.GetFlag(FlagZ), "Bit 0 should be found set")

	cpu.BIT_3_A()
	assert.False(t, cpu.GetFlag(FlagZ), "Bit 3 should be found set")

	cpu.BIT_1_A() // This bit should be clear
	assert.True(t, cpu.GetFlag(FlagZ), "Bit 1 should be found clear")

	// Clear some bits
	cpu.RES_3_A()
	cpu.RES_5_A()

	expectedAfterReset := uint8(0x01 | 0x80) // Only bits 0,7 remain = 0x81
	assert.Equal(t, expectedAfterReset, cpu.A, "RES operations should clear correct bits")
}

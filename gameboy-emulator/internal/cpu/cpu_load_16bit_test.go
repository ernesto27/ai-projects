package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLD_BC_nn(t *testing.T) {
	cpu := NewCPU()

	// Test: Load 16-bit immediate into BC
	low := uint8(0x34)
	high := uint8(0x12)

	// Execute instruction
	cycles := cpu.LD_BC_nn(low, high)

	// Verify BC register was loaded correctly
	assert.Equal(t, low, cpu.C, "C register should contain low byte")
	assert.Equal(t, high, cpu.B, "B register should contain high byte")
	assert.Equal(t, uint16(0x1234), cpu.GetBC(), "BC should contain 0x1234")

	// Verify cycle count
	assert.Equal(t, uint8(12), cycles, "LD_BC_nn should take 12 cycles")

	// Test: Different values
	testCases := []struct {
		name     string
		low      uint8
		high     uint8
		expected uint16
	}{
		{"Load 0x0000", 0x00, 0x00, 0x0000},
		{"Load 0xFFFF", 0xFF, 0xFF, 0xFFFF},
		{"Load 0x5500", 0x00, 0x55, 0x5500},
		{"Load 0x00AA", 0xAA, 0x00, 0x00AA},
		{"Load 0x8000", 0x00, 0x80, 0x8000},
		{"Load 0x0001", 0x01, 0x00, 0x0001},
		{"Load 0x8080", 0x80, 0x80, 0x8080},
		{"Load 0x1234", 0x34, 0x12, 0x1234},
		{"Load 0xABCD", 0xCD, 0xAB, 0xABCD},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cycles := cpu.LD_BC_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.C, "C should contain low byte 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.B, "B should contain high byte 0x%02X", tc.high)
			assert.Equal(t, tc.expected, cpu.GetBC(), "BC should contain 0x%04X", tc.expected)
			assert.Equal(t, uint8(12), cycles, "Should always take 12 cycles")
		})
	}

	// Test: Verify other registers are unchanged
	cpu.A = 0x11
	cpu.D = 0x22
	cpu.E = 0x33
	cpu.H = 0x44
	cpu.L = 0x55
	cpu.SP = 0x6666
	cpu.PC = 0x7777

	cpu.LD_BC_nn(0x99, 0x88)

	assert.Equal(t, uint8(0x11), cpu.A, "A register should be unchanged")
	assert.Equal(t, uint8(0x22), cpu.D, "D register should be unchanged")
	assert.Equal(t, uint8(0x33), cpu.E, "E register should be unchanged")
	assert.Equal(t, uint8(0x44), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.L, "L register should be unchanged")
	assert.Equal(t, uint16(0x6666), cpu.SP, "SP should be unchanged")
	assert.Equal(t, uint16(0x7777), cpu.PC, "PC should be unchanged")

	// Verify BC was updated
	assert.Equal(t, uint8(0x99), cpu.C, "C should contain new low byte")
	assert.Equal(t, uint8(0x88), cpu.B, "B should contain new high byte")
	assert.Equal(t, uint16(0x8899), cpu.GetBC(), "BC should contain 0x8899")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_BC_nn(0x77, 0x66)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Little-endian byte order verification
	testEndianness := []struct {
		name  string
		low   uint8
		high  uint8
		bcHex string
	}{
		{"0x1234 = low:0x34, high:0x12", 0x34, 0x12, "0x1234"},
		{"0xABCD = low:0xCD, high:0xAB", 0xCD, 0xAB, "0xABCD"},
		{"0x8000 = low:0x00, high:0x80", 0x00, 0x80, "0x8000"},
		{"0x00FF = low:0xFF, high:0x00", 0xFF, 0x00, "0x00FF"},
	}

	for _, tc := range testEndianness {
		t.Run("Endianness_"+tc.name, func(t *testing.T) {
			cpu.LD_BC_nn(tc.low, tc.high)

			// Verify individual bytes
			assert.Equal(t, tc.low, cpu.C, "C (low byte) should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.B, "B (high byte) should be 0x%02X", tc.high)

			// Verify combined 16-bit value
			combined := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, combined, cpu.GetBC(), "BC should equal %s", tc.bcHex)
		})
	}

	// Test: Edge cases and boundary values
	edgeCases := []struct {
		name string
		low  uint8
		high uint8
	}{
		{"Minimum value", 0x00, 0x00},
		{"Maximum value", 0xFF, 0xFF},
		{"High byte only", 0x00, 0xFF},
		{"Low byte only", 0xFF, 0x00},
		{"All bits set in low", 0xFF, 0x55},
		{"All bits set in high", 0x55, 0xFF},
		{"Alternating pattern", 0xAA, 0x55},
		{"Power of 2 boundary", 0x00, 0x80},
	}

	for _, tc := range edgeCases {
		t.Run("Edge_"+tc.name, func(t *testing.T) {
			initialFlags := cpu.F // Save flags

			cycles := cpu.LD_BC_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.C, "C should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.B, "B should be 0x%02X", tc.high)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should be unchanged")
		})
	}
}

func TestLD_DE_nn(t *testing.T) {
	cpu := NewCPU()

	// Test: Load 16-bit immediate into DE
	low := uint8(0x78)
	high := uint8(0x56)

	// Execute instruction
	cycles := cpu.LD_DE_nn(low, high)

	// Verify DE register was loaded correctly
	assert.Equal(t, low, cpu.E, "E register should contain low byte")
	assert.Equal(t, high, cpu.D, "D register should contain high byte")
	assert.Equal(t, uint16(0x5678), cpu.GetDE(), "DE should contain 0x5678")

	// Verify cycle count
	assert.Equal(t, uint8(12), cycles, "LD_DE_nn should take 12 cycles")

	// Test: Different values
	testCases := []struct {
		name     string
		low      uint8
		high     uint8
		expected uint16
	}{
		{"Load 0x0000", 0x00, 0x00, 0x0000},
		{"Load 0xFFFF", 0xFF, 0xFF, 0xFFFF},
		{"Load 0xAA00", 0x00, 0xAA, 0xAA00},
		{"Load 0x0055", 0x55, 0x00, 0x0055},
		{"Load 0x4000", 0x00, 0x40, 0x4000},
		{"Load 0x0080", 0x80, 0x00, 0x0080},
		{"Load 0xC0C0", 0xC0, 0xC0, 0xC0C0},
		{"Load 0x9876", 0x76, 0x98, 0x9876},
		{"Load 0xDEAD", 0xAD, 0xDE, 0xDEAD},
		{"Load 0xBEEF", 0xEF, 0xBE, 0xBEEF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cycles := cpu.LD_DE_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.E, "E should contain low byte 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.D, "D should contain high byte 0x%02X", tc.high)
			assert.Equal(t, tc.expected, cpu.GetDE(), "DE should contain 0x%04X", tc.expected)
			assert.Equal(t, uint8(12), cycles, "Should always take 12 cycles")
		})
	}

	// Test: Verify other registers are unchanged
	cpu.A = 0xAA
	cpu.B = 0xBB
	cpu.C = 0xCC
	cpu.D = 0xDD
	cpu.H = 0xEE
	cpu.L = 0xFF
	cpu.SP = 0x1111
	cpu.PC = 0x2222

	cpu.LD_DE_nn(0x33, 0x44)

	assert.Equal(t, uint8(0xAA), cpu.A, "A register should be unchanged")
	assert.Equal(t, uint8(0xBB), cpu.B, "B register should be unchanged")
	assert.Equal(t, uint8(0xCC), cpu.C, "C register should be unchanged")
	assert.Equal(t, uint8(0xEE), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0xFF), cpu.L, "L register should be unchanged")
	assert.Equal(t, uint16(0x1111), cpu.SP, "SP should be unchanged")
	assert.Equal(t, uint16(0x2222), cpu.PC, "PC should be unchanged")

	// Verify DE was updated
	assert.Equal(t, uint8(0x33), cpu.E, "E should contain new low byte")
	assert.Equal(t, uint8(0x44), cpu.D, "D should contain new high byte")
	assert.Equal(t, uint16(0x4433), cpu.GetDE(), "DE should contain 0x4433")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_DE_nn(0x99, 0x88)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Little-endian byte order verification
	testEndianness := []struct {
		name  string
		low   uint8
		high  uint8
		deHex string
	}{
		{"0x5678 = low:0x78, high:0x56", 0x78, 0x56, "0x5678"},
		{"0xDEAD = low:0xAD, high:0xDE", 0xAD, 0xDE, "0xDEAD"},
		{"0x4000 = low:0x00, high:0x40", 0x00, 0x40, "0x4000"},
		{"0x00C0 = low:0xC0, high:0x00", 0xC0, 0x00, "0x00C0"},
	}

	for _, tc := range testEndianness {
		t.Run("Endianness_"+tc.name, func(t *testing.T) {
			cpu.LD_DE_nn(tc.low, tc.high)

			// Verify individual bytes
			assert.Equal(t, tc.low, cpu.E, "E (low byte) should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.D, "D (high byte) should be 0x%02X", tc.high)

			// Verify combined 16-bit value
			combined := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, combined, cpu.GetDE(), "DE should equal %s", tc.deHex)
		})
	}

	// Test: Verify BC and DE independence
	cpu.LD_BC_nn(0x11, 0x22) // BC = 0x2211
	cpu.LD_DE_nn(0x33, 0x44) // DE = 0x4433

	// Verify BC wasn't affected by DE load
	assert.Equal(t, uint16(0x2211), cpu.GetBC(), "BC should remain 0x2211")
	assert.Equal(t, uint8(0x11), cpu.C, "C should remain 0x11")
	assert.Equal(t, uint8(0x22), cpu.B, "B should remain 0x22")

	// Verify DE was set correctly
	assert.Equal(t, uint16(0x4433), cpu.GetDE(), "DE should be 0x4433")
	assert.Equal(t, uint8(0x33), cpu.E, "E should be 0x33")
	assert.Equal(t, uint8(0x44), cpu.D, "D should be 0x44")

	// Test: Edge cases and boundary values
	edgeCases := []struct {
		name string
		low  uint8
		high uint8
	}{
		{"Minimum value", 0x00, 0x00},
		{"Maximum value", 0xFF, 0xFF},
		{"High byte only", 0x00, 0xFF},
		{"Low byte only", 0xFF, 0x00},
		{"Alternating bits low", 0xAA, 0x55},
		{"Alternating bits high", 0x55, 0xAA},
		{"Power of 2 low", 0x80, 0x00},
		{"Power of 2 high", 0x00, 0x80},
		{"All ones low", 0xFF, 0x00},
		{"All ones high", 0x00, 0xFF},
	}

	for _, tc := range edgeCases {
		t.Run("Edge_"+tc.name, func(t *testing.T) {
			initialFlags := cpu.F // Save flags

			cycles := cpu.LD_DE_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.E, "E should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.D, "D should be 0x%02X", tc.high)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should be unchanged")
		})
	}
}

func TestLD_HL_nn(t *testing.T) {
	cpu := NewCPU()

	// Test: Load 16-bit immediate into HL
	low := uint8(0xBC)
	high := uint8(0x9A)

	// Execute instruction
	cycles := cpu.LD_HL_nn(low, high)

	// Verify HL register was loaded correctly
	assert.Equal(t, low, cpu.L, "L register should contain low byte")
	assert.Equal(t, high, cpu.H, "H register should contain high byte")
	assert.Equal(t, uint16(0x9ABC), cpu.GetHL(), "HL should contain 0x9ABC")

	// Verify cycle count
	assert.Equal(t, uint8(12), cycles, "LD_HL_nn should take 12 cycles")

	// Test: Different values
	testCases := []struct {
		name     string
		low      uint8
		high     uint8
		expected uint16
	}{
		{"Load 0x0000", 0x00, 0x00, 0x0000},
		{"Load 0xFFFF", 0xFF, 0xFF, 0xFFFF},
		{"Load 0x8000", 0x00, 0x80, 0x8000},
		{"Load 0x007F", 0x7F, 0x00, 0x007F},
		{"Load 0x1000", 0x00, 0x10, 0x1000},
		{"Load 0x00F0", 0xF0, 0x00, 0x00F0},
		{"Load 0x5A5A", 0x5A, 0x5A, 0x5A5A},
		{"Load 0xFEDC", 0xDC, 0xFE, 0xFEDC},
		{"Load 0xCAFE", 0xFE, 0xCA, 0xCAFE},
		{"Load 0xBABE", 0xBE, 0xBA, 0xBABE},
		{"Load 0x1337", 0x37, 0x13, 0x1337},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cycles := cpu.LD_HL_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.L, "L should contain low byte 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.H, "H should contain high byte 0x%02X", tc.high)
			assert.Equal(t, tc.expected, cpu.GetHL(), "HL should contain 0x%04X", tc.expected)
			assert.Equal(t, uint8(12), cycles, "Should always take 12 cycles")
		})
	}

	// Test: Verify other registers are unchanged
	cpu.A = 0x11
	cpu.B = 0x22
	cpu.C = 0x33
	cpu.D = 0x44
	cpu.E = 0x55
	cpu.SP = 0x6666
	cpu.PC = 0x7777

	cpu.LD_HL_nn(0x88, 0x99)

	assert.Equal(t, uint8(0x11), cpu.A, "A register should be unchanged")
	assert.Equal(t, uint8(0x22), cpu.B, "B register should be unchanged")
	assert.Equal(t, uint8(0x33), cpu.C, "C register should be unchanged")
	assert.Equal(t, uint8(0x44), cpu.D, "D register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.E, "E register should be unchanged")
	assert.Equal(t, uint16(0x6666), cpu.SP, "SP should be unchanged")
	assert.Equal(t, uint16(0x7777), cpu.PC, "PC should be unchanged")

	// Verify HL was updated
	assert.Equal(t, uint8(0x88), cpu.L, "L should contain new low byte")
	assert.Equal(t, uint8(0x99), cpu.H, "H should contain new high byte")
	assert.Equal(t, uint16(0x9988), cpu.GetHL(), "HL should contain 0x9988")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_HL_nn(0xAA, 0xBB)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Little-endian byte order verification
	testEndianness := []struct {
		name  string
		low   uint8
		high  uint8
		hlHex string
	}{
		{"0x1234 = low:0x34, high:0x12", 0x34, 0x12, "0x1234"},
		{"0xABCD = low:0xCD, high:0xAB", 0xCD, 0xAB, "0xABCD"},
		{"0x8080 = low:0x80, high:0x80", 0x80, 0x80, "0x8080"},
		{"0x00FF = low:0xFF, high:0x00", 0xFF, 0x00, "0x00FF"},
		{"0xFF00 = low:0x00, high:0xFF", 0x00, 0xFF, "0xFF00"},
	}

	for _, tc := range testEndianness {
		t.Run("Endianness_"+tc.name, func(t *testing.T) {
			cpu.LD_HL_nn(tc.low, tc.high)

			// Verify individual bytes
			assert.Equal(t, tc.low, cpu.L, "L (low byte) should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.H, "H (high byte) should be 0x%02X", tc.high)

			// Verify combined 16-bit value
			combined := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, combined, cpu.GetHL(), "HL should equal %s", tc.hlHex)
		})
	}

	// Test: Verify BC, DE, and HL independence
	cpu.LD_BC_nn(0x11, 0x22) // BC = 0x2211
	cpu.LD_DE_nn(0x33, 0x44) // DE = 0x4433
	cpu.LD_HL_nn(0x55, 0x66) // HL = 0x6655

	// Verify BC wasn't affected
	assert.Equal(t, uint16(0x2211), cpu.GetBC(), "BC should remain 0x2211")
	assert.Equal(t, uint8(0x11), cpu.C, "C should remain 0x11")
	assert.Equal(t, uint8(0x22), cpu.B, "B should remain 0x22")

	// Verify DE wasn't affected
	assert.Equal(t, uint16(0x4433), cpu.GetDE(), "DE should remain 0x4433")
	assert.Equal(t, uint8(0x33), cpu.E, "E should remain 0x33")
	assert.Equal(t, uint8(0x44), cpu.D, "D should remain 0x44")

	// Verify HL was set correctly
	assert.Equal(t, uint16(0x6655), cpu.GetHL(), "HL should be 0x6655")
	assert.Equal(t, uint8(0x55), cpu.L, "L should be 0x55")
	assert.Equal(t, uint8(0x66), cpu.H, "H should be 0x66")

	// Test: Memory-related addresses (common HL use cases)
	memoryAddresses := []struct {
		name string
		low  uint8
		high uint8
		addr uint16
		desc string
	}{
		{"ROM start", 0x00, 0x00, 0x0000, "Start of ROM"},
		{"ROM bank 1", 0x00, 0x40, 0x4000, "Start of ROM bank 1"},
		{"VRAM start", 0x00, 0x80, 0x8000, "Start of VRAM"},
		{"WRAM High", 0x00, 0xD0, 0xD000, "Start of WRAM High"},
		{"Work RAM", 0x00, 0xC0, 0xC000, "Start of Work RAM"},
		{"High RAM", 0x80, 0xFF, 0xFF80, "Start of High RAM"},
		{"Stack area", 0xFE, 0xFF, 0xFFFE, "Default stack pointer"},
	}

	for _, tc := range memoryAddresses {
		t.Run("Memory_"+tc.name, func(t *testing.T) {
			cycles := cpu.LD_HL_nn(tc.low, tc.high)

			assert.Equal(t, tc.addr, cpu.GetHL(), "HL should point to %s (0x%04X)", tc.desc, tc.addr)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
		})
	}

	// Test: Edge cases and boundary values
	edgeCases := []struct {
		name string
		low  uint8
		high uint8
	}{
		{"Minimum value", 0x00, 0x00},
		{"Maximum value", 0xFF, 0xFF},
		{"High byte only", 0x00, 0xFF},
		{"Low byte only", 0xFF, 0x00},
		{"Video RAM start", 0x00, 0x80},
		{"Sprite table", 0x00, 0xFE},
		{"Interrupt vectors", 0x00, 0x00},
		{"Cartridge header", 0x00, 0x01},
		{"Nintendo logo", 0x04, 0x01},
		{"Title area", 0x34, 0x01},
	}

	for _, tc := range edgeCases {
		t.Run("Edge_"+tc.name, func(t *testing.T) {
			initialFlags := cpu.F // Save flags

			cycles := cpu.LD_HL_nn(tc.low, tc.high)

			assert.Equal(t, tc.low, cpu.L, "L should be 0x%02X", tc.low)
			assert.Equal(t, tc.high, cpu.H, "H should be 0x%02X", tc.high)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should be unchanged")
		})
	}
}

func TestLD_SP_nn(t *testing.T) {
	cpu := NewCPU()

	// Test: Load 16-bit immediate into SP
	low := uint8(0xF0)
	high := uint8(0xFE)

	// Execute instruction
	cycles := cpu.LD_SP_nn(low, high)

	// Verify SP register was loaded correctly
	assert.Equal(t, uint16(0xFEF0), cpu.SP, "SP should contain 0xFEF0")

	// Verify cycle count
	assert.Equal(t, uint8(12), cycles, "LD_SP_nn should take 12 cycles")

	// Test: Different values
	testCases := []struct {
		name     string
		low      uint8
		high     uint8
		expected uint16
	}{
		{"Load 0x0000", 0x00, 0x00, 0x0000},
		{"Load 0xFFFF", 0xFF, 0xFF, 0xFFFF},
		{"Load 0x8000", 0x00, 0x80, 0x8000},
		{"Load 0x00FF", 0xFF, 0x00, 0x00FF},
		{"Load 0x4000", 0x00, 0x40, 0x4000},
		{"Load 0x0080", 0x80, 0x00, 0x0080},
		{"Load 0xC000", 0x00, 0xC0, 0xC000},
		{"Load 0x1234", 0x34, 0x12, 0x1234},
		{"Load 0xABCD", 0xCD, 0xAB, 0xABCD},
		{"Load 0xDEAD", 0xAD, 0xDE, 0xDEAD},
		{"Load 0xBEEF", 0xEF, 0xBE, 0xBEEF},
		{"Load 0xFFFE", 0xFE, 0xFF, 0xFFFE}, // Common stack start
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cycles := cpu.LD_SP_nn(tc.low, tc.high)

			assert.Equal(t, tc.expected, cpu.SP, "SP should contain 0x%04X", tc.expected)
			assert.Equal(t, uint8(12), cycles, "Should always take 12 cycles")
		})
	}

	// Test: Verify other registers are unchanged
	cpu.A = 0x11
	cpu.B = 0x22
	cpu.C = 0x33
	cpu.D = 0x44
	cpu.E = 0x55
	cpu.H = 0x66
	cpu.L = 0x77
	cpu.PC = 0x8888

	cpu.LD_SP_nn(0xAA, 0xBB)

	assert.Equal(t, uint8(0x11), cpu.A, "A register should be unchanged")
	assert.Equal(t, uint8(0x22), cpu.B, "B register should be unchanged")
	assert.Equal(t, uint8(0x33), cpu.C, "C register should be unchanged")
	assert.Equal(t, uint8(0x44), cpu.D, "D register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.E, "E register should be unchanged")
	assert.Equal(t, uint8(0x66), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x77), cpu.L, "L register should be unchanged")
	assert.Equal(t, uint16(0x8888), cpu.PC, "PC should be unchanged")

	// Verify SP was updated
	assert.Equal(t, uint16(0xBBAA), cpu.SP, "SP should contain 0xBBAA")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_SP_nn(0x99, 0x77)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Little-endian byte order verification
	testEndianness := []struct {
		name  string
		low   uint8
		high  uint8
		spHex string
	}{
		{"0x1234 = low:0x34, high:0x12", 0x34, 0x12, "0x1234"},
		{"0xABCD = low:0xCD, high:0xAB", 0xCD, 0xAB, "0xABCD"},
		{"0x8000 = low:0x00, high:0x80", 0x00, 0x80, "0x8000"},
		{"0x00FF = low:0xFF, high:0x00", 0xFF, 0x00, "0x00FF"},
		{"0xFFFE = low:0xFE, high:0xFF", 0xFE, 0xFF, "0xFFFE"},
	}

	for _, tc := range testEndianness {
		t.Run("Endianness_"+tc.name, func(t *testing.T) {
			cpu.LD_SP_nn(tc.low, tc.high)

			// Verify combined 16-bit value
			combined := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, combined, cpu.SP, "SP should equal %s", tc.spHex)
		})
	}

	// Test: Stack pointer specific edge cases
	stackEdgeCases := []struct {
		name string
		low  uint8
		high uint8
		desc string
	}{
		{"Stack top", 0xFE, 0xFF, "Common Game Boy stack initialization"},
		{"Stack empty", 0x00, 0x00, "Stack pointer at memory start"},
		{"High RAM", 0x00, 0xFF, "Stack in high RAM area"},
		{"Work RAM", 0x00, 0xC0, "Stack in work RAM"},
		{"Echo RAM", 0x00, 0xE0, "Stack in echo RAM area"},
		{"OAM boundary", 0x00, 0xFE, "Stack near OAM"},
		{"I/O boundary", 0xFF, 0xFE, "Stack at I/O boundary"},
		{"Interrupt boundary", 0x00, 0x80, "Stack at interrupt boundary"},
	}

	for _, tc := range stackEdgeCases {
		t.Run("Stack_"+tc.name, func(t *testing.T) {
			initialFlags := cpu.F // Save flags

			cycles := cpu.LD_SP_nn(tc.low, tc.high)

			expected := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, expected, cpu.SP, "SP should be 0x%04X (%s)", expected, tc.desc)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should be unchanged")
		})
	}

	// Test: Verify SP independence from other register pairs
	cpu.B = 0x12
	cpu.C = 0x34
	cpu.D = 0x56
	cpu.E = 0x78
	cpu.H = 0x9A
	cpu.L = 0xBC

	cpu.LD_SP_nn(0xEF, 0xCD)

	// Verify SP was set correctly
	assert.Equal(t, uint16(0xCDEF), cpu.SP, "SP should be 0xCDEF")

	// Verify other register pairs unchanged
	assert.Equal(t, uint16(0x1234), cpu.GetBC(), "BC should be unchanged")
	assert.Equal(t, uint16(0x5678), cpu.GetDE(), "DE should be unchanged")
	assert.Equal(t, uint16(0x9ABC), cpu.GetHL(), "HL should be unchanged")

	// Test: Boundary value testing
	boundaryTests := []struct {
		name string
		low  uint8
		high uint8
	}{
		{"Minimum value", 0x00, 0x00},
		{"Maximum value", 0xFF, 0xFF},
		{"High byte only", 0x00, 0xFF},
		{"Low byte only", 0xFF, 0x00},
		{"Powers of 2", 0x00, 0x80},
		{"Sign bit set", 0x80, 0x80},
		{"Alternating bits", 0xAA, 0x55},
		{"Common stack start", 0xFE, 0xFF},
	}

	for _, tc := range boundaryTests {
		t.Run("Boundary_"+tc.name, func(t *testing.T) {
			initialFlags := cpu.F // Save flags

			cycles := cpu.LD_SP_nn(tc.low, tc.high)

			expected := uint16(tc.high)<<8 | uint16(tc.low)
			assert.Equal(t, expected, cpu.SP, "SP should be 0x%04X", expected)
			assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should be unchanged")
		})
	}
}

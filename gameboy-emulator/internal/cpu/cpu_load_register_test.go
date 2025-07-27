package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLD_C_A tests the LD C,A instruction
func TestLD_C_A(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from A to C
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.A = value
		cpu.C = 0x99 // Different value in C
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD C,A instruction
		cycles := cpu.LD_C_A()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD C,A should take 4 cycles")

		// C register should now contain A's value
		assert.Equal(t, value, cpu.C, "C register should contain A's value")

		// A register should be unchanged (source remains intact)
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_A_D tests the LD A,D instruction
func TestLD_A_D(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from D to A
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.D = value
		cpu.A = 0x99 // Different value in A
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialD := cpu.D
		initialB := cpu.B
		initialC := cpu.C
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD A,D instruction
		cycles := cpu.LD_A_D()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD A,D should take 4 cycles")

		// A register should now contain D's value
		assert.Equal(t, value, cpu.A, "A register should contain D's value")

		// D register should be unchanged (source remains intact)
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_D_A tests the LD D,A instruction
func TestLD_D_A(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from A to D
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.A = value
		cpu.D = 0x99 // Different value in D
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD D,A instruction
		cycles := cpu.LD_D_A()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD D,A should take 4 cycles")

		// D register should now contain A's value
		assert.Equal(t, value, cpu.D, "D register should contain A's value")

		// A register should be unchanged (source remains intact)
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_B_D tests the LD B,D instruction
func TestLD_B_D(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from D to B
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.D = value
		cpu.B = 0x99 // Different value in B
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD B,D instruction
		cycles := cpu.LD_B_D()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD B,D should take 4 cycles")

		// B register should now contain D's value
		assert.Equal(t, value, cpu.B, "B register should contain D's value")

		// D register should be unchanged (source remains intact)
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_D_B tests the LD D,B instruction
func TestLD_D_B(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from B to D
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.B = value
		cpu.D = 0x99 // Different value in D
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD D,B instruction
		cycles := cpu.LD_D_B()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD D,B should take 4 cycles")

		// D register should now contain B's value
		assert.Equal(t, value, cpu.D, "D register should contain B's value")

		// B register should be unchanged (source remains intact)
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_C_E tests the LD C,E instruction
func TestLD_C_E(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from E to C
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.E = value
		cpu.C = 0x99 // Different value in C
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD C,E instruction
		cycles := cpu.LD_C_E()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD C,E should take 4 cycles")

		// C register should now contain E's value
		assert.Equal(t, value, cpu.C, "C register should contain E's value")

		// E register should be unchanged (source remains intact)
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_A_E tests the LD A,E instruction
func TestLD_A_E(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from E to A
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.E = value
		cpu.A = 0x99 // Different value in A
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialE := cpu.E
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD A,E instruction
		cycles := cpu.LD_A_E()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD A,E should take 4 cycles")

		// A register should now contain E's value
		assert.Equal(t, value, cpu.A, "A register should contain E's value")

		// E register should be unchanged (source remains intact)
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_E_A tests the LD E,A instruction
func TestLD_E_A(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from A to E
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.A = value
		cpu.E = 0x99 // Different value in E
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD E,A instruction
		cycles := cpu.LD_E_A()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD E,A should take 4 cycles")

		// E register should now contain A's value
		assert.Equal(t, value, cpu.E, "E register should contain A's value")

		// A register should be unchanged (source remains intact)
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_B_E tests the LD B,E instruction
func TestLD_B_E(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from E to B
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.E = value
		cpu.B = 0x99 // Different value in B
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialA := cpu.A
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD B,E instruction
		cycles := cpu.LD_B_E()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD B,E should take 4 cycles")

		// B register should now contain E's value
		assert.Equal(t, value, cpu.B, "B register should contain E's value")

		// E register should be unchanged (source remains intact)
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")

		// All other registers and flags should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

func TestLD_HL_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test: Store A to memory at HL
	cpu.A = 0x42
	cpu.SetHL(0x8000)

	// Execute instruction
	cycles := cpu.LD_HL_A(mmu)

	// Verify memory was written correctly
	storedValue := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x42), storedValue, "Memory at HL should contain value from A")

	// Verify cycle count
	assert.Equal(t, uint8(8), cycles, "LD_HL_A should take 8 cycles")

	// Test: Different values and addresses
	testCases := []struct {
		name     string
		aValue   uint8
		hlAddr   uint16
		expected uint8
	}{
		{"Store 0x00", 0x00, 0x8001, 0x00},
		{"Store 0xFF", 0xFF, 0x8002, 0xFF},
		{"Store 0x55", 0x55, 0x9000, 0x55},
		{"Store 0xAA", 0xAA, 0x9FFF, 0xAA},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.A = tc.aValue
			cpu.SetHL(tc.hlAddr)

			cycles := cpu.LD_HL_A(mmu)

			storedValue := mmu.ReadByte(tc.hlAddr)
			assert.Equal(t, tc.expected, storedValue,
				"Memory at 0x%04X should contain 0x%02X", tc.hlAddr, tc.expected)
			assert.Equal(t, uint8(8), cycles, "Should always take 8 cycles")
		})
	}

	// Test: Verify A register is unchanged
	originalA := uint8(0x77)
	cpu.A = originalA
	cpu.SetHL(0x8500)

	cpu.LD_HL_A(mmu)

	assert.Equal(t, originalA, cpu.A, "A register should remain unchanged")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_HL_A(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")
}

func TestLD_A_BC(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test: Load A from memory at BC
	cpu.SetBC(0x8000)
	mmu.WriteByte(0x8000, 0x42)

	// Execute instruction
	cycles := cpu.LD_A_BC(mmu)

	// Verify A register was loaded correctly
	assert.Equal(t, uint8(0x42), cpu.A, "A register should contain value from memory at BC")

	// Verify cycle count
	assert.Equal(t, uint8(8), cycles, "LD_A_BC should take 8 cycles")

	// Test: Different values and addresses
	testCases := []struct {
		name     string
		bcAddr   uint16
		memValue uint8
		expected uint8
	}{
		{"Load 0x00", 0x8001, 0x00, 0x00},
		{"Load 0xFF", 0x8002, 0xFF, 0xFF},
		{"Load 0x55", 0x9000, 0x55, 0x55},
		{"Load 0xAA", 0x9FFF, 0xAA, 0xAA},
		{"Load from WRAM area", 0xC100, 0x33, 0x33},
		{"Load from high memory", 0xFF80, 0x77, 0x77},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.SetBC(tc.bcAddr)
			mmu.WriteByte(tc.bcAddr, tc.memValue)

			cycles := cpu.LD_A_BC(mmu)

			assert.Equal(t, tc.expected, cpu.A,
				"A register should contain 0x%02X from memory at BC (0x%04X)", tc.expected, tc.bcAddr)
			assert.Equal(t, uint8(8), cycles, "Should always take 8 cycles")
		})
	}

	// Test: Verify BC register is unchanged
	originalBC := uint16(0x8500)
	cpu.SetBC(originalBC)
	mmu.WriteByte(originalBC, 0x99)

	cpu.LD_A_BC(mmu)

	assert.Equal(t, originalBC, cpu.GetBC(), "BC register should remain unchanged")
	// Test: Verify other registers are unchanged
	cpu.D = 0x33
	cpu.E = 0x44
	cpu.H = 0x55
	cpu.L = 0x66

	cpu.SetBC(0x8600)
	mmu.WriteByte(0x8600, 0x77)

	cpu.LD_A_BC(mmu)

	assert.Equal(t, uint8(0x77), cpu.A, "A should be loaded with memory value")
	assert.Equal(t, uint16(0x8600), cpu.GetBC(), "BC register should remain unchanged")
	assert.Equal(t, uint8(0x33), cpu.D, "D register should be unchanged")
	assert.Equal(t, uint8(0x44), cpu.E, "E register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x66), cpu.L, "L register should be unchanged")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_A_BC(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")
}

func TestLD_A_DE(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test: Load A from memory at DE
	cpu.SetDE(0x8000)
	mmu.WriteByte(0x8000, 0x42)

	// Execute instruction
	cycles := cpu.LD_A_DE(mmu)

	// Verify A register was loaded correctly
	assert.Equal(t, uint8(0x42), cpu.A, "A register should contain value from memory at DE")

	// Verify cycle count
	assert.Equal(t, uint8(8), cycles, "LD_A_DE should take 8 cycles")

	// Test: Different values and addresses
	testCases := []struct {
		name     string
		deAddr   uint16
		memValue uint8
		expected uint8
	}{
		{"Load 0x00", 0x8001, 0x00, 0x00},
		{"Load 0xFF", 0x8002, 0xFF, 0xFF},
		{"Load 0x55", 0x9000, 0x55, 0x55},
		{"Load 0xAA", 0x9FFF, 0xAA, 0xAA},
		{"Load from WRAM area", 0xC200, 0x88, 0x88},
		{"Load from high memory", 0xFF90, 0x99, 0x99},
		{"Load from WRAM high", 0xD000, 0xCC, 0xCC},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.SetDE(tc.deAddr)
			mmu.WriteByte(tc.deAddr, tc.memValue)

			cycles := cpu.LD_A_DE(mmu)

			assert.Equal(t, tc.expected, cpu.A,
				"A register should contain 0x%02X from memory at DE (0x%04X)", tc.expected, tc.deAddr)
			assert.Equal(t, uint8(8), cycles, "Should always take 8 cycles")
		})
	}

	// Test: Verify DE register is unchanged
	originalDE := uint16(0x8500)
	cpu.SetDE(originalDE)
	mmu.WriteByte(originalDE, 0x99)

	cpu.LD_A_DE(mmu)

	assert.Equal(t, originalDE, cpu.GetDE(), "DE register should remain unchanged")

	// Test: Verify other registers are unchanged
	cpu.B = 0x11
	cpu.C = 0x22
	cpu.H = 0x55
	cpu.L = 0x66

	cpu.SetDE(0x8600)
	mmu.WriteByte(0x8600, 0x77)

	cpu.LD_A_DE(mmu)

	assert.Equal(t, uint8(0x77), cpu.A, "A should be loaded with memory value")
	assert.Equal(t, uint16(0x8600), cpu.GetDE(), "DE register should remain unchanged")
	assert.Equal(t, uint8(0x11), cpu.B, "B register should be unchanged")
	assert.Equal(t, uint8(0x22), cpu.C, "C register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x66), cpu.L, "L register should be unchanged")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_A_DE(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Load from different memory regions in sequence
	memoryRegions := []struct {
		name    string
		address uint16
		value   uint8
	}{
		// Skip ROM and External RAM as they route to cartridge
		{"VRAM", 0x8000, 0x03},
		{"Work RAM", 0xC000, 0x05},
		{"Work RAM High", 0xD000, 0x04},
		{"High RAM", 0xFF80, 0x06},
		{"I/O Register", 0xFF40, 0x07},
	}

	for _, region := range memoryRegions {
		t.Run("Sequential_"+region.name, func(t *testing.T) {
			cpu.SetDE(region.address)
			mmu.WriteByte(region.address, region.value)

			cycles := cpu.LD_A_DE(mmu)

			assert.Equal(t, region.value, cpu.A,
				"Should load 0x%02X from %s at 0x%04X", region.value, region.name, region.address)
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		})
	}
}

func TestLD_BC_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test: Store A to memory at BC
	cpu.A = 0x42
	cpu.SetBC(0x8000)

	// Execute instruction
	cycles := cpu.LD_BC_A(mmu)

	// Verify memory was written correctly
	storedValue := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x42), storedValue, "Memory at BC should contain value from A")

	// Verify cycle count
	assert.Equal(t, uint8(8), cycles, "LD_BC_A should take 8 cycles")

	// Test: Different values and addresses
	testCases := []struct {
		name     string
		aValue   uint8
		bcAddr   uint16
		expected uint8
	}{
		{"Store 0x00", 0x00, 0x8001, 0x00},
		{"Store 0xFF", 0xFF, 0x8002, 0xFF},
		{"Store 0x55", 0x55, 0x9000, 0x55},
		{"Store 0xAA", 0xAA, 0x9FFF, 0xAA},
		{"Store to Work RAM", 0x33, 0xC000, 0x33},
		{"Store to High RAM", 0x77, 0xFF80, 0x77},
		{"Store to WRAM high", 0xCC, 0xD000, 0xCC},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.A = tc.aValue
			cpu.SetBC(tc.bcAddr)

			cycles := cpu.LD_BC_A(mmu)

			storedValue := mmu.ReadByte(tc.bcAddr)
			assert.Equal(t, tc.expected, storedValue,
				"Memory at BC (0x%04X) should contain 0x%02X", tc.bcAddr, tc.expected)
			assert.Equal(t, uint8(8), cycles, "Should always take 8 cycles")
		})
	}

	// Test: Verify A register is unchanged
	originalA := uint8(0x77)
	cpu.A = originalA
	cpu.SetBC(0x8500)

	cpu.LD_BC_A(mmu)

	assert.Equal(t, originalA, cpu.A, "A register should remain unchanged")

	// Test: Verify BC register is unchanged
	originalBC := uint16(0x8600)
	cpu.SetBC(originalBC)
	cpu.A = 0x88

	cpu.LD_BC_A(mmu)

	assert.Equal(t, originalBC, cpu.GetBC(), "BC register should remain unchanged")

	// Test: Verify other registers are unchanged
	cpu.D = 0x33
	cpu.E = 0x44
	cpu.H = 0x55
	cpu.L = 0x66

	cpu.A = 0x99
	cpu.SetBC(0x8700)

	cpu.LD_BC_A(mmu)

	assert.Equal(t, uint8(0x99), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint16(0x8700), cpu.GetBC(), "BC register should remain unchanged")
	assert.Equal(t, uint8(0x33), cpu.D, "D register should be unchanged")
	assert.Equal(t, uint8(0x44), cpu.E, "E register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x66), cpu.L, "L register should be unchanged")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_BC_A(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")
}

func TestLD_DE_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test: Store A to memory at DE
	cpu.A = 0x42
	cpu.SetDE(0x8000)

	// Execute instruction
	cycles := cpu.LD_DE_A(mmu)

	// Verify memory was written correctly
	storedValue := mmu.ReadByte(0x8000)
	assert.Equal(t, uint8(0x42), storedValue, "Memory at DE should contain value from A")

	// Verify cycle count
	assert.Equal(t, uint8(8), cycles, "LD_DE_A should take 8 cycles")

	// Test: Different values and addresses
	testCases := []struct {
		name     string
		aValue   uint8
		deAddr   uint16
		expected uint8
	}{
		{"Store 0x00", 0x00, 0x8001, 0x00},
		{"Store 0xFF", 0xFF, 0x8002, 0xFF},
		{"Store 0x55", 0x55, 0x9000, 0x55},
		{"Store 0xAA", 0xAA, 0x9FFF, 0xAA},
		{"Store to Work RAM", 0x33, 0xC000, 0x33},
		{"Store to High RAM", 0x77, 0xFF80, 0x77},
		{"Store to WRAM high", 0xCC, 0xD000, 0xCC},
		{"Store to Echo RAM", 0x88, 0xE000, 0x88},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu.A = tc.aValue
			cpu.SetDE(tc.deAddr)

			cycles := cpu.LD_DE_A(mmu)

			storedValue := mmu.ReadByte(tc.deAddr)
			assert.Equal(t, tc.expected, storedValue,
				"Memory at DE (0x%04X) should contain 0x%02X", tc.deAddr, tc.expected)
			assert.Equal(t, uint8(8), cycles, "Should always take 8 cycles")
		})
	}

	// Test: Verify A register is unchanged
	originalA := uint8(0x77)
	cpu.A = originalA
	cpu.SetDE(0x8500)

	cpu.LD_DE_A(mmu)

	assert.Equal(t, originalA, cpu.A, "A register should remain unchanged")

	// Test: Verify DE register is unchanged
	originalDE := uint16(0x8600)
	cpu.SetDE(originalDE)
	cpu.A = 0x88

	cpu.LD_DE_A(mmu)

	assert.Equal(t, originalDE, cpu.GetDE(), "DE register should remain unchanged")

	// Test: Verify other registers are unchanged
	cpu.B = 0x11
	cpu.C = 0x22
	cpu.H = 0x55
	cpu.L = 0x66

	cpu.A = 0x99
	cpu.SetDE(0x8700)

	cpu.LD_DE_A(mmu)

	assert.Equal(t, uint8(0x99), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint16(0x8700), cpu.GetDE(), "DE register should remain unchanged")
	assert.Equal(t, uint8(0x11), cpu.B, "B register should be unchanged")
	assert.Equal(t, uint8(0x22), cpu.C, "C register should be unchanged")
	assert.Equal(t, uint8(0x55), cpu.H, "H register should be unchanged")
	assert.Equal(t, uint8(0x66), cpu.L, "L register should be unchanged")

	// Test: Verify flags are unaffected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_DE_A(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should be unaffected")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should be unaffected")

	// Test: Store to different memory regions sequentially
	memoryRegions := []struct {
		name    string
		address uint16
		value   uint8
	}{
		{"Work RAM", 0xC000, 0x01},
		{"VRAM", 0x8000, 0x02},
		{"WRAM High", 0xD000, 0x03},
		{"High RAM", 0xFF80, 0x04},
		{"Echo RAM", 0xE000, 0x05},
		{"Work RAM Bank 1", 0xD000, 0x06},
	}

	for _, region := range memoryRegions {
		t.Run("Sequential_"+region.name, func(t *testing.T) {
			cpu.A = region.value
			cpu.SetDE(region.address)

			cycles := cpu.LD_DE_A(mmu)

			storedValue := mmu.ReadByte(region.address)
			assert.Equal(t, region.value, storedValue,
				"Should store 0x%02X to %s at 0x%04X", region.value, region.name, region.address)
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		})
	}

	// Test: Complete memory operation sequence
	testSequence := []struct {
		addr  uint16
		value uint8
	}{
		{0x8000, 0x11},
		{0xC000, 0x22},
		{0xD000, 0x33},
		{0xFF80, 0x44},
	}

	for i, step := range testSequence {
		t.Run(fmt.Sprintf("Sequence_Step_%d", i+1), func(t *testing.T) {
			cpu.A = step.value
			cpu.SetDE(step.addr)

			cycles := cpu.LD_DE_A(mmu)

			storedValue := mmu.ReadByte(step.addr)
			assert.Equal(t, step.value,
				storedValue, "Step %d: Memory at 0x%04X should contain 0x%02X", i+1, step.addr, step.value)
			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
		})
	}
}

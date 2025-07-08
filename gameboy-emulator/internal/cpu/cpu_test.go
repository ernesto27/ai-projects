package cpu

import (
	"testing"

	"gameboy-emulator/internal/memory"

	"github.com/stretchr/testify/assert"
)

// TestNewCPU tests CPU initialization
func TestNewCPU(t *testing.T) {
	cpu := NewCPU()

	// Test initial register values (Game Boy boot state)
	assert.Equal(t, uint8(0x01), cpu.A, "A register should be 0x01")
	assert.Equal(t, uint8(0xB0), cpu.F, "F register should be 0xB0")
	assert.Equal(t, uint8(0x00), cpu.B, "B register should be 0x00")
	assert.Equal(t, uint8(0x13), cpu.C, "C register should be 0x13")
	assert.Equal(t, uint8(0x00), cpu.D, "D register should be 0x00")
	assert.Equal(t, uint8(0xD8), cpu.E, "E register should be 0xD8")
	assert.Equal(t, uint8(0x01), cpu.H, "H register should be 0x01")
	assert.Equal(t, uint8(0x4D), cpu.L, "L register should be 0x4D")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be 0xFFFE")
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should be 0x0100")
	assert.False(t, cpu.Halted, "CPU should not be halted")
	assert.False(t, cpu.Stopped, "CPU should not be stopped")
}

// TestAFRegisterPair tests AF register pair operations
func TestAFRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting AF pair
	cpu.A = 0x12
	cpu.F = 0x34
	assert.Equal(t, uint16(0x1234), cpu.GetAF(), "AF pair should combine A and F registers")

	// Test setting AF pair
	cpu.SetAF(0x5678)
	assert.Equal(t, uint8(0x56), cpu.A, "A should be set to high byte")
	assert.Equal(t, uint8(0x78), cpu.F, "F should be set to low byte")
}

// TestBCRegisterPair tests BC register pair operations
func TestBCRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting BC pair
	cpu.B = 0xAB
	cpu.C = 0xCD
	assert.Equal(t, uint16(0xABCD), cpu.GetBC(), "BC pair should combine B and C registers")

	// Test setting BC pair
	cpu.SetBC(0x1234)
	assert.Equal(t, uint8(0x12), cpu.B, "B should be set to high byte")
	assert.Equal(t, uint8(0x34), cpu.C, "C should be set to low byte")
}

// TestDERegisterPair tests DE register pair operations
func TestDERegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting DE pair
	cpu.D = 0xEF
	cpu.E = 0x01
	assert.Equal(t, uint16(0xEF01), cpu.GetDE(), "DE pair should combine D and E registers")

	// Test setting DE pair
	cpu.SetDE(0x9876)
	assert.Equal(t, uint8(0x98), cpu.D, "D should be set to high byte")
	assert.Equal(t, uint8(0x76), cpu.E, "E should be set to low byte")
}

// TestHLRegisterPair tests HL register pair operations
func TestHLRegisterPair(t *testing.T) {
	cpu := NewCPU()

	// Test getting HL pair
	cpu.H = 0x42
	cpu.L = 0x24
	assert.Equal(t, uint16(0x4224), cpu.GetHL(), "HL pair should combine H and L registers")

	// Test setting HL pair
	cpu.SetHL(0xBEEF)
	assert.Equal(t, uint8(0xBE), cpu.H, "H should be set to high byte")
	assert.Equal(t, uint8(0xEF), cpu.L, "L should be set to low byte")
}

// TestFlagOperations tests flag register operations
func TestFlagOperations(t *testing.T) {
	cpu := NewCPU()

	// Test setting flags
	cpu.SetFlag(FlagZ, true)
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")

	cpu.SetFlag(FlagN, true)
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")

	cpu.SetFlag(FlagH, true)
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set")

	cpu.SetFlag(FlagC, true)
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")

	// Test clearing flags
	cpu.SetFlag(FlagZ, false)
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be cleared")

	cpu.SetFlag(FlagN, false)
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be cleared")

	cpu.SetFlag(FlagH, false)
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be cleared")

	cpu.SetFlag(FlagC, false)
	assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be cleared")
}

// TestFlagConstants tests that flag constants have correct values
func TestFlagConstants(t *testing.T) {
	assert.Equal(t, 0x80, int(FlagZ), "Zero flag constant should be 0x80")
	assert.Equal(t, 0x40, int(FlagN), "Subtract flag constant should be 0x40")
	assert.Equal(t, 0x20, int(FlagH), "Half-carry flag constant should be 0x20")
	assert.Equal(t, 0x10, int(FlagC), "Carry flag constant should be 0x10")
}

// TestMultipleFlags tests setting multiple flags at once
func TestMultipleFlags(t *testing.T) {
	cpu := NewCPU()

	// Start with a clean flag register for this test
	cpu.F = 0x00

	// Set multiple flags
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagC, true)

	// Check both are set
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be set")

	// Check other flags are not affected
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")
}

// TestRegisterPairBoundaries tests edge cases with register pairs
func TestRegisterPairBoundaries(t *testing.T) {
	cpu := NewCPU()

	// Test maximum values
	cpu.SetAF(0xFFFF)
	assert.Equal(t, uint8(0xFF), cpu.A, "A should be set to 0xFF")
	assert.Equal(t, uint8(0xFF), cpu.F, "F should be set to 0xFF")

	// Test minimum values
	cpu.SetBC(0x0000)
	assert.Equal(t, uint8(0x00), cpu.B, "B should be set to 0x00")
	assert.Equal(t, uint8(0x00), cpu.C, "C should be set to 0x00")

	// Test mixed values
	cpu.SetDE(0x00FF)
	assert.Equal(t, uint8(0x00), cpu.D, "D should be set to 0x00")
	assert.Equal(t, uint8(0xFF), cpu.E, "E should be set to 0xFF")

	cpu.SetHL(0xFF00)
	assert.Equal(t, uint8(0xFF), cpu.H, "H should be set to 0xFF")
	assert.Equal(t, uint8(0x00), cpu.L, "L should be set to 0x00")
}

// TestCPUReset tests CPU reset functionality
func TestCPUReset(t *testing.T) {
	cpu := NewCPU()

	// Modify CPU state
	cpu.A = 0xFF
	cpu.PC = 0x1234
	cpu.SP = 0x5678
	cpu.Halted = true
	cpu.Stopped = true

	// Reset CPU
	cpu.Reset()

	// Verify reset to initial state
	assert.Equal(t, uint8(0x01), cpu.A, "A should be reset to 0x01")
	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should be reset to 0x0100")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be reset to 0xFFFE")
	assert.False(t, cpu.Halted, "CPU should not be halted after reset")
	assert.False(t, cpu.Stopped, "CPU should not be stopped after reset")
}

// TestNOPInstruction tests the NOP instruction
func TestNOPInstruction(t *testing.T) {
	cpu := NewCPU()

	// Store initial state
	initialPC := cpu.PC
	initialA := cpu.A
	initialF := cpu.F
	initialSP := cpu.SP

	// Execute NOP instruction
	cycles := cpu.NOP()

	// NOP should take 4 cycles
	assert.Equal(t, uint8(4), cycles, "NOP should take 4 cycles")

	// NOP should not change any registers
	assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged after NOP")
	assert.Equal(t, initialA, cpu.A, "A should be unchanged after NOP")
	assert.Equal(t, initialF, cpu.F, "F should be unchanged after NOP")
	assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged after NOP")
}

// TestLD_A_n tests the LD A,n instruction
func TestLD_A_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into A
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialF := cpu.F
		initialBC := cpu.GetBC()
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD A,n instruction
		cycles := cpu.LD_A_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD A,n should take 8 cycles")

		// A register should contain the loaded value
		assert.Equal(t, value, cpu.A, "A register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialBC, cpu.GetBC(), "BC register pair should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestINC_A tests the INC A instruction with various flag conditions
func TestINC_A(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.A = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_A()

	assert.Equal(t, uint8(4), cycles, "INC A should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.A, "A should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.A = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_A()

	assert.Equal(t, uint8(0x00), cpu.A, "A should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.A = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_A()

	assert.Equal(t, uint8(0x10), cpu.A, "A should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: No half-carry (0x0E -> 0x0F)
	cpu.A = 0x0E
	cpu.F = 0x00 // Clear all flags
	cpu.INC_A()

	assert.Equal(t, uint8(0x0F), cpu.A, "A should increment from 0x0E to 0x0F")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear on 0x0E->0x0F")

	// Test case 5: Carry flag preservation (INC A doesn't affect carry)
	cpu.A = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_A()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC A")
}

// TestDEC_A tests the DEC A instruction with various flag conditions
func TestDEC_A(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.A = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_A()

	if cycles != 4 {
		t.Errorf("Expected DEC A to take 4 cycles, got %d", cycles)
	}
	if cpu.A != 0x41 {
		t.Errorf("Expected A=0x41 after DEC A, got A=0x%02X", cpu.A)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after normal decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC A")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after normal decrement")
	}

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.A = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_A()

	if cpu.A != 0x00 {
		t.Errorf("Expected A=0x00 after 0x01 decrement, got A=0x%02X", cpu.A)
	}
	if !cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be set after 0x01->0x00 decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC A")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after 0x01->0x00 decrement")
	}

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.A = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_A()

	if cpu.A != 0xFF {
		t.Errorf("Expected A=0xFF after 0x00 decrement, got A=0x%02X", cpu.A)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x00->0xFF decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC A")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0x00->0xFF decrement")
	}

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.A = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_A()

	if cpu.A != 0x0F {
		t.Errorf("Expected A=0x0F after 0x10 decrement, got A=0x%02X", cpu.A)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x10->0x0F decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC A")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0x10->0x0F decrement")
	}

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.A = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_A()

	if cpu.A != 0x0E {
		t.Errorf("Expected A=0x0E after 0x0F decrement, got A=0x%02X", cpu.A)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x0F->0x0E decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC A")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after 0x0F->0x0E decrement")
	}

	// Test case 6: Carry flag preservation (DEC A doesn't affect carry)
	cpu.A = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_A()

	if !cpu.GetFlag(FlagC) {
		t.Error("Expected Carry flag to be preserved after DEC A")
	}
}

// TestLD_B_n tests the LD B,n instruction
func TestLD_B_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into B
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialF := cpu.F
		initialC := cpu.C
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD B,n instruction
		cycles := cpu.LD_B_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD B,n should take 8 cycles")

		// B register should contain the loaded value
		assert.Equal(t, value, cpu.B, "B register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestINC_B tests the INC B instruction with various flag conditions
func TestINC_B(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.B = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_B()

	if cycles != 4 {
		t.Errorf("Expected INC B to take 4 cycles, got %d", cycles)
	}
	if cpu.B != 0x43 {
		t.Errorf("Expected B=0x43 after INC B, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after normal increment")
	}
	if cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be clear after INC B")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after normal increment")
	}

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.B = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_B()

	if cpu.B != 0x00 {
		t.Errorf("Expected B=0x00 after 0xFF increment, got B=0x%02X", cpu.B)
	}
	if !cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be set after 0xFF->0x00 increment")
	}
	if cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be clear after INC B")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0xFF->0x00 increment")
	}

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.B = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_B()

	if cpu.B != 0x10 {
		t.Errorf("Expected B=0x10 after 0x0F increment, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x0F->0x10 increment")
	}
	if cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be clear after INC B")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0x0F->0x10 increment")
	}

	// Test case 4: Other registers unchanged
	cpu.A = 0x99
	cpu.C = 0x88
	cpu.B = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	initialA := cpu.A
	initialC := cpu.C

	cpu.INC_B()

	// Other registers should be unchanged
	if cpu.A != initialA {
		t.Errorf("Expected A unchanged after INC B, got A=0x%02X, expected 0x%02X", cpu.A, initialA)
	}
	if cpu.C != initialC {
		t.Errorf("Expected C unchanged after INC B, got C=0x%02X, expected 0x%02X", cpu.C, initialC)
	}
	if !cpu.GetFlag(FlagC) {
		t.Error("Expected Carry flag to be preserved after INC B")
	}
}

// TestDEC_B tests the DEC B instruction with various flag conditions
func TestDEC_B(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.B = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_B()

	if cycles != 4 {
		t.Errorf("Expected DEC B to take 4 cycles, got %d", cycles)
	}
	if cpu.B != 0x41 {
		t.Errorf("Expected B=0x41 after DEC B, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after normal decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC B")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after normal decrement")
	}

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.B = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_B()

	if cpu.B != 0x00 {
		t.Errorf("Expected B=0x00 after 0x01 decrement, got B=0x%02X", cpu.B)
	}
	if !cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be set after 0x01->0x00 decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC B")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after 0x01->0x00 decrement")
	}

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.B = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_B()

	if cpu.B != 0xFF {
		t.Errorf("Expected B=0xFF after 0x00 decrement, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x00->0xFF decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC B")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0x00->0xFF decrement")
	}

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.B = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_B()

	if cpu.B != 0x0F {
		t.Errorf("Expected B=0x0F after 0x10 decrement, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x10->0x0F decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC B")
	}
	if !cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be set after 0x10->0x0F decrement")
	}

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.B = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_B()

	if cpu.B != 0x0E {
		t.Errorf("Expected B=0x0E after 0x0F decrement, got B=0x%02X", cpu.B)
	}
	if cpu.GetFlag(FlagZ) {
		t.Error("Expected Zero flag to be clear after 0x0F->0x0E decrement")
	}
	if !cpu.GetFlag(FlagN) {
		t.Error("Expected Subtract flag to be set after DEC B")
	}
	if cpu.GetFlag(FlagH) {
		t.Error("Expected Half-carry flag to be clear after 0x0F->0x0E decrement")
	}

	// Test case 6: Carry flag preservation (DEC B doesn't affect carry)
	cpu.B = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_B()

	if !cpu.GetFlag(FlagC) {
		t.Error("Expected Carry flag to be preserved after DEC B")
	}
}

// TestLD_A_B tests the LD A,B instruction
func TestLD_A_B(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from B to A
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.B = value
		cpu.A = 0x99 // Different value in A
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC
		initialB := cpu.B

		// Execute LD A,B instruction
		cycles := cpu.LD_A_B()

		// Should take 4 cycles
		if cycles != 4 {
			t.Errorf("Expected LD A,B to take 4 cycles, got %d", cycles)
		}

		// A register should now contain B's value
		if cpu.A != value {
			t.Errorf("Expected A=0x%02X after LD A,B, got A=0x%02X", value, cpu.A)
		}

		// B register should be unchanged (source remains intact)
		if cpu.B != initialB {
			t.Errorf("Expected B unchanged after LD A,B, got B=0x%02X, expected 0x%02X", cpu.B, initialB)
		}

		// All other registers and flags should be unchanged
		if cpu.C != initialC {
			t.Errorf("Expected C unchanged after LD A,B, got C=0x%02X, expected 0x%02X", cpu.C, initialC)
		}
		if cpu.D != initialD {
			t.Errorf("Expected D unchanged after LD A,B, got D=0x%02X, expected 0x%02X", cpu.D, initialD)
		}
		if cpu.E != initialE {
			t.Errorf("Expected E unchanged after LD A,B, got E=0x%02X, expected 0x%02X", cpu.E, initialE)
		}
		if cpu.H != initialH {
			t.Errorf("Expected H unchanged after LD A,B, got H=0x%02X, expected 0x%02X", cpu.H, initialH)
		}
		if cpu.L != initialL {
			t.Errorf("Expected L unchanged after LD A,B, got L=0x%02X, expected 0x%02X", cpu.L, initialL)
		}
		if cpu.F != initialF {
			t.Errorf("Expected F unchanged after LD A,B, got F=0x%02X, expected 0x%02X", cpu.F, initialF)
		}
		if cpu.SP != initialSP {
			t.Errorf("Expected SP unchanged after LD A,B, got SP=0x%04X, expected 0x%04X", cpu.SP, initialSP)
		}
		if cpu.PC != initialPC {
			t.Errorf("Expected PC unchanged after LD A,B, got PC=0x%04X, expected 0x%04X", cpu.PC, initialPC)
		}
	}
}

// TestLD_B_A tests the LD B,A instruction
func TestLD_B_A(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from A to B
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.A = value
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

		// Execute LD B,A instruction
		cycles := cpu.LD_B_A()

		// Should take 4 cycles
		if cycles != 4 {
			t.Errorf("Expected LD B,A to take 4 cycles, got %d", cycles)
		}

		// B register should now contain A's value
		if cpu.B != value {
			t.Errorf("Expected B=0x%02X after LD B,A, got B=0x%02X", value, cpu.B)
		}

		// A register should be unchanged (source remains intact)
		if cpu.A != initialA {
			t.Errorf("Expected A unchanged after LD B,A, got A=0x%02X, expected 0x%02X", cpu.A, initialA)
		}

		// All other registers and flags should be unchanged
		if cpu.C != initialC {
			t.Errorf("Expected C unchanged after LD B,A, got C=0x%02X, expected 0x%02X", cpu.C, initialC)
		}
		if cpu.D != initialD {
			t.Errorf("Expected D unchanged after LD B,A, got D=0x%02X, expected 0x%02X", cpu.D, initialD)
		}
		if cpu.E != initialE {
			t.Errorf("Expected E unchanged after LD B,A, got E=0x%02X, expected 0x%02X", cpu.E, initialE)
		}
		if cpu.H != initialH {
			t.Errorf("Expected H unchanged after LD B,A, got H=0x%02X, expected 0x%02X", cpu.H, initialH)
		}
		if cpu.L != initialL {
			t.Errorf("Expected L unchanged after LD B,A, got L=0x%02X, expected 0x%02X", cpu.L, initialL)
		}
		if cpu.F != initialF {
			t.Errorf("Expected F unchanged after LD B,A, got F=0x%02X, expected 0x%02X", cpu.F, initialF)
		}
		if cpu.SP != initialSP {
			t.Errorf("Expected SP unchanged after LD B,A, got SP=0x%04X, expected 0x%04X", cpu.SP, initialSP)
		}
		if cpu.PC != initialPC {
			t.Errorf("Expected PC unchanged after LD B,A, got PC=0x%04X, expected 0x%04X", cpu.PC, initialPC)
		}
	}
}

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

// TestLD_A_C tests the LD A,C instruction
func TestLD_A_C(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from C to A
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.C = value
		cpu.A = 0x99 // Different value in A
		cpu.F = 0x50 // Set some flags

		// Store initial state (other registers and flags should be unchanged)
		initialC := cpu.C
		initialB := cpu.B
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD A,C instruction
		cycles := cpu.LD_A_C()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD A,C should take 4 cycles")

		// A register should now contain C's value
		assert.Equal(t, value, cpu.A, "A register should contain C's value")

		// C register should be unchanged (source remains intact)
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")

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

// TestLD_C_n tests the LD C,n instruction
func TestLD_C_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into C
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD C,n instruction
		cycles := cpu.LD_C_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD C,n should take 8 cycles")

		// C register should contain the loaded value
		assert.Equal(t, value, cpu.C, "C register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
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

// TestINC_C tests the INC C instruction with various flag conditions
func TestINC_C(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.C = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_C()

	assert.Equal(t, uint8(4), cycles, "INC C should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.C, "C should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.C = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_C()

	assert.Equal(t, uint8(0x00), cpu.C, "C should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.C = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_C()

	assert.Equal(t, uint8(0x10), cpu.C, "C should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: Carry flag preservation (INC C does not affect carry)
	cpu.C = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_C()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC C")

	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.B = 0x88
	cpu.C = 0x42
	initialA := cpu.A
	initialB := cpu.B

	cpu.INC_C()

	assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
	assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
}

// TestDEC_C tests the DEC C instruction with various flag conditions
func TestDEC_C(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.C = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_C()

	assert.Equal(t, uint8(4), cycles, "DEC C should take 4 cycles")
	assert.Equal(t, uint8(0x41), cpu.C, "C should be decremented to 0x41")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.C = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_C()

	assert.Equal(t, uint8(0x00), cpu.C, "C should become 0x00 after 0x01 decrement")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after 0x01->0x00 decrement")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.C = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_C()

	assert.Equal(t, uint8(0xFF), cpu.C, "C should wrap to 0xFF after 0x00 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x00->0xFF decrement")

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.C = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_C()

	assert.Equal(t, uint8(0x0F), cpu.C, "C should become 0x0F after 0x10 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x10->0x0F decrement")

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.C = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_C()

	assert.Equal(t, uint8(0x0E), cpu.C, "C should become 0x0E after 0x0F decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear after 0x0F->0x0E decrement")

	// Test case 6: Carry flag preservation (DEC C doesn't affect carry)
	cpu.C = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_C()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after DEC C")
}

// TestLD_D_n tests the LD D,n instruction
func TestLD_D_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into D
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD D,n instruction
		cycles := cpu.LD_D_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD D,n should take 8 cycles")

		// D register should contain the loaded value
		assert.Equal(t, value, cpu.D, "D register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
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

// TestINC_D tests the INC D instruction with various flag conditions
func TestINC_D(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.D = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_D()

	assert.Equal(t, uint8(4), cycles, "INC D should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.D, "D should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.D = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_D()

	assert.Equal(t, uint8(0x00), cpu.D, "D should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.D = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_D()

	assert.Equal(t, uint8(0x10), cpu.D, "D should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: Carry flag preservation (INC D does not affect carry)
	cpu.D = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_D()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC D")

	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.B = 0x88
	cpu.C = 0x77
	cpu.D = 0x42
	initialA := cpu.A
	initialB := cpu.B
	initialC := cpu.C

	cpu.INC_D()

	assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
	assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
	assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
}

// TestDEC_D tests the DEC D instruction with various flag conditions
func TestDEC_D(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.D = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_D()

	assert.Equal(t, uint8(4), cycles, "DEC D should take 4 cycles")
	assert.Equal(t, uint8(0x41), cpu.D, "D should be decremented to 0x41")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.D = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_D()

	assert.Equal(t, uint8(0x00), cpu.D, "D should become 0x00 after 0x01 decrement")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after 0x01->0x00 decrement")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.D = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_D()

	assert.Equal(t, uint8(0xFF), cpu.D, "D should wrap to 0xFF after 0x00 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x00->0xFF decrement")

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.D = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_D()

	assert.Equal(t, uint8(0x0F), cpu.D, "D should become 0x0F after 0x10 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x10->0x0F decrement")

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.D = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_D()

	assert.Equal(t, uint8(0x0E), cpu.D, "D should become 0x0E after 0x0F decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear after 0x0F->0x0E decrement")

	// Test case 6: Carry flag preservation (DEC D doesn't affect carry)
	cpu.D = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_D()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after DEC D")
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

// TestLD_E_n tests the LD E,n instruction
func TestLD_E_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into E
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD E,n instruction
		cycles := cpu.LD_E_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD E,n should take 8 cycles")

		// E register should contain the loaded value
		assert.Equal(t, value, cpu.E, "E register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
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

// TestINC_E tests the INC E instruction with various flag conditions
func TestINC_E(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.E = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_E()

	assert.Equal(t, uint8(4), cycles, "INC E should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.E, "E should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.E = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_E()

	assert.Equal(t, uint8(0x00), cpu.E, "E should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.E = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_E()

	assert.Equal(t, uint8(0x10), cpu.E, "E should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: Carry flag preservation (INC E does not affect carry)
	cpu.E = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_E()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC E")

	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.B = 0x88
	cpu.C = 0x77
	cpu.D = 0x66
	cpu.E = 0x42
	initialA := cpu.A
	initialB := cpu.B
	initialC := cpu.C
	initialD := cpu.D

	cpu.INC_E()

	assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
	assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
	assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
	assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
}

// TestDEC_E tests the DEC E instruction with various flag conditions
func TestDEC_E(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.E = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_E()

	assert.Equal(t, uint8(4), cycles, "DEC E should take 4 cycles")
	assert.Equal(t, uint8(0x41), cpu.E, "E should be decremented to 0x41")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.E = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_E()

	assert.Equal(t, uint8(0x00), cpu.E, "E should become 0x00 after 0x01 decrement")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after 0x01->0x00 decrement")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)

	cpu.E = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_E()

	assert.Equal(t, uint8(0xFF), cpu.E, "E should wrap to 0xFF after 0x00 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x00->0xFF decrement")

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.E = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_E()

	assert.Equal(t, uint8(0x0F), cpu.E, "E should become 0x0F after 0x10 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x10->0x0F decrement")

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.E = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_E()

	assert.Equal(t, uint8(0x0E), cpu.E, "E should become 0x0E after 0x0F decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear after 0x0F->0x0E decrement")

	// Test case 6: Carry flag preservation (DEC E doesn't affect carry)
	cpu.E = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_E()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after DEC E")
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

// TestLD_E_B tests the LD E,B instruction
func TestLD_E_B(t *testing.T) {
	cpu := NewCPU()

	// Test copying different values from B to E
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Set up initial state
		cpu.B = value
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

		// Execute LD E,B instruction
		cycles := cpu.LD_E_B()

		// Should take 4 cycles
		assert.Equal(t, uint8(4), cycles, "LD E,B should take 4 cycles")

		// E register should now contain B's value
		assert.Equal(t, value, cpu.E, "E register should contain B's value")

		// B register should be unchanged (source remains intact)
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")

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

// TestLD_H_n tests the LD H,n instruction
func TestLD_H_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into H
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD H,n instruction
		cycles := cpu.LD_H_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD H,n should take 8 cycles")

		// H register should contain the loaded value
		assert.Equal(t, value, cpu.H, "H register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestINC_H tests the INC H instruction with various flag conditions
func TestINC_H(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.H = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_H()

	assert.Equal(t, uint8(4), cycles, "INC H should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.H, "H should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.H = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_H()

	assert.Equal(t, uint8(0x00), cpu.H, "H should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.H = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_H()

	assert.Equal(t, uint8(0x10), cpu.H, "H should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: Carry flag preservation (INC H does not affect carry)
	cpu.H = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_H()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC H")

	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.B = 0x88
	cpu.C = 0x77
	cpu.D = 0x66
	cpu.E = 0x55
	cpu.H = 0x42
	initialA := cpu.A
	initialB := cpu.B
	initialC := cpu.C
	initialD := cpu.D
	initialE := cpu.E

	cpu.INC_H()

	assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
	assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
	assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
	assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
	assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
}

// TestDEC_H tests the DEC H instruction with various flag conditions
func TestDEC_H(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.H = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_H()

	assert.Equal(t, uint8(4), cycles, "DEC H should take 4 cycles")
	assert.Equal(t, uint8(0x41), cpu.H, "H should be decremented to 0x41")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.H = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_H()

	assert.Equal(t, uint8(0x00), cpu.H, "H should become 0x00 after 0x01 decrement")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after 0x01->0x00 decrement")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.H = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_H()

	assert.Equal(t, uint8(0xFF), cpu.H, "H should wrap to 0xFF after 0x00 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x00->0xFF decrement")

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.H = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_H()

	assert.Equal(t, uint8(0x0F), cpu.H, "H should become 0x0F after 0x10 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x10->0x0F decrement")

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.H = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_H()

	assert.Equal(t, uint8(0x0E), cpu.H, "H should become 0x0E after 0x0F decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear after 0x0F->0x0E decrement")

	// Test case 6: Carry flag preservation (DEC H doesn't affect carry)
	cpu.H = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_H()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after DEC H")
}

// TestLD_A_H tests the LD A,(HL) instruction
func TestLD_A_HL(t *testing.T) {
	// Test loading A from memory at address stored in HL
	t.Run("Basic Memory Load", func(t *testing.T) {
		cpu := NewCPU()
		mmu := memory.NewMMU()

		// Set up test scenario
		testAddress := uint16(0x8000) // VRAM address
		testValue := uint8(0x42)      // Value to load

		// Set HL to point to test address
		cpu.SetHL(testAddress)

		// Write test value to memory
		mmu.WriteByte(testAddress, testValue)

		// Set A to different value to verify the load
		cpu.A = 0x00

		// Execute LD A,(HL) instruction
		cycles := cpu.LD_A_HL(mmu)

		// Verify the result
		assert.Equal(t, testValue, cpu.A, "A should contain the value loaded from memory")
		assert.Equal(t, uint8(8), cycles, "LD A,(HL) should take 8 cycles")
	})

	t.Run("Load from Different Memory Regions", func(t *testing.T) {
		cpu := NewCPU()
		mmu := memory.NewMMU()

		// Test loading from different memory regions
		testCases := []struct {
			name     string
			address  uint16
			value    uint8
			expected uint8
		}{
			{"VRAM", 0x8000, 0x12, 0x12},
			{"WRAM", 0xC000, 0x34, 0x34},
			{"HRAM", 0xFF80, 0x56, 0x56},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Set up test
				cpu.SetHL(tc.address)
				mmu.WriteByte(tc.address, tc.value)
				cpu.A = 0x00 // Clear A register

				// Execute instruction
				cycles := cpu.LD_A_HL(mmu)

				// Verify result
				assert.Equal(t, tc.expected, cpu.A, "A should contain the loaded value")
				assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
			})
		}
	})

	t.Run("Flags Unaffected", func(t *testing.T) {
		cpu := NewCPU()
		mmu := memory.NewMMU()

		// Set up test
		testAddress := uint16(0x8000)
		testValue := uint8(0x00) // Load zero to test flag behavior

		cpu.SetHL(testAddress)
		mmu.WriteByte(testAddress, testValue)

		// Set flags to known state
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)

		// Store the F register value AFTER setting flags
		originalF := cpu.F

		// Execute instruction
		cpu.LD_A_HL(mmu)

		// Verify flags are unchanged (LD instructions do not affect flags)
		assert.Equal(t, originalF, cpu.F, "F register should be unchanged")
		assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
	})

	t.Run("Register Preservation", func(t *testing.T) {
		cpu := NewCPU()
		mmu := memory.NewMMU()

		// Set up test
		testAddress := uint16(0x8000)
		testValue := uint8(0x99)

		cpu.SetHL(testAddress)
		mmu.WriteByte(testAddress, testValue)

		// Store initial register values
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute instruction
		cpu.LD_A_HL(mmu)

		// Verify only A changed
		assert.Equal(t, testValue, cpu.A, "A should contain the loaded value")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	})
}

// TestLD_L_n tests the LD_L_n instruction (0x2E)
func TestLD_L_n(t *testing.T) {
	t.Run("Load immediate value into L register", func(t *testing.T) {
		cpu := NewCPU()

		// Test various immediate values
		testCases := []struct {
			name           string
			value          uint8
			expectedL      uint8
			expectedCycles uint8
		}{
			{"Load 0x00", 0x00, 0x00, 8},
			{"Load 0x42", 0x42, 0x42, 8},
			{"Load 0xFF", 0xFF, 0xFF, 8},
			{"Load 0x80", 0x80, 0x80, 8},
			{"Load 0x7F", 0x7F, 0x7F, 8},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Set L to different value first
				cpu.L = ^tc.value // Bitwise NOT to ensure it changes

				// Execute instruction
				cycles := cpu.LD_L_n(tc.value)

				// Verify results
				assert.Equal(t, tc.expectedL, cpu.L, "L register should contain the immediate value")
				assert.Equal(t, tc.expectedCycles, cycles, "Should take 8 cycles")
			})
		}
	})

	t.Run("Flags unaffected", func(t *testing.T) {
		cpu := NewCPU()

		// Set all flags to test they remain unchanged
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)
		originalF := cpu.F

		// Execute instruction
		cpu.LD_L_n(0x55)

		// Verify flags are unchanged (LD instructions do not affect flags)
		assert.Equal(t, originalF, cpu.F, "F register should be unchanged")
		assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
	})
	t.Run("Register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Store initial values
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute instruction
		cpu.LD_L_n(0x99)

		// Verify only L changed
		assert.Equal(t, uint8(0x99), cpu.L, "L should contain the loaded value")
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	})
}

// TestINC_L tests the INC_L instruction (0x2C)
func TestINC_L(t *testing.T) {
	t.Run("Basic increment operations", func(t *testing.T) {
		cpu := NewCPU()

		testCases := []struct {
			name           string
			initialL       uint8
			expectedL      uint8
			expectedZ      bool
			expectedN      bool
			expectedH      bool
			expectedCycles uint8
		}{
			{"Increment 0x00", 0x00, 0x01, false, false, false, 4},
			{"Increment 0x0F (half-carry)", 0x0F, 0x10, false, false, true, 4},
			{"Increment 0x7F", 0x7F, 0x80, false, false, true, 4},
			{"Increment 0xFF (wrap to zero)", 0xFF, 0x00, true, false, true, 4},
			{"Increment 0x42", 0x42, 0x43, false, false, false, 4},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Set up initial state
				cpu.L = tc.initialL
				cpu.F = 0x00 // Clear all flags

				// Execute instruction
				cycles := cpu.INC_L()

				// Verify results
				assert.Equal(t, tc.expectedL, cpu.L, "L register should be incremented correctly")
				assert.Equal(t, tc.expectedCycles, cycles, "Should take 4 cycles")
				assert.Equal(t, tc.expectedZ, cpu.GetFlag(FlagZ), "Z flag should be set correctly")
				assert.Equal(t, tc.expectedN, cpu.GetFlag(FlagN), "N flag should be clear (addition)")
				assert.Equal(t, tc.expectedH, cpu.GetFlag(FlagH), "H flag should be set correctly")
			})
		}
	})

	t.Run("Carry flag preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Test with C flag set
		cpu.L = 0x42
		cpu.SetFlag(FlagC, true)

		cycles := cpu.INC_L()

		// Verify C flag is preserved
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be preserved")
		assert.Equal(t, uint8(0x43), cpu.L, "L should be incremented")
		assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")

		// Test with C flag clear
		cpu.L = 0x42
		cpu.SetFlag(FlagC, false)

		cpu.INC_L()

		// Verify C flag remains clear
		assert.False(t, cpu.GetFlag(FlagC), "C flag should remain clear")
	})

	t.Run("Register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Store initial values
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute instruction
		cpu.INC_L()

		// Verify other registers unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	})
}

// TestDEC_L tests the DEC_L instruction (0x2D)
func TestDEC_L(t *testing.T) {
	t.Run("Basic decrement operations", func(t *testing.T) {
		cpu := NewCPU()

		testCases := []struct {
			name           string
			initialL       uint8
			expectedL      uint8
			expectedZ      bool
			expectedN      bool
			expectedH      bool
			expectedCycles uint8
		}{
			{"Decrement 0x01", 0x01, 0x00, true, true, false, 4},
			{"Decrement 0x00 (wrap to 0xFF)", 0x00, 0xFF, false, true, true, 4},
			{"Decrement 0x10 (half-carry)", 0x10, 0x0F, false, true, true, 4},
			{"Decrement 0x80", 0x80, 0x7F, false, true, true, 4},
			{"Decrement 0x42", 0x42, 0x41, false, true, false, 4},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Set up initial state
				cpu.L = tc.initialL
				cpu.F = 0x00 // Clear all flags

				// Execute instruction
				cycles := cpu.DEC_L()

				// Verify results
				assert.Equal(t, tc.expectedL, cpu.L, "L register should be decremented correctly")
				assert.Equal(t, tc.expectedCycles, cycles, "Should take 4 cycles")
				assert.Equal(t, tc.expectedZ, cpu.GetFlag(FlagZ), "Z flag should be set correctly")
				assert.Equal(t, tc.expectedN, cpu.GetFlag(FlagN), "N flag should be set (subtraction)")
				assert.Equal(t, tc.expectedH, cpu.GetFlag(FlagH), "H flag should be set correctly")
			})
		}
	})

	t.Run("Carry flag preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Test with C flag set
		cpu.L = 0x42
		cpu.SetFlag(FlagC, true)

		cycles := cpu.DEC_L()

		// Verify C flag is preserved
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be preserved")
		assert.Equal(t, uint8(0x41), cpu.L, "L should be decremented")
		assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")

		// Test with C flag clear
		cpu.L = 0x42
		cpu.SetFlag(FlagC, false)

		cpu.DEC_L()

		// Verify C flag remains clear
		assert.False(t, cpu.GetFlag(FlagC), "C flag should remain clear")
	})

	t.Run("Register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Store initial values
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute instruction
		cpu.DEC_L()

		// Verify other registers unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	})
}

// TestLD_A_L tests the LD_A_L instruction (0x7D)
func TestLD_A_L(t *testing.T) {
	t.Run("Copy L to A register", func(t *testing.T) {
		cpu := NewCPU()

		testCases := []struct {
			name           string
			initialL       uint8
			expectedA      uint8
			expectedCycles uint8
		}{
			{"Copy 0x00", 0x00, 0x00, 4},
			{"Copy 0x42", 0x42, 0x42, 4},
			{"Copy 0xFF", 0xFF, 0xFF, 4},
			{"Copy 0x80", 0x80, 0x80, 4},
			{"Copy 0x7F", 0x7F, 0x7F, 4},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Set up initial state
				cpu.L = tc.initialL
				cpu.A = ^tc.initialL // Set A to different value

				// Execute instruction
				cycles := cpu.LD_A_L()

				// Verify results
				assert.Equal(t, tc.expectedA, cpu.A, "A register should contain L's value")
				assert.Equal(t, tc.expectedCycles, cycles, "Should take 4 cycles")
				assert.Equal(t, tc.initialL, cpu.L, "L register should be unchanged (source)")
			})
		}
	})

	t.Run("Flags unaffected", func(t *testing.T) {
		cpu := NewCPU()

		// Set all flags and test they remain unchanged
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)
		originalF := cpu.F

		cpu.L = 0x55
		cpu.LD_A_L()

		assert.Equal(t, originalF, cpu.F, "F register should be unchanged")
		assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagN), "N flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagH), "H flag should be unchanged")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should be unchanged")
	})

	t.Run("Register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Store initial values
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialSP := cpu.SP
		initialPC := cpu.PC

		cpu.L = 0x99
		cpu.LD_A_L()

		// Verify other registers unchanged
		assert.Equal(t, uint8(0x99), cpu.A, "A should contain L's value")
		assert.Equal(t, uint8(0x99), cpu.L, "L should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	})
}

// TestLD_B_L tests the LD_B_L instruction (0x45)
func TestLD_B_L(t *testing.T) {
	t.Run("Copy L to B register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.L = value
			cpu.B = ^value // Set B to different value

			cycles := cpu.LD_B_L()

			assert.Equal(t, value, cpu.B, "B register should contain L's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.L, "L register should be unchanged (source)")
		}
	})

	t.Run("Flags and other registers preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Set flags and store initial values
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagC, true)
		initialF := cpu.F
		initialA := cpu.A
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H

		cpu.L = 0x77
		cpu.LD_B_L()

		// Verify results
		assert.Equal(t, uint8(0x77), cpu.B, "B should contain L's value")
		assert.Equal(t, initialF, cpu.F, "Flags should be unchanged")
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
	})
}

// TestLD_C_L tests the LD_C_L instruction (0x4D)
func TestLD_C_L(t *testing.T) {
	t.Run("Copy L to C register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.L = value
			cpu.C = ^value // Set C to different value

			cycles := cpu.LD_C_L()

			assert.Equal(t, value, cpu.C, "C register should contain L's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.L, "L register should be unchanged (source)")
		}
	})

	t.Run("Comprehensive register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Store all initial values
		initialA := cpu.A
		initialB := cpu.B
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H
		initialF := cpu.F

		cpu.L = 0x88
		cpu.LD_C_L()

		// Verify only C changed
		assert.Equal(t, uint8(0x88), cpu.C, "C should contain L's value")
		assert.Equal(t, uint8(0x88), cpu.L, "L should be unchanged")
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
	})
}

// TestLD_L_A tests the LD_L_A instruction (0x6F)
func TestLD_L_A(t *testing.T) {
	t.Run("Copy A to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.A = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_A()

			assert.Equal(t, value, cpu.L, "L register should contain A's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.A, "A register should be unchanged (source)")
		}
	})

	t.Run("Flag and register preservation", func(t *testing.T) {
		cpu := NewCPU()

		// Set flags and store initial values
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		initialF := cpu.F
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialH := cpu.H

		cpu.A = 0x33
		cpu.LD_L_A()

		// Verify results
		assert.Equal(t, uint8(0x33), cpu.L, "L should contain A's value")
		assert.Equal(t, uint8(0x33), cpu.A, "A should be unchanged")
		assert.Equal(t, initialF, cpu.F, "Flags should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
	})
}

// TestLD_L_B tests the LD_L_B instruction (0x68)
func TestLD_L_B(t *testing.T) {
	t.Run("Copy B to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.B = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_B()

			assert.Equal(t, value, cpu.L, "L register should contain B's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.B, "B register should be unchanged (source)")
		}
	})
}

// TestLD_L_C tests the LD_L_C instruction (0x69)
func TestLD_L_C(t *testing.T) {
	t.Run("Copy C to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.C = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_C()

			assert.Equal(t, value, cpu.L, "L register should contain C's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.C, "C register should be unchanged (source)")
		}
	})
}

// TestLD_L_D tests the LD_L_D instruction (0x6A)
func TestLD_L_D(t *testing.T) {
	t.Run("Copy D to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.D = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_D()

			assert.Equal(t, value, cpu.L, "L register should contain D's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.D, "D register should be unchanged (source)")
		}
	})
}

// TestLD_L_E tests the LD_L_E instruction (0x6B)
func TestLD_L_E(t *testing.T) {
	t.Run("Copy E to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.E = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_E()

			assert.Equal(t, value, cpu.L, "L register should contain E's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.E, "E register should be unchanged (source)")
		}
	})
}

// TestLD_L_H tests the LD_L_H instruction (0x6C)
func TestLD_L_H(t *testing.T) {
	t.Run("Copy H to L register", func(t *testing.T) {
		cpu := NewCPU()

		testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F, 0x55, 0xAA}

		for _, value := range testValues {
			cpu.H = value
			cpu.L = ^value // Set L to different value

			cycles := cpu.LD_L_H()

			assert.Equal(t, value, cpu.L, "L register should contain H's value")
			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, value, cpu.H, "H register should be unchanged (source)")
		}
	})

	t.Run("HL register pair interaction", func(t *testing.T) {
		cpu := NewCPU()

		// Test that LD_L_H works correctly for HL register pair
		cpu.H = 0x12
		cpu.L = 0x34

		// Initial HL should be 0x1234
		assert.Equal(t, uint16(0x1234), cpu.GetHL(), "Initial HL should be 0x1234")

		// Execute LD_L_H (copy H to L)
		cpu.LD_L_H()

		// Now L should be 0x12, H should still be 0x12
		assert.Equal(t, uint8(0x12), cpu.L, "L should contain H's value")
		assert.Equal(t, uint8(0x12), cpu.H, "H should be unchanged")

		// HL should now be 0x1212
		assert.Equal(t, uint16(0x1212), cpu.GetHL(), "HL should be 0x1212 after LD_L_H")
	})
}

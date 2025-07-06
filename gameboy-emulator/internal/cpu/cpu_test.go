package cpu

import (
	"testing"
	
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
	
	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.C = 0x88
	cpu.B = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	initialA := cpu.A
	initialC := cpu.C
	
	cpu.DEC_B()
	
	// Other registers should be unchanged
	if cpu.A != initialA {
		t.Errorf("Expected A unchanged after DEC B, got A=0x%02X, expected 0x%02X", cpu.A, initialA)
	}
	if cpu.C != initialC {
		t.Errorf("Expected C unchanged after DEC B, got C=0x%02X, expected 0x%02X", cpu.C, initialC)
	}
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
		cpu.A = 0x99  // Different value in A
		cpu.F = 0x50  // Set some flags
		
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
		cpu.B = 0x99  // Different value in B
		cpu.F = 0x50  // Set some flags
		
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
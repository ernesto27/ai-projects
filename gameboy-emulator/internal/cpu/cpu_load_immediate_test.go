package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

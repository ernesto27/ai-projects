package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// TestINC_L tests the INC L instruction with various flag conditions
func TestINC_L(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal increment (no flags set)
	cpu.L = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.INC_L()

	assert.Equal(t, uint8(4), cycles, "INC L should take 4 cycles")
	assert.Equal(t, uint8(0x43), cpu.L, "L should be incremented to 0x43")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0xFF -> 0x00)
	cpu.L = 0xFF
	cpu.F = 0x00 // Clear all flags
	cpu.INC_L()

	assert.Equal(t, uint8(0x00), cpu.L, "L should wrap to 0x00 after 0xFF increment")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after overflow")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after overflow")

	// Test case 3: Half-carry flag set (0x0F -> 0x10)
	cpu.L = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.INC_L()

	assert.Equal(t, uint8(0x10), cpu.L, "L should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be clear")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set on 0x0F->0x10")

	// Test case 4: Carry flag preservation (INC L does not affect carry)
	cpu.L = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.INC_L()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after INC L")

	// Test case 5: Other registers unchanged
	cpu.A = 0x99
	cpu.B = 0x88
	cpu.C = 0x77
	cpu.D = 0x66
	cpu.E = 0x55
	cpu.H = 0x80 // HL = 0x8055
	cpu.L = 0x55
	initialA := cpu.A
	initialB := cpu.B
	initialC := cpu.C
	initialD := cpu.D
	initialE := cpu.E

	cpu.INC_L()

	assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
	assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
	assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
	assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
	assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
}

// TestDEC_L tests the DEC L instruction with various flag conditions
func TestDEC_L(t *testing.T) {
	cpu := NewCPU()

	// Test case 1: Normal decrement (subtract flag set)
	cpu.L = 0x42
	cpu.F = 0x00 // Clear all flags
	cycles := cpu.DEC_L()

	assert.Equal(t, uint8(4), cycles, "DEC L should take 4 cycles")
	assert.Equal(t, uint8(0x41), cpu.L, "L should be decremented to 0x41")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 2: Zero flag set (0x01 -> 0x00)
	cpu.L = 0x01
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_L()

	assert.Equal(t, uint8(0x00), cpu.L, "L should become 0x00 after 0x01 decrement")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set after 0x01->0x00 decrement")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear")

	// Test case 3: Half-carry flag set (0x00 -> 0xFF)
	cpu.L = 0x00
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_L()

	assert.Equal(t, uint8(0xFF), cpu.L, "L should wrap to 0xFF after 0x00 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x00->0xFF decrement")

	// Test case 4: Half-carry flag set (0x10 -> 0x0F)
	cpu.L = 0x10
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_L()

	assert.Equal(t, uint8(0x0F), cpu.L, "L should become 0x0F after 0x10 decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set after 0x10->0x0F decrement")

	// Test case 5: No half-carry (0x0F -> 0x0E)
	cpu.L = 0x0F
	cpu.F = 0x00 // Clear all flags
	cpu.DEC_L()

	assert.Equal(t, uint8(0x0E), cpu.L, "L should become 0x0E after 0x0F decrement")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be clear")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be clear after 0x0F->0x0E decrement")

	// Test case 6: Carry flag preservation (DEC L doesn't affect carry)
	cpu.L = 0x42
	cpu.SetFlag(FlagC, true) // Set carry flag
	cpu.DEC_L()

	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved after DEC L")
}

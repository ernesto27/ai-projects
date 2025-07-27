package cpu

import (
	"testing"
)

func TestJP_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: JP 0x1234 (little-endian: 0x34, 0x12)
	// Use WRAM for instruction data instead of ROM
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x34) // Low byte
	mmu.WriteByte(0xC001, 0x12) // High byte

	cycles := cpu.JP_nn(mmu)

	// Verify jump executed correctly
	if cpu.PC != 0x1234 {
		t.Errorf("Expected PC=0x1234, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
}

func TestJR_n_Forward(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: JR +5 (0x05)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x05) // Positive offset

	cycles := cpu.JR_n(mmu)

	// PC should be 0xC001 (after reading offset) + 0x05 = 0xC006
	if cpu.PC != 0xC006 {
		t.Errorf("Expected PC=0xC006, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_n_Backward(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: JR -3 (0xFD in two's complement)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0xFD) // -3 in two's complement

	cycles := cpu.JR_n(mmu)

	// PC should be 0xC001 (after reading offset) + (-3) = 0xBFFE  
	if cpu.PC != 0xBFFE {
		t.Errorf("Expected PC=0xBFFE, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_NZ_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag clear (jump should be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x34) // Low byte
	mmu.WriteByte(0xC001, 0x12) // High byte

	cycles := cpu.JP_NZ_nn(mmu)

	// Verify jump executed
	if cpu.PC != 0x1234 {
		t.Errorf("Expected PC=0x1234, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
}

func TestJP_NZ_nn_NoJump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag set (jump should NOT be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x34) // Low byte
	mmu.WriteByte(0xC001, 0x12) // High byte

	cycles := cpu.JP_NZ_nn(mmu)

	// Verify jump NOT executed (PC should be 0xC002 after reading address)
	if cpu.PC != 0xC002 {
		t.Errorf("Expected PC=0xC002, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_Z_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag set (jump should be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x34) // Low byte
	mmu.WriteByte(0xC001, 0x12) // High byte

	cycles := cpu.JP_Z_nn(mmu)

	// Verify jump executed
	if cpu.PC != 0x1234 {
		t.Errorf("Expected PC=0x1234, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
}

func TestJP_Z_nn_NoJump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag clear (jump should NOT be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x34) // Low byte
	mmu.WriteByte(0xC001, 0x12) // High byte

	cycles := cpu.JP_Z_nn(mmu)

	// Verify jump NOT executed (PC should be 0xC002 after reading address)
	if cpu.PC != 0xC002 {
		t.Errorf("Expected PC=0xC002, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_NC_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Carry flag clear (jump should be taken)
	cpu.SetFlag(FlagC, false)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x78) // Low byte
	mmu.WriteByte(0xC001, 0x56) // High byte

	cycles := cpu.JP_NC_nn(mmu)

	// Verify jump executed
	if cpu.PC != 0x5678 {
		t.Errorf("Expected PC=0x5678, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
}

func TestJP_C_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Carry flag set (jump should be taken)
	cpu.SetFlag(FlagC, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0xBC) // Low byte
	mmu.WriteByte(0xC001, 0x9A) // High byte

	cycles := cpu.JP_C_nn(mmu)

	// Verify jump executed
	if cpu.PC != 0x9ABC {
		t.Errorf("Expected PC=0x9ABC, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
}

func TestJR_NZ_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag clear (jump should be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x10) // +16 offset

	cycles := cpu.JR_NZ_n(mmu)

	// PC should be 0xC001 (after reading offset) + 16 = 0xC011
	if cpu.PC != 0xC011 {
		t.Errorf("Expected PC=0xC011, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_NZ_n_NoJump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag set (jump should NOT be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x10) // +16 offset

	cycles := cpu.JR_NZ_n(mmu)

	// PC should be 0xC001 (after reading offset, no jump)
	if cpu.PC != 0xC001 {
		t.Errorf("Expected PC=0xC001, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 8 {
		t.Errorf("Expected 8 cycles, got %d", cycles)
	}
}

func TestJR_Z_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Zero flag set (jump should be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0xF0) // -16 offset (0xF0 = -16 in two's complement)

	cycles := cpu.JR_Z_n(mmu)

	// PC should be 0xC001 (after reading offset) + (-16) = 0xBFF1  
	if cpu.PC != 0xBFF1 {
		t.Errorf("Expected PC=0xBFF1, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_NC_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Carry flag clear (jump should be taken)
	cpu.SetFlag(FlagC, false)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x05) // +5 offset

	cycles := cpu.JR_NC_n(mmu)

	// PC should be 0xC001 (after reading offset) + 5 = 0xC006
	if cpu.PC != 0xC006 {
		t.Errorf("Expected PC=0xC006, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_C_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set up test: Carry flag set (jump should be taken)
	cpu.SetFlag(FlagC, true)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0xFE) // -2 offset (0xFE = -2 in two's complement)

	cycles := cpu.JR_C_n(mmu)

	// PC should be 0xC001 (after reading offset) + (-2) = 0xBFFF
	if cpu.PC != 0xBFFF {
		t.Errorf("Expected PC=0xBFFF, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_HL(t *testing.T) {
	cpu := NewCPU()

	// Set up test: HL = 0xABCD
	cpu.SetHL(0xABCD)
	cpu.PC = 0xC000

	cycles := cpu.JP_HL()

	// Verify jump to HL address
	if cpu.PC != 0xABCD {
		t.Errorf("Expected PC=0xABCD, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 4 {
		t.Errorf("Expected 4 cycles, got %d", cycles)
	}
}

// Test edge cases and boundary conditions
func TestJR_EdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test maximum positive offset (+127)
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x7F) // +127

	cycles := cpu.JR_n(mmu)

	// PC should be 0xC001 + 127 = 0xC080
	if cpu.PC != 0xC080 {
		t.Errorf("Expected PC=0xC080, got PC=0x%04X", cpu.PC)
	}

	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}

	// Test maximum negative offset (-128)
	cpu.PC = 0xC200  
	mmu.WriteByte(0xC200, 0x80) // -128

	cycles = cpu.JR_n(mmu)

	// PC should be 0xC201 + (-128) = 0xC181
	if cpu.PC != 0xC181 {
		t.Errorf("Expected PC=0xC181, got PC=0x%04X", cpu.PC)
	}

	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

// Test that flags are not affected by jump instructions
func TestJumpInstructions_FlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags to specific values
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)

	// Store original flag values
	originalZ := cpu.GetFlag(FlagZ)
	originalN := cpu.GetFlag(FlagN)
	originalH := cpu.GetFlag(FlagH)
	originalC := cpu.GetFlag(FlagC)

	// Execute JP_nn
	cpu.PC = 0xC000
	mmu.WriteByte(0xC000, 0x00)
	mmu.WriteByte(0xC001, 0x02)
	cpu.JP_nn(mmu)

	// Verify flags are unchanged
	if cpu.GetFlag(FlagZ) != originalZ || cpu.GetFlag(FlagN) != originalN ||
		cpu.GetFlag(FlagH) != originalH || cpu.GetFlag(FlagC) != originalC {
		t.Error("Jump instructions should not affect flags")
	}
}

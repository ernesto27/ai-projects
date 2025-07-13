package cpu

import (
	"testing"

	"gameboy-emulator/internal/memory"
)

func TestJP_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: JP 0x1234 (little-endian: 0x34, 0x12)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x34) // Low byte
	mmu.WriteByte(0x0101, 0x12) // High byte

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
	mmu := memory.NewMMU()

	// Set up test: JR +5 (0x05)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x05) // Positive offset

	cycles := cpu.JR_n(mmu)

	// PC should be 0x0101 (after reading offset) + 0x05 = 0x0106
	if cpu.PC != 0x0106 {
		t.Errorf("Expected PC=0x0106, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_n_Backward(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: JR -3 (0xFD in two's complement)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0xFD) // -3 in two's complement

	cycles := cpu.JR_n(mmu)

	// PC should be 0x0101 (after reading offset) + (-3) = 0x00FE
	if cpu.PC != 0x00FE {
		t.Errorf("Expected PC=0x00FE, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_NZ_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Zero flag clear (jump should be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x34) // Low byte
	mmu.WriteByte(0x0101, 0x12) // High byte

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
	mmu := memory.NewMMU()

	// Set up test: Zero flag set (jump should NOT be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x34) // Low byte
	mmu.WriteByte(0x0101, 0x12) // High byte

	cycles := cpu.JP_NZ_nn(mmu)

	// Verify jump NOT executed (PC should be 0x0102 after reading address)
	if cpu.PC != 0x0102 {
		t.Errorf("Expected PC=0x0102, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_Z_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Zero flag set (jump should be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x34) // Low byte
	mmu.WriteByte(0x0101, 0x12) // High byte

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
	mmu := memory.NewMMU()

	// Set up test: Zero flag clear (jump should NOT be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x34) // Low byte
	mmu.WriteByte(0x0101, 0x12) // High byte

	cycles := cpu.JP_Z_nn(mmu)

	// Verify jump NOT executed (PC should be 0x0102 after reading address)
	if cpu.PC != 0x0102 {
		t.Errorf("Expected PC=0x0102, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJP_NC_nn_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Carry flag clear (jump should be taken)
	cpu.SetFlag(FlagC, false)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x78) // Low byte
	mmu.WriteByte(0x0101, 0x56) // High byte

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
	mmu := memory.NewMMU()

	// Set up test: Carry flag set (jump should be taken)
	cpu.SetFlag(FlagC, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0xBC) // Low byte
	mmu.WriteByte(0x0101, 0x9A) // High byte

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
	mmu := memory.NewMMU()

	// Set up test: Zero flag clear (jump should be taken)
	cpu.SetFlag(FlagZ, false)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x10) // +16 offset

	cycles := cpu.JR_NZ_n(mmu)

	// PC should be 0x0101 (after reading offset) + 16 = 0x0111
	if cpu.PC != 0x0111 {
		t.Errorf("Expected PC=0x0111, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_NZ_n_NoJump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Zero flag set (jump should NOT be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x10) // +16 offset

	cycles := cpu.JR_NZ_n(mmu)

	// PC should be 0x0101 (after reading offset, no jump)
	if cpu.PC != 0x0101 {
		t.Errorf("Expected PC=0x0101, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for not taken jump
	if cycles != 8 {
		t.Errorf("Expected 8 cycles, got %d", cycles)
	}
}

func TestJR_Z_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Zero flag set (jump should be taken)
	cpu.SetFlag(FlagZ, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0xF0) // -16 offset (0xF0 = -16 in two's complement)

	cycles := cpu.JR_Z_n(mmu)

	// PC should be 0x0101 (after reading offset) + (-16) = 0x00F1
	if cpu.PC != 0x00F1 {
		t.Errorf("Expected PC=0x00F1, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_NC_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Carry flag clear (jump should be taken)
	cpu.SetFlag(FlagC, false)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x05) // +5 offset

	cycles := cpu.JR_NC_n(mmu)

	// PC should be 0x0101 (after reading offset) + 5 = 0x0106
	if cpu.PC != 0x0106 {
		t.Errorf("Expected PC=0x0106, got PC=0x%04X", cpu.PC)
	}

	// Verify cycle count for taken jump
	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

func TestJR_C_n_Jump(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up test: Carry flag set (jump should be taken)
	cpu.SetFlag(FlagC, true)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0xFE) // -2 offset (0xFE = -2 in two's complement)

	cycles := cpu.JR_C_n(mmu)

	// PC should be 0x0101 (after reading offset) + (-2) = 0x00FF
	if cpu.PC != 0x00FF {
		t.Errorf("Expected PC=0x00FF, got PC=0x%04X", cpu.PC)
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
	cpu.PC = 0x0100

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
	mmu := memory.NewMMU()

	// Test maximum positive offset (+127)
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x7F) // +127

	cycles := cpu.JR_n(mmu)

	// PC should be 0x0101 + 127 = 0x0180
	if cpu.PC != 0x0180 {
		t.Errorf("Expected PC=0x0180, got PC=0x%04X", cpu.PC)
	}

	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}

	// Test maximum negative offset (-128)
	cpu.PC = 0x0200
	mmu.WriteByte(0x0200, 0x80) // -128

	cycles = cpu.JR_n(mmu)

	// PC should be 0x0201 + (-128) = 0x0181
	if cpu.PC != 0x0181 {
		t.Errorf("Expected PC=0x0181, got PC=0x%04X", cpu.PC)
	}

	if cycles != 12 {
		t.Errorf("Expected 12 cycles, got %d", cycles)
	}
}

// Test that flags are not affected by jump instructions
func TestJumpInstructions_FlagsUnaffected(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

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
	cpu.PC = 0x0100
	mmu.WriteByte(0x0100, 0x00)
	mmu.WriteByte(0x0101, 0x02)
	cpu.JP_nn(mmu)

	// Verify flags are unchanged
	if cpu.GetFlag(FlagZ) != originalZ || cpu.GetFlag(FlagN) != originalN ||
		cpu.GetFlag(FlagH) != originalH || cpu.GetFlag(FlagC) != originalC {
		t.Error("Jump instructions should not affect flags")
	}
}

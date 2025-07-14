package cpu

import (
	"testing"

	"gameboy-emulator/internal/memory"

	"github.com/stretchr/testify/assert"
)

// ================================
// PUSH Operation Tests
// ================================

func TestPUSH_BC(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test basic PUSH BC operation
	cpu.B = 0x12
	cpu.C = 0x34
	cpu.SP = 0xFFFE

	cycles := cpu.PUSH_BC(mmu)

	// Verify stack pointer decremented by 2
	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should be decremented by 2")

	// Verify B pushed to SP+1 (high byte)
	assert.Equal(t, uint8(0x12), mmu.ReadByte(0xFFFD), "Register B should be pushed to high address")

	// Verify C pushed to SP (low byte)
	assert.Equal(t, uint8(0x34), mmu.ReadByte(0xFFFC), "Register C should be pushed to low address")

	// Verify cycle count
	assert.Equal(t, uint8(16), cycles, "PUSH BC should take 16 cycles")

	// Verify BC register values unchanged
	assert.Equal(t, uint8(0x12), cpu.B, "Register B should be unchanged after PUSH")
	assert.Equal(t, uint8(0x34), cpu.C, "Register C should be unchanged after PUSH")
}

func TestPUSH_DE(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.D = 0x56
	cpu.E = 0x78
	cpu.SP = 0xFFFE

	cycles := cpu.PUSH_DE(mmu)

	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should be decremented by 2")
	assert.Equal(t, uint8(0x56), mmu.ReadByte(0xFFFD), "Register D should be pushed to high address")
	assert.Equal(t, uint8(0x78), mmu.ReadByte(0xFFFC), "Register E should be pushed to low address")
	assert.Equal(t, uint8(16), cycles, "PUSH DE should take 16 cycles")
}

func TestPUSH_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.H = 0x9A
	cpu.L = 0xBC
	cpu.SP = 0xFFFE

	cycles := cpu.PUSH_HL(mmu)

	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should be decremented by 2")
	assert.Equal(t, uint8(0x9A), mmu.ReadByte(0xFFFD), "Register H should be pushed to high address")
	assert.Equal(t, uint8(0xBC), mmu.ReadByte(0xFFFC), "Register L should be pushed to low address")
	assert.Equal(t, uint8(16), cycles, "PUSH HL should take 16 cycles")
}

func TestPUSH_AF(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.A = 0xDE
	cpu.F = 0xF0
	cpu.SP = 0xFFFE

	cycles := cpu.PUSH_AF(mmu)

	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should be decremented by 2")
	assert.Equal(t, uint8(0xDE), mmu.ReadByte(0xFFFD), "Register A should be pushed to high address")
	assert.Equal(t, uint8(0xF0), mmu.ReadByte(0xFFFC), "Register F should be pushed to low address")
	assert.Equal(t, uint8(16), cycles, "PUSH AF should take 16 cycles")
}

// ================================
// POP Operation Tests
// ================================

func TestPOP_BC(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Setup stack with test data
	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0x34) // Low byte (C)
	mmu.WriteByte(0xFFFD, 0x12) // High byte (B)

	cycles := cpu.POP_BC(mmu)

	// Verify registers loaded correctly
	assert.Equal(t, uint8(0x12), cpu.B, "Register B should be loaded from stack")
	assert.Equal(t, uint8(0x34), cpu.C, "Register C should be loaded from stack")

	// Verify stack pointer incremented by 2
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be incremented by 2")

	// Verify cycle count
	assert.Equal(t, uint8(12), cycles, "POP BC should take 12 cycles")
}

func TestPOP_DE(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0x78) // Low byte (E)
	mmu.WriteByte(0xFFFD, 0x56) // High byte (D)

	cycles := cpu.POP_DE(mmu)

	assert.Equal(t, uint8(0x56), cpu.D, "Register D should be loaded from stack")
	assert.Equal(t, uint8(0x78), cpu.E, "Register E should be loaded from stack")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be incremented by 2")
	assert.Equal(t, uint8(12), cycles, "POP DE should take 12 cycles")
}

func TestPOP_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0xBC) // Low byte (L)
	mmu.WriteByte(0xFFFD, 0x9A) // High byte (H)

	cycles := cpu.POP_HL(mmu)

	assert.Equal(t, uint8(0x9A), cpu.H, "Register H should be loaded from stack")
	assert.Equal(t, uint8(0xBC), cpu.L, "Register L should be loaded from stack")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be incremented by 2")
	assert.Equal(t, uint8(12), cycles, "POP HL should take 12 cycles")
}

func TestPOP_AF(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0xF0) // Low byte (F) - flags
	mmu.WriteByte(0xFFFD, 0xDE) // High byte (A)

	cycles := cpu.POP_AF(mmu)

	assert.Equal(t, uint8(0xDE), cpu.A, "Register A should be loaded from stack")
	assert.Equal(t, uint8(0xF0), cpu.F, "Register F should be loaded from stack")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be incremented by 2")
	assert.Equal(t, uint8(12), cycles, "POP AF should take 12 cycles")
}

// ================================
// CALL Operation Tests
// ================================

func TestCALL_nn(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Setup: PC=0x8000, SP=0xFFFE, target address=0x4000
	cpu.PC = 0x8000
	cpu.SP = 0xFFFE

	// Write target address in memory (little-endian)
	mmu.WriteByte(0x8000, 0x00) // Low byte of 0x4000
	mmu.WriteByte(0x8001, 0x40) // High byte of 0x4000

	cycles := cpu.CALL_nn(mmu)

	// Verify PC jumped to target address
	assert.Equal(t, uint16(0x4000), cpu.PC, "PC should jump to target address")

	// Verify return address (0x8002) pushed to stack
	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should be decremented by 2")

	// Check return address on stack (little-endian)
	returnLow := mmu.ReadByte(0xFFFC)
	returnHigh := mmu.ReadByte(0xFFFD)
	returnAddr := uint16(returnHigh)<<8 | uint16(returnLow)

	assert.Equal(t, uint16(0x8002), returnAddr, "Return address should be pushed to stack")
	assert.Equal(t, uint8(24), cycles, "CALL nn should take 24 cycles")
}

func TestCALL_NZ_nn_WhenZClear(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.PC = 0x8000
	cpu.SP = 0xFFFE
	cpu.SetFlag(FlagZ, false) // Clear Z flag

	// Write target address
	mmu.WriteByte(0x8000, 0x00)
	mmu.WriteByte(0x8001, 0x40)

	cycles := cpu.CALL_NZ_nn(mmu)

	// Should execute call since Z is clear
	assert.Equal(t, uint16(0x4000), cpu.PC, "PC should jump to target address when Z is clear")
	assert.Equal(t, uint8(24), cycles, "CALL NZ should take 24 cycles when condition is true")
}

func TestCALL_NZ_nn_WhenZSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.PC = 0x8000
	cpu.SP = 0xFFFE
	cpu.SetFlag(FlagZ, true) // Set Z flag

	// Write target address
	mmu.WriteByte(0x8000, 0x00)
	mmu.WriteByte(0x8001, 0x40)

	cycles := cpu.CALL_NZ_nn(mmu)

	// Should NOT execute call since Z is set
	assert.Equal(t, uint16(0x8002), cpu.PC, "PC should skip call when Z is set")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be unchanged when call is skipped")
	assert.Equal(t, uint8(12), cycles, "CALL NZ should take 12 cycles when condition is false")
}

func TestCALL_Z_nn_WhenZSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.PC = 0x8000
	cpu.SP = 0xFFFE
	cpu.SetFlag(FlagZ, true) // Set Z flag

	mmu.WriteByte(0x8000, 0x00)
	mmu.WriteByte(0x8001, 0x40)

	cycles := cpu.CALL_Z_nn(mmu)

	// Should execute call since Z is set
	assert.Equal(t, uint16(0x4000), cpu.PC, "PC should jump to target address when Z is set")
	assert.Equal(t, uint8(24), cycles, "CALL Z should take 24 cycles when condition is true")
}

func TestCALL_NC_nn_WhenCClear(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.PC = 0x8000
	cpu.SP = 0xFFFE
	cpu.SetFlag(FlagC, false) // Clear C flag

	mmu.WriteByte(0x8000, 0x00)
	mmu.WriteByte(0x8001, 0x40)

	cycles := cpu.CALL_NC_nn(mmu)

	// Should execute call since C is clear
	assert.Equal(t, uint16(0x4000), cpu.PC, "PC should jump to target address when C is clear")
	assert.Equal(t, uint8(24), cycles, "CALL NC should take 24 cycles when condition is true")
}

func TestCALL_C_nn_WhenCSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.PC = 0x8000
	cpu.SP = 0xFFFE
	cpu.SetFlag(FlagC, true) // Set C flag

	mmu.WriteByte(0x8000, 0x00)
	mmu.WriteByte(0x8001, 0x40)

	cycles := cpu.CALL_C_nn(mmu)

	// Should execute call since C is set
	assert.Equal(t, uint16(0x4000), cpu.PC, "PC should jump to target address when C is set")
	assert.Equal(t, uint8(24), cycles, "CALL C should take 24 cycles when condition is true")
}

// ================================
// RET Operation Tests
// ================================

func TestRET(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Setup stack with return address 0x8003
	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0x03) // Low byte
	mmu.WriteByte(0xFFFD, 0x80) // High byte

	cycles := cpu.RET(mmu)

	// Verify PC set to return address
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set to return address")

	// Verify SP incremented by 2
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should be incremented by 2")

	assert.Equal(t, uint8(16), cycles, "RET should take 16 cycles")
}

func TestRET_NZ_WhenZClear(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	cpu.SetFlag(FlagZ, false) // Clear Z flag
	mmu.WriteByte(0xFFFC, 0x03)
	mmu.WriteByte(0xFFFD, 0x80)

	cycles := cpu.RET_NZ(mmu)

	// Should execute return since Z is clear
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set correctly")

	if cycles != 20 {
		t.Errorf("Expected 20 cycles, got %d", cycles)
	}
}

func TestRET_NZ_WhenZSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	originalPC := cpu.PC
	cpu.SP = 0xFFFC
	cpu.SetFlag(FlagZ, true) // Set Z flag

	cycles := cpu.RET_NZ(mmu)

	// Should NOT execute return since Z is set
	assert.Equal(t, uint16(originalPC), cpu.PC, "PC should be set correctly")

	if cpu.SP != 0xFFFC {
		t.Errorf("Expected SP unchanged, got SP=0x%04X", cpu.SP)
	}

	if cycles != 8 {
		t.Errorf("Expected 8 cycles, got %d", cycles)
	}
}

func TestRET_Z_WhenZSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	cpu.SetFlag(FlagZ, true) // Set Z flag
	mmu.WriteByte(0xFFFC, 0x03)
	mmu.WriteByte(0xFFFD, 0x80)

	cycles := cpu.RET_Z(mmu)

	// Should execute return since Z is set
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set correctly")

	if cycles != 20 {
		t.Errorf("Expected 20 cycles, got %d", cycles)
	}
}

func TestRET_NC_WhenCClear(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	cpu.SetFlag(FlagC, false) // Clear C flag
	mmu.WriteByte(0xFFFC, 0x03)
	mmu.WriteByte(0xFFFD, 0x80)

	cycles := cpu.RET_NC(mmu)

	// Should execute return since C is clear
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set correctly")

	if cycles != 20 {
		t.Errorf("Expected 20 cycles, got %d", cycles)
	}
}

func TestRET_C_WhenCSet(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFC
	cpu.SetFlag(FlagC, true) // Set C flag
	mmu.WriteByte(0xFFFC, 0x03)
	mmu.WriteByte(0xFFFD, 0x80)

	cycles := cpu.RET_C(mmu)

	// Should execute return since C is set
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set correctly")

	if cycles != 20 {
		t.Errorf("Expected 20 cycles, got %d", cycles)
	}
}

func TestRETI(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Setup stack with return address
	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0x03)
	mmu.WriteByte(0xFFFD, 0x80)

	cycles := cpu.RETI(mmu)

	// Should behave like regular RET
	assert.Equal(t, uint16(0x8003), cpu.PC, "PC should be set correctly")

	if cpu.SP != 0xFFFE {
		t.Errorf("Expected SP=0xFFFE, got SP=0x%04X", cpu.SP)
	}

	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}

	// TODO: Test interrupt enable when interrupt system is implemented
}

// ================================
// Stack Helper Method Tests (Phase 1)
// ================================

func TestPushByte(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set initial stack pointer
	cpu.SP = 0xFFFE

	// Push a byte
	cpu.pushByte(mmu, 0x42)

	// Verify stack pointer decremented
	assert.Equal(t, uint16(0xFFFD), cpu.SP, "Stack pointer should decrement after pushByte")

	// Verify byte written to stack
	assert.Equal(t, uint8(0x42), mmu.ReadByte(0xFFFD), "Byte should be written to stack location")

	// Push another byte
	cpu.pushByte(mmu, 0x84)

	// Verify stack pointer decremented again
	assert.Equal(t, uint16(0xFFFC), cpu.SP, "Stack pointer should decrement again")

	// Verify second byte written
	assert.Equal(t, uint8(0x84), mmu.ReadByte(0xFFFC), "Second byte should be written to new stack location")

	// Verify first byte still there
	assert.Equal(t, uint8(0x42), mmu.ReadByte(0xFFFD), "First byte should still be on stack")
}

func TestPopByte(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up stack with data
	cpu.SP = 0xFFFC
	mmu.WriteByte(0xFFFC, 0x84)
	mmu.WriteByte(0xFFFD, 0x42)

	// Pop first byte (should be 0x84)
	value1 := cpu.popByte(mmu)
	assert.Equal(t, uint8(0x84), value1, "Should pop the correct byte")
	assert.Equal(t, uint16(0xFFFD), cpu.SP, "Stack pointer should increment after popByte")

	// Pop second byte (should be 0x42)
	value2 := cpu.popByte(mmu)
	assert.Equal(t, uint8(0x42), value2, "Should pop the second byte")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "Stack pointer should increment again")
}

func TestPushPopByteRoundTrip(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SP = 0xFFFE
	originalSP := cpu.SP

	// Test round-trip with various values
	testValues := []uint8{0x00, 0x42, 0xFF, 0x80, 0x7F}

	for _, testValue := range testValues {
		// Push the value
		cpu.pushByte(mmu, testValue)

		// Pop the value back
		poppedValue := cpu.popByte(mmu)

		// Verify round-trip
		assert.Equal(t, testValue, poppedValue, "Push/pop round-trip should preserve value")
		assert.Equal(t, originalSP, cpu.SP, "Stack pointer should return to original position")
	}
}

func TestGetStackTop(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test with empty stack
	cpu.SP = 0xFFFE
	top := cpu.getStackTop(mmu)
	assert.Equal(t, uint8(0), top, "Empty stack should return 0")

	// Test with data on stack
	cpu.SP = 0xFFFD
	mmu.WriteByte(0xFFFD, 0x99)
	top = cpu.getStackTop(mmu)
	assert.Equal(t, uint8(0x99), top, "Should return top of stack value")
	assert.Equal(t, uint16(0xFFFD), cpu.SP, "getStackTop should not modify SP")
}

func TestGetStackDepth(t *testing.T) {
	cpu := NewCPU()

	// Test empty stack
	cpu.SP = 0xFFFE
	depth := cpu.getStackDepth()
	assert.Equal(t, uint16(0), depth, "Empty stack should have depth 0")

	// Test with some items
	cpu.SP = 0xFFFC
	depth = cpu.getStackDepth()
	assert.Equal(t, uint16(2), depth, "Stack with 2 bytes should have depth 2")

	// Test with many items
	cpu.SP = 0xFF80
	depth = cpu.getStackDepth()
	assert.Equal(t, uint16(126), depth, "Stack should calculate correct depth")
}

func TestIsStackEmpty(t *testing.T) {
	cpu := NewCPU()

	// Test empty stack
	cpu.SP = 0xFFFE
	assert.True(t, cpu.isStackEmpty(), "Stack should be empty at initial position")

	// Test with overflow condition
	cpu.SP = 0xFFFF
	assert.True(t, cpu.isStackEmpty(), "Stack should be empty with SP overflow")

	// Test non-empty stack
	cpu.SP = 0xFFFD
	assert.False(t, cpu.isStackEmpty(), "Stack should not be empty with items")

	cpu.SP = 0xFF80
	assert.False(t, cpu.isStackEmpty(), "Stack should not be empty with many items")
}

func TestStackHelperIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Start with empty stack
	cpu.SP = 0xFFFE
	assert.True(t, cpu.isStackEmpty(), "Should start with empty stack")
	assert.Equal(t, uint16(0), cpu.getStackDepth(), "Should start with zero depth")

	// Push some bytes using helper
	cpu.pushByte(mmu, 0x11)
	assert.False(t, cpu.isStackEmpty(), "Should not be empty after push")
	assert.Equal(t, uint16(1), cpu.getStackDepth(), "Should have depth 1")
	assert.Equal(t, uint8(0x11), cpu.getStackTop(mmu), "Top should be the pushed value")

	cpu.pushByte(mmu, 0x22)
	assert.Equal(t, uint16(2), cpu.getStackDepth(), "Should have depth 2")
	assert.Equal(t, uint8(0x22), cpu.getStackTop(mmu), "Top should be the last pushed value")

	// Pop bytes back
	value1 := cpu.popByte(mmu)
	assert.Equal(t, uint8(0x22), value1, "Should pop last pushed value first")
	assert.Equal(t, uint16(1), cpu.getStackDepth(), "Should have depth 1 after pop")

	value2 := cpu.popByte(mmu)
	assert.Equal(t, uint8(0x11), value2, "Should pop first pushed value last")
	assert.True(t, cpu.isStackEmpty(), "Should be empty after popping all")
	assert.Equal(t, uint16(0), cpu.getStackDepth(), "Should have zero depth")
}

package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for Memory Increment/Decrement Operations ===

func TestINC_HL_mem(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test 1: Normal increment in memory
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x42)

	cycles := cpu.INC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "INC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x43), mmu.ReadByte(0x8000), "Memory value should be incremented")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL register should remain unchanged")

	// Check flags
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should not be set")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be cleared")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should not be set")

	// Test 2: Zero result
	mmu.WriteByte(0x8000, 0xFF)
	cycles = cpu.INC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "INC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x00), mmu.ReadByte(0x8000), "Memory should wrap from 0xFF to 0x00")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be cleared")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set (0xFF -> 0x00)")

	// Test 3: Half-carry detection
	mmu.WriteByte(0x8000, 0x0F)
	cycles = cpu.INC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "INC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x10), mmu.ReadByte(0x8000), "Memory should increment from 0x0F to 0x10")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should not be set")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be cleared")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set (0x0F -> 0x10)")

	// Test 4: Carry flag preservation
	cpu.SetFlag(FlagC, true) // Set carry flag
	mmu.WriteByte(0x8000, 0x42)
	cpu.INC_HL_mem(mmu)
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

func TestDEC_HL_mem(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test 1: Normal decrement in memory
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x42)

	cycles := cpu.DEC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "DEC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x41), mmu.ReadByte(0x8000), "Memory value should be decremented")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL register should remain unchanged")

	// Check flags
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should not be set")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should not be set")

	// Test 2: Zero result
	mmu.WriteByte(0x8000, 0x01)
	cycles = cpu.DEC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "DEC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x00), mmu.ReadByte(0x8000), "Memory should decrement from 0x01 to 0x00")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should not be set")

	// Test 3: Underflow wrap-around
	mmu.WriteByte(0x8000, 0x00)
	cycles = cpu.DEC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "DEC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0x8000), "Memory should wrap from 0x00 to 0xFF")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should not be set")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set (0x00 -> 0xFF)")

	// Test 4: Half-carry detection
	mmu.WriteByte(0x8000, 0x10)
	cycles = cpu.DEC_HL_mem(mmu)

	assert.Equal(t, uint8(12), cycles, "DEC (HL) should take 12 cycles")
	assert.Equal(t, uint8(0x0F), mmu.ReadByte(0x8000), "Memory should decrement from 0x10 to 0x0F")
	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should not be set")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be set")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set (0x10 -> 0x0F)")

	// Test 5: Carry flag preservation
	cpu.SetFlag(FlagC, true) // Set carry flag
	mmu.WriteByte(0x8000, 0x42)
	cpu.DEC_HL_mem(mmu)
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

// === Tests for Memory Store Operations ===

func TestLD_HL_mem_n(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test 1: Store immediate value to memory
	cpu.SetHL(0x8000)
	cycles := cpu.LD_HL_mem_n(mmu, 0x99)

	assert.Equal(t, uint8(12), cycles, "LD (HL),n should take 12 cycles")
	assert.Equal(t, uint8(0x99), mmu.ReadByte(0x8000), "Memory should contain the immediate value")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL register should remain unchanged")

	// Test 2: Store different values
	cpu.SetHL(0x9000)
	cycles = cpu.LD_HL_mem_n(mmu, 0x00)

	assert.Equal(t, uint8(12), cycles, "LD (HL),n should take 12 cycles")
	assert.Equal(t, uint8(0x00), mmu.ReadByte(0x9000), "Memory should contain 0x00")

	cycles = cpu.LD_HL_mem_n(mmu, 0xFF)
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0x9000), "Memory should contain 0xFF")

	// Test 3: Flags should not be affected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_HL_mem_n(mmu, 0x42)

	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

// === Tests for Memory Load Operations ===

func TestLD_B_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test 1: Load from memory to register B
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x55)
	cpu.B = 0x00 // Clear B first

	cycles := cpu.LD_B_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD B,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0x55), cpu.B, "Register B should contain the memory value")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL register should remain unchanged")
	assert.Equal(t, uint8(0x55), mmu.ReadByte(0x8000), "Memory should remain unchanged")

	// Test 2: Load different values
	mmu.WriteByte(0x8000, 0xFF)
	cycles = cpu.LD_B_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD B,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0xFF), cpu.B, "Register B should contain 0xFF")

	// Test 3: Flags should not be affected
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.LD_B_HL(mmu)

	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

func TestLD_C_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0xAA)
	cpu.C = 0x00

	cycles := cpu.LD_C_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD C,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0xAA), cpu.C, "Register C should contain the memory value")
}

func TestLD_D_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0xBB)
	cpu.D = 0x00

	cycles := cpu.LD_D_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD D,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0xBB), cpu.D, "Register D should contain the memory value")
}

func TestLD_E_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0xCC)
	cpu.E = 0x00

	cycles := cpu.LD_E_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD E,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0xCC), cpu.E, "Register E should contain the memory value")
}

func TestLD_H_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test: This is an interesting case - we read from the address formed by H and L,
	// then store the result in H
	cpu.SetHL(0x8000) // H=0x80, L=0x00
	mmu.WriteByte(0x8000, 0xDD)

	cycles := cpu.LD_H_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD H,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0xDD), cpu.H, "Register H should contain the memory value")
	assert.Equal(t, uint8(0x00), cpu.L, "Register L should remain unchanged")
	// Note: After this operation, HL = 0xDD00 (new H value + old L value)
	assert.Equal(t, uint16(0xDD00), cpu.GetHL(), "HL should now be 0xDD00")
}

func TestLD_L_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test: Similar to LD_H_HL - we read from the address formed by H and L,
	// then store the result in L
	cpu.SetHL(0x8000) // H=0x80, L=0x00
	mmu.WriteByte(0x8000, 0xEE)

	cycles := cpu.LD_L_HL(mmu)

	assert.Equal(t, uint8(8), cycles, "LD L,(HL) should take 8 cycles")
	assert.Equal(t, uint8(0x80), cpu.H, "Register H should remain unchanged")
	assert.Equal(t, uint8(0xEE), cpu.L, "Register L should contain the memory value")
	// Note: After this operation, HL = 0x80EE (old H value + new L value)
	assert.Equal(t, uint16(0x80EE), cpu.GetHL(), "HL should now be 0x80EE")
}

// === Tests for Wrapper Functions ===

func TestWrapINC_HL_mem(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x10)

	cycles, err := wrapINC_HL_mem(cpu, mmu)

	assert.NoError(t, err, "wrapINC_HL_mem should not return error")
	assert.Equal(t, uint8(12), cycles, "wrapINC_HL_mem should return 12 cycles")
	assert.Equal(t, uint8(0x11), mmu.ReadByte(0x8000), "Memory should be incremented")

	// Test that parameters don't matter (wrapper ignores them)
	cycles, err = wrapINC_HL_mem(cpu, mmu, 0x42, 0x43)
	assert.NoError(t, err, "wrapINC_HL_mem should work with parameters")
	assert.Equal(t, uint8(12), cycles, "wrapINC_HL_mem should return 12 cycles with params")
}

func TestWrapLD_HL_mem_n(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)

	// Test normal operation
	cycles, err := wrapLD_HL_mem_n(cpu, mmu, 0x77)

	assert.NoError(t, err, "wrapLD_HL_mem_n should not return error")
	assert.Equal(t, uint8(12), cycles, "wrapLD_HL_mem_n should return 12 cycles")
	assert.Equal(t, uint8(0x77), mmu.ReadByte(0x8000), "Memory should contain immediate value")

	// Test error handling - no parameters
	_, err = wrapLD_HL_mem_n(cpu, mmu)
	assert.Error(t, err, "wrapLD_HL_mem_n should return error with no parameters")
	assert.Contains(t, err.Error(), "requires 1 parameter, got 0")
}

func TestWrapLD_B_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x33)
	cpu.B = 0x00

	cycles, err := wrapLD_B_HL(cpu, mmu)

	assert.NoError(t, err, "wrapLD_B_HL should not return error")
	assert.Equal(t, uint8(8), cycles, "wrapLD_B_HL should return 8 cycles")
	assert.Equal(t, uint8(0x33), cpu.B, "Register B should contain memory value")
}

// === Integration Tests ===

func TestMemoryOperationsIntegration(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test a sequence of memory operations
	cpu.SetHL(0x8000)

	// Store a value, then increment it, then load it to a register
	cpu.LD_HL_mem_n(mmu, 0x10) // Store 0x10 at (HL)
	assert.Equal(t, uint8(0x10), mmu.ReadByte(0x8000))

	cpu.INC_HL_mem(mmu) // Increment memory to 0x11
	assert.Equal(t, uint8(0x11), mmu.ReadByte(0x8000))

	cpu.LD_B_HL(mmu) // Load memory value into B
	assert.Equal(t, uint8(0x11), cpu.B)

	// Test the reverse - decrement memory, then load to different register
	cpu.DEC_HL_mem(mmu) // Decrement memory to 0x10
	assert.Equal(t, uint8(0x10), mmu.ReadByte(0x8000))

	cpu.LD_C_HL(mmu) // Load memory value into C
	assert.Equal(t, uint8(0x10), cpu.C)
}

func TestMemoryOperationsBoundaryConditions(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test operations at different memory addresses
	addresses := []uint16{0x0000, 0x7FFF, 0x8000, 0x9FFF, 0xC000, 0xFFFE}

	for _, addr := range addresses {
		t.Run(fmt.Sprintf("Memory operations at 0x%04X", addr), func(t *testing.T) {
			cpu.SetHL(addr)

			// Store, increment, load sequence
			cpu.LD_HL_mem_n(mmu, 0x42)
			assert.Equal(t, uint8(0x42), mmu.ReadByte(addr), "Should store at address 0x%04X", addr)

			cpu.INC_HL_mem(mmu)
			assert.Equal(t, uint8(0x43), mmu.ReadByte(addr), "Should increment at address 0x%04X", addr)

			cpu.LD_D_HL(mmu)
			assert.Equal(t, uint8(0x43), cpu.D, "Should load from address 0x%04X", addr)
		})
	}
}

func TestAllMemoryLoadOperations(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test all LD r,(HL) operations work correctly
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x99)

	loadTests := []struct {
		name     string
		loadFunc func() uint8
		checkReg func() uint8
	}{
		{"LD B,(HL)", func() uint8 { return cpu.LD_B_HL(mmu) }, func() uint8 { return cpu.B }},
		{"LD C,(HL)", func() uint8 { return cpu.LD_C_HL(mmu) }, func() uint8 { return cpu.C }},
		{"LD D,(HL)", func() uint8 { return cpu.LD_D_HL(mmu) }, func() uint8 { return cpu.D }},
		{"LD E,(HL)", func() uint8 { return cpu.LD_E_HL(mmu) }, func() uint8 { return cpu.E }},
		{"LD H,(HL)", func() uint8 { return cpu.LD_H_HL(mmu) }, func() uint8 { return cpu.H }},
		{"LD L,(HL)", func() uint8 { return cpu.LD_L_HL(mmu) }, func() uint8 { return cpu.L }},
	}

	for _, test := range loadTests {
		t.Run(test.name, func(t *testing.T) {
			// Reset HL for H and L tests since they modify the address
			cpu.SetHL(0x8000)
			mmu.WriteByte(0x8000, 0x99)

			cycles := test.loadFunc()
			assert.Equal(t, uint8(8), cycles, "%s should take 8 cycles", test.name)
			assert.Equal(t, uint8(0x99), test.checkReg(), "%s should load 0x99 into register", test.name)
		})
	}
}

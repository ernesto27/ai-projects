package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Step 1: Test for wrapNOP Function ===
// Let's test our first wrapper function to make sure it works correctly

func TestWrapNOP(t *testing.T) {
	// === Setup ===
	// Create a CPU instance - like getting a worker ready for work
	cpu := NewCPU()

	// Create an MMU instance - like setting up the filing cabinet
	mmu := memory.NewMMU()

	// Store the initial state to compare later
	initialA := cpu.A
	initialF := cpu.F
	initialPC := cpu.PC

	// === Test 1: Basic Functionality ===
	// Call our wrapper function
	cycles, err := wrapNOP(cpu, mmu)

	// Check that it worked correctly - much cleaner with assert!
	assert.NoError(t, err, "wrapNOP should not return an error")
	assert.Equal(t, uint8(4), cycles, "wrapNOP should return 4 cycles")

	// === Test 2: State Preservation ===
	// NOP should not change any registers
	assert.Equal(t, initialA, cpu.A, "NOP should not change A register")
	assert.Equal(t, initialF, cpu.F, "NOP should not change F register")
	assert.Equal(t, initialPC, cpu.PC, "NOP should not change PC register")

	// === Test 3: Parameters Don't Matter ===
	// NOP should work the same regardless of what parameters we pass
	cycles1, err1 := wrapNOP(cpu, mmu)             // No parameters
	cycles2, err2 := wrapNOP(cpu, mmu, 0x42)       // One parameter
	cycles3, err3 := wrapNOP(cpu, mmu, 0x42, 0x43) // Two parameters

	// All should work the same
	assert.NoError(t, err1, "wrapNOP should work with no parameters")
	assert.NoError(t, err2, "wrapNOP should work with one parameter")
	assert.NoError(t, err3, "wrapNOP should work with two parameters")

	assert.Equal(t, uint8(4), cycles1, "wrapNOP should return 4 cycles with no params")
	assert.Equal(t, uint8(4), cycles2, "wrapNOP should return 4 cycles with one param")
	assert.Equal(t, uint8(4), cycles3, "wrapNOP should return 4 cycles with two params")
}

// === Step 2: Comparison Test ===
// Let's make sure our wrapper behaves exactly like the original function

func TestWrapNOPVsOriginal(t *testing.T) {
	// Create two identical CPUs
	cpu1 := NewCPU()
	cpu2 := NewCPU()
	mmu := memory.NewMMU()

	// Call the original function on cpu1
	originalCycles := cpu1.NOP()

	// Call our wrapper function on cpu2
	wrapperCycles, err := wrapNOP(cpu2, mmu)

	// They should behave exactly the same
	assert.NoError(t, err, "wrapper should not return error when original doesn't")
	assert.Equal(t, originalCycles, wrapperCycles, "wrapper should return same cycles as original")

	// Both CPUs should be in the same state
	assert.Equal(t, cpu1.A, cpu2.A, "Both CPUs should have same A register")
	assert.Equal(t, cpu1.F, cpu2.F, "Both CPUs should have same F register")
	assert.Equal(t, cpu1.PC, cpu2.PC, "Both CPUs should have same PC register")
}

// === Step 2: Test for wrapINC_A Function ===
// Let's test our second wrapper function - this one actually changes CPU state!

func TestWrapINC_A(t *testing.T) {
	// === Setup ===
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set A to a known value to test increment
	cpu.A = 0x42
	expectedA := cpu.A + 1 // After increment, A should be 0x43

	// === Test 1: Basic Functionality ===
	cycles, err := wrapINC_A(cpu, mmu)

	// Check that it worked correctly
	assert.NoError(t, err, "wrapINC_A should not return an error")
	assert.Equal(t, uint8(4), cycles, "wrapINC_A should return 4 cycles")
	assert.Equal(t, expectedA, cpu.A, "INC_A should increment A register by 1")

	// === Test 2: Flag Effects ===
	// INC_A affects flags, so let's test that

	// Test Zero flag when result is zero
	cpu.A = 0xFF // When we increment 0xFF, it wraps to 0x00
	_, err = wrapINC_A(cpu, mmu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x00), cpu.A, "0xFF + 1 should wrap to 0x00")
	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")

	// Test Half-carry flag
	cpu.A = 0x0F // When we increment 0x0F, we get half-carry from bit 3 to 4
	_, err = wrapINC_A(cpu, mmu)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0x10), cpu.A, "0x0F + 1 should be 0x10")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be set when carry from bit 3 to 4")

	// Test N flag is always cleared for INC
	assert.False(t, cpu.GetFlag(FlagN), "N flag should always be cleared for INC")

	// === Test 3: Parameters Don't Matter ===
	cpu.A = 0x10                         // Reset to known value
	cycles1, err1 := wrapINC_A(cpu, mmu) // No parameters

	cpu.A = 0x10                               // Reset to same value
	cycles2, err2 := wrapINC_A(cpu, mmu, 0x42) // One parameter (ignored)

	// Both should work the same
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, uint8(4), cycles1)
	assert.Equal(t, uint8(4), cycles2)
	// Note: We can't compare A values here because both calls modified A
}

func TestWrapINC_AVsOriginal(t *testing.T) {
	// Create two identical CPUs
	cpu1 := NewCPU()
	cpu2 := NewCPU()
	mmu := memory.NewMMU()

	// Set both to same starting value
	cpu1.A = 0x42
	cpu2.A = 0x42

	// Call the original function on cpu1
	originalCycles := cpu1.INC_A()

	// Call our wrapper function on cpu2
	wrapperCycles, err := wrapINC_A(cpu2, mmu)

	// They should behave exactly the same
	assert.NoError(t, err, "wrapper should not return error when original doesn't")
	assert.Equal(t, originalCycles, wrapperCycles, "wrapper should return same cycles as original")

	// Both CPUs should be in the same state after the operation
	assert.Equal(t, cpu1.A, cpu2.A, "Both CPUs should have same A register")
	assert.Equal(t, cpu1.F, cpu2.F, "Both CPUs should have same F (flags) register")
}

package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for Easy Wrapper Functions ===
// These test all the wrapper functions that don't need MMU or parameters

// === Test Decrement Wrapper Functions ===

func TestWrapDEC_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set A to a known value
	cpu.A = 0x42
	expectedA := cpu.A - 1

	// Test the wrapper
	cycles, err := wrapDEC_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles, "DEC_A should return 4 cycles")
	assert.Equal(t, expectedA, cpu.A, "A should be decremented by 1")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should be set for DEC")
}

func TestWrapDEC_B(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.B = 0x42
	expectedB := cpu.B - 1

	cycles, err := wrapDEC_B(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedB, cpu.B)
	assert.True(t, cpu.GetFlag(FlagN))
}

func TestWrapDEC_C(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.C = 0x42
	expectedC := cpu.C - 1

	cycles, err := wrapDEC_C(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedC, cpu.C)
	assert.True(t, cpu.GetFlag(FlagN))
}

func TestWrapDEC_D(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.D = 0x42
	expectedD := cpu.D - 1

	cycles, err := wrapDEC_D(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedD, cpu.D)
	assert.True(t, cpu.GetFlag(FlagN))
}

func TestWrapDEC_E(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.E = 0x42
	expectedE := cpu.E - 1

	cycles, err := wrapDEC_E(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedE, cpu.E)
	assert.True(t, cpu.GetFlag(FlagN))
}

func TestWrapDEC_H(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.H = 0x42
	expectedH := cpu.H - 1

	cycles, err := wrapDEC_H(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedH, cpu.H)
	assert.True(t, cpu.GetFlag(FlagN))
}

func TestWrapDEC_L(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.L = 0x42
	expectedL := cpu.L - 1

	cycles, err := wrapDEC_L(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedL, cpu.L)
	assert.True(t, cpu.GetFlag(FlagN))
}

// === Test Increment Wrapper Functions ===

func TestWrapINC_B(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.B = 0x42
	expectedB := cpu.B + 1

	cycles, err := wrapINC_B(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedB, cpu.B)
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear for INC")
}

func TestWrapINC_C(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.C = 0x42
	expectedC := cpu.C + 1

	cycles, err := wrapINC_C(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedC, cpu.C)
	assert.False(t, cpu.GetFlag(FlagN))
}

func TestWrapINC_D(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.D = 0x42
	expectedD := cpu.D + 1

	cycles, err := wrapINC_D(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedD, cpu.D)
	assert.False(t, cpu.GetFlag(FlagN))
}

func TestWrapINC_E(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.E = 0x42
	expectedE := cpu.E + 1

	cycles, err := wrapINC_E(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedE, cpu.E)
	assert.False(t, cpu.GetFlag(FlagN))
}

func TestWrapINC_H(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.H = 0x42
	expectedH := cpu.H + 1

	cycles, err := wrapINC_H(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedH, cpu.H)
	assert.False(t, cpu.GetFlag(FlagN))
}

func TestWrapINC_L(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.L = 0x42
	expectedL := cpu.L + 1

	cycles, err := wrapINC_L(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, expectedL, cpu.L)
	assert.False(t, cpu.GetFlag(FlagN))
}

// === Test Register Load Wrapper Functions ===

func TestWrapLD_A_B(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.A = 0x00
	cpu.B = 0x42

	cycles, err := wrapLD_A_B(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x42), cpu.A, "A should get B's value")
	assert.Equal(t, uint8(0x42), cpu.B, "B should remain unchanged")
}

func TestWrapLD_A_C(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.A = 0x00
	cpu.C = 0x55

	cycles, err := wrapLD_A_C(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x55), cpu.A)
	assert.Equal(t, uint8(0x55), cpu.C)
}

func TestWrapLD_B_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.A = 0x99
	cpu.B = 0x00

	cycles, err := wrapLD_B_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x99), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint8(0x99), cpu.B, "B should get A's value")
}

// === Test Arithmetic Wrapper Functions ===

func TestWrapADD_A_B(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.A = 0x10
	cpu.B = 0x20

	cycles, err := wrapADD_A_B(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x30), cpu.A, "A should contain sum")
	assert.Equal(t, uint8(0x20), cpu.B, "B should remain unchanged")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear for ADD")
}

func TestWrapADD_A_A(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.A = 0x21

	cycles, err := wrapADD_A_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x42), cpu.A, "A should contain double its original value")
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be clear for ADD")
}

// === Comparison Test: Wrapper vs Original ===
// This tests that our wrappers behave exactly like the original functions

func TestWrappersVsOriginals(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*CPU)
		wrapper  func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		original func(*CPU) uint8
	}{
		{
			name:     "DEC_A comparison",
			setup:    func(cpu *CPU) { cpu.A = 0x42 },
			wrapper:  wrapDEC_A,
			original: (*CPU).DEC_A,
		},
		{
			name:     "INC_B comparison",
			setup:    func(cpu *CPU) { cpu.B = 0x42 },
			wrapper:  wrapINC_B,
			original: (*CPU).INC_B,
		},
		{
			name:     "LD_A_C comparison",
			setup:    func(cpu *CPU) { cpu.A = 0x00; cpu.C = 0x99 },
			wrapper:  wrapLD_A_C,
			original: (*CPU).LD_A_C,
		},
		{
			name:     "ADD_A_B comparison",
			setup:    func(cpu *CPU) { cpu.A = 0x10; cpu.B = 0x20 },
			wrapper:  wrapADD_A_B,
			original: (*CPU).ADD_A_B,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create two identical CPUs
			cpu1 := NewCPU()
			cpu2 := NewCPU()
			mmu := createTestMMU()

			// Set them up identically
			tt.setup(cpu1)
			tt.setup(cpu2)

			// Call original on cpu1, wrapper on cpu2
			originalCycles := tt.original(cpu1)
			wrapperCycles, err := tt.wrapper(cpu2, mmu)

			// They should behave identically
			assert.NoError(t, err)
			assert.Equal(t, originalCycles, wrapperCycles, "cycles should match")

			// All registers should be identical
			assert.Equal(t, cpu1.A, cpu2.A, "A register should match")
			assert.Equal(t, cpu1.B, cpu2.B, "B register should match")
			assert.Equal(t, cpu1.C, cpu2.C, "C register should match")
			assert.Equal(t, cpu1.D, cpu2.D, "D register should match")
			assert.Equal(t, cpu1.E, cpu2.E, "E register should match")
			assert.Equal(t, cpu1.F, cpu2.F, "F (flags) register should match")
			assert.Equal(t, cpu1.H, cpu2.H, "H register should match")
			assert.Equal(t, cpu1.L, cpu2.L, "L register should match")
		})
	}
}

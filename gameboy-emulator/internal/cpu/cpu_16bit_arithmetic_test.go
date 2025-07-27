package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for 16-bit Increment Operations ===

func TestINC_BC(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal increment
	cpu.SetBC(0x1234)
	cycles := cpu.INC_BC()

	assert.Equal(t, uint8(8), cycles, "INC BC should take 8 cycles")
	assert.Equal(t, uint16(0x1235), cpu.GetBC(), "BC should be incremented by 1")

	// Test 2: Wrap around from 0xFFFF to 0x0000
	cpu.SetBC(0xFFFF)
	cycles = cpu.INC_BC()

	assert.Equal(t, uint8(8), cycles, "INC BC should take 8 cycles")
	assert.Equal(t, uint16(0x0000), cpu.GetBC(), "BC should wrap from 0xFFFF to 0x0000")

	// Test 3: Flags should not be affected
	cpu.SetBC(0x5678)
	// Set all flags to test they're preserved
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.INC_BC()

	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

func TestINC_DE(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal increment
	cpu.SetDE(0xABCD)
	cycles := cpu.INC_DE()

	assert.Equal(t, uint8(8), cycles, "INC DE should take 8 cycles")
	assert.Equal(t, uint16(0xABCE), cpu.GetDE(), "DE should be incremented by 1")

	// Test 2: Wrap around from 0xFFFF to 0x0000
	cpu.SetDE(0xFFFF)
	cycles = cpu.INC_DE()

	assert.Equal(t, uint8(8), cycles, "INC DE should take 8 cycles")
	assert.Equal(t, uint16(0x0000), cpu.GetDE(), "DE should wrap from 0xFFFF to 0x0000")

	// Test 3: Flags should not be affected
	cpu.SetDE(0x1111)
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)

	cpu.INC_DE()

	assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be preserved")
	assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be preserved")
	assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be preserved")
	assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

func TestINC_HL(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal increment
	cpu.SetHL(0x8000)
	cycles := cpu.INC_HL()

	assert.Equal(t, uint8(8), cycles, "INC HL should take 8 cycles")
	assert.Equal(t, uint16(0x8001), cpu.GetHL(), "HL should be incremented by 1")

	// Test 2: Test edge case - increment across byte boundary
	cpu.SetHL(0x00FF)
	cycles = cpu.INC_HL()

	assert.Equal(t, uint8(8), cycles, "INC HL should take 8 cycles")
	assert.Equal(t, uint16(0x0100), cpu.GetHL(), "HL should increment across byte boundary")

	// Test 3: Maximum value wrap
	cpu.SetHL(0xFFFF)
	cpu.INC_HL()
	assert.Equal(t, uint16(0x0000), cpu.GetHL(), "HL should wrap from 0xFFFF to 0x0000")
}

func TestINC_SP(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal increment
	cpu.SP = 0xFFFE // Typical starting stack pointer
	cycles := cpu.INC_SP()

	assert.Equal(t, uint8(8), cycles, "INC SP should take 8 cycles")
	assert.Equal(t, uint16(0xFFFF), cpu.SP, "SP should be incremented by 1")

	// Test 2: Wrap around from 0xFFFF to 0x0000
	cpu.SP = 0xFFFF
	cycles = cpu.INC_SP()

	assert.Equal(t, uint8(8), cycles, "INC SP should take 8 cycles")
	assert.Equal(t, uint16(0x0000), cpu.SP, "SP should wrap from 0xFFFF to 0x0000")

	// Test 3: Increment from 0x0000
	cpu.SP = 0x0000
	cpu.INC_SP()
	assert.Equal(t, uint16(0x0001), cpu.SP, "SP should increment from 0x0000 to 0x0001")
}

// === Tests for 16-bit Decrement Operations ===

func TestDEC_BC(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal decrement
	cpu.SetBC(0x1234)
	cycles := cpu.DEC_BC()

	assert.Equal(t, uint8(8), cycles, "DEC BC should take 8 cycles")
	assert.Equal(t, uint16(0x1233), cpu.GetBC(), "BC should be decremented by 1")

	// Test 2: Wrap around from 0x0000 to 0xFFFF
	cpu.SetBC(0x0000)
	cycles = cpu.DEC_BC()

	assert.Equal(t, uint8(8), cycles, "DEC BC should take 8 cycles")
	assert.Equal(t, uint16(0xFFFF), cpu.GetBC(), "BC should wrap from 0x0000 to 0xFFFF")

	// Test 3: Flags should not be affected
	cpu.SetBC(0x5678)
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	cpu.DEC_BC()

	assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagN), "Subtract flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagH), "Half-carry flag should be preserved")
	assert.True(t, cpu.GetFlag(FlagC), "Carry flag should be preserved")
}

func TestDEC_DE(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal decrement
	cpu.SetDE(0xABCD)
	cycles := cpu.DEC_DE()

	assert.Equal(t, uint8(8), cycles, "DEC DE should take 8 cycles")
	assert.Equal(t, uint16(0xABCC), cpu.GetDE(), "DE should be decremented by 1")

	// Test 2: Decrement across byte boundary
	cpu.SetDE(0x0100)
	cycles = cpu.DEC_DE()

	assert.Equal(t, uint8(8), cycles, "DEC DE should take 8 cycles")
	assert.Equal(t, uint16(0x00FF), cpu.GetDE(), "DE should decrement across byte boundary")

	// Test 3: Wrap around from 0x0000 to 0xFFFF
	cpu.SetDE(0x0000)
	cpu.DEC_DE()
	assert.Equal(t, uint16(0xFFFF), cpu.GetDE(), "DE should wrap from 0x0000 to 0xFFFF")
}

func TestDEC_HL(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal decrement
	cpu.SetHL(0x8001)
	cycles := cpu.DEC_HL()

	assert.Equal(t, uint8(8), cycles, "DEC HL should take 8 cycles")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL should be decremented by 1")

	// Test 2: Decrement from 0x0001 to 0x0000
	cpu.SetHL(0x0001)
	cycles = cpu.DEC_HL()

	assert.Equal(t, uint8(8), cycles, "DEC HL should take 8 cycles")
	assert.Equal(t, uint16(0x0000), cpu.GetHL(), "HL should decrement to 0x0000")

	// Test 3: Wrap around
	cpu.SetHL(0x0000)
	cpu.DEC_HL()
	assert.Equal(t, uint16(0xFFFF), cpu.GetHL(), "HL should wrap from 0x0000 to 0xFFFF")
}

func TestDEC_SP(t *testing.T) {
	cpu := NewCPU()

	// Test 1: Normal decrement
	cpu.SP = 0xFFFE
	cycles := cpu.DEC_SP()

	assert.Equal(t, uint8(8), cycles, "DEC SP should take 8 cycles")
	assert.Equal(t, uint16(0xFFFD), cpu.SP, "SP should be decremented by 1")

	// Test 2: Wrap around from 0x0000 to 0xFFFF
	cpu.SP = 0x0000
	cycles = cpu.DEC_SP()

	assert.Equal(t, uint8(8), cycles, "DEC SP should take 8 cycles")
	assert.Equal(t, uint16(0xFFFF), cpu.SP, "SP should wrap from 0x0000 to 0xFFFF")

	// Test 3: Decrement from 0x0001
	cpu.SP = 0x0001
	cpu.DEC_SP()
	assert.Equal(t, uint16(0x0000), cpu.SP, "SP should decrement from 0x0001 to 0x0000")
}

// === Tests for Wrapper Functions ===

func TestWrapINC_BC(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.SetBC(0x1000)
	cycles, err := wrapINC_BC(cpu, mmu)

	assert.NoError(t, err, "wrapINC_BC should not return error")
	assert.Equal(t, uint8(8), cycles, "wrapINC_BC should return 8 cycles")
	assert.Equal(t, uint16(0x1001), cpu.GetBC(), "BC should be incremented")

	// Test that parameters don't matter (wrapper ignores them)
	cycles, err = wrapINC_BC(cpu, mmu, 0x42, 0x43)
	assert.NoError(t, err, "wrapINC_BC should work with parameters")
	assert.Equal(t, uint8(8), cycles, "wrapINC_BC should return 8 cycles with params")
}

func TestWrapDEC_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	cpu.SetHL(0x2000)
	cycles, err := wrapDEC_HL(cpu, mmu)

	assert.NoError(t, err, "wrapDEC_HL should not return error")
	assert.Equal(t, uint8(8), cycles, "wrapDEC_HL should return 8 cycles")
	assert.Equal(t, uint16(0x1FFF), cpu.GetHL(), "HL should be decremented")
}

// === Integration Tests ===

func Test16BitArithmeticIntegration(t *testing.T) {
	cpu := NewCPU()

	// Test increment/decrement round trip
	originalBC := uint16(0x1234)
	cpu.SetBC(originalBC)

	// Increment then decrement should return to original value
	cpu.INC_BC()
	cpu.DEC_BC()
	assert.Equal(t, originalBC, cpu.GetBC(), "INC then DEC should return to original value")

	// Test multiple operations
	cpu.SetHL(0x8000)
	cpu.INC_HL()
	cpu.INC_HL()
	cpu.INC_HL()
	assert.Equal(t, uint16(0x8003), cpu.GetHL(), "Three INC_HL should increment by 3")

	cpu.DEC_HL()
	cpu.DEC_HL()
	assert.Equal(t, uint16(0x8001), cpu.GetHL(), "Two DEC_HL should decrement by 2")
}

func Test16BitArithmeticBoundaryConditions(t *testing.T) {
	cpu := NewCPU()

	// Test all possible wrap-around scenarios
	testCases := []struct {
		name        string
		initial     uint16
		operation   func()
		expected    uint16
		description string
	}{
		{
			name:        "INC_BC wrap",
			initial:     0xFFFF,
			operation:   func() { cpu.INC_BC() },
			expected:    0x0000,
			description: "BC increment wrap-around",
		},
		{
			name:        "DEC_DE wrap",
			initial:     0x0000,
			operation:   func() { cpu.DEC_DE() },
			expected:    0xFFFF,
			description: "DE decrement wrap-around",
		},
		{
			name:        "INC_SP boundary",
			initial:     0x7FFF,
			operation:   func() { cpu.INC_SP() },
			expected:    0x8000,
			description: "SP increment across sign boundary",
		},
		{
			name:        "DEC_HL boundary",
			initial:     0x8000,
			operation:   func() { cpu.DEC_HL() },
			expected:    0x7FFF,
			description: "HL decrement across sign boundary",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the specific register
			switch tc.name {
			case "INC_BC wrap":
				cpu.SetBC(tc.initial)
			case "DEC_DE wrap":
				cpu.SetDE(tc.initial)
			case "INC_SP boundary":
				cpu.SP = tc.initial
			case "DEC_HL boundary":
				cpu.SetHL(tc.initial)
			}

			// Perform the operation
			tc.operation()

			// Check the result
			switch tc.name {
			case "INC_BC wrap":
				assert.Equal(t, tc.expected, cpu.GetBC(), tc.description)
			case "DEC_DE wrap":
				assert.Equal(t, tc.expected, cpu.GetDE(), tc.description)
			case "INC_SP boundary":
				assert.Equal(t, tc.expected, cpu.SP, tc.description)
			case "DEC_HL boundary":
				assert.Equal(t, tc.expected, cpu.GetHL(), tc.description)
			}
		})
	}
}

func TestAllWrapperFunctions(t *testing.T) {
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test all wrapper functions work correctly
	wrapperTests := []struct {
		name    string
		wrapper func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		setup   func()
		check   func() bool
	}{
		{
			name:    "wrapINC_BC",
			wrapper: wrapINC_BC,
			setup:   func() { cpu.SetBC(0x1000) },
			check:   func() bool { return cpu.GetBC() == 0x1001 },
		},
		{
			name:    "wrapINC_DE",
			wrapper: wrapINC_DE,
			setup:   func() { cpu.SetDE(0x2000) },
			check:   func() bool { return cpu.GetDE() == 0x2001 },
		},
		{
			name:    "wrapINC_HL",
			wrapper: wrapINC_HL,
			setup:   func() { cpu.SetHL(0x3000) },
			check:   func() bool { return cpu.GetHL() == 0x3001 },
		},
		{
			name:    "wrapINC_SP",
			wrapper: wrapINC_SP,
			setup:   func() { cpu.SP = 0x4000 },
			check:   func() bool { return cpu.SP == 0x4001 },
		},
		{
			name:    "wrapDEC_BC",
			wrapper: wrapDEC_BC,
			setup:   func() { cpu.SetBC(0x1000) },
			check:   func() bool { return cpu.GetBC() == 0x0FFF },
		},
		{
			name:    "wrapDEC_DE",
			wrapper: wrapDEC_DE,
			setup:   func() { cpu.SetDE(0x2000) },
			check:   func() bool { return cpu.GetDE() == 0x1FFF },
		},
		{
			name:    "wrapDEC_HL",
			wrapper: wrapDEC_HL,
			setup:   func() { cpu.SetHL(0x3000) },
			check:   func() bool { return cpu.GetHL() == 0x2FFF },
		},
		{
			name:    "wrapDEC_SP",
			wrapper: wrapDEC_SP,
			setup:   func() { cpu.SP = 0x4000 },
			check:   func() bool { return cpu.SP == 0x3FFF },
		},
	}

	for _, test := range wrapperTests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			cycles, err := test.wrapper(cpu, mmu)

			assert.NoError(t, err, "%s should not return error", test.name)
			assert.Equal(t, uint8(8), cycles, "%s should return 8 cycles", test.name)
			assert.True(t, test.check(), "%s should modify register correctly", test.name)
		})
	}
}

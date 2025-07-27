package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOROperations tests all OR instruction variants with comprehensive scenarios
// OR operations perform bitwise OR between A and operand, storing result in A
// All OR operations follow the same flag pattern:
// Z: Set if result is zero, N: Always reset, H: Always reset, C: Always reset

func TestOR_A_A(t *testing.T) {
	t.Run("OR A,A - Zero value", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00

		cycles := cpu.OR_A_A()

		assert.Equal(t, uint8(0x00), cpu.A, "A should remain 0x00")
		assert.Equal(t, uint8(4), cycles, "OR A,A should take 4 cycles")

		// Flag verification: Z=1, N=0, H=0, C=0
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,A - Non-zero value", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x42

		cycles := cpu.OR_A_A()

		assert.Equal(t, uint8(0x42), cpu.A, "A should remain 0x42")
		assert.Equal(t, uint8(4), cycles, "OR A,A should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is non-zero")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,A - Maximum value", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF

		cycles := cpu.OR_A_A()

		assert.Equal(t, uint8(0xFF), cpu.A, "A should remain 0xFF")
		assert.Equal(t, uint8(4), cycles, "OR A,A should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xFF")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_B(t *testing.T) {
	t.Run("OR A,B - Both zero", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00
		cpu.B = 0x00

		cycles := cpu.OR_A_B()

		assert.Equal(t, uint8(0x00), cpu.A, "A should be 0x00")
		assert.Equal(t, uint8(0x00), cpu.B, "B should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,B should take 4 cycles")

		// Flag verification: Z=1, N=0, H=0, C=0
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,B - Non-overlapping bits", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x0F // Binary: 00001111
		cpu.B = 0xF0 // Binary: 11110000

		cycles := cpu.OR_A_B()

		assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0x0F | 0xF0)")
		assert.Equal(t, uint8(0xF0), cpu.B, "B should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,B should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xFF")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,B - Overlapping bits", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x33 // Binary: 00110011
		cpu.B = 0x55 // Binary: 01010101

		cycles := cpu.OR_A_B()

		assert.Equal(t, uint8(0x77), cpu.A, "A should be 0x77 (0x33 | 0x55)")
		assert.Equal(t, uint8(0x55), cpu.B, "B should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,B should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is non-zero")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,B - Flag preservation test", func(t *testing.T) {
		cpu := NewCPU()
		// Set all flags initially
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)

		cpu.A = 0x42
		cpu.B = 0x18

		cycles := cpu.OR_A_B()

		assert.Equal(t, uint8(0x5A), cpu.A, "A should be 0x5A (0x42 | 0x18)")
		assert.Equal(t, uint8(4), cycles, "OR A,B should take 4 cycles")

		// Verify OR operation overwrites all flags according to its rules
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset (result is 0x5A)")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_C(t *testing.T) {
	t.Run("OR A,C - Bit setting pattern", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x80 // Binary: 10000000
		cpu.C = 0x01 // Binary: 00000001

		cycles := cpu.OR_A_C()

		assert.Equal(t, uint8(0x81), cpu.A, "A should be 0x81 (0x80 | 0x01)")
		assert.Equal(t, uint8(0x01), cpu.C, "C should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,C should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0x81")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_D(t *testing.T) {
	t.Run("OR A,D - Mixed bit pattern", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xAA // Binary: 10101010
		cpu.D = 0x55 // Binary: 01010101

		cycles := cpu.OR_A_D()

		assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0xAA | 0x55)")
		assert.Equal(t, uint8(0x55), cpu.D, "D should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,D should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xFF")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_E(t *testing.T) {
	t.Run("OR A,E - Identical values", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x7E
		cpu.E = 0x7E

		cycles := cpu.OR_A_E()

		assert.Equal(t, uint8(0x7E), cpu.A, "A should remain 0x7E (0x7E | 0x7E)")
		assert.Equal(t, uint8(0x7E), cpu.E, "E should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,E should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0x7E")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_H(t *testing.T) {
	t.Run("OR A,H - High nibble operation", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x08 // Binary: 00001000
		cpu.H = 0x04 // Binary: 00000100

		cycles := cpu.OR_A_H()

		assert.Equal(t, uint8(0x0C), cpu.A, "A should be 0x0C (0x08 | 0x04)")
		assert.Equal(t, uint8(0x04), cpu.H, "H should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,H should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0x0C")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_L(t *testing.T) {
	t.Run("OR A,L - Low nibble operation", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x02 // Binary: 00000010
		cpu.L = 0x01 // Binary: 00000001

		cycles := cpu.OR_A_L()

		assert.Equal(t, uint8(0x03), cpu.A, "A should be 0x03 (0x02 | 0x01)")
		assert.Equal(t, uint8(0x01), cpu.L, "L should remain unchanged")
		assert.Equal(t, uint8(4), cycles, "OR A,L should take 4 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0x03")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_HL(t *testing.T) {
	t.Run("OR A,(HL) - Memory operation", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Set up HL to point to an address
		cpu.SetHL(0x8000)

		// Store a value in memory at HL
		mmu.WriteByte(0x8000, 0x3C)

		// Set A to a value
		cpu.A = 0xC3

		cycles := cpu.OR_A_HL(mmu)

		assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0xC3 | 0x3C)")
		assert.Equal(t, uint8(8), cycles, "OR A,(HL) should take 8 cycles")

		// Verify memory wasn't modified
		assert.Equal(t, uint8(0x3C), mmu.ReadByte(0x8000), "Memory at HL should remain unchanged")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xFF")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,(HL) - Zero result from memory", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Set up HL to point to an address
		cpu.SetHL(0x9000)

		// Store zero in memory at HL
		mmu.WriteByte(0x9000, 0x00)

		// Set A to zero
		cpu.A = 0x00

		cycles := cpu.OR_A_HL(mmu)

		assert.Equal(t, uint8(0x00), cpu.A, "A should be 0x00 (0x00 | 0x00)")
		assert.Equal(t, uint8(8), cycles, "OR A,(HL) should take 8 cycles")

		// Flag verification: Z=1, N=0, H=0, C=0
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

func TestOR_A_n(t *testing.T) {
	t.Run("OR A,n - Immediate bit setting", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x08 // Binary: 00001000

		cycles := cpu.OR_A_n(0x80) // Binary: 10000000

		assert.Equal(t, uint8(0x88), cpu.A, "A should be 0x88 (0x08 | 0x80)")
		assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0x88")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,n - Zero result", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x00

		cycles := cpu.OR_A_n(0x00)

		assert.Equal(t, uint8(0x00), cpu.A, "A should be 0x00 (0x00 | 0x00)")
		assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

		// Flag verification: Z=1, N=0, H=0, C=0
		assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set when result is 0")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,n - Maximum value", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0xFF

		cycles := cpu.OR_A_n(0xFF)

		assert.Equal(t, uint8(0xFF), cpu.A, "A should be 0xFF (0xFF | 0xFF)")
		assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xFF")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})

	t.Run("OR A,n - Practical bit masking", func(t *testing.T) {
		cpu := NewCPU()
		cpu.A = 0x42 // Binary: 01000010

		// Set bits 0 and 7
		cycles := cpu.OR_A_n(0x81) // Binary: 10000001

		assert.Equal(t, uint8(0xC3), cpu.A, "A should be 0xC3 (0x42 | 0x81)")
		assert.Equal(t, uint8(8), cycles, "OR A,n should take 8 cycles")

		// Flag verification: Z=0, N=0, H=0, C=0
		assert.False(t, cpu.GetFlag(FlagZ), "Zero flag should be reset when result is 0xC3")
		assert.False(t, cpu.GetFlag(FlagN), "Subtract flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagH), "Half-carry flag should be reset for OR operations")
		assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be reset for OR operations")
	})
}

// TestOROperationsEdgeCases tests edge cases and boundary conditions
func TestOROperationsEdgeCases(t *testing.T) {
	t.Run("OR operations - All variants produce same result", func(t *testing.T) {
		// Test that all OR variants with the same values produce identical results
		testValue := uint8(0x5A)

		// Test OR A,A
		cpu1 := NewCPU()
		cpu1.A = testValue
		cpu1.OR_A_A()

		// Test OR A,B with B = A
		cpu2 := NewCPU()
		cpu2.A = testValue
		cpu2.B = testValue
		cpu2.OR_A_B()

		assert.Equal(t, cpu1.A, cpu2.A, "OR A,A and OR A,B should produce same result when B=A")
		assert.Equal(t, cpu1.F, cpu2.F, "Flag states should be identical")
	})

	t.Run("OR operations - Bit pattern verification", func(t *testing.T) {
		// Test specific bit patterns that are common in Game Boy programming
		testCases := []struct {
			a        uint8
			operand  uint8
			expected uint8
			desc     string
		}{
			{0x00, 0x01, 0x01, "Set bit 0"},
			{0x01, 0x02, 0x03, "Set bit 1"},
			{0x03, 0x04, 0x07, "Set bit 2"},
			{0x07, 0x08, 0x0F, "Set bit 3"},
			{0x0F, 0x10, 0x1F, "Set bit 4"},
			{0x1F, 0x20, 0x3F, "Set bit 5"},
			{0x3F, 0x40, 0x7F, "Set bit 6"},
			{0x7F, 0x80, 0xFF, "Set bit 7"},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				cpu := NewCPU()
				cpu.A = tc.a

				result := cpu.OR_A_n(tc.operand)

				assert.Equal(t, tc.expected, cpu.A, "OR result should match expected value")
				assert.Equal(t, uint8(8), result, "OR A,n should take 8 cycles")
				assert.Equal(t, tc.expected == 0, cpu.GetFlag(FlagZ), "Zero flag should match result")
			})
		}
	})
}

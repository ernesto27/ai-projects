package cpu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDAA(t *testing.T) {
	tests := []struct {
		name    string
		a       uint8
		n       bool  // Was last operation subtraction?
		h       bool  // Half-carry from last operation
		c       bool  // Carry from last operation
		wantA   uint8
		wantZ   bool
		wantN   bool  // Should remain unchanged
		wantH   bool  // Should be false after DAA
		wantC   bool
	}{
		// Addition cases (N=0)
		{
			name:  "BCD addition: 09 + 01 = 10",
			a:     0x0A, // Result of 0x09 + 0x01 in binary
			n:     false,
			h:     false,
			c:     false,
			wantA: 0x10, // Corrected BCD
			wantZ: false,
			wantN: false,
			wantH: false,
			wantC: false,
		},
		{
			name:  "BCD addition with half-carry: 08 + 05 = 13",
			a:     0x0D, // Result of 0x08 + 0x05 in binary
			n:     false,
			h:     true, // Half-carry occurred
			c:     false,
			wantA: 0x13, // Corrected BCD
			wantZ: false,
			wantN: false,
			wantH: false,
			wantC: false,
		},
		{
			name:  "BCD addition with carry: 99 + 01 = 00",
			a:     0x9A, // Result of 0x99 + 0x01 in binary
			n:     false,
			h:     false,
			c:     true, // Carry occurred
			wantA: 0x00, // Corrected BCD (with carry out)
			wantZ: true,
			wantN: false,
			wantH: false,
			wantC: true,
		},
		{
			name:  "BCD addition: 50 + 49 = 99",
			a:     0x99, // Result is already correct BCD
			n:     false,
			h:     false,
			c:     false,
			wantA: 0x99, // No correction needed
			wantZ: false,
			wantN: false,
			wantH: false,
			wantC: false,
		},
		// Subtraction cases (N=1)
		{
			name:  "BCD subtraction with half-borrow",
			a:     0x00, // Result after subtraction with borrow
			n:     true,
			h:     true, // Half-borrow occurred
			c:     false,
			wantA: 0xFA, // Corrected BCD
			wantZ: false,
			wantN: true,
			wantH: false,
			wantC: false,
		},
		{
			name:  "BCD subtraction with carry (borrow)",
			a:     0x00, // Result after subtraction with borrow
			n:     true,
			h:     false,
			c:     true, // Carry (borrow) occurred
			wantA: 0xA0, // Corrected BCD
			wantZ: false,
			wantN: true,
			wantH: false,
			wantC: true,
		},
		{
			name:  "Zero result",
			a:     0x00,
			n:     false,
			h:     false,
			c:     false,
			wantA: 0x00,
			wantZ: true,
			wantN: false,
			wantH: false,
			wantC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a
			cpu.SetFlag(FlagN, tt.n)
			cpu.SetFlag(FlagH, tt.h)
			cpu.SetFlag(FlagC, tt.c)

			cycles := cpu.DAA()

			assert.Equal(t, uint8(4), cycles, "DAA should take 4 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register result mismatch")
			assert.Equal(t, tt.wantZ, cpu.GetFlag(FlagZ), "Z flag mismatch")
			assert.Equal(t, tt.wantN, cpu.GetFlag(FlagN), "N flag should be preserved")
			assert.Equal(t, tt.wantH, cpu.GetFlag(FlagH), "H flag should be false after DAA")
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC), "C flag mismatch")
		})
	}
}

func TestCPL(t *testing.T) {
	tests := []struct {
		name    string
		a       uint8
		wantA   uint8
		initialZ bool
		initialC bool
	}{
		{
			name:     "Complement 0x00",
			a:        0x00,
			wantA:    0xFF,
			initialZ: false,
			initialC: false,
		},
		{
			name:     "Complement 0xFF",
			a:        0xFF,
			wantA:    0x00,
			initialZ: true,
			initialC: true,
		},
		{
			name:     "Complement 0x42",
			a:        0x42, // 0b01000010
			wantA:    0xBD, // 0b10111101
			initialZ: false,
			initialC: false,
		},
		{
			name:     "Complement 0xAA",
			a:        0xAA, // 0b10101010
			wantA:    0x55, // 0b01010101
			initialZ: true,
			initialC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a
			cpu.SetFlag(FlagZ, tt.initialZ)
			cpu.SetFlag(FlagC, tt.initialC)

			cycles := cpu.CPL()

			assert.Equal(t, uint8(4), cycles, "CPL should take 4 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register result mismatch")
			assert.Equal(t, tt.initialZ, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
			assert.True(t, cpu.GetFlag(FlagN), "N flag should be set")
			assert.True(t, cpu.GetFlag(FlagH), "H flag should be set")
			assert.Equal(t, tt.initialC, cpu.GetFlag(FlagC), "C flag should be unchanged")
		})
	}
}

func TestSCF(t *testing.T) {
	tests := []struct {
		name     string
		initialZ bool
		initialN bool
		initialH bool
		initialC bool
	}{
		{
			name:     "All flags initially false",
			initialZ: false,
			initialN: false,
			initialH: false,
			initialC: false,
		},
		{
			name:     "All flags initially true",
			initialZ: true,
			initialN: true,
			initialH: true,
			initialC: true,
		},
		{
			name:     "Mixed initial flags",
			initialZ: true,
			initialN: false,
			initialH: true,
			initialC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetFlag(FlagZ, tt.initialZ)
			cpu.SetFlag(FlagN, tt.initialN)
			cpu.SetFlag(FlagH, tt.initialH)
			cpu.SetFlag(FlagC, tt.initialC)

			cycles := cpu.SCF()

			assert.Equal(t, uint8(4), cycles, "SCF should take 4 cycles")
			assert.Equal(t, tt.initialZ, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
			assert.False(t, cpu.GetFlag(FlagN), "N flag should be cleared")
			assert.False(t, cpu.GetFlag(FlagH), "H flag should be cleared")
			assert.True(t, cpu.GetFlag(FlagC), "C flag should be set")
		})
	}
}

func TestCCF(t *testing.T) {
	tests := []struct {
		name     string
		initialZ bool
		initialN bool
		initialH bool
		initialC bool
		wantC    bool
	}{
		{
			name:     "Flip carry from false to true",
			initialZ: false,
			initialN: false,
			initialH: false,
			initialC: false,
			wantC:    true,
		},
		{
			name:     "Flip carry from true to false",
			initialZ: true,
			initialN: true,
			initialH: true,
			initialC: true,
			wantC:    false,
		},
		{
			name:     "Flip carry with mixed flags",
			initialZ: true,
			initialN: false,
			initialH: true,
			initialC: false,
			wantC:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetFlag(FlagZ, tt.initialZ)
			cpu.SetFlag(FlagN, tt.initialN)
			cpu.SetFlag(FlagH, tt.initialH)
			cpu.SetFlag(FlagC, tt.initialC)

			cycles := cpu.CCF()

			assert.Equal(t, uint8(4), cycles, "CCF should take 4 cycles")
			assert.Equal(t, tt.initialZ, cpu.GetFlag(FlagZ), "Z flag should be unchanged")
			assert.False(t, cpu.GetFlag(FlagN), "N flag should be cleared")
			assert.False(t, cpu.GetFlag(FlagH), "H flag should be cleared")
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC), "C flag should be flipped")
		})
	}
}

func TestFlagOperationsSequence(t *testing.T) {
	// Test a realistic sequence of flag operations
	cpu := NewCPU()
	
	// Simulate BCD arithmetic sequence
	cpu.A = 0x09
	cpu.B = 0x01
	
	// ADD A,B (would set appropriate flags)
	cpu.A = cpu.A + cpu.B  // 0x0A
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)
	
	// DAA to correct BCD result
	cpu.DAA()
	assert.Equal(t, uint8(0x10), cpu.A, "BCD correction should give 0x10")
	
	// CPL to invert all bits
	cpu.CPL()
	assert.Equal(t, uint8(0xEF), cpu.A, "CPL should invert 0x10 to 0xEF")
	assert.True(t, cpu.GetFlag(FlagN), "CPL should set N flag")
	assert.True(t, cpu.GetFlag(FlagH), "CPL should set H flag")
	
	// SCF to set carry
	cpu.SCF()
	assert.True(t, cpu.GetFlag(FlagC), "SCF should set carry flag")
	assert.False(t, cpu.GetFlag(FlagN), "SCF should clear N flag")
	assert.False(t, cpu.GetFlag(FlagH), "SCF should clear H flag")
	
	// CCF to flip carry
	cpu.CCF()
	assert.False(t, cpu.GetFlag(FlagC), "CCF should flip carry flag to false")
	assert.False(t, cpu.GetFlag(FlagN), "CCF should clear N flag")
	assert.False(t, cpu.GetFlag(FlagH), "CCF should clear H flag")
	
	// CCF again to flip back
	cpu.CCF()
	assert.True(t, cpu.GetFlag(FlagC), "CCF should flip carry flag back to true")
}

func TestBCDRealWorldScenario(t *testing.T) {
	// Test realistic BCD usage: adding scores in a game
	cpu := NewCPU()
	
	// Player has score 0x99 (99 points in BCD)
	// They get 0x02 (2 more points in BCD)
	// Result should be 0x01 with carry (101 points, displaying as "01" with hundreds digit elsewhere)
	
	cpu.A = 0x99
	cpu.B = 0x02
	
	// Simulate ADD A,B
	result := uint16(cpu.A) + uint16(cpu.B) // 0x9B
	cpu.A = uint8(result)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagH, false) // 0x9 + 0x2 = 0xB > 0x9, so half-carry
	cpu.SetFlag(FlagC, false) // 0x99 + 0x02 = 0x9B < 0x100, so no carry
	
	// But wait, we need to check half-carry properly for BCD
	// In BCD context, half-carry means lower nibble > 9 after addition
	lowerNibble := (cpu.A & 0x0F)
	if lowerNibble > 0x09 {
		cpu.SetFlag(FlagH, true)
	}
	
	// Check if upper nibble > 9 (which it is: 0x9B means upper nibble B = 11)
	if cpu.A > 0x99 {
		cpu.SetFlag(FlagC, true)
	}
	
	cpu.DAA()
	assert.Equal(t, uint8(0x01), cpu.A, "BCD 99 + 02 should give 01 with carry")
	assert.True(t, cpu.GetFlag(FlagC), "Should have carry flag set for overflow")
}
package cpu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRLCA(t *testing.T) {
	tests := []struct {
		name    string
		a       uint8
		wantA   uint8
		wantC   bool
	}{
		{
			name:  "Rotate 0b10110101",
			a:     0b10110101, // 0xB5
			wantA: 0b01101011, // 0x6B - bit 7 moved to bit 0
			wantC: true,       // bit 7 was 1
		},
		{
			name:  "Rotate 0b01110101", 
			a:     0b01110101, // 0x75
			wantA: 0b11101010, // 0xEA - bit 7 moved to bit 0
			wantC: false,      // bit 7 was 0
		},
		{
			name:  "Rotate 0x00",
			a:     0x00,
			wantA: 0x00,
			wantC: false,
		},
		{
			name:  "Rotate 0xFF",
			a:     0xFF,
			wantA: 0xFF, // All bits set, rotation doesn't change it
			wantC: true, // bit 7 was 1
		},
		{
			name:  "Rotate 0x80 (only bit 7 set)",
			a:     0x80,
			wantA: 0x01, // bit 7 goes to bit 0
			wantC: true,
		},
		{
			name:  "Rotate 0x01 (only bit 0 set)",
			a:     0x01,
			wantA: 0x02, // shift left by 1
			wantC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a

			cycles := cpu.RLCA()

			assert.Equal(t, uint8(4), cycles, "Should take 4 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register result mismatch")
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC), "Carry flag mismatch")
			
			// RLCA always clears Z, N, H flags
			assert.False(t, cpu.GetFlag(FlagZ), "Z flag should be 0")
			assert.False(t, cpu.GetFlag(FlagN), "N flag should be 0")
			assert.False(t, cpu.GetFlag(FlagH), "H flag should be 0")
		})
	}
}

func TestRRCA(t *testing.T) {
	tests := []struct {
		name    string
		a       uint8
		wantA   uint8
		wantC   bool
	}{
		{
			name:  "Rotate 0b10110101",
			a:     0b10110101, // 0xB5
			wantA: 0b11011010, // 0xDA - bit 0 moved to bit 7
			wantC: true,       // bit 0 was 1
		},
		{
			name:  "Rotate 0b10110100",
			a:     0b10110100, // 0xB4
			wantA: 0b01011010, // 0x5A - bit 0 moved to bit 7
			wantC: false,      // bit 0 was 0
		},
		{
			name:  "Rotate 0x01 (only bit 0 set)",
			a:     0x01,
			wantA: 0x80, // bit 0 goes to bit 7
			wantC: true,
		},
		{
			name:  "Rotate 0x02 (only bit 1 set)",
			a:     0x02,
			wantA: 0x01, // shift right by 1
			wantC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a

			cycles := cpu.RRCA()

			assert.Equal(t, uint8(4), cycles)
			assert.Equal(t, tt.wantA, cpu.A)
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC))
			assert.False(t, cpu.GetFlag(FlagZ))
			assert.False(t, cpu.GetFlag(FlagN))
			assert.False(t, cpu.GetFlag(FlagH))
		})
	}
}

func TestRLA(t *testing.T) {
	tests := []struct {
		name      string
		a         uint8
		carryIn   bool
		wantA     uint8
		wantC     bool
	}{
		{
			name:    "Rotate with carry=0",
			a:       0b10110101,
			carryIn: false,
			wantA:   0b01101010, // Shift left, carry becomes bit 0
			wantC:   true,       // bit 7 was 1
		},
		{
			name:    "Rotate with carry=1",
			a:       0b10110100,
			carryIn: true,
			wantA:   0b01101001, // Shift left, carry(1) becomes bit 0
			wantC:   true,       // bit 7 was 1
		},
		{
			name:    "Rotate 0x7F with carry=1",
			a:       0x7F,       // 0b01111111
			carryIn: true,
			wantA:   0xFF,       // 0b11111111 - shift left with carry in bit 0
			wantC:   false,      // bit 7 was 0
		},
		{
			name:    "Rotate 0x80 with carry=0",
			a:       0x80,       // 0b10000000
			carryIn: false,
			wantA:   0x00,       // 0b00000000 - shift left with no carry
			wantC:   true,       // bit 7 was 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a
			cpu.SetFlag(FlagC, tt.carryIn)

			cycles := cpu.RLA()

			assert.Equal(t, uint8(4), cycles)
			assert.Equal(t, tt.wantA, cpu.A)
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC))
			assert.False(t, cpu.GetFlag(FlagZ))
			assert.False(t, cpu.GetFlag(FlagN))
			assert.False(t, cpu.GetFlag(FlagH))
		})
	}
}

func TestRRA(t *testing.T) {
	tests := []struct {
		name      string
		a         uint8
		carryIn   bool
		wantA     uint8
		wantC     bool
	}{
		{
			name:    "Rotate with carry=0",
			a:       0b10110101,
			carryIn: false,
			wantA:   0b01011010, // Shift right, carry becomes bit 7
			wantC:   true,       // bit 0 was 1
		},
		{
			name:    "Rotate with carry=1",
			a:       0b10110100,
			carryIn: true,
			wantA:   0b11011010, // Shift right, carry(1) becomes bit 7
			wantC:   false,      // bit 0 was 0
		},
		{
			name:    "Rotate 0xFE with carry=1",
			a:       0xFE,       // 0b11111110
			carryIn: true,
			wantA:   0xFF,       // 0b11111111 - shift right with carry in bit 7
			wantC:   false,      // bit 0 was 0
		},
		{
			name:    "Rotate 0x01 with carry=0",
			a:       0x01,       // 0b00000001
			carryIn: false,
			wantA:   0x00,       // 0b00000000 - shift right with no carry
			wantC:   true,       // bit 0 was 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.A = tt.a
			cpu.SetFlag(FlagC, tt.carryIn)

			cycles := cpu.RRA()

			assert.Equal(t, uint8(4), cycles)
			assert.Equal(t, tt.wantA, cpu.A)
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC))
			assert.False(t, cpu.GetFlag(FlagZ))
			assert.False(t, cpu.GetFlag(FlagN))
			assert.False(t, cpu.GetFlag(FlagH))
		})
	}
}

func TestRotationSequence(t *testing.T) {
	// Test a sequence of rotations to verify they work together correctly
	cpu := NewCPU()
	cpu.A = 0b10110101 // Start with this pattern

	// RLCA: 10110101 -> 01101011, C=1
	cpu.RLCA()
	assert.Equal(t, uint8(0b01101011), cpu.A)
	assert.True(t, cpu.GetFlag(FlagC))

	// RLA: 01101011 -> 11010110 (carry=1 goes to bit 0), C=0
	cpu.RLA()
	assert.Equal(t, uint8(0b11010111), cpu.A) // bit 0 = old carry (1)
	assert.False(t, cpu.GetFlag(FlagC))       // bit 7 was 0

	// RRCA: 11010111 -> 11101011, C=1
	cpu.RRCA()
	assert.Equal(t, uint8(0b11101011), cpu.A)
	assert.True(t, cpu.GetFlag(FlagC))

	// RRA: 11101011 -> 11110101 (carry=1 goes to bit 7), C=1
	cpu.RRA()
	assert.Equal(t, uint8(0b11110101), cpu.A)
	assert.True(t, cpu.GetFlag(FlagC)) // bit 0 was 1
}
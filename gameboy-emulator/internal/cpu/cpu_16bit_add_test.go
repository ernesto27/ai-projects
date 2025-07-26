package cpu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestADD_HL_BC(t *testing.T) {
	tests := []struct {
		name     string
		hl       uint16
		bc       uint16
		wantHL   uint16
		wantH    bool // Half-carry flag
		wantC    bool // Carry flag
		wantN    bool // N flag (should always be false)
	}{
		{
			name:   "Basic addition",
			hl:     0x1234,
			bc:     0x0056,
			wantHL: 0x128A,
			wantH:  false,
			wantC:  false,
			wantN:  false,
		},
		{
			name:   "Addition with carry",
			hl:     0xFFFF,
			bc:     0x0001,
			wantHL: 0x0000,
			wantH:  true,  // Carry from bit 11
			wantC:  true,  // Carry from bit 15
			wantN:  false,
		},
		{
			name:   "Addition with half-carry only",
			hl:     0x0FFF,
			bc:     0x0001,
			wantHL: 0x1000,
			wantH:  true,  // Carry from bit 11 to bit 12
			wantC:  false, // No carry from bit 15
			wantN:  false,
		},
		{
			name:   "Zero + zero",
			hl:     0x0000,
			bc:     0x0000,
			wantHL: 0x0000,
			wantH:  false,
			wantC:  false,
			wantN:  false,
		},
		{
			name:   "Maximum values",
			hl:     0x8000,
			bc:     0x8000,
			wantHL: 0x0000,
			wantH:  false,
			wantC:  true, // Overflow
			wantN:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetHL(tt.hl)
			cpu.SetBC(tt.bc)

			cycles := cpu.ADD_HL_BC()

			assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
			assert.Equal(t, tt.wantHL, cpu.GetHL(), "HL result mismatch")
			assert.Equal(t, tt.wantH, cpu.GetFlag(FlagH), "Half-carry flag mismatch")
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC), "Carry flag mismatch")
			assert.Equal(t, tt.wantN, cpu.GetFlag(FlagN), "N flag should always be false")
		})
	}
}

func TestADD_HL_DE(t *testing.T) {
	cpu := NewCPU()
	cpu.SetHL(0x1000)
	cpu.SetDE(0x0234)

	cycles := cpu.ADD_HL_DE()

	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x1234), cpu.GetHL())
	assert.False(t, cpu.GetFlag(FlagN), "N flag should be false")
}

func TestADD_HL_HL(t *testing.T) {
	tests := []struct {
		name   string
		hl     uint16
		wantHL uint16
		wantH  bool
		wantC  bool
	}{
		{
			name:   "Double small value",
			hl:     0x1234,
			wantHL: 0x2468,
			wantH:  false,
			wantC:  false,
		},
		{
			name:   "Double with overflow",
			hl:     0x8000,
			wantHL: 0x0000,
			wantH:  false,
			wantC:  true,
		},
		{
			name:   "Double with half-carry",
			hl:     0x0800,
			wantHL: 0x1000,
			wantH:  true, // 0x800 + 0x800 = 0x1000, carry from bit 11
			wantC:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetHL(tt.hl)

			cycles := cpu.ADD_HL_HL()

			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.wantHL, cpu.GetHL())
			assert.Equal(t, tt.wantH, cpu.GetFlag(FlagH))
			assert.Equal(t, tt.wantC, cpu.GetFlag(FlagC))
			assert.False(t, cpu.GetFlag(FlagN))
		})
	}
}

func TestADD_HL_SP(t *testing.T) {
	cpu := NewCPU()
	cpu.SetHL(0x2000)
	cpu.SP = 0x1000

	cycles := cpu.ADD_HL_SP()

	assert.Equal(t, uint8(8), cycles)
	assert.Equal(t, uint16(0x3000), cpu.GetHL())
	assert.False(t, cpu.GetFlag(FlagN))
}

func TestADD_HL_PreservesZeroFlag(t *testing.T) {
	// Test that 16-bit ADD instructions don't affect the Z flag
	cpu := NewCPU()
	
	// Set Z flag to true initially
	cpu.SetFlag(FlagZ, true)
	cpu.SetHL(0x1000)
	cpu.SetBC(0x2000)

	cpu.ADD_HL_BC()

	// Z flag should remain unchanged
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should be preserved")
}

func TestHalfCarryCalculation(t *testing.T) {
	// Test the half-carry calculation in detail
	tests := []struct {
		name string
		hl   uint16
		bc   uint16
		wantH bool
	}{
		{"No half-carry: 0x0000 + 0x0000", 0x0000, 0x0000, false},
		{"No half-carry: 0x0123 + 0x0456", 0x0123, 0x0456, false},
		{"Half-carry: 0x0FFF + 0x0001", 0x0FFF, 0x0001, true},
		{"Half-carry: 0x0800 + 0x0800", 0x0800, 0x0800, true},
		{"Half-carry: 0x07FF + 0x0801", 0x07FF, 0x0801, true},
		{"No half-carry: 0x1000 + 0x1000", 0x1000, 0x1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			cpu.SetHL(tt.hl)
			cpu.SetBC(tt.bc)

			cpu.ADD_HL_BC()

			assert.Equal(t, tt.wantH, cpu.GetFlag(FlagH), "Half-carry calculation incorrect")
		})
	}
}
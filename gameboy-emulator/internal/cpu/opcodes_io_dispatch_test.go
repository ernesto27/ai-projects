package cpu

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"gameboy-emulator/internal/memory"
)

func TestIOOpcodeDispatch(t *testing.T) {
	tests := []struct {
		name       string
		opcode     uint8
		params     []uint8
		setupCPU   func(*CPU)
		setupMMU   func(*memory.MMU)
		verifyCPU  func(*testing.T, *CPU)
		verifyMMU  func(*testing.T, *memory.MMU)
		wantCycles uint8
	}{
		{
			name:   "LDH (n),A - opcode 0xE0",
			opcode: 0xE0,
			params: []uint8{0x40}, // LCD Control register offset
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x91 // LCD on, BG on, sprites on
			},
			setupMMU: func(mmu *memory.MMU) {
				// No setup needed
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x91), cpu.A, "A should be unchanged")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				stored := mmu.ReadByte(0xFF40)
				assert.Equal(t, uint8(0x91), stored, "Data should be at 0xFF40")
			},
			wantCycles: 12,
		},
		{
			name:   "LDH A,(n) - opcode 0xF0",
			opcode: 0xF0,
			params: []uint8{0x44}, // LCD Y-coordinate register offset
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Clear A
			},
			setupMMU: func(mmu *memory.MMU) {
				mmu.WriteByte(0xFF44, 0x90) // VBlank period
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x90), cpu.A, "A should contain 0x90")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				// MMU unchanged
			},
			wantCycles: 12,
		},
		{
			name:   "LD (C),A - opcode 0xE2",
			opcode: 0xE2,
			params: []uint8{}, // No immediate parameters
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x42
				cpu.C = 0x26 // Sound on/off register offset
			},
			setupMMU: func(mmu *memory.MMU) {
				// No setup needed
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x42), cpu.A, "A should be unchanged")
				assert.Equal(t, uint8(0x26), cpu.C, "C should be unchanged")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				stored := mmu.ReadByte(0xFF26)
				assert.Equal(t, uint8(0x42), stored, "Data should be at 0xFF26")
			},
			wantCycles: 8,
		},
		{
			name:   "LD A,(C) - opcode 0xF2",
			opcode: 0xF2,
			params: []uint8{}, // No immediate parameters
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Clear A
				cpu.C = 0x00 // Joypad register offset
			},
			setupMMU: func(mmu *memory.MMU) {
				mmu.WriteByte(0xFF00, 0x0F) // All buttons pressed
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x0F), cpu.A, "A should contain 0x0F")
				assert.Equal(t, uint8(0x00), cpu.C, "C should be unchanged")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				// MMU unchanged
			},
			wantCycles: 8,
		},
		{
			name:   "LD (nn),A - opcode 0xEA",
			opcode: 0xEA,
			params: []uint8{0x00, 0x80}, // 0x8000 (VRAM start, little-endian)
			setupCPU: func(cpu *CPU) {
				cpu.A = 0xAA // Sprite data
			},
			setupMMU: func(mmu *memory.MMU) {
				// No setup needed
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0xAA), cpu.A, "A should be unchanged")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				stored := mmu.ReadByte(0x8000)
				assert.Equal(t, uint8(0xAA), stored, "Data should be at 0x8000")
			},
			wantCycles: 16,
		},
		{
			name:   "LD A,(nn) - opcode 0xFA",
			opcode: 0xFA,
			params: []uint8{0x10, 0xC0}, // 0xC010 (WRAM, little-endian)
			setupCPU: func(cpu *CPU) {
				cpu.A = 0x00 // Clear A
			},
			setupMMU: func(mmu *memory.MMU) {
				mmu.WriteByte(0xC010, 0x55) // Player state data
			},
			verifyCPU: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x55), cpu.A, "A should contain 0x55")
			},
			verifyMMU: func(t *testing.T, mmu *memory.MMU) {
				// MMU unchanged
			},
			wantCycles: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()

			// Setup test conditions
			tt.setupCPU(cpu)
			tt.setupMMU(mmu)

			// Execute the instruction through opcode dispatch
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			// Verify no errors
			assert.NoError(t, err, "ExecuteInstruction should not return error")

			// Verify cycle count
			assert.Equal(t, tt.wantCycles, cycles, "Cycle count should match expected")

			// Verify CPU state
			tt.verifyCPU(t, cpu)

			// Verify MMU state
			tt.verifyMMU(t, mmu)
		})
	}
}

func TestIOOpcodeImplementation(t *testing.T) {
	// Test that all I/O opcodes are now implemented
	ioOpcodes := []uint8{0xE0, 0xF0, 0xE2, 0xF2, 0xEA, 0xFA}

	for _, opcode := range ioOpcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X should be implemented", opcode), func(t *testing.T) {
			assert.True(t, IsOpcodeImplemented(opcode), "I/O opcode should be implemented")
		})
	}
}

func TestIOOpcodeParameterValidation(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Test parameter validation for each I/O instruction
	tests := []struct {
		name       string
		opcode     uint8
		params     []uint8
		wantError  bool
		errorCheck func(error) bool
	}{
		{
			name:      "LDH (n),A with correct params",
			opcode:    0xE0,
			params:    []uint8{0x40},
			wantError: false,
		},
		{
			name:      "LDH (n),A with no params",
			opcode:    0xE0,
			params:    []uint8{},
			wantError: true,
		},
		{
			name:      "LDH (n),A with too many params",
			opcode:    0xE0,
			params:    []uint8{0x40, 0x50},
			wantError: true,
		},
		{
			name:      "LD (C),A with correct params",
			opcode:    0xE2,
			params:    []uint8{},
			wantError: false,
		},
		{
			name:      "LD (C),A with extra params",
			opcode:    0xE2,
			params:    []uint8{0x40},
			wantError: true,
		},
		{
			name:      "LD (nn),A with correct params",
			opcode:    0xEA,
			params:    []uint8{0x00, 0x80},
			wantError: false,
		},
		{
			name:      "LD (nn),A with one param",
			opcode:    0xEA,
			params:    []uint8{0x80},
			wantError: true,
		},
		{
			name:      "LD (nn),A with too many params",
			opcode:    0xEA,
			params:    []uint8{0x00, 0x80, 0x90},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cpu.ExecuteInstruction(mmu, tt.opcode, tt.params...)

			if tt.wantError {
				assert.Error(t, err, "Should return error for invalid parameters")
			} else {
				assert.NoError(t, err, "Should not return error for valid parameters")
			}
		})
	}
}
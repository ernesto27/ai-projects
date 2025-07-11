package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWrapL_RegisterOperations tests all L register wrapper functions
// These tests verify that the wrapper functions correctly call the underlying CPU methods
// and work properly with the opcode dispatch system

// TestWrapLD_A_L tests the wrapLD_A_L wrapper function (0x7D)
func TestWrapLD_A_L(t *testing.T) {
	tests := []struct {
		name     string
		initialL uint8
		initialA uint8
		wantA    uint8
		wantL    uint8
	}{
		{
			name:     "Wrapper correctly calls LD_A_L",
			initialL: 0x42,
			initialA: 0x00,
			wantA:    0x42,
			wantL:    0x42,
		},
		{
			name:     "Wrapper handles zero value",
			initialL: 0x00,
			initialA: 0xFF,
			wantA:    0x00,
			wantL:    0x00,
		},
		{
			name:     "Wrapper handles max value",
			initialL: 0xFF,
			initialA: 0x00,
			wantA:    0xFF,
			wantL:    0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()
			
			cpu.L = tt.initialL
			cpu.A = tt.initialA
			
			cycles, err := wrapLD_A_L(cpu, mmu)
			
			assert.NoError(t, err, "Wrapper should not return an error")
			assert.Equal(t, uint8(4), cycles, "Wrapper should return correct cycle count")
			assert.Equal(t, tt.wantA, cpu.A, "Register A should be updated correctly")
			assert.Equal(t, tt.wantL, cpu.L, "Register L should remain unchanged")
		})
	}
}

// TestWrapLD_B_L tests the wrapLD_B_L wrapper function (0x45)
func TestWrapLD_B_L(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.L = 0x78
	cpu.B = 0x12
	
	cycles, err := wrapLD_B_L(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x78), cpu.B, "B should get L's value")
	assert.Equal(t, uint8(0x78), cpu.L, "L should remain unchanged")
}

// TestWrapLD_C_L tests the wrapLD_C_L wrapper function (0x4D)
func TestWrapLD_C_L(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.L = 0x56
	cpu.C = 0x13 // Default value from NewCPU
	
	cycles, err := wrapLD_C_L(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x56), cpu.C, "C should get L's value")
	assert.Equal(t, uint8(0x56), cpu.L, "L should remain unchanged")
}

// TestWrapLD_L_A tests the wrapLD_L_A wrapper function (0x6F)
func TestWrapLD_L_A(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.A = 0x34
	cpu.L = 0x4D // Default value from NewCPU
	
	cycles, err := wrapLD_L_A(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x34), cpu.L, "L should get A's value")
	assert.Equal(t, uint8(0x34), cpu.A, "A should remain unchanged")
}

// TestWrapLD_L_B tests the wrapLD_L_B wrapper function (0x68)
func TestWrapLD_L_B(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.B = 0x90
	cpu.L = 0x00
	
	cycles, err := wrapLD_L_B(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x90), cpu.L, "L should get B's value")
	assert.Equal(t, uint8(0x90), cpu.B, "B should remain unchanged")
}

// TestWrapLD_L_C tests the wrapLD_L_C wrapper function (0x69)
func TestWrapLD_L_C(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.C = 0xAB
	cpu.L = 0xFF
	
	cycles, err := wrapLD_L_C(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0xAB), cpu.L, "L should get C's value")
	assert.Equal(t, uint8(0xAB), cpu.C, "C should remain unchanged")
}

// TestWrapLD_L_D tests the wrapLD_L_D wrapper function (0x6A)
func TestWrapLD_L_D(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.D = 0xCD
	cpu.L = 0x11
	
	cycles, err := wrapLD_L_D(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0xCD), cpu.L, "L should get D's value")
	assert.Equal(t, uint8(0xCD), cpu.D, "D should remain unchanged")
}

// TestWrapLD_L_E tests the wrapLD_L_E wrapper function (0x6B)
func TestWrapLD_L_E(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.E = 0xEF
	cpu.L = 0x22
	
	cycles, err := wrapLD_L_E(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0xEF), cpu.L, "L should get E's value")
	assert.Equal(t, uint8(0xEF), cpu.E, "E should remain unchanged")
}

// TestWrapLD_L_H tests the wrapLD_L_H wrapper function (0x6C)
func TestWrapLD_L_H(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	cpu.H = 0x12
	cpu.L = 0x34
	
	cycles, err := wrapLD_L_H(cpu, mmu)
	
	assert.NoError(t, err)
	assert.Equal(t, uint8(4), cycles)
	assert.Equal(t, uint8(0x12), cpu.L, "L should get H's value")
	assert.Equal(t, uint8(0x12), cpu.H, "H should remain unchanged")
	assert.Equal(t, uint16(0x1212), cpu.GetHL(), "HL should now be 0x1212")
}

// TestL_WrapperVsOriginal_Comparison tests that wrapper functions produce identical results to direct calls
func TestL_WrapperVsOriginal_Comparison(t *testing.T) {
	tests := []struct {
		name         string
		opcode       uint8
		wrapperFunc  func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		directFunc   func(*CPU) uint8
		setupFunc    func(*CPU)
		checkFunc    func(*testing.T, *CPU, *CPU)
	}{
		{
			name:        "LD A,L comparison",
			opcode:      0x7D,
			wrapperFunc: wrapLD_A_L,
			directFunc:  (*CPU).LD_A_L,
			setupFunc: func(cpu *CPU) {
				cpu.L = 0x42
				cpu.A = 0x00
			},
			checkFunc: func(t *testing.T, cpu1, cpu2 *CPU) {
				assert.Equal(t, cpu1.A, cpu2.A, "A register should match")
				assert.Equal(t, cpu1.L, cpu2.L, "L register should match")
			},
		},
		{
			name:        "LD B,L comparison",
			opcode:      0x45,
			wrapperFunc: wrapLD_B_L,
			directFunc:  (*CPU).LD_B_L,
			setupFunc: func(cpu *CPU) {
				cpu.L = 0x78
				cpu.B = 0x12
			},
			checkFunc: func(t *testing.T, cpu1, cpu2 *CPU) {
				assert.Equal(t, cpu1.B, cpu2.B, "B register should match")
				assert.Equal(t, cpu1.L, cpu2.L, "L register should match")
			},
		},
		{
			name:        "LD L,A comparison",
			opcode:      0x6F,
			wrapperFunc: wrapLD_L_A,
			directFunc:  (*CPU).LD_L_A,
			setupFunc: func(cpu *CPU) {
				cpu.A = 0x34
				cpu.L = 0x56
			},
			checkFunc: func(t *testing.T, cpu1, cpu2 *CPU) {
				assert.Equal(t, cpu1.A, cpu2.A, "A register should match")
				assert.Equal(t, cpu1.L, cpu2.L, "L register should match")
			},
		},
		{
			name:        "LD L,H comparison",
			opcode:      0x6C,
			wrapperFunc: wrapLD_L_H,
			directFunc:  (*CPU).LD_L_H,
			setupFunc: func(cpu *CPU) {
				cpu.H = 0x12
				cpu.L = 0x34
			},
			checkFunc: func(t *testing.T, cpu1, cpu2 *CPU) {
				assert.Equal(t, cpu1.H, cpu2.H, "H register should match")
				assert.Equal(t, cpu1.L, cpu2.L, "L register should match")
				assert.Equal(t, cpu1.GetHL(), cpu2.GetHL(), "HL register pair should match")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test wrapper function
			cpu1 := NewCPU()
			mmu := memory.NewMMU()
			tt.setupFunc(cpu1)
			
			cycles1, err := tt.wrapperFunc(cpu1, mmu)
			assert.NoError(t, err, "Wrapper function should not return an error")
			
			// Test direct function call
			cpu2 := NewCPU()
			tt.setupFunc(cpu2)
			
			cycles2 := tt.directFunc(cpu2)
			
			// Compare results
			assert.Equal(t, cycles2, cycles1, "Cycle counts should match")
			tt.checkFunc(t, cpu1, cpu2)
			
			// Verify flags are identical
			assert.Equal(t, cpu1.F, cpu2.F, "Flags should match")
		})
	}
}

// TestL_DispatchIntegration tests the L register operations through the opcode dispatch system
func TestL_DispatchIntegration(t *testing.T) {
	tests := []struct {
		name    string
		opcode  uint8
		setup   func(*CPU)
		verify  func(*testing.T, *CPU)
	}{
		{
			name:   "Dispatch LD A,L (0x7D)",
			opcode: 0x7D,
			setup: func(cpu *CPU) {
				cpu.L = 0x99
				cpu.A = 0x00
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x99), cpu.A, "A should get L's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD B,L (0x45)",
			opcode: 0x45,
			setup: func(cpu *CPU) {
				cpu.L = 0x88
				cpu.B = 0x11
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x88), cpu.B, "B should get L's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD C,L (0x4D)",
			opcode: 0x4D,
			setup: func(cpu *CPU) {
				cpu.L = 0x77
				cpu.C = 0x22
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x77), cpu.C, "C should get L's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,A (0x6F)",
			opcode: 0x6F,
			setup: func(cpu *CPU) {
				cpu.A = 0x66
				cpu.L = 0x33
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x66), cpu.L, "L should get A's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,B (0x68)",
			opcode: 0x68,
			setup: func(cpu *CPU) {
				cpu.B = 0x55
				cpu.L = 0x44
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x55), cpu.L, "L should get B's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,C (0x69)",
			opcode: 0x69,
			setup: func(cpu *CPU) {
				cpu.C = 0x44
				cpu.L = 0x55
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x44), cpu.L, "L should get C's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,D (0x6A)",
			opcode: 0x6A,
			setup: func(cpu *CPU) {
				cpu.D = 0x33
				cpu.L = 0x66
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x33), cpu.L, "L should get D's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,E (0x6B)",
			opcode: 0x6B,
			setup: func(cpu *CPU) {
				cpu.E = 0x22
				cpu.L = 0x77
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x22), cpu.L, "L should get E's value via dispatch")
			},
		},
		{
			name:   "Dispatch LD L,H (0x6C)",
			opcode: 0x6C,
			setup: func(cpu *CPU) {
				cpu.H = 0x11
				cpu.L = 0x88
			},
			verify: func(t *testing.T, cpu *CPU) {
				assert.Equal(t, uint8(0x11), cpu.L, "L should get H's value via dispatch")
				assert.Equal(t, uint16(0x1111), cpu.GetHL(), "HL should be 0x1111")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := memory.NewMMU()
			
			tt.setup(cpu)
			
			// Execute instruction through dispatch system
			cycles, err := cpu.ExecuteInstruction(mmu, tt.opcode)
			
			assert.NoError(t, err, "Dispatch should not return an error")
			assert.Equal(t, uint8(4), cycles, "All L register ops should take 4 cycles")
			
			tt.verify(t, cpu)
		})
	}
}

// TestL_OpcodeImplementationStatus verifies that all L register opcodes are now implemented
func TestL_OpcodeImplementationStatus(t *testing.T) {
	expectedOpcodes := []uint8{
		0x45, // LD B,L
		0x4D, // LD C,L
		0x68, // LD L,B
		0x69, // LD L,C
		0x6A, // LD L,D
		0x6B, // LD L,E
		0x6C, // LD L,H
		0x6F, // LD L,A
		0x7D, // LD A,L
	}

	for _, opcode := range expectedOpcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X should be implemented", opcode), func(t *testing.T) {
			implemented := IsOpcodeImplemented(opcode)
			assert.True(t, implemented, "Opcode 0x%02X should be implemented in dispatch table", opcode)
		})
	}
	
	// Test that they actually work
	cpu := NewCPU()
	mmu := memory.NewMMU()
	
	for _, opcode := range expectedOpcodes {
		t.Run(fmt.Sprintf("Opcode 0x%02X should execute without error", opcode), func(t *testing.T) {
			cpu.Reset() // Reset to known state
			
			cycles, err := cpu.ExecuteInstruction(mmu, opcode)
			
			assert.NoError(t, err, "Opcode 0x%02X should execute without error", opcode)
			assert.Equal(t, uint8(4), cycles, "Opcode 0x%02X should take 4 cycles", opcode)
		})
	}
}

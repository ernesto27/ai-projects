package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLD_A_L tests the LD A,L instruction (0x7D)
// This instruction copies the value from register L to register A
// It's commonly used to get the low byte of an address into the accumulator
func TestLD_A_L(t *testing.T) {
	tests := []struct {
		name     string
		initialL uint8
		initialA uint8
		wantA    uint8
		wantL    uint8 // L should remain unchanged
	}{
		{
			name:     "Copy 0x00 from L to A",
			initialL: 0x00,
			initialA: 0xFF, // Different value to ensure copy happens
			wantA:    0x00,
			wantL:    0x00,
		},
		{
			name:     "Copy 0xFF from L to A",
			initialL: 0xFF,
			initialA: 0x00, // Different value to ensure copy happens
			wantA:    0xFF,
			wantL:    0xFF,
		},
		{
			name:     "Copy 0x55 from L to A",
			initialL: 0x55,
			initialA: 0xAA, // Opposite pattern to ensure copy happens
			wantA:    0x55,
			wantL:    0x55,
		},
		{
			name:     "Copy 0xAA from L to A",
			initialL: 0xAA,
			initialA: 0x55, // Opposite pattern to ensure copy happens
			wantA:    0xAA,
			wantL:    0xAA,
		},
		{
			name:     "Copy same value (no change test)",
			initialL: 0x42,
			initialA: 0x42, // Same value to test it still works
			wantA:    0x42,
			wantL:    0x42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			// Set initial state
			cpu.L = tt.initialL
			cpu.A = tt.initialA
			
			// Preserve initial flags for testing
			initialFlags := cpu.F
			
			// Execute instruction
			cycles := cpu.LD_A_L()
			
			// Verify results
			assert.Equal(t, tt.wantA, cpu.A, "Register A should contain the value from L")
			assert.Equal(t, tt.wantL, cpu.L, "Register L should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD A,L should take 4 cycles")
			
			// Verify flags are not affected (register loads don't change flags)
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected by register load")
			
			// Verify other registers are not affected
			assert.Equal(t, uint8(0x00), cpu.B, "Register B should not be affected")
			assert.Equal(t, uint8(0x13), cpu.C, "Register C should not be affected")
			assert.Equal(t, uint8(0x00), cpu.D, "Register D should not be affected")
			assert.Equal(t, uint8(0xD8), cpu.E, "Register E should not be affected")
			assert.Equal(t, uint8(0x01), cpu.H, "Register H should not be affected")
		})
	}
}

// TestLD_B_L tests the LD B,L instruction (0x45)
// This instruction copies the value from register L to register B
// Often used for preserving the low byte of HL while using B for other operations
func TestLD_B_L(t *testing.T) {
	tests := []struct {
		name     string
		initialL uint8
		initialB uint8
		wantB    uint8
		wantL    uint8 // L should remain unchanged
	}{
		{
			name:     "Copy 0x00 from L to B",
			initialL: 0x00,
			initialB: 0xFF,
			wantB:    0x00,
			wantL:    0x00,
		},
		{
			name:     "Copy 0xFF from L to B",
			initialL: 0xFF,
			initialB: 0x00,
			wantB:    0xFF,
			wantL:    0xFF,
		},
		{
			name:     "Copy 0x34 from L to B (typical low byte)",
			initialL: 0x34,
			initialB: 0x12,
			wantB:    0x34,
			wantL:    0x34,
		},
		{
			name:     "Copy 0x80 from L to B (high bit set)",
			initialL: 0x80,
			initialB: 0x7F,
			wantB:    0x80,
			wantL:    0x80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			// Set initial state
			cpu.L = tt.initialL
			cpu.B = tt.initialB
			initialFlags := cpu.F
			
			// Execute instruction
			cycles := cpu.LD_B_L()
			
			// Verify results
			assert.Equal(t, tt.wantB, cpu.B, "Register B should contain the value from L")
			assert.Equal(t, tt.wantL, cpu.L, "Register L should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD B,L should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_C_L tests the LD C,L instruction (0x4D)
// This instruction copies the value from register L to register C
// Often used when you need the low byte in register C for I/O operations
func TestLD_C_L(t *testing.T) {
	tests := []struct {
		name     string
		initialL uint8
		initialC uint8
		wantC    uint8
		wantL    uint8
	}{
		{
			name:     "Copy I/O port address low byte",
			initialL: 0x40, // Common I/O port low byte
			initialC: 0x00,
			wantC:    0x40,
			wantL:    0x40,
		},
		{
			name:     "Copy 0xFF from L to C",
			initialL: 0xFF,
			initialC: 0x13, // Default C value from NewCPU
			wantC:    0xFF,
			wantL:    0xFF,
		},
		{
			name:     "Copy zero value",
			initialL: 0x00,
			initialC: 0xAA,
			wantC:    0x00,
			wantL:    0x00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.L = tt.initialL
			cpu.C = tt.initialC
			initialFlags := cpu.F
			
			cycles := cpu.LD_C_L()
			
			assert.Equal(t, tt.wantC, cpu.C, "Register C should contain the value from L")
			assert.Equal(t, tt.wantL, cpu.L, "Register L should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD C,L should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_A tests the LD L,A instruction (0x6F)
// This instruction copies the value from register A to register L
// This is very common - setting the low byte of HL from a calculated value in A
func TestLD_L_A(t *testing.T) {
	tests := []struct {
		name     string
		initialA uint8
		initialL uint8
		wantL    uint8
		wantA    uint8 // A should remain unchanged
	}{
		{
			name:     "Set low byte of address from A",
			initialA: 0x34, // Calculated low byte
			initialL: 0x00,
			wantL:    0x34,
			wantA:    0x34,
		},
		{
			name:     "Copy 0x00 from A to L",
			initialA: 0x00,
			initialL: 0xFF,
			wantL:    0x00,
			wantA:    0x00,
		},
		{
			name:     "Copy 0xFF from A to L",
			initialA: 0xFF,
			initialL: 0x4D, // Default L value from NewCPU
			wantL:    0xFF,
			wantA:    0xFF,
		},
		{
			name:     "Copy calculated result to low byte",
			initialA: 0x80, // Result of some calculation
			initialL: 0x12,
			wantL:    0x80,
			wantA:    0x80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.A = tt.initialA
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_A()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from A")
			assert.Equal(t, tt.wantA, cpu.A, "Register A should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD L,A should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_B tests the LD L,B instruction (0x68)
// This instruction copies the value from register B to register L
// Used when constructing addresses or moving data between register pairs
func TestLD_L_B(t *testing.T) {
	tests := []struct {
		name     string
		initialB uint8
		initialL uint8
		wantL    uint8
		wantB    uint8
	}{
		{
			name:     "Transfer from B to L for address construction",
			initialB: 0x56,
			initialL: 0x00,
			wantL:    0x56,
			wantB:    0x56,
		},
		{
			name:     "Copy zero from B to L",
			initialB: 0x00,
			initialL: 0xFF,
			wantL:    0x00,
			wantB:    0x00,
		},
		{
			name:     "Copy 0xFF from B to L",
			initialB: 0xFF,
			initialL: 0x4D,
			wantL:    0xFF,
			wantB:    0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.B = tt.initialB
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_B()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from B")
			assert.Equal(t, tt.wantB, cpu.B, "Register B should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD L,B should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_C tests the LD L,C instruction (0x69)
// This instruction copies the value from register C to register L
// Common when using C for I/O and then transferring result to address calculation
func TestLD_L_C(t *testing.T) {
	tests := []struct {
		name     string
		initialC uint8
		initialL uint8
		wantL    uint8
		wantC    uint8
	}{
		{
			name:     "Transfer I/O result from C to L",
			initialC: 0x20, // I/O port result
			initialL: 0x00,
			wantL:    0x20,
			wantC:    0x20,
		},
		{
			name:     "Copy port address low byte",
			initialC: 0x44,
			initialL: 0x12,
			wantL:    0x44,
			wantC:    0x44,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.C = tt.initialC
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_C()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from C")
			assert.Equal(t, tt.wantC, cpu.C, "Register C should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD L,C should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_D tests the LD L,D instruction (0x6A)
// This instruction copies the value from register D to register L
// Used in address manipulation when combining DE and HL register pairs
func TestLD_L_D(t *testing.T) {
	tests := []struct {
		name     string
		initialD uint8
		initialL uint8
		wantL    uint8
		wantD    uint8
	}{
		{
			name:     "Transfer high byte of DE to low byte of HL",
			initialD: 0x12, // High byte of some address
			initialL: 0x34,
			wantL:    0x12,
			wantD:    0x12,
		},
		{
			name:     "Copy zero from D to L",
			initialD: 0x00,
			initialL: 0xFF,
			wantL:    0x00,
			wantD:    0x00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.D = tt.initialD
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_D()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from D")
			assert.Equal(t, tt.wantD, cpu.D, "Register D should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD L,D should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_E tests the LD L,E instruction (0x6B)
// This instruction copies the value from register E to register L
// Often used to transfer low bytes between DE and HL register pairs
func TestLD_L_E(t *testing.T) {
	tests := []struct {
		name     string
		initialE uint8
		initialL uint8
		wantL    uint8
		wantE    uint8
	}{
		{
			name:     "Transfer low byte from DE to HL",
			initialE: 0x78, // Low byte of DE
			initialL: 0x56,
			wantL:    0x78,
			wantE:    0x78,
		},
		{
			name:     "Copy byte for address calculation",
			initialE: 0xCD,
			initialL: 0xAB,
			wantL:    0xCD,
			wantE:    0xCD,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.E = tt.initialE
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_E()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from E")
			assert.Equal(t, tt.wantE, cpu.E, "Register E should remain unchanged")
			assert.Equal(t, uint8(4), cycles, "LD L,E should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestLD_L_H tests the LD L,H instruction (0x6C)
// This instruction copies the value from register H to register L
// This swaps high and low bytes within the HL register pair
func TestLD_L_H(t *testing.T) {
	tests := []struct {
		name     string
		initialH uint8
		initialL uint8
		wantL    uint8
		wantH    uint8
		wantHL   uint16 // Final HL value for verification
	}{
		{
			name:     "Duplicate high byte to low byte",
			initialH: 0x12,
			initialL: 0x34,
			wantL:    0x12, // L gets H's value
			wantH:    0x12, // H remains the same
			wantHL:   0x1212, // HL becomes 0x1212
		},
		{
			name:     "Copy 0xFF from H to L",
			initialH: 0xFF,
			initialL: 0x00,
			wantL:    0xFF,
			wantH:    0xFF,
			wantHL:   0xFFFF,
		},
		{
			name:     "Copy 0x00 from H to L",
			initialH: 0x00,
			initialL: 0xAA,
			wantL:    0x00,
			wantH:    0x00,
			wantHL:   0x0000,
		},
		{
			name:     "Pattern replication test",
			initialH: 0xA5,
			initialL: 0x5A,
			wantL:    0xA5,
			wantH:    0xA5,
			wantHL:   0xA5A5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			
			cpu.H = tt.initialH
			cpu.L = tt.initialL
			initialFlags := cpu.F
			
			cycles := cpu.LD_L_H()
			
			assert.Equal(t, tt.wantL, cpu.L, "Register L should contain the value from H")
			assert.Equal(t, tt.wantH, cpu.H, "Register H should remain unchanged")
			assert.Equal(t, tt.wantHL, cpu.GetHL(), "HL register pair should have expected combined value")
			assert.Equal(t, uint8(4), cycles, "LD L,H should take 4 cycles")
			assert.Equal(t, initialFlags, cpu.F, "Flags should not be affected")
		})
	}
}

// TestL_RegisterOperations_Integration tests multiple L register operations together
// This verifies that the L register operations work correctly in combination
func TestL_RegisterOperations_Integration(t *testing.T) {
	t.Run("Sequential L register operations", func(t *testing.T) {
		cpu := NewCPU()
		
		// Test sequence: A -> L -> B -> L -> C -> L -> A
		// This tests a common pattern in Game Boy programming
		
		// Step 1: Load a value into A, then transfer to L
		cpu.A = 0x42
		cycles1 := cpu.LD_L_A() // A(0x42) -> L
		assert.Equal(t, uint8(0x42), cpu.L)
		assert.Equal(t, uint8(4), cycles1)
		
		// Step 2: Copy L to B
		cycles2 := cpu.LD_B_L() // L(0x42) -> B
		assert.Equal(t, uint8(0x42), cpu.B)
		assert.Equal(t, uint8(4), cycles2)
		
		// Step 3: Change L, then copy to C
		cpu.L = 0x84 // Simulate some address calculation
		cycles3 := cpu.LD_C_L() // L(0x84) -> C
		assert.Equal(t, uint8(0x84), cpu.C)
		assert.Equal(t, uint8(4), cycles3)
		
		// Step 4: Copy L back to A
		cycles4 := cpu.LD_A_L() // L(0x84) -> A
		assert.Equal(t, uint8(0x84), cpu.A)
		assert.Equal(t, uint8(4), cycles4)
		
		// Verify final state
		assert.Equal(t, uint8(0x84), cpu.A, "A should have final L value")
		assert.Equal(t, uint8(0x42), cpu.B, "B should preserve first L value")
		assert.Equal(t, uint8(0x84), cpu.C, "C should have final L value")
		assert.Equal(t, uint8(0x84), cpu.L, "L should have final value")
	})
	
	t.Run("HL register pair manipulation", func(t *testing.T) {
		cpu := NewCPU()
		
		// Set up HL register pair
		cpu.H = 0x80
		cpu.L = 0x00
		assert.Equal(t, uint16(0x8000), cpu.GetHL(), "Initial HL should be 0x8000")
		
		// Test LD L,H (duplicate high byte to low byte)
		cycles := cpu.LD_L_H()
		assert.Equal(t, uint8(0x80), cpu.L, "L should now equal H")
		assert.Equal(t, uint8(0x80), cpu.H, "H should remain unchanged")
		assert.Equal(t, uint16(0x8080), cpu.GetHL(), "HL should now be 0x8080")
		assert.Equal(t, uint8(4), cycles)
		
		// Test copying L to other registers
		cycles1 := cpu.LD_A_L()
		cycles2 := cpu.LD_B_L()
		cycles3 := cpu.LD_C_L()
		
		assert.Equal(t, uint8(0x80), cpu.A, "A should have L value")
		assert.Equal(t, uint8(0x80), cpu.B, "B should have L value")
		assert.Equal(t, uint8(0x80), cpu.C, "C should have L value")
		assert.Equal(t, uint8(4), cycles1)
		assert.Equal(t, uint8(4), cycles2)
		assert.Equal(t, uint8(4), cycles3)
	})
}

// TestL_RegisterOperations_EdgeCases tests edge cases and boundary conditions
func TestL_RegisterOperations_EdgeCases(t *testing.T) {
	t.Run("Boundary values", func(t *testing.T) {
		cpu := NewCPU()
		
		// Test with 0x00
		cpu.A = 0x00
		cpu.LD_L_A()
		assert.Equal(t, uint8(0x00), cpu.L)
		
		cpu.LD_B_L()
		assert.Equal(t, uint8(0x00), cpu.B)
		
		// Test with 0xFF
		cpu.A = 0xFF
		cpu.LD_L_A()
		assert.Equal(t, uint8(0xFF), cpu.L)
		
		cpu.LD_C_L()
		assert.Equal(t, uint8(0xFF), cpu.C)
		
		// Test with powers of 2
		for i := uint8(0); i < 8; i++ {
			value := uint8(1 << i)
			cpu.A = value
			cpu.LD_L_A()
			assert.Equal(t, value, cpu.L, "Power of 2 value should be copied correctly")
			
			cpu.LD_A_L()
			assert.Equal(t, value, cpu.A, "Power of 2 value should be copied back correctly")
		}
	})
	
	t.Run("Flag preservation", func(t *testing.T) {
		cpu := NewCPU()
		
		// Set all flags
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)
		initialFlags := cpu.F
		
		// Perform various L register operations
		cpu.A = 0x55
		cpu.LD_L_A()
		assert.Equal(t, initialFlags, cpu.F, "LD L,A should preserve flags")
		
		cpu.LD_A_L()
		assert.Equal(t, initialFlags, cpu.F, "LD A,L should preserve flags")
		
		cpu.LD_B_L()
		assert.Equal(t, initialFlags, cpu.F, "LD B,L should preserve flags")
		
		cpu.LD_C_L()
		assert.Equal(t, initialFlags, cpu.F, "LD C,L should preserve flags")
		
		cpu.B = 0xAA
		cpu.LD_L_B()
		assert.Equal(t, initialFlags, cpu.F, "LD L,B should preserve flags")
		
		cpu.C = 0x33
		cpu.LD_L_C()
		assert.Equal(t, initialFlags, cpu.F, "LD L,C should preserve flags")
		
		cpu.D = 0xCC
		cpu.LD_L_D()
		assert.Equal(t, initialFlags, cpu.F, "LD L,D should preserve flags")
		
		cpu.E = 0x99
		cpu.LD_L_E()
		assert.Equal(t, initialFlags, cpu.F, "LD L,E should preserve flags")
		
		cpu.H = 0x66
		cpu.LD_L_H()
		assert.Equal(t, initialFlags, cpu.F, "LD L,H should preserve flags")
	})
}

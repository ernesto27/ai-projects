package cpu

import (
	"gameboy-emulator/internal/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

// === Tests for Memory/MMU Wrapper Functions ===
// These test wrapper functions that need MMU access but no parameter extraction

// === Test Memory Load Wrapper Functions ===

func TestWrapLD_A_HL(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up HL to point to a memory location
	cpu.SetHL(0x8000)
	mmu.WriteByte(0x8000, 0x42)

	// Test the wrapper
	cycles, err := wrapLD_A_HL(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD A,(HL) should return 8 cycles")
	assert.Equal(t, uint8(0x42), cpu.A, "A should contain value from memory")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL should remain unchanged")
}

func TestWrapLD_A_BC(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up BC to point to a memory location
	cpu.SetBC(0x9000)
	mmu.WriteByte(0x9000, 0x55)

	// Test the wrapper
	cycles, err := wrapLD_A_BC(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD A,(BC) should return 8 cycles")
	assert.Equal(t, uint8(0x55), cpu.A, "A should contain value from memory")
	assert.Equal(t, uint16(0x9000), cpu.GetBC(), "BC should remain unchanged")
}

func TestWrapLD_A_DE(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up DE to point to a memory location
	cpu.SetDE(0xA000)
	mmu.WriteByte(0xA000, 0xAA)

	// Test the wrapper
	cycles, err := wrapLD_A_DE(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD A,(DE) should return 8 cycles")
	assert.Equal(t, uint8(0xAA), cpu.A, "A should contain value from memory")
	assert.Equal(t, uint16(0xA000), cpu.GetDE(), "DE should remain unchanged")
}

// === Test Memory Store Wrapper Functions ===

func TestWrapLD_HL_A(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up A with a value and HL pointing to memory
	cpu.A = 0x99
	cpu.SetHL(0x8000)

	// Test the wrapper
	cycles, err := wrapLD_HL_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD (HL),A should return 8 cycles")
	assert.Equal(t, uint8(0x99), mmu.ReadByte(0x8000), "Memory should contain A's value")
	assert.Equal(t, uint8(0x99), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL should remain unchanged")
}

func TestWrapLD_BC_A(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up A with a value and BC pointing to memory
	cpu.A = 0x77
	cpu.SetBC(0x9000)

	// Test the wrapper
	cycles, err := wrapLD_BC_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD (BC),A should return 8 cycles")
	assert.Equal(t, uint8(0x77), mmu.ReadByte(0x9000), "Memory should contain A's value")
	assert.Equal(t, uint8(0x77), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint16(0x9000), cpu.GetBC(), "BC should remain unchanged")
}

func TestWrapLD_DE_A(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	// Set up A with a value and DE pointing to memory
	cpu.A = 0x33
	cpu.SetDE(0xA000)

	// Test the wrapper
	cycles, err := wrapLD_DE_A(cpu, mmu)

	assert.NoError(t, err)
	assert.Equal(t, uint8(8), cycles, "LD (DE),A should return 8 cycles")
	assert.Equal(t, uint8(0x33), mmu.ReadByte(0xA000), "Memory should contain A's value")
	assert.Equal(t, uint8(0x33), cpu.A, "A should remain unchanged")
	assert.Equal(t, uint16(0xA000), cpu.GetDE(), "DE should remain unchanged")
}

// === Test Memory Operations with Different Memory Regions ===

func TestMemoryRegionOperations(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	memoryTests := []struct {
		name    string
		address uint16
		value   uint8
		region  string
	}{
		{"ROM Bank 0", 0x0100, 0x10, "ROM"},
		{"ROM Bank 1", 0x4000, 0x40, "ROM"},
		{"VRAM", 0x8000, 0x80, "VRAM"},
		{"External RAM", 0xA000, 0xA0, "External RAM"},
		{"Work RAM", 0xC000, 0xC0, "Work RAM"},
		{"High RAM", 0xFF80, 0xFF, "High RAM"},
	}

	for _, tt := range memoryTests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset CPU
			cpu.Reset()

			// Test load operations
			cpu.SetHL(tt.address)
			mmu.WriteByte(tt.address, tt.value)

			cycles, err := wrapLD_A_HL(cpu, mmu)
			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.value, cpu.A, "Should load correct value from %s", tt.region)

			// Test store operations
			cpu.A = tt.value + 1 // Different value
			cycles, err = wrapLD_HL_A(cpu, mmu)
			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.value+1, mmu.ReadByte(tt.address), "Should store correct value to %s", tt.region)
		})
	}
}

// === Test Load/Store Round-Trip Operations ===

func TestMemoryRoundTrip(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	testCases := []struct {
		name          string
		setRegPair    func(*CPU, uint16)
		wrapLoad      func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		wrapStore     func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		address       uint16
		originalValue uint8
		newValue      uint8
	}{
		{
			name:          "HL round-trip",
			setRegPair:    (*CPU).SetHL,
			wrapLoad:      wrapLD_A_HL,
			wrapStore:     wrapLD_HL_A,
			address:       0x8000,
			originalValue: 0x42,
			newValue:      0x84,
		},
		{
			name:          "BC round-trip",
			setRegPair:    (*CPU).SetBC,
			wrapLoad:      wrapLD_A_BC,
			wrapStore:     wrapLD_BC_A,
			address:       0x9000,
			originalValue: 0x55,
			newValue:      0xAA,
		},
		{
			name:          "DE round-trip",
			setRegPair:    (*CPU).SetDE,
			wrapLoad:      wrapLD_A_DE,
			wrapStore:     wrapLD_DE_A,
			address:       0xA000,
			originalValue: 0x33,
			newValue:      0xCC,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Reset CPU
			cpu.Reset()

			// Set up register pair and memory
			tt.setRegPair(cpu, tt.address)
			mmu.WriteByte(tt.address, tt.originalValue)

			// Load from memory to A
			cycles, err := tt.wrapLoad(cpu, mmu)
			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.originalValue, cpu.A, "Should load original value")

			// Change A and store back to memory
			cpu.A = tt.newValue
			cycles, err = tt.wrapStore(cpu, mmu)
			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.newValue, mmu.ReadByte(tt.address), "Should store new value")

			// Load again to verify
			cycles, err = tt.wrapLoad(cpu, mmu)
			assert.NoError(t, err)
			assert.Equal(t, uint8(8), cycles)
			assert.Equal(t, tt.newValue, cpu.A, "Should load the stored value")
		})
	}
}

// === Comparison Test: Memory Wrappers vs Original Functions ===

func TestMemoryWrappersVsOriginals(t *testing.T) {
	tests := []struct {
		name     string
		wrapper  func(*CPU, memory.MemoryInterface, ...uint8) (uint8, error)
		original func(*CPU, memory.MemoryInterface) uint8
		setup    func(*CPU, memory.MemoryInterface)
		check    func(*CPU, *CPU, memory.MemoryInterface) bool
	}{
		{
			name:     "LD A,(HL) comparison",
			wrapper:  wrapLD_A_HL,
			original: (*CPU).LD_A_HL,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.SetHL(0x8000)
				mmu.WriteByte(0x8000, 0x42)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A
			},
		},
		{
			name:     "LD (HL),A comparison",
			wrapper:  wrapLD_HL_A,
			original: (*CPU).LD_HL_A,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.A = 0x99
				cpu.SetHL(0x8000)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A && cpu1.GetHL() == cpu2.GetHL()
			},
		},
		{
			name:     "LD A,(BC) comparison",
			wrapper:  wrapLD_A_BC,
			original: (*CPU).LD_A_BC,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.SetBC(0x9000)
				mmu.WriteByte(0x9000, 0x55)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A
			},
		},
		{
			name:     "LD (BC),A comparison",
			wrapper:  wrapLD_BC_A,
			original: (*CPU).LD_BC_A,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.A = 0x77
				cpu.SetBC(0x9000)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A && cpu1.GetBC() == cpu2.GetBC()
			},
		},
		{
			name:     "LD A,(DE) comparison",
			wrapper:  wrapLD_A_DE,
			original: (*CPU).LD_A_DE,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.SetDE(0xA000)
				mmu.WriteByte(0xA000, 0xAA)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A
			},
		},
		{
			name:     "LD (DE),A comparison",
			wrapper:  wrapLD_DE_A,
			original: (*CPU).LD_DE_A,
			setup: func(cpu *CPU, mmu memory.MemoryInterface) {
				cpu.A = 0x33
				cpu.SetDE(0xA000)
			},
			check: func(cpu1, cpu2 *CPU, mmu memory.MemoryInterface) bool {
				return cpu1.A == cpu2.A && cpu1.GetDE() == cpu2.GetDE()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create two identical CPUs and MMUs
			cpu1 := NewCPU()
			cpu2 := NewCPU()
			mmu1 := memory.NewMMU()
			mmu2 := memory.NewMMU()

			// Set them up identically
			tt.setup(cpu1, mmu1)
			tt.setup(cpu2, mmu2)

			// Call original on cpu1, wrapper on cpu2
			originalCycles := tt.original(cpu1, mmu1)
			wrapperCycles, err := tt.wrapper(cpu2, mmu2)

			// They should behave identically
			assert.NoError(t, err)
			assert.Equal(t, originalCycles, wrapperCycles, "cycles should match")
			assert.True(t, tt.check(cpu1, cpu2, mmu1), "CPU state should match")

			// Memory should also be identical
			for addr := uint16(0x8000); addr <= 0xA000; addr += 0x1000 {
				assert.Equal(t, mmu1.ReadByte(addr), mmu2.ReadByte(addr), "Memory at 0x%04X should match", addr)
			}
		})
	}
}

// === Test Edge Cases and Error Conditions ===

func TestMemoryEdgeCases(t *testing.T) {
	cpu := NewCPU()
	mmu := memory.NewMMU()

	t.Run("Memory boundaries", func(t *testing.T) {
		// Test at memory boundaries
		boundaries := []uint16{0x0000, 0x7FFF, 0x8000, 0x9FFF, 0xA000, 0xBFFF, 0xC000, 0xDFFF, 0xFE00, 0xFEFF, 0xFF00, 0xFF7F, 0xFF80, 0xFFFE}

		for _, addr := range boundaries {
			cpu.Reset()
			cpu.SetHL(addr)
			cpu.A = 0x42

			// Test store operation
			cycles, err := wrapLD_HL_A(cpu, mmu)
			assert.NoError(t, err, "Store should work at boundary 0x%04X", addr)
			assert.Equal(t, uint8(8), cycles)

			// Test load operation
			cycles, err = wrapLD_A_HL(cpu, mmu)
			assert.NoError(t, err, "Load should work at boundary 0x%04X", addr)
			assert.Equal(t, uint8(8), cycles)
		}
	})

	t.Run("Value preservation", func(t *testing.T) {
		// Test that all other registers are preserved
		cpu.Reset()
		cpu.A = 0x11
		cpu.B = 0x22
		cpu.C = 0x33
		cpu.D = 0x44
		cpu.E = 0x55
		cpu.H = 0x80 // Setting H to 0x80
		cpu.L = 0x00 // Setting L to 0x00 (so HL = 0x8000)
		cpu.F = 0x80

		originalB, originalC, originalD, originalE, originalH, originalL, originalF := cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L, cpu.F

		// HL is already 0x8000, so no need to call SetHL
		mmu.WriteByte(0x8000, 0x99)

		// Load should only affect A
		_, err := wrapLD_A_HL(cpu, mmu)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x99), cpu.A, "A should be updated")
		assert.Equal(t, originalB, cpu.B, "B should be preserved")
		assert.Equal(t, originalC, cpu.C, "C should be preserved")
		assert.Equal(t, originalD, cpu.D, "D should be preserved")
		assert.Equal(t, originalE, cpu.E, "E should be preserved")
		assert.Equal(t, originalH, cpu.H, "H should be preserved")
		assert.Equal(t, originalL, cpu.L, "L should be preserved")
		assert.Equal(t, originalF, cpu.F, "F should be preserved")
	})
}

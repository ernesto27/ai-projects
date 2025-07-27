package cpu

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/memory"

	"github.com/stretchr/testify/assert"
)

// createTestMMU creates an MMU with a dummy MBC for testing
func createTestMMU() *memory.MMU {
	// Create a simple ROM-only cartridge for testing
	romData := make([]byte, 32*1024)
	
	// Add minimal header
	copy(romData[0x0134:], "TEST")
	romData[0x0147] = uint8(cartridge.ROM_ONLY)
	romData[0x0148] = 0x00 // 32KB
	romData[0x0149] = 0x00 // No RAM
	
	// Calculate checksum
	var checksum uint8 = 0
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum = checksum - romData[addr] - 1
	}
	romData[0x014D] = checksum
	
	cart, _ := cartridge.NewCartridge(romData)
	mbc, _ := cartridge.CreateMBC(cart)
	return memory.NewMMU(mbc)
}

// TestLD_D_n tests the LD D,n instruction
func TestLD_D_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into D
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialE := cpu.E
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD D,n instruction
		cycles := cpu.LD_D_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD D,n should take 8 cycles")

		// D register should contain the loaded value
		assert.Equal(t, value, cpu.D, "D register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_E_n tests the LD E,n instruction
func TestLD_E_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into E
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80, 0x55, 0xAA}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialH := cpu.H
		initialL := cpu.L
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD E,n instruction
		cycles := cpu.LD_E_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD E,n should take 8 cycles")

		// E register should contain the loaded value
		assert.Equal(t, value, cpu.E, "E register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialH, cpu.H, "H register should be unchanged")
		assert.Equal(t, initialL, cpu.L, "L register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// TestLD_L_n tests the LD L,n instruction
func TestLD_L_n(t *testing.T) {
	cpu := NewCPU()

	// Test loading different values into L
	testValues := []uint8{0x00, 0x42, 0xFF, 0x01, 0x80}

	for _, value := range testValues {
		// Store initial state (other registers should be unchanged)
		initialA := cpu.A
		initialB := cpu.B
		initialC := cpu.C
		initialD := cpu.D
		initialE := cpu.E
		initialF := cpu.F
		initialSP := cpu.SP
		initialPC := cpu.PC

		// Execute LD L,n instruction
		cycles := cpu.LD_L_n(value)

		// Should take 8 cycles
		assert.Equal(t, uint8(8), cycles, "LD L,n should take 8 cycles")

		// L register should contain the loaded value
		assert.Equal(t, value, cpu.L, "L register should contain the loaded value")

		// Other registers should be unchanged
		assert.Equal(t, initialA, cpu.A, "A register should be unchanged")
		assert.Equal(t, initialB, cpu.B, "B register should be unchanged")
		assert.Equal(t, initialC, cpu.C, "C register should be unchanged")
		assert.Equal(t, initialD, cpu.D, "D register should be unchanged")
		assert.Equal(t, initialE, cpu.E, "E register should be unchanged")
		assert.Equal(t, initialF, cpu.F, "F register should be unchanged")
		assert.Equal(t, initialSP, cpu.SP, "SP should be unchanged")
		assert.Equal(t, initialPC, cpu.PC, "PC should be unchanged")
	}
}

// === Memory Load Instruction Tests ===

func TestLD_A_HL(t *testing.T) {
	// Test LD_A_HL - Load A from memory at HL (0x7E)

	t.Run("Basic memory read", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Set up test scenario
		cpu.SetHL(0x8000)           // HL points to VRAM start
		mmu.WriteByte(0x8000, 0x42) // Write test value to memory
		cpu.A = 0x00                // Clear A register

		// Execute instruction
		cycles := cpu.LD_A_HL(mmu)

		// Verify results
		assert.Equal(t, uint8(0x42), cpu.A, "A should contain value read from memory")
		assert.Equal(t, uint8(8), cycles, "LD_A_HL should take 8 cycles")
		assert.Equal(t, uint16(0x8000), cpu.GetHL(), "HL should remain unchanged")
	})

	t.Run("Read from different memory regions", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Test reading from Work RAM (WRAM)
		cpu.SetHL(0xC000)           // HL points to WRAM start
		mmu.WriteByte(0xC000, 0x99) // Write test value
		cpu.A = 0x00

		cycles := cpu.LD_A_HL(mmu)

		assert.Equal(t, uint8(0x99), cpu.A, "Should read from WRAM")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	})

	t.Run("Read zero value", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Memory is initialized to zero by default
		cpu.SetHL(0xD000) // HL points to middle of WRAM
		cpu.A = 0xFF      // Set A to non-zero value

		cycles := cpu.LD_A_HL(mmu)

		assert.Equal(t, uint8(0x00), cpu.A, "Should read zero from uninitialized memory")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	})

	t.Run("Read maximum value", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		cpu.SetHL(0xE000)           // HL points to Echo RAM
		mmu.WriteByte(0xE000, 0xFF) // Write maximum 8-bit value
		cpu.A = 0x00

		cycles := cpu.LD_A_HL(mmu)

		assert.Equal(t, uint8(0xFF), cpu.A, "Should read maximum value 0xFF")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	})

	t.Run("Flags are not affected", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Set all flags to specific values
		cpu.SetFlag(FlagZ, true)
		cpu.SetFlag(FlagN, true)
		cpu.SetFlag(FlagH, true)
		cpu.SetFlag(FlagC, true)

		cpu.SetHL(0xC100)
		mmu.WriteByte(0xC100, 0x55)

		cpu.LD_A_HL(mmu)

		// Verify flags remain unchanged
		assert.True(t, cpu.GetFlag(FlagZ), "Z flag should remain set")
		assert.True(t, cpu.GetFlag(FlagN), "N flag should remain set")
		assert.True(t, cpu.GetFlag(FlagH), "H flag should remain set")
		assert.True(t, cpu.GetFlag(FlagC), "C flag should remain set")
	})

	t.Run("Other registers preserved", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		// Set all registers to known values
		cpu.B = 0x11
		cpu.C = 0x22
		cpu.D = 0x33
		cpu.E = 0x44
		cpu.H = 0x80 // HL = 0x8055
		cpu.L = 0x55
		cpu.SP = 0xFFFE
		cpu.PC = 0x0150

		mmu.WriteByte(0x8055, 0x99)

		cpu.LD_A_HL(mmu)

		// Verify other registers unchanged
		assert.Equal(t, uint8(0x11), cpu.B, "B register should be preserved")
		assert.Equal(t, uint8(0x22), cpu.C, "C register should be preserved")
		assert.Equal(t, uint8(0x33), cpu.D, "D register should be preserved")
		assert.Equal(t, uint8(0x44), cpu.E, "E register should be preserved")
		assert.Equal(t, uint8(0x80), cpu.H, "H register should be preserved")
		assert.Equal(t, uint8(0x55), cpu.L, "L register should be preserved")
		assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be preserved")
		assert.Equal(t, uint16(0x0150), cpu.PC, "PC should be preserved")

		// Verify A changed to expected value
		assert.Equal(t, uint8(0x99), cpu.A, "A should contain value from memory")
	})

	t.Run("Edge case - read from VRAM start", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		cpu.SetHL(0x8000)           // HL points to start of VRAM
		mmu.WriteByte(0x8000, 0x31) // Write value to VRAM area
		cpu.A = 0x00

		cycles := cpu.LD_A_HL(mmu)

		assert.Equal(t, uint8(0x31), cpu.A, "Should read from VRAM area")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	})

	t.Run("Edge case - read from address 0xFFFF", func(t *testing.T) {
		cpu := NewCPU()
		mmu := createTestMMU()

		cpu.SetHL(0xFFFF)           // HL points to end of address space
		mmu.WriteByte(0xFFFF, 0x88) // Write value to high RAM
		cpu.A = 0x00

		cycles := cpu.LD_A_HL(mmu)

		assert.Equal(t, uint8(0x88), cpu.A, "Should read from high RAM")
		assert.Equal(t, uint8(8), cycles, "Should take 8 cycles")
	})
}

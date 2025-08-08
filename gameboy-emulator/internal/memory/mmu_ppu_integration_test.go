package memory

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/joypad"
	"gameboy-emulator/internal/ppu"

	"github.com/stretchr/testify/assert"
)

// TestPPUMMUIntegration tests complete PPU-MMU integration
func TestPPUMMUIntegration(t *testing.T) {
	// Create test components
	dummyMBC := &cartridge.MBC0{}
	interruptController := interrupt.NewInterruptController()
	mmu := NewMMU(dummyMBC, interruptController, joypad.NewJoypad())
	ppuInstance := ppu.NewPPU()
	
	// Connect PPU to MMU
	mmu.SetPPU(ppuInstance)

	t.Run("PPU registers accessible through MMU", func(t *testing.T) {
		// Test LCDC register (0xFF40)
		mmu.WriteByte(0xFF40, 0x91)
		assert.Equal(t, uint8(0x91), mmu.ReadByte(0xFF40), "LCDC register should be accessible through MMU")
		
		// Test STAT register (0xFF41)
		mmu.WriteByte(0xFF41, 0x80)
		result := mmu.ReadByte(0xFF41)
		// STAT register has read-only bits, so we check that write had some effect
		assert.NotEqual(t, uint8(0x00), result, "STAT register should be accessible through MMU")
		
		// Test scroll registers
		mmu.WriteByte(0xFF42, 0x10) // SCY
		mmu.WriteByte(0xFF43, 0x20) // SCX
		assert.Equal(t, uint8(0x10), mmu.ReadByte(0xFF42), "SCY register should work")
		assert.Equal(t, uint8(0x20), mmu.ReadByte(0xFF43), "SCX register should work")
		
		// Test LY register is read-only
		mmu.WriteByte(0xFF44, 0x50)
		ly := mmu.ReadByte(0xFF44)
		assert.NotEqual(t, uint8(0x50), ly, "LY register should be read-only")
		
		// Test LYC register
		mmu.WriteByte(0xFF45, 0x30)
		assert.Equal(t, uint8(0x30), mmu.ReadByte(0xFF45), "LYC register should work")
		
		// Test palette registers
		mmu.WriteByte(0xFF47, 0xE4) // BGP
		mmu.WriteByte(0xFF48, 0xD2) // OBP0
		mmu.WriteByte(0xFF49, 0xA4) // OBP1
		assert.Equal(t, uint8(0xE4), mmu.ReadByte(0xFF47), "BGP register should work")
		assert.Equal(t, uint8(0xD2), mmu.ReadByte(0xFF48), "OBP0 register should work")
		assert.Equal(t, uint8(0xA4), mmu.ReadByte(0xFF49), "OBP1 register should work")
		
		// Test window registers
		mmu.WriteByte(0xFF4A, 0x40) // WY
		mmu.WriteByte(0xFF4B, 0x50) // WX
		assert.Equal(t, uint8(0x40), mmu.ReadByte(0xFF4A), "WY register should work")
		assert.Equal(t, uint8(0x50), mmu.ReadByte(0xFF4B), "WX register should work")
	})

	t.Run("VRAM access routed to PPU", func(t *testing.T) {
		// Test VRAM write/read through MMU
		testData := []struct {
			address uint16
			value   uint8
		}{
			{0x8000, 0x42}, // Start of VRAM
			{0x8100, 0x24}, // Middle of VRAM
			{0x9FFF, 0x99}, // End of VRAM
		}

		for _, td := range testData {
			mmu.WriteByte(td.address, td.value)
			result := mmu.ReadByte(td.address)
			assert.Equal(t, td.value, result, "VRAM should be accessible at address 0x%04X", td.address)
		}
	})

	t.Run("OAM access routed to PPU", func(t *testing.T) {
		// Advance PPU to a mode where OAM is accessible (H-Blank or V-Blank)
		// PPU starts in OAM Scan mode which blocks OAM access
		for ppuInstance.GetCurrentMode() == ppu.ModeOAMScan || ppuInstance.GetCurrentMode() == ppu.ModeDrawing {
			ppuInstance.Update(1)
		}
		
		// Test OAM write/read through MMU
		testData := []struct {
			address uint16
			value   uint8
		}{
			{0xFE00, 0x80}, // Start of OAM
			{0xFE50, 0x40}, // Middle of OAM
			{0xFE9F, 0x20}, // End of OAM
		}

		for _, td := range testData {
			mmu.WriteByte(td.address, td.value)
			result := mmu.ReadByte(td.address)
			assert.Equal(t, td.value, result, "OAM should be accessible at address 0x%04X", td.address)
		}
	})

	t.Run("VRAM access restrictions during Drawing mode", func(t *testing.T) {
		// Set PPU to Drawing mode (Mode 3) by directly manipulating PPU state
		// We'll use reflection or create a test helper, but for now let's use Update cycles
		ppuInstance.Reset() // Start fresh
		
		// Test that VRAM is normally accessible (not in Drawing mode initially)
		mmu.WriteByte(0x8000, 0xAA)
		result := mmu.ReadByte(0x8000)
		assert.Equal(t, uint8(0xAA), result, "VRAM should be accessible when not in Drawing mode")
		
		// Force PPU into Drawing mode by advancing cycles
		// PPU starts in OAM Scan (80 cycles), then Drawing mode starts
		for i := 0; i < 85; i++ { // Advance past OAM Scan into Drawing
			ppuInstance.Update(1) // 1 cycle per update
		}
		
		// Check that PPU is in Drawing mode
		assert.Equal(t, ppu.ModeDrawing, ppuInstance.GetCurrentMode(), "PPU should be in Drawing mode")
		
		// Now VRAM reads should return 0xFF
		result = mmu.ReadByte(0x8000)
		assert.Equal(t, uint8(0xFF), result, "VRAM read should return 0xFF during Drawing mode")
		
		// VRAM writes should be ignored
		mmu.WriteByte(0x8000, 0xBB)
		// Advance to non-Drawing mode
		for ppuInstance.GetCurrentMode() == ppu.ModeDrawing {
			ppuInstance.Update(1)
		}
		result = mmu.ReadByte(0x8000)
		assert.Equal(t, uint8(0xAA), result, "VRAM write should be ignored during Drawing mode")
	})

	t.Run("OAM access restrictions during Drawing and OAM Scan modes", func(t *testing.T) {
		ppuInstance.Reset()
		
		// Test that OAM is accessible in H-Blank
		ppuInstance.Update(1) // Advance to ensure we're not in restricted mode initially
		for ppuInstance.GetCurrentMode() != ppu.ModeHBlank && ppuInstance.GetCurrentMode() != ppu.ModeVBlank {
			ppuInstance.Update(1)
		}
		
		mmu.WriteByte(0xFE00, 0xCC)
		result := mmu.ReadByte(0xFE00)
		assert.Equal(t, uint8(0xCC), result, "OAM should be accessible in H-Blank/V-Blank")
		
		// Force into OAM Scan mode (Mode 2)
		ppuInstance.Reset()
		assert.Equal(t, ppu.ModeOAMScan, ppuInstance.GetCurrentMode(), "PPU should start in OAM Scan mode")
		
		// OAM reads should return 0xFF during OAM Scan
		result = mmu.ReadByte(0xFE00)
		assert.Equal(t, uint8(0xFF), result, "OAM read should return 0xFF during OAM Scan mode")
		
		// OAM writes should be ignored during OAM Scan
		mmu.WriteByte(0xFE00, 0xDD)
		// Advance to accessible mode
		for ppuInstance.GetCurrentMode() == ppu.ModeOAMScan || ppuInstance.GetCurrentMode() == ppu.ModeDrawing {
			ppuInstance.Update(1)
		}
		result = mmu.ReadByte(0xFE00)
		assert.Equal(t, uint8(0xCC), result, "OAM write should be ignored during OAM Scan mode")
	})

	t.Run("Direct PPU access works correctly", func(t *testing.T) {
		// Test that direct PPU access works alongside MMU routing
		ppuInstance.WriteVRAM(0x8200, 0xEE)
		assert.Equal(t, uint8(0xEE), ppuInstance.ReadVRAM(0x8200), "Direct PPU VRAM access should work")
		
		// Verify MMU sees the same data
		assert.Equal(t, uint8(0xEE), mmu.ReadByte(0x8200), "MMU should see data written directly to PPU")
		
		// Test same for OAM
		ppuInstance.WriteOAM(0xFE10, 0xFF)
		assert.Equal(t, uint8(0xFF), ppuInstance.ReadOAM(0xFE10), "Direct PPU OAM access should work")
		assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xFE10), "MMU should see data written directly to PPU OAM")
	})

	t.Run("Invalid VRAM/OAM addresses handled correctly", func(t *testing.T) {
		// Test invalid VRAM addresses
		result := ppuInstance.ReadVRAM(0x7FFF) // Below VRAM range
		assert.Equal(t, uint8(0xFF), result, "Invalid VRAM address should return 0xFF")
		
		result = ppuInstance.ReadVRAM(0xA000) // Above VRAM range
		assert.Equal(t, uint8(0xFF), result, "Invalid VRAM address should return 0xFF")
		
		// Test invalid OAM addresses
		result = ppuInstance.ReadOAM(0xFDFF) // Below OAM range
		assert.Equal(t, uint8(0xFF), result, "Invalid OAM address should return 0xFF")
		
		result = ppuInstance.ReadOAM(0xFEA0) // Above OAM range
		assert.Equal(t, uint8(0xFF), result, "Invalid OAM address should return 0xFF")
		
		// Invalid writes should be ignored (no crash/panic)
		ppuInstance.WriteVRAM(0x7FFF, 0x42) // Should not crash
		ppuInstance.WriteOAM(0xFEA0, 0x42)  // Should not crash
	})
}

// TestPPUMMUWithoutPPU tests MMU behavior when no PPU is connected
func TestPPUMMUWithoutPPU(t *testing.T) {
	dummyMBC := &cartridge.MBC0{}
	interruptController := interrupt.NewInterruptController()
	mmu := NewMMU(dummyMBC, interruptController, joypad.NewJoypad())
	// Don't call SetPPU - test fallback behavior

	t.Run("VRAM falls back to internal memory without PPU", func(t *testing.T) {
		mmu.WriteByte(0x8000, 0x42)
		result := mmu.ReadByte(0x8000)
		assert.Equal(t, uint8(0x42), result, "VRAM should fall back to internal memory without PPU")
	})

	t.Run("OAM falls back to internal memory without PPU", func(t *testing.T) {
		mmu.WriteByte(0xFE00, 0x84)
		result := mmu.ReadByte(0xFE00)
		assert.Equal(t, uint8(0x84), result, "OAM should fall back to internal memory without PPU")
	})

	t.Run("PPU registers return default values without PPU", func(t *testing.T) {
		// PPU register access should return 0x00 or default internal memory values
		result := mmu.ReadByte(0xFF40) // LCDC
		assert.Equal(t, uint8(0x00), result, "LCDC should return default without PPU")
		
		// Writes should go to internal memory
		mmu.WriteByte(0xFF40, 0x91)
		result = mmu.ReadByte(0xFF40)
		assert.Equal(t, uint8(0x91), result, "LCDC write should go to internal memory without PPU")
	})
}
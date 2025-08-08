package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gameboy-emulator/internal/memory"
	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/joypad"
)

// TestPPUMMUIntegration tests the complete PPU-MMU integration
func TestPPUMMUIntegration(t *testing.T) {
	// Create a mock ROM for testing
	mockROM := make([]byte, 32*1024) // 32KB ROM
	mockROM[0x0147] = 0x00 // ROM only cartridge
	
	// Create cartridge and interrupt controller
	cart, err := cartridge.LoadROMFromBytes(mockROM, "test-rom")
	assert.NoError(t, err, "Should create test cartridge")
	
	mbc, err := cartridge.CreateMBC(cart)
	assert.NoError(t, err, "Should create MBC")
	interruptController := interrupt.NewInterruptController()
	
	// Create MMU and PPU
	mmu := memory.NewMMU(mbc, interruptController, joypad.NewJoypad())
	ppu := NewPPU()
	
	// Connect PPU to MMU
	mmu.SetPPU(ppu)
	
	// Test LCDC register access through MMU
	mmu.WriteByte(0xFF40, 0x11) // Disable LCD (bit 7 = 0)
	assert.Equal(t, uint8(0x11), mmu.ReadByte(0xFF40), "LCDC should be readable through MMU")
	assert.Equal(t, uint8(0x11), ppu.GetLCDC(), "PPU LCDC should be updated")
	assert.False(t, ppu.IsLCDEnabled(), "LCD should be disabled")
	
	// Test STAT register access through MMU
	mmu.WriteByte(0xFF41, 0x78) // Enable all STAT interrupts
	stat := mmu.ReadByte(0xFF41)
	assert.Equal(t, uint8(0x78), stat&0x78, "STAT interrupt enables should be set")
	assert.Equal(t, uint8(0x00), stat&0x03, "Mode should be H-Blank when LCD disabled")
	
	// Test palette registers through MMU
	mmu.WriteByte(0xFF47, 0x1B) // Set BGP to inverted palette
	assert.Equal(t, uint8(0x1B), mmu.ReadByte(0xFF47), "BGP should be readable through MMU")
	assert.Equal(t, uint8(0x1B), ppu.GetBGP(), "PPU BGP should be updated")
	
	mmu.WriteByte(0xFF48, 0x30) // Set OBP0
	assert.Equal(t, uint8(0x30), mmu.ReadByte(0xFF48), "OBP0 should be readable through MMU")
	assert.Equal(t, uint8(0x30), ppu.GetOBP0(), "PPU OBP0 should be updated")
	
	mmu.WriteByte(0xFF49, 0xFC) // Set OBP1
	assert.Equal(t, uint8(0xFC), mmu.ReadByte(0xFF49), "OBP1 should be readable through MMU")
	assert.Equal(t, uint8(0xFC), ppu.GetOBP1(), "PPU OBP1 should be updated")
	
	// Test scroll registers through MMU
	mmu.WriteByte(0xFF42, 123) // Set SCY
	assert.Equal(t, uint8(123), mmu.ReadByte(0xFF42), "SCY should be readable through MMU")
	assert.Equal(t, uint8(123), ppu.GetSCY(), "PPU SCY should be updated")
	
	mmu.WriteByte(0xFF43, 234) // Set SCX
	assert.Equal(t, uint8(234), mmu.ReadByte(0xFF43), "SCX should be readable through MMU")
	assert.Equal(t, uint8(234), ppu.GetSCX(), "PPU SCX should be updated")
	
	// Test LYC register through MMU
	mmu.WriteByte(0xFF45, 100) // Set LYC
	assert.Equal(t, uint8(100), mmu.ReadByte(0xFF45), "LYC should be readable through MMU")
	assert.Equal(t, uint8(100), ppu.GetLYC(), "PPU LYC should be updated")
	
	// Test window registers through MMU
	mmu.WriteByte(0xFF4A, 144) // Set WY
	assert.Equal(t, uint8(144), mmu.ReadByte(0xFF4A), "WY should be readable through MMU")
	assert.Equal(t, uint8(144), ppu.GetWY(), "PPU WY should be updated")
	
	mmu.WriteByte(0xFF4B, 167) // Set WX
	assert.Equal(t, uint8(167), mmu.ReadByte(0xFF4B), "WX should be readable through MMU")
	assert.Equal(t, uint8(167), ppu.GetWX(), "PPU WX should be updated")
	
	// Test LY register is read-only through MMU
	originalLY := mmu.ReadByte(0xFF44)
	mmu.WriteByte(0xFF44, 99) // Try to write to LY (should be ignored)
	assert.Equal(t, originalLY, mmu.ReadByte(0xFF44), "LY should remain unchanged (read-only)")
	assert.Equal(t, originalLY, ppu.GetLY(), "PPU LY should remain unchanged")
}

// TestPPURegisterConstants tests that PPU constants match MMU constants
func TestPPURegisterConstants(t *testing.T) {
	// These constants should match between PPU and MMU packages
	assert.Equal(t, uint16(0xFF40), LCDCAddress, "LCDC address constant should match")
	assert.Equal(t, uint16(0xFF41), STATAddress, "STAT address constant should match")
	assert.Equal(t, uint16(0xFF42), SCYAddress, "SCY address constant should match")
	assert.Equal(t, uint16(0xFF43), SCXAddress, "SCX address constant should match")
	assert.Equal(t, uint16(0xFF44), LYAddress, "LY address constant should match")
	assert.Equal(t, uint16(0xFF45), LYCAddress, "LYC address constant should match")
	assert.Equal(t, uint16(0xFF47), BGPAddress, "BGP address constant should match")
	assert.Equal(t, uint16(0xFF48), OBP0Address, "OBP0 address constant should match")
	assert.Equal(t, uint16(0xFF49), OBP1Address, "OBP1 address constant should match")
	assert.Equal(t, uint16(0xFF4A), WYAddress, "WY address constant should match")
	assert.Equal(t, uint16(0xFF4B), WXAddress, "WX address constant should match")
}

// TestLCDEnableDisableFlow tests the LCD enable/disable workflow
func TestLCDEnableDisableFlow(t *testing.T) {
	// Create a mock ROM for testing
	mockROM := make([]byte, 32*1024) // 32KB ROM
	mockROM[0x0147] = 0x00 // ROM only cartridge
	
	// Create cartridge and interrupt controller
	cart, err := cartridge.LoadROMFromBytes(mockROM, "test-rom")
	assert.NoError(t, err, "Should create test cartridge")
	
	mbc, err := cartridge.CreateMBC(cart)
	assert.NoError(t, err, "Should create MBC")
	interruptController := interrupt.NewInterruptController()
	
	// Create MMU and PPU
	mmu := memory.NewMMU(mbc, interruptController, joypad.NewJoypad())
	ppu := NewPPU()
	mmu.SetPPU(ppu)
	
	// PPU starts with LCD enabled
	assert.True(t, ppu.IsLCDEnabled(), "LCD should start enabled")
	assert.Equal(t, ModeOAMScan, ppu.GetCurrentMode(), "Should start in OAM scan mode")
	
	// Disable LCD through MMU
	mmu.WriteByte(0xFF40, 0x11) // LCDC with LCD disabled (bit 7 = 0)
	assert.False(t, ppu.IsLCDEnabled(), "LCD should be disabled")
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should reset when LCD disabled")
	assert.Equal(t, ModeHBlank, ppu.GetCurrentMode(), "Should be in H-Blank when LCD disabled")
	
	// Re-enable LCD through MMU
	mmu.WriteByte(0xFF40, 0x91) // LCDC with LCD enabled (bit 7 = 1)
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled")
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should reset when LCD enabled")
	assert.Equal(t, ModeOAMScan, ppu.GetCurrentMode(), "Should start in OAM scan when LCD enabled")
}

// TestSTATInterruptIntegration tests STAT interrupt functionality
func TestSTATInterruptIntegration(t *testing.T) {
	// Create a mock ROM for testing
	mockROM := make([]byte, 32*1024) // 32KB ROM
	mockROM[0x0147] = 0x00 // ROM only cartridge
	
	// Create cartridge and interrupt controller
	cart, err := cartridge.LoadROMFromBytes(mockROM, "test-rom")
	assert.NoError(t, err, "Should create test cartridge")
	
	mbc, err := cartridge.CreateMBC(cart)
	assert.NoError(t, err, "Should create MBC")
	interruptController := interrupt.NewInterruptController()
	
	// Create MMU and PPU
	mmu := memory.NewMMU(mbc, interruptController, joypad.NewJoypad())
	ppu := NewPPU()
	mmu.SetPPU(ppu)
	
	// Enable LYC=LY interrupt through MMU
	mmu.WriteByte(0xFF41, 1<<STATLYCInterrupt) // Enable LYC interrupt
	mmu.WriteByte(0xFF45, 100) // Set LYC to 100
	
	// Set LY to match LYC (this would normally happen during PPU update)
	ppu.LY = 100
	shouldInterrupt := ppu.updateLYCFlag()
	
	assert.True(t, shouldInterrupt, "Should trigger interrupt when LY=LYC and interrupt enabled")
	
	// Check that STAT register shows LYC flag set
	stat := mmu.ReadByte(0xFF41)
	assert.Equal(t, uint8(1<<STATLYCFlag), stat&(1<<STATLYCFlag), "LYC flag should be set in STAT")
	
	// Clear LYC match
	ppu.LY = 50
	shouldInterrupt = ppu.updateLYCFlag()
	
	assert.False(t, shouldInterrupt, "Should not trigger interrupt when LY≠LYC")
	
	// Check that STAT register shows LYC flag cleared
	stat = mmu.ReadByte(0xFF41)
	assert.Equal(t, uint8(0), stat&(1<<STATLYCFlag), "LYC flag should be cleared in STAT")
}

// TestPaletteIntegration tests complete palette workflow through MMU
func TestPaletteIntegration(t *testing.T) {
	// Create a mock ROM for testing
	mockROM := make([]byte, 32*1024) // 32KB ROM
	mockROM[0x0147] = 0x00 // ROM only cartridge
	
	// Create cartridge and interrupt controller
	cart, err := cartridge.LoadROMFromBytes(mockROM, "test-rom")
	assert.NoError(t, err, "Should create test cartridge")
	
	mbc, err := cartridge.CreateMBC(cart)
	assert.NoError(t, err, "Should create MBC")
	interruptController := interrupt.NewInterruptController()
	
	// Create MMU and PPU
	mmu := memory.NewMMU(mbc, interruptController, joypad.NewJoypad())
	ppu := NewPPU()
	mmu.SetPPU(ppu)
	
	// Set palettes through MMU
	mmu.WriteByte(0xFF47, 0x39) // BGP: 0→1, 1→2, 2→3, 3→0
	mmu.WriteByte(0xFF48, 0x1B) // OBP0: 0→3, 1→2, 2→1, 3→0
	mmu.WriteByte(0xFF49, 0xE4) // OBP1: identity palette
	
	// Test palette color conversion
	assert.Equal(t, uint8(1), ppu.GetBGColor(0), "BG color 0 should map to 1")
	assert.Equal(t, uint8(3), ppu.GetSpriteColor(0, 0), "Sprite color 0 should map to 3 with OBP0")
	assert.Equal(t, uint8(0), ppu.GetSpriteColor(0, 1), "Sprite color 0 should map to 0 with OBP1")
	
	// Test RGB conversion
	bgRGB := ppu.GetBGColorRGB(0, true) // Color 0 → 1 → Game Boy color 1
	expectedRGB := GetRGBColor(1, true)
	assert.Equal(t, expectedRGB, bgRGB, "BG RGB should match expected Game Boy color")
	
	// Test transparency
	assert.True(t, IsColorTransparent(0), "Color 0 should be transparent for sprites")
}
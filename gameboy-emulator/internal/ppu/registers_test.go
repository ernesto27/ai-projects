package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestLCDCRegister tests all LCDC register functionality
func TestLCDCRegister(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial state
	assert.Equal(t, uint8(0x91), ppu.GetLCDC(), "LCDC should have initial value 0x91")
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled initially")
	
	// Test LCD disable
	ppu.SetLCDC(0x11) // LCD disabled (bit 7 = 0)
	assert.Equal(t, uint8(0x11), ppu.GetLCDC(), "LCDC should be updated")
	assert.False(t, ppu.IsLCDEnabled(), "LCD should be disabled")
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should reset to 0 when LCD disabled")
	assert.Equal(t, ModeHBlank, ppu.GetCurrentMode(), "Mode should be H-Blank when LCD disabled")
	
	// Test LCD re-enable
	ppu.SetLCDC(0x91) // LCD enabled (bit 7 = 1)
	assert.Equal(t, uint8(0x91), ppu.GetLCDC(), "LCDC should be updated")
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled")
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should reset to 0 when LCD enabled")
	assert.Equal(t, ModeOAMScan, ppu.GetCurrentMode(), "Mode should be OAM scan when LCD enabled")
	
	// Test individual LCDC bits through helper methods
	ppu.SetLCDC(0b11111111) // All bits set
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled")
	assert.True(t, ppu.IsWindowEnabled(), "Window should be enabled")
	assert.True(t, ppu.IsSpriteEnabled(), "Sprites should be enabled")
	assert.True(t, ppu.IsBGEnabled(), "Background should be enabled")
	assert.Equal(t, uint8(16), ppu.GetSpriteSize(), "Sprite size should be 16 (8x16)")
	
	ppu.SetLCDC(0b10000000) // Only LCD enabled
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled")
	assert.False(t, ppu.IsWindowEnabled(), "Window should be disabled")
	assert.False(t, ppu.IsSpriteEnabled(), "Sprites should be disabled")
	assert.False(t, ppu.IsBGEnabled(), "Background should be disabled")
	assert.Equal(t, uint8(8), ppu.GetSpriteSize(), "Sprite size should be 8 (8x8)")
}

// TestSTATRegister tests all STAT register functionality
func TestSTATRegister(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial state
	assert.Equal(t, uint8(0x02), ppu.GetSTAT(), "STAT should have initial value 0x02 (OAMScan mode)")
	
	// Test STAT write (only bits 6-3 are writable)
	ppu.SetSTAT(0xFF) // Try to write all bits
	stat := ppu.GetSTAT()
	// Should preserve mode bits (1-0) and LYC flag (2), update interrupt enables (6-3)
	assert.Equal(t, uint8(0x78), stat&0x78, "Interrupt enable bits should be set")
	// PPU starts in OAMScan mode (2), so mode bits should be 0x02
	assert.Equal(t, uint8(0x02), stat&0x03, "Mode bits should be preserved (OAMScan = 2)")
	
	// Test STAT mode updates
	ppu.Mode = ModeVBlank
	ppu.updateSTATMode()
	assert.Equal(t, uint8(0x01), ppu.GetSTAT()&0x03, "Mode bits should reflect V-Blank")
	
	ppu.Mode = ModeOAMScan
	ppu.updateSTATMode()
	assert.Equal(t, uint8(0x02), ppu.GetSTAT()&0x03, "Mode bits should reflect OAM Scan")
	
	ppu.Mode = ModeDrawing
	ppu.updateSTATMode()
	assert.Equal(t, uint8(0x03), ppu.GetSTAT()&0x03, "Mode bits should reflect Drawing")
	
	ppu.Mode = ModeHBlank
	ppu.updateSTATMode()
	assert.Equal(t, uint8(0x00), ppu.GetSTAT()&0x03, "Mode bits should reflect H-Blank")
}

// TestLYCComparison tests LYC=LY functionality
func TestLYCComparison(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial state
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should start at 0")
	assert.Equal(t, uint8(0), ppu.GetLYC(), "LYC should start at 0")
	
	// Set LYC to different value
	ppu.SetLYC(100)
	assert.Equal(t, uint8(100), ppu.GetLYC(), "LYC should be updated")
	assert.False(t, ppu.updateLYCFlag(), "Should not trigger interrupt when LY != LYC")
	assert.Equal(t, uint8(0), ppu.GetSTAT()&(1<<STATLYCFlag), "LYC flag should be clear")
	
	// Enable LYC interrupt
	ppu.SetSTAT(1 << STATLYCInterrupt)
	assert.False(t, ppu.updateLYCFlag(), "Should not trigger interrupt when LY != LYC")
	
	// Set LY to match LYC
	ppu.LY = 100
	assert.True(t, ppu.updateLYCFlag(), "Should trigger interrupt when LY == LYC and interrupt enabled")
	assert.Equal(t, uint8(1<<STATLYCFlag), ppu.GetSTAT()&(1<<STATLYCFlag), "LYC flag should be set")
	
	// Disable LYC interrupt
	ppu.SetSTAT(0)
	assert.False(t, ppu.updateLYCFlag(), "Should not trigger interrupt when interrupt disabled")
	assert.Equal(t, uint8(1<<STATLYCFlag), ppu.GetSTAT()&(1<<STATLYCFlag), "LYC flag should still be set")
	
	// Clear match condition
	ppu.LY = 50
	assert.False(t, ppu.updateLYCFlag(), "Should not trigger interrupt when LY != LYC")
	assert.Equal(t, uint8(0), ppu.GetSTAT()&(1<<STATLYCFlag), "LYC flag should be clear")
}

// TestSTATInterrupts tests STAT interrupt conditions
func TestSTATInterrupts(t *testing.T) {
	ppu := NewPPU()
	
	// Test no interrupts enabled
	ppu.SetSTAT(0x00)
	ppu.Mode = ModeHBlank
	assert.False(t, ppu.ShouldTriggerSTATInterrupt(), "No interrupt when none enabled")
	
	// Test H-Blank interrupt
	ppu.SetSTAT(1 << STATMode0Interrupt)
	ppu.Mode = ModeHBlank
	assert.True(t, ppu.ShouldTriggerSTATInterrupt(), "Should trigger H-Blank interrupt")
	ppu.Mode = ModeVBlank
	assert.False(t, ppu.ShouldTriggerSTATInterrupt(), "Should not trigger in V-Blank mode")
	
	// Test V-Blank interrupt
	ppu.SetSTAT(1 << STATMode1Interrupt)
	ppu.Mode = ModeVBlank
	assert.True(t, ppu.ShouldTriggerSTATInterrupt(), "Should trigger V-Blank interrupt")
	ppu.Mode = ModeOAMScan
	assert.False(t, ppu.ShouldTriggerSTATInterrupt(), "Should not trigger in OAM scan mode")
	
	// Test OAM interrupt
	ppu.SetSTAT(1 << STATMode2Interrupt)
	ppu.Mode = ModeOAMScan
	assert.True(t, ppu.ShouldTriggerSTATInterrupt(), "Should trigger OAM interrupt")
	ppu.Mode = ModeDrawing
	assert.False(t, ppu.ShouldTriggerSTATInterrupt(), "Should not trigger in drawing mode")
	
	// Test Drawing mode (no interrupt)
	ppu.SetSTAT(0xFF) // All interrupts enabled
	ppu.Mode = ModeDrawing
	assert.False(t, ppu.ShouldTriggerSTATInterrupt(), "Drawing mode never triggers STAT interrupt")
}

// TestScrollRegisters tests SCX and SCY registers
func TestScrollRegisters(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial values
	assert.Equal(t, uint8(0), ppu.GetSCX(), "SCX should start at 0")
	assert.Equal(t, uint8(0), ppu.GetSCY(), "SCY should start at 0")
	
	// Test SCX
	ppu.SetSCX(123)
	assert.Equal(t, uint8(123), ppu.GetSCX(), "SCX should be updated")
	
	// Test SCY
	ppu.SetSCY(234)
	assert.Equal(t, uint8(234), ppu.GetSCY(), "SCY should be updated")
	
	// Test boundary values
	ppu.SetSCX(255)
	ppu.SetSCY(255)
	assert.Equal(t, uint8(255), ppu.GetSCX(), "SCX should handle max value")
	assert.Equal(t, uint8(255), ppu.GetSCY(), "SCY should handle max value")
}

// TestWindowRegisters tests WX and WY registers
func TestWindowRegisters(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial values
	assert.Equal(t, uint8(0), ppu.GetWX(), "WX should start at 0")
	assert.Equal(t, uint8(0), ppu.GetWY(), "WY should start at 0")
	
	// Test WX
	ppu.SetWX(167) // Window X coordinate
	assert.Equal(t, uint8(167), ppu.GetWX(), "WX should be updated")
	
	// Test WY
	ppu.SetWY(144) // Window Y coordinate
	assert.Equal(t, uint8(144), ppu.GetWY(), "WY should be updated")
	
	// Test boundary values
	ppu.SetWX(255)
	ppu.SetWY(255)
	assert.Equal(t, uint8(255), ppu.GetWX(), "WX should handle max value")
	assert.Equal(t, uint8(255), ppu.GetWY(), "WY should handle max value")
}

// TestPaletteRegisters tests BGP, OBP0, and OBP1 registers
func TestPaletteRegisters(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial values (0xE4 = standard Game Boy palette)
	assert.Equal(t, uint8(0xE4), ppu.GetBGP(), "BGP should start at 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.GetOBP0(), "OBP0 should start at 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.GetOBP1(), "OBP1 should start at 0xE4")
	
	// Test BGP
	ppu.SetBGP(0x1B) // Custom palette: 00-01-10-11
	assert.Equal(t, uint8(0x1B), ppu.GetBGP(), "BGP should be updated")
	
	// Test OBP0
	ppu.SetOBP0(0x30) // Custom palette
	assert.Equal(t, uint8(0x30), ppu.GetOBP0(), "OBP0 should be updated")
	
	// Test OBP1
	ppu.SetOBP1(0xFC) // Custom palette
	assert.Equal(t, uint8(0xFC), ppu.GetOBP1(), "OBP1 should be updated")
	
	// Test all possible values
	for i := 0; i <= 255; i++ {
		value := uint8(i)
		ppu.SetBGP(value)
		ppu.SetOBP0(value)
		ppu.SetOBP1(value)
		assert.Equal(t, value, ppu.GetBGP(), "BGP should handle all values")
		assert.Equal(t, value, ppu.GetOBP0(), "OBP0 should handle all values")
		assert.Equal(t, value, ppu.GetOBP1(), "OBP1 should handle all values")
	}
}

// TestLYRegisterReadOnly tests that LY register is read-only
func TestLYRegisterReadOnly(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial value
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should start at 0")
	
	// LY can only be modified internally, not through SetLY
	// (There's no SetLY method, which is correct - LY is read-only)
	
	// Test internal LY modification
	ppu.LY = 100
	assert.Equal(t, uint8(100), ppu.GetLY(), "LY should be readable after internal modification")
	
	// Test LY boundaries
	ppu.LY = 0
	assert.Equal(t, uint8(0), ppu.GetLY(), "LY should handle 0")
	
	ppu.LY = 153
	assert.Equal(t, uint8(153), ppu.GetLY(), "LY should handle max scanline 153")
}

// TestRegisterConstants tests that register constants are correct
func TestRegisterConstants(t *testing.T) {
	// Test LCDC bit positions
	assert.Equal(t, uint8(7), LCDCLCDEnable, "LCD Enable bit should be 7")
	assert.Equal(t, uint8(6), LCDCWindowTileMap, "Window Tile Map bit should be 6")
	assert.Equal(t, uint8(5), LCDCWindowEnable, "Window Enable bit should be 5")
	assert.Equal(t, uint8(4), LCDCBGWindowTileData, "BG Window Tile Data bit should be 4")
	assert.Equal(t, uint8(3), LCDCBGTileMap, "BG Tile Map bit should be 3")
	assert.Equal(t, uint8(2), LCDCSpriteSize, "Sprite Size bit should be 2")
	assert.Equal(t, uint8(1), LCDCSpriteEnable, "Sprite Enable bit should be 1")
	assert.Equal(t, uint8(0), LCDCBGPriority, "BG Priority bit should be 0")
	
	// Test STAT bit positions
	assert.Equal(t, uint8(6), STATLYCInterrupt, "LYC Interrupt bit should be 6")
	assert.Equal(t, uint8(5), STATMode2Interrupt, "Mode 2 Interrupt bit should be 5")
	assert.Equal(t, uint8(4), STATMode1Interrupt, "Mode 1 Interrupt bit should be 4")
	assert.Equal(t, uint8(3), STATMode0Interrupt, "Mode 0 Interrupt bit should be 3")
	assert.Equal(t, uint8(2), STATLYCFlag, "LYC Flag bit should be 2")
	
	// Test register addresses
	assert.Equal(t, uint16(0xFF40), LCDCAddress, "LCDC address should be 0xFF40")
	assert.Equal(t, uint16(0xFF41), STATAddress, "STAT address should be 0xFF41")
	assert.Equal(t, uint16(0xFF42), SCYAddress, "SCY address should be 0xFF42")
	assert.Equal(t, uint16(0xFF43), SCXAddress, "SCX address should be 0xFF43")
	assert.Equal(t, uint16(0xFF44), LYAddress, "LY address should be 0xFF44")
	assert.Equal(t, uint16(0xFF45), LYCAddress, "LYC address should be 0xFF45")
	assert.Equal(t, uint16(0xFF47), BGPAddress, "BGP address should be 0xFF47")
	assert.Equal(t, uint16(0xFF48), OBP0Address, "OBP0 address should be 0xFF48")
	assert.Equal(t, uint16(0xFF49), OBP1Address, "OBP1 address should be 0xFF49")
	assert.Equal(t, uint16(0xFF4A), WYAddress, "WY address should be 0xFF4A")
	assert.Equal(t, uint16(0xFF4B), WXAddress, "WX address should be 0xFF4B")
}

// TestRegisterReset tests that Reset properly initializes all registers
func TestRegisterReset(t *testing.T) {
	ppu := NewPPU()
	
	// Modify all registers
	ppu.SetLCDC(0x00)
	ppu.SetSTAT(0xFF)
	ppu.SetSCX(100)
	ppu.SetSCY(200)
	ppu.LY = 50
	ppu.SetLYC(75)
	ppu.SetWX(150)
	ppu.SetWY(125)
	ppu.SetBGP(0x12)
	ppu.SetOBP0(0x34)
	ppu.SetOBP1(0x56)
	
	// Reset PPU
	ppu.Reset()
	
	// Verify all registers are back to initial state
	assert.Equal(t, uint8(0x91), ppu.GetLCDC(), "LCDC should be reset to 0x91")
	assert.Equal(t, uint8(0x00), ppu.GetSTAT(), "STAT should be reset to 0x00")
	assert.Equal(t, uint8(0x00), ppu.GetSCX(), "SCX should be reset to 0")
	assert.Equal(t, uint8(0x00), ppu.GetSCY(), "SCY should be reset to 0")
	assert.Equal(t, uint8(0x00), ppu.GetLY(), "LY should be reset to 0")
	assert.Equal(t, uint8(0x00), ppu.GetLYC(), "LYC should be reset to 0")
	assert.Equal(t, uint8(0x00), ppu.GetWX(), "WX should be reset to 0")
	assert.Equal(t, uint8(0x00), ppu.GetWY(), "WY should be reset to 0")
	assert.Equal(t, uint8(0xE4), ppu.GetBGP(), "BGP should be reset to 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.GetOBP0(), "OBP0 should be reset to 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.GetOBP1(), "OBP1 should be reset to 0xE4")
	assert.True(t, ppu.IsLCDEnabled(), "LCD should be enabled after reset")
	assert.Equal(t, ModeOAMScan, ppu.GetCurrentMode(), "Mode should be OAM scan after reset")
}
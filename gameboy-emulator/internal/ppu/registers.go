// Package ppu - LCD register management for Game Boy PPU
// Implements LCDC, STAT, LY, LYC and palette register access with authentic Game Boy behavior

package ppu

// LCDC (LCD Control) register bit positions - 0xFF40
const (
	LCDCLCDEnable    uint8 = 7 // Bit 7: LCD Enable (0=Off, 1=On)
	LCDCWindowTileMap uint8 = 6 // Bit 6: Window Tile Map (0=9800-9BFF, 1=9C00-9FFF)
	LCDCWindowEnable  uint8 = 5 // Bit 5: Window Enable (0=Off, 1=On)
	LCDCBGWindowTileData uint8 = 4 // Bit 4: BG & Window Tile Data (0=8800-97FF, 1=8000-8FFF)
	LCDCBGTileMap     uint8 = 3 // Bit 3: BG Tile Map (0=9800-9BFF, 1=9C00-9FFF)
	LCDCSpriteSize    uint8 = 2 // Bit 2: Sprite Size (0=8x8, 1=8x16)
	LCDCSpriteEnable  uint8 = 1 // Bit 1: Sprite Enable (0=Off, 1=On)
	LCDCBGPriority    uint8 = 0 // Bit 0: BG & Window Priority (0=Off, 1=On)
)

// STAT (LCD Status) register bit positions - 0xFF41
const (
	STATLYCInterrupt   uint8 = 6 // Bit 6: LYC=LY Interrupt Enable (0=Off, 1=On)
	STATMode2Interrupt uint8 = 5 // Bit 5: Mode 2 OAM Interrupt Enable
	STATMode1Interrupt uint8 = 4 // Bit 4: Mode 1 V-Blank Interrupt Enable  
	STATMode0Interrupt uint8 = 3 // Bit 3: Mode 0 H-Blank Interrupt Enable
	STATLYCFlag        uint8 = 2 // Bit 2: LYC=LY Flag (0=Different, 1=Equal) - Read Only
	// Bits 1-0: Mode Flag (00=H-Blank, 01=V-Blank, 10=OAM, 11=Drawing) - Read Only
)

// LCD Register memory addresses
const (
	LCDCAddress uint16 = 0xFF40 // LCD Control register
	STATAddress uint16 = 0xFF41 // LCD Status register
	SCYAddress  uint16 = 0xFF42 // Background scroll Y
	SCXAddress  uint16 = 0xFF43 // Background scroll X
	LYAddress   uint16 = 0xFF44 // Current scanline (read-only)
	LYCAddress  uint16 = 0xFF45 // LY Compare register
	WYAddress   uint16 = 0xFF4A // Window Y position
	WXAddress   uint16 = 0xFF4B // Window X position
	BGPAddress  uint16 = 0xFF47 // Background palette data
	OBP0Address uint16 = 0xFF48 // Object palette 0 data
	OBP1Address uint16 = 0xFF49 // Object palette 1 data
)

// =============================================================================
// LCDC Register Methods
// =============================================================================

// SetLCDC writes to the LCD Control register (0xFF40)
// Handles LCD enable/disable and updates internal state
func (ppu *PPU) SetLCDC(value uint8) {
	oldLCDEnabled := ppu.LCDEnabled
	oldWindowEnabled := ppu.IsWindowEnabled()
	
	ppu.LCDC = value
	ppu.LCDEnabled = (value & (1 << LCDCLCDEnable)) != 0
	newWindowEnabled := ppu.IsWindowEnabled()
	
	// When LCD is disabled, reset PPU state
	if oldLCDEnabled && !ppu.LCDEnabled {
		ppu.LY = 0
		ppu.Cycles = 0
		ppu.Mode = ModeHBlank
		ppu.updateSTATMode()
		
		// Reset window state when LCD is disabled
		if ppu.windowRenderer != nil {
			ppu.windowRenderer.ResetWindowState()
		}
	}
	
	// When LCD is re-enabled, start in OAM scan mode
	if !oldLCDEnabled && ppu.LCDEnabled {
		ppu.LY = 0
		ppu.Cycles = 0
		ppu.Mode = ModeOAMScan
		ppu.updateSTATMode()
	}
	
	// When window is disabled or enabled, reset window state
	if oldWindowEnabled != newWindowEnabled && ppu.windowRenderer != nil {
		ppu.windowRenderer.ResetWindowState()
	}
}

// GetLCDC reads the LCD Control register (0xFF40)
func (ppu *PPU) GetLCDC() uint8 {
	return ppu.LCDC
}

// IsLCDEnabled is already defined in ppu.go

// IsWindowEnabled returns true if window is enabled (LCDC bit 5)
func (ppu *PPU) IsWindowEnabled() bool {
	return (ppu.LCDC & (1 << LCDCWindowEnable)) != 0
}

// IsSpriteEnabled returns true if sprites are enabled (LCDC bit 1)
func (ppu *PPU) IsSpriteEnabled() bool {
	return (ppu.LCDC & (1 << LCDCSpriteEnable)) != 0
}

// IsBGEnabled returns true if background is enabled (LCDC bit 0)
func (ppu *PPU) IsBGEnabled() bool {
	return (ppu.LCDC & (1 << LCDCBGPriority)) != 0
}

// GetSpriteSize returns sprite height (8 for 8x8, 16 for 8x16)
func (ppu *PPU) GetSpriteSize() uint8 {
	if (ppu.LCDC & (1 << LCDCSpriteSize)) != 0 {
		return 16 // 8x16 sprites
	}
	return 8 // 8x8 sprites
}

// =============================================================================
// STAT Register Methods  
// =============================================================================

// SetSTAT writes to the LCD Status register (0xFF41)
// Only bits 6-3 (interrupt enables) are writable, bits 2-0 are read-only
func (ppu *PPU) SetSTAT(value uint8) {
	// Preserve read-only bits (2-0: LYC flag and mode)
	readOnlyBits := ppu.STAT & 0x07
	// Update writable bits (6-3: interrupt enables)
	writableBits := value & 0x78
	ppu.STAT = writableBits | readOnlyBits
}

// GetSTAT reads the LCD Status register (0xFF41)
func (ppu *PPU) GetSTAT() uint8 {
	return ppu.STAT
}

// updateSTATMode updates the mode bits (1-0) in STAT register
// This is called internally when PPU mode changes
func (ppu *PPU) updateSTATMode() {
	// Clear mode bits and set new mode
	ppu.STAT = (ppu.STAT & 0xFC) | uint8(ppu.Mode)
}

// updateLYCFlag updates the LYC=LY flag (bit 2) in STAT register
// Returns true if LYC=LY and interrupt should be triggered
func (ppu *PPU) updateLYCFlag() bool {
	lycMatch := (ppu.LY == ppu.LYC)
	
	if lycMatch {
		ppu.STAT |= (1 << STATLYCFlag) // Set LYC flag
	} else {
		ppu.STAT &= ^uint8(1 << STATLYCFlag) // Clear LYC flag
	}
	
	// Return true if interrupt should be triggered
	return lycMatch && ((ppu.STAT & (1 << STATLYCInterrupt)) != 0)
}

// ShouldTriggerSTATInterrupt checks if current STAT conditions should trigger interrupt
func (ppu *PPU) ShouldTriggerSTATInterrupt() bool {
	switch ppu.Mode {
	case ModeHBlank:
		return (ppu.STAT & (1 << STATMode0Interrupt)) != 0
	case ModeVBlank:
		return (ppu.STAT & (1 << STATMode1Interrupt)) != 0
	case ModeOAMScan:
		return (ppu.STAT & (1 << STATMode2Interrupt)) != 0
	case ModeDrawing:
		return false // Mode 3 has no interrupt
	default:
		return false
	}
}

// =============================================================================
// LY and LYC Register Methods
// =============================================================================

// GetLY reads the current scanline register (0xFF44) - Read Only
func (ppu *PPU) GetLY() uint8 {
	return ppu.LY
}

// SetLYC writes to the LY Compare register (0xFF45)
func (ppu *PPU) SetLYC(value uint8) {
	ppu.LYC = value
	// Update LYC flag immediately when LYC changes
	ppu.updateLYCFlag()
}

// GetLYC reads the LY Compare register (0xFF45)
func (ppu *PPU) GetLYC() uint8 {
	return ppu.LYC
}

// =============================================================================
// Scroll Register Methods
// =============================================================================

// SetSCY writes to background scroll Y register (0xFF42)
func (ppu *PPU) SetSCY(value uint8) {
	ppu.SCY = value
}

// GetSCY reads background scroll Y register (0xFF42)
func (ppu *PPU) GetSCY() uint8 {
	return ppu.SCY
}

// SetSCX writes to background scroll X register (0xFF43)  
func (ppu *PPU) SetSCX(value uint8) {
	ppu.SCX = value
}

// GetSCX reads background scroll X register (0xFF43)
func (ppu *PPU) GetSCX() uint8 {
	return ppu.SCX
}

// =============================================================================
// Window Register Methods
// =============================================================================

// SetWY writes to window Y position register (0xFF4A)
func (ppu *PPU) SetWY(value uint8) {
	ppu.WY = value
}

// GetWY reads window Y position register (0xFF4A)
func (ppu *PPU) GetWY() uint8 {
	return ppu.WY
}

// SetWX writes to window X position register (0xFF4B)
func (ppu *PPU) SetWX(value uint8) {
	ppu.WX = value
}

// GetWX reads window X position register (0xFF4B)
func (ppu *PPU) GetWX() uint8 {
	return ppu.WX
}

// GetWindowTileMapSelect returns true if window uses tile map 1 (LCDC bit 6)
func (ppu *PPU) GetWindowTileMapSelect() bool {
	return (ppu.LCDC & (1 << LCDCWindowTileMap)) != 0
}

// GetBGWindowTileDataSelect returns true if BG & window use tile data method 1 (LCDC bit 4)
func (ppu *PPU) GetBGWindowTileDataSelect() bool {
	return (ppu.LCDC & (1 << LCDCBGWindowTileData)) != 0
}

// =============================================================================
// Palette Register Methods
// =============================================================================

// SetBGP writes to background palette register (0xFF47)
func (ppu *PPU) SetBGP(value uint8) {
	ppu.BGP = value
}

// GetBGP reads background palette register (0xFF47)
func (ppu *PPU) GetBGP() uint8 {
	return ppu.BGP
}

// SetOBP0 writes to object palette 0 register (0xFF48)
func (ppu *PPU) SetOBP0(value uint8) {
	ppu.OBP0 = value
}

// GetOBP0 reads object palette 0 register (0xFF48)
func (ppu *PPU) GetOBP0() uint8 {
	return ppu.OBP0
}

// SetOBP1 writes to object palette 1 register (0xFF49)
func (ppu *PPU) SetOBP1(value uint8) {
	ppu.OBP1 = value
}

// GetOBP1 reads object palette 1 register (0xFF49)
func (ppu *PPU) GetOBP1() uint8 {
	return ppu.OBP1
}

// =============================================================================
// Background-Specific Register Methods
// =============================================================================

// IsBackgroundEnabled returns true if background is enabled (LCDC bit 0)
func (ppu *PPU) IsBackgroundEnabled() bool {
	return (ppu.LCDC & (1 << LCDCBGPriority)) != 0
}

// GetScrollX returns the background scroll X coordinate
func (ppu *PPU) GetScrollX() uint8 {
	return ppu.GetSCX()
}

// GetScrollY returns the background scroll Y coordinate  
func (ppu *PPU) GetScrollY() uint8 {
	return ppu.GetSCY()
}

// IsBackgroundTileMap1 returns true if background uses tile map 1 (LCDC bit 3)
func (ppu *PPU) IsBackgroundTileMap1() bool {
	return (ppu.LCDC & (1 << LCDCBGTileMap)) != 0
}

// IsBackgroundTileData1 returns true if background uses tile data method 1 (LCDC bit 4)
func (ppu *PPU) IsBackgroundTileData1() bool {
	return (ppu.LCDC & (1 << LCDCBGWindowTileData)) != 0
}

// GetBackgroundPalette returns the background palette for color conversion
func (ppu *PPU) GetBackgroundPalette() [4]uint8 {
	return DecodePalette(ppu.GetBGP())
}
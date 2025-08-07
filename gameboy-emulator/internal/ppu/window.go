// Package ppu - Game Boy window rendering system
// Implements scanline-based window rendering with tile maps and positioning
// The window is a second background layer that appears above background but below sprites

package ppu

import "fmt"

// WindowRenderer handles Game Boy window rendering
// The window system provides a secondary background layer that can be positioned anywhere on screen
// and has priority over the background layer but below sprites
type WindowRenderer struct {
	ppu           *PPU           // Reference to parent PPU for registers and framebuffer
	vramInterface VRAMInterface  // Interface for accessing VRAM and tile data
	
	// Window-specific state
	windowLineCounter uint8 // Internal line counter for window rendering (resets when window is disabled)
	isWindowActive    bool  // True if window was visible on current frame
}

// NewWindowRenderer creates a new window renderer
func NewWindowRenderer(ppu *PPU, vramInterface VRAMInterface) *WindowRenderer {
	return &WindowRenderer{
		ppu:               ppu,
		vramInterface:     vramInterface,
		windowLineCounter: 0,
		isWindowActive:    false,
	}
}

// RenderWindowScanline renders window pixels for a single scanline
// This is called during PPU Drawing mode after background rendering but before sprites
// The window has priority over background but below sprites
func (wr *WindowRenderer) RenderWindowScanline(scanline uint8) {
	// Skip rendering if window is disabled
	if !wr.ppu.IsWindowEnabled() {
		return
	}
	
	// Skip rendering if scanline is out of visible range
	if scanline >= ScreenHeight {
		return
	}
	
	// Check if window is visible on this scanline
	if !wr.isWindowVisibleOnScanline(scanline) {
		return
	}
	
	// Window is visible, mark as active and render pixels
	wr.isWindowActive = true
	
	// Get window position (WX represents screen X - 7, WY represents screen Y directly)
	windowX := int(wr.ppu.GetWX()) - 7  // Convert WX to actual screen X position
	
	// Calculate window internal coordinates
	windowScanlineY := wr.windowLineCounter  // Use internal line counter, not screen scanline
	
	// Determine which tile row we're in and pixel row within that tile
	tileY := int(windowScanlineY) / TileHeight
	pixelRowInTile := int(windowScanlineY) % TileHeight
	
	// Get window tile map base address
	tileMapBase := wr.getWindowTileMapBase()
	
	// Render window pixels starting from the visible X position
	startX := windowX
	if startX < 0 {
		startX = 0  // Clip to screen bounds
	}
	
	for screenX := startX; screenX < ScreenWidth; screenX++ {
		// Calculate window internal X coordinate
		windowInternalX := screenX - windowX
		if windowInternalX < 0 {
			continue  // Skip pixels before window starts
		}
		
		// Determine which tile column we're in and pixel column within that tile
		tileX := windowInternalX / TileWidth
		pixelColInTile := windowInternalX % TileWidth
		
		// Calculate tile map address for this tile position
		tileMapAddress := tileMapBase + uint16(tileY*32 + tileX)
		
		// Read tile index from tile map
		tileIndex := wr.vramInterface.ReadVRAM(tileMapAddress)
		
		// Get tile data address based on tile index and addressing mode
		tileDataAddress := wr.getTileDataAddress(tileIndex)
		
		// Calculate address for specific row of tile data (2 bytes per row)
		rowDataAddress := tileDataAddress + uint16(pixelRowInTile*2)
		
		// Read tile data for this row (2 bytes: low and high color bits)
		lowByte := wr.vramInterface.ReadVRAM(rowDataAddress)
		highByte := wr.vramInterface.ReadVRAM(rowDataAddress + 1)
		
		// Extract color for this pixel position within the tile
		bitPosition := 7 - pixelColInTile  // Game Boy uses MSB first
		colorBit0 := (lowByte >> bitPosition) & 1
		colorBit1 := (highByte >> bitPosition) & 1
		rawColor := (colorBit1 << 1) | colorBit0
		
		// Apply background palette (window uses same palette as background)
		finalColor := wr.applyBackgroundPalette(rawColor)
		
		// Set pixel in framebuffer (window has priority over background)
		wr.ppu.SetPixel(screenX, int(scanline), finalColor)
	}
	
	// Increment window line counter for next scanline
	wr.windowLineCounter++
}

// isWindowVisibleOnScanline determines if the window should be rendered on the given scanline
// Window is visible if it's enabled and the current scanline is >= WY
func (wr *WindowRenderer) isWindowVisibleOnScanline(scanline uint8) bool {
	windowY := wr.ppu.GetWY()
	return scanline >= windowY
}

// getWindowTileMapBase returns the base address for the window tile map
// Based on LCDC bit 6: 0 = 0x9800-0x9BFF, 1 = 0x9C00-0x9FFF
func (wr *WindowRenderer) getWindowTileMapBase() uint16 {
	if wr.ppu.GetWindowTileMapSelect() {
		return 0x9C00  // Upper tile map
	}
	return 0x9800      // Lower tile map (default)
}

// getTileDataAddress calculates the address of tile data based on tile index and addressing mode
// Game Boy has two tile data addressing modes controlled by LCDC bit 4
func (wr *WindowRenderer) getTileDataAddress(tileIndex uint8) uint16 {
	if wr.ppu.GetBGWindowTileDataSelect() {
		// Mode 1: 0x8000-0x8FFF, unsigned indexing (0-255)
		return 0x8000 + uint16(tileIndex)*16
	} else {
		// Mode 0: 0x8800-0x97FF, signed indexing (-128 to 127)
		// Tile index 0-127 maps to 0x9000-0x97FF
		// Tile index 128-255 maps to 0x8800-0x8FFF
		if tileIndex < 128 {
			return 0x9000 + uint16(tileIndex)*16
		} else {
			return 0x8800 + uint16(tileIndex-128)*16
		}
	}
}

// applyBackgroundPalette applies the background palette to a raw color value
// Window uses the same palette as background (BGP register)
func (wr *WindowRenderer) applyBackgroundPalette(rawColor uint8) uint8 {
	palette := wr.ppu.GetBGP()
	
	switch rawColor {
	case 0:
		return (palette >> 0) & 0x03
	case 1:
		return (palette >> 2) & 0x03
	case 2:
		return (palette >> 4) & 0x03
	case 3:
		return (palette >> 6) & 0x03
	default:
		return ColorWhite
	}
}

// ResetWindowState resets window-specific state (called when window is disabled or LCD is reset)
// This resets the internal window line counter
func (wr *WindowRenderer) ResetWindowState() {
	wr.windowLineCounter = 0
	wr.isWindowActive = false
}

// GetWindowLineCounter returns the current window line counter (for debugging/testing)
func (wr *WindowRenderer) GetWindowLineCounter() uint8 {
	return wr.windowLineCounter
}

// IsWindowActive returns true if window was visible on current frame (for debugging/testing)
func (wr *WindowRenderer) IsWindowActive() bool {
	return wr.isWindowActive
}

// ValidateWindowPosition checks if window position is valid and returns adjusted values if needed
// This is mainly for debugging and development purposes
func (wr *WindowRenderer) ValidateWindowPosition() (bool, string) {
	wx := wr.ppu.GetWX()
	wy := wr.ppu.GetWY()
	
	var issues []string
	
	// Check for common window positioning issues
	if wx < 7 {
		issues = append(issues, fmt.Sprintf("WX=%d is less than 7, window will not be visible", wx))
	}
	if wx > 166 {
		issues = append(issues, fmt.Sprintf("WX=%d is greater than 166, window extends beyond screen", wx))
	}
	if wy > 143 {
		issues = append(issues, fmt.Sprintf("WY=%d is greater than 143, window below visible area", wy))
	}
	
	if len(issues) > 0 {
		return false, fmt.Sprintf("Window position issues: %v", issues)
	}
	
	return true, "Window position is valid"
}
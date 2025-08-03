// Package ppu - Game Boy background rendering system
// Implements scanline-based background rendering with tile maps and scrolling

package ppu

import "fmt"

// BackgroundRenderer handles Game Boy background rendering
// Renders one scanline at a time using tiles from VRAM
type BackgroundRenderer struct {
	ppu           *PPU           // Reference to parent PPU for registers and framebuffer
	vramInterface VRAMInterface // Interface for accessing VRAM and tile data
}

// NewBackgroundRenderer creates a new background renderer
func NewBackgroundRenderer(ppu *PPU, vramInterface VRAMInterface) *BackgroundRenderer {
	return &BackgroundRenderer{
		ppu:           ppu,
		vramInterface: vramInterface,
	}
}

// RenderBackgroundScanline renders a single scanline of background graphics
// This is called during PPU Drawing mode for each visible scanline (0-143)
func (br *BackgroundRenderer) RenderBackgroundScanline(scanline uint8) {
	// Skip rendering if background is disabled
	if !br.ppu.IsBackgroundEnabled() {
		br.clearScanline(scanline)
		return
	}
	
	// Skip rendering if scanline is out of visible range
	if scanline >= ScreenHeight {
		return
	}
	
	// Apply scrolling to determine starting position in background map
	scrollX := br.ppu.GetScrollX()
	scrollY := br.ppu.GetScrollY()
	
	// Calculate the background map Y coordinate for this scanline
	backgroundY := uint8((int(scrollY) + int(scanline)) % 256)
	
	// Determine which tile row we're in and pixel row within that tile
	tileY := int(backgroundY) / TileHeight
	pixelRowInTile := int(backgroundY) % TileHeight
	
	// Render all 160 visible pixels for this scanline
	for screenX := 0; screenX < ScreenWidth; screenX++ {
		// Calculate background map X coordinate with scrolling
		backgroundX := uint8((int(scrollX) + screenX) % 256)
		
		// Determine which tile column and pixel column within tile  
		tileX := int(backgroundX) / TileWidth
		pixelColInTile := int(backgroundX) % TileWidth
		
		// Get the tile for this position
		tile := br.fetchBackgroundTile(tileX, tileY)
		if tile == nil {
			// Use empty tile if fetch fails
			tile = NewTile()
		}
		
		// Get the pixel color from the tile
		tileColor := tile.GetPixel(pixelColInTile, pixelRowInTile)
		
		// Apply background palette to convert tile color to final color
		finalColor := br.applyBackgroundPalette(tileColor)
		
		// Write pixel to framebuffer
		br.ppu.SetPixel(screenX, int(scanline), finalColor)
	}
}

// fetchBackgroundTile retrieves a tile from the background tile map
// Uses current LCDC settings to determine tile map and data selection
func (br *BackgroundRenderer) fetchBackgroundTile(tileX, tileY int) *Tile {
	// Ensure tile coordinates are within tile map bounds (32x32)
	tileX = tileX % TileMapWidth
	tileY = tileY % TileMapHeight
	if tileX < 0 {
		tileX += TileMapWidth
	}
	if tileY < 0 {
		tileY += TileMapHeight
	}
	
	// Determine which tile map to use based on LCDC bit 3
	useMap1 := br.ppu.IsBackgroundTileMap1()
	
	// Calculate tile map address
	var tileMapAddress uint16
	if useMap1 {
		// Background map 1 at 0x9C00
		tileMapAddress = BackgroundMap1Start + uint16(tileY*TileMapWidth+tileX)
	} else {
		// Background map 0 at 0x9800
		tileMapAddress = BackgroundMap0Start + uint16(tileY*TileMapWidth+tileX)
	}
	
	// Read tile index from tile map
	tileIndex := br.vramInterface.ReadVRAM(tileMapAddress)
	
	// Determine tile data addressing mode based on LCDC bit 4
	useSignedMode := !br.ppu.IsBackgroundTileData1()
	
	// Calculate tile data address
	var tileDataAddress uint16
	if useSignedMode {
		// $8800 method: signed addressing
		signedIndex := int8(tileIndex)
		tileDataAddress = 0x9000 + uint16(signedIndex)*TileSize
	} else {
		// $8000 method: unsigned addressing  
		tileDataAddress = TilePatternTable0Start + uint16(tileIndex)*TileSize
	}
	
	// Read tile data from VRAM (16 bytes per tile)
	var tileData TileData
	for i := 0; i < TileSize; i++ {
		tileData[i] = br.vramInterface.ReadVRAM(tileDataAddress + uint16(i))
	}
	
	// Convert raw tile data to tile structure
	return NewTileFromData(tileData)
}

// applyBackgroundPalette converts tile color (0-3) to final color using BGP register
func (br *BackgroundRenderer) applyBackgroundPalette(tileColor uint8) uint8 {
	// Clamp tile color to valid range
	if tileColor > 3 {
		tileColor = 3
	}
	
	// Get background palette from PPU
	bgPalette := br.ppu.GetBackgroundPalette()
	
	// Apply palette transformation
	return ApplyPalette(tileColor, bgPalette)
}

// clearScanline fills a scanline with the background color (usually white/color 0)
// Used when background is disabled via LCDC bit 0
func (br *BackgroundRenderer) clearScanline(scanline uint8) {
	if scanline >= ScreenHeight {
		return
	}
	
	// Fill scanline with color 0 (white)
	for x := 0; x < ScreenWidth; x++ {
		br.ppu.SetPixel(x, int(scanline), ColorWhite)
	}
}

// RenderFullBackground renders the entire background (for debugging/testing)
// Not used during normal emulation - only for development and testing
func (br *BackgroundRenderer) RenderFullBackground() {
	for scanline := uint8(0); scanline < ScreenHeight; scanline++ {
		br.RenderBackgroundScanline(scanline)
	}
}

// GetBackgroundPixel returns the background color at screen coordinates (for sprite priority)
// This is used by the sprite system to determine priority interactions
func (br *BackgroundRenderer) GetBackgroundPixel(screenX, screenY int) uint8 {
	// Bounds check
	if screenX < 0 || screenX >= ScreenWidth || screenY < 0 || screenY >= ScreenHeight {
		return ColorWhite
	}
	
	// If background is disabled, return white
	if !br.ppu.IsBackgroundEnabled() {
		return ColorWhite
	}
	
	// Apply scrolling to get background coordinates
	scrollX := br.ppu.GetScrollX()
	scrollY := br.ppu.GetScrollY()
	
	backgroundX := uint8((int(scrollX) + screenX) % 256)
	backgroundY := uint8((int(scrollY) + screenY) % 256)
	
	// Calculate tile coordinates
	tileX := int(backgroundX) / TileWidth
	tileY := int(backgroundY) / TileHeight
	pixelX := int(backgroundX) % TileWidth
	pixelY := int(backgroundY) % TileHeight
	
	// Fetch tile and get pixel
	tile := br.fetchBackgroundTile(tileX, tileY)
	if tile == nil {
		return ColorWhite
	}
	
	tileColor := tile.GetPixel(pixelX, pixelY)
	return br.applyBackgroundPalette(tileColor)
}

// IsBackgroundPixelTransparent checks if a background pixel is transparent (color 0)
// Used for sprite priority calculations
func (br *BackgroundRenderer) IsBackgroundPixelTransparent(screenX, screenY int) bool {
	// Bounds check
	if screenX < 0 || screenX >= ScreenWidth || screenY < 0 || screenY >= ScreenHeight {
		return true // Out of bounds is considered transparent
	}
	
	// If background is disabled, everything is transparent
	if !br.ppu.IsBackgroundEnabled() {
		return true
	}
	
	// Get the raw tile color (before palette application)
	scrollX := br.ppu.GetScrollX()
	scrollY := br.ppu.GetScrollY()
	
	backgroundX := uint8((int(scrollX) + screenX) % 256)
	backgroundY := uint8((int(scrollY) + screenY) % 256)
	
	tileX := int(backgroundX) / TileWidth
	tileY := int(backgroundY) / TileHeight
	pixelX := int(backgroundX) % TileWidth
	pixelY := int(backgroundY) % TileHeight
	
	tile := br.fetchBackgroundTile(tileX, tileY)
	if tile == nil {
		return true
	}
	
	// Tile color 0 is transparent
	tileColor := tile.GetPixel(pixelX, pixelY)
	return tileColor == 0
}

// =============================================================================
// Debug and Analysis Functions
// =============================================================================

// GetVisibleTiles returns information about tiles visible on current screen
// Useful for debugging and development tools
func (br *BackgroundRenderer) GetVisibleTiles() []TileInfo {
	if !br.ppu.IsBackgroundEnabled() {
		return nil
	}
	
	var visibleTiles []TileInfo
	scrollX := br.ppu.GetScrollX()
	scrollY := br.ppu.GetScrollY()
	
	// Calculate range of tiles that might be visible
	startTileX := int(scrollX) / TileWidth
	startTileY := int(scrollY) / TileHeight
	
	// We need to check tiles that might be partially visible
	for tileY := startTileY; tileY <= startTileY+ScreenTilesHeight+1; tileY++ {
		for tileX := startTileX; tileX <= startTileX+ScreenTilesWidth+1; tileX++ {
			// Wrap coordinates to tile map bounds
			wrappedX := tileX % TileMapWidth
			wrappedY := tileY % TileMapHeight
			if wrappedX < 0 {
				wrappedX += TileMapWidth
			}
			if wrappedY < 0 {
				wrappedY += TileMapHeight
			}
			
			// Get tile index
			useMap1 := br.ppu.IsBackgroundTileMap1()
			var tileMapAddress uint16
			if useMap1 {
				tileMapAddress = BackgroundMap1Start + uint16(wrappedY*TileMapWidth+wrappedX)
			} else {
				tileMapAddress = BackgroundMap0Start + uint16(wrappedY*TileMapWidth+wrappedX)
			}
			
			tileIndex := br.vramInterface.ReadVRAM(tileMapAddress)
			
			visibleTiles = append(visibleTiles, TileInfo{
				TileX:     wrappedX,
				TileY:     wrappedY,
				TileIndex: tileIndex,
				ScreenX:   tileX*TileWidth - int(scrollX),
				ScreenY:   tileY*TileHeight - int(scrollY),
			})
		}
	}
	
	return visibleTiles
}

// TileInfo contains information about a tile's position and index
type TileInfo struct {
	TileX     int    // Position in tile map (0-31)
	TileY     int    // Position in tile map (0-31)
	TileIndex uint8  // Tile index from tile map
	ScreenX   int    // Screen position (may be negative or > screen size)
	ScreenY   int    // Screen position (may be negative or > screen size)
}

// AnalyzeBackground returns statistics about current background state
func (br *BackgroundRenderer) AnalyzeBackground() map[string]interface{} {
	analysis := make(map[string]interface{})
	
	analysis["backgroundEnabled"] = br.ppu.IsBackgroundEnabled()
	analysis["scrollX"] = br.ppu.GetScrollX()
	analysis["scrollY"] = br.ppu.GetScrollY()
	analysis["tileMap1Selected"] = br.ppu.IsBackgroundTileMap1()
	analysis["tileData1Selected"] = br.ppu.IsBackgroundTileData1()
	
	if br.ppu.IsBackgroundEnabled() {
		visibleTiles := br.GetVisibleTiles()
		analysis["visibleTileCount"] = len(visibleTiles)
		
		// Count unique tile indices
		uniqueTiles := make(map[uint8]bool)
		for _, tile := range visibleTiles {
			uniqueTiles[tile.TileIndex] = true
		}
		analysis["uniqueTileCount"] = len(uniqueTiles)
		
		// Find most common tile index
		tileFreq := make(map[uint8]int)
		for _, tile := range visibleTiles {
			tileFreq[tile.TileIndex]++
		}
		
		maxFreq := 0
		mostCommonTile := uint8(0)
		for index, freq := range tileFreq {
			if freq > maxFreq {
				maxFreq = freq
				mostCommonTile = index
			}
		}
		
		analysis["mostCommonTile"] = mostCommonTile
		analysis["mostCommonTileFreq"] = maxFreq
	}
	
	return analysis
}

// String returns a human-readable description of the background renderer state
func (br *BackgroundRenderer) String() string {
	if !br.ppu.IsBackgroundEnabled() {
		return "Background Renderer: DISABLED"
	}
	
	scrollX := br.ppu.GetScrollX()
	scrollY := br.ppu.GetScrollY()
	tileMap := "Map 0"
	if br.ppu.IsBackgroundTileMap1() {
		tileMap = "Map 1"
	}
	tileData := "$8000"
	if !br.ppu.IsBackgroundTileData1() {
		tileData = "$8800"
	}
	
	return fmt.Sprintf("Background Renderer: ENABLED | Scroll: (%d,%d) | %s | %s method", 
		scrollX, scrollY, tileMap, tileData)
}

// ValidateRenderer performs consistency checks on the background renderer
func (br *BackgroundRenderer) ValidateRenderer() []string {
	var issues []string
	
	// Check if renderer has required dependencies
	if br.ppu == nil {
		issues = append(issues, "PPU reference is nil")
	}
	if br.vramInterface == nil {
		issues = append(issues, "VRAM interface is nil")
	}
	
	// Check if PPU has required register functionality
	if br.ppu != nil {
		// Try to read scroll registers (should not crash)
		_ = br.ppu.GetScrollX()
		_ = br.ppu.GetScrollY()
		_ = br.ppu.IsBackgroundEnabled()
	}
	
	return issues
}
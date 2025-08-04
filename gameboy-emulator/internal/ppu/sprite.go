// Package ppu implements sprite (OAM) rendering for the Game Boy emulator.
// This file handles Object Attribute Memory parsing, sprite rendering,
// and integration with the PPU graphics pipeline.
package ppu

import "sort"

// OAM (Object Attribute Memory) constants
const (
	// OAM memory layout
	OAMStartAddress = 0xFE00 // Start of OAM in memory map
	OAMEndAddress   = 0xFE9F // End of OAM in memory map
	OAMSize         = 160    // Total OAM size in bytes (40 sprites × 4 bytes)
	
	// Sprite limits and dimensions
	MaxSprites           = 40 // Total sprites in OAM
	MaxSpritesPerLine    = 10 // Maximum sprites rendered per scanline
	SpriteWidth          = 8  // Sprite width in pixels
	SpriteHeight8x8      = 8  // Standard sprite height
	SpriteHeight8x16     = 16 // Tall sprite height (when LCDC.2 = 1)
	SpriteBytesPerSprite = 4  // Bytes per sprite in OAM
	
	// Sprite positioning offsets (Game Boy uses offset coordinates)
	SpriteYOffset = 16 // Y position offset (sprite Y=16 appears at screen Y=0)
	SpriteXOffset = 8  // X position offset (sprite X=8 appears at screen X=0)
)

// Sprite attribute flags (byte 3 of each OAM entry)
const (
	SpriteFlagPriority = 0x80 // Bit 7: 0=Above background, 1=Behind background colors 1-3
	SpriteFlagFlipY    = 0x40 // Bit 6: Vertical flip
	SpriteFlagFlipX    = 0x20 // Bit 5: Horizontal flip
	SpriteFlagPalette  = 0x10 // Bit 4: Palette (0=OBP0, 1=OBP1)
	// Bits 3-0: Unused in DMG (Game Boy), used for VRAM bank and palette in CGB
)

// Sprite represents a single Game Boy sprite with parsed OAM attributes
type Sprite struct {
	// Raw OAM data (4 bytes per sprite)
	Y        uint8 // Y position (0-255, with 16-pixel offset)
	X        uint8 // X position (0-255, with 8-pixel offset)  
	TileID   uint8 // Tile number (0-255, references tile data in VRAM)
	Flags    uint8 // Attribute flags (priority, flip, palette)
	
	// Parsed attributes for easier access
	ScreenY      int  // Actual screen Y position (Y - SpriteYOffset)
	ScreenX      int  // Actual screen X position (X - SpriteXOffset)
	Priority     bool // true = behind background, false = above background
	FlipX        bool // Horizontal flip
	FlipY        bool // Vertical flip
	PaletteNum   uint8 // Palette number (0=OBP0, 1=OBP1)
	OAMIndex     uint8 // Original index in OAM (for priority sorting)
}

// NewSprite creates a Sprite from raw OAM data at the specified OAM index
func NewSprite(oamData [4]uint8, oamIndex uint8) *Sprite {
	sprite := &Sprite{
		Y:        oamData[0],
		X:        oamData[1], 
		TileID:   oamData[2],
		Flags:    oamData[3],
		OAMIndex: oamIndex,
	}
	
	// Parse screen positions (apply Game Boy coordinate offsets)
	sprite.ScreenY = int(sprite.Y) - SpriteYOffset
	sprite.ScreenX = int(sprite.X) - SpriteXOffset
	
	// Parse attribute flags
	sprite.Priority = (sprite.Flags & SpriteFlagPriority) != 0
	sprite.FlipY = (sprite.Flags & SpriteFlagFlipY) != 0
	sprite.FlipX = (sprite.Flags & SpriteFlagFlipX) != 0
	sprite.PaletteNum = 0
	if (sprite.Flags & SpriteFlagPalette) != 0 {
		sprite.PaletteNum = 1
	}
	
	return sprite
}

// IsVisible returns true if the sprite is potentially visible on the given scanline
// considering the sprite's Y position and height
func (s *Sprite) IsVisible(scanline uint8, spriteHeight int) bool {
	// Check if sprite's Y range intersects with the current scanline
	spriteTop := s.ScreenY
	spriteBottom := spriteTop + spriteHeight - 1
	
	// Sprite is visible if the scanline falls within its Y range
	return int(scanline) >= spriteTop && int(scanline) <= spriteBottom
}

// GetTileRow returns the row within the sprite's tile data for the given scanline
// accounting for vertical flipping
func (s *Sprite) GetTileRow(scanline uint8, spriteHeight int) int {
	row := int(scanline) - s.ScreenY
	
	// Apply vertical flip if enabled
	if s.FlipY {
		row = spriteHeight - 1 - row
	}
	
	return row
}

// SpriteRenderer handles sprite rendering and OAM management for the PPU
type SpriteRenderer struct {
	ppu           *PPU           // Reference to parent PPU
	vramInterface VRAMInterface  // Interface for accessing VRAM and OAM
	
	// Sprite data cache
	sprites [MaxSprites]*Sprite // All 40 sprites parsed from OAM
	
	// Per-scanline sprite data
	visibleSprites [MaxSpritesPerLine]*Sprite // Sprites visible on current scanline
	spriteCount    int                        // Number of sprites on current scanline
}

// NewSpriteRenderer creates a new sprite renderer
func NewSpriteRenderer(ppu *PPU, vramInterface VRAMInterface) *SpriteRenderer {
	return &SpriteRenderer{
		ppu:           ppu,
		vramInterface: vramInterface,
	}
}

// ScanOAM scans Object Attribute Memory and builds the sprite cache
// This should be called during PPU Mode 2 (OAM Scan)
func (sr *SpriteRenderer) ScanOAM() {
	// Read all 40 sprites from OAM memory
	for i := 0; i < MaxSprites; i++ {
		oamAddress := uint16(OAMStartAddress + i*SpriteBytesPerSprite)
		
		// Read 4 bytes of OAM data for this sprite
		oamData := [4]uint8{
			sr.vramInterface.ReadOAM(oamAddress),
			sr.vramInterface.ReadOAM(oamAddress + 1),
			sr.vramInterface.ReadOAM(oamAddress + 2),
			sr.vramInterface.ReadOAM(oamAddress + 3),
		}
		
		// Create sprite object
		sr.sprites[i] = NewSprite(oamData, uint8(i))
	}
}

// GetSpritesForScanline finds and sorts sprites visible on the given scanline
// Returns up to 10 sprites according to Game Boy sprite limits
func (sr *SpriteRenderer) GetSpritesForScanline(scanline uint8) []*Sprite {
	// Clear previous scanline data
	sr.spriteCount = 0
	for i := range sr.visibleSprites {
		sr.visibleSprites[i] = nil
	}
	
	// Determine sprite height from LCDC register
	spriteHeight := SpriteHeight8x8
	if sr.ppu.GetSpriteSize() == 16 {
		spriteHeight = SpriteHeight8x16
	}
	
	// Collect visible sprites for this scanline
	var candidates []*Sprite
	for i := 0; i < MaxSprites; i++ {
		sprite := sr.sprites[i]
		if sprite != nil && sprite.IsVisible(scanline, spriteHeight) {
			candidates = append(candidates, sprite)
		}
	}
	
	// Sort sprites by Game Boy priority rules:
	// 1. Primary: X coordinate (smaller X = higher priority)
	// 2. Secondary: OAM index (smaller index = higher priority)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].X == candidates[j].X {
			return candidates[i].OAMIndex < candidates[j].OAMIndex
		}
		return candidates[i].X < candidates[j].X
	})
	
	// Limit to maximum sprites per scanline
	maxSprites := len(candidates)
	if maxSprites > MaxSpritesPerLine {
		maxSprites = MaxSpritesPerLine
	}
	
	// Copy to visible sprites array
	for i := 0; i < maxSprites; i++ {
		sr.visibleSprites[i] = candidates[i]
		sr.spriteCount++
	}
	
	return candidates[:maxSprites]
}

// RenderSpriteScanline renders all visible sprites for the current scanline
// This integrates with the existing PPU rendering pipeline
func (sr *SpriteRenderer) RenderSpriteScanline(scanline uint8) {
	// Skip sprite rendering if sprites are disabled
	if !sr.ppu.GetSpritesEnabled() {
		return
	}
	
	// Get sprites for this scanline
	sprites := sr.GetSpritesForScanline(scanline)
	
	// Render each sprite (in reverse order for proper priority)
	for i := len(sprites) - 1; i >= 0; i-- {
		sr.renderSprite(sprites[i], scanline)
	}
}

// renderSprite renders a single sprite on the given scanline
func (sr *SpriteRenderer) renderSprite(sprite *Sprite, scanline uint8) {
	// Determine sprite height
	spriteHeight := SpriteHeight8x8
	if sr.ppu.GetSpriteSize() == 16 {
		spriteHeight = SpriteHeight8x16
	}
	
	// Get the tile row for this scanline
	tileRow := sprite.GetTileRow(scanline, spriteHeight)
	if tileRow < 0 || tileRow >= spriteHeight {
		return // Invalid row
	}
	
	// Get tile ID (for 8x16 sprites, clear bit 0 for top tile)
	tileID := sprite.TileID
	if spriteHeight == SpriteHeight8x16 {
		if tileRow >= 8 {
			tileID |= 0x01  // Bottom tile (set bit 0)
			tileRow -= 8    // Adjust row for bottom tile
		} else {
			tileID &= 0xFE  // Top tile (clear bit 0)
		}
	}
	
	// Read tile data from VRAM (sprites always use tile data area 0x8000-0x8FFF)
	tileAddress := uint16(0x8000 + uint16(tileID)*16 + uint16(tileRow)*2)
	lowByte := sr.vramInterface.ReadVRAM(tileAddress)
	highByte := sr.vramInterface.ReadVRAM(tileAddress + 1)
	
	// Render each pixel of the sprite row
	for pixelX := 0; pixelX < SpriteWidth; pixelX++ {
		// Calculate screen X position
		screenX := sprite.ScreenX + pixelX
		
		// Skip pixels outside screen bounds
		if screenX < 0 || screenX >= ScreenWidth {
			continue
		}
		
		// Get pixel bit position (accounting for horizontal flip)
		bitPos := pixelX
		if sprite.FlipX {
			bitPos = 7 - pixelX
		}
		
		// Extract color from tile data
		colorBit0 := (lowByte >> (7 - bitPos)) & 1
		colorBit1 := (highByte >> (7 - bitPos)) & 1
		color := colorBit1<<1 | colorBit0
		
		// Skip transparent pixels (color 0)
		if color == 0 {
			continue
		}
		
		// Apply sprite-to-background priority
		if sr.shouldRenderSpritePixel(sprite, screenX, int(scanline), color) {
			// Apply sprite palette
			finalColor := sr.applySpritepalette(color, sprite.PaletteNum)
			sr.ppu.SetPixel(screenX, int(scanline), finalColor)
		}
	}
}

// shouldRenderSpritePixel determines if a sprite pixel should be rendered
// based on sprite-to-background priority rules
func (sr *SpriteRenderer) shouldRenderSpritePixel(sprite *Sprite, x, y int, spriteColor uint8) bool {
	// Get current background pixel color
	bgColor := sr.ppu.GetPixel(x, y)
	
	// Sprite priority rules:
	// - If sprite priority flag is 0: sprite appears above background
	// - If sprite priority flag is 1: sprite appears behind background colors 1-3, above color 0
	if sprite.Priority {
		// Behind background: only render if background is color 0 (white/transparent)
		return bgColor == ColorWhite
	} else {
		// Above background: always render sprite (except transparent sprite pixels, handled earlier)
		return true
	}
}

// applySpritepalette applies the sprite palette (OBP0 or OBP1) to a color value
func (sr *SpriteRenderer) applySpritepalette(color uint8, paletteNum uint8) uint8 {
	var palette uint8
	if paletteNum == 0 {
		palette = sr.ppu.OBP0
	} else {
		palette = sr.ppu.OBP1
	}
	
	// Apply palette mapping (same as background palette system)
	switch color {
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
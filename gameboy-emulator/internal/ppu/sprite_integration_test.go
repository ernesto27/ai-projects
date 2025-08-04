package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestSpriteBackgroundPriority tests sprite rendering priority with background
func TestSpriteBackgroundPriority(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Set up background tile data (solid color 3: low=0xFF, high=0xFF)
	backgroundTile := TileData{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // All pixels color 3 (black)
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	vram.SetTileData(0x8010, backgroundTile) // Tile 1 address
	
	// Set up sprite tile data (checkerboard pattern)
	spriteTile := TileData{
		0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, // Alternating pattern
		0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55,
	}
	vram.SetTileData(0x8020, spriteTile) // Tile 2 address
	
	// Set up background tile map to use tile 1
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1) // Fill background with tile 1
	}
	
	tests := []struct {
		name              string
		spritePriority    bool
		expectedOverride  bool // true if sprite should override background
	}{
		{
			name:             "Sprite above background (priority=false)",
			spritePriority:   false,
			expectedOverride: true,
		},
		{
			name:             "Sprite behind background (priority=true)",
			spritePriority:   true,
			expectedOverride: false, // Background color 3 should remain
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear framebuffer
			ppu.Reset()
			
			// Set up sprite with different priority settings
			spriteFlags := uint8(0x00)
			if tt.spritePriority {
				spriteFlags |= SpriteFlagPriority
			}
			vram.SetSprite(0, 32, 32, 2, spriteFlags) // Y=32, X=32, Tile=2
			
			// Enable background and sprites
			ppu.LCDC = 0x93 // Background and sprites enabled
			
			// Render one scanline where sprite is visible
			scanline := uint8(16) // ScreenY = 32-16 = 16
			
			// First render background
			if ppu.backgroundRenderer != nil {
				ppu.backgroundRenderer.RenderBackgroundScanline(scanline)
			}
			
			// Then render sprites
			if ppu.spriteRenderer != nil {
				ppu.spriteRenderer.ScanOAM()
				ppu.spriteRenderer.RenderSpriteScanline(scanline)
			}
			
			// Check pixel at sprite position (X=32, ScreenX=24)
			pixelColor := ppu.GetPixel(24, 16)
			
			if tt.expectedOverride {
				// Sprite should override background - expect sprite color
				// The exact color depends on the sprite pattern and palette
				assert.NotEqual(t, uint8(3), pixelColor, "Sprite should override background")
			} else {
				// Background should remain - expect background color 3
				assert.Equal(t, uint8(3), pixelColor, "Background should remain visible")
			}
		})
	}
}

// TestSpriteTransparency tests that sprite color 0 pixels are transparent
func TestSpriteTransparency(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Set up background tile data (solid color 2: low=0xAA, high=0x00)
	backgroundTile := TileData{
		0xAA, 0x00, 0xAA, 0x00, 0xAA, 0x00, 0xAA, 0x00, // All pixels color 2
		0xAA, 0x00, 0xAA, 0x00, 0xAA, 0x00, 0xAA, 0x00,
	}
	vram.SetTileData(0x8010, backgroundTile) // Tile 1 address
	
	// Set up sprite tile data with transparent pixels (color 0)
	spriteTile := TileData{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // All pixels color 0 (transparent)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	vram.SetTileData(0x8020, spriteTile) // Tile 2 address
	
	// Set up background tile map
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Set up sprite
	vram.SetSprite(0, 32, 32, 2, 0x00) // Transparent sprite
	
	// Enable background and sprites
	ppu.LCDC = 0x93
	
	// Render scanline
	scanline := uint8(16)
	if ppu.backgroundRenderer != nil {
		ppu.backgroundRenderer.RenderBackgroundScanline(scanline)
	}
	
	// Check what background actually rendered first
	bgColor := ppu.GetPixel(24, 16)
	
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.ScanOAM()
		ppu.spriteRenderer.RenderSpriteScanline(scanline)
	}
	
	// Check that background color remains (sprite is transparent)
	pixelColor := ppu.GetPixel(24, 16) // Sprite position
	assert.Equal(t, bgColor, pixelColor, "Background should show through transparent sprite")
}

// TestSpriteFlipping tests horizontal and vertical sprite flipping functionality
func TestSpriteFlipping(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create an asymmetric sprite pattern for testing flips
	// Pattern: half black (color 3), half transparent (color 0) - horizontally
	// Top half: black, bottom half: transparent - vertically
	spriteTile := TileData{
		// Top half rows (black on left, transparent on right)
		0xF0, 0xF0, // Row 0: 11110000 | 11110000 = left color 3, right color 0
		0xF0, 0xF0, // Row 1
		0xF0, 0xF0, // Row 2
		0xF0, 0xF0, // Row 3
		// Bottom half rows (transparent on left, black on right)  
		0x0F, 0x0F, // Row 4: 00001111 | 00001111 = left color 0, right color 3
		0x0F, 0x0F, // Row 5
		0x0F, 0x0F, // Row 6
		0x0F, 0x0F, // Row 7
	}
	vram.SetTileData(0x8010, spriteTile) // Tile 1 address
	
	tests := []struct {
		name   string
		flags  uint8
		testX  int // X position to test
		testY  int // Y position to test (relative to sprite)
		expect uint8 // Expected color
	}{
		{
			name:   "No flip - top-left (should be black)",
			flags:  0x00,
			testX:  24, // Left side of sprite
			testY:  0,  // Top row  
			expect: 3,  // Black
		},
		{
			name:   "No flip - top-right (should be white)", 
			flags:  0x00,
			testX:  28, // Right side of sprite
			testY:  0,  // Top row
			expect: 0,  // White/transparent - background will show
		},
		{
			name:   "Horizontal flip - top-left (should be white)",
			flags:  SpriteFlagFlipX,
			testX:  24, // Left side of sprite (was right before flip)
			testY:  0,  // Top row
			expect: 0,  // White/transparent
		},
		{
			name:   "Vertical flip - bottom-left (should be black)", 
			flags:  SpriteFlagFlipY,
			testX:  24, // Left side
			testY:  7,  // Bottom row (was top before flip)
			expect: 3,  // Black
		},
		{
			name:   "Both flips - bottom-right (should be black)",
			flags:  SpriteFlagFlipX | SpriteFlagFlipY,
			testX:  28, // Right side (was left before horizontal flip)
			testY:  7,  // Bottom row (was top before vertical flip)
			expect: 3,  // Black
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear framebuffer to white background
			ppu.Reset()
			
			// Set up sprite with flip flags
			vram.SetSprite(0, 32, 32, 1, tt.flags) // Y=32, X=32, Tile=1
			
			// Enable sprites only (no background interference)
			ppu.LCDC = 0x82 // Only sprites enabled
			
			// Render the sprite scanline
			scanline := uint8(16 + tt.testY) // Base sprite Y + row offset
			if ppu.spriteRenderer != nil {
				ppu.spriteRenderer.ScanOAM()
				ppu.spriteRenderer.RenderSpriteScanline(scanline)
			}
			
			// Check the pixel color
			pixelColor := ppu.GetPixel(tt.testX, int(scanline))
			
			if tt.expect == 0 {
				// For transparent pixels, background (white) should remain
				assert.Equal(t, uint8(ColorWhite), pixelColor, "Transparent sprite pixel should show background")
			} else {
				// For non-transparent pixels, just verify it's not transparent
				assert.NotEqual(t, uint8(ColorWhite), pixelColor, "Non-transparent sprite pixel should not be white")
			}
		})
	}
}

// TestSpritePalettes tests sprite palette application (OBP0 vs OBP1)
func TestSpritePalettes(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Set up different sprite palettes
	ppu.OBP0 = 0xE4 // Standard palette: 0→0, 1→1, 2→2, 3→3
	ppu.OBP1 = 0x1B // Inverted palette: 0→3, 1→2, 2→1, 3→0
	
	// Create sprite tile with color 1 pixels
	spriteTile := TileData{
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, // All pixels color 1
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF,
	}
	vram.SetTileData(0x8010, spriteTile) // Tile 1 address
	
	tests := []struct {
		name           string
		paletteFlag    uint8
		expectedColor  uint8
	}{
		{
			name:          "OBP0 palette (color 1 → 1)",
			paletteFlag:   0x00, // Use OBP0
			expectedColor: 1,
		},
		{
			name:          "OBP1 palette (color 1 → 2)", 
			paletteFlag:   SpriteFlagPalette, // Use OBP1
			expectedColor: 2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear framebuffer
			ppu.Reset()
			
			// Set up sprite with palette selection
			vram.SetSprite(0, 32, 32, 1, tt.paletteFlag)
			
			// Enable sprites only
			ppu.LCDC = 0x82
			
			// Render sprite
			scanline := uint8(16)
			if ppu.spriteRenderer != nil {
				ppu.spriteRenderer.ScanOAM()
				ppu.spriteRenderer.RenderSpriteScanline(scanline)
			}
			
			// Check palette-mapped color
			pixelColor := ppu.GetPixel(24, 16)
			assert.Equal(t, tt.expectedColor, pixelColor, "Sprite palette mapping incorrect")
		})
	}
}

// TestSprite8x16Mode tests 8x16 sprite rendering
func TestSprite8x16Mode(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Enable 8x16 sprite mode
	ppu.LCDC = 0x86 // Sprites enabled + 8x16 mode
	
	// Set up tile data for 8x16 sprite (uses tiles 0x42 and 0x43)
	topTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // Top tile: color 1 pattern
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	bottomTile := TileData{
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, // Bottom tile: color 2 pattern
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF,
	}
	
	vram.SetTileData(0x8000+0x42*16, topTile)    // Even tile ID for top
	vram.SetTileData(0x8000+0x43*16, bottomTile) // Odd tile ID for bottom
	
	// Set up 8x16 sprite (tile ID bit 0 will be ignored)
	vram.SetSprite(0, 32, 32, 0x42, 0x00) // Base tile 0x42
	
	// Render top half of sprite (rows 0-7)
	scanline := uint8(16) // Top row of sprite
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.ScanOAM()
		ppu.spriteRenderer.RenderSpriteScanline(scanline)
	}
	
	// Check top half uses top tile pattern (color 1)
	topColor := ppu.GetPixel(24, 16)
	assert.Equal(t, uint8(1), topColor, "Top half should use top tile")
	
	// Clear and render bottom half of sprite (rows 8-15)
	ppu.Reset()
	scanline = uint8(24) // Bottom row of sprite
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.ScanOAM()
		ppu.spriteRenderer.RenderSpriteScanline(scanline)
	}
	
	// Check bottom half uses bottom tile pattern (color 2)
	bottomColor := ppu.GetPixel(24, 24)
	assert.Equal(t, uint8(2), bottomColor, "Bottom half should use bottom tile")
}

// TestSpriteOffScreenClipping tests proper clipping of sprites at screen edges
func TestSpriteOffScreenClipping(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create sprite tile with solid color
	spriteTile := TileData{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // All pixels color 3
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	vram.SetTileData(0x8010, spriteTile) // Tile 1 address
	
	tests := []struct {
		name          string
		spriteX       uint8
		expectedLeft  bool // Should left edge be visible?
		expectedRight bool // Should right edge be visible?
	}{
		{
			name:          "Sprite fully on screen",
			spriteX:       32,  // ScreenX = 24
			expectedLeft:  true,
			expectedRight: true,
		},
		{
			name:          "Sprite partially off left edge",
			spriteX:       4,   // ScreenX = -4 (left half off-screen)
			expectedLeft:  false,
			expectedRight: true,
		},
		{
			name:          "Sprite partially off right edge", 
			spriteX:       164, // ScreenX = 156 (right half off-screen)
			expectedLeft:  true,
			expectedRight: false,
		},
		{
			name:          "Sprite completely off left",
			spriteX:       0,   // ScreenX = -8 (completely off-screen)
			expectedLeft:  false,
			expectedRight: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear framebuffer
			ppu.Reset()
			
			// Set up sprite at test position
			vram.SetSprite(0, 32, tt.spriteX, 1, 0x00) // Y=32, variable X
			
			// Enable sprites only
			ppu.LCDC = 0x82
			
			// Render sprite
			scanline := uint8(16)
			if ppu.spriteRenderer != nil {
				ppu.spriteRenderer.ScanOAM()
				ppu.spriteRenderer.RenderSpriteScanline(scanline)
			}
			
			// Check left edge visibility
			leftColor := ppu.GetPixel(0, 16) // Leftmost screen pixel
			if tt.expectedLeft {
				assert.NotEqual(t, uint8(ColorWhite), leftColor, "Left edge should be visible")
			} else {
				assert.Equal(t, uint8(ColorWhite), leftColor, "Left edge should be clipped")
			}
			
			// Check right edge visibility
			rightColor := ppu.GetPixel(ScreenWidth-1, 16) // Rightmost screen pixel
			if tt.expectedRight {
				assert.NotEqual(t, uint8(ColorWhite), rightColor, "Right edge should be visible")
			} else {
				assert.Equal(t, uint8(ColorWhite), rightColor, "Right edge should be clipped")
			}
		})
	}
}

// TestFullPPUIntegrationWithSprites tests complete PPU pipeline with sprites
func TestFullPPUIntegrationWithSprites(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Set up background tiles
	backgroundTile := TileData{
		0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, // Checkerboard background
		0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA,
	}
	vram.SetTileData(0x8010, backgroundTile) // Tile 1 address
	
	// Set up sprite tiles
	spriteTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // Sprite pattern
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8020, spriteTile) // Tile 2 address
	
	// Set up background tile map
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Set up multiple sprites
	vram.SetSprite(0, 32, 32, 2, 0x00)                    // Normal priority sprite
	vram.SetSprite(1, 32, 48, 2, SpriteFlagPriority)      // Behind background sprite
	vram.SetSprite(2, 32, 24, 2, SpriteFlagFlipX)         // Flipped sprite
	
	// Enable all graphics features
	ppu.LCDC = 0x93 // Background and sprites enabled
	
	// Simulate PPU update cycle for one scanline
	scanline := uint8(16)
	ppu.LY = scanline
	
	// Simulate OAM scan phase
	ppu.Mode = ModeOAMScan
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.ScanOAM()
	}
	
	// Simulate drawing phase
	ppu.Mode = ModeDrawing
	if ppu.backgroundRenderer != nil {
		ppu.backgroundRenderer.RenderBackgroundScanline(scanline)
	}
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.RenderSpriteScanline(scanline)
	}
	
	// Verify sprites were found and processed
	sprites := ppu.spriteRenderer.GetSpritesForScanline(scanline)
	assert.Equal(t, 3, len(sprites), "Should find 3 sprites on scanline")
	
	// Verify pixel rendering worked
	// Check that at least some pixels were modified from the initial white
	pixelsChanged := 0
	for x := 0; x < ScreenWidth; x++ {
		if ppu.GetPixel(x, int(scanline)) != ColorWhite {
			pixelsChanged++
		}
	}
	assert.Greater(t, pixelsChanged, 0, "Some pixels should have been rendered")
	
	// Check specific sprite positions have expected colors
	sprite0Color := ppu.GetPixel(24, int(scanline)) // Sprite 0 position (X=32, ScreenX=24)
	sprite1Color := ppu.GetPixel(40, int(scanline)) // Sprite 1 position (X=48, ScreenX=40)
	sprite2Color := ppu.GetPixel(16, int(scanline)) // Sprite 2 position (X=24, ScreenX=16)
	
	// Verify colors are not all the same (indicating rendering occurred)
	colors := []uint8{sprite0Color, sprite1Color, sprite2Color}
	allSame := true
	for i := 1; i < len(colors); i++ {
		if colors[i] != colors[0] {
			allSame = false
			break
		}
	}
	assert.False(t, allSame, "Sprite positions should have varied colors indicating proper rendering")
}
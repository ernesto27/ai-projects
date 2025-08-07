package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestWindowSpriteInteraction tests window rendering with sprites
// Priority should be: Background < Window < Sprites
func TestWindowSpriteInteraction(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create background tile (color 1)
	backgroundTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // All pixels color 1
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8000, backgroundTile) // Tile 0
	
	// Create window tile (color 2)
	windowTile := TileData{
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, // All pixels color 2
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF,
	}
	vram.SetTileData(0x8010, windowTile) // Tile 1
	
	// Create sprite tile (color 3)
	spriteTile := TileData{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // All pixels color 3
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	vram.SetTileData(0x8020, spriteTile) // Tile 2
	
	// Set up tile maps
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 0) // Background uses tile 0
	}
	for i := uint16(0x9C00); i < 0x9C00+32*32; i++ {
		vram.WriteVRAM(i, 1) // Window uses tile 1
	}
	
	// Position window at (40, 30) and sprite at (50, 30)
	ppu.SetWX(47) // Screen X = 47-7 = 40
	ppu.SetWY(30)
	vram.SetSprite(0, 46, 58, 2, 0x00) // Y=46, X=58, Tile=2 -> Screen pos (42, 50)
	
	// Enable all layers
	ppu.SetLCDC(0xF3) // Background + Window + Sprites enabled, Window tile map 1
	
	// Clear and set up
	ppu.Reset()
	ppu.SetLCDC(0xF3)
	ppu.SetWX(47)
	ppu.SetWY(30)
	
	// Render scanline 35 (window and sprite should be visible)
	scanline := uint8(35)
	
	// Render all layers in proper order
	ppu.backgroundRenderer.RenderBackgroundScanline(scanline)
	ppu.windowRenderer.RenderWindowScanline(scanline)
	ppu.spriteRenderer.ScanOAM()
	ppu.spriteRenderer.RenderSpriteScanline(scanline)
	
	// Test layer priorities:
	// X=30: Background only (color 1)
	bgPixel := ppu.GetPixel(30, int(scanline))
	assert.Equal(t, uint8(1), bgPixel, "Background area should show background color")
	
	// X=45: Window over background (color 2)
	windowPixel := ppu.GetPixel(45, int(scanline))
	assert.Equal(t, uint8(2), windowPixel, "Window area should show window color")
	
	// X=50: Sprite over window (color 3)
	spritePixel := ppu.GetPixel(50, int(scanline))
	assert.Equal(t, uint8(3), spritePixel, "Sprite area should show sprite color")
}

// TestWindowEdgeCases tests various edge cases for window positioning
func TestWindowEdgeCases(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	renderer := ppu.GetWindowRenderer()
	
	// Create window tile
	windowTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // Color 1
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8010, windowTile)
	
	// Set up window tile map
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	tests := []struct {
		name     string
		wx       uint8
		wy       uint8
		scanline uint8
		testX    int
		expected uint8
	}{
		{
			name:     "Window at left edge",
			wx:       7,  // Screen X = 0
			wy:       50,
			scanline: 60,
			testX:    0,
			expected: 1, // Window color
		},
		{
			name:     "Window partially off left",
			wx:       5,  // Screen X = -2
			wy:       50,
			scanline: 60,
			testX:    0,
			expected: 1, // Window color (clipped)
		},
		{
			name:     "Window at right edge",
			wx:       159, // Screen X = 152
			wy:       50,
			scanline: 60,
			testX:    155,
			expected: 1, // Window color
		},
		{
			name:     "Window beyond right edge",
			wx:       170, // Screen X = 163
			wy:       50,
			scanline: 60,
			testX:    155,
			expected: 0, // No window visible (white)
		},
		{
			name:     "Window at top edge",
			wx:       50,
			wy:       0,
			scanline: 0,
			testX:    45,
			expected: 1, // Window color
		},
		{
			name:     "Window beyond bottom",
			wx:       50,
			wy:       150,
			scanline: 100,
			testX:    45,
			expected: 0, // No window visible (white)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear and set up
			ppu.Reset()
			ppu.SetLCDC(0xB1) // Background + Window enabled
			ppu.SetWX(tt.wx)
			ppu.SetWY(tt.wy)
			
			// Render window
			renderer.RenderWindowScanline(tt.scanline)
			
			// Check result
			pixel := ppu.GetPixel(tt.testX, int(tt.scanline))
			assert.Equal(t, tt.expected, pixel, "Window edge case should render correctly")
		})
	}
}

// TestWindowLineCounterAdvanced tests complex window line counter scenarios
func TestWindowLineCounterAdvanced(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	renderer := ppu.GetWindowRenderer()
	
	// Set up window tile
	windowTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8010, windowTile)
	
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Test 1: Window line counter increments only when window is visible
	ppu.Reset()
	ppu.SetLCDC(0xB1)
	ppu.SetWX(50)
	ppu.SetWY(100) // Window starts at Y=100
	
	// Render scanlines 90-110
	for scanline := uint8(90); scanline <= 110; scanline++ {
		renderer.RenderWindowScanline(scanline)
	}
	
	// Window should be visible on scanlines 100-110 (11 scanlines)
	assert.Equal(t, uint8(11), renderer.GetWindowLineCounter(), "Window line counter should be 11")
	
	// Test 2: Window line counter resets when window is disabled
	ppu.SetLCDC(0x91) // Disable window
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should reset when disabled")
	
	// Test 3: Window line counter resets when window is re-enabled
	ppu.SetLCDC(0xB1) // Re-enable window
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should be 0 after re-enable")
}

// TestFullPPUPipelineWithWindow tests complete PPU pipeline including window
func TestFullPPUPipelineWithWindow(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create different colored tiles
	bgTile := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00} // Color 1
	windowTile := TileData{0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF} // Color 2
	spriteTile := TileData{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // Color 3
	
	vram.SetTileData(0x8000, bgTile)
	vram.SetTileData(0x8010, windowTile)
	vram.SetTileData(0x8020, spriteTile)
	
	// Set up tile maps
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 0) // Background
	}
	for i := uint16(0x9C00); i < 0x9C00+32*32; i++ {
		vram.WriteVRAM(i, 1) // Window
	}
	
	// Set up sprite
	vram.SetSprite(0, 46, 66, 2, 0x00) // Sprite at screen position (30, 58)
	
	// Configure PPU
	ppu.Reset()
	ppu.SetLCDC(0xF3) // All features enabled, window uses tile map 1
	ppu.SetWX(57)     // Window at screen X=50
	ppu.SetWY(30)     // Window at screen Y=30
	
	// Simulate PPU rendering pipeline for scanline 35
	scanline := uint8(35)
	ppu.LY = scanline
	
	// This simulates what happens during PPU ModeDrawing
	if ppu.backgroundRenderer != nil {
		ppu.backgroundRenderer.RenderBackgroundScanline(scanline)
	}
	if ppu.windowRenderer != nil {
		ppu.windowRenderer.RenderWindowScanline(scanline)
	}
	if ppu.spriteRenderer != nil {
		ppu.spriteRenderer.ScanOAM()
		ppu.spriteRenderer.RenderSpriteScanline(scanline)
	}
	
	// Verify layering at different X positions
	bgOnly := ppu.GetPixel(30, int(scanline))      // Background only
	windowOnly := ppu.GetPixel(55, int(scanline))  // Window over background
	spriteArea := ppu.GetPixel(58, int(scanline))  // Sprite over window
	
	assert.Equal(t, uint8(1), bgOnly, "Background-only area should show background")
	assert.Equal(t, uint8(2), windowOnly, "Window area should show window")
	assert.Equal(t, uint8(3), spriteArea, "Sprite area should show sprite (highest priority)")
	
	// Check window state
	assert.Equal(t, uint8(1), ppu.windowRenderer.GetWindowLineCounter(), "Window line counter should advance")
	assert.True(t, ppu.windowRenderer.IsWindowActive(), "Window should be marked active")
}

// TestWindowPPUModeIntegration tests window behavior during actual PPU mode transitions
func TestWindowPPUModeIntegration(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Set up basic window
	windowTile := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00}
	vram.SetTileData(0x8010, windowTile)
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Configure PPU for window rendering
	ppu.Reset()
	ppu.SetLCDC(0xB1) // Background + Window enabled
	ppu.SetWX(50)
	ppu.SetWY(50)
	
	// Simulate PPU mode transitions (this tests the integrated rendering pipeline)
	// Normally this would be called by the PPU Update method
	for scanline := uint8(45); scanline < 60; scanline++ {
		ppu.LY = scanline
		ppu.Mode = ModeDrawing
		
		// This simulates what happens in the actual PPU.Update method
		if ppu.backgroundRenderer != nil {
			ppu.backgroundRenderer.RenderBackgroundScanline(ppu.LY)
		}
		if ppu.windowRenderer != nil {
			ppu.windowRenderer.RenderWindowScanline(ppu.LY)
		}
		if ppu.spriteRenderer != nil {
			ppu.spriteRenderer.ScanOAM()
			ppu.spriteRenderer.RenderSpriteScanline(ppu.LY)
		}
	}
	
	// Window should have been visible for scanlines 50-59 (10 lines)
	assert.Equal(t, uint8(10), ppu.windowRenderer.GetWindowLineCounter(), "Window line counter should be 10")
	
	// Check that pixels were actually rendered
	pixelCount := 0
	for x := 43; x < ScreenWidth; x++ { // Window starts at X=43 (50-7)
		if ppu.GetPixel(x, 55) != ColorWhite {
			pixelCount++
		}
	}
	assert.Greater(t, pixelCount, 0, "Window pixels should have been rendered")
}
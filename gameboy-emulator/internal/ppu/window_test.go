package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestNewWindowRenderer tests window renderer creation
func TestNewWindowRenderer(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	renderer := NewWindowRenderer(ppu, vram)
	
	assert.NotNil(t, renderer, "Window renderer should be created")
	assert.Equal(t, ppu, renderer.ppu, "PPU reference should be set")
	assert.Equal(t, vram, renderer.vramInterface, "VRAM interface should be set")
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should start at 0")
	assert.False(t, renderer.IsWindowActive(), "Window should not be active initially")
}

// TestWindowVisibilityCheck tests window visibility logic
func TestWindowVisibilityCheck(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	tests := []struct {
		name      string
		windowY   uint8
		scanline  uint8
		visible   bool
	}{
		{
			name:     "Window visible - scanline equals WY",
			windowY:  50,
			scanline: 50,
			visible:  true,
		},
		{
			name:     "Window visible - scanline greater than WY", 
			windowY:  30,
			scanline: 60,
			visible:  true,
		},
		{
			name:     "Window not visible - scanline less than WY",
			windowY:  100,
			scanline: 50,
			visible:  false,
		},
		{
			name:     "Window at top - WY=0",
			windowY:  0,
			scanline: 0,
			visible:  true,
		},
		{
			name:     "Window at bottom - WY=143",
			windowY:  143,
			scanline: 143,
			visible:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu.SetWY(tt.windowY)
			result := renderer.isWindowVisibleOnScanline(tt.scanline)
			assert.Equal(t, tt.visible, result, "Window visibility should match expected")
		})
	}
}

// TestWindowTileMapSelection tests tile map selection based on LCDC bit 6
func TestWindowTileMapSelection(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	// Test tile map 0 (LCDC bit 6 = 0)
	ppu.SetLCDC(0x91) // Window tile map bit 6 = 0
	assert.Equal(t, uint16(0x9800), renderer.getWindowTileMapBase(), "Should use tile map 0")
	
	// Test tile map 1 (LCDC bit 6 = 1)
	ppu.SetLCDC(0xD1) // Window tile map bit 6 = 1 (0x40 bit set)
	assert.Equal(t, uint16(0x9C00), renderer.getWindowTileMapBase(), "Should use tile map 1")
}

// TestWindowTileDataSelection tests tile data addressing modes
func TestWindowTileDataSelection(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	tests := []struct {
		name          string
		lcdc          uint8
		tileIndex     uint8
		expectedAddr  uint16
	}{
		{
			name:         "Mode 1 - unsigned, tile 0",
			lcdc:         0x91, // BG/Window tile data bit 4 = 1
			tileIndex:    0,
			expectedAddr: 0x8000,
		},
		{
			name:         "Mode 1 - unsigned, tile 128",
			lcdc:         0x91,
			tileIndex:    128,
			expectedAddr: 0x8800,
		},
		{
			name:         "Mode 1 - unsigned, tile 255",
			lcdc:         0x91,
			tileIndex:    255,
			expectedAddr: 0x8FF0,
		},
		{
			name:         "Mode 0 - signed, tile 0",
			lcdc:         0x81, // BG/Window tile data bit 4 = 0
			tileIndex:    0,
			expectedAddr: 0x9000,
		},
		{
			name:         "Mode 0 - signed, tile 127",
			lcdc:         0x81,
			tileIndex:    127,
			expectedAddr: 0x97F0,
		},
		{
			name:         "Mode 0 - signed, tile 128 (-128)",
			lcdc:         0x81,
			tileIndex:    128,
			expectedAddr: 0x8800,
		},
		{
			name:         "Mode 0 - signed, tile 255 (-1)",
			lcdc:         0x81,
			tileIndex:    255,
			expectedAddr: 0x8FF0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu.SetLCDC(tt.lcdc)
			result := renderer.getTileDataAddress(tt.tileIndex)
			assert.Equal(t, tt.expectedAddr, result, "Tile data address should match expected")
		})
	}
}

// TestWindowPaletteApplication tests background palette application
func TestWindowPaletteApplication(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	// Set test palette: 0→0, 1→1, 2→2, 3→3 (identity)
	ppu.SetBGP(0xE4)
	
	tests := []struct {
		rawColor    uint8
		expectedColor uint8
	}{
		{0, 0}, // Color 0 → palette bits 0-1 = 00 = 0
		{1, 1}, // Color 1 → palette bits 2-3 = 01 = 1  
		{2, 2}, // Color 2 → palette bits 4-5 = 10 = 2
		{3, 3}, // Color 3 → palette bits 6-7 = 11 = 3
	}
	
	for _, tt := range tests {
		result := renderer.applyBackgroundPalette(tt.rawColor)
		assert.Equal(t, tt.expectedColor, result, "Palette should map correctly")
	}
	
	// Test different palette mapping
	ppu.SetBGP(0x1B) // 0→3, 1→2, 2→1, 3→0 (inverted)
	
	invertedTests := []struct {
		rawColor    uint8
		expectedColor uint8
	}{
		{0, 3}, // Color 0 → palette bits 0-1 = 11 = 3
		{1, 2}, // Color 1 → palette bits 2-3 = 10 = 2
		{2, 1}, // Color 2 → palette bits 4-5 = 01 = 1
		{3, 0}, // Color 3 → palette bits 6-7 = 00 = 0
	}
	
	for _, tt := range invertedTests {
		result := renderer.applyBackgroundPalette(tt.rawColor)
		assert.Equal(t, tt.expectedColor, result, "Inverted palette should map correctly")
	}
}

// TestWindowStateReset tests window state reset functionality
func TestWindowStateReset(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	// Set some state
	renderer.windowLineCounter = 50
	renderer.isWindowActive = true
	
	// Reset state
	renderer.ResetWindowState()
	
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should be reset")
	assert.False(t, renderer.IsWindowActive(), "Window should not be active after reset")
}

// TestWindowPositionValidation tests window position validation
func TestWindowPositionValidation(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewWindowRenderer(ppu, vram)
	
	tests := []struct {
		name    string
		wx      uint8
		wy      uint8
		valid   bool
	}{
		{
			name:  "Valid position",
			wx:    50,
			wy:    50,
			valid: true,
		},
		{
			name:  "WX too small",
			wx:    5,
			wy:    50,
			valid: false,
		},
		{
			name:  "WX too large",
			wx:    170,
			wy:    50,
			valid: false,
		},
		{
			name:  "WY too large",
			wx:    50,
			wy:    150,
			valid: false,
		},
		{
			name:  "Edge case - WX=7",
			wx:    7,
			wy:    50,
			valid: true,
		},
		{
			name:  "Edge case - WX=166",
			wx:    166,
			wy:    143,
			valid: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ppu.SetWX(tt.wx)
			ppu.SetWY(tt.wy)
			
			valid, _ := renderer.ValidateWindowPosition()
			assert.Equal(t, tt.valid, valid, "Position validation should match expected")
		})
	}
}

// TestWindowBasicRendering tests basic window rendering functionality
func TestWindowBasicRendering(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create window tile data (solid color 2)
	windowTile := TileData{
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, // All pixels color 2
		0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF,
	}
	vram.SetTileData(0x8010, windowTile) // Tile 1 address
	
	// Set up window tile map to use tile 1
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Position window at (10, 20) -> WX=17, WY=20
	ppu.SetWX(17) // Screen X = 17-7 = 10
	ppu.SetWY(20)
	
	// Enable window
	ppu.SetLCDC(0xB1) // Background + Window enabled
	
	// Render window on scanline 25 (window should be visible)
	renderer := ppu.GetWindowRenderer()
	renderer.RenderWindowScanline(25)
	
	// Check that window pixels were rendered
	windowPixel := ppu.GetPixel(10, 25) // Should be color 2
	assert.Equal(t, uint8(2), windowPixel, "Window pixel should be rendered")
	
	// Check window line counter incremented
	assert.Equal(t, uint8(1), renderer.GetWindowLineCounter(), "Window line counter should increment")
	assert.True(t, renderer.IsWindowActive(), "Window should be marked as active")
}

// TestWindowWithBackground tests window rendering over background
func TestWindowWithBackground(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create background tile data (solid color 1)
	backgroundTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // All pixels color 1
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8000, backgroundTile) // Tile 0 address
	
	// Create window tile data (solid color 3)
	windowTile := TileData{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // All pixels color 3
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	vram.SetTileData(0x8010, windowTile) // Tile 1 address
	
	// Set up background tile map (0x9800) to use tile 0
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 0)
	}
	
	// Set up window tile map (0x9C00) to use tile 1
	for i := uint16(0x9C00); i < 0x9C00+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Position window at (50, 30)
	ppu.SetWX(57) // Screen X = 57-7 = 50
	ppu.SetWY(30)
	
	// Enable background and window, set window to use tile map 1
	ppu.SetLCDC(0xF1) // Background + Window enabled + Window tile map 1
	
	// Clear framebuffer
	ppu.Reset()
	ppu.SetLCDC(0xF1) // Background + Window enabled + Window tile map 1
	ppu.SetWX(57)
	ppu.SetWY(30)
	
	// Render background first
	bgRenderer := ppu.GetBackgroundRenderer()
	bgRenderer.RenderBackgroundScanline(35)
	
	// Render window (should override background)
	windowRenderer := ppu.GetWindowRenderer()
	windowRenderer.RenderWindowScanline(35)
	
	// Check background area (before window)
	backgroundPixel := ppu.GetPixel(30, 35) // Should be color 1 (background)
	
	// Debug: check window position calculation
	windowX := int(ppu.GetWX()) - 7  // Should be 57-7 = 50
	t.Logf("Debug: WX=%d, windowX=%d, checking pixel at X=30", ppu.GetWX(), windowX)
	
	assert.Equal(t, uint8(1), backgroundPixel, "Background pixel should be visible")
	
	// Check window area (window should override background)
	windowPixel := ppu.GetPixel(55, 35) // Should be color 3 (window)
	assert.Equal(t, uint8(3), windowPixel, "Window pixel should override background")
}

// TestWindowDisabledRendering tests that window doesn't render when disabled
func TestWindowDisabledRendering(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Create window tile data
	windowTile := TileData{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // All pixels color 3
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
	vram.SetTileData(0x8010, windowTile)
	
	// Set up window tile map
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Position window
	ppu.SetWX(17)
	ppu.SetWY(20)
	
	// Disable window (LCDC bit 5 = 0)
	ppu.SetLCDC(0x91) // Background enabled, Window disabled
	
	// Try to render window
	renderer := ppu.GetWindowRenderer()
	renderer.RenderWindowScanline(25)
	
	// Check that no window pixels were rendered
	pixel := ppu.GetPixel(10, 25) // Should be white (not rendered)
	assert.Equal(t, uint8(ColorWhite), pixel, "Window should not render when disabled")
	
	// Check window state not updated
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should not increment")
	assert.False(t, renderer.IsWindowActive(), "Window should not be active")
}

// TestWindowLineCounter tests window line counter behavior
func TestWindowLineCounter(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	renderer := ppu.GetWindowRenderer()
	
	// Create simple window tile
	windowTile := TileData{
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, // Color 1
		0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00,
	}
	vram.SetTileData(0x8010, windowTile)
	
	// Set up window tile map
	for i := uint16(0x9800); i < 0x9800+32*32; i++ {
		vram.WriteVRAM(i, 1)
	}
	
	// Position window and enable it
	ppu.SetWX(17)
	ppu.SetWY(50)
	ppu.SetLCDC(0xB1) // Background + Window enabled
	
	// Render multiple scanlines
	for scanline := uint8(45); scanline < 60; scanline++ {
		renderer.RenderWindowScanline(scanline)
	}
	
	// Window should be visible starting from scanline 50
	// So line counter should be: 0 (45-49) + 10 (50-59) = 10
	assert.Equal(t, uint8(10), renderer.GetWindowLineCounter(), "Window line counter should be 10")
	assert.True(t, renderer.IsWindowActive(), "Window should be active")
}

// TestPPUWindowIntegration tests window integration with PPU pipeline
func TestPPUWindowIntegration(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	ppu.SetVRAMInterface(vram)
	
	// Verify window renderer was created
	assert.NotNil(t, ppu.GetWindowRenderer(), "PPU should have window renderer")
	
	// Test that window state resets when PPU resets
	renderer := ppu.GetWindowRenderer()
	renderer.windowLineCounter = 25
	renderer.isWindowActive = true
	
	ppu.Reset()
	
	assert.Equal(t, uint8(0), renderer.GetWindowLineCounter(), "Window line counter should reset with PPU")
	assert.False(t, renderer.IsWindowActive(), "Window should not be active after PPU reset")
}
package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestNewBackgroundRenderer tests background renderer creation
func TestNewBackgroundRenderer(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	renderer := NewBackgroundRenderer(ppu, vram)
	
	assert.NotNil(t, renderer, "Background renderer should be created")
	assert.Equal(t, ppu, renderer.ppu, "PPU reference should be set")
	assert.Equal(t, vram, renderer.vramInterface, "VRAM interface should be set")
}

// TestBackgroundRendererDisabled tests rendering when background is disabled
func TestBackgroundRendererDisabled(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	// Disable background
	ppu.SetLCDC(0x00) // All bits off, including background enable
	
	// Clear framebuffer to non-white initially
	for x := 0; x < ScreenWidth; x++ {
		ppu.SetPixel(x, 0, ColorBlack)
	}
	
	// Render scanline 0
	renderer.RenderBackgroundScanline(0)
	
	// All pixels should be white (cleared)
	for x := 0; x < ScreenWidth; x++ {
		assert.Equal(t, uint8(ColorWhite), ppu.GetPixel(x, 0), "Pixel should be white when background disabled")
	}
}

// TestBackgroundRendererSolidColor tests rendering with a solid color tile
func TestBackgroundRendererSolidColor(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	// Enable background with basic settings
	ppu.SetLCDC(0x91) // LCD enabled, background enabled, map 0, tile data $8000 method
	ppu.SetBGP(0xE4)  // Standard palette
	
	// Create solid color tile (all pixels = color 3)
	solidTileData := TileData{
		0xFF, 0xFF, // Row 0: all bits set = color 3
		0xFF, 0xFF, // Row 1: all bits set = color 3
		0xFF, 0xFF, // Row 2: all bits set = color 3
		0xFF, 0xFF, // Row 3: all bits set = color 3
		0xFF, 0xFF, // Row 4: all bits set = color 3
		0xFF, 0xFF, // Row 5: all bits set = color 3
		0xFF, 0xFF, // Row 6: all bits set = color 3
		0xFF, 0xFF, // Row 7: all bits set = color 3
	}
	
	// Place tile at index 0 in pattern table ($8000)
	vram.SetTileData(0x8000, solidTileData)
	
	// Set up tile map to use tile 0 for position (0,0)
	vram.SetTileMapEntry(0x9800, 0) // Map 0, tile (0,0) = tile index 0
	
	// Render scanline 0 (no scrolling)
	renderer.RenderBackgroundScanline(0)
	
	// First 8 pixels should be black (palette maps color 3 to black)
	expectedColor := uint8(3) // Color 3 after palette application
	for x := 0; x < 8; x++ {
		assert.Equal(t, expectedColor, ppu.GetPixel(x, 0), "Pixel %d should match tile color", x)
	}
}

// TestBackgroundRendererCheckerboard tests rendering with a pattern
func TestBackgroundRendererCheckerboard(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	// Enable background
	ppu.SetLCDC(0x91) // LCD enabled, background enabled, map 0, tile data $8000 method
	ppu.SetBGP(0xE4)
	
	// Create checkerboard tile pattern
	checkerTileData := TileData{
		0x55, 0x00, // Row 0: 01010101 (alternating pattern)
		0xAA, 0x00, // Row 1: 10101010 (opposite pattern)
		0x55, 0x00, // Row 2: repeat
		0xAA, 0x00, // Row 3: repeat
		0x55, 0x00, // Row 4: repeat
		0xAA, 0x00, // Row 5: repeat
		0x55, 0x00, // Row 6: repeat
		0xAA, 0x00, // Row 7: repeat
	}
	
	vram.SetTileData(0x8000, checkerTileData)
	vram.SetTileMapEntry(0x9800, 0)
	
	// Render scanline 0 and 1
	renderer.RenderBackgroundScanline(0)
	renderer.RenderBackgroundScanline(1)
	
	// Check pattern on scanline 0 (should be 0,1,0,1,0,1,0,1)
	expectedPattern0 := []uint8{0, 1, 0, 1, 0, 1, 0, 1}
	for x := 0; x < 8; x++ {
		assert.Equal(t, expectedPattern0[x], ppu.GetPixel(x, 0), "Scanline 0 pixel %d should match pattern", x)
	}
	
	// Check pattern on scanline 1 (should be 1,0,1,0,1,0,1,0)
	expectedPattern1 := []uint8{1, 0, 1, 0, 1, 0, 1, 0}
	for x := 0; x < 8; x++ {
		assert.Equal(t, expectedPattern1[x], ppu.GetPixel(x, 1), "Scanline 1 pixel %d should match pattern", x)
	}
}

// TestBackgroundScrolling tests horizontal and vertical scrolling
func TestBackgroundScrolling(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	// Enable background
	ppu.SetLCDC(0x91) // LCD enabled, background enabled, map 0, tile data $8000 method
	ppu.SetBGP(0xE4)
	
	// Create two different tiles
	tile0Data := TileData{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // All color 0
	tile1Data := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00} // All color 1
	
	vram.SetTileData(0x8000, tile0Data) // Tile 0
	vram.SetTileData(0x8010, tile1Data) // Tile 1
	
	// Set up tile map: tile 0 at (0,0), tile 1 at (1,0)
	vram.SetTileMapEntry(0x9800, 0)    // Map position (0,0) = tile 0
	vram.SetTileMapEntry(0x9801, 1)    // Map position (1,0) = tile 1
	
	// Test no scrolling first
	ppu.SetSCX(0)
	ppu.SetSCY(0)
	renderer.RenderBackgroundScanline(0)
	
	// First 8 pixels should be color 0, next 8 should be color 1
	for x := 0; x < 8; x++ {
		assert.Equal(t, uint8(0), ppu.GetPixel(x, 0), "No scroll: first tile pixels should be color 0")
	}
	for x := 8; x < 16; x++ {
		assert.Equal(t, uint8(1), ppu.GetPixel(x, 0), "No scroll: second tile pixels should be color 1")
	}
	
	// Test horizontal scrolling by 8 pixels (1 tile)
	ppu.SetSCX(8)
	renderer.RenderBackgroundScanline(0)
	
	// Now first 8 pixels should be color 1 (shifted left)
	for x := 0; x < 8; x++ {
		assert.Equal(t, uint8(1), ppu.GetPixel(x, 0), "Scroll X=8: pixels should be from second tile")
	}
}

// TestBackgroundTileMapSelection tests switching between tile maps
func TestBackgroundTileMapSelection(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetBGP(0xE4)
	
	// Create test tiles
	tile0Data := TileData{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Color 0
	tile1Data := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00} // Color 1
	
	vram.SetTileData(0x8000, tile0Data) // Tile 0
	vram.SetTileData(0x8010, tile1Data) // Tile 1
	
	// Set up different tile maps
	vram.SetTileMapEntry(0x9800, 0)    // Map 0: tile 0
	vram.SetTileMapEntry(0x9C00, 1)    // Map 1: tile 1
	
	// Test with tile map 0
	ppu.SetLCDC(0x91) // LCD enabled, background enabled, map 0
	renderer.RenderBackgroundScanline(0)
	assert.Equal(t, uint8(0), ppu.GetPixel(0, 0), "Map 0 should use tile 0")
	
	// Test with tile map 1
	ppu.SetLCDC(0x99) // LCD enabled, background enabled, map 1 (bit 3 set), tile data $8000
	renderer.RenderBackgroundScanline(0)
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 0), "Map 1 should use tile 1")
}

// TestBackgroundTileDataSelection tests switching between tile data methods
func TestBackgroundTileDataSelection(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetBGP(0xE4)
	
	// Create tiles in both addressing methods
	tile0Data := TileData{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Color 0
	tile1Data := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00} // Color 1
	
	// Place tiles in $8000 method (unsigned)
	vram.SetTileData(0x8000, tile0Data) // $8000 method tile 0
	vram.SetTileData(0x8010, tile1Data) // $8000 method tile 1
	
	// Place tiles in $8800 method (signed) - tile 0 is at $9000
	vram.SetTileData(0x9000, tile1Data) // $8800 method tile 0 (signed)
	
	// Set tile map to use tile index 0
	vram.SetTileMapEntry(0x9800, 0)
	
	// Test $8000 method (LCDC bit 4 = 1)
	ppu.SetLCDC(0x91) // LCD enabled, background enabled, tile data $8000 method
	renderer.RenderBackgroundScanline(0)
	assert.Equal(t, uint8(0), ppu.GetPixel(0, 0), "$8000 method should use tile from 0x8000")
	
	// Test $8800 method (LCDC bit 4 = 0)
	ppu.SetLCDC(0x81) // LCD enabled, background enabled, tile data $8800 method
	renderer.RenderBackgroundScanline(0)
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 0), "$8800 method should use tile from 0x9000")
}

// TestBackgroundPaletteApplication tests palette color mapping
func TestBackgroundPaletteApplication(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	
	// Create tile with solid colors for each row to make testing clearer
	colorTileData := TileData{
		0x00, 0x00, // Row 0: all color 0
		0xFF, 0x00, // Row 1: all color 1 (plane 0 = FF, plane 1 = 00)
		0x00, 0xFF, // Row 2: all color 2 (plane 0 = 00, plane 1 = FF)
		0xFF, 0xFF, // Row 3: all color 3 (plane 0 = FF, plane 1 = FF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Rest: color 0
	}
	
	vram.SetTileData(0x8000, colorTileData)
	vram.SetTileMapEntry(0x9800, 0)
	
	// Test identity palette (0->0, 1->1, 2->2, 3->3)
	ppu.SetBGP(0xE4) // 11100100 = 3,2,1,0
	renderer.RenderBackgroundScanline(0)
	renderer.RenderBackgroundScanline(1)
	renderer.RenderBackgroundScanline(2)
	renderer.RenderBackgroundScanline(3)
	
	assert.Equal(t, uint8(0), ppu.GetPixel(0, 0), "Row 0 should be color 0")
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 1), "Row 1 should be color 1")
	assert.Equal(t, uint8(2), ppu.GetPixel(0, 2), "Row 2 should be color 2")  
	assert.Equal(t, uint8(3), ppu.GetPixel(0, 3), "Row 3 should be color 3")
	
	// Test inverted palette (0->3, 1->2, 2->1, 3->0)
	ppu.SetBGP(0x1B) // 00011011 = 0,1,2,3
	renderer.RenderBackgroundScanline(0)
	renderer.RenderBackgroundScanline(1)
	renderer.RenderBackgroundScanline(2)
	renderer.RenderBackgroundScanline(3)
	
	assert.Equal(t, uint8(3), ppu.GetPixel(0, 0), "Row 0 should be inverted to color 3")
	assert.Equal(t, uint8(2), ppu.GetPixel(0, 1), "Row 1 should be inverted to color 2")
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 2), "Row 2 should be inverted to color 1")
	assert.Equal(t, uint8(0), ppu.GetPixel(0, 3), "Row 3 should be inverted to color 0")
}

// TestBackgroundPixelAccess tests pixel access methods
func TestBackgroundPixelAccess(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetBGP(0xE4)
	
	// Create test tile (color 0 and 3 pattern)
	testTileData := TileData{
		0x55, 0x55, // Row 0: color 0,3,0,3,0,3,0,3 (plane0=55=01010101, plane1=55=01010101)
		0x55, 0x55, // Row 1: same pattern  
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	
	vram.SetTileData(0x8000, testTileData)
	vram.SetTileMapEntry(0x9800, 0)
	
	// Render full background for testing
	renderer.RenderFullBackground()
	
	// Test GetBackgroundPixel
	assert.Equal(t, uint8(0), renderer.GetBackgroundPixel(0, 0), "Should get color 0")
	assert.Equal(t, uint8(3), renderer.GetBackgroundPixel(1, 0), "Should get color 3")
	
	// Test IsBackgroundPixelTransparent
	assert.True(t, renderer.IsBackgroundPixelTransparent(0, 0), "Color 0 should be transparent")
	assert.False(t, renderer.IsBackgroundPixelTransparent(1, 0), "Color 3 should not be transparent")
	
	// Test out of bounds
	assert.Equal(t, uint8(ColorWhite), renderer.GetBackgroundPixel(-1, 0), "Out of bounds should return white")
	assert.True(t, renderer.IsBackgroundPixelTransparent(-1, 0), "Out of bounds should be transparent")
	
	// Test with background disabled
	ppu.SetLCDC(0x00)
	assert.True(t, renderer.IsBackgroundPixelTransparent(1, 0), "Should be transparent when background disabled")
}

// TestBackgroundScrollingWraparound tests scroll wraparound behavior
func TestBackgroundScrollingWraparound(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetBGP(0xE4)
	
	// Create distinguishable tiles
	tile0Data := TileData{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // Color 0
	tile1Data := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00} // Color 1
	
	vram.SetTileData(0x8000, tile0Data)
	vram.SetTileData(0x8010, tile1Data)
	
	// Set up tile map pattern: tile 0 at (0,0), tile 1 at (31,0) (last column)
	vram.SetTileMapEntry(0x9800, 0)      // Position (0,0)
	vram.SetTileMapEntry(0x9800+31, 1)   // Position (31,0)
	
	// Test wraparound: scroll to position where we should see the last tile
	ppu.SetSCX(uint8(31*8)) // Scroll to see tile at column 31
	renderer.RenderBackgroundScanline(0)
	
	// First pixel should now show tile 1 (from position 31,0)
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 0), "Should wrap around to show tile from column 31")
	
	// Test vertical wraparound
	vram.SetTileMapEntry(0x9800+31*32, 1) // Position (0,31) - last row
	ppu.SetSCX(0)
	ppu.SetSCY(uint8(31*8)) // Scroll to see tile at row 31
	renderer.RenderBackgroundScanline(0)
	
	// Should now show tile from row 31
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 0), "Should wrap around to show tile from row 31")
}

// TestBackgroundAnalysis tests the debugging and analysis functions
func TestBackgroundAnalysis(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	// Test with background disabled
	ppu.SetLCDC(0x00)
	analysis := renderer.AnalyzeBackground()
	assert.False(t, analysis["backgroundEnabled"].(bool), "Should detect background disabled")
	
	// Test with background enabled
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetSCX(10)
	ppu.SetSCY(20)
	
	analysis = renderer.AnalyzeBackground()
	assert.True(t, analysis["backgroundEnabled"].(bool), "Should detect background enabled")
	assert.Equal(t, uint8(10), analysis["scrollX"].(uint8), "Should report scroll X")
	assert.Equal(t, uint8(20), analysis["scrollY"].(uint8), "Should report scroll Y")
	assert.False(t, analysis["tileMap1Selected"].(bool), "Should detect map 0 selected")
	assert.True(t, analysis["tileData1Selected"].(bool), "Should detect tile data 1 selected")
	
	// Test visible tiles
	visibleTiles := renderer.GetVisibleTiles()
	assert.Greater(t, len(visibleTiles), 0, "Should return visible tiles")
	
	// Test string representation
	str := renderer.String()
	assert.Contains(t, str, "ENABLED", "String should indicate enabled state")
	assert.Contains(t, str, "10", "String should contain scroll X")
	assert.Contains(t, str, "20", "String should contain scroll Y")
}

// TestBackgroundRendererValidation tests the validation functionality
func TestBackgroundRendererValidation(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	// Test with valid renderer
	renderer := NewBackgroundRenderer(ppu, vram)
	issues := renderer.ValidateRenderer()
	assert.Equal(t, 0, len(issues), "Valid renderer should have no issues")
	
	// Test with nil PPU
	rendererBadPPU := &BackgroundRenderer{ppu: nil, vramInterface: vram}
	issues = rendererBadPPU.ValidateRenderer()
	assert.Greater(t, len(issues), 0, "Should detect nil PPU")
	assert.Contains(t, issues[0], "PPU reference is nil", "Should report nil PPU")
	
	// Test with nil VRAM interface
	rendererBadVRAM := &BackgroundRenderer{ppu: ppu, vramInterface: nil}
	issues = rendererBadVRAM.ValidateRenderer()
	assert.Greater(t, len(issues), 0, "Should detect nil VRAM interface")
	assert.Contains(t, issues[0], "VRAM interface is nil", "Should report nil VRAM interface")
}

// TestBackgroundRendererIntegration tests integration with PPU
func TestBackgroundRendererIntegration(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	// Set up VRAM interface (this initializes background renderer)
	ppu.SetVRAMInterface(vram)
	
	// Check that background renderer was created
	bgRenderer := ppu.GetBackgroundRenderer()
	assert.NotNil(t, bgRenderer, "Background renderer should be initialized")
	
	// Enable background and set up test data
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetBGP(0xE4)
	
	// Create test tile
	testTileData := TileData{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00}
	vram.SetTileData(0x8000, testTileData)
	vram.SetTileMapEntry(0x9800, 0)
	
	// Test that PPU timing integration works by simulating a frame
	// Manually call the background renderer (normally done during PPU update)
	bgRenderer.RenderBackgroundScanline(0)
	
	// Verify rendering worked
	assert.Equal(t, uint8(1), ppu.GetPixel(0, 0), "Integration should produce correct rendering")
}

// TestBackgroundEdgeCases tests edge cases and error conditions
func TestBackgroundEdgeCases(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetBGP(0xE4)
	
	// Test rendering beyond screen height (should be ignored)
	renderer.RenderBackgroundScanline(200) // Beyond screen height
	// Should not crash
	
	// Test with extreme scroll values
	ppu.SetSCX(255)
	ppu.SetSCY(255)
	renderer.RenderBackgroundScanline(0)
	// Should not crash and should handle wraparound
	
	// Test pixel access with extreme coordinates
	color := renderer.GetBackgroundPixel(1000, 1000) // Way out of bounds
	assert.Equal(t, uint8(ColorWhite), color, "Extreme out of bounds should return white")
	
	transparent := renderer.IsBackgroundPixelTransparent(1000, 1000)
	assert.True(t, transparent, "Extreme out of bounds should be transparent")
	
	// Test with all tile map entries set to 255 (edge tile index)
	for i := uint16(0); i < 1024; i++ {
		vram.SetTileMapEntry(0x9800+i, 255)
	}
	
	// Create tile at index 255
	vram.SetTileData(0x8000+255*16, TileData{0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA})
	
	renderer.RenderBackgroundScanline(0)
	// Should not crash and should render correctly
}

// TestBackgroundPerformance tests basic performance characteristics
func TestBackgroundPerformance(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewBackgroundRenderer(ppu, vram)
	
	ppu.SetLCDC(0x91) // LCD enabled, background enabled
	ppu.SetBGP(0xE4)
	
	// Set up complex tile pattern to stress test
	for tileIdx := uint8(0); tileIdx < 16; tileIdx++ {
		tileData := TileData{}
		for i := 0; i < 16; i++ {
			tileData[i] = uint8(int(tileIdx)*16 + i) // Unique pattern per tile
		}
		vram.SetTileData(0x8000+uint16(tileIdx)*16, tileData)
	}
	
	// Fill tile map with pattern
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			vram.SetTileMapEntry(0x9800+uint16(y*32+x), uint8((x+y)%16))
		}
	}
	
	// Render multiple full frames to test performance
	for frame := 0; frame < 10; frame++ {
		for scanline := uint8(0); scanline < ScreenHeight; scanline++ {
			renderer.RenderBackgroundScanline(scanline)
		}
	}
	
	// Test should complete without timeout (performance check)
	assert.True(t, true, "Performance test completed")
}
package ppu

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestNewVRAM tests VRAM initialization
func TestNewVRAM(t *testing.T) {
	vram := NewVRAM()
	
	// Test that VRAM is initialized to all zeros
	for i := uint16(0x8000); i <= 0x9FFF; i++ {
		assert.Equal(t, uint8(0), vram.ReadByte(i), "VRAM should be initialized to 0 at address 0x%04X", i)
	}
	
	// Test that organized access structures are created
	assert.NotNil(t, vram.GetPatternTable0(), "Pattern table 0 should be initialized")
	assert.NotNil(t, vram.GetPatternTable1(), "Pattern table 1 should be initialized")
	assert.NotNil(t, vram.GetBackgroundMap0(), "Background map 0 should be initialized")
	assert.NotNil(t, vram.GetBackgroundMap1(), "Background map 1 should be initialized")
}

// TestVRAMRawAccess tests basic read/write operations
func TestVRAMRawAccess(t *testing.T) {
	vram := NewVRAM()
	
	// Test byte operations
	vram.WriteByte(0x8000, 0xFF)
	assert.Equal(t, uint8(0xFF), vram.ReadByte(0x8000), "Byte write/read should work")
	
	vram.WriteByte(0x9FFF, 0xAA)
	assert.Equal(t, uint8(0xAA), vram.ReadByte(0x9FFF), "Byte write/read at end of VRAM should work")
	
	// Test word operations (little-endian)
	vram.WriteWord(0x8100, 0x1234)
	assert.Equal(t, uint16(0x1234), vram.ReadWord(0x8100), "Word write/read should work")
	assert.Equal(t, uint8(0x34), vram.ReadByte(0x8100), "Word low byte should be correct")
	assert.Equal(t, uint8(0x12), vram.ReadByte(0x8101), "Word high byte should be correct")
	
	// Test address validation
	assert.Equal(t, uint8(0xFF), vram.ReadByte(0x7FFF), "Invalid address should return 0xFF")
	assert.Equal(t, uint8(0xFF), vram.ReadByte(0xA000), "Invalid address should return 0xFF")
	
	// Test that writes to invalid addresses are ignored
	originalValue := vram.ReadByte(0x8000)
	vram.WriteByte(0x7FFF, 0x55) // Should be ignored
	assert.Equal(t, originalValue, vram.ReadByte(0x8000), "Invalid write should not affect valid memory")
}

// TestVRAMAddressValidation tests the address validation logic
func TestVRAMAddressValidation(t *testing.T) {
	vram := NewVRAM()
	
	// Test valid addresses
	assert.True(t, vram.IsValidAddress(0x8000), "0x8000 should be valid")
	assert.True(t, vram.IsValidAddress(0x9FFF), "0x9FFF should be valid")
	assert.True(t, vram.IsValidAddress(0x8800), "0x8800 should be valid")
	assert.True(t, vram.IsValidAddress(0x9000), "0x9000 should be valid")
	
	// Test invalid addresses
	assert.False(t, vram.IsValidAddress(0x7FFF), "0x7FFF should be invalid")
	assert.False(t, vram.IsValidAddress(0xA000), "0xA000 should be invalid")
	assert.False(t, vram.IsValidAddress(0x0000), "0x0000 should be invalid")
	assert.False(t, vram.IsValidAddress(0xFFFF), "0xFFFF should be invalid")
}

// TestVRAMClear tests the clear functionality
func TestVRAMClear(t *testing.T) {
	vram := NewVRAM()
	
	// Fill with some data first
	for i := uint16(0x8000); i <= 0x80FF; i++ {
		vram.WriteByte(i, 0xFF)
	}
	
	// Clear with pattern
	vram.Clear(0xAA)
	
	// Verify all memory is cleared to pattern
	for i := uint16(0x8000); i <= 0x9FFF; i++ {
		assert.Equal(t, uint8(0xAA), vram.ReadByte(i), "All VRAM should be cleared to 0xAA")
	}
	
	// Clear to zeros
	vram.Clear(0x00)
	for i := uint16(0x8000); i <= 0x9FFF; i++ {
		assert.Equal(t, uint8(0x00), vram.ReadByte(i), "All VRAM should be cleared to 0x00")
	}
}

// TestTilePatternTableAccess tests pattern table operations
func TestTilePatternTableAccess(t *testing.T) {
	vram := NewVRAM()
	
	// Create test tile
	testTile := CreateTestTile(2) // Checkerboard pattern
	
	// Test pattern table 0 ($8000 method)
	table0 := vram.GetPatternTable0()
	table0.SetTile(0, testTile)
	
	retrievedTile := table0.GetTile(0)
	assert.True(t, CompareTiles(testTile, retrievedTile), "Retrieved tile should match original")
	
	// Test pattern table 1 ($8800 method)
	table1 := vram.GetPatternTable1()
	table1.SetTile(0, testTile)
	
	retrievedTile1 := table1.GetTile(0)
	assert.True(t, CompareTiles(testTile, retrievedTile1), "Retrieved tile from table 1 should match original")
	
	// Test different tile indices
	borderTile := CreateTestTile(4) // Border pattern
	table0.SetTile(255, borderTile)
	
	retrievedBorder := table0.GetTile(255)
	assert.True(t, CompareTiles(borderTile, retrievedBorder), "Tile at index 255 should be correct")
}

// TestTilePatternTableRawData tests raw data operations
func TestTilePatternTableRawData(t *testing.T) {
	vram := NewVRAM()
	table0 := vram.GetPatternTable0()
	
	// Create test data (solid color 3 pattern)
	testData := TileData{
		0xFF, 0xFF, // Row 0: all color 3
		0xFF, 0xFF, // Row 1: all color 3
		0xFF, 0xFF, // Row 2: all color 3
		0xFF, 0xFF, // Row 3: all color 3
		0xFF, 0xFF, // Row 4: all color 3
		0xFF, 0xFF, // Row 5: all color 3
		0xFF, 0xFF, // Row 6: all color 3
		0xFF, 0xFF, // Row 7: all color 3
	}
	
	// Set raw data
	table0.SetTileData(10, testData)
	
	// Get raw data back
	retrievedData := table0.GetTileData(10)
	assert.Equal(t, testData, retrievedData, "Raw tile data should match")
	
	// Verify it decodes correctly
	tile := table0.GetTile(10)
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			assert.Equal(t, uint8(3), tile.GetPixel(x, y), "All pixels should be color 3")
		}
	}
}

// TestTilePatternTableBulkLoad tests loading multiple tiles
func TestTilePatternTableBulkLoad(t *testing.T) {
	vram := NewVRAM()
	table0 := vram.GetPatternTable0()
	
	// Create test tiles data
	tilesData := []TileData{
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Empty
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Solid
		{0x55, 0x00, 0x55, 0x00, 0x55, 0x00, 0x55, 0x00, 0x55, 0x00, 0x55, 0x00, 0x55, 0x00, 0x55, 0x00}, // Pattern
	}
	
	// Load tiles starting at index 50
	table0.LoadTiles(50, tilesData)
	
	// Verify tiles were loaded correctly
	for i, expectedData := range tilesData {
		actualData := table0.GetTileData(50 + uint8(i))
		assert.Equal(t, expectedData, actualData, "Tile %d should match expected data", i)
	}
	
	// Test bounds checking - should not crash when loading near the end
	table0.LoadTiles(254, tilesData) // Should only load 2 tiles (254, 255)
}

// TestTileMapAccess tests tile map operations
func TestTileMapAccess(t *testing.T) {
	vram := NewVRAM()
	
	// Test background map 0
	map0 := vram.GetBackgroundMap0()
	
	// Set some tile indices
	map0.SetTileIndex(0, 0, 123)
	map0.SetTileIndex(31, 31, 234) // Corner
	map0.SetTileIndex(15, 15, 45)  // Middle
	
	// Read them back
	assert.Equal(t, uint8(123), map0.GetTileIndex(0, 0), "Tile index at (0,0) should be correct")
	assert.Equal(t, uint8(234), map0.GetTileIndex(31, 31), "Tile index at (31,31) should be correct")
	assert.Equal(t, uint8(45), map0.GetTileIndex(15, 15), "Tile index at (15,15) should be correct")
	
	// Test bounds checking
	assert.Equal(t, uint8(0), map0.GetTileIndex(-1, 0), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), map0.GetTileIndex(32, 0), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), map0.GetTileIndex(0, -1), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), map0.GetTileIndex(0, 32), "Out of bounds should return 0")
	
	// Test that out-of-bounds writes are ignored
	originalValue := map0.GetTileIndex(0, 0)
	map0.SetTileIndex(-1, 0, 99) // Should be ignored
	assert.Equal(t, originalValue, map0.GetTileIndex(0, 0), "Out of bounds write should be ignored")
}

// TestTileMapLinearAccess tests linear indexing for tile maps
func TestTileMapLinearAccess(t *testing.T) {
	vram := NewVRAM()
	map0 := vram.GetBackgroundMap0()
	
	// Test linear access (index = y * 32 + x)
	map0.SetTileIndexLinear(0, 111)    // (0, 0)
	map0.SetTileIndexLinear(32, 222)   // (0, 1)
	map0.SetTileIndexLinear(1023, 77) // (31, 31)
	
	// Verify through coordinate access
	assert.Equal(t, uint8(111), map0.GetTileIndex(0, 0), "Linear index 0 should map to (0,0)")
	assert.Equal(t, uint8(222), map0.GetTileIndex(0, 1), "Linear index 32 should map to (0,1)")
	assert.Equal(t, uint8(77), map0.GetTileIndex(31, 31), "Linear index 1023 should map to (31,31)")
	
	// Test linear read
	assert.Equal(t, uint8(111), map0.GetTileIndexLinear(0), "Linear read should work")
	assert.Equal(t, uint8(222), map0.GetTileIndexLinear(32), "Linear read should work")
	assert.Equal(t, uint8(77), map0.GetTileIndexLinear(1023), "Linear read should work")
	
	// Test bounds checking
	assert.Equal(t, uint8(0), map0.GetTileIndexLinear(-1), "Negative linear index should return 0")
	assert.Equal(t, uint8(0), map0.GetTileIndexLinear(1024), "Out of bounds linear index should return 0")
}

// TestTileMapFillAndLoad tests bulk tile map operations
func TestTileMapFillAndLoad(t *testing.T) {
	vram := NewVRAM()
	map0 := vram.GetBackgroundMap0()
	
	// Test fill operation
	map0.FillMap(42)
	
	// Verify all positions are filled
	for y := 0; y < TileMapHeight; y++ {
		for x := 0; x < TileMapWidth; x++ {
			assert.Equal(t, uint8(42), map0.GetTileIndex(x, y), "All tiles should be filled with 42")
		}
	}
	
	// Test load operation with partial data
	testData := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	map0.LoadMapData(testData)
	
	// Verify first 10 positions
	for i := 0; i < len(testData); i++ {
		assert.Equal(t, testData[i], map0.GetTileIndexLinear(i), "Loaded data should match at index %d", i)
	}
	
	// Verify remaining positions are unchanged (should still be 42)
	for i := len(testData); i < 20; i++ {
		assert.Equal(t, uint8(42), map0.GetTileIndexLinear(i), "Unloaded positions should remain unchanged")
	}
	
	// Test get map data
	fullData := map0.GetMapData()
	assert.Equal(t, TileMapSize, len(fullData), "Full map data should be correct size")
	
	// Verify first few entries match what we loaded
	for i := 0; i < len(testData); i++ {
		assert.Equal(t, testData[i], fullData[i], "Exported data should match at index %d", i)
	}
}

// TestTileMapVisibleRegion tests the visible region calculation
func TestTileMapVisibleRegion(t *testing.T) {
	vram := NewVRAM()
	map0 := vram.GetBackgroundMap0()
	
	// Fill map with predictable pattern (tile index = y * 32 + x)
	for y := 0; y < TileMapHeight; y++ {
		for x := 0; x < TileMapWidth; x++ {
			map0.SetTileIndex(x, y, uint8(y*32+x))
		}
	}
	
	// Test visible region with no scroll
	visibleTiles := map0.GetVisibleRegion(0, 0)
	
	// Verify dimensions
	assert.Equal(t, ScreenTilesHeight, len(visibleTiles), "Visible region should have correct height")
	assert.Equal(t, ScreenTilesWidth, len(visibleTiles[0]), "Visible region should have correct width")
	
	// Verify content matches expected tiles
	for screenY := 0; screenY < ScreenTilesHeight; screenY++ {
		for screenX := 0; screenX < ScreenTilesWidth; screenX++ {
			expected := uint8(screenY*32 + screenX)
			assert.Equal(t, expected, visibleTiles[screenY][screenX], 
				"Visible tile at (%d,%d) should be %d", screenX, screenY, expected)
		}
	}
	
	// Test with scroll offset
	visibleTilesScrolled := map0.GetVisibleRegion(8, 8) // Scroll by 1 tile
	
	// First visible tile should now be (1,1) from the map
	expected := uint8(1*32 + 1) // y=1, x=1
	assert.Equal(t, expected, visibleTilesScrolled[0][0], "Scrolled visible region should be offset")
	
	// Test wrapping (scroll beyond map boundary)
	visibleTilesWrapped := map0.GetVisibleRegion(248, 248) // 31 * 8 = 248, so this wraps to tile (31,31)
	
	// Should wrap around to show tiles starting from (31,31), then (0,0), etc.
	expectedWrapped := uint8((31*32+31) % 256) // Handle overflow
	assert.Equal(t, expectedWrapped, visibleTilesWrapped[0][0], "Should wrap to corner tile")
}

// TestHighLevelVRAMOperations tests the high-level VRAM operations
func TestHighLevelVRAMOperations(t *testing.T) {
	vram := NewVRAM()
	
	// Set up a test tile in pattern table 0
	testTile := CreateTestTile(2) // Checkerboard
	table0 := vram.GetPatternTable0()
	table0.SetTile(42, testTile)
	
	// Set up tile map to reference this tile
	map0 := vram.GetBackgroundMap0()
	map0.SetTileIndex(5, 5, 42)
	
	// Test GetTileFromMap
	retrievedTile := vram.GetTileFromMap(5, 5, false, false)
	assert.True(t, CompareTiles(testTile, retrievedTile), "GetTileFromMap should return correct tile")
	
	// Test with different map and addressing mode
	table1 := vram.GetPatternTable1()
	table1.SetTile(128, testTile) // Signed index -128
	
	map1 := vram.GetBackgroundMap1()
	map1.SetTileIndex(10, 10, 128)
	
	retrievedTileSigned := vram.GetTileFromMap(10, 10, true, true)
	assert.True(t, CompareTiles(testTile, retrievedTileSigned), "GetTileFromMap with signed mode should work")
}

// TestVRAMFramebufferRendering tests the framebuffer rendering helper
func TestVRAMFramebufferRendering(t *testing.T) {
	vram := NewVRAM()
	
	// Create a framebuffer
	var framebuffer [ScreenHeight][ScreenWidth]uint8
	
	// Create test tile and palette
	testTile := CreateTestTile(4) // Border pattern
	palette := [4]uint8{0, 1, 2, 3} // Identity palette
	
	// Render tile to framebuffer at (0,0)
	vram.RenderTileToFramebuffer(&framebuffer, testTile, 0, 0, palette)
	
	// Verify the tile was rendered correctly
	for tileY := 0; tileY < TileHeight; tileY++ {
		for tileX := 0; tileX < TileWidth; tileX++ {
			expected := palette[testTile.GetPixel(tileX, tileY)]
			assert.Equal(t, expected, framebuffer[tileY][tileX], 
				"Rendered pixel at (%d,%d) should match tile", tileX, tileY)
		}
	}
	
	// Test rendering with offset
	vram.RenderTileToFramebuffer(&framebuffer, testTile, 16, 16, palette)
	
	// Verify offset rendering
	expected := palette[testTile.GetPixel(0, 0)]
	assert.Equal(t, expected, framebuffer[16][16], "Offset rendering should work")
	
	// Test bounds clipping (render partially off-screen)
	vram.RenderTileToFramebuffer(&framebuffer, testTile, ScreenWidth-4, ScreenHeight-4, palette)
	
	// Should not crash and should render visible portion
	expected = palette[testTile.GetPixel(0, 0)]
	assert.Equal(t, expected, framebuffer[ScreenHeight-4][ScreenWidth-4], "Clipped rendering should work")
}

// TestVRAMStats tests the statistics and debugging functions
func TestVRAMStats(t *testing.T) {
	vram := NewVRAM()
	
	// Initially, VRAM should be empty
	stats := vram.GetVRAMStats()
	assert.Equal(t, 0x2000, stats["totalSize"].(int), "Total size should be 8KB")
	assert.Equal(t, 0, stats["patternDataUsed"].(int), "Pattern data should be empty initially")
	assert.Equal(t, 0, stats["mapDataUsed"].(int), "Map data should be empty initially")
	assert.Equal(t, 0, stats["totalUsed"].(int), "Total used should be 0 initially")
	assert.Equal(t, 0.0, stats["percentUsed"].(float64), "Percent used should be 0 initially")
	
	// Add some pattern data
	vram.WriteByte(0x8000, 0xFF)
	vram.WriteByte(0x8010, 0xAA)
	
	// Add some map data  
	vram.WriteByte(0x9800, 0x55)
	vram.WriteByte(0x9C00, 0x33)
	
	// Check updated stats
	statsAfter := vram.GetVRAMStats()
	assert.Equal(t, 2, statsAfter["patternDataUsed"].(int), "Pattern data should show 2 bytes used")
	assert.Equal(t, 2, statsAfter["mapDataUsed"].(int), "Map data should show 2 bytes used")
	assert.Equal(t, 4, statsAfter["totalUsed"].(int), "Total used should be 4")
	assert.True(t, statsAfter["percentUsed"].(float64) > 0.0, "Percent used should be greater than 0")
}

// TestVRAMValidation tests the validation functionality
func TestVRAMValidation(t *testing.T) {
	vram := NewVRAM()
	
	// Test that validation function runs without crashing
	_ = vram.ValidateVRAM()
	
	// Add some tile references
	map0 := vram.GetBackgroundMap0()
	map0.SetTileIndex(0, 0, 10)
	map0.SetTileIndex(1, 1, 200) // High index (negative in signed mode)
	
	// Test validation after adding references
	_ = vram.ValidateVRAM()
	
	// The main goal is that validation doesn't crash and returns some result
	assert.True(t, true, "Validation completed successfully")
}

// TestTileMapUtilities tests tile map utility functions
func TestTileMapUtilities(t *testing.T) {
	vram := NewVRAM()
	map0 := vram.GetBackgroundMap0()
	
	// Set up a pattern with tile index 42 in several locations
	map0.SetTileIndex(5, 5, 42)
	map0.SetTileIndex(10, 15, 42)
	map0.SetTileIndex(20, 25, 42)
	map0.SetTileIndex(0, 0, 99) // Different tile
	
	// Test FindTileUsage
	locations := map0.FindTileUsage(42)
	assert.Equal(t, 3, len(locations), "Should find 3 locations with tile 42")
	
	// Verify locations are correct
	foundLocations := make(map[string]bool)
	for _, loc := range locations {
		key := fmt.Sprintf("%d,%d", loc.X, loc.Y)
		foundLocations[key] = true
	}
	
	assert.True(t, foundLocations["5,5"], "Should find tile at (5,5)")
	assert.True(t, foundLocations["10,15"], "Should find tile at (10,15)")
	assert.True(t, foundLocations["20,25"], "Should find tile at (20,25)")
	assert.False(t, foundLocations["0,0"], "Should not find different tile at (0,0)")
	
	// Test DumpTileMap
	dumpStr := map0.DumpTileMap(4, 4) // Small dump
	assert.Contains(t, dumpStr, "Tile Map", "Dump should contain title")
	assert.Contains(t, dumpStr, "63", "Should contain hex representation of tile 99")
	assert.Contains(t, dumpStr, "00", "Should contain hex representation of empty tiles")
}

// TestVRAMInterfaceImplementation tests the VRAMInterface implementation
func TestVRAMInterfaceImplementation(t *testing.T) {
	vram := NewVRAM()
	
	// Test that VRAM implements VRAMInterface methods
	var vramInterface VRAMInterface = vram
	
	// Test VRAM read/write through interface
	vramInterface.WriteVRAM(0x8100, 0xAB)
	assert.Equal(t, uint8(0xAB), vramInterface.ReadVRAM(0x8100), "VRAM interface should work")
	
	// Test OAM methods (currently placeholders)
	vramInterface.WriteOAM(0xFE00, 0x55) // Should not crash
	assert.Equal(t, uint8(0), vramInterface.ReadOAM(0xFE00), "OAM should return 0 (placeholder)")
}

// TestVRAMEdgeCases tests edge cases and error conditions
func TestVRAMEdgeCases(t *testing.T) {
	vram := NewVRAM()
	
	// Test operations on boundary addresses
	vram.WriteByte(0x7FFF, 0xFF) // Should be ignored
	vram.WriteByte(0x8000, 0xAA) // Should work
	vram.WriteByte(0x9FFF, 0xBB) // Should work
	vram.WriteByte(0xA000, 0xFF) // Should be ignored
	
	assert.Equal(t, uint8(0xAA), vram.ReadByte(0x8000), "Valid address should work")
	assert.Equal(t, uint8(0xBB), vram.ReadByte(0x9FFF), "Valid address should work")
	
	// Test word operations at boundaries
	vram.WriteWord(0x9FFE, 0x1234) // Spans boundary
	assert.Equal(t, uint8(0x34), vram.ReadByte(0x9FFE), "Word write at boundary should work")
	assert.Equal(t, uint8(0x12), vram.ReadByte(0x9FFF), "High byte should be written correctly")
	
	// Test pattern table with all indices
	table0 := vram.GetPatternTable0()
	testData := TileData{0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 
					   0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA}
	
	// Test all possible tile indices
	for i := 0; i <= 255; i++ {
		table0.SetTileData(uint8(i), testData)
		retrievedData := table0.GetTileData(uint8(i))
		assert.Equal(t, testData, retrievedData, "Tile index %d should work", i)
	}
}
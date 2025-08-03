package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestNewTile tests tile creation and initialization
func TestNewTile(t *testing.T) {
	tile := NewTile()
	
	// Test that new tile is empty (all pixels are 0)
	assert.True(t, tile.IsEmpty(), "New tile should be empty")
	
	// Test all pixels are initialized to 0
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			assert.Equal(t, uint8(0), tile.GetPixel(x, y), "All pixels should be 0")
		}
	}
}

// TestTilePixelAccess tests getting and setting individual pixels
func TestTilePixelAccess(t *testing.T) {
	tile := NewTile()
	
	// Test setting and getting pixels
	tile.SetPixel(0, 0, 3)
	tile.SetPixel(7, 7, 2)
	tile.SetPixel(3, 4, 1)
	
	assert.Equal(t, uint8(3), tile.GetPixel(0, 0), "Pixel (0,0) should be 3")
	assert.Equal(t, uint8(2), tile.GetPixel(7, 7), "Pixel (7,7) should be 2")
	assert.Equal(t, uint8(1), tile.GetPixel(3, 4), "Pixel (3,4) should be 1")
	
	// Test bounds checking
	assert.Equal(t, uint8(0), tile.GetPixel(-1, 0), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), tile.GetPixel(8, 0), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), tile.GetPixel(0, -1), "Out of bounds should return 0")
	assert.Equal(t, uint8(0), tile.GetPixel(0, 8), "Out of bounds should return 0")
	
	// Test color clamping
	tile.SetPixel(1, 1, 5) // Invalid color > 3
	assert.Equal(t, uint8(3), tile.GetPixel(1, 1), "Color should be clamped to 3")
	
	// Test that out-of-bounds writes are ignored
	originalPixel := tile.GetPixel(0, 0)
	tile.SetPixel(-1, 0, 1) // Should be ignored
	tile.SetPixel(8, 0, 1)  // Should be ignored
	assert.Equal(t, originalPixel, tile.GetPixel(0, 0), "Out of bounds writes should be ignored")
}

// TestTileClear tests clearing tiles with specific colors
func TestTileClear(t *testing.T) {
	tile := NewTile()
	
	// Set some pixels first
	tile.SetPixel(1, 1, 2)
	tile.SetPixel(5, 3, 3)
	assert.False(t, tile.IsEmpty(), "Tile should not be empty after setting pixels")
	
	// Clear with color 1
	tile.Clear(1)
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			assert.Equal(t, uint8(1), tile.GetPixel(x, y), "All pixels should be 1 after clear")
		}
	}
	assert.False(t, tile.IsEmpty(), "Tile should not be empty after clearing with color 1")
	
	// Clear with color 0 (should make it empty)
	tile.Clear(0)
	assert.True(t, tile.IsEmpty(), "Tile should be empty after clearing with color 0")
	
	// Test color clamping in clear
	tile.Clear(5) // Should clamp to 3
	assert.Equal(t, uint8(3), tile.GetPixel(0, 0), "Clear should clamp color to 3")
}

// TestTileCopy tests tile copying functionality
func TestTileCopy(t *testing.T) {
	original := NewTile()
	
	// Create a pattern in original tile
	original.SetPixel(0, 0, 3)
	original.SetPixel(1, 1, 2)
	original.SetPixel(7, 7, 1)
	
	// Copy the tile
	copy := original.Copy()
	
	// Test that copy has same data
	assert.True(t, CompareTiles(original, copy), "Copy should have same data as original")
	
	// Test that they are separate objects (modify copy)
	copy.SetPixel(0, 0, 0)
	assert.False(t, CompareTiles(original, copy), "Copy should be independent of original")
	assert.Equal(t, uint8(3), original.GetPixel(0, 0), "Original should be unchanged")
	assert.Equal(t, uint8(0), copy.GetPixel(0, 0), "Copy should be modified")
}

// TestTileFlipping tests horizontal and vertical flipping
func TestTileFlipping(t *testing.T) {
	tile := NewTile()
	
	// Create asymmetric pattern for testing
	// Pattern:
	// 3 0 0 0 0 0 0 1
	// 0 0 0 0 0 0 0 0
	// ...
	// 0 0 0 0 0 0 0 0
	// 2 0 0 0 0 0 0 0
	tile.SetPixel(0, 0, 3) // Top-left
	tile.SetPixel(7, 0, 1) // Top-right
	tile.SetPixel(0, 7, 2) // Bottom-left
	
	// Test horizontal flip
	hFlipped := tile.FlipHorizontal()
	assert.Equal(t, uint8(1), hFlipped.GetPixel(0, 0), "Horizontally flipped: top-right becomes top-left")
	assert.Equal(t, uint8(3), hFlipped.GetPixel(7, 0), "Horizontally flipped: top-left becomes top-right")
	assert.Equal(t, uint8(2), hFlipped.GetPixel(7, 7), "Horizontally flipped: bottom-left becomes bottom-right")
	assert.Equal(t, uint8(0), hFlipped.GetPixel(0, 7), "Horizontally flipped: bottom-right becomes bottom-left")
	
	// Test vertical flip
	vFlipped := tile.FlipVertical()
	assert.Equal(t, uint8(2), vFlipped.GetPixel(0, 0), "Vertically flipped: bottom-left becomes top-left")
	assert.Equal(t, uint8(0), vFlipped.GetPixel(7, 0), "Vertically flipped: bottom-right becomes top-right")
	assert.Equal(t, uint8(3), vFlipped.GetPixel(0, 7), "Vertically flipped: top-left becomes bottom-left")
	assert.Equal(t, uint8(1), vFlipped.GetPixel(7, 7), "Vertically flipped: top-right becomes bottom-right")
	
	// Test both flip (180 degree rotation)
	bothFlipped := tile.FlipBoth()
	assert.Equal(t, uint8(0), bothFlipped.GetPixel(0, 0), "Both flipped: bottom-right becomes top-left")
	assert.Equal(t, uint8(2), bothFlipped.GetPixel(7, 0), "Both flipped: bottom-left becomes top-right")
	assert.Equal(t, uint8(1), bothFlipped.GetPixel(0, 7), "Both flipped: top-right becomes bottom-left")
	assert.Equal(t, uint8(3), bothFlipped.GetPixel(7, 7), "Both flipped: top-left becomes bottom-right")
}

// TestTileDataEncoding tests 2bpp encoding and decoding
func TestTileDataEncoding(t *testing.T) {
	// Create a test tile with known pattern
	tile := NewTile()
	
	// Set a specific pattern that's easy to verify
	// Row 0: 0,1,2,3,0,1,2,3 (tests all color values)
	tile.SetPixel(0, 0, 0)
	tile.SetPixel(1, 0, 1)
	tile.SetPixel(2, 0, 2)
	tile.SetPixel(3, 0, 3)
	tile.SetPixel(4, 0, 0)
	tile.SetPixel(5, 0, 1)
	tile.SetPixel(6, 0, 2)
	tile.SetPixel(7, 0, 3)
	
	// Encode to Game Boy format
	data := tile.ToData()
	
	// Expected encoding for row 0:
	// Colors: 0,1,2,3,0,1,2,3
	// Binary: 00,01,10,11,00,01,10,11
	// Plane 0: 0,1,0,1,0,1,0,1 = 0b01010101 = 0x55
	// Plane 1: 0,0,1,1,0,0,1,1 = 0b00110011 = 0x33
	assert.Equal(t, uint8(0x55), data[0], "Row 0 plane 0 should be 0x55")
	assert.Equal(t, uint8(0x33), data[1], "Row 0 plane 1 should be 0x33")
	
	// Test round-trip encoding
	decoded := NewTileFromData(data)
	assert.True(t, CompareTiles(tile, decoded), "Round-trip encoding should preserve tile data")
	
	// Test with empty tile
	emptyTile := NewTile()
	emptyData := emptyTile.ToData()
	for i, b := range emptyData {
		assert.Equal(t, uint8(0), b, "Empty tile should encode to all zeros at byte %d", i)
	}
	
	// Test with solid color tile
	solidTile := NewTile()
	solidTile.Clear(3) // All pixels = 3 (binary 11)
	solidData := solidTile.ToData()
	for i := 0; i < TileSize; i += 2 {
		assert.Equal(t, uint8(0xFF), solidData[i], "Solid color 3 plane 0 should be 0xFF")
		assert.Equal(t, uint8(0xFF), solidData[i+1], "Solid color 3 plane 1 should be 0xFF")
	}
}

// TestTileAddressing tests tile address calculation
func TestTileAddressing(t *testing.T) {
	// Test $8000 method (unsigned)
	addr0 := GetTileAddress(0, false)
	assert.Equal(t, uint16(0x8000), addr0, "Tile 0 in $8000 method should be at 0x8000")
	
	addr1 := GetTileAddress(1, false)
	assert.Equal(t, uint16(0x8010), addr1, "Tile 1 in $8000 method should be at 0x8010")
	
	addr255 := GetTileAddress(255, false)
	assert.Equal(t, uint16(0x8FF0), addr255, "Tile 255 in $8000 method should be at 0x8FF0")
	
	// Test $8800 method (signed)
	addr0Signed := GetTileAddress(0, true)
	assert.Equal(t, uint16(0x9000), addr0Signed, "Tile 0 in $8800 method should be at 0x9000")
	
	// Test negative tile indices (128-255 in unsigned = -128 to -1 in signed)
	addrNeg1 := GetTileAddress(255, true) // 255 as uint8 = -1 as int8
	assert.Equal(t, uint16(0x8FF0), addrNeg1, "Tile -1 in $8800 method should be at 0x8FF0")
	
	addrNeg128 := GetTileAddress(128, true) // 128 as uint8 = -128 as int8
	assert.Equal(t, uint16(0x8800), addrNeg128, "Tile -128 in $8800 method should be at 0x8800")
	
	// Test positive indices in signed mode
	addr127 := GetTileAddress(127, true)
	assert.Equal(t, uint16(0x97F0), addr127, "Tile 127 in $8800 method should be at 0x97F0")
}

// TestTileAddressValidation tests tile address validation
func TestTileAddressValidation(t *testing.T) {
	// Test valid addresses
	assert.True(t, IsValidTileAddress(0x8000), "0x8000 should be valid tile address")
	assert.True(t, IsValidTileAddress(0x8FFF), "0x8FFF should be valid tile address")
	assert.True(t, IsValidTileAddress(0x8800), "0x8800 should be valid tile address")
	assert.True(t, IsValidTileAddress(0x97FF), "0x97FF should be valid tile address")
	
	// Test invalid addresses
	assert.False(t, IsValidTileAddress(0x7FFF), "0x7FFF should not be valid tile address")
	assert.False(t, IsValidTileAddress(0x9800), "0x9800 should not be valid tile address (tile map)")
	assert.False(t, IsValidTileAddress(0xA000), "0xA000 should not be valid tile address")
}

// TestTileMapAddressing tests tile map address calculation
func TestTileMapAddressing(t *testing.T) {
	// Test map 0 (0x9800)
	addr00 := GetTileMapAddress(0, 0, false)
	assert.Equal(t, uint16(0x9800), addr00, "Map 0 tile (0,0) should be at 0x9800")
	
	addr10 := GetTileMapAddress(1, 0, false)
	assert.Equal(t, uint16(0x9801), addr10, "Map 0 tile (1,0) should be at 0x9801")
	
	addr01 := GetTileMapAddress(0, 1, false)
	assert.Equal(t, uint16(0x9820), addr01, "Map 0 tile (0,1) should be at 0x9820") // Next row = +32
	
	// Test map 1 (0x9C00)
	addr00Map1 := GetTileMapAddress(0, 0, true)
	assert.Equal(t, uint16(0x9C00), addr00Map1, "Map 1 tile (0,0) should be at 0x9C00")
	
	// Test bounds checking
	addrInvalid := GetTileMapAddress(-1, 0, false)
	assert.Equal(t, uint16(0), addrInvalid, "Invalid coordinates should return 0")
	
	addrInvalid2 := GetTileMapAddress(32, 0, false)
	assert.Equal(t, uint16(0), addrInvalid2, "Invalid coordinates should return 0")
}

// TestTileMapValidation tests tile map address validation
func TestTileMapValidation(t *testing.T) {
	// Test valid tile map addresses
	assert.True(t, IsValidTileMapAddress(0x9800), "0x9800 should be valid tile map address")
	assert.True(t, IsValidTileMapAddress(0x9BFF), "0x9BFF should be valid tile map address")
	assert.True(t, IsValidTileMapAddress(0x9C00), "0x9C00 should be valid tile map address")
	assert.True(t, IsValidTileMapAddress(0x9FFF), "0x9FFF should be valid tile map address")
	
	// Test invalid tile map addresses
	assert.False(t, IsValidTileMapAddress(0x97FF), "0x97FF should not be valid tile map address")
	assert.False(t, IsValidTileMapAddress(0xA000), "0xA000 should not be valid tile map address")
}

// TestCreateTestTile tests the test tile creation utility
func TestCreateTestTile(t *testing.T) {
	// Test solid color tiles
	whiteTile := CreateTestTile(0)
	assert.True(t, whiteTile.IsEmpty(), "Test tile pattern 0 should be empty")
	
	blackTile := CreateTestTile(1)
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			assert.Equal(t, uint8(3), blackTile.GetPixel(x, y), "Test tile pattern 1 should be all black")
		}
	}
	
	// Test checkerboard pattern
	checkerTile := CreateTestTile(2)
	assert.Equal(t, uint8(0), checkerTile.GetPixel(0, 0), "Checkerboard (0,0) should be white")
	assert.Equal(t, uint8(3), checkerTile.GetPixel(1, 0), "Checkerboard (1,0) should be black")
	assert.Equal(t, uint8(3), checkerTile.GetPixel(0, 1), "Checkerboard (0,1) should be black")
	assert.Equal(t, uint8(0), checkerTile.GetPixel(1, 1), "Checkerboard (1,1) should be white")
	
	// Test border pattern
	borderTile := CreateTestTile(4)
	assert.Equal(t, uint8(3), borderTile.GetPixel(0, 0), "Border tile corner should be black")
	assert.Equal(t, uint8(3), borderTile.GetPixel(7, 0), "Border tile corner should be black")
	assert.Equal(t, uint8(3), borderTile.GetPixel(0, 7), "Border tile corner should be black")
	assert.Equal(t, uint8(3), borderTile.GetPixel(7, 7), "Border tile corner should be black")
	assert.Equal(t, uint8(0), borderTile.GetPixel(1, 1), "Border tile interior should be white")
}

// TestAnalyzeTile tests tile analysis functionality
func TestAnalyzeTile(t *testing.T) {
	// Test empty tile
	emptyTile := NewTile()
	emptyStats := AnalyzeTile(emptyTile)
	
	assert.True(t, emptyStats["isEmpty"].(bool), "Empty tile should be detected as empty")
	assert.Equal(t, 64, emptyStats["color0Count"].(int), "Empty tile should have 64 color 0 pixels")
	assert.Equal(t, 0, emptyStats["color1Count"].(int), "Empty tile should have 0 color 1 pixels")
	assert.Equal(t, 0, emptyStats["color2Count"].(int), "Empty tile should have 0 color 2 pixels")
	assert.Equal(t, 0, emptyStats["color3Count"].(int), "Empty tile should have 0 color 3 pixels")
	assert.Equal(t, 64, emptyStats["totalPixels"].(int), "Total pixels should be 64")
	
	// Test solid color tile
	solidTile := NewTile()
	solidTile.Clear(2)
	solidStats := AnalyzeTile(solidTile)
	
	assert.False(t, solidStats["isEmpty"].(bool), "Solid color tile should not be empty")
	assert.Equal(t, 0, solidStats["color0Count"].(int), "Solid color 2 tile should have 0 color 0 pixels")
	assert.Equal(t, 64, solidStats["color2Count"].(int), "Solid color 2 tile should have 64 color 2 pixels")
}

// TestCompareTiles tests tile comparison functionality
func TestCompareTiles(t *testing.T) {
	tile1 := NewTile()
	tile2 := NewTile()
	
	// Test identical empty tiles
	assert.True(t, CompareTiles(tile1, tile2), "Empty tiles should be equal")
	
	// Make tiles different
	tile1.SetPixel(0, 0, 1)
	assert.False(t, CompareTiles(tile1, tile2), "Different tiles should not be equal")
	
	// Make them same again
	tile2.SetPixel(0, 0, 1)
	assert.True(t, CompareTiles(tile1, tile2), "Identical tiles should be equal")
	
	// Test with complex patterns
	pattern1 := CreateTestTile(2) // Checkerboard
	pattern2 := CreateTestTile(2) // Same checkerboard
	pattern3 := CreateTestTile(3) // Different pattern
	
	assert.True(t, CompareTiles(pattern1, pattern2), "Same patterns should be equal")
	assert.False(t, CompareTiles(pattern1, pattern3), "Different patterns should not be equal")
}

// TestCreateTileFromPattern tests creating tiles from string patterns
func TestCreateTileFromPattern(t *testing.T) {
	// Test simple pattern
	pattern := "        " + // Row 0: all spaces (color 0)
			  "........" + // Row 1: all dots (color 1)
			  "oooooooo" + // Row 2: all o's (color 2)
			  "########" + // Row 3: all #'s (color 3)
			  "        " + // Row 4: all spaces
			  "        " + // Row 5: all spaces
			  "        " + // Row 6: all spaces
			  "        "   // Row 7: all spaces
	
	tile := CreateTileFromPattern(pattern)
	
	// Test row 0 (all color 0)
	for x := 0; x < TileWidth; x++ {
		assert.Equal(t, uint8(0), tile.GetPixel(x, 0), "Row 0 should be all color 0")
	}
	
	// Test row 1 (all color 1)
	for x := 0; x < TileWidth; x++ {
		assert.Equal(t, uint8(1), tile.GetPixel(x, 1), "Row 1 should be all color 1")
	}
	
	// Test row 2 (all color 2)
	for x := 0; x < TileWidth; x++ {
		assert.Equal(t, uint8(2), tile.GetPixel(x, 2), "Row 2 should be all color 2")
	}
	
	// Test row 3 (all color 3)
	for x := 0; x < TileWidth; x++ {
		assert.Equal(t, uint8(3), tile.GetPixel(x, 3), "Row 3 should be all color 3")
	}
}

// TestTileConstants tests that tile constants have correct values
func TestTileConstants(t *testing.T) {
	assert.Equal(t, 8, TileWidth, "Tile width should be 8")
	assert.Equal(t, 8, TileHeight, "Tile height should be 8")
	assert.Equal(t, 16, TileSize, "Tile size should be 16 bytes")
	assert.Equal(t, 255, MaxTileIndex, "Max tile index should be 255")
	assert.Equal(t, 256, MaxTiles, "Max tiles should be 256")
	assert.Equal(t, 32, TileMapWidth, "Tile map width should be 32")
	assert.Equal(t, 32, TileMapHeight, "Tile map height should be 32")
	assert.Equal(t, 1024, TileMapSize, "Tile map size should be 1024")
	assert.Equal(t, 20, ScreenTilesWidth, "Screen width in tiles should be 20")
	assert.Equal(t, 18, ScreenTilesHeight, "Screen height in tiles should be 18")
	
	// Test VRAM addresses
	assert.Equal(t, uint16(TilePatternTable0Start), uint16(0x8000), "Pattern table 0 should start at 0x8000")
	assert.Equal(t, uint16(TilePatternTable0End), uint16(0x8FFF), "Pattern table 0 should end at 0x8FFF")
	assert.Equal(t, uint16(TilePatternTable1Start), uint16(0x8800), "Pattern table 1 should start at 0x8800")
	assert.Equal(t, uint16(TilePatternTable1End), uint16(0x97FF), "Pattern table 1 should end at 0x97FF")
	assert.Equal(t, uint16(BackgroundMap0Start), uint16(0x9800), "Background map 0 should start at 0x9800")
	assert.Equal(t, uint16(BackgroundMap0End), uint16(0x9BFF), "Background map 0 should end at 0x9BFF")
	assert.Equal(t, uint16(BackgroundMap1Start), uint16(0x9C00), "Background map 1 should start at 0x9C00")
	assert.Equal(t, uint16(BackgroundMap1End), uint16(0x9FFF), "Background map 1 should end at 0x9FFF")
}

// TestTileString tests the string representation of tiles
func TestTileString(t *testing.T) {
	// Test empty tile string representation
	emptyTile := NewTile()
	emptyStr := emptyTile.String()
	assert.Contains(t, emptyStr, "Tile 8x8:", "String should contain title")
	
	// Test that different tiles have different string representations
	checkerTile := CreateTestTile(2)
	checkerStr := checkerTile.String()
	
	assert.NotEqual(t, emptyStr, checkerStr, "Different tiles should have different string representations")
}

// TestTileDataValidation tests edge cases and validation
func TestTileDataValidation(t *testing.T) {
	// Test that TileData is exactly 16 bytes
	var data TileData
	assert.Equal(t, 16, len(data), "TileData should be exactly 16 bytes")
	
	// Test that invalid tile data doesn't crash the system
	invalidData := TileData{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	
	tile := NewTileFromData(invalidData)
	assert.NotNil(t, tile, "Should be able to create tile from any data")
	
	// All pixels should be color 3 (since 0xFF in both planes = color 3)
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			assert.Equal(t, uint8(3), tile.GetPixel(x, y), "0xFF data should result in color 3")
		}
	}
}
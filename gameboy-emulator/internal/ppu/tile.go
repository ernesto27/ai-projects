// Package ppu - Game Boy tile system implementation
// Handles 8x8 pixel tiles, the fundamental building blocks of Game Boy graphics

package ppu

// Game Boy tile system constants
const (
	// Tile dimensions
	TileWidth  = 8  // Each tile is 8 pixels wide
	TileHeight = 8  // Each tile is 8 pixels tall
	TileSize   = 16 // Each tile is 16 bytes in Game Boy 2bpp format
	
	// Tile limits
	MaxTileIndex = 255 // Maximum tile index (0-255)
	MaxTiles     = 256 // Maximum number of tiles in a pattern table
	
	// VRAM memory layout constants
	TilePatternTable0Start = 0x8000 // Start of tile pattern table 0 ($8000 method)
	TilePatternTable0End   = 0x8FFF // End of tile pattern table 0
	TilePatternTable1Start = 0x8800 // Start of tile pattern table 1 ($8800 method) 
	TilePatternTable1End   = 0x97FF // End of tile pattern table 1
	
	// Tile map constants
	BackgroundMap0Start = 0x9800 // Background tile map 0 (32x32 grid)
	BackgroundMap0End   = 0x9BFF // End of background map 0
	BackgroundMap1Start = 0x9C00 // Background tile map 1 (32x32 grid)
	BackgroundMap1End   = 0x9FFF // End of background map 1
	
	// Tile map dimensions
	TileMapWidth  = 32 // Tile maps are 32 tiles wide
	TileMapHeight = 32 // Tile maps are 32 tiles tall
	TileMapSize   = TileMapWidth * TileMapHeight // 1024 tile indices per map
	
	// Screen dimensions in tiles
	ScreenTilesWidth  = 20 // 160 pixels / 8 = 20 tiles across
	ScreenTilesHeight = 18 // 144 pixels / 8 = 18 tiles down
)

// Tile represents a decoded 8x8 pixel Game Boy tile
// Each pixel can have a value from 0-3 representing one of 4 colors
type Tile struct {
	// Pixels stores the 8x8 grid of pixel values
	// [row][column] format where each value is 0-3
	Pixels [TileHeight][TileWidth]uint8
}

// TileData represents raw Game Boy tile data in 2bpp format
// Each tile is stored as 16 bytes in VRAM using 2 bits per pixel
type TileData [TileSize]uint8

// NewTile creates a new tile with all pixels set to 0 (transparent/white)
func NewTile() *Tile {
	return &Tile{
		Pixels: [TileHeight][TileWidth]uint8{}, // All pixels default to 0
	}
}

// NewTileFromData creates a tile by decoding Game Boy 2bpp data
func NewTileFromData(data TileData) *Tile {
	tile := NewTile()
	tile.LoadFromData(data)
	return tile
}

// GetPixel returns the color value (0-3) at the specified position
// Returns 0 if coordinates are out of bounds
func (t *Tile) GetPixel(x, y int) uint8 {
	if x < 0 || x >= TileWidth || y < 0 || y >= TileHeight {
		return 0
	}
	return t.Pixels[y][x]
}

// SetPixel sets the color value (0-3) at the specified position
// Does nothing if coordinates are out of bounds or color > 3
func (t *Tile) SetPixel(x, y int, color uint8) {
	if x < 0 || x >= TileWidth || y < 0 || y >= TileHeight {
		return
	}
	if color > 3 {
		color = 3 // Clamp to valid color range
	}
	t.Pixels[y][x] = color
}

// Clear sets all pixels in the tile to the specified color (0-3)
func (t *Tile) Clear(color uint8) {
	if color > 3 {
		color = 3 // Clamp to valid color range
	}
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			t.Pixels[y][x] = color
		}
	}
}

// IsEmpty returns true if all pixels are 0 (transparent/white)
func (t *Tile) IsEmpty() bool {
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			if t.Pixels[y][x] != 0 {
				return false
			}
		}
	}
	return true
}

// Copy creates a deep copy of the tile
func (t *Tile) Copy() *Tile {
	newTile := NewTile()
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			newTile.Pixels[y][x] = t.Pixels[y][x]
		}
	}
	return newTile
}

// LoadFromData decodes Game Boy 2bpp tile data into the tile
// Game Boy format: 16 bytes total, 2 bytes per row (8 rows)
// Each pixel = bit from plane 0 + (bit from plane 1 << 1)
func (t *Tile) LoadFromData(data TileData) {
	for row := 0; row < TileHeight; row++ {
		// Each row uses 2 bytes: plane 0 and plane 1
		plane0 := data[row*2]     // Even bytes = bit plane 0
		plane1 := data[row*2+1]   // Odd bytes = bit plane 1
		
		// Decode each pixel in the row (8 pixels)
		for col := 0; col < TileWidth; col++ {
			// Extract bit for this pixel from each plane
			// Bits are stored MSB first (bit 7 = leftmost pixel)
			bit0 := (plane0 >> (7 - col)) & 1
			bit1 := (plane1 >> (7 - col)) & 1
			
			// Combine bits to get final color (0-3)
			t.Pixels[row][col] = bit0 + (bit1 << 1)
		}
	}
}

// ToData encodes the tile into Game Boy 2bpp format
// Converts 8x8 pixel data back to 16 bytes for storage in VRAM
func (t *Tile) ToData() TileData {
	var data TileData
	
	for row := 0; row < TileHeight; row++ {
		var plane0, plane1 uint8
		
		// Encode each pixel in the row
		for col := 0; col < TileWidth; col++ {
			color := t.Pixels[row][col]
			
			// Extract bits from color value
			bit0 := color & 1         // Low bit
			bit1 := (color >> 1) & 1  // High bit
			
			// Set bits in planes (MSB first)
			if bit0 != 0 {
				plane0 |= 1 << (7 - col)
			}
			if bit1 != 0 {
				plane1 |= 1 << (7 - col)
			}
		}
		
		// Store the two planes for this row
		data[row*2] = plane0     // Even byte = plane 0
		data[row*2+1] = plane1   // Odd byte = plane 1
	}
	
	return data
}

// FlipHorizontal creates a horizontally flipped copy of the tile
// Used for sprite flipping effects
func (t *Tile) FlipHorizontal() *Tile {
	flipped := NewTile()
	for row := 0; row < TileHeight; row++ {
		for col := 0; col < TileWidth; col++ {
			// Mirror horizontally: column 0 becomes 7, 1 becomes 6, etc.
			flipped.Pixels[row][col] = t.Pixels[row][TileWidth-1-col]
		}
	}
	return flipped
}

// FlipVertical creates a vertically flipped copy of the tile
// Used for sprite flipping effects
func (t *Tile) FlipVertical() *Tile {
	flipped := NewTile()
	for row := 0; row < TileHeight; row++ {
		for col := 0; col < TileWidth; col++ {
			// Mirror vertically: row 0 becomes 7, 1 becomes 6, etc.
			flipped.Pixels[row][col] = t.Pixels[TileHeight-1-row][col]
		}
	}
	return flipped
}

// FlipBoth creates a copy flipped both horizontally and vertically
// Equivalent to 180-degree rotation
func (t *Tile) FlipBoth() *Tile {
	flipped := NewTile()
	for row := 0; row < TileHeight; row++ {
		for col := 0; col < TileWidth; col++ {
			// Mirror both axes
			flipped.Pixels[row][col] = t.Pixels[TileHeight-1-row][TileWidth-1-col]
		}
	}
	return flipped
}

// String returns a human-readable representation of the tile
// Uses characters to represent different color values
func (t *Tile) String() string {
	colorChars := []rune{' ', '░', '▒', '█'} // 0=space, 1=light, 2=medium, 3=dark
	
	result := "Tile 8x8:\n"
	for row := 0; row < TileHeight; row++ {
		for col := 0; col < TileWidth; col++ {
			color := t.Pixels[row][col]
			if color > 3 {
				color = 3 // Safety clamp
			}
			result += string(colorChars[color])
		}
		result += "\n"
	}
	return result
}

// =============================================================================
// Tile Address Calculation Functions
// =============================================================================

// GetTileAddress calculates the VRAM address for a tile index
// Game Boy has two addressing modes: $8000 (unsigned) and $8800 (signed)
func GetTileAddress(index uint8, useSignedMode bool) uint16 {
	if useSignedMode {
		// $8800 method: treat index as signed byte (-128 to +127)
		// Tile 0 is at $9000, negative indices go downward, positive go upward
		signedIndex := int8(index)
		return 0x9000 + uint16(signedIndex)*TileSize
	} else {
		// $8000 method: treat index as unsigned byte (0-255)
		// Tile 0 is at $8000, indices increase upward
		return TilePatternTable0Start + uint16(index)*TileSize
	}
}

// GetTileIndexFromAddress calculates the tile index from a VRAM address
// Returns the tile index and whether it uses signed mode
func GetTileIndexFromAddress(address uint16) (uint8, bool) {
	// Check if address is in $8000 table range
	if address >= TilePatternTable0Start && address <= TilePatternTable0End {
		// $8000 method (unsigned)
		tileIndex := (address - TilePatternTable0Start) / TileSize
		return uint8(tileIndex), false
	}
	
	// Check if address is in $8800 table range
	if address >= TilePatternTable1Start && address <= TilePatternTable1End {
		// $8800 method (signed)
		offset := int16(address - 0x9000) // Offset from middle of table
		tileIndex := offset / TileSize
		return uint8(tileIndex), true
	}
	
	// Invalid address, return 0
	return 0, false
}

// IsValidTileAddress checks if an address is within VRAM tile pattern area
func IsValidTileAddress(address uint16) bool {
	return (address >= TilePatternTable0Start && address <= TilePatternTable0End) ||
		   (address >= TilePatternTable1Start && address <= TilePatternTable1End)
}

// =============================================================================
// Tile Map Functions
// =============================================================================

// GetTileMapAddress calculates the address of a tile map entry
// mapSelect: false = map 0 ($9800), true = map 1 ($9C00)
func GetTileMapAddress(x, y int, mapSelect bool) uint16 {
	if x < 0 || x >= TileMapWidth || y < 0 || y >= TileMapHeight {
		return 0 // Invalid coordinates
	}
	
	offset := uint16(y*TileMapWidth + x)
	
	if mapSelect {
		return BackgroundMap1Start + offset // Map 1 at $9C00
	} else {
		return BackgroundMap0Start + offset // Map 0 at $9800
	}
}

// IsValidTileMapAddress checks if an address is within tile map area
func IsValidTileMapAddress(address uint16) bool {
	return (address >= BackgroundMap0Start && address <= BackgroundMap0End) ||
		   (address >= BackgroundMap1Start && address <= BackgroundMap1End)
}

// =============================================================================
// Debugging and Utility Functions
// =============================================================================

// CreateTestTile creates a tile with a specific pattern for testing
func CreateTestTile(pattern uint8) *Tile {
	tile := NewTile()
	
	switch pattern {
	case 0: // Solid color 0 (white/transparent)
		tile.Clear(0)
	case 1: // Solid color 3 (black)
		tile.Clear(3)
	case 2: // Checkerboard pattern
		for y := 0; y < TileHeight; y++ {
			for x := 0; x < TileWidth; x++ {
				if (x+y)%2 == 0 {
					tile.Pixels[y][x] = 0
				} else {
					tile.Pixels[y][x] = 3
				}
			}
		}
	case 3: // Gradient pattern
		for y := 0; y < TileHeight; y++ {
			for x := 0; x < TileWidth; x++ {
				tile.Pixels[y][x] = uint8((x + y) % 4)
			}
		}
	case 4: // Border pattern
		for y := 0; y < TileHeight; y++ {
			for x := 0; x < TileWidth; x++ {
				if x == 0 || x == TileWidth-1 || y == 0 || y == TileHeight-1 {
					tile.Pixels[y][x] = 3 // Border
				} else {
					tile.Pixels[y][x] = 0 // Interior
				}
			}
		}
	default:
		tile.Clear(0)
	}
	
	return tile
}

// AnalyzeTile returns statistics about a tile's content
func AnalyzeTile(tile *Tile) map[string]interface{} {
	colorCounts := [4]int{}
	
	// Count pixels of each color
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			color := tile.Pixels[y][x]
			if color <= 3 {
				colorCounts[color]++
			}
		}
	}
	
	return map[string]interface{}{
		"isEmpty":     tile.IsEmpty(),
		"color0Count": colorCounts[0], // White/transparent
		"color1Count": colorCounts[1], // Light gray
		"color2Count": colorCounts[2], // Dark gray  
		"color3Count": colorCounts[3], // Black
		"totalPixels": TileWidth * TileHeight,
	}
}

// CompareTiles returns true if two tiles have identical pixel data
func CompareTiles(tile1, tile2 *Tile) bool {
	for y := 0; y < TileHeight; y++ {
		for x := 0; x < TileWidth; x++ {
			if tile1.Pixels[y][x] != tile2.Pixels[y][x] {
				return false
			}
		}
	}
	return true
}

// CreateTileFromPattern creates a tile from a string pattern for testing
// Each character represents a pixel: ' '=0, '.'=1, 'o'=2, '#'=3
func CreateTileFromPattern(pattern string) *Tile {
	tile := NewTile()
	
	// Remove newlines and convert to single string
	cleanPattern := ""
	for _, char := range pattern {
		if char != '\n' && char != '\r' {
			cleanPattern += string(char)
		}
	}
	
	// Fill tile from pattern (up to 64 characters)
	for i, char := range cleanPattern {
		if i >= TileWidth*TileHeight {
			break
		}
		
		x := i % TileWidth
		y := i / TileWidth
		
		var color uint8
		switch char {
		case ' ':
			color = 0
		case '.':
			color = 1
		case 'o':
			color = 2
		case '#':
			color = 3
		default:
			color = 0
		}
		
		tile.SetPixel(x, y, color)
	}
	
	return tile
}
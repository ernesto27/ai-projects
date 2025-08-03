// Package ppu - VRAM (Video RAM) organization and management
// Handles tile pattern tables and tile maps as organized in Game Boy memory

package ppu

import "fmt"

// VRAM represents the Game Boy's Video RAM organization
// Total size: 8KB (0x8000-0x9FFF) containing tile patterns and tile maps
type VRAM struct {
	// Raw VRAM data (8KB total)
	data [0x2000]uint8 // 0x8000-0x9FFF mapped to 0x0000-0x1FFF
	
	// Organized access to VRAM regions
	patternTable0 *TilePatternTable // $8000 method tiles (0x0000-0x0FFF in data)
	patternTable1 *TilePatternTable // $8800 method tiles (0x0800-0x17FF in data)
	backgroundMap0 *TileMap         // Background map 0 (0x1800-0x1BFF in data)
	backgroundMap1 *TileMap         // Background map 1 (0x1C00-0x1FFF in data)
}

// TilePatternTable manages a collection of 256 tiles
// Each tile is 16 bytes, so table is 4KB total
type TilePatternTable struct {
	vram     *VRAM  // Reference to parent VRAM
	baseAddr uint16 // Base address in VRAM space (0x8000 or 0x8800)
}

// TileMap represents a 32×32 grid of tile indices
// Each entry is 1 byte (tile index), so map is 1KB total
type TileMap struct {
	vram     *VRAM  // Reference to parent VRAM
	baseAddr uint16 // Base address in VRAM space (0x9800 or 0x9C00)
}

// NewVRAM creates a new VRAM instance with organized access structures
func NewVRAM() *VRAM {
	vram := &VRAM{
		data: [0x2000]uint8{}, // Initialize all VRAM to 0
	}
	
	// Create organized access structures
	vram.patternTable0 = &TilePatternTable{
		vram:     vram,
		baseAddr: 0x8000,
	}
	
	vram.patternTable1 = &TilePatternTable{
		vram:     vram,
		baseAddr: 0x8800,
	}
	
	vram.backgroundMap0 = &TileMap{
		vram:     vram,
		baseAddr: 0x9800,
	}
	
	vram.backgroundMap1 = &TileMap{
		vram:     vram,
		baseAddr: 0x9C00,
	}
	
	return vram
}

// =============================================================================
// VRAM Raw Access Methods
// =============================================================================

// ReadByte reads a byte from VRAM at the specified Game Boy address
func (v *VRAM) ReadByte(address uint16) uint8 {
	if !v.IsValidAddress(address) {
		return 0xFF // Invalid address returns 0xFF
	}
	
	// Convert Game Boy address (0x8000-0x9FFF) to array index (0x0000-0x1FFF)
	index := address - 0x8000
	return v.data[index]
}

// WriteByte writes a byte to VRAM at the specified Game Boy address
func (v *VRAM) WriteByte(address uint16, value uint8) {
	if !v.IsValidAddress(address) {
		return // Ignore writes to invalid addresses
	}
	
	// Convert Game Boy address (0x8000-0x9FFF) to array index (0x0000-0x1FFF)
	index := address - 0x8000
	v.data[index] = value
}

// ReadWord reads a 16-bit word from VRAM (little-endian)
func (v *VRAM) ReadWord(address uint16) uint16 {
	low := uint16(v.ReadByte(address))
	high := uint16(v.ReadByte(address + 1))
	return (high << 8) | low
}

// WriteWord writes a 16-bit word to VRAM (little-endian)
func (v *VRAM) WriteWord(address uint16, value uint16) {
	v.WriteByte(address, uint8(value&0xFF))         // Low byte
	v.WriteByte(address+1, uint8((value>>8)&0xFF)) // High byte
}

// IsValidAddress checks if an address is within VRAM range
func (v *VRAM) IsValidAddress(address uint16) bool {
	return address >= 0x8000 && address <= 0x9FFF
}

// Clear fills all VRAM with the specified value
func (v *VRAM) Clear(value uint8) {
	for i := range v.data {
		v.data[i] = value
	}
}

// =============================================================================
// Tile Pattern Table Methods
// =============================================================================

// GetPatternTable0 returns the tile pattern table using $8000 addressing
func (v *VRAM) GetPatternTable0() *TilePatternTable {
	return v.patternTable0
}

// GetPatternTable1 returns the tile pattern table using $8800 addressing
func (v *VRAM) GetPatternTable1() *TilePatternTable {
	return v.patternTable1
}

// GetTile reads a tile from the pattern table using the specified addressing mode
func (t *TilePatternTable) GetTile(index uint8) *Tile {
	address := GetTileAddress(index, t.baseAddr == 0x8800)
	
	// Read 16 bytes of tile data
	var tileData TileData
	for i := 0; i < TileSize; i++ {
		tileData[i] = t.vram.ReadByte(address + uint16(i))
	}
	
	// Decode into tile structure
	return NewTileFromData(tileData)
}

// SetTile writes a tile to the pattern table using the specified addressing mode
func (t *TilePatternTable) SetTile(index uint8, tile *Tile) {
	address := GetTileAddress(index, t.baseAddr == 0x8800)
	
	// Encode tile to Game Boy format
	tileData := tile.ToData()
	
	// Write 16 bytes of tile data
	for i := 0; i < TileSize; i++ {
		t.vram.WriteByte(address+uint16(i), tileData[i])
	}
}

// GetTileData reads raw tile data (16 bytes) from the pattern table
func (t *TilePatternTable) GetTileData(index uint8) TileData {
	address := GetTileAddress(index, t.baseAddr == 0x8800)
	
	var tileData TileData
	for i := 0; i < TileSize; i++ {
		tileData[i] = t.vram.ReadByte(address + uint16(i))
	}
	
	return tileData
}

// SetTileData writes raw tile data (16 bytes) to the pattern table
func (t *TilePatternTable) SetTileData(index uint8, data TileData) {
	address := GetTileAddress(index, t.baseAddr == 0x8800)
	
	for i := 0; i < TileSize; i++ {
		t.vram.WriteByte(address+uint16(i), data[i])
	}
}

// LoadTiles loads multiple tiles from raw data array
func (t *TilePatternTable) LoadTiles(startIndex uint8, tilesData []TileData) {
	for i, tileData := range tilesData {
		if int(startIndex)+i > MaxTileIndex {
			break // Don't exceed tile table bounds
		}
		t.SetTileData(startIndex+uint8(i), tileData)
	}
}

// =============================================================================
// Tile Map Methods
// =============================================================================

// GetBackgroundMap0 returns background tile map 0 (0x9800)
func (v *VRAM) GetBackgroundMap0() *TileMap {
	return v.backgroundMap0
}

// GetBackgroundMap1 returns background tile map 1 (0x9C00)
func (v *VRAM) GetBackgroundMap1() *TileMap {
	return v.backgroundMap1
}

// GetTileIndex reads a tile index from the map at the specified coordinates
func (tm *TileMap) GetTileIndex(x, y int) uint8 {
	if x < 0 || x >= TileMapWidth || y < 0 || y >= TileMapHeight {
		return 0 // Return 0 for out-of-bounds coordinates
	}
	
	address := GetTileMapAddress(x, y, tm.baseAddr == 0x9C00)
	return tm.vram.ReadByte(address)
}

// SetTileIndex writes a tile index to the map at the specified coordinates
func (tm *TileMap) SetTileIndex(x, y int, index uint8) {
	if x < 0 || x >= TileMapWidth || y < 0 || y >= TileMapHeight {
		return // Ignore out-of-bounds writes
	}
	
	address := GetTileMapAddress(x, y, tm.baseAddr == 0x9C00)
	tm.vram.WriteByte(address, index)
}

// GetTileIndexLinear reads a tile index using linear addressing (0-1023)
func (tm *TileMap) GetTileIndexLinear(index int) uint8 {
	if index < 0 || index >= TileMapSize {
		return 0
	}
	
	x := index % TileMapWidth
	y := index / TileMapWidth
	return tm.GetTileIndex(x, y)
}

// SetTileIndexLinear writes a tile index using linear addressing (0-1023)
func (tm *TileMap) SetTileIndexLinear(index int, tileIndex uint8) {
	if index < 0 || index >= TileMapSize {
		return
	}
	
	x := index % TileMapWidth
	y := index / TileMapWidth
	tm.SetTileIndex(x, y, tileIndex)
}

// FillMap fills the entire tile map with the specified tile index
func (tm *TileMap) FillMap(tileIndex uint8) {
	for y := 0; y < TileMapHeight; y++ {
		for x := 0; x < TileMapWidth; x++ {
			tm.SetTileIndex(x, y, tileIndex)
		}
	}
}

// LoadMapData loads tile map data from a byte array
func (tm *TileMap) LoadMapData(data []uint8) {
	maxLength := TileMapSize
	if len(data) < maxLength {
		maxLength = len(data)
	}
	
	for i := 0; i < maxLength; i++ {
		tm.SetTileIndexLinear(i, data[i])
	}
}

// GetMapData exports tile map data as a byte array
func (tm *TileMap) GetMapData() []uint8 {
	data := make([]uint8, TileMapSize)
	
	for i := 0; i < TileMapSize; i++ {
		data[i] = tm.GetTileIndexLinear(i)
	}
	
	return data
}

// GetVisibleRegion returns the tile indices for the visible screen area
// scrollX, scrollY: background scroll position
// Returns 20×18 array of tile indices for the visible screen
func (tm *TileMap) GetVisibleRegion(scrollX, scrollY uint8) [ScreenTilesHeight][ScreenTilesWidth]uint8 {
	var visibleTiles [ScreenTilesHeight][ScreenTilesWidth]uint8
	
	// Calculate starting tile coordinates (wrap around 32×32 map)
	startTileX := int(scrollX) / TileWidth
	startTileY := int(scrollY) / TileHeight
	
	for screenY := 0; screenY < ScreenTilesHeight; screenY++ {
		for screenX := 0; screenX < ScreenTilesWidth; screenX++ {
			// Calculate map coordinates with wrapping
			mapX := (startTileX + screenX) % TileMapWidth
			mapY := (startTileY + screenY) % TileMapHeight
			
			visibleTiles[screenY][screenX] = tm.GetTileIndex(mapX, mapY)
		}
	}
	
	return visibleTiles
}

// =============================================================================
// High-Level VRAM Operations
// =============================================================================

// GetTileFromMap retrieves a decoded tile from a map coordinate
// mapSelect: false = map 0, true = map 1
// useSignedMode: tile addressing mode for pattern table
func (v *VRAM) GetTileFromMap(mapX, mapY int, mapSelect bool, useSignedMode bool) *Tile {
	// Get tile map
	var tileMap *TileMap
	if mapSelect {
		tileMap = v.backgroundMap1
	} else {
		tileMap = v.backgroundMap0
	}
	
	// Get tile index from map
	tileIndex := tileMap.GetTileIndex(mapX, mapY)
	
	// Get pattern table
	var patternTable *TilePatternTable
	if useSignedMode {
		patternTable = v.patternTable1
	} else {
		patternTable = v.patternTable0
	}
	
	// Get decoded tile
	return patternTable.GetTile(tileIndex)
}

// RenderTileToFramebuffer renders a tile to the PPU framebuffer
// This is a helper method for future background/sprite rendering
func (v *VRAM) RenderTileToFramebuffer(framebuffer *[ScreenHeight][ScreenWidth]uint8, 
									   tile *Tile, screenX, screenY int, palette [4]uint8) {
	for tileY := 0; tileY < TileHeight; tileY++ {
		for tileX := 0; tileX < TileWidth; tileX++ {
			// Calculate screen coordinates
			pixelX := screenX + tileX
			pixelY := screenY + tileY
			
			// Check bounds
			if pixelX < 0 || pixelX >= ScreenWidth || pixelY < 0 || pixelY >= ScreenHeight {
				continue
			}
			
			// Get pixel color and apply palette
			rawColor := tile.GetPixel(tileX, tileY)
			finalColor := palette[rawColor]
			
			// Write to framebuffer
			framebuffer[pixelY][pixelX] = finalColor
		}
	}
}

// =============================================================================
// Debugging and Analysis Functions
// =============================================================================

// GetVRAMStats returns statistics about VRAM usage
func (v *VRAM) GetVRAMStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Count non-zero bytes in each region
	var patternBytes, mapBytes int
	
	// Pattern table region (0x8000-0x97FF)
	for i := 0x0000; i < 0x1800; i++ {
		if v.data[i] != 0 {
			patternBytes++
		}
	}
	
	// Tile map region (0x9800-0x9FFF)  
	for i := 0x1800; i < 0x2000; i++ {
		if v.data[i] != 0 {
			mapBytes++
		}
	}
	
	stats["totalSize"] = len(v.data)
	stats["patternDataUsed"] = patternBytes
	stats["mapDataUsed"] = mapBytes
	stats["totalUsed"] = patternBytes + mapBytes
	stats["percentUsed"] = float64(patternBytes+mapBytes) / float64(len(v.data)) * 100.0
	
	return stats
}

// DumpTileMap returns a string representation of a tile map for debugging
func (tm *TileMap) DumpTileMap(rows, cols int) string {
	if rows > TileMapHeight {
		rows = TileMapHeight
	}
	if cols > TileMapWidth {
		cols = TileMapWidth
	}
	
	result := fmt.Sprintf("Tile Map (showing %dx%d):\n", cols, rows)
	
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			tileIndex := tm.GetTileIndex(x, y)
			result += fmt.Sprintf("%02X ", tileIndex)
		}
		result += "\n"
	}
	
	return result
}

// FindTileUsage finds all map locations that use a specific tile index
func (tm *TileMap) FindTileUsage(targetIndex uint8) []struct{ X, Y int } {
	var locations []struct{ X, Y int }
	
	for y := 0; y < TileMapHeight; y++ {
		for x := 0; x < TileMapWidth; x++ {
			if tm.GetTileIndex(x, y) == targetIndex {
				locations = append(locations, struct{ X, Y int }{X: x, Y: y})
			}
		}
	}
	
	return locations
}

// ValidateVRAM performs consistency checks on VRAM data
func (v *VRAM) ValidateVRAM() []string {
	var issues []string
	
	// Check for common issues or inconsistencies
	// This is mainly for debugging and development
	
	// Count unique tiles in use
	usedTiles := make(map[uint8]bool)
	for _, tileMap := range []*TileMap{v.backgroundMap0, v.backgroundMap1} {
		for y := 0; y < TileMapHeight; y++ {
			for x := 0; x < TileMapWidth; x++ {
				index := tileMap.GetTileIndex(x, y)
				usedTiles[index] = true
			}
		}
	}
	
	if len(usedTiles) == 0 {
		issues = append(issues, "No tiles are referenced by tile maps")
	}
	
	// Check for potentially invalid tile indices in signed mode
	for index := range usedTiles {
		if index > 127 {
			// In signed mode, indices 128-255 are negative (-128 to -1)
			// This isn't necessarily an error, but worth noting
		}
	}
	
	return issues
}

// =============================================================================
// VRAM Interface Implementation for PPU
// =============================================================================

// The VRAM struct can be used to implement the VRAMInterface needed by PPU

// ReadVRAM implements VRAMInterface.ReadVRAM
func (v *VRAM) ReadVRAM(address uint16) uint8 {
	return v.ReadByte(address)
}

// WriteVRAM implements VRAMInterface.WriteVRAM  
func (v *VRAM) WriteVRAM(address uint16, value uint8) {
	v.WriteByte(address, value)
}

// ReadOAM implements VRAMInterface.ReadOAM (OAM is separate from VRAM)
// This is a placeholder - OAM will be implemented separately
func (v *VRAM) ReadOAM(address uint16) uint8 {
	// OAM is at 0xFE00-0xFE9F, separate from VRAM
	// Return 0 for now - will be implemented when we add sprite support
	return 0
}

// WriteOAM implements VRAMInterface.WriteOAM (OAM is separate from VRAM)
// This is a placeholder - OAM will be implemented separately
func (v *VRAM) WriteOAM(address uint16, value uint8) {
	// OAM is at 0xFE00-0xFE9F, separate from VRAM
	// Do nothing for now - will be implemented when we add sprite support
}